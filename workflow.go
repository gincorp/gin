package main

import (
	"encoding/json"
)

type Workflow struct {
	Name      string
	Steps     []Step
	Variables map[string]string
}

func ParseWorkflow(data string) (w Workflow, err error) {
	err = json.Unmarshal([]byte(data), &w)

	return
}
