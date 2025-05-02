package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

// Consumer represents an AMQP consumer with a connection and queue name.
type Consumer struct {
	conn    *amqp.Connection
	queName string
}

// NewConsumer creates a new Consumer instance with the given AMQP connection and performs initial setup.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup initializes the consumer by opening a channel and declaring the exchange. Returns an error if the process fails.
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload represents a structured message with a name and associated data.
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen subscribes the consumer to the specified topics, binds them to an exchange, and processes incoming messages.
// Returns an error if the channel setup, queue declaration, or message consumption encounters an issue.
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)

		}
	}()

	fmt.Printf(" [*] Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)

	<-forever

	return nil

}

// handlePayload processes a given payload by performing actions based on its Name field, such as logging or handling "auth".
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// logEvent sends a JSON-formatted log entry to a remote logging service and returns an error if the operation fails.
// The function takes a Payload struct as input, marshals it to JSON, and sends it via an HTTP POST request.
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
