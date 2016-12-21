package workflow

import (
	"encoding/json"
	"time"
)

// Runner is a stateful representation of a Running Workflow
type Runner struct {
	EndTime      time.Time
	ErrorMessage string
	Last         string
	StartTime    time.Time
	State        string
	UUID         string
	Variables    map[string]interface{}
	Workflow     Workflow
}

// NewRunner initialises and Return a Runner
func NewRunner(uuid string, wf Workflow) (wfr Runner) {
	wfr.UUID = uuid
	wfr.Workflow = wf
	wfr.Variables = make(map[string]interface{})
	wfr.Variables["Defaults"] = wf.Variables

	return
}

// ParseRunner returns a parsed Runner from a string
func ParseRunner(data string) (wfr Runner, err error) {
	if data == "" {
		return
	}
	err = json.Unmarshal([]byte(data), &wfr)

	return
}

// Start puts a Running Workflow into a started state
func (wfr *Runner) Start() {
	wfr.StartTime = time.Now()
	wfr.State = "started"
}

// Next returns the next Step of a Runner
func (wfr *Runner) Next() (s Step, done bool) {
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
func (wfr *Runner) Current() (i int, s Step) {
	for i, s = range wfr.Workflow.Steps {
		if s.Name == wfr.Last {
			return
		}
	}

	return
}

// Fail will set state to "failed" and end the workflow runner
func (wfr *Runner) Fail(msg string) {
	wfr.ErrorMessage = msg
	wfr.endWithState("failed")
}

// End will set state to "ended" and end the workflow runner
func (wfr *Runner) End() {
	wfr.endWithState("ended")
}

func (wfr *Runner) endWithState(state string) {
	wfr.EndTime = time.Now()
	wfr.State = state
}
