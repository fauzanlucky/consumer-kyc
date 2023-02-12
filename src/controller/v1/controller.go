package v1

import (
	"encoding/json"

	"github.com/fauzanlucky/consumer-kyc/src/constant"
	httpAdditionalselfie "github.com/fauzanlucky/consumer-kyc/src/entity/v1/http/additionalselfie"
	httpKYC "github.com/fauzanlucky/consumer-kyc/src/entity/v1/http/kyc"
	publisher "github.com/fauzanlucky/consumer-kyc/src/publisher/v1"
	"github.com/fauzanlucky/consumer-kyc/src/service/v1/additionalselfie"
	"github.com/fauzanlucky/consumer-kyc/src/service/v1/kyc"
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
	publisher.CallbackRoute.RoutingKey = d.ReplyTo
	publisher.CallbackRoute.QueueName = d.ReplyTo

	var auth string
	if val, ok := request.Headers["Authorization"]; ok && val.(string) != "" {
		auth = val.(string)
	} else {
		logger.Errorf(nil, "no authorization found", constant.ErrNoAuthorizationFound)
		publishCallback(constant.ErrNoAuthorizationFound)
		d.Ack(false)
	}
	client, err := jwt.ExtractClient(auth)
	if err != nil {
		logger.Errorf(nil, "Extract client", err)
		publishCallback(err)
		d.Ack(false)
	}
	body := httpKYC.KYCRequest{}
	err = json.Unmarshal(request.Body, &body)
	if err != nil {
		logger.Errorf(nil, "Unmarshal body", err)
		publishCallback(err)
		d.Ack(false)
	}
	err = validation.Validator.Struct(&body)
	if err != nil {
		logger.Errorf(nil, "Unmarshal body", err)
		publishCallback(err)
		d.Ack(false)
	}

	tempData, _ := json.Marshal(body.Data)
	switch body.KYCType {
	case "user", "creator":
		data := httpKYC.UpdateKYC{}
		_ = json.Unmarshal(tempData, &data)
		if data.KYCID = aes.DecryptCMS(data.EncKYCID); data.KYCID <= 0 {
			logger.Errorf(nil, "decrypt ID", err)
			err = constant.ErrInvalidID
			publishCallback(err)
			break
		}

		err = ctrl.kycSvc.UpdateKYC(data.KYCID, data.Status, data.Reason, client.Email)
		if err != nil {
			logger.Errorf(nil, "Update KYC", err)
			publishCallback(err)
			d.Ack(false)
			break
		}
	case "retake":
		data := httpAdditionalselfie.UpdateKYCRetake{}
		_ = json.Unmarshal(tempData, &data)
		if data.AdditionalSelfieID = aes.DecryptCMS(data.EncAdditionalSelfieID); data.AdditionalSelfieID <= 0 {
			err = constant.ErrInvalidID
			break
		}

		err = ctrl.additionalselfieSvc.UpdateRetake(data.AdditionalSelfieID, data.Status, data.Reason, client.Email)
		if err != nil {
			logger.Errorf(nil, "Update KYC", err)
			publishCallback(err)
			d.Ack(false)
			break
		}
	}

	publishCallback(err)
	d.Ack(false)
}

func publishCallback(err error) {
	data, _ := json.Marshal(err.Error())
	publisher.CallbackRoute.Publish(&publisher.Publish{
		Body: string(data),
	}, constant.RMQMessageBasePriority)
}
