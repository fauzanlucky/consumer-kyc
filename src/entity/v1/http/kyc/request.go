package kyc

type UpdateKYCRequest struct {
	KYCType string `json:"type" validate:"required,oneof=creator retake retake-similar user user-similar"`

	KYCID            int `json:"-"`
	SimilarAccountID int `json:"-"`

	// Values below should be pre-validated by client
	EncKYCID            string `json:"kyc_id" validate:"required"`
	EncSimilarAccountID string `json:"similar_account_id"`
	Status              string `json:"status" validate:"required"`
	Reason              string `json:"reason" validate:"omitempty,max=255"`
}
