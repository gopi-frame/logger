package lumberjack

import (
	"io"

	"github.com/gopi-frame/logger"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/natefinch/lumberjack.v2"
)

var driverName = "lumberjack"

func init() {
	if driverName != "" {
		logger.RegisterWriter(driverName, new(Driver))
	}
}

type Driver struct {}

func (d Driver) Open(options map[string]any) (io.WriteCloser, error) {
	logger := new(lumberjack.Logger)
	err := mapstructure.WeakDecode(options, logger)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
