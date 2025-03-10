// infrastructure/rabbitmq_consumer.go
package infrastructure

import (
	"consumer_api/src/products/application"
	"consumer_api/src/products/domain"
	_ "consumer_api/src/products/infraestructure/broker"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"sync"
	_ "time"

	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Permite todas las conexiones (ajustar según sea necesario)
}

// Lista de conexiones WebSocket activas
var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

type RabbitMQConsumer struct {
	channel   *amqp.Channel
	queueName string
	useCase   *application.ProductProcessingUseCase
}

func (r *RabbitMQConsumer) ProcessProduct(product *domain.Product) error {
	log.Printf("Procesando producto: %v", product)
	return nil
}

// NewRabbitMQConsumer crea un nuevo consumidor de RabbitMQ.
func NewRabbitMQConsumer(conn *amqp.Connection, queueName string, useCase *application.ProductProcessingUseCase) (*RabbitMQConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{channel: ch, queueName: queueName, useCase: useCase}, nil
}
func (r *RabbitMQConsumer) Consume(pb *PublisherController) error {
	msgs, err := r.channel.Consume(
		r.queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var product domain.Product
		err := json.Unmarshal(msg.Body, &product)
		if err != nil {
			log.Printf("Error al deserializar mensaje: %v", err)
			continue
		}

		// Procesar el producto
		err = r.ProcessProduct(&product)
		pb.Publisher()
		if err != nil {
			log.Printf("Error procesando el producto: %v", err)
			continue
		}

		// 🔥 Enviar mensaje al frontend mediante WebSocket
		notifyClients("Producto agregado con exito: " + product.Nombre)
	}

	return nil
}

// este funciona sin ws
/*func (r *RabbitMQConsumer) Consume( pb *PublisherController) error {
	// Consumir los mensajes de la cola
	msgs, err := r.channel.Consume(
		r.queueName, // Nombre de la cola
		"",          // Consumer (vacío para no especificar un consumidor único)
		true,        // Auto ack
		false,       // No exclusivo
		false,       // No esperar
		false,       // No arguments
		nil,         // No additional args
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		var product domain.Product
		err := json.Unmarshal(msg.Body, &product)
		if err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}
		// Procesar el producto
		err = r.ProcessProduct(&product)
		pb.Publisher()
		if err != nil {
			log.Printf("Error al procesar el producto: %v", err)
			continue
		}
		// Notificar al frontend (WebSocket)
		/*if err := r.NotifyFrontend("Producto agregado: " + product.Nombre); err != nil {
			log.Printf("Error notificando al frontend: %v", err)
		}
	}
	return nil
}*/
func notifyClients(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error enviando mensaje WebSocket: %v", err)
			client.Close()
			delete(clients, client) // Elimina clientes desconectados
		}
	}
}
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar conexión WebSocket: %v", err)
		return
	}

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	log.Println("Nueva conexión WebSocket establecida.")
	notifyClients("Conexion establecida Bienvenido")
}

// Función para notificar al frontend
/*func (r *RabbitMQConsumer) NotifyFrontend(message string) error {
	// Crear una conexión WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3001/web_socket", nil)
	if err != nil {
		return fmt.Errorf("Error al conectar con WebSocket: %v", err)
	}
	defer conn.Close()

	log.Printf("Mensaje de la cola recibido: %s", message)

	// Intentar enviar el mensaje varias veces si la conexión no está disponible
	for attempts := 0; attempts < 3; attempts++ {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err == nil {
			log.Printf("Mensaje enviado correctamente: %s", message)
			return nil
		}

		log.Printf("Error al enviar mensaje (intento %d): %v", attempts+1, err)
		time.Sleep(2 * time.Second) // Esperar un poco antes de reintentar
	}

	return fmt.Errorf("Error al enviar mensaje WebSocket después de varios intentos: %v", err)
}
*/
