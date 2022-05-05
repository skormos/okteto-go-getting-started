# Sample Okteto Project

This serves as an example of how to leverage Okteto to build a web-based service based in Kubernetes.

## Endpoints

### API
The endpoints defined under `/api` is defined via OpenAPI spec, located in this project at `<ROOT>/api/sample.openapi3.json`. This can be rendered using [Swagger's editor](https://editor.swagger.io). 

### Ping
The `/ping` endpoint was added to test [Prometheus integration](https://prometheus.io/docs/tutorials/instrumenting_http_server_in_go/), and to refactor the handlers to support the generated API and a new endpoint at the same level at `/api`.

### Metrics
The `/metrics` endpoint supports Prometheus.

## Make

### General Notes
The `make` targets are patterned for common tasks for a Go project. Most are rather self-explanatory, and should work locally or on an Okteto development environment. 

### Specific Targets

#### Generating APIs
Any changes to the OpenAPI spec definition, requires that `make gen` is run so the API interface is updated.

#### Deploying the Docker Image
Once a change is made to the functionality in the server code, the Docker image needs to be updated and deployed. This can be done with `make build-image`.

# Getting Started on Okteto with Go

[![Develop on Okteto](https://okteto.com/develop-okteto.svg)](https://cloud.okteto.com/deploy?repository=https://github.com/okteto/go-getting-started)

This example shows how to use [Okteto](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using Kubernetes manifests.

This is the application used for the [How to Develop and Debug Go Applications in Kubernetes](https://okteto.com/blog/how-to-develop-go-apps-in-kubernetes/) blog post.
