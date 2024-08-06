# Logger
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgopi-frame%2Flogger.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgopi-frame%2Flogger?ref=badge_shield)

Package logger is a package for managing logger drivers and creating logger instances.

## Installation

```shell
go get -u github.com/gopi-frame/logger
```

## Import

```go
import "github.com/gopi-frame/logger"
```

## Usage

```go
package main

import (
	"github.com/gopi-frame/logger"
	
	_ "github.com/gopi-frame/logger/driver/zap"
)

func main() {
	log, err := logger.Open("zap", map[string]any{
		"Level": "debug",
    })
}
```

## Drivers

- [zap](driver/zap/README.md)
- [slog](driver/slog/README.md)

## How to create a custom driver

To create a custom driver, just implement
the [logger.Driver](https://pkg.go.dev/github.com/gopi-frame/contract/logger#Driver) interface
and register it using [logger.Register](https://pkg.go.dev/github.com/gopi-frame/logger#Register)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgopi-frame%2Flogger.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgopi-frame%2Flogger?ref=badge_large)