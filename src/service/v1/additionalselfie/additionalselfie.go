package additionalselfie

import (
	"encoding/json"

	"github.com/forkyid/consumer-kyc-update/src/constant"
	dbAdditionalselfie "github.com/forkyid/consumer-kyc-update/src/entity/v1/db/additionalselfie"
	"github.com/forkyid/consumer-kyc-update/src/repository/v1/additionalselfie"
	"github.com/pkg/errors"
)

type Service struct {
	repo additionalselfie.Repositorier
}

func NewService(
	repo additionalselfie.Repositorier,
) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
	isRetakeUnprocessed(id int) (isUnprocessed bool, err error)
	UpdateRetake(id int, status, reason, processorEmail string) (err error)
	UpdateRetakeSimilar(id, similarAccountID int, status, processorEmail string) (err error)
}

func (svc *Service) isRetakeUnprocessed(id int) (isUnprocessed bool, err error) {
	additionalselfie := dbAdditionalselfie.AdditionalSelfie{}
	additionalselfie, err = svc.repo.GetKYCRetakeByID(id)
	if err != nil {
		err = errors.Wrap(err, "get retake")
		return
	}
	additionalselfiePayload := dbAdditionalselfie.AdditionalSelfiePayload{}
	err = json.Unmarshal(additionalselfie.Payload, &additionalselfiePayload)
	if err != nil {
		err = errors.Wrap(err, "unmarshal retake payload")
		return
	}
	if additionalselfiePayload.ProcessedAt == "" && additionalselfiePayload.ProcessedBy == "" {
		isUnprocessed = true
	}
	return
}

func (svc *Service) UpdateRetake(id int, status, reason, processorEmail string) (err error) {
	isUnprocessed, err := svc.isRetakeUnprocessed(id)
	if err != nil {
		err = errors.Wrap(err, "validate kyc unprocessed")
		return
	}
	if !isUnprocessed {
		err = constant.ErrKYCAlreadyHandled
		return
	}

	err = svc.repo.UpdateRetakePayload(id, status, reason, processorEmail)
	if err != nil {
		err = errors.Wrap(err, "update KYC payload")
		return
	}

	return
}

func (svc *Service) UpdateRetakeSimilar(id, accountID int, status, processorEmail string) (err error) {
	additionalselfie := dbAdditionalselfie.AdditionalSelfie{}
	additionalselfie, err = svc.repo.GetKYCRetakeByID(id)
	if err != nil {
		err = errors.Wrap(err, "get retake")
		return
	}
	additionalselfiePayload := dbAdditionalselfie.AdditionalSelfiePayload{}
	err = json.Unmarshal(additionalselfie.Payload, &additionalselfiePayload)
	if err != nil {
		err = errors.Wrap(err, "unmarshal retake payload")
		return
	}
	similarAccountIdx := -1
	for index, similarAccount := range additionalselfiePayload.Data.SimilarAccounts {
		if similarAccount.Data.AccountID == accountID {
			similarAccountIdx = index
		}
	}
	if similarAccountIdx == -1 {
		err = constant.ErrSimilarIDNotFound
		return
	}
	err = svc.repo.UpdateRetakeSimilarPayload(id, similarAccountIdx, status, processorEmail)
	return
}
