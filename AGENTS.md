# Agent Guidelines for learn-pub-sub-starter

This project is a Go-based RabbitMQ pub/sub learning starter project. It uses AMQP for messaging with RabbitMQ.

## Project Structure

```
.
├── cmd/
│   ├── server/main.go       # Server entry point
│   └── client/main.go       # Client entry point
├── internal/
│   ├── pubsub/              # RabbitMQ pub/sub utilities
│   │   ├── publish.go
│   │   ├── subscribe.go
│   │   ├── consume.go
│   │   └── utils.go
│   ├── routing/            # Exchange, queue, and routing key constants
│   │   ├── routing.go
│   │   └── models.go
│   └── gamelogic/          # Game logic (CLI interface)
└── go.mod
```

## Build, Lint, and Test Commands

### Building

```bash
# Build all binaries
go build ./...

# Build server
go build -o server ./cmd/server

# Build client
go build -o client ./cmd/client
```

### Running

```bash
# Run server
go run ./cmd/server

# Run client
go run ./cmd/client
```

### Testing

This project currently has no test files.

```bash
# Run all tests (none exist currently)
go test ./...

# Run tests for a specific package
go test ./internal/pubsub

# Run tests verbosely
go test -v ./...

# Run a single test by name
go test -v -run TestName ./...
```

### Linting

```bash
# Run go vet
go vet ./...

# Format code
go fmt ./...

# Run golint (if installed)
golint ./...

# Run staticcheck (if installed)
staticcheck ./...
```

## Code Style Guidelines

### General

- Use Go 1.22.1 or later (see go.mod)
- Run `go fmt ./...` before committing
- Run `go vet ./...` to catch common errors

### Imports

- Use the short alias `amqp "github.com/rabbitmq/amqp091-go"` for the RabbitMQ client
- Group imports: standard library first, then third-party
- Use blank imports sparingly

Example:
```go
import (
    "context"
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"

    "github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
    "github.com/jondatkins/learn-pub-sub-starter/internal/routing"
)
```

### Formatting

- Use `gofmt` (standard Go formatter)
- 4-space indentation (enforced by gofmt)
- No unnecessary blank lines within function bodies
- Group related code with blank lines between logical sections

### Naming Conventions

- **Packages**: lowercase, short, descriptive (e.g., `pubsub`, `routing`, `gamelogic`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interfaces**: PascalCase, typically with `-er` suffix (e.g., `Publisher`, `Subscriber`)
- **Structs**: PascalCase
- **Struct fields**: camelCase, even when exported for JSON serialization

Example for enum-like constants:
```go
type SimpleQueueType string

const (
    QueueTypeDurable   SimpleQueueType = "durable"
    QueueTypeTransient SimpleQueueType = "transient"
)
```

### Types

- Use generics where appropriate (e.g., `PublishJSON[T any]`)
- Use meaningful type names rather than primitives where clarity benefits
- Prefer `error` over returning error codes

### Error Handling

- Always check and handle errors, don't ignore with `_`
- Return errors up the stack rather than logging and continuing when possible
- Use `fmt.Errorf` with `%w` for wrapped errors:
  ```go
  return fmt.Errorf("failed to publish: %w", err)
  ```
- For fatal errors during initialization, use `log.Fatal` or `log.Fatalf`
- For recoverable errors in goroutines, log and continue

### RabbitMQ Patterns

- Always declare exchanges before publishing to them
- Use `ExchangeDeclare` with appropriate type (direct, topic, fanout)
- Declare queues before binding them
- Use `QueueBind` to bind queues to exchanges with routing keys
- Close channels and connections when done
- Use `PublishWithContext` for publishing with context support
- Use generic `PublishJSON` and `SubscribeJSON` functions for type-safe messaging
- Always acknowledge messages with `delivery.Ack()` or reject with `delivery.Nack()`

### Concurrency

- Use goroutines for consuming messages
- Use channels for communication between goroutines
- Ensure proper synchronization when sharing data between goroutines
- Use `context.Context` for cancellation and timeouts

### Logging

- Use `log` package for application logging
- Use `fmt` package for user-facing output (CLI applications)
- Include relevant context in log messages using `%+v` or similar

### Comments

- Comment exported functions and types to help users
- Use complete sentences with proper punctuation
- Update comments when changing code

### Testing (when added)

- Use table-driven tests where appropriate
- Name test files `*_test.go`
- Use descriptive test names: `TestFunctionName_Scenario_ExpectedBehavior`
- Use `t.Fatal` or `t.Fatalf` when a test cannot continue
- Use `t.Log` for test-specific debug output
