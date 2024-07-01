package zap

import (
	"github.com/gopi-frame/logger"
	"github.com/gopi-frame/logger/driver"
)

var driverName = "zap"

func init() {
	if driverName != "" {
		logger.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (c *Driver) Open(options map[string]any) (driver.Logger, error) {
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewLogger(cfg)
}
