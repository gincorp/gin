package main

import (
    "encoding/json"
    "log"
)

type JobManager struct {
    JobList map[string]func(JobNotification)(map[string]interface {}, error)
}

type JobNotification struct {
    Context map[string]string
    Name string
    Type string
}

func NewJobManager() (j JobManager) {
    j.JobList = make(map[string]func(JobNotification)(map[string]interface {}, error))
    j.JobList["post-to-web"] = doWebCall
    j.JobList["get-from-web"] = doWebCall

    return
}

func (j JobManager) Consume(body string) (output map[string]interface{}, err error) {
    jn := j.ParseBody(body)
    output, err = j.JobList[jn.Type](jn)

    return
}

func (j JobManager) ParseBody(b string) (n JobNotification) {
    if err := json.Unmarshal([]byte(b), &n); err != nil {
        log.Println(err)
    }

    return
}
