[![Build Status](https://travis-ci.org/gincorp/gin.svg?branch=master)](https://travis-ci.org/gincorp/gin)
[![GoDoc](https://godoc.org/github.com/gincorp/gin?status.svg)](https://godoc.org/github.com/gincorp/gin)
[![Go Report Card](https://goreportcard.com/badge/github.com/gincorp/gin)](https://goreportcard.com/report/github.com/gincorp/gin)


![gin gopher](doc/gopher-text_small.png)

gin is a simple distributed workflow engine, including task brokerage and statefulness, written in go.

This document provides some contextual data for the project. Please see [Getting Started](GETTING_STARTED.md) for a deeper dive.

```bash
$ go get -u github.com/gincorp/gin
$ # master node
$ gin -mode master -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0
$ # job node
$ gin -mode job -amqp amqp://guest:guest@localhost/vhost
$ # api node
$ gin -mode api -amqp amqp://guest:guest@localhost/vhost -redis redis://localhost:6379/0 -host 0.0.0.0 -port 8080
```

## Architecture

![naive gin architecture diagram](doc/workflow-engine.png)

A gin cluster contains a number of components; two of which are necessary:

* One or more master components
* One or more job components

There may also be, and practically always are, at least one API component too.

These components are called 'Nodes'; further information of which may be found in the wiki.


### Third Party Dependencies

There are also a couple of dependencies

* RabbitMQ
* Redis


### Terminology

A gin cluster runs one or more **workflows**.

A workflow contains a name, some optional **default variables** and a set of **tasks**.

A set of variables is a key value pair of `string: string`.

A task contains a **job function type**, a name and **contextual data**

A set of contextual data is much the same as default variables, but the values may be golang `text/template` data.


#### A Workflow Definition

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

### Valuable information

| who       | what |
|-----------|------|
| dockerhub | https://hub.docker.com/r/gincorp/gin/   |
| circleci  | https://circleci.com/gh/gincorp/gin   |
| licence   | MIT   |


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
