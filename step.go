package main

import (
    "bytes"
    "encoding/json"
    "text/template"
)

type Step struct {
    Context map[string]string
    Name string
    Register string
    Type string
    UUID string
}

func (s *Step)Compile(v map[string]interface{}) (*Step, error) {
    var err error

    for varKey,varValue := range s.Context {
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

func (s *Step)Json() (j []byte, err error) {
    return json.Marshal(s)
}
