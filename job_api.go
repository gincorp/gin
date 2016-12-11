package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (j JobManager) StartAPI() {
	http.HandleFunc("/mon", j.MonRoute)

	http.ListenAndServe(":8000", nil)
}

func (j JobManager) MonRoute(w http.ResponseWriter, r *http.Request) {
	var o interface{}

	status := http.StatusOK

	if r.Method == "GET" {
		o = NewMon()
	} else {
		o = ErrorResponse{fmt.Sprintf("Method %q not allowed", r.Method)}
		status = http.StatusMethodNotAllowed
	}

	json, err := json.Marshal(o)
	if err != nil {
		json = []byte(err.Error())
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, string(json))

}
