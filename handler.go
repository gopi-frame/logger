package logger

import (
	"fmt"
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/exception"
	"io"
)

var handlers = kv.NewMap[string, func(config map[string]any) (io.WriteCloser, error)]()

func RegisterHandler(handlerName string, creator func(config map[string]any) (io.WriteCloser, error)) {
	handlers.Lock()
	defer handlers.Unlock()
	if handlers.ContainsKey(handlerName) {
		panic(exception.NewArgumentException("handlerName", handlerName, fmt.Sprintf("duplicate handler \"%s\"", handlerName)))
	}
	handlers.Set(handlerName, creator)
}

func CreateHandler(handlerName string, config map[string]any) (io.WriteCloser, error) {
	handlers.RLock()
	handler, ok := handlers.Get(handlerName)
	handlers.RUnlock()
	if ok {
		return handler(config)
	}
	return nil, NewUnknownHandlerException(handlerName)
}
