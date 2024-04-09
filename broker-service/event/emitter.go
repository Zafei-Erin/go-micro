package event

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, nil
	}

	return emitter, nil
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.PublishWithContext(context.Background(), "logs_topic", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(event),
	})

	if err != nil {
		return err
	}
	return nil
}
