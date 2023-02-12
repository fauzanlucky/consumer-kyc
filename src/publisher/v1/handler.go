package v1

import "github.com/fauzanlucky/consumer-kyc/src/constant"

var (
	CallbackRoute = Route{
		ExchangeName: constant.KYCExchangeName,
		ExchangeType: constant.KYCExchangeType,
	}
)
