package member

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	Username string
	Email    string
}
