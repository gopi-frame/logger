package file

import (
	"io"

	"github.com/gopi-frame/logger"
)

var driverName = "file"

func init() {
	if driverName != "" {
		logger.RegisterWriter(driverName, new(Driver))
	}
}

type Driver struct{}

func (d Driver) Open(options map[string]any) (io.WriteCloser, error) {
	config, err := UnmarshalOptions(options)
	if err != nil {
		return nil, err
	}
	return NewFileWriter(config)
}
