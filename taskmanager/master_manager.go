package taskmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gincorp/gin/datastore"
	"github.com/gincorp/gin/workflow"

	"github.com/fatih/structs"
)

// DataStore provides access to workflows, runners and so on
type DataStore interface {
	LoadWorkflowRunner(string) (workflow.Runner, error)
	DumpWorkflowRunner(workflow.Runner) error
	LoadWorkflow(string) (workflow.Workflow, error)
}

// MasterManager is a container for Master Taskmanager configuration
type MasterManager struct {
	datastore DataStore
}

// NewMasterManager returns an initialised Master Taskmanager
func NewMasterManager(redisURI string) (m MasterManager) {
	var err error

	if m.datastore, err = datastore.NewDatastore(redisURI); err != nil {
		log.Fatal(err)
	}

	return
}

// Consume handles json from the message queue; for a Master node these will be responses.
// Parse messages, update Workflow contexts, write to database and call next step
func (m MasterManager) Consume(body string) (output map[string]interface{}, err error) {
	var b interface{}
	var uuid string
	var wfr workflow.Runner

	if err = json.Unmarshal([]byte(body), &b); err != nil {
		return
	}

	input := b.(map[string]interface{})

	if input["InitWorkflow"] != nil {
		req := input["InitWorkflow"].(map[string]interface{})

		uuid, err = m.load(input["UUID"].(string), req["Name"].(string), req["Variables"])
		if err != nil {
			return
		}

		if wfr, err = m.datastore.LoadWorkflowRunner(uuid); err != nil {
			return
		}

	} else {
		uuid = input["UUID"].(string)
		if wfr, err = m.datastore.LoadWorkflowRunner(uuid); err != nil {
			return
		}

		idx, step := wfr.Current()

		step.SetStatus(input)
		wfr.Workflow.Steps[idx] = step

		switch input["Register"].(type) {
		case string:
			register := input["Register"].(string)

			switch input["Data"].(type) {
			case map[string]interface{}:
				data := input["Data"].(map[string]interface{})
				wfr.Variables[register] = data

			default:
				log.Println("Not registering input: got garbage back")
			}
		}

		if input["Failed"].(bool) {
			wfr.Fail(fmt.Sprintf("Step %q failed. See below", wfr.Last))
			m.datastore.DumpWorkflowRunner(wfr)
			return
		}

	}

	m.datastore.DumpWorkflowRunner(wfr)

	s, done := m.proceed(wfr.UUID)
	if !done {
		output = structs.Map(s)
	}

	return
}

// load a workflow from storage and create a WorkflowRunner state machine
func (m MasterManager) load(u, name string, variables interface{}) (uuid string, err error) {
	wf, err := m.datastore.LoadWorkflow(name)
	if err != nil {
		return
	}

	wfr := workflow.NewRunner(u, wf)

	switch variables.(type) {
	case map[string]interface{}:
		wfr.Variables["Runtime"] = variables
	}

	wfr.Start()

	m.datastore.DumpWorkflowRunner(wfr)

	return wfr.UUID, nil
}

// proceed will, should there be a next step in the workflow, compile step templates
// and push the step to the emssage queue. Otherwise it'll set the relevant Workflow Runner
// to completed.
func (m MasterManager) proceed(uuid string) (step workflow.Step, done bool) {
	wfr, err := m.datastore.LoadWorkflowRunner(uuid)
	if err != nil {
		log.Print(err)
		return
	}

	step, done = wfr.Next()

	if done {
		wfr.End()
	} else {
		err := step.Compile(wfr.Variables)
		if err != nil {
			wfr.Fail(fmtError(step, err))
			m.datastore.DumpWorkflowRunner(wfr)

			return
		}

		step.UUID = wfr.UUID

		wfr.Last = step.Name
	}

	m.datastore.DumpWorkflowRunner(wfr)
	return step, done
}

func fmtError(step workflow.Step, err error) string {
	return fmt.Sprintf("%s: %s", step.Name, err.Error())
}
