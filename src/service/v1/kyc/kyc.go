package kyc

import (
	"encoding/json"

	"github.com/forkyid/consumer-kyc-update/src/constant"
	dbKYC "github.com/forkyid/consumer-kyc-update/src/entity/v1/db/kyc"
	"github.com/forkyid/consumer-kyc-update/src/repository/v1/kyc"
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
	validateKYC(kycID int) (isUnprocessed, isCreator bool, err error)
	UpdateKYC(kycID int, kycType, status, reason, processorEmail string) (err error)
	UpdateKYCUserSimilar(kycID, similarAccountID int, status, processorEmail string) (err error)
}

func (svc *Service) validateKYC(kycID int) (isUnprocessed, isCreatorKYC bool, err error) {
	kyc := dbKYC.KYC{}
	kyc, err = svc.repo.GetKYCByID(kycID)
	if err != nil {
		err = errors.Wrap(err, "get KYC")
		return
	}
	kycPayload := dbKYC.KYCUserPayload{}
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
	isCreatorKYC = kycPayload.Data.IsCreator
	return
}

func (svc *Service) UpdateKYC(kycID int, kycType, status, reason, processorEmail string) (err error) {
	var isUnprocessed bool
	var isCreator bool
	isUnprocessed, isCreator, err = svc.validateKYC(kycID)
	if err != nil {
		err = errors.Wrap(err, "validate kyc unprocessed")
		return
	}
	if !isUnprocessed {
		err = constant.ErrKYCAlreadyHandled
		return
	}
	if (kycType == constant.KYCTypeUser) == isCreator {
		err = constant.ErrInvalidID
		return
	}
	if status == constant.KYCStatusVerifiedCheck {
		status = constant.KYCStatusApproved
	}

	err = svc.repo.UpdateKYCPayload(kycID, status, reason, processorEmail)
	if err != nil {
		err = errors.Wrap(err, "update KYC payload")
		return
	}

	return
}

func (svc *Service) UpdateKYCUserSimilar(kycID, similarAccountID int, status, processorEmail string) (err error) {
	kyc, err := svc.repo.GetKYCByID(kycID)
	if err != nil {
		err = errors.Wrap(err, "get kyc by id")
		return
	}
	kycPayload := dbKYC.KYCUserPayload{}
	err = json.Unmarshal(kyc.Payload, &kycPayload)
	if err != nil {
		err = errors.Wrap(err, "unmarshal payload")
		return
	}
	similarAccountIdx := -1
	for index, similarAccounts := range kycPayload.Data.SimilarAccounts {
		if similarAccounts.Data.AccountID == similarAccountID {
			similarAccountIdx = index
		}
	}
	if similarAccountIdx == -1 {
		err = constant.ErrSimilarIDNotFound
		return
	}
	err = svc.repo.UpdateKYCSimilar(kycID, similarAccountIdx, status, processorEmail)
	return
}
