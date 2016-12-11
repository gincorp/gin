package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StartAPI ...
// Listen on 0.0.0.0:8000 to requests for job nodes
// Used for node metadata and monitoring
func (j JobManager) StartAPI() {
	http.HandleFunc("/mon", j.monRoute)

	http.ListenAndServe(":8000", nil)
}

func (j JobManager) monRoute(w http.ResponseWriter, r *http.Request) {
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
