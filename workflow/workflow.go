package workflow

import (
    "encoding/json"
)

// Workflow ...
// Raw workflow configuration container,
// reflects config in storage and without a state machine
type Workflow struct {
    Name      string
    Steps     []Step
    Variables map[string]string
}

// ParseWorkflow ...
// Return a Workflow from a textual representation from storage
func ParseWorkflow(data string) (w Workflow, err error) {
    err = json.Unmarshal([]byte(data), &w)

    return
}
