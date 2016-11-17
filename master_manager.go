package main

import (
    "log"
)

type MasterManager struct {

}

func NewMasterManager() (m MasterManager) {
    return
}

func (m MasterManager) Consume(body string) (output map[string]interface{}, err error) {
    log.Println(body)

    return
}
