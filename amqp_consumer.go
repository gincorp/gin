package main

import (
    "fmt"
    "log"

    "github.com/streadway/amqp"
)


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

func NewConsumer(uri, key string) (*Consumer) {
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
