package workflow

import (
    "bytes"
    "encoding/json"
    "text/template"
)

// Step ...
// Step configuration container
type Step struct {
    Context      map[string]string
    Duration     string
    End          string
    ErrorMessage string
    Failed       bool
    Name         string
    Register     string
    Start        string
    Type         string
    UUID         string
}

// Compile ...
// Compile, in place, `Step.Context` entry templates with
// state data from a WorkflowRunner
func (s *Step) Compile(v map[string]interface{}) (err error) {
    for varKey, varValue := range s.Context {
        var buf bytes.Buffer

        tmpl := template.Must(template.New("stepContext").Parse(varValue))

        err = tmpl.Execute(&buf, v)
        if err != nil {
            return
        }

        s.Context[varKey] = buf.String()
    }

    return
}

// SetStatus receives data from job nodes and updates compiled
// Step data within a Workflow Runner for added metadata visibility
func (s *Step) SetStatus(m map[string]interface{}) {
    s.Duration = m["Duration"].(string)
    s.Start = m["Start"].(string)
    s.End = m["End"].(string)
    s.Failed = m["Failed"].(bool)

    if s.Failed {
        s.ErrorMessage = m["ErrorMessage"].(string)
    }
}

// JSON ...
// JSON representation of a `Step`
func (s *Step) JSON() (j []byte, err error) {
    return json.Marshal(s)
}
