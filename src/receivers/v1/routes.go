package v1

import (
	"os"
	"sync"

	controller "github.com/forkyid/go-consumer-boilerplate/src/controller/v1"
	"github.com/forkyid/go-consumer-boilerplate/src/database"
	examplerepository "github.com/forkyid/go-consumer-boilerplate/src/repository/v1/example"
	"github.com/forkyid/go-consumer-boilerplate/src/service/v1/example"
	"gorm.io/gorm"
)

var master *gorm.DB
var slave *gorm.DB

func initMaster(wg *sync.WaitGroup) {
	defer wg.Done()
	master = database.DBMaster()
}

func initSlave(wg *sync.WaitGroup) {
	defer wg.Done()
	slave = database.DBSlave()
}

func getRoute() controller.Route {
	wg := &sync.WaitGroup{}

	// note: delete database-related file if not used (repository, service layers)
	wg.Add(2)
	go initMaster(wg)
	go initSlave(wg)
	wg.Wait()

	db := database.DB{
		Master: master,
		Slave:  slave,
	}

	exampleRepo := examplerepository.NewRepository(db)

	return controller.Route{
		ExchangeName: os.Getenv("EXCHANGE_NAME"),
		QueueName:    os.Getenv("QUEUE_NAME"),
		RouteName:    os.Getenv("ROUTE_NAME"),
		Handler: controller.NewController(
			example.NewService(exampleRepo),
		).Handler,
	}
}

func Start() {
	route := getRoute()
	route.Consume()
}
