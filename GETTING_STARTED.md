Getting Started
==

This document will configure a simple gin cluster, with a workflow and custom jobs, via dependencies running via containers.

Pre-requisities
--

This guide requires redis and rabbit amqp. To get these going via docker, lazily a direct mapping of ports, the following commands may be run:

```bash
$ docker run -p 5672:5672 -p 15672:15672 -d --hostname some-rabbit --name some-rabbit rabbitmq:3-management
$ docker run -p 6379 --name wf-redis redis
```

Configuring a Workflow
--

In the directory [`example/create_workflow`](examples/workflow) there is a sample configuration script. This script will generate some json and stick it in our redis server. This json, though, can be generated any way- it just has to be valid and in redis under the correct key.

This workflow definition will run a job function called `get-financial` and store the output of this under the key `finance`. It will then email this information using the template:

```golang
emailBody := []string{
    "Greetings,",
    "",
    "The markets:",
    "{{ range .finance.prices }} {{ .Name }}: {{ .Price }} \r\n{{ end }}",
    "",
    "Regards,",
    "gin",
}
```

(Which is joined in the `send-email step`)

### Steps

(**Note:** This workflow assumes an email will be sent via gmail. The default values, specifically lines 53-54, must be changed prior to the below to use something else.)

1. Build the tool: `go build`
1. Run the tool, linking with the redis we spun up: `./create-workflow -redis redis://localhost`
1. Verify the workflow:
    1. Login to redis: `docker run -ti --link some-redis:some-redis redis redis-cli -h some-redis`
    1. Validate key `"workflow.Send_daily_email"` exists by running: `keys *`
    1. Validate there is a load of json there by running `get  "workflow.Send_daily_email"`


Running Workflows
--

To run a workflow we'll need to run:

1. A Master node
1. A Job node with access to the job functions we've put in our workflow steps; and
1. An API node to kick off and monitor out builds

The master node that ships with gin is fine for our purposes; a master node is merely a workflow broker and state manager. This is very generic and very complete. Ditto the API node; the default api is a very simple beast.

Job nodes, though, are different; the default job functions in gin are very basic and sparse. Out of the box we provide:

1. `post-to-web` - post some data to an endpoint somewhere
1. `get-from-web` - get some data from an endpoint somewhere
1. `log` - dump some data into the job node's log

In [`example/jobs`](examples/jobs) there is a simple implementation of how one would integrate jobs. In `main.go` we initialise a new job node:

```golang
n := node.NewNode(*amqpURI, "", "job")
```

We assign some functions:

```golang
jobManager := taskmanager.NewJobManager()
jobManager.AddJob("get-financial", getCurrencyPrices)
jobManager.AddJob("send-email", sendEmail)

n.TaskManager = jobManager
```

And then we start the rabbitmq consumer:

```golang
log.Fatal(n.ConsumerLoop())
```

The functions `get-financial` and `send-email` may be found in `functions.go`. The signature for these functions may be found in `taskmanager/job_manager.go`:

```golang
func(JobNotification) (map[string]interface{}, error)
```

Which is to say: a valid function receives a `workflow.JobNotification` and returns some mapped data of `string` to anything, and an error.

### Running

(**Note:** The following tasks will run in the foreground. You may either use job control to send them to the background, run them in different terminals/tmux windows, whatever)

#### Master Node:

`gin -mode job -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0`

#### API Node:

`gin -mode api -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0 -host 0.0.0.0 -port 8080`

#### Job Node:

(**Note**: This step expects to send email, and will require a username and password. It is all geared up to send via gmail. You'll need to have changed the workflow, as per the instruction further up, should you be using something else)

1. Navigate to `examples/jobs`
1. Build example code: `go build`
1. Set some environment variables (See Lines: 73-74, we use these to avoid storing login data)
   1. `export MAIL_USERNAME=foo@gmail.com`
   1. `export MAIL_PASSWORD=secure`
1. Run the node: `./jobs -amqp amqp://guest:guest@localhost/vhost`

### Starting the workflow

A workflow is initialised and run via a crafted message to the master node. We can do this via the api:

```bash
$ curl -X POST -d '{"Name": "Send daily email", "Variables": {"mail_to": "something@example.com"} }' localhost:8080/wf/
{"UUID": "5dba2200-c7f2-4bd7-bda5-35ad3ff6eed5"}
```

The `Variables` object is completely optional; though the workflow we've created uses it. The above key is available in only the running workflow as:

```golang
"{{ .Runtime.mail_to }}"
```

Which complements `{{.Defaults}}` which are configured when they're created.

### Viewing/ Monitoring a workflow

```bash
$ curl localhost:8080/wf/5dba2200-c7f2-4bd7-bda5-35ad3ff6eed5
```

```json
{
    "EndTime": "2016-12-17T15:17:56.451905674Z",
    "ErrorMessage": "",
    "Last": "Send Email",
    "StartTime": "2016-12-17T15:16:22.66100022Z",
    "State": "ended",
    "UUID": "5dba2200-c7f2-4bd7-bda5-35ad3ff6eed5",
    "Variables": {
        "": {},
        "Defaults": {
            "mail_from": "foo@example.com",
            "mail_host": "smtp.gmail.com",
            "mail_port": "587",
            "mail_to": "foo@example.com"
        },
        "finance": {
            "prices": [{
                "Name": "$/€",
                "Price": "0.956600"
            }, {
                "Name": "$/£",
                "Price": "0.800740"
            }, {
                "Name": "Gold per ounce",
                "Price": "0.000756"
            }]
        }
    },
    "Workflow": {
        "Name": "Send daily email",
        "Steps": [{
            "Context": null,
            "Duration": "81 ms",
            "End": "2016-12-17T15:17:54",
            "ErrorMessage": "",
            "Failed": false,
            "Name": "Get Financial Information",
            "Register": "finance",
            "Start": "2016-12-17T15:17:54",
            "Type": "get-financial",
            "UUID": ""
        }, {
            "Context": {
                "body": "Greetings,\r\n\r\nThe markets:\r\n $/€: 0.956600 \r\n $/£: 0.800740 \r\n Gold per ounce: 0.000756 \r\n\r\n\r\nRegards,\r\ngin",
                "from": "foo@example.com",
                "host": "smtp.gmail.com",
                "port": "587",
                "subject": "Daily Update Email",
                "to": "foo@example.com"
            },
            "Duration": "1855 ms",
            "End": "2016-12-17T15:17:56",
            "ErrorMessage": "",
            "Failed": false,
            "Name": "Send Email",
            "Register": "",
            "Start": "2016-12-17T15:17:54",
            "Type": "send-email",
            "UUID": ""
        }],
        "Variables": {
            "mail_from": "foo@example.com",
            "mail_host": "smtp.gmail.com",
            "mail_port": "587",
            "mail_to": "foo@example.com"
        }
    }
}
```

This JSON is a copy of the initial workflow definition as created earlier, with the addition of contextual and state data; including compiled step data.
