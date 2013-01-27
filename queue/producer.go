package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Publish a Protocol Buffer to an AMQP Exchange
//
// @param {String} key
// @param {[]byte} body
// @param {*amqp.Connection} connection
// @return {error}
// @api public
func Publish(key string, body []byte, connection *amqp.Connection) error {
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	defer channel.Close()

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64), make(chan uint64))
		defer confirmOne(ack, nack)
	}

	if err = channel.Publish(
		exchangeName, // publish to an exchange
		key,          // routing to 0 or more queues
		true,         // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
			Priority:        0, // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}
