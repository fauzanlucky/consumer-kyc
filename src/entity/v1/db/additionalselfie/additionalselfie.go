package additionalselfie

import (
	"time"

	"gorm.io/gorm"
)

type AdditionalSelfie struct {
	gorm.Model
	Payload []byte `gorm:"column:payload;type:jsonb"`
}

type AdditionalSelfiePayload struct {
	Data        AdditionalSelfieData       `json:"data"`
	Inferences  AdditionalSelfieInferences `json:"inferences"`
	ProcessedAt string                     `json:"processed_at"`
	ProcessedBy string                     `json:"processed_by"`
}

type AdditionalSelfieData struct {
	Email                   string                    `json:"email"`
	Reason                  string                    `json:"reason"`
	Status                  string                    `json:"status"`
	Nickname                string                    `json:"nickname"`
	Username                string                    `json:"username"`
	MemberID                int                       `json:"member_id"`
	AccountID               int                       `json:"account_id"`
	SelfiePath              string                    `json:"selfie_path"`
	RetakeCount             int                       `json:"retake_count"`
	RetakeStatus            string                    `json:"retake_status"`
	VideoLeftPath           string                    `json:"video_left_path"`
	VideoFrontPath          string                    `json:"video_front_path"`
	VideoRightPath          string                    `json:"video_right_path"`
	AdditionalSelfiePath    string                    `json:"additional_selfie_path"`
	AdditionalSelfieFileID  int                       `json:"additional_selfie_file_id"`
	VideoLeftThumbnailPath  string                    `json:"video_left_thumbnail_path"`
	VideoFrontThumbnailPath string                    `json:"video_front_thumbnail_path"`
	VideoRightThumbnailPath string                    `json:"video_right_thumbnail_path"`
	SimilarAccounts         []AdditionalSelfiePayload `json:"similar_accounts"`
	CreatedAt               time.Time                 `json:"created_at"`
}

type AdditionalSelfieHeadPose struct {
	Left    float64 `json:"left"`
	Right   float64 `json:"right"`
	Forward float64 `json:"forward"`
}

type AdditionalSelfieAntiSpoofing struct {
	Video            float64 `json:"video"`
	Selfie           float64 `json:"selfie"`
	SpoofCount       int     `json:"spoof_count"`
	AdditionalSelfie float64 `json:"additional_selfie"`
}

type AdditionalSelfieFaceDetection struct {
	Video            float64 `json:"video"`
	Selfie           float64 `json:"selfie"`
	AdditionalSelfie float64 `json:"additional_selfie"`
}

type AdditionalSelfieInferences struct {
	HeadPose       AdditionalSelfieHeadPose      `json:"head_pose"`
	AntiSpoofing   AdditionalSelfieAntiSpoofing  `json:"anti_spoofing"`
	FaceDetection  AdditionalSelfieFaceDetection `json:"face_detection"`
	FaceSimilarity float64                       `json:"face_similarity"`
}

func (AdditionalSelfie) TableName() string {
	return "additional_selfies"
}
