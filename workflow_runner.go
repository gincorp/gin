package main

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
)

// WorkflowRunner ...
// Stateful representation of a Running Workflow
type WorkflowRunner struct {
	State     string
	EndTime   time.Time
	Last      string
	StartTime time.Time
	UUID      string
	Variables map[string]interface{}
	Workflow  Workflow
}

// NewWorkflowRunner ...
// Initialise and Return a WorkflowRunner
func NewWorkflowRunner(wf Workflow) (wfr WorkflowRunner) {
	wfr.UUID = uuid.NewV4().String()
	wfr.Workflow = wf
	wfr.Variables = make(map[string]interface{})
	wfr.Variables["Defaults"] = wf.Variables

	return
}

// ParseWorkflowRunner ...
// Parse a Running Workflow from a stored state
func ParseWorkflowRunner(data string) (wfr WorkflowRunner, err error) {
	err = json.Unmarshal([]byte(data), &wfr)

	return
}

// Start ...
// Put a Running Workflow into a started state
func (wfr *WorkflowRunner) Start() {
	wfr.StartTime = time.Now()
	wfr.State = "started"
}

// Next ...
// Return, should there be one, the next step of a Running Workflow
func (wfr *WorkflowRunner) Next() (s Step, done bool) {
	var idx int
	wfr.State = "running"

	if wfr.Last == "" {
		return wfr.Workflow.Steps[0], false
	}

	for idx, s = range wfr.Workflow.Steps {
		if s.Name == wfr.Last {
			break
		}
	}

	if idx+1 >= len(wfr.Workflow.Steps) {
		return s, true
	}

	return wfr.Workflow.Steps[idx+1], false
}

// Current returns the current step. It is used, mainly,
// after a step has returned to add extra data
func (wfr *WorkflowRunner) Current() (i int, s Step) {
	for i, s = range wfr.Workflow.Steps {
		if s.Name == wfr.Last {
			return
		}
	}

	return
}

// Fail will set state to "failed" and end the workflow runner
func (wfr *WorkflowRunner) Fail() {
	wfr.endWithState("failed")
}

// End will set state to "ended" and end the workflow runner
func (wfr *WorkflowRunner) End() {
	wfr.endWithState("ended")
}

func (wfr *WorkflowRunner) endWithState(state string) {
	wfr.EndTime = time.Now()
	wfr.State = state
}
