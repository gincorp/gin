package main

import (
    "flag"
    "log"
)

var (
    amqpUri *string
    redisUri *string
    mode *string
)

func init() {
    amqpUri = flag.String("amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")
    redisUri = flag.String("redis", "redis://localhost:6379/0", "URI of redis node")
    mode = flag.String("mode", "job", "mode with which to run")

    flag.Parse()

}

func main() {
    log.Printf("Using %q", *amqpUri)
    log.Printf("Running in %q mode", *mode)

    switch *mode {
    case "job", "master":
        node := NewNode(*amqpUri, *mode)

        if err := node.ConsumerLoop(); err != nil {
            log.Fatal(err)
        }

    default:
        log.Fatalf("Do not recognise mode '%q'", *mode)
    }
}
