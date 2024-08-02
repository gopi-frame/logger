package slog

import (
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/logger/driver/slog.driverName=custom`
var driverName = "slog"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		logger.Register(driverName, new(Driver))
	}
}

// Driver is a slog logger driver.
type Driver struct{}

// Open opens a slog logger.
func (d Driver) Open(options map[string]any) (loggercontract.Logger, error) {
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewLogger(cfg)
}
