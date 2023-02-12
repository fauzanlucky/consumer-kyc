package kyc

import (
	"encoding/json"

	"github.com/fauzanlucky/consumer-kyc/src/constant"
	dbKYC "github.com/fauzanlucky/consumer-kyc/src/entity/v1/db/kyc"
	"github.com/fauzanlucky/consumer-kyc/src/repository/v1/kyc"
	"github.com/pkg/errors"
)

type Service struct {
	repo kyc.Repositorier
}

func NewService(
	repo kyc.Repositorier,
) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
	isKYCUnprocessed(kycID int) (isUnprocessed bool, err error)
	UpdateKYC(kycID int, status, reason, processorEmail string) (err error)
}

func (svc *Service) isKYCUnprocessed(kycID int) (isUnprocessed bool, err error) {
	kyc := dbKYC.KYC{}
	kyc, err = svc.repo.GetKYCByID(kycID)
	if err != nil {
		err = errors.Wrap(err, "get KYC")
		return
	}
	kycPayload := dbKYC.KYCPayload{}
	err = json.Unmarshal(kyc.Payload, &kycPayload)
	if err != nil {
		err = errors.Wrap(err, "unmarshal KYC payload")
		return
	}
	if kycPayload.ProcessedAt != "" || kycPayload.ProcessedBy != "" {
		isUnprocessed = false
		return
	} else {
		isUnprocessed = true
	}
	return
}

func (svc *Service) UpdateKYC(kycID int, status, reason, processorEmail string) (err error) {
	var isUnprocessed bool
	isUnprocessed, err = svc.isKYCUnprocessed(kycID)
	if err != nil {
		err = errors.Wrap(err, "validate kyc unprocessed")
		return
	}
	if !isUnprocessed {
		err = constant.ErrKYCAlreadyHandled
		return
	}

	err = svc.repo.UpdateKYCPayload(kycID, status, reason, processorEmail)
	if err != nil {
		err = errors.Wrap(err, "update KYC payload")
		return
	}

	return
}
