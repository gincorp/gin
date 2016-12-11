package main

import (
	"encoding/json"
	"log"
)

// MasterManager ...
// Container for Master Task manager configuration
type MasterManager struct {
	datastore Datastore
}

// NewMasterManager ...
// Initialise and return a Master Task Manager
func NewMasterManager() (m MasterManager) {
	var err error

	if m.datastore, err = NewDatastore(*redisURI); err != nil {
		log.Fatal(err)
	}

	return
}

// Consume ...
// Handle json from the message queue; for a Master node these will be responses.
// Parse messages, update Workflow contexts, write to database and call next step
func (m MasterManager) Consume(body string) (output map[string]interface{}, err error) {
	var b interface{}
	var wfr WorkflowRunner

	if err = json.Unmarshal([]byte(body), &b); err != nil {
		return
	}

	output = b.(map[string]interface{})
	uuid := output["UUID"].(string)
	if wfr, err = m.datastore.LoadWorkflowRunner(uuid); err != nil {
		return
	}

	switch output["Register"].(type) {
	case string:
		register := output["Register"].(string)

		switch output["Data"].(type) {
		case map[string]interface{}:
			data := output["Data"].(map[string]interface{})
			wfr.Variables[register] = data

		default:
			log.Println("Not registering output: got garbage back")
		}
	}

	m.datastore.DumpWorkflowRunner(wfr)
	m.Continue(wfr.UUID)

	return
}

// Load ...
// Load a workflow from storage and create a WorkflowRunner state machine
func (m MasterManager) Load(name string) (uuid string, err error) {
	wf, err := m.datastore.LoadWorkflow(name)
	if err != nil {
		return
	}

	wfr := NewWorkflowRunner(wf)
	wfr.Start()

	m.datastore.DumpWorkflowRunner(wfr)

	return wfr.UUID, nil
}

// Continue ...
// Should there be a next step in the workflow, compile step templates
// and push the step to the emssage queue
func (m MasterManager) Continue(uuid string) {
	wfr, err := m.datastore.LoadWorkflowRunner(uuid)
	if err != nil {
		log.Print(err)
		return
	}

	step, done := wfr.Next()

	if done {
		wfr.End()
	} else {
		compiledStep, err := step.Compile(wfr.Variables)
		if err != nil {
			log.Printf("workflow %s failed to compile step %s: %q",
				wfr.Workflow.Name,
				step.Name,
				err.Error(),
			)
			return
		}

		compiledStep.UUID = wfr.UUID

		j, err := compiledStep.JSON()
		if err != nil {
			log.Print(err)
			return
		}

		if err := node.Producer.send(j); err != nil {
			log.Fatal(err)
		}

		wfr.Last = compiledStep.Name
	}

	m.datastore.DumpWorkflowRunner(wfr)
}
