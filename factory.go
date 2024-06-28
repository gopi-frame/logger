package logger

import (
	"fmt"
	"io"
	"sort"

	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/logger/contract"
	"github.com/gopi-frame/support/maps"
)

var drivers = maps.NewMap[string, contract.Driver]()

func Register(driverName string, driver contract.Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if driver == nil {
		panic(exception.NewEmptyArgumentException("driver"))
	}
	if _, dup := drivers.Get(driverName); dup {
		panic(exception.NewArgumentException("driverName", driver, fmt.Sprintf("duplicate driver \"%s\"", driverName)))
	}
	drivers.Set(driverName, driver)
}

func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	sort.Strings(list)
	return list
}

func Open(driverName string, options map[string]any) (contract.Logger, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	return driver.Open(options)
}

var writers = maps.NewMap[string, contract.Writer]()

func RegisterWriter(driverName string, writer contract.Writer) {
	writers.Lock()
	defer writers.Unlock()
	if _, dup := writers.Get(driverName); dup {
		panic(exception.NewArgumentException("driverName", driverName, fmt.Sprintf("duplicate writer \"%s\"", driverName)))
	}
	writers.Set(driverName, writer)
}

func Writers() []string {
	writers.RLock()
	list := writers.Keys()
	defer writers.RUnlock()
	sort.Strings(list)
	return list
}

func Writer(driverName string, options map[string]any) (io.WriteCloser, error) {
	writers.Lock()
	w, ok := writers.Get(driverName)
	writers.Unlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown writer driver \"%s\"", driverName))
	}
	return w.Open(options)
}
