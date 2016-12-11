package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/redis.v5"
)

type Datastore struct {
	db *redis.Client
}

func NewDatastore(uri string) (d Datastore, err error) {
	var opts *redis.Options

	if opts, err = redis.ParseURL(uri); err != nil {
		return
	}

	d.db = redis.NewClient(opts)

	_, err = d.db.Ping().Result()
	return
}

func (d Datastore) LoadWorkflow(name string) (wf Workflow, err error) {
	var config string

	if config, err = d.load(wfConfigName(name)); err != nil {
		return
	}

	return ParseWorkflow(config)
}

func (d Datastore) LoadWorkflowRunner(uuid string) (wfr WorkflowRunner, err error) {
	var config string

	if config, err = d.load(wfStateName(uuid)); err != nil {
		return
	}

	return ParseWorkflowRunner(config)
}

func (d Datastore) DumpWorkflowRunner(wfr WorkflowRunner) error {
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
