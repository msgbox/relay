package worker

import (
	"fmt"
	"github.com/msgbox/relay/messages"
	"github.com/msgbox/relay/queue"
	"github.com/streadway/amqp"
	"log"
)

// An Incoming Worker should listen for new items on the incoming
// queue and process them as necessary

func Receive(tag string, connection *amqp.Connection) {

	// Create a Consumer
	c, err := queue.BuildConsumer(tag, connection)
	if err != nil {
		fmt.Errorf("Consumer Error: %s", err)
	}

	// Start the Worker
	deliveries, err := c.Channel.Consume(
		"incoming_messages", // name
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

	go handleIncoming(deliveries, c.Done)

}

// Handle Incoming Items off the queue
func handleIncoming(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		msg := messages.Parse(d.Body)

		// For now just print them
		// TO-DO: Write items to a database or storage
		log.Printf(
			"Receiving Message:\n ID: %i \n To: %t \n From: %f \n Created At: %t \n",
			msg.GetId(),
			msg.GetReceiver(),
			msg.GetCreator(),
			msg.GetCreatedAt(),
		)

		d.Ack(true)
	}

	done <- nil
}
