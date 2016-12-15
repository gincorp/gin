package taskmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// JobManager contains configuration for Job Task Managers
type JobManager struct {
	// Map job aliases to functions providing that job.
	// Jobs receive a notification providing data from the queue and return
	// maps for json to handle
	JobList map[string]func(JobNotification) (map[string]interface{}, error)
}

// JobNotification is a container for Unmarshal'd json from the message queue
type JobNotification struct {
	Context  map[string]string
	Name     string
	Register string
	Type     string
	UUID     string
}

// NewJobManager returns a `JobManager` to route jobs from the queue
func NewJobManager() (j JobManager) {
	j.JobList = make(map[string]func(JobNotification) (map[string]interface{}, error))

	j.AddJob("post-to-web", doWebCall)
	j.AddJob("get-from-web", doWebCall)
	j.AddJob("log", logOutput)

	return
}

// AddJob updates j.JobList with the key 'key' with the value of a function to call
func (j *JobManager) AddJob(key string, f func(JobNotification) (map[string]interface{}, error)) {
	j.JobList[key] = f
}

// Consume handles json from the message queue. It formats it correctly,
// route the job, and returns output and metadata
func (j JobManager) Consume(body string) (output map[string]interface{}, err error) {
	jn := j.parseBody(body)

	output = make(map[string]interface{})

	output["UUID"] = jn.UUID
	output["Register"] = jn.Register
	output["Failed"] = false

	start := time.Now()
	output["Data"], err = j.JobList[jn.Type](jn)
	end := time.Now()

	if err != nil {
		output["Failed"] = true
		output["ErrorMessage"] = err.Error()
	}

	format := "2006-01-02T15:04:05"
	output["Start"] = start.Format(format)
	output["End"] = end.Format(format)
	output["Duration"] = fmt.Sprintf("%d ms", (end.UnixNano()-start.UnixNano())/1000000)

	return
}

func (j JobManager) parseBody(b string) (n JobNotification) {
	if err := json.Unmarshal([]byte(b), &n); err != nil {
		log.Println(err)
	}

	return
}
