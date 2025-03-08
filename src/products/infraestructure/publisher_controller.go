package infrastructure

import (
	"consumer_api/src/products/infraestructure/broker"
	"consumer_api/src/products/infraestructure/handlers"
)

type PublisherController struct {
	PublisherHandler *handlers.PublisherHandler
}

func NewController(
	publisher *broker.RabbitMQPublisher, // 🔹 Se agrega el publisher de RabbitMQ
) *PublisherController {
	// 🔹 Se pasa el publisher al handler de creación
	createHandler := handlers.NewCreateProductHandler(publisher)
	return &PublisherController{
		PublisherHandler: createHandler,
	}
}

func (pc *PublisherController) Publisher() {
	pc.PublisherHandler.Handle()
}
