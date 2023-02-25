package v1

import (
	"sync"

	"github.com/forkyid/consumer-kyc-update/src/constant"
	controller "github.com/forkyid/consumer-kyc-update/src/controller/v1"
	"github.com/forkyid/consumer-kyc-update/src/database"
	additionalselfieRepository "github.com/forkyid/consumer-kyc-update/src/repository/v1/additionalselfie"
	kycRepository "github.com/forkyid/consumer-kyc-update/src/repository/v1/kyc"
	"github.com/forkyid/consumer-kyc-update/src/service/v1/additionalselfie"
	"github.com/forkyid/consumer-kyc-update/src/service/v1/kyc"
	"gorm.io/gorm"
)

var cmsMain *gorm.DB
var cmsReplica *gorm.DB
var giftshopMain *gorm.DB
var giftshopReplica *gorm.DB

func initCMSMain(wg *sync.WaitGroup) {
	defer wg.Done()
	cmsMain = database.DBCMSMain()
}
func initCMSReplica(wg *sync.WaitGroup) {
	defer wg.Done()
	cmsReplica = database.DBCMSReplica()
}
func initGiftshopMain(wg *sync.WaitGroup) {
	defer wg.Done()
	giftshopMain = database.DBGiftshopMain()
}
func initGiftshopReplica(wg *sync.WaitGroup) {
	defer wg.Done()
	giftshopReplica = database.DBGiftshopReplica()
}

func getRoute() controller.Route {
	wg := &sync.WaitGroup{}

	// note: delete database-related file if not used (repository, service layers)
	wg.Add(4)
	go initCMSMain(wg)
	go initCMSReplica(wg)
	go initGiftshopMain(wg)
	go initGiftshopReplica(wg)
	wg.Wait()

	db := database.DB{
		CMSMain:         cmsMain,
		CMSReplica:      cmsReplica,
		GiftshopMain:    giftshopMain,
		GiftshopReplica: giftshopReplica,
	}

	kycRepo := kycRepository.NewRepository(db)
	additionalselfieRepo := additionalselfieRepository.NewRepository(db)

	return controller.Route{
		ExchangeName: constant.KYCExchangeName,
		ExchangeType: constant.KYCExchangeType,
		QueueName:    constant.KYCQueueName,
		RoutingKey:   constant.KYCRoutingKey,
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
