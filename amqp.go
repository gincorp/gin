package main

import (
    "fmt"
    "github.com/streadway/amqp"
    "log"
)

type TaskManager interface {
    Consume(string)(map[string]interface{}, error)
}

type Consumer struct {
    channel *amqp.Channel
    conn    *amqp.Connection
    done    chan error
    queue   string
    tag     string
    tm TaskManager
}

func NewConsumer(uri, key string) (*Consumer, error) {
    exchangeName := "workflow.exchange"

    var tm TaskManager
    switch key {
    case "job":
        tm = NewJobManager()
    case "master":
        tm = NewMasterManager()
    default:
        log.Fatalf("Key %q is invalid", key)
    }

    c := &Consumer{
        channel: nil,
        conn:    nil,
        done:    make(chan error),
        queue:   key,
        tag:     key,
        tm:      tm,
    }

    var err error

    c.conn, err = amqp.Dial(uri)
    if err != nil {
        return nil, fmt.Errorf("Dial: %s", err)
    }

    go func() {
        fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
    }()

    log.Printf("got Connection, getting Channel")
    c.channel, err = c.conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("Channel: %s", err)
    }

    log.Printf("got Channel, declaring Exchange (%q)", exchangeName)
    if err = c.channel.ExchangeDeclare(
        exchangeName,                      // name of the exchange
        "direct",                          // type
        true,                              // durable
        false,                             // delete when complete
        false,                             // internal
        false,                             // noWait
        nil,                               // arguments
    ); err != nil {
        return nil, fmt.Errorf("Exchange Declare: %s", err)
    }

    log.Printf("declared Exchange, declaring Queue (%q)", c.queue)
    queue, err := c.channel.QueueDeclare(
        c.queue,                       // name of the queue
        true,                          // durable
        false,                         // delete when usused
        false,                         // exclusive
        false,                         // noWait
        nil,                           // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Declare: %s", err)
    }

    log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
        queue.Name, queue.Messages, queue.Consumers, key)

    if err = c.channel.QueueBind(
        queue.Name,                     // name of the queue
        key,                            // bindingKey
        exchangeName,                   // sourceExchange
        false,                          // noWait
        nil,                            // arguments
    ); err != nil {
        return nil, fmt.Errorf("Queue Bind: %s", err)
    }

    log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
    deliveries, err := c.channel.Consume(
        queue.Name, // name
        c.tag,      // consumerTag,
        false,      // noAck
        false,      // exclusive
        false,      // noLocal
        false,      // noWait
        nil,        // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Consume: %s", err)
    }

    go c.handle(deliveries, c.done)
    select{}
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

func (c *Consumer) handle(deliveries <-chan amqp.Delivery, done chan error) {
    for d := range deliveries {
        log.Printf("[%v] : %q received %q", d.DeliveryTag, c.queue, d.Body)

        if output, err := c.tm.Consume( string(d.Body) ); err != nil {
            log.Printf("[%v] : errors %q", d.DeliveryTag, err)
            d.Ack(false)
        } else {
            log.Printf("[%v] : returned %q", d.DeliveryTag, output)
            d.Ack(true)
        }
    }
    log.Printf("handle: deliveries channel closed")
    done <- nil
}
