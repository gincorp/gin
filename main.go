package main

import (
	"flag"
	"log"

	"github.com/gincorp/gin/api"
	"github.com/gincorp/gin/node"
)

var (
	amqpURI  *string
	mode     *string
	redisURI *string
	host     *string
	port     *int
)

func init() {
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")
	redisURI = flag.String("redis", "redis://localhost:6379/0", "URI of redis node")
	mode = flag.String("mode", "job", "mode with which to run")
	host = flag.String("host", "0.0.0.0", "In api mode; host to bind to")
	port = flag.Int("port", 8080, "In api mode; port to listen on")

	flag.Parse()

}

func main() {
	log.Printf("Using %q", *amqpURI)
	log.Printf("Running in %q mode", *mode)

	switch *mode {
	case "job", "master":
		n := node.NewNode(*amqpURI, *redisURI, *mode)

		if err := n.ConsumerLoop(); err != nil {
			log.Fatal(err)
		}

	case "api":
		a, err := api.NewAPI(*amqpURI, *redisURI, *host, *port)
		if err != nil {
			log.Fatal(err)
		}

		a.Start()

	default:
		log.Fatalf("Do not recognise mode '%q'", *mode)
	}
}
