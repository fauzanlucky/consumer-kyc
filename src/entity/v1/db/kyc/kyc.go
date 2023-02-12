package kyc

import (
	"gorm.io/gorm"
)

type KYC struct {
	gorm.Model
	Payload []byte `gorm:"column:payload;type:jsonb"`
}

// Payload and data struct does not enlist all fields
// the following structs should be used for read-only
type KYCPayload struct {
	ProcessedAt string `json:"processed_at"`
	ProcessedBy string `json:"processed_by"`
}

func (KYC) TableName() string {
	return "kyc"
}
