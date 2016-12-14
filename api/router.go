package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

func (a API) monRoute(w http.ResponseWriter, r *http.Request) {
	var o interface{}

	status := http.StatusOK

	if r.Method == "GET" {
		o = NewMon()
	} else {
		o = ErrorResponse{fmt.Sprintf("Method %q not allowed", r.Method)}
		status = http.StatusMethodNotAllowed
	}

	j, err := json.Marshal(o)
	if err != nil {
		j = []byte(err.Error())
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, string(j))
}

func (a API) wfRoute(w http.ResponseWriter, r *http.Request) {
	var o interface{}

	status := http.StatusOK

	if r.Method == "GET" {
		uuid := r.URL.Path[len("/wf/"):]

		o = a.getWF(uuid)
		switch o.(type) {
		case ErrorResponse:
			status = http.StatusBadRequest
		}
	} else if r.Method == "POST" {
		defer r.Body.Close()
		sr := StarterRequest{}
		json.NewDecoder(r.Body).Decode(&sr)

		o = a.startWF(sr)
		switch o.(type) {
		case ErrorResponse:
			status = http.StatusBadRequest
		}

		log.Println(o)
	} else {
		o = ErrorResponse{fmt.Sprintf("Method %q not allowed", r.Method)}
		status = http.StatusMethodNotAllowed
	}

	j, err := json.Marshal(o)
	if err != nil {
		j = []byte(err.Error())
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, string(j))
}

func (a API) getWF(uuid string) (wf interface{}) {
	wfr, err := a.datastore.LoadWorkflowRunner(uuid)

	if err != nil {
		return ErrorResponse{err.Error()}
	}
	return wfr
}

func (a API) startWF(s StarterRequest) (sr interface{}) {
	uuid := uuid.NewV4().String()
	o := StartWorkflow{s, time.Now(), uuid}

	j, err := json.Marshal(o)
	if err != nil {
		return ErrorResponse{err.Error()}
	}

	err = a.producer.Send(j)
	if err != nil {
		return ErrorResponse{err.Error()}
	}

	return StarterResponse{uuid}
}
