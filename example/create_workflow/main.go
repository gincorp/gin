// Copyright 2016, gincorp.
//
// This project is under the MIT licence;
// found in LICENCE.md in the parent directory.
//
// This script is expected to be run via go run wf.go
// It configures a simple workflow designed to be run with the
// job manager implemenetation in the parent directory.
//
//
// This workflow will compile certain interesting bits of info:
// weather forecasts, news headlines, exchange rates. It'll then
// email this information to a comfigured email address.
//
// ENV Vars:
// SENDER_ADDRESS=''     - the sender address on an email
// RECEIPIENT_ADDRESS='' - the receipient of the email

package main

import (
    "flag"
    "log"
    "os"
    "strings"

    "github.com/gincorp/gin/datastore"
    "github.com/gincorp/gin/workflow"
)

var (
    redisURI *string
)

func init() {
    redisURI = flag.String("redis", "redis://localhost:6379/0", "URI of redis node")
    flag.Parse()
}

func main() {
    d, err := datastore.NewDatastore(*redisURI)
    if err != nil {
        panic(err)
    }

    log.Print(d.SaveWorkflow(sendDailyEmail(), true))
}

func sendDailyEmail() workflow.Workflow {
    // Grab some information about some stuff, email it

    vars := make(map[string]string)
    vars["mail_host"] = "smtp.gmail.com"
    vars["mail_port"] = "587"
    vars["mail_from"] = os.Getenv("SENDER_ADDRESS")
    vars["mail_to"] = os.Getenv("RECEIPIENT_ADDRESS")

    emailBody := []string{
        "Greetings,",
        "",
        "The markets:",
        "{{ range .finance.prices }} {{ .Name }}: {{ .Price }} \r\n{{ end }}",
        "",
        "Regards,",
        "gin",
    }

    return workflow.Workflow{
        Name:      "Send daily email",
        Variables: vars,
        Steps: []workflow.Step{
            workflow.Step{
                Name:     "Get Financial Information",
                Type:     "get-financial",
                Register: "finance",
            },
            workflow.Step{
                Name: "Send Email",
                Type: "send-email",
                Context: map[string]string{
                    "host": "{{ .Defaults.mail_host }}",
                    "port": "{{ .Defaults.mail_port }}",

                    "from":    "{{ .Defaults.mail_from }}",
                    "to":      "{{ .Runtime.mail_to }}",
                    "subject": "Daily Update Email",
                    "body":    strings.Join(emailBody[:], "\r\n"),
                },
            },
        },
    }
}
