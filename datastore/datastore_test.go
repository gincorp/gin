package datastore

import (
	"testing"
	"time"

	"github.com/gincorp/gin/workflow"
	"gopkg.in/redis.v5"
)

type TestConnector struct{}

func (tc TestConnector) Get(s string) *redis.StringCmd { return &redis.StringCmd{} }
func (tc TestConnector) Set(s string, d interface{}, t time.Duration) *redis.StatusCmd {
	return &redis.StatusCmd{}
}
func (tc TestConnector) Ping() *redis.StatusCmd { return &redis.StatusCmd{} }

func TestNewDatastore(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"initialises with a valid redis uri", args{uri: "redis://localhost"}, false},
		{"errors with an invalid redis uri", args{uri: "some-nonsense"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatastore(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatastore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDatastore_LoadWorkflow(t *testing.T) {
	var wf interface{}
	var err error

	db := &TestConnector{}
	d := &Datastore{db: db}

	t.Run("Return a workflow from datastore", func(t *testing.T) {
		wf, err = d.LoadWorkflow("workflow_name")
		if err != nil {
			t.Errorf("datastore.LoadWorkflow() error = %v", err)
		}

		switch wf.(type) {
		case workflow.Workflow:
		default:
			t.Errorf("datastore.LoadWorkflow() returned %T, expected workflow.Workflow", wf)
		}
	})
}

func TestDatastore_SaveWorkflow(t *testing.T) {
	type fields struct {
		db connector
	}
	type args struct {
		w         workflow.Workflow
		overwrite bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"saves a workflow", fields{TestConnector{}}, args{w: workflow.Workflow{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Datastore{
				db: tt.fields.db,
			}
			if err := d.SaveWorkflow(tt.args.w, tt.args.overwrite); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.SaveWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatastore_LoadWorkflowRunner(t *testing.T) {
	type fields struct {
		db connector
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"loads a workflow", fields{TestConnector{}}, args{uuid: "abcd-1234-foo"}, false},
	}
	for _, tt := range tests {
		var wfr interface{}
		var err error

		t.Run(tt.name, func(t *testing.T) {
			d := Datastore{
				db: tt.fields.db,
			}
			wfr, err = d.LoadWorkflowRunner(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Datastore.LoadWorkflowRunner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			switch wfr.(type) {
			case workflow.Runner:
			default:
				t.Errorf("Datastore.LoadWorkflowRunner() error Expected a workflow.Runner, received %T", wfr)
			}
		})
	}
}

func TestDatastore_DumpWorkflowRunner(t *testing.T) {
	type fields struct {
		db connector
	}
	type args struct {
		wfr workflow.Runner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"saves a workflow runner", fields{TestConnector{}}, args{workflow.Runner{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Datastore{
				db: tt.fields.db,
			}
			if err := d.DumpWorkflowRunner(tt.args.wfr); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.DumpWorkflowRunner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
