package main

import (
    "encoding/json"
    "time"

    "github.com/satori/go.uuid"
)

type WorkflowRunner struct {
    EndTime time.Time
    Last string
    StartTime time.Time
    UUID string
    Variables map[string]interface{}
    Workflow Workflow
}

//type Variables map[string]interface{}

func NewWorkflowRunner(wf Workflow) (wfr WorkflowRunner) {
    wfr.UUID = uuid.NewV4().String()
    wfr.Workflow = wf
    wfr.Variables = make(map[string]interface{})
    wfr.Variables["Defaults"] = wf.Variables

    return
}

func ParseWorkflowRunner(data string)(wfr WorkflowRunner, err error) {
    err = json.Unmarshal([]byte(data), &wfr)

    return
}

func (wfr *WorkflowRunner)Start() {
    wfr.StartTime = time.Now()
}

func (wfr *WorkflowRunner)Next() (s Step, done bool) {
    var idx int

    if wfr.Last == "" {
        return wfr.Workflow.Steps[0], false
    }

    for idx,s = range wfr.Workflow.Steps {
        if s.Name == wfr.Last {break}
    }

    if idx + 1 >= len(wfr.Workflow.Steps) {
        return s, true
    }

    return wfr.Workflow.Steps[idx + 1], false
}

func (wfr *WorkflowRunner)End() {
    wfr.EndTime = time.Now()
}
