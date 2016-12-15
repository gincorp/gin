package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gincorp/gin/workflow"

	"gopkg.in/redis.v5"
)

// Datastore handles connections to and from datastores
type Datastore struct {
	db *redis.Client
}

// NewDatastore will create and test a connection into storage
func NewDatastore(uri string) (d Datastore, err error) {
	var opts *redis.Options

	if opts, err = redis.ParseURL(uri); err != nil {
		return
	}

	d.db = redis.NewClient(opts)

	_, err = d.db.Ping().Result()
	return
}

// LoadWorkflow returns a Workflow object from a workflow name
func (d Datastore) LoadWorkflow(name string) (wf workflow.Workflow, err error) {
	var config string

	if config, err = d.load(wfConfigName(name)); err != nil {
		return
	}

	return workflow.ParseWorkflow(config)
}

// SaveWorkflow saves a new workflow or overwrites an existing (if overwrite=true)
func (d Datastore) SaveWorkflow(w workflow.Workflow, overwrite bool) error {
	j, err := json.Marshal(w)
	if err != nil {
		return err
	}

	s, err := d.load(wfConfigName(w.Name))

	if err != redis.Nil && err != nil {
		return err
	}

	if s != "" && overwrite == false {
		return errors.New(fmt.Sprintf("Refusing to overwrite workflow %q", w.Name))
	}

	return d.db.Set(wfConfigName(w.Name), j, 0).Err()
}

// LoadWorkflowRunner returns a WorkflowRunner; a parsed and compiled workflow
// with a simple state machine
func (d Datastore) LoadWorkflowRunner(uuid string) (wfr workflow.Runner, err error) {
	var config string

	if config, err = d.load(wfStateName(uuid)); err != nil {
		return
	}

	return workflow.ParseRunner(config)
}

// DumpWorkflowRunner dumps a running `WorkflowRunner`'s state to storage
func (d Datastore) DumpWorkflowRunner(wfr workflow.Runner) error {
	j, err := json.Marshal(wfr)
	if err != nil {
		return err
	}

	return d.db.Set(wfStateName(wfr.UUID), j, 0).Err()
}

func (d Datastore) load(key string) (string, error) {
	return d.db.Get(key).Result()
}

func normaliseName(wfName string) string {
	return strings.Replace(wfName, " ", "_", -1)
}

func wfConfigName(wfName string) string {
	return fmt.Sprintf("workflow.%s", normaliseName(wfName))
}

func wfStateName(uuid string) string {
	return fmt.Sprintf("state.%s", uuid)
}
