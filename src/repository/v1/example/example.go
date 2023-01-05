package example

import (
	"github.com/forkyid/go-consumer-boilerplate/src/database"
	"gorm.io/gorm"
)

type Repository struct {
	master *gorm.DB
	slave  *gorm.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{
		master: db.Master,
		slave:  db.Slave,
	}
}

type Repositorier interface {
	UpdateEmail(memberID int, email string) (err error)
}

func (repo *Repository) UpdateEmail(memberID int, email string) (err error) {
	query := repo.master.Begin()
	query = query.Where("id", memberID)
	err = query.Update("email", email).Error

	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}
