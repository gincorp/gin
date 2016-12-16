gin
==

gin is a simple distributed workflow engine, including task brokerage and staefulness, written in go.

It contains three components; of which there may be a multitude of instances:

* One or more master Nodes
* One or more job Nodes
* One or more API Nodes

And a couple of dependencies:

* RabbitMQ
* Redis

A `workflow` is made up of `tasks` and `variables` and collated in a `runner` which stores state via a simple state machine.

Master and Job nodes may be run right out of the box. Job nodes, however, come with very few configured jobs out the box and so it becomes simpler to create a project which initialises a JobManager object.

| who       | what |
|-----------|------|
| dockerhub | https://hub.docker.com/r/gincorp/gin/   |
| circleci  | https://circleci.com/gh/gincorp/gin   |
| licence   | MIT   |



Master Nodes
--

Master nodes are run with:

```
$ gin -mode master -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0
```

Or via docker:

```
$ docker run -p8080:8080 jspc/gin -mode master -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0
```

Master nodes compile and broker workflow tasks via a runner.

Job Nodes
--

Job nodes are run with:

```
$ gin -mode job -amqp amqp://guest:guest@localhost/vhost
```

Or via docker:

```
$ docker run -p8080:8080 jspc/gin -mode job -amqp amqp://guest:guest@localhost/vhost
```

Job nodes receive formatted tasks via rabbitmq, execute the task and return outputs and metadata.

API Nodes
--

API nodes are run with:

```
$ gin -mode job -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0 -host 0.0.0.0 -port 8080
```

Or via docker:

```
$ docker run -p8080:8080 jspc/gin -mode job -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/ -host 0.0.0.0 -port 8080
```

API nodes provide an interface to the engine; for starting and configuring workflows.

**Note**: There is no auth-n, auth-z nor much of anything in front of the API. You should either use an API gateway or router in front of the tool. It probably wont be included in this project; it seems out of scope.


Architecture
--

![naive gin architecture diagram](img/workflow-engine.png)


Workflows
--

A workflow definition looks like:

```json
{
    "name": "do stuff",
    "variables": {
        "echo_url": "http://172.17.0.1:8000/some-endpoint",
        "content_type": "application/data"
    },
    "steps": [{
        "name": "Test Stuff",
        "type": "post-to-web",
        "context": {
            "url": "{{.Defaults.echo_url}}",
            "content-type": "{{.Defaults.content_type}}",
            "data": "foo=bar"
        },
        "register": "echo_data"
    }, {
        "name": "Templatery Stuff",
        "type": "log",
        "context": {
            "message": "{{.echo_data.Host}}"
        },
        "register": "something"
    }]
}
```

The initial set of variables are accessible via the `.Defaults` key in context templates. These templates allow:

* DRYing up of step data
* The ability to rely on the output of a previous step

To access the output of a previous step, one must `register` the step with a unique name. This will be recognisable by ansible users, for example (in fact: thats where I stole the idea).

A workflow is kicked off via an api call:

```bash
curl -X POST -d '{"Name": "nontrivial workflow", "Variables":{"foo": "bar"} }' localhost:8080/wf/
```

The object `Variables` may be omitted: this is useful for data which is not known until runtime, and is made availabled via the `.Rumtime` key in context templates; in an identical fashion to how `.Defaults` works.

The request mints, and returns, a UUID. A workflow runner is attached to this uuid. This uuid can be used to view the status of a workflow; like the follow failed workflow:

```bash
$ curl localhost:8080/wf/7423bbb5-b020-4410-91b6-a5e2754e6a47
```

Returning:

```json
{
  "EndTime": "2016-12-14T11:40:54.367619801Z",
  "ErrorMessage": "Step \"Test Stuff\" failed. See below",
  "Last": "Test Stuff",
  "StartTime": "2016-12-14T11:34:32.703863903Z",
  "State": "failed",
  "UUID": "7423bbb5-b020-4410-91b6-a5e2754e6a47",
  "Variables": {
    "Defaults": {
      "content_type": "application/data",
      "echo_url": "http://172.17.0.1:8000/some-endpoint"
    }
  },
  "Workflow": {
    "Name": "do stuff",
    "Steps": [
      {
        "Context": {
          "content-type": "application/data",
          "data": "foo=bar",
          "url": "http://172.17.0.1:8000/some-endpoint"
        },
        "Duration": "2 ms",
        "End": "2016-12-14T11:36:06",
        "ErrorMessage": "Post http://172.17.0.1:8000/some-endpoint: dial tcp 172.17.0.1:8000: getsockopt: connection refused",
        "Failed": true,
        "Name": "Test Stuff",
        "Register": "echo_data",
        "Start": "2016-12-14T11:36:06",
        "Type": "post-to-web",
        "UUID": ""
      },
      {
        "Context": {
          "message": "{{.echo_data.Host}}"
        },
        "Duration": "",
        "End": "",
        "ErrorMessage": "",
        "Failed": false,
        "Name": "Templatery Stuff",
        "Register": "something",
        "Start": "",
        "Type": "log",
        "UUID": ""
      }
    ],
    "Variables": {
      "content_type": "application/data",
      "echo_url": "http://172.17.0.1:8000/some-endpoint"
    }
  }
}
```

Here we see the first step has been compiled, yet the second step has npt. This is because a step is compiled right before it is run, to allow it access to the latest data. We can also see the runner has failed; it could not get access to a service and thus bombs out.


Licence
--

MIT License

Copyright (c) 2016 jspc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
