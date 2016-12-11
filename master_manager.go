package main

import (
    "encoding/json"
    "log"
)

type MasterManager struct {
    Datastore Datastore
}

func NewMasterManager() (m MasterManager) {
    var err error

    if m.Datastore, err = NewDatastore(*redisUri); err != nil {
        log.Fatal(err)
    }

    return
}

func (m MasterManager) Consume(body string) (output map[string]interface{}, err error) {
    // Parse body into some object
    // Lookup workflow runner by parsed body's UUID
    // if object.Register != "" then add to wfr.Variables[object.Register]
    // dump back to datastore
    // Call m.continue()

    var b interface{}
    var wfr WorkflowRunner

    if err = json.Unmarshal([]byte(body), &b); err != nil {
        return
    }

    output = b.(map[string]interface{})
    uuid := output["UUID"].(string)
    if wfr, err = m.Datastore.LoadWorkflowRunner(uuid); err != nil {
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

    m.Datastore.DumpWorkflowRunner(wfr)
    m.Continue(wfr.UUID)

    return
}

func (m MasterManager) Load(name string) (uuid string, err error){
    wf, err := m.Datastore.LoadWorkflow(name)
    if err != nil {
        return
    }

    wfr := NewWorkflowRunner(wf)
    wfr.Start()

    m.Datastore.DumpWorkflowRunner(wfr)

    return wfr.UUID, nil
}

func (m MasterManager) Continue(uuid string) {
    wfr, err := m.Datastore.LoadWorkflowRunner(uuid)
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

        j, err := compiledStep.Json()
        if err != nil {
            log.Print(err)
            return
        }

        if err := node.Producer.send(j); err != nil {
            log.Fatal(err)
        }

        wfr.Last = compiledStep.Name
        m.Datastore.DumpWorkflowRunner(wfr)
    }
}
