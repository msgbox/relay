package queue

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Tag     string
	Done    chan error
}

// Build a Consumer to use in Workers that will consume an AMQP queue
//
// @param {String} tag(consumer_tag)
// @param {*amqp.Connection} connection
// @return {*Consumer}, {error}
// @api public
func BuildConsumer(tag string, connection *amqp.Connection) (*Consumer, error) {
	c := &Consumer{
		Conn:    connection,
		Channel: nil,
		Tag:     tag,
		Done:    make(chan error),
	}

	var err error

	c.Channel, err = connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	return c, nil
}

// Shutdown a Consumer by cancelling the delivery channel
func (c *Consumer) Shutdown() error {

	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.Tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	// wait for handle() to exit
	return <-c.Done
}
