package node

import (
    "fmt"
    "log"

    "github.com/streadway/amqp"
)

// Producer configuration for sending messages via rabbit MQ
type Producer struct {
    exch string
    key  string
    uri  string
}

// NewProducer creates configuration container *Producer
func NewProducer(uri, key string) *Producer {
    exchangeName := "workflow.exchange"

    p := &Producer{
        exch: exchangeName,
        key:  key,
        uri:  uri,
    }

    return p
}

// Send a payload via rabbit amqp
func (p *Producer) Send(body []byte) error {
    connection, err := amqp.Dial(p.uri)
    if err != nil {
        return fmt.Errorf("Dial: %s", err)
    }

    defer connection.Close()

    channel, err := connection.Channel()
    if err != nil {
        return fmt.Errorf("Channel: %s", err)
    }

    if err := channel.Confirm(false); err != nil {
        return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
    }

    confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

    defer confirmOne(confirms)

    if err = channel.Publish(
        p.exch, // publish to an exchange
        p.key,  // routing to 0 or more queues
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            Headers:         amqp.Table{},
            ContentType:     "text/plain",
            ContentEncoding: "",
            Body:            body,
            DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
            Priority:        0,              // 0-9
        },
    ); err != nil {
        return fmt.Errorf("Exchange Publish: %s", err)
    }

    return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
    if confirmed := <-confirms; confirmed.Ack {
        log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
    } else {
        log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
    }

}
