package v1

import "github.com/forkyid/consumer-kyc-update/src/constant"

var (
	CallbackRoute = Route{
		ExchangeName: constant.KYCExchangeName,
		ExchangeType: constant.KYCExchangeType,
	}
)
