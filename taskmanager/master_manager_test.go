package taskmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/gincorp/gin/datastore"
	"github.com/gincorp/gin/workflow"
)

type TestDataStore struct{ latest string }

func (ds TestDataStore) LoadWorkflowRunner(u string) (workflow.Runner, error) {
	wf, _ := ds.LoadWorkflow("abc")
	return workflow.Runner{time.Now(), "", ds.latest, time.Now(), "", "efgh-456", map[string]interface{}{}, wf}, nil
}

func (ds TestDataStore) DumpWorkflowRunner(w workflow.Runner) error { return nil }

func (ds TestDataStore) LoadWorkflow(n string) (workflow.Workflow, error) {
	s1 := workflow.Step{Name: "A step", Type: "nil", Context: map[string]string{}}
	s2 := workflow.Step{Name: "Last step", Type: "nil", Context: map[string]string{}}

	return workflow.Workflow{"abc", []workflow.Step{s1, s2}, map[string]string{}}, nil
}

func TestNewMasterManager(t *testing.T) {
	var mm interface{}
	var ds interface{}

	t.Run("Initialises and returns a MasterManager", func(t *testing.T) {
		mm = NewMasterManager("redis://localhost")

		switch mm.(type) {
		case MasterManager:
		default:
			t.Errorf("NewMasterManager() error = Received %T, expected MasterManager", mm)
		}
	})

	t.Run("Initialises with a valid Datastore", func(t *testing.T) {
		ds = mm.(MasterManager).datastore
		switch ds.(type) {
		case datastore.Datastore:
		default:
			t.Errorf("NewMasterManager().datastore error = Received %T, expected datastore.Datasotr", mm.(MasterManager).datastore)
		}
	})
}

func TestMasterManager_Consume(t *testing.T) {
	initWorkflow := `
{
  "InitWorkflow": {
    "Name": "abc",
    "Variables": {}
  },
  "UUID": "efgh-456"
}
`

	startedWorkflowRegister := `
{
  "UUID": "efgh-456",
  "Register": "name",
  "Data": {
    "apple": "red",
    "banana": "yellow"
  },
  "Failed": false,
  "Start": "",
  "End": "",
  "Duration": "1 ms"
}
`

	startedWorkflowNoRegister := `
{
  "UUID": "efgh-456",
  "Data": {
    "pear": "green",
    "grape": "red"
  },
  "Failed": false,
  "Start": "",
  "End": "",
  "Duration": "1 ms"
}
`

	startedWorkflowFailed := `
{
  "UUID": "efgh-456",
  "Data": {
    "orange": "orange",
    "kiwi": "brown"
  },
  "Failed": true,
  "ErrorMessage": "Errors",
  "Start": "",
  "End": "",
  "Duration": "1 ms"
}
`
	ds1 := TestDataStore{}
	ds2 := TestDataStore{"A step"}
	ds3 := TestDataStore{"Last step"}
	tests := []struct {
		name      string
		ds        DataStore
		body      string
		wantErr   bool
		failedWF  bool
		shouldEnd bool
	}{
		// I suspect these tests are going to be really shit
		{"Initialise a workflow", ds1, initWorkflow, false, false, false},
		{"Continue a new wf with a register", ds1, startedWorkflowRegister, false, false, false},
		{"Continue a new wf without a register", ds1, startedWorkflowNoRegister, false, false, false},
		{"Continue a new wf with a failed field", ds1, startedWorkflowFailed, false, true, true},

		{"Continue a started wf with a register", ds2, startedWorkflowRegister, false, false, false},
		{"Continue a started wf without a register", ds2, startedWorkflowNoRegister, false, false, false},
		{"Continue a started wf with a failed field", ds2, startedWorkflowFailed, false, true, true},

		{"Finalise a started wf with a register", ds3, startedWorkflowRegister, false, false, true},
		{"Finalise a started wf without a register", ds3, startedWorkflowNoRegister, false, false, true},
		{"Finalise a started wf with a failed field", ds3, startedWorkflowFailed, false, true, true},
	}
	for _, tt := range tests {
		m := MasterManager{
			datastore: tt.ds,
		}
		gotOutput, err := m.Consume(tt.body)

		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("MasterManager.Consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

		if !tt.failedWF && !tt.shouldEnd {
			// A lot of what consume does will depend massively on the state of a runner. Should the runner have failed,
			// for example, it wont bother returning the next step

			t.Run(fmt.Sprintf("%s - uuid", tt.name), func(t *testing.T) {
				switch gotOutput["UUID"].(type) {
				case string:
				default:
					t.Errorf("MasterManager.Consume() UUID error = received type %T, expected string", gotOutput["UUID"])
					return
				}

				if gotOutput["UUID"].(string) != "efgh-456" {
					t.Errorf("MasterManager.Consume() error = received uuid %V, want 'efgh-456'", gotOutput["UUID"])
					return
				}
			})
		} else {
			t.Run(fmt.Sprintf("%s - empty data", tt.name), func(t *testing.T) {
				if len(gotOutput) > 0 {
					t.Errorf("MasterManager.Consume() should return a map of zero length. Got len = %d", len(gotOutput))
					return
				}
			})
		}
	}
}
