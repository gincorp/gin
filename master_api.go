package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StarterRequest ...
// Placeholder for incoming 'start workflow' requests
type StarterRequest struct {
	Name string
}

// StarterResponse ...
// Placeholder for outgoing 'start workflow' responses
type StarterResponse struct {
	UUID string
}

// ErrorResponse ...
// Placeholder for outgoing errors in responses
type ErrorResponse struct {
	Message string
}

// StartAPI ...
// Listen on 0.0.0.0:8080 to requests for master nodes
// Used for node metadata, monitoring, and
// starting and looking up workflows
func (m MasterManager) StartAPI() {
	http.HandleFunc("/mon/", m.monRoute)
	http.HandleFunc("/wf/", m.wfRoute)

	http.ListenAndServe(":8080", nil)
}

func (m MasterManager) monRoute(w http.ResponseWriter, r *http.Request) {
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

func (m MasterManager) wfRoute(w http.ResponseWriter, r *http.Request) {
	var o interface{}

	status := http.StatusOK

	if r.Method == "GET" {
		uuid := r.URL.Path[len("/wf/"):]

		o = m.getWF(uuid)
		switch o.(type) {
		case ErrorResponse:
			status = http.StatusBadRequest
		}
	} else if r.Method == "POST" {
		defer r.Body.Close()
		sr := StarterRequest{}
		json.NewDecoder(r.Body).Decode(&sr)

		o = m.startWF(sr)
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

func (m MasterManager) getWF(uuid string) (wf interface{}) {
	wfr, err := m.datastore.LoadWorkflowRunner(uuid)

	if err != nil {
		return ErrorResponse{err.Error()}
	}
	return wfr
}

func (m MasterManager) startWF(s StarterRequest) (sr interface{}) {
	uuid, err := m.Load(s.Name)

	if err != nil {
		sr = ErrorResponse{err.Error()}
	} else {
		sr = StarterResponse{uuid}
		m.Continue(uuid) // Start first step
	}

	return
}
