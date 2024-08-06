# slog
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/logger/driver/slog.svg)](https://pkg.go.dev/github.com/gopi-frame/logger/driver/slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/logger/driver/slog)](https://goreportcard.com/report/github.com/gopi-frame/logger/driver/slog)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package slog implements a logger using the [slog](https://pkg.go.dev/log/slog) package.

## Installation

```shell
go get -u github.com/gopi-frame/logger/driver/slog
```

## Import

```go
import "github.com/gopi-frame/logger/driver/slog"
```

## Quick Start

```go
package main

import "github.com/gopi-frame/logger/driver/slog"

func main() {
    l, err := slog.NewLogger(nil)
    if err != nil {
        panic(err)
    }
    l.Info("Hello world")
}
```

## Advance Usage

```go
package main

import (
    "context"
    "github.com/gopi-frame/logger"
    "github.com/gopi-frame/logger/driver/slog"
    loggercontract "github.com/gopi-frame/contract/logger"
)

func main() {
    cfg := slog.NewConfig()
    cfg.Apply(
        slog.WithLevel(slog.LevelDebug)
    slog.AddSource(),
        slog.WithFields(map[string]any{
            "key": "value",
        }),
)
    l, err := slog.NewLogger(cfg)
    if err != nil {
        panic(err)
    }
    l.Info("Hello world")
    // log a formatted message
    l.Infof("This is a INFO level formatted message: %s", "Hello World")
    // log a message with context fields
    ctx := logger.WithValue(context.Background(), map[string]any{
        "key": "value",
    })
    l.WithContext(ctx).Info("Hello World")
    // create a child logger with new level
    childLog := l.WithLevel(loggercontract.LevelDebug)
    childLog.Debug("This is a DEBUG level message")
}
```

