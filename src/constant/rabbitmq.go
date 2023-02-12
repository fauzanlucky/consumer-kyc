package constant

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	ExchangeError          = "Failed to declare an exchange"
	RMQMessageBasePriority = 10
)

var (
	_ = godotenv.Load()

	KYCExchangeName = os.Getenv("RMQ_BE_KYC_EXCHANGE")
	KYCExchangeType = os.Getenv("RMQ_BE_KYC_EXCHANGE_TYPE")
	KYCRoutingKey   = os.Getenv("RMQ_BE_KYC_EXCHANGE_TYPE")
)
