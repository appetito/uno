# Uno - Microservices Framework for NATS

Uno is a lightweight, efficient microservices framework built on top of NATS, designed to simplify the development of distributed applications. It provides a clean and intuitive API for building scalable microservices with built-in monitoring, tracing, and observability features.

Uno is result of rethink and rework of NATS `micro` framework, major additions are:

- [x] Typed endpoints (automatic marshaling/unmarshaling)
- [x] Interceptors
- [x] Auto generated client
- [x] Scaffolding (code generator)

## Features

- 🚀 Easy-to-use API for building NATS-based microservices
- 📊 Built-in Prometheus integration
- 🔍 OpenTelemetry integration for distributed tracing
- 🎯 Flexible endpoint routing and grouping
- 🔄 Automatic request handling and response management
- 🛠 Scaffolding and Client generation based on YAML service definition
- 🔌 Extensible through interceptors
- 🕒 Context support (deadlines)

## Project status

Uno is currently in the alpha stage and is still under development. We welcome contributions and feedback to improve the framework.

## Installation

```bash
go get github.com/appetito/uno
```

## Quick Start

```go
package main

import (
    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"
)

func main() {
    // Connect to NATS
    nc, _ := nats.Connect(nats.DefaultURL)
    
    // Create a new service
    svc, _ := uno.AddService(nc, uno.Config{
        Name:    "my-service",
        Version: "1.0.0",
    })

    // Add an endpoint
    svc.AddEndpoint("hello", uno.HandlerFunc(func(req uno.Request) {
        req.Respond([]byte("Hello, World!"))
    }))

    // Start the service
    svc.ServeForever()
}
```

Run and call service with nats CLI:

```bash
nats req "hello" 'Hey'
```

## Features

### Endpoint Groups

```go
group := svc.AddGroup("users")
group.AddEndpoint("add", uno.HandlerFunc(AddUserHandler))
group.AddEndpoint("list", uno.HandlerFunc(ListUsersHandler))
```

### Request Handling
Type-safe request handling with automatic marshaling/unmarshaling:


```go

type User struct {
    ID    string	`json:"id"`
    Name  string 	`json:"name"`
    Email string    `json:"email"`
}

type CreateUserRequest struct {
    Name  string	`json:"name"`
    Email string	`json:"email"`
}


svc.AddEndpoint("create-user", uno.AsStructHandler[CreateUserRequest](func(req uno.Request, request CreateUserRequest) {
    // Handle user creation
    user := User{
        ID:    uuid.New().String(),
        Name:  request.Name,
        Email: request.Email,
    }
    req.RespondJSON(user)
}))
```


### Interceptors (Middleware)

Interceptors in Uno work similarly to middleware in other frameworks, allowing you to add cross-cutting concerns to your service endpoints. They can modify or enhance the request/response flow, add logging, handle errors, or perform any other operations before or after request processing.

#### Built-in Interceptors

### 

### Monitoring

Built-in support for service statistics and Prometheus metrics.

### Tracing

Built-in opentelemetry integration.


### Code Generation

Uno includes powerful code generation capabilities that help you quickly scaffold services and generate type-safe clients. The framework can generate service code from YAML definitions.

#### Service Definition

Define your service in uno.yaml YAML file:

```yaml
namespace: myapp
name: user-service
package: github.com/myorg/userservice

types:
  - name: User
    fields:
      - name: ID
        type: string
      - name: Name
        type: string
      - name: Email
        type: string

  - name: CreateUserRequest
    fields:
      - name: Name
        type: string
      - name: Email
        type: string

endpoints:
  - name: CreateUser
    description: Creates a new user
    request: CreateUserRequest
    response: User
```

#### Generate Code

```bash
# Generate full project structure
uno init -f uno.yaml

# Generate only API package
uno apigen -f uno.yaml
```

The generator creates:

- Type-safe request/response structs
- Client package with strongly-typed methods
- Service interface and implementation stubs
- API documentation (TODO)

Generated code includes:

- Type definitions
- Service interface
- Client implementation
- Request/response handling
- Validation helpers (TODO)

## Examples

examples [README](./examples/README.md)