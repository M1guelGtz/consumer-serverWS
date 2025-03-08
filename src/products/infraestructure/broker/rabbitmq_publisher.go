package broker

import (
    "github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
    channel *amqp.Channel
    queueName string
}

func NewRabbitMQPublisher(conn *amqp.Connection, queueName string) (*RabbitMQPublisher, error) {
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    _, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        return nil, err
    }

    return &RabbitMQPublisher{channel: ch, queueName: queueName}, nil
}

func (r *RabbitMQPublisher) Publish(message []byte) error {
    return r.channel.Publish("", r.queueName, false, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(message),
    })
}
