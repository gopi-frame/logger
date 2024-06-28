package zap

import (
	"github.com/gopi-frame/logger"
	"github.com/gopi-frame/logger/contract"
)

var driverName = "zap"

func init() {
	if driverName != "" {
		logger.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (c *Driver) Open(options map[string]any) (contract.Logger, error) {
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewLogger(cfg)
}
