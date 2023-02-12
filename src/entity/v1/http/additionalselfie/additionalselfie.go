package additionalselfie

type UpdateKYCRetake struct {
	AdditionalSelfieID int `json:"-"`

	EncAdditionalSelfieID string `json:"additional_selfie_id" validate:"required"`
	Status                string `json:"status" validate:"required,oneof=approved blocked declined replace_original"`
	Reason                string `json:"reason" validate:"omitempty,max=255"`
}
