package node

import (
	"github.com/streadway/amqp"
)

// Consumer provides a container for `Consumer` configuration and run time values
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

// NewConsumer will, given the URI of a rabbitMQ instance and a key with which to consume from,
// generate a Consumer for a node to receive messages germane to their operation
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
