package v1

import (
	"os"
	"sync"

	controller "github.com/fauzanlucky/consumer-kyc/src/controller/v1"
	"github.com/fauzanlucky/consumer-kyc/src/database"
	additionalselfieRepository "github.com/fauzanlucky/consumer-kyc/src/repository/v1/additionalselfie"
	kycRepository "github.com/fauzanlucky/consumer-kyc/src/repository/v1/kyc"
	"github.com/fauzanlucky/consumer-kyc/src/service/v1/additionalselfie"
	"github.com/fauzanlucky/consumer-kyc/src/service/v1/kyc"
	"gorm.io/gorm"
)

var main *gorm.DB
var replica *gorm.DB

func initMain(wg *sync.WaitGroup) {
	defer wg.Done()
	main = database.DBMain()
}

func initReplica(wg *sync.WaitGroup) {
	defer wg.Done()
	replica = database.DBReplica()
}

func getRoute() controller.Route {
	wg := &sync.WaitGroup{}

	// note: delete database-related file if not used (repository, service layers)
	wg.Add(2)
	go initMain(wg)
	go initReplica(wg)
	wg.Wait()

	db := database.DB{
		Main:    main,
		Replica: replica,
	}

	kycRepo := kycRepository.NewRepository(db)
	additionalselfieRepo := additionalselfieRepository.NewRepository(db)

	return controller.Route{
		ExchangeName: os.Getenv("RMQ_BE_KYC_EXCHANGE_NAME"),
		QueueName:    os.Getenv("RMQ_BE_KYC_QUEUE_NAME"),
		RoutingKey:   os.Getenv("RMQ_BE_KYC_ROUTE_NAME"),
		Handler: controller.NewController(
			kyc.NewService(kycRepo),
			additionalselfie.NewService(additionalselfieRepo),
		).Handler,
	}
}

func Start() {
	route := getRoute()
	route.Consume()
}
