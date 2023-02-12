package additionalselfie

import (
	"time"

	"github.com/fauzanlucky/consumer-kyc/src/database"
	dbAdditionalselfie "github.com/fauzanlucky/consumer-kyc/src/entity/v1/db/additionalselfie"
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
	GetKYCRetakeByID(id int) (result dbAdditionalselfie.Additionalselfie, err error)
	UpdateRetakePayload(kycID int, status, reason, processorEmail string) (err error)
}

func (repo *Repository) GetKYCRetakeByID(id int) (result dbAdditionalselfie.Additionalselfie, err error) {
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
