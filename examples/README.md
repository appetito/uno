# Examples

This repository contains examples of how to use the Uno framework.

We have:
- Greeter service 
- Greetings analytics service

Greeter service has API to greet users, it sends greetings to stream `GREETS` and analytics service has API to get stats about greeted users.

## Pre-requisites

- [NATS Server](https://nats.io/)
- [NATS CLI tool](https://github.com/nats-io/natscli)
- [Go](https://go.dev/)

We need to have [NATS Server](https://nats.io/) running on port 4222 with JetStream enabled.

Also we need `GREETS` stream and `ga` consumer created:

```bash
nats stream create --storage memory --subjects greets --defaults GREETS
nats consumer create --pull --defaults GREETS ga
```

## APIs

API of greeter service:

```go
Greet(GreetRequest) -> Greeting
```

API of greeter analytics service:

```go
GetUsersStats(GetUsersStatsRequest) -> UserStats
TopGreetedUsers(TopGreetedUsersRequest) -> []UserStats
```

## Running

Run greeter service: 

```bash
cd examples/greeter
go run cmd/main.go
```

Run greeter analytics service:

```bash
cd examples/greetanalytics
go run cmd/main.go
```

Greet some users:

```bash
nats req "example.Greeter.Greet" '{"Name":"Alice"}'
nats req "example.Greeter.Greet" '{"Name":"Alice"}'
nats req "example.Greeter.Greet" '{"Name":"Bob"}'
```

Check greeter analytics service:

```bash
nats req "example.GreetAnalytics.TopGreetedUsers" '{"Count":5}
nats req "example.GreetAnalytics.GetUsersStats" '{"Name":"Alice"}'
```

Add some load:

```bash
go run examples/test.go --workers=100 --greets=10000
```

Output should be something like this:
```bash
Starting 100 workers with 10000 greetings
Elapsed: 2.320689625s
Average greets per second: 4308.967123
Top greeted users: [{Bob 1077} {David 1053} {Grace 1020} {Ivan 1001} {Eve 993} {Jack 987} {Alice 987} {Helen 968} {Frank 962} {Charlie 954}]
```

## Play with Uno

Both services generated with Uno codegen tool. Services definitions is in the `uno.yaml` files. You can extend/update services API and re-generate code with:

```bash
uno apigen -f uno.yaml
```

## Getting Started with Examples

These example services showcase Uno's code generation capabilities. Each service is defined in its respective `uno.yaml` file, which serves as the source of truth for the service's API and types.

To modify or extend a service:

1. Edit the service definition in `uno.yaml`
2. Regenerate the service API code:
```bash
# Update API code only (after making changes)
uno apigen -f uno.yaml
```
3. Update service implementation to support the new API
4. Run and fun