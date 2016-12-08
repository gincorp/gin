package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/streadway/amqp"
)

type Node struct {
    Consumer *Consumer
    Producer *Producer
}

var (
    consumerKey, producerKey string
)

func NewNode(uri, nodeMode string) (n Node) {
    switch nodeMode {
    case "job":
        consumerKey = "job"
        producerKey = "master"
    case "master":
        consumerKey = "master"
        producerKey = "job"
    }


    c := NewConsumer(uri, consumerKey)
    p := NewProducer(uri, producerKey)

    n.Consumer = c
    n.Producer = p
    return
}

func (n *Node)ConsumerLoop() (err error) {
    if n.Consumer.conn, err = amqp.Dial(n.Consumer.uri); err != nil {
        return fmt.Errorf("Dial: %s", err)
    }

    go func() {
        fmt.Printf("closing: %s", <-n.Consumer.conn.NotifyClose(make(chan *amqp.Error)))
    }()

    log.Printf("got Connection, getting Channel")
    if n.Consumer.channel, err = n.Consumer.conn.Channel(); err != nil {
        return fmt.Errorf("Channel: %s", err)
    }

    log.Printf("got Channel, declaring Exchange (%q)", n.Consumer.exch)
    if err = n.Consumer.channel.ExchangeDeclare(
        n.Consumer.exch,                            // name of the exchange
        "direct",                          // type
        true,                              // durable
        false,                             // delete when complete
        false,                             // internal
        false,                             // noWait
        nil,                               // arguments
    ); err != nil {
        return fmt.Errorf("Exchange Declare: %s", err)
    }

    log.Printf("declared Exchange, declaring Queue (%q)", n.Consumer.queue)
    queue, err := n.Consumer.channel.QueueDeclare(
        n.Consumer.queue,                       // name of the queue
        true,                          // durable
        false,                         // delete when usused
        false,                         // exclusive
        false,                         // noWait
        nil,                           // arguments
    )
    if err != nil {
        return fmt.Errorf("Queue Declare: %s", err)
    }

    log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
        queue.Name, queue.Messages, queue.Consumers, n.Consumer.key)

    if err = n.Consumer.channel.QueueBind(
        queue.Name,                     // name of the queue
        n.Consumer.key,                            // bindingKey
        n.Consumer.exch,                         // sourceExchange
        false,                          // noWait
        nil,                            // arguments
    ); err != nil {
        return fmt.Errorf("Queue Bind: %s", err)
    }

    log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", n.Consumer.tag)
    deliveries, err := n.Consumer.channel.Consume(
        queue.Name, // name
        n.Consumer.tag,      // consumerTag,
        false,      // noAck
        false,      // exclusive
        false,      // noLocal
        false,      // noWait
        nil,        // arguments
    )
    if err != nil {
        return fmt.Errorf("Queue Consume: %s", err)
    }

    go n.Consume(deliveries, n.Consumer.done)

    select{}
}

func (n *Node)Consume(deliveries <-chan amqp.Delivery, done chan error) {
    for d := range deliveries {
        log.Printf("[%v] : %q received %q", d.DeliveryTag, n.Consumer.queue, d.Body)

        if output, err := n.Consumer.tm.Consume( string(d.Body) ); err != nil {
            log.Printf("[%v] : errors %q", d.DeliveryTag, err)

            d.Ack(false)
        } else {
            if n.Consumer.tm.ShouldRespond() {
                go func() {
                    log.Printf("[%v] : responding with %q", d.DeliveryTag, output)

                    if err := n.Deliver(output); err != nil {
                        log.Printf("[%v] : response errored: %q", d.DeliveryTag, err)

                        d.Ack(false)
                    } else {
                        log.Printf("[%v] : responded", d.DeliveryTag)

                        d.Ack(true)
                    }
                }()
            } else {
                d.Ack(true)
            }
        }
    }
    log.Printf("handle: deliveries channel closed")
    done <- nil
}

func (n *Node)Deliver(message interface{}) error {
    j, err := json.Marshal(message)

    if err != nil {
        return err
    }

    return n.Producer.send(j)
}
