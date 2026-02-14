package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ColaVerificacion        = "verify.request"
	TimeoutVerificacion     = 10 * time.Second
	TimeoutConexionRabbitMQ = 5 * time.Second
)

type SolicitudVerificacion struct {
	Documento interface{} `json:"documento"`
	Firma     string      `json:"firma"`
}

type RespuestaVerificacion struct {
	Valido  bool   `json:"valido"`
	Mensaje string `json:"mensaje"`
}

type ClienteRPC struct {
	conexion *amqp.Connection
	canal    *amqp.Channel
}

func NuevoClienteRPC(url string) (*ClienteRPC, error) {
	conexion, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a RabbitMQ: %w", err)
	}

	canal, err := conexion.Channel()
	if err != nil {
		conexion.Close()
		return nil, fmt.Errorf("error al abrir canal: %w", err)
	}

	return &ClienteRPC{
		conexion: conexion,
		canal:    canal,
	}, nil
}

func (c *ClienteRPC) Cerrar() {
	if c.canal != nil {
		c.canal.Close()
	}
	if c.conexion != nil {
		c.conexion.Close()
	}
}

func (c *ClienteRPC) VerificarDocumento(ctx context.Context, solicitud *SolicitudVerificacion) (*RespuestaVerificacion, error) {
	colaRespuesta, err := c.canal.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear cola de respuesta: %w", err)
	}

	mensajes, err := c.canal.Consume(
		colaRespuesta.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error al consumir cola de respuesta: %w", err)
	}

	correlationId := uuid.New().String()

	cuerpo, err := json.Marshal(solicitud)
	if err != nil {
		return nil, fmt.Errorf("error al serializar solicitud: %w", err)
	}

	err = c.canal.PublishWithContext(
		ctx,
		"",
		ColaVerificacion,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       colaRespuesta.Name,
			Body:          cuerpo,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error al publicar mensaje: %w", err)
	}

	timeout := time.After(TimeoutVerificacion)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout esperando respuesta de MS2")

		case <-ctx.Done():
			return nil, ctx.Err()

		case mensaje := <-mensajes:
			if mensaje.CorrelationId == correlationId {
				var respuesta RespuestaVerificacion
				if err := json.Unmarshal(mensaje.Body, &respuesta); err != nil {
					return nil, fmt.Errorf("error al deserializar respuesta: %w", err)
				}
				return &respuesta, nil
			}
		}
	}
}
