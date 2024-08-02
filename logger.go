// Package logger is a package for managing logger drivers and creating loggers.
package logger

import (
	"context"
	"fmt"
	"sort"

	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/exception"
)

var drivers = kv.NewMap[string, logger.Driver]()

// Register registers a new logger driver.
// If a driver with the same name already exists, it panics.
func Register(driverName string, driver logger.Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if driver == nil {
		panic(exception.NewEmptyArgumentException("driver"))
	}
	if drivers.ContainsKey(driverName) {
		panic(exception.NewArgumentException("driverName", driver, fmt.Sprintf("duplicate driver \"%s\"", driverName)))
	}
	drivers.Set(driverName, driver)
}

// Drivers returns a list of registered logger drivers.
func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	sort.Strings(list)
	return list
}

// Open opens a new logger using the given driver name and options.
func Open(driverName string, options map[string]any) (logger.Logger, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	return driver.Open(options)
}

var ctxValueKey = struct {
	key string
}{
	key: "valueKey",
}

// WithValue returns a new context that carries value.
func WithValue(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, ctxValueKey, value)
}

// GetValue returns the value stored in ctx, if any.
func GetValue(ctx context.Context) any {
	return ctx.Value(ctxValueKey)
}
