package additionalselfie

import (
	"fmt"
	"time"

	"github.com/forkyid/consumer-kyc-update/src/database"
	dbAdditionalselfie "github.com/forkyid/consumer-kyc-update/src/entity/v1/db/additionalselfie"
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
	GetKYCRetakeByID(id int) (result dbAdditionalselfie.AdditionalSelfie, err error)
	UpdateRetakePayload(kycID int, status, reason, processorEmail string) (err error)
	UpdateRetakeSimilarPayload(kycID, similarAccountIdx int, status, processorEmail string) (err error)
}

func (repo *Repository) GetKYCRetakeByID(id int) (result dbAdditionalselfie.AdditionalSelfie, err error) {
	query := repo.main.Model(&result)
	err = query.Take(&result, id).Error
	return
}

func (repo *Repository) UpdateRetakePayload(id int, status, reason, processorEmail string) (err error) {
	query := repo.main.Begin()
	query = query.Exec(`
		UPDATE additional_selfies SET payload = 
		jsonb_set(payload || '{"processed_at":"?","processed_by":"?"}', '{data}', payload->'data' || '{"status":"?","reason":"?"}')
		WHERE id = ?
	`, time.Now().Format(time.RFC3339), processorEmail, status, reason)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	err = query.Commit().Error
	return
}

func (repo *Repository) UpdateRetakeSimilarPayload(id, similarAccountIdx int, status, processorEmail string) (err error) {
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
	UPDATE additional_selfies
	SET payload = %s
	WHERE id = %d
	`, updatedStatusPayload, id)
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
