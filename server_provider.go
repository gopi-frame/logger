package logger

import (
	"reflect"

	"github.com/gopi-frame/contract/container"
	"github.com/gopi-frame/contract/logger"
)

// ServerProvider logger server provider
type ServerProvider struct{}

// Register register
func (s *ServerProvider) Register(c container.Container) {
	c.Bind(reflect.TypeFor[Logger]().String(), func(c container.Container) any {
		return NewManager()
	})
	c.Alias(reflect.TypeFor[Logger]().String(), "logger")
	c.Alias(reflect.TypeFor[Logger]().String(), reflect.TypeFor[logger.Logger]().String())
}
