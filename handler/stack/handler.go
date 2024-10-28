package stack

import (
	"errors"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"io"
	"strings"
)

var handlerName = "stack"

//goland:noinspection GoBoolExpressions
func init() {
	if handlerName != "" {
		logger.RegisterHandler(handlerName, func(config map[string]any) (io.WriteCloser, error) {
			return NewStackHandlerFromConfig(config)
		})
	}
}

type StackHandler struct {
	handlers     []io.WriteCloser
	breakOnError bool
}

// NewStackHandler creates a new stack handler.
func NewStackHandler(handlers ...io.WriteCloser) *StackHandler {
	return &StackHandler{
		handlers: handlers,
	}
}

func NewStackHandlerFromConfig(config map[string]any) (*StackHandler, error) {
	var cfg struct {
		BreakOnError bool
		Handlers     []map[string]any
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &cfg,
		WeaklyTypedInput: true,
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName) || strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))
		},
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			env.ExpandStringWithEnvHookFunc(),
			env.ExpandSliceWithEnvHookFunc(),
			env.ExpandStringKeyMapWithEnvHookFunc(),
			mapstructure.StringToBasicTypeHookFunc(),
		),
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	var handlers []io.WriteCloser
	for _, handler := range cfg.Handlers {
		driver := handler["driver"].(string)
		h, err := logger.CreateHandler(driver, handler)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, h)
	}
	if len(handlers) == 0 {
		return &StackHandler{
			breakOnError: true,
		}, nil
	}
	handler := NewStackHandler(handlers...)
	handler.breakOnError = cfg.BreakOnError
	return handler, nil
}

func (h *StackHandler) Write(p []byte) (n int, err error) {
	var errs []error
	for _, handler := range h.handlers {
		n, err = handler.Write(p)
		if err != nil && h.breakOnError {
			return
		} else if err != nil {
			errs = append(errs, err)
		}
	}
	err = errors.Join(errs...)
	return
}

func (h *StackHandler) Close() error {
	var errs []error
	for _, handler := range h.handlers {
		if err := handler.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	h.handlers = nil
	return errors.Join(errs...)
}
