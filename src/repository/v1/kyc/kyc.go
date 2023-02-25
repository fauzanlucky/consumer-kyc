package kyc

import (
	"fmt"
	"time"

	"github.com/forkyid/consumer-kyc-update/src/database"
	dbKYC "github.com/forkyid/consumer-kyc-update/src/entity/v1/db/kyc"
	"gorm.io/gorm"
)

type Repository struct {
	main    *gorm.DB
	replica *gorm.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{
		main:    db.CMSMain,
		replica: db.CMSReplica,
	}
}

type Repositorier interface {
	GetKYCByID(kycID int) (result dbKYC.KYC, err error)
	UpdateKYCPayload(kycID int, status, reason, processorEmail string) (err error)
	UpdateKYCSimilar(kycID, similarAccountIdx int, status, processorEmail string) (err error)
}

func (repo *Repository) GetKYCByID(kycID int) (result dbKYC.KYC, err error) {
	query := repo.main.Model(&result)
	err = query.Take(&result, kycID).Error
	return
}

func (repo *Repository) UpdateKYCPayload(kycID int, status, reason, processorEmail string) (err error) {
	query := repo.main.Begin()
	updateQueryString := fmt.Sprintf(`
		UPDATE kyc SET payload = 
		jsonb_set(payload || '{"processed_at":"%s","processed_by":"%s"}', '{data}', payload->'data' || '{"status":"%s","reason":"%s"}')
		WHERE id = %d
	`, time.Now().Format(time.RFC3339), processorEmail, status, reason, kycID)
	query = query.Exec(updateQueryString)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error
	return
}

func (repo *Repository) UpdateKYCSimilar(kycID, similarAccountIdx int, status, processorEmail string) (err error) {
	path := fmt.Sprintf(`{data, similar_accounts, %d}`, similarAccountIdx)
	newValue := fmt.Sprintf(`{"processed_by": "%s", "processed_at": "%s"}`, processorEmail, time.Now().Format(time.RFC3339))
	updatedProcessedPayload := fmt.Sprintf(`
	jsonb_set(payload, '%s', payload->'data'->'similar_accounts'->%d || ('%s')::jsonb)
	`, path, similarAccountIdx, newValue)
	updatedStatusPayload := fmt.Sprintf(`jsonb_set(
		%s, 
		'{data, similar_accounts, %d, data, status}', 
		'"%s"'::jsonb)`, updatedProcessedPayload, similarAccountIdx, status)
	updateQuery := fmt.Sprintf(`
	UPDATE kyc
	SET payload = %s
	WHERE id = %d
	`, updatedStatusPayload, kycID)
	query := repo.main.Begin()
	query = query.Exec(updateQuery)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error

	return
}
