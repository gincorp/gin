package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type StarterRequest struct {
    Name string
}

type StarterResponse struct {
    UUID string
}

type ErrorResponse struct {
    Message string
}

func (m MasterManager)StartAPI() {
    http.HandleFunc("/mon/", m.MonRoute)
    http.HandleFunc("/wf/", m.WFRoute)

    http.ListenAndServe(":8080", nil)
}

func (m MasterManager)MonRoute(w http.ResponseWriter, r *http.Request) {
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
        j =[]byte( err.Error() )
    }

    w.WriteHeader(status)
    fmt.Fprintf(w, string(j))
}

func (m MasterManager)WFRoute(w http.ResponseWriter, r *http.Request) {
    var o interface{}

    status := http.StatusOK

    if r.Method == "GET" {
        uuid := r.URL.Path[len("/wf/"):]

        o = m.GetWF(uuid)
        switch o.(type) {
        case ErrorResponse:
            status = http.StatusBadRequest
        }
    } else if r.Method == "POST" {
        defer r.Body.Close()
        sr := StarterRequest{}
        json.NewDecoder(r.Body).Decode(&sr)

        o = m.StartWF(sr)
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

func (m MasterManager)StartWF(s StarterRequest) (sr interface{}) {
    uuid, err := m.Load(s.Name)

    if err != nil {
        sr = ErrorResponse{err.Error()}
    } else {
        sr = StarterResponse{uuid}
        m.Continue(uuid)    // Start first step
    }

    return
}

func (m MasterManager)GetWF(uuid string) (wf interface{}) {
    wfr, err := m.Datastore.LoadWorkflowRunner(uuid)

    if err != nil {
        return ErrorResponse{err.Error()}
    }
    return wfr
}
