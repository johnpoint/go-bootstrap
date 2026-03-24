# go-bootstrap

A lightweight, modular Go framework library for building production-ready backend services with standardized component lifecycle management, HTTP API servers, and message queue integration.

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Features

- ✅ **Component Lifecycle Management** - Standardized initialization and lifecycle for all components
- ✅ **HTTP API Server** - Built on Gin framework with flexible routing and middleware support
- ✅ **RabbitMQ Integration** - Producer/consumer patterns with configuration-driven setup
- ✅ **Error Handling** - Custom error types with codes and messages for structured error handling
- ✅ **Cryptographic Utilities** - AES encryption and MD5 hashing utilities
- ✅ **Configurable Logging** - JSON/Text format logging with multiple severity levels
- ✅ **Multiple Codecs** - JSON, form-data, octet-stream, and text request/response handling
- ✅ **Global & Local Components** - Flexible component registration patterns

## Installation

```bash
go get github.com/johnpoint/go-bootstrap/v2
```

## Quick Start

### Basic Bootstrap Setup

```go
package main

import (
    "log/slog"
    "github.com/johnpoint/go-bootstrap/v2/core"
    ginBoot "github.com/johnpoint/go-bootstrap/v2/gin"
)

func main() {
    boot := core.NewBoot(
        core.WithComponents(
            ginBoot.NewApiServer("0.0.0.0:8888"),
        ),
        core.SetLoggerType(core.LoggerTypeText),
        core.Level(slog.LevelDebug),
    )
    
    boot.Init()
}
```

### Creating an Endpoint

```go
import (
    "io"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/johnpoint/go-bootstrap/v2/gin"
)

type HelloRequest struct {
    Name string `json:"name"`
}

type HelloResponse struct {
    Message string `json:"message"`
}

type HelloEndpoint struct{}

func (e *HelloEndpoint) Method() string { return http.MethodGet }
func (e *HelloEndpoint) Path() string { return "/hello" }

func (e *HelloEndpoint) HandlerFunc() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Handle request
    }
}

func (e *HelloEndpoint) NewEncoder(w io.Writer) gin.Encoder {
    return gin.NewJsonEncoder(w)
}

func (e *HelloEndpoint) NewDecoder(r *http.Request) gin.Decoder {
    return gin.NewJsonDecoder(r)
}
```

## Package Structure

### `berror` - Error Handling

Custom error type for structured error handling with codes and messages.

```go
type BErr struct {
    Code      int       // Error code
    Message   string    // User-facing message
    ErrorInfo error     // Internal error details
}

// Helper functions
WarpErr(code int, msg string, err error) *BErr
DecodeErr(err interface{}) (int, string)
GetErrCode(err error) int
GetErrMessage(err error) string
```

### `core` - Bootstrapping Framework

Core framework for managing component lifecycle and application initialization.

```go
// Component interface - implement this for custom components
type Component interface {
    Init(ctx context.Context) error
}

// Bootstrap setup
boot := core.NewBoot(
    core.WithComponents(...),
    core.SetLoggerType(core.LoggerTypeJSON),
    core.Level(slog.LevelInfo),
)
boot.Init()
```

**Configuration Options:**
- `WithComponents(components ...Component)` - Register components
- `SetLoggerType(loggerType string)` - JSON or Text logging
- `Level(level slog.Level)` - Log level (Debug, Info, Warn, Error)

### `gin` - HTTP API Server

HTTP server built on Gin framework with flexible routing, middleware, and multiple codec support.

```go
server := ginBoot.NewApiServer(
    "0.0.0.0:8080",
    middlewares...,
)

// Register endpoints
server.RegisterEndpoints(
    &HelloEndpoint{},
    &GetUserEndpoint{},
)
```

**Features:**
- Built-in health check endpoint: `GET /` returns `{"state": "OK", "time": <timestamp>}`
- Multiple codec support: JSON, form-data, octet-stream, text
- Per-endpoint middleware support
- Global middleware support
- Automatic error handling and logging

**Available Endpoint Types:**
- `GetEndpoint` - HTTP GET
- `PostEndpoint` - HTTP POST
- `PutEndpoint` - HTTP PUT
- `DeleteEndpoint` - HTTP DELETE
- `PatchEndpoint` - HTTP PATCH
- `HeadEndpoint` - HTTP HEAD
- `OptionsEndpoint` - HTTP OPTIONS
- `AnyEndpoint` - Matches any HTTP method

### `rabbitmq` - Message Queue Integration

RabbitMQ producer/consumer integration with configuration-driven setup.

```go
// Configuration
config := &rabbitmq.Config{
    Address:         "amqp://guest:guest@localhost:5672/",
    ExchangeName:    "my.exchange",
    ExchangeKind:    "direct",
    QueueName:       "my.queue",
    BindKey:         "my.key",
    DeliveryMode:    2, // 2 = persistent
    PrefetchCount:   1,
    ChannelNum:      2,
}

// Producer
mq := rabbitmq.New().
    SetContext(ctx).
    SetConfig(config).
    SetHandle(func(msg []byte) error {
        // Handle message
        return nil
    })
mq.Send([]byte("hello"))

// Consumer
mq.Consume()
```

### `utils` - Utility Functions

Cryptographic and utility functions for common operations.

**Encryption:**
```go
import "github.com/johnpoint/go-bootstrap/v2/utils"

// AES encryption/decryption
ciphertext := utils.AesEncrypt(plaintext, key, iv)
plaintext := utils.AesDecrypt(ciphertext, key, iv)

// With Base64 encoding
encoded := utils.EncryptByAes(plaintext, key, iv)
decoded := utils.DecryptByAes(encoded, key, iv)

// MD5 hashing
hash := utils.Md5Hash("data")

// Random generation
randomStr := utils.RandomStr(length)
randomInt := utils.RandomInt(min, max)
```

## Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/rabbitmq/amqp091-go` - RabbitMQ client
- `github.com/go-playground/validator/v10` - Validation
- `github.com/google/uuid` - UUID generation
- `github.com/json-iterator/go` - JSON parsing
- `github.com/stretchr/testify` - Testing utilities

See `go.mod` for the complete list of dependencies.

## Examples

Complete examples can be found in the `test/` directory:

- `test/core_test.go` - Bootstrap initialization and server setup example

## Configuration

### Logger Configuration

```go
boot := core.NewBoot(
    core.SetLoggerType(core.LoggerTypeJSON), // or LoggerTypeText
    core.Level(slog.LevelDebug),             // Log level
)
```

**Logger Types:**
- `LoggerTypeJSON` - JSON structured logging
- `LoggerTypeText` - Human-readable text logging

**Log Levels:**
- `slog.LevelDebug` - Debug messages
- `slog.LevelInfo` - Informational messages
- `slog.LevelWarn` - Warning messages
- `slog.LevelError` - Error messages

### RabbitMQ Configuration

```go
type Config struct {
    Address        string // Connection URL
    ExchangeName   string // Exchange name
    ExchangeKind   string // direct, topic, fanout, headers
    QueueName      string // Queue name
    BindKey        string // Routing key
    DeliveryMode   uint8  // 2 = persistent
    PrefetchCount  int    // Prefetch count
    ChannelNum     int    // Number of channels
}
```

## Project Layout

```
go-bootstrap/
├── berror/           # Custom error handling
├── core/             # Bootstrapping framework
├── gin/              # HTTP API server
├── middleware/       # HTTP middleware components
├── rabbitmq/         # RabbitMQ integration
├── utils/            # Utility functions (crypto, random)
├── test/             # Examples and tests
├── go.mod            # Module definition
├── go.sum            # Dependency checksums
└── README.md         # This file
```

## Version

**Current Version:** v2  
**Go Version Requirement:** 1.25.0+

## Community & Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/johnpoint/go-bootstrap).

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Built with ❤️ for Go developers**
