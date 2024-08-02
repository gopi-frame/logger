package zap

import (
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/logger/driver/zap.driverName=custom`
var driverName = "zap"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		logger.Register(driverName, new(Driver))
	}
}

// Driver is a zap logger driver.
type Driver struct{}

// Open opens a zap logger.
func (c *Driver) Open(options map[string]any) (loggercontract.Logger, error) {
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewLogger(cfg)
}
