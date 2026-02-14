package messaging

import (
	"encoding/json"
	"log"
	"ms1-documents/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type messagePublisher struct {
	messaging *config.Messaging
}

type DocumentCreatedMessage struct {
	DocumentID string `json:"documentId"`
	UUID       string `json:"uuid"`
}

func NewMessagePublisher(messaging *config.Messaging) Publisher {
	return &messagePublisher{
		messaging: messaging,
	}
}

func (p *messagePublisher) PublishDocumentCreated(documentID, uuid string) error {
	message := DocumentCreatedMessage{
		DocumentID: documentID,
		UUID:       uuid,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.messaging.Channel.Publish(
		"documents",
		"documents.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Mensaje publicado: %s", body)
	return nil
}
