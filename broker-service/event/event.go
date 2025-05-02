package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares a topic exchange named "logs_topic" with durable and non-auto-deleted settings.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

// declareRandomQueue declares a non-durable, exclusive, auto-deleted queue with a generated name on the provided channel.
// Returns the declared queue and any error encountered during the declaration process.
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
