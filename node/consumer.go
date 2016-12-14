package node

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Consumer ...
// Container for `Consumer` configuration and run time values
type Consumer struct {
	channel *amqp.Channel
	conn    *amqp.Connection
	done    chan error
	exch    string
	key     string
	queue   string
	tag     string
	uri     string
}

// NewConsumer ...
// Given the URI of a rabbitMQ instance and a key with which to consume from,
// generate a Consumer; `Node`s use this receive messages
func NewConsumer(uri, key string) *Consumer {
	exchangeName := "workflow.exchange"

	c := &Consumer{
		channel: nil,
		conn:    nil,
		done:    make(chan error),
		exch:    exchangeName,
		key:     key,
		queue:   key,
		tag:     key,
		uri:     uri,
	}

	return c
}

// Shutdown ...
// Close AMPQ connections
func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}
