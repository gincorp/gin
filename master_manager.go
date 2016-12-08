package main

import (
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
    log.Println(body)

    return
}
