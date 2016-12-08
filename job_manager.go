package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
)

type JobManager struct {
    JobList map[string]func(JobNotification)(map[string]interface {}, error)
}

type JobNotification struct {
    Context map[string]string
    Name string
    Type string
    UUID string
}

func NewJobManager() (j JobManager) {
    j.JobList = make(map[string]func(JobNotification)(map[string]interface {}, error))
    j.JobList["post-to-web"] = doWebCall
    j.JobList["get-from-web"] = doWebCall

    return
}

func (j JobManager) Consume(body string) (output map[string]interface{}, err error) {
    jn := j.ParseBody(body)

    output = make(map[string]interface{})

    output["UUID"] = jn.UUID

    start := time.Now().UnixNano()
    output["Data"], err = j.JobList[jn.Type](jn)
    end := time.Now().UnixNano()

    output["Duration"] = fmt.Sprintf("%d ms", (end - start) / 1000000)

    return
}

func (j JobManager) ParseBody(b string) (n JobNotification) {
    if err := json.Unmarshal([]byte(b), &n); err != nil {
        log.Println(err)
    }

    return
}
