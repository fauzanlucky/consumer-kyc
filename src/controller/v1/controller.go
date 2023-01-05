package v1

import (
	"encoding/json"
	"log"

	"github.com/forkyid/go-consumer-boilerplate/src/entity/v1/http/member"
	"github.com/forkyid/go-consumer-boilerplate/src/service/v1/example"
	"github.com/streadway/amqp"
)

type Controller struct {
	// put necessary services here (or empty it)
	svc example.Servicer
}

func NewController(
	svc example.Servicer,
) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (ctrl *Controller) Handler(request *Request, d *amqp.Delivery) {
	data := member.Member{}
	err := json.Unmarshal(request.Body, &data)
	if err != nil {
		log.Println(err.Error())
		d.Ack(false) // d.Nack(false, true) if you want to requeue it.
		return
	}

	// do your magic!
	log.Println("priority:", d.Priority)
	log.Printf("%+v\n", data)

	d.Ack(false)
}
