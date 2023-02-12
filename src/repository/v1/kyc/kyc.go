package kyc

import (
	"time"

	"github.com/fauzanlucky/consumer-kyc/src/database"
	dbKYC "github.com/fauzanlucky/consumer-kyc/src/entity/v1/db/kyc"
	"gorm.io/gorm"
)

type Repository struct {
	main    *gorm.DB
	replica *gorm.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{
		main:    db.Main,
		replica: db.Replica,
	}
}

type Repositorier interface {
	GetKYCByID(kycID int) (result dbKYC.KYC, err error)
	UpdateKYCPayload(kycID int, status, reason, processorEmail string) (err error)
}

func (repo *Repository) GetKYCByID(kycID int) (result dbKYC.KYC, err error) {
	query := repo.main.Model(&result)
	err = query.Take(&result, kycID).Error
	return
}

func (repo *Repository) UpdateKYCPayload(kycID int, status, reason, processorEmail string) (err error) {
	query := repo.main.Begin()
	query = query.Exec(`
		UPDATE kyc SET payload = 
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
