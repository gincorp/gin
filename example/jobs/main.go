package main

import (
    "flag"
    "log"

    "github.com/gincorp/gin/node"
    "github.com/gincorp/gin/taskmanager"
)

var (
    amqpURI *string
)

func init() {
    amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")

    flag.Parse()

}

func main() {
    log.Printf("Using %q", *amqpURI)

    n := node.NewNode(*amqpURI, "", "job")

    jobManager := taskmanager.NewJobManager()
    jobManager.AddJob("get-financial", getCurrencyPrices)
    jobManager.AddJob("send-email", sendEmail)

    n.TaskManager = jobManager

    log.Fatal(n.ConsumerLoop())
}
