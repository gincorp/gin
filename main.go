package main

import (
    "flag"
    // "fmt"
    "log"
    // "net/http"
    // "time"

    // "github.com/boltdb/bolt"
    // "github.com/gorilla/context"
    // "golang.org/x/crypto/ssh/terminal"
)

var (
    amqpUri *string
    mode *string
)

func init() {
    amqpUri = flag.String("amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")
    mode = flag.String("mode", "job", "mode with which to run")

    flag.Parse()

}

func main() {
    log.Printf("Using %q", *amqpUri)
    log.Printf("Running in %q mode", *mode)

    switch *mode {
    case "job", "master":
        _, err := NewConsumer(*amqpUri, *mode)

        if err != nil {
            log.Fatal(err.Error())
        }

    default:
        log.Fatalf("Do not recognise mode '%q'", *mode)
    }
}
