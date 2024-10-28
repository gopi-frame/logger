package lumberjack

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"strings"
)

var handlerName = "lumberjack"

//goland:noinspection GoBoolExpressions
func init() {
	if handlerName != "" {
		logger.RegisterHandler(handlerName, func(config map[string]any) (io.WriteCloser, error) {
			return NewLumberjackHandlerFromConfig(config)
		})
	}
}

type LumberjackHandler = lumberjack.Logger

func NewLumberjackHandler(filename string) *LumberjackHandler {
	return &lumberjack.Logger{
		Filename: filename,
	}
}

func NewLumberjackHandlerFromConfig(config map[string]any) (*LumberjackHandler, error) {
	var handler LumberjackHandler
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &handler,
		WeaklyTypedInput: true,
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName) || strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))
		},
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			env.ExpandStringWithEnvHookFunc(),
			env.ExpandSliceWithEnvHookFunc(),
			env.ExpandStringKeyMapWithEnvHookFunc(),
			mapstructure.StringToBasicTypeHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
		),
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return &handler, nil
}
