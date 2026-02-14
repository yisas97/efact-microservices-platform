package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Messaging struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func NewMessaging(uri string) (*Messaging, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"documents",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"documents.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		"documents.created",
		"documents.created",
		"documents",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	log.Println("Conectado a RabbitMQ")
	return &Messaging{
		Channel:    ch,
		Connection: conn,
	}, nil
}

func (m *Messaging) Disconnect() {
	if m.Channel != nil {
		m.Channel.Close()
	}
	if m.Connection != nil {
		m.Connection.Close()
	}
}
