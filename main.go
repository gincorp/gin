package main

import (
	"flag"
	"log"
)

var (
	amqpURI  *string
	mode     *string
	redisURI *string

	node Node
)

func init() {
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")
	redisURI = flag.String("redis", "redis://localhost:6379/0", "URI of redis node")
	mode = flag.String("mode", "job", "mode with which to run")

	flag.Parse()

}

func main() {
	log.Printf("Using %q", *amqpURI)
	log.Printf("Running in %q mode", *mode)

	switch *mode {
	case "job", "master":
		node = NewNode(*amqpURI, *mode)

		go node.TaskManager.StartAPI()

		if err := node.ConsumerLoop(); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("Do not recognise mode '%q'", *mode)
	}
}
