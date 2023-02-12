package additionalselfie

import (
	"gorm.io/gorm"
)

type Additionalselfie struct {
	gorm.Model
	Payload []byte `gorm:"column:payload;type:jsonb"`
}

// Payload and data struct does not enlist all fields
// the following structs should be used for read-only
type AdditionalselfiePayload struct {
	ProcessedAt string `json:"processed_at"`
	ProcessedBy string `json:"processed_by"`
}

func (Additionalselfie) TableName() string {
	return "additional_selfies"
}
