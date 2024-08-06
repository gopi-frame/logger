# zap
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/logger/driver/zap.svg)](https://pkg.go.dev/github.com/gopi-frame/logger/driver/zap)
[![codecov](https://codecov.io/gh/gopi-frame/logger/graph/badge.svg?token=9EUOUXQ6PD&flag=zap)](https://codecov.io/gh/gopi-frame/logger?flag=zap)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/logger/driver/zap)](https://goreportcard.com/report/github.com/gopi-frame/logger/driver/zap)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package zap is an implementation for [gopi-frame/logger](../../README.md) based on [uber-go/zap](https://github.com/uber-go/zap).
And this package can be used as a stand-alone package.
## Installation
```shell
go get -u github.com/gopi-frame/logger/driver/zap
```

## Import
```go
import "github.com/gopi-frame/logger/driver/zap"
```

## Quick Start

```go
package main

import "github.com/gopi-frame/logger/driver/zap"

func main() {
    // create a logger with default configuration
    log, err := zap.NewLogger(nil) 
    if err != nil {
        panic(err)
    }
    log.Debug("This is a DEBUG level message")
    log.Info("This is a INFO level message")
    log.Warn("This is a WARN level message with context fields", map[string]any{"key": "value"})
    log.Error("This is a ERROR level message")
    log.Panic("This is a PANIC level message")
    log.Fatal("This is a FATAL level message")
    // log a formatted message
    log.Debugf("This is a DEBUG Level formatted message: %s", "Hello World")
    log.Infof("This is a INFO level formatted message: %s", "Hello World")
    log.Warnf("This is a WARN level formatted message: %s", "Hello World")
    log.Errorf("This is a ERROR level formatted message: %s", "Hello World")
    log.Panicf("This is a PANIC level formatted message: %s", "Hello World")
    log.Fatalf("This is a FATAL level formatted message: %s", "Hello World")
}
```

## Advance Usage

```go
package main

import (
    "context"
    "github.com/gopi-frame/logger"
    "github.com/gopi-frame/logger/driver/zap"
    "go.uber.org/zap/zapcore"
    
    loggercontract "github.com/gopi-frame/contract/logger"
)

func main() {
    cfg := zap.NewConfig()
    if err := cfg.Apply(
        zap.Level(zapcore.InfoLevel),
        zap.Encoder(zap.EncoderJSON),
        zap.DurationEncoder(zapcore.ISO8601DurationEncoder),
        zap.TimeEncoder(zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")),
    ); err != nil {
        panic(err)
    }
    log, err := zap.NewLogger(cfg)
    if err != nil {
        panic(err)
    }
    log.Debug("This is a DEBUG level message")
    log.Info("This is a INFO level message")
    log.Warn("This is a WARN level message with context fields", map[string]any{"key": "value"})
    log.Error("This is a ERROR level message")
    log.Panic("This is a PANIC level message")
    log.Fatal("This is a FATAL level message")
    // log a formatted message
    log.Debugf("This is a DEBUG Level formatted message: %s", "Hello World")
    log.Infof("This is a INFO level formatted message: %s", "Hello World")
    log.Warnf("This is a WARN level formatted message: %s", "Hello World")
    log.Errorf("This is a ERROR level formatted message: %s", "Hello World")
    log.Panicf("This is a PANIC level formatted message: %s", "Hello World")
    log.Fatalf("This is a FATAL level formatted message: %s", "Hello World")
    // log a message with context
    ctx := logger.WithValue(context.Background(), map[string]any{"key": "value"})
    log.WithContext(ctx).Info("This is a INFO level message with context fields")
    // create a child logger with new level
    childLog := log.WithLevel(loggercontract.LevelDebug)
    childLog.Debug("This is a DEBUG level message")
    
}
```

## Use With Driver

```go
package main

import "github.com/gopi-frame/logger/driver/zap"
//import "github.com/gopi-frame/logger" 

func main() {
	var options = map[string]any{
		"level":   "warn",
		"encoder": "json",
	}
	log, err := new(zap.Driver).Open(options) // or logger.Open("zap", options)
    if err != nil {
		panic(err)
    }
	log.Debug("This is a DEBUG level message")
	log.Info("This is a INFO level message")
	log.Warn("This is a WARN level message with context fields", map[string]any{"key": "value"})
	log.Error("This is a ERROR level message")
	log.Panic("This is a PANIC level message")
	log.Fatal("This is a FATAL level message")
}
```

## Options

This package uses [mapstructure](https://github.com/go-viper/mapstructure/v2) to parse options.

For more information about options, please
see [logger/zap.Config](https://pkg.go.dev/github.com/gopi-frame/logger/zap#Config)

Example:
```go
var options = map[string]any{
	"level": "warn",
	"encoder": "json",
	"encoderConfig": map[string]any{
		"timeKey": "ts",
    },
}
```

### Level

Option `level` is used to set the log level, default is `warn`.
Message with level lower than the level will be ignored.

```go
var options = map[string]any{
	"level": "info",
}
```

For more information about level, please 
see [zapcore.Level](https://pkg.go.dev/github.com/uber-go/zap/zapcore#Level).

For more information about how the level is parsed, please
see [logger/zap.DecodeLevelHook](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DecodeLevelHook)

#### Available config values for level

| Type           | Value     | Decoded value                                                                 | 
|----------------|-----------|-------------------------------------------------------------------------------|
| string\|number | debug, -1 | [zapcore.DebugLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#DebugLevel)   |
| string\|number | info, 0   | [zapcore.InfoLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#InfoLevel)     |
| string\|number | warn, 1   | [zapcore.WarnLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#WarnLevel)     |
| string\|number | error, 2  | [zapcore.ErrorLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#ErrorLevel)   |
| string\|number | dpanic, 3 | [zapcore.DPanicLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#DPanicLevel) |
| string\|number | panic, 4  | [zapcore.PanicLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#PanicLevel)   |
| string\|number | fatal, 5  | [zapcore.FatalLevel](https://pkg.go.dev/go.uber.org/zap/zapcore#FatalLevel)   |

### Development

Option `development` is used to set development mode, default is `false`.

```go
var options = map[string]any{
	"development": true,
}
```

### Fields

Option `fields` is used to add fields to the logger, default is `nil`.

```go
var options = map[string]any{
	"fields": map[string]any{
	    "key": "value",	
    },
}
```

### Caller

Option `caller` is used to add caller to the logger, default is `false`.

```go
var options = map[string]any{
	"caller": true,
}
```

### CallerSkip

Option `callerSkip` is used to set the extra number of callers to skip, default is `0`.

```go
var options = map[string]any{
	"callerSkip": 1,
}
```

### Encoder

Option `encoder` is used to set the encoder type, `json` and `text` are available, default is `json`.

```go
var options = map[string]any{
    "encoder": "json",
}
```

### EncoderConfig

Option `encoderConfig` is used to set the encoder config.

#### MessageKey

Option `messageKey` is used to set the message key, default is `message`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"messageKey": "message",
    },
}
```

#### LevelKey

Option `levelKey` is used to set the level key, default is `level`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"levelKey": "level",
	}
}
```

#### TimeKey

Option `timeKey` is used to set the time key, default is `time`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"timeKey": "time",
	},
}
```

#### NameKey

Option `nameKey` is used to set the name key, default is `name`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"nameKey": "name",
    },
}
```

#### CallerKey

Option `callerKey` is used to set the caller key, default is `caller`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "callerKey": "caller",	
    },
}
```

#### FunctionKey

Option `functionKey` is used to set the function key, default is `function`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "functionKey": "function",	
    },
}
```

#### StacktraceKey

Option `stacktraceKey` is used to set the stacktrace key, default is `stacktrace`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "stacktraceKey": "stacktrace",	
    },
}
```

#### LevelEncoder

Option `LevelEncoder` is used to set the level encoder, default is `lowercase` which means `LowercaseLevelEncoder`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"levelEncoder": "lowercase",
	},
}
```

For more information about level encoder, please
see [zapcore.LevelEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#LevelEncoder).

For more information about how the level encoder is parsed, please
see [logger/zap.DecodeLevelEncoderHook](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DecodeLevelEncoderHook)

**Available config value for level encoder**

| Type   | Value        | Decoded Value                                                                                               |
|--------|--------------|-------------------------------------------------------------------------------------------------------------|
| string | lowercase    | [zapcore.LowercaseLevelEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#LowercaseLevelEncoder)           |
| string | capital      | [zapcore.CapitalLevelEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#CapitalLevelEncoder)               |
| string | color        | [zapcore.LowercaseColorLevelEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#LowercaseColorLevelEncoder) |
| string | capitalColor | [zapcore.CapitalColorLevelEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#CapitalColorLevelEncoder)     |

#### DurationEncoder

Option `durationEncoder` is used to set the duration encoder, default is `string` which means `StringDurationEncoder`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
		"durationEncoder": "string",
	},
}
```

For more information about duration encoder, please
see [zapcore.DurationEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#DurationEncoder).

For more information about how the duration encoder is parsed, please
see [logger/zap.DecodeDurationEncoderHook](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DecodeDurationEncoderHook)

**Available config value for duration encoder**

| Type           | Value   | Decoded Value                                                                                       |
|----------------|---------|-----------------------------------------------------------------------------------------------------|
| string\|number | seconds | [zapcore.SecondsDurationEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#SecondsDurationEncoder) |
| string\|number | string  | [zapcore.StringDurationEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#StringDurationEncoder)   |
| string\|number | nanos   | [zapcore.NanosDurationEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#NanosDurationEncoder)     |
| string\|number | ms      | [zapcore.MillisDurationEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#MillisDurationEncoder)   |

#### TimeEncoder

Option `timeEncoder` is used to set the time encoder, default is `rfc3339` which means `RFC3339TimeEncoder`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "timeEncoder": "rfc3339",	
    },
}
```

For more information about time encoder, please
see [zapcore.TimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#TimeEncoder).

For more information about how the time encoder is parsed, please
see [logger/zap.DecodeTimeEncoderHook](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DecodeTimeEncoderHook)

**Available config value for time encoder**

| Type   | Value       | Decoded Value                                                                                                           |
|--------|-------------|-------------------------------------------------------------------------------------------------------------------------|
| string | rfc3339nano | [zapcore.RFC3339NanoTimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#RFC3339NanoTimeEncoder)                     |
| string | rfc3339     | [zapcore.RFC3339TimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#RFC3339TimeEncoder)                             |
| string | iso8601     | [zapcore.ISO8601TimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#ISO8601TimeEncoder)                             |
| string | millis      | [zapcore.EpochMillisTimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#EpochMillisTimeEncoder)                     |
| string | nanos       | [zapcore.EpochNanosTimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#EpochNanosTimeEncoder#EpochNanosTimeEncoder) |
| string | timestamp   | [zapcore.EpochTimeEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#EpochTimeEncoder)                                 |

#### How to encode time with custom layout?
If you want to encode time with custom layout, use a map which contains a key named "layout" to specify the
layout.

```go
var options = map[string]any{
	"timeEncoder": map[string]any{
	    "layout": "2006-01-02 15:04:05",	
    }
}
```

#### CallerEncoder

Option `callerEncoder` is used to set the caller encoder, default is `short` which means `ShortCallerEncoder`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "callerEncoder": "short",	
    },
}
```

For more information about caller encoder, please
see [zapcore.CallerEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#CallerEncoder).

For more information about how the caller encoder is parsed, please
see [logger/zap.DecodeCallerEncoderHook](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DecodeCallerEncoderHook)

**Available config value for caller encoder**

| Type   | Value | Decoded Value                                                                               |
|--------|-------|---------------------------------------------------------------------------------------------|
| string | full  | [zapcore.FullCallerEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#FullCallerEncoder)   |
| string | short | [zapcore.ShortCallerEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#ShortCallerEncoder) |

#### NameEncoder

Option `nameEncoder` is used to set the name encoder, default is `full` which means `FullNameEncoder`.

```go
var options = map[string]any{
	"encoderConfig": map[string]any{
	    "nameEncoder": "full",	
    },
}
```

For more information about name encoder, please
see [zapcore.NameEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#NameEncoder).

**Available config value for name encoder**

| Type   | Value | Decoded Value                                                                         |
|--------|-------|---------------------------------------------------------------------------------------|
| string | full  | [zapcore.FullNameEncoder](https://pkg.go.dev/go.uber.org/zap/zapcore#FullNameEncoder) |

#### LineEnding

Option `lineEnding` is used to set the line ending, default is `\n`.

```go
var options = map[string]any{
	"lineEnding": "\n",
}
```

#### SkipLineEnding

Option `skipLineEnding` is used to skip the line ending, default is `false`.

```go
var options = map[string]any{
	"skipLineEnding": false,
}
```

### Writers

Option `writers` is used to set the writers.

When `writers` is not set, [DefaultWriter](https://pkg.go.dev/github.com/gopi-frame/logger/zap#DefaultWriter) 
will be used as the writer.

When `writers` is set, the key of the map should be the writer [Driver](https://pkg.go.dev/github.com/gopi-frame/contract/writer#Driver), 
and the value should be the configuration of the writer.

For more information about the writer configuration, please
see [https://github.com/gopi-frame/writer](https://github.com/gopi-frame/writer)

### Hooks

Option `hooks` is used to set the hooks.

For more information, please
see [zap.Hooks](https://pkg.go.dev/go.uber.org/zap#Hooks)

### Stacktrace

Option `stacktrace` is used to set the level which should record the stacktrace.

```go
var options = map[string]any{
	"stacktrace": zapcore.ErrorLevel,
}
```

For more information, please
see [zap.AddStacktrace](https://pkg.go.dev/go.uber.org/zap#AddStacktrace)

### IncreaseLevel

Option `increaseLevel` is used to increase the level of the logger.

```go
var options = map[string]any{
	"increaseLevel": zapcore.ErrorLevel,
}
```

For more information, please
see [zap.IncreaseLevel](https://pkg.go.dev/go.uber.org/zap#IncreaseLevel)


### PanicHook

Option `panicHook` is used to set the hook after the panic level log is written.

For more information, please
see [zap.WithPanicHook](https://pkg.go.dev/go.uber.org/zap#WithPanicHook)

### FatalHook

Option `fatalHook` is used to set the hook after the fatal level log is written.

For more information, please
see [zap.WithFatalHook](https://pkg.go.dev/go.uber.org/zap#WithFatalHook)


### Clock

Option `clock` is used to set the clock.

For more information, please
see [zap.WithClock](https://pkg.go.dev/go.uber.org/zap#WithClock)
