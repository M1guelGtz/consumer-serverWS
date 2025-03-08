package handlers

import (
	"consumer_api/src/products/domain"
	"consumer_api/src/products/infraestructure/broker"
	"encoding/json"
	"log"
)

type PublisherHandler struct {
	Publisher *broker.RabbitMQPublisher
}

func NewCreateProductHandler(publisher *broker.RabbitMQPublisher) *PublisherHandler {
	return &PublisherHandler{Publisher: publisher}
}

func (h *PublisherHandler) Handle() {
	var product domain.Product

	// Crear un objeto JSON con los detalles del producto
	message := map[string]interface{}{
		"nombre":   product.Nombre,
		"cantidad": product.Cantidad,
		"precio":   product.Precio,
	}

	// Serializar el objeto JSON a bytes
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Error al serializar el mensaje:", err)
		return
	}

	// Publicar el mensaje en RabbitMQ
	err = h.Publisher.Publish(messageBytes)
	log.Println("mensaje publicado a la cola 2")
	if err != nil {
		log.Println("Error publicando mensaje en RabbitMQ:", err)
	}
}
