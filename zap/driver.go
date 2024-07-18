package zap

import (
	lc "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
)

var driverName = "zap"

func init() {
	if driverName != "" {
		logger.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (c *Driver) Open(options map[string]any) (lc.Logger, error) {
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewLogger(cfg)
}
