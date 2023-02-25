package v1

import (
	"encoding/json"

	"github.com/forkyid/consumer-kyc-update/src/constant"
	httpCallback "github.com/forkyid/consumer-kyc-update/src/entity/v1/http/callback"
	httpKYC "github.com/forkyid/consumer-kyc-update/src/entity/v1/http/kyc"
	publisher "github.com/forkyid/consumer-kyc-update/src/publisher/v1"
	"github.com/forkyid/consumer-kyc-update/src/service/v1/additionalselfie"
	"github.com/forkyid/consumer-kyc-update/src/service/v1/kyc"
	"github.com/forkyid/go-utils/v1/aes"
	"github.com/forkyid/go-utils/v1/jwt"
	"github.com/forkyid/go-utils/v1/logger"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/streadway/amqp"
)

type Controller struct {
	kycSvc              kyc.Servicer
	additionalselfieSvc additionalselfie.Servicer
}

func NewController(
	kycSvc kyc.Servicer,
	additionalselfieSvc additionalselfie.Servicer,
) *Controller {
	return &Controller{
		kycSvc:              kycSvc,
		additionalselfieSvc: additionalselfieSvc,
	}
}

func (ctrl *Controller) Handler(request *Request, d *amqp.Delivery) {
	publisher.CallbackRoute.RoutingKey = d.CorrelationId
	publisher.CallbackRoute.QueueName = d.ReplyTo

	var auth string
	if val, ok := request.Headers["Authorization"]; ok && val.(string) != "" {
		auth = val.(string)
	} else {
		err := constant.ErrNoAuthorizationFound
		publishCallback(err)
		d.Ack(false)
		return
	}
	client, err := jwt.ExtractClient(auth)
	if err != nil {
		publishCallback(err)
		d.Ack(false)
		return
	}
	body := httpKYC.UpdateKYCRequest{}
	err = json.Unmarshal(request.Body, &body)
	if err != nil {
		publishCallback(err)
		d.Ack(false)
		return
	}
	err = validation.Validator.Struct(&body)
	if err != nil {
		publishCallback(err)
		d.Ack(false)
		return
	}
	if body.KYCID = aes.DecryptCMS(body.EncKYCID); body.KYCID <= 0 {
		err = constant.ErrInvalidID
		publishCallback(err)
		d.Ack(false)
		return
	}
	if body.SimilarAccountID = aes.DecryptCMS(body.EncSimilarAccountID); (body.KYCType == constant.KYCTypeUserSimilar || body.KYCType == constant.KYCTypeRetakeSimilar) && body.SimilarAccountID <= 0 {
		err = constant.ErrInvalidID
		publishCallback(err)
		d.Ack(false)
		return
	}
	switch body.KYCType {
	case constant.KYCTypeUser, constant.KYCTypeCreator:
		err = ctrl.kycSvc.UpdateKYC(body.KYCID, body.KYCType, body.Status, body.Reason, client.Email)

	case constant.KYCTypeRetake:
		err = ctrl.additionalselfieSvc.UpdateRetake(body.KYCID, body.Status, body.Reason, client.Email)

	case constant.KYCTypeRetakeSimilar:
		err = ctrl.additionalselfieSvc.UpdateRetakeSimilar(body.KYCID, body.SimilarAccountID, body.Status, client.Email)

	case constant.KYCTypeUserSimilar:
		err = ctrl.kycSvc.UpdateKYCUserSimilar(body.KYCID, body.SimilarAccountID, body.Status, client.Email)
	}

	publishCallback(err)
	d.Ack(false)
}

func publishCallback(err error) {
	payload := []byte("{}")
	if err != nil {
		logger.Errorf(nil, "", err)
		payload, _ = json.Marshal(httpCallback.CallbackRequest{
			ErrorMessage: err.Error(),
		})
	}
	publisher.CallbackRoute.Publish(&publisher.Publish{
		Body: string(payload),
	}, constant.RMQMessageBasePriority)
}
