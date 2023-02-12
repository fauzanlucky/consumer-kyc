package kyc

type KYCRequest struct {
	KYCType string            `json:"type" validate:"required,oneof=creator retake user"`
	Data    map[string]string `json:"data" validate:"required"`
}

type UpdateKYC struct {
	KYCID int `json:"-"`

	EncKYCID string `json:"kyc_id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=approved blocked need_retake not_completed verified_check"`
	Reason   string `json:"reason" validate:"omitempty,max=255"`
}
