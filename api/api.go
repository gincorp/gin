package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jspc/workflow-engine/datastore"
	"github.com/jspc/workflow-engine/node"
)

// StarterRequest is a placeholder for incoming 'start workflow' requests
type StarterRequest struct {
	Name      string
	Variables map[string]interface{}
}

// StarterResponse is a placeholder for outgoing 'start workflow' responses
type StarterResponse struct {
	UUID string
}

// ErrorResponse is a placeholder for outgoing errors in responses
type ErrorResponse struct {
	Message string
}

// StartWorkflow contains the message structure that kicking off a workflow
// expects in a master TaskManager
type StartWorkflow struct {
	InitWorkflow StarterRequest
	Time         time.Time
	UUID         string
}

// API contains routing and messagign capabilities. Valid HTTP requests
// to api.API are routed via rabbitmq to either kick off or report upon
// workflows.
type API struct {
	datastore datastore.Datastore
	listener  string
	producer  *node.Producer
}

var (
	routingKey = "master"
)

// NewAPI receives an ampquri for routing messages, a redisURI for
// configuring workflows, and a host and a port
// for listening on.
func NewAPI(ampqURI, redisURI, host string, port int) (a API, err error) {
	a.producer = node.NewProducer(ampqURI, routingKey)
	a.listener = fmt.Sprintf("%s:%d", host, port)

	a.datastore, err = datastore.NewDatastore(redisURI)

	return
}

// Start starts a listener on API.listener
// Used for routing metadata, monitoring, and
// workflow requests.
func (a API) Start() {
	http.HandleFunc("/mon/", a.monRoute)
	http.HandleFunc("/wf/", a.wfRoute)

	log.Fatal(http.ListenAndServe(a.listener, nil))
}
