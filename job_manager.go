package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// JobManager ...
// Job configuration
type JobManager struct {
	// Map job aliases to functions providing that job.
	// Jobs receive a notification providing data from the queue and return
	// maps for json to handle
	JobList map[string]func(JobNotification) (map[string]interface{}, error)
}

// JobNotification ...
// Data container for Unmarshal'd json from the message queue
type JobNotification struct {
	Context  map[string]string
	Name     string
	Register string
	Type     string
	UUID     string
}

// NewJobManager ...
// Return a `JobManager` to route jobs fromt the queue
func NewJobManager() (j JobManager) {
	j.JobList = make(map[string]func(JobNotification) (map[string]interface{}, error))
	j.JobList["post-to-web"] = doWebCall
	j.JobList["get-from-web"] = doWebCall
	j.JobList["log"] = logOutput

	return
}

// Consume ...
// Handle json from the message queue. Format it correctly, route the job, and
// return output and metadata
func (j JobManager) Consume(body string) (output map[string]interface{}, err error) {
	jn := j.parseBody(body)

	output = make(map[string]interface{})

	output["UUID"] = jn.UUID
	output["Register"] = jn.Register

	start := time.Now().UnixNano()
	output["Data"], err = j.JobList[jn.Type](jn)
	end := time.Now().UnixNano()

	output["Duration"] = fmt.Sprintf("%d ms", (end-start)/1000000)

	return
}

func (j JobManager) parseBody(b string) (n JobNotification) {
	if err := json.Unmarshal([]byte(b), &n); err != nil {
		log.Println(err)
	}

	return
}
