package main

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
)

// WorkflowRunner ...
// Stateful representation of a Running Workflow
type WorkflowRunner struct {
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
}

// Next ...
// Return, should there be one, the next step of a Running Workflow
func (wfr *WorkflowRunner) Next() (s Step, done bool) {
	var idx int

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

// End ...
// Put a Running Workflow into an ended state
func (wfr *WorkflowRunner) End() {
	wfr.EndTime = time.Now()
}
