package additionalselfie

import (
	"encoding/json"

	"github.com/fauzanlucky/consumer-kyc/src/constant"
	dbAdditionalselfie "github.com/fauzanlucky/consumer-kyc/src/entity/v1/db/additionalselfie"
	"github.com/fauzanlucky/consumer-kyc/src/repository/v1/additionalselfie"
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
}

func (svc *Service) isRetakeUnprocessed(id int) (isUnprocessed bool, err error) {
	additionalselfie := dbAdditionalselfie.Additionalselfie{}
	additionalselfie, err = svc.repo.GetKYCRetakeByID(id)
	if err != nil {
		err = errors.Wrap(err, "get retake")
		return
	}
	additionalselfiePayload := dbAdditionalselfie.AdditionalselfiePayload{}
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
	var isUnprocessed bool
	isUnprocessed, err = svc.isRetakeUnprocessed(id)
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
