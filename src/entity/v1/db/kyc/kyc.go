package kyc

import (
	"time"

	"gorm.io/gorm"
)

type KYC struct {
	gorm.Model
	Payload []byte `gorm:"column:payload;type:jsonb"`
}

type KYCCreatorPayload struct {
	Data        KYCCreatorData `json:"data"`
	Inferences  KYCInferences  `json:"inferences"`
	ProcessedAt string         `json:"processed_at"`
	ProcessedBy string         `json:"processed_by"`
}

type KYCUserPayload struct {
	Data        KYCUserData   `json:"data"`
	Inferences  KYCInferences `json:"inferences"`
	ProcessedAt string        `json:"processed_at"`
	ProcessedBy string        `json:"processed_by"`
}

type KYCSimilarAccount struct {
	VideoFrontPath          string `json:"video_front_path"`
	VideoLeftPath           string `json:"video_left_path"`
	VideoRightPath          string `json:"video_right_path"`
	SelfiePath              string `json:"selfie_path"`
	VideoFrontThumbnailPath string `json:"video_front_thumbnail_path"`
	VideoLeftThumbnailPath  string `json:"video_left_thumbnail_path"`
	VideoRightThumbnailPath string `json:"video_right_thumbnail_path"`
	MemberID                int    `json:"member_id"`
	AccountStatus           string `json:"account_status"`
}

type KYCData struct {
	Email                   string              `json:"email"`
	MemberID                int                 `json:"member_id"`
	AccountID               int                 `json:"account_id"`
	Nickname                string              `json:"nickname"`
	Username                string              `json:"username"`
	Status                  string              `json:"status"`
	Reason                  string              `json:"reason"`
	IsCreator               bool                `json:"is_creator"`
	VideoFrontPath          string              `json:"video_front_path"`
	VideoLeftPath           string              `json:"video_left_path"`
	VideoRightPath          string              `json:"video_right_path"`
	SelfiePath              string              `json:"selfie_path"`
	VideoFrontThumbnailPath string              `json:"video_front_thumbnail_path"`
	VideoLeftThumbnailPath  string              `json:"video_left_thumbnail_path"`
	VideoRightThumbnailPath string              `json:"video_right_thumbnail_path"`
	SimilarAccounts         []KYCCreatorPayload `json:"similar_accounts"`
	CreatedAt               time.Time           `json:"created_at"`
}

type KYCCreatorData struct {
	Name              string `json:"name"`
	CreatorID         int    `json:"creator_id"`
	DateOfBirth       string `json:"date_of_birth"`
	Address           string `json:"address"`
	PostalCode        string `json:"postal_code"`
	DocumentType      string `json:"document_type"`
	DocumentNumber    string `json:"document_number"`
	DocumentPath      string `json:"document_path"`
	IsCreator         bool   `json:"is_creator"`
	KYCUserSelfiePath string `json:"kyc_user_selfie_path"`
	KYCData
}

type KYCUserData struct {
	KYCData
}

type KYCInferences struct {
	AntiSpoofing   KYCInferenceDetail `json:"anti_spoofing"`
	FaceDetection  KYCInferenceDetail `json:"face_detection"`
	HeadPose       KYCHeadPoseDetail  `json:"head_pose"`
	FaceSimilarity float64            `json:"face_similarity"`
}

type KYCInferenceDetail struct {
	SelfiePercent        float64 `json:"selfie"`
	VideoPercent         float64 `json:"video"`
	KYCUserSelfiePercent float64 `json:"kyc_user_selfie"`
}

type KYCHeadPoseDetail struct {
	ForwardPercent float64 `json:"forward"`
	LeftPercent    float64 `json:"left"`
	RightPercent   float64 `json:"right"`
}

func (KYC) TableName() string {
	return "kyc"
}
