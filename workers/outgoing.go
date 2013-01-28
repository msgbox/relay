package worker

import (
	"fmt"
	"github.com/msgbox/relay/messages"
	"github.com/msgbox/relay/queue"
	"github.com/streadway/amqp"
	"log"
)

// An Outgoing Worker should listen for new items on the outgoing
// queue and process them as necessary

func Send(tag string, connection *amqp.Connection) {

	// Create a Consumer
	c, err := queue.BuildConsumer(tag, connection)
	if err != nil {
		fmt.Errorf("Consumer Error: %s", err)
	}

	// Start the Worker
	deliveries, err := c.Channel.Consume(
		"outgoing_messages", // name
		c.Tag,               // consumerTag,
		true,                // noAck
		false,               // exclusive
		false,               // noLocal
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		fmt.Errorf("Queue Consume: %s", err)
	}

	go handleOutgoing(deliveries, c.Done)

}

// Handle Outgoing messages from the queue
func handleOutgoing(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		msg := messages.Parse(d.Body)

		// For now just print them
		// TODO: Send them to the receiver via a socket connection
		log.Printf(
			"Sending Message:\n ID: %i \n To: %t \n From: %f \n At: %t \n",
			msg.GetId(),
			msg.GetReceiver(),
			msg.GetCreator(),
			msg.GetCreatedAt(),
		)

		d.Ack(true)
	}

	done <- nil
}
