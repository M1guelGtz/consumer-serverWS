package infrastructure

import (
	"consumer_api/src/products/application"
	"consumer_api/src/products/infraestructure/broker"
	"log"
	"net/http"

	_ "github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func Init() {
	http.HandleFunc("/web_socket", WebSocketHandler)

	go func() {
		log.Println("Servidor WebSocket escuchando en :3500")
		log.Fatal(http.ListenAndServe(":3500", nil))
	}()
	conn, err := amqp.Dial("amqp://miguel:7s0725FLU2@50.19.162.154:5672/")
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()
	productConsumer := &RabbitMQConsumer{}
	publisher, err := broker.NewRabbitMQPublisher(conn, "Q2")
	if err != nil {
		log.Fatal("Error al crear el publicador de RabbitMQ", err)
	}
	controller := NewController(publisher)
	useCase := application.NewProductProcessingUseCase(productConsumer)
	rabbitMQConsumer, err := NewRabbitMQConsumer(conn, "Q1", useCase)
	if err != nil {
		log.Fatalf("Error al crear el consumidor de RabbitMQ: %v", err)
	}

	go func() {
		err := rabbitMQConsumer.Consume(controller)
		controller.Publisher()
		if err != nil {
			log.Fatalf("Error al consumir los mensajes de RabbitMQ: %v", err)
		}

	}()

	select {}
}
