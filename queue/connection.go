package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var (
	uri          = "amqp://guest:guest@localhost:5672/"
	exchangeName = "msgbox"
	exchangeType = "direct"
	reliable     = true
)

// Create a AMQP Connection
// This allows the service to publish and consume from a
// queue. It should be closed using a defer after being opened.
//
// @return {*amqp.Connection}, {error v}
// @api public
func Connect() (*amqp.Connection, error) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("Dial: %s", err)
		return nil, fmt.Errorf("Dial: %s", err)
	}

	// Declare the AMQP Topology to ensure it's
	// setup to handle the correct exchanges and queues
	t_err := declareTopology(connection)
	if t_err != nil {
		return nil, fmt.Errorf("Topology: %s", err)
	}

	return connection, nil
}

// Declare a Common Topology to ensure the structure exists
// and is correct before publishing or consuming against it
//
// @param {&amqp.Connection} connection
// @return {Error}
// @api private
func declareTopology(connection *amqp.Connection) error {

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	defer channel.Close()

	if err = channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Create the "incoming_messages" queue
	// This holds all the unprocessed incoming messages
	// received by the relay
	i_state, err := channel.QueueDeclare(
		"incoming_messages", // name of the queue
		true,                // durable
		false,               // delete when usused
		false,               // exclusive
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}
	log.Printf("Declared Queue: %s", *&i_state.Messages)

	if err = channel.QueueBind(
		"incoming_messages", // name of the queue
		"incoming",          // bindingKey
		exchangeName,        // sourceExchange
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	// Create the "outgoing_messages" queue
	// This holds all the unprocessed outgoing messages
	// that have not been sent yet
	o_state, err := channel.QueueDeclare(
		"outgoing_messages", // name of the queue
		true,                // durable
		false,               // delete when usused
		false,               // exclusive
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}
	log.Printf("Declared Queue: %s", *&o_state.Messages)

	if err = channel.QueueBind(
		"outgoing_messages", // name of the queue
		"outgoing",          // bindingKey
		exchangeName,        // sourceExchange
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	return nil
}
