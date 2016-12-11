package main

import (
	"bytes"
	"encoding/json"
	"text/template"
)

// Step ...
// Step configuration container
type Step struct {
	Context  map[string]string
	Name     string
	Register string
	Type     string
	UUID     string
}

// Compile ...
// Compile, in place, `Step.Context` entry templates with
// state data from a WorkflowRunner
func (s *Step) Compile(v map[string]interface{}) (*Step, error) {
	var err error

	for varKey, varValue := range s.Context {
		var buf bytes.Buffer

		tmpl := template.Must(template.New("stepContext").Parse(varValue))

		err = tmpl.Execute(&buf, v)
		if err != nil {
			return s, err
		}

		s.Context[varKey] = buf.String()
	}

	return s, nil
}

// JSON ...
// JSON representation of a `Step`
func (s *Step) JSON() (j []byte, err error) {
	return json.Marshal(s)
}
