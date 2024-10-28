package slog

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config is the configuration for the [Logger].
type Config struct {
	Level        Level          `json:"level" yaml:"level" toml:"level"`
	Fields       map[string]any `json:"fields" yaml:"fields" toml:"fields"`
	Encoder      string         `json:"encoder" yaml:"encoder" toml:"encoder"`
	AddSource    bool           `json:"addSource" yaml:"addSource" toml:"addSource"`
	PanicOnFatal bool           `json:"panicOnFatal" yaml:"panicOnFatal" toml:"panicOnFatal"`
	Handler      string         `json:"handler" yaml:"handler" toml:"handler"`
	HandlerWith  map[string]any `json:"handlerWith" yaml:"handlerWith" toml:"handlerWith"`
}

// NewConfig creates a new [Config] instance with default values.
func NewConfig() *Config {
	return &Config{
		Level:   Level{slog.LevelDebug},
		Fields:  make(map[string]any),
		Encoder: EncoderJSON,
	}
}

// SlogHandler creates a new slog handler.
func (c *Config) SlogHandler() (slog.Handler, error) {
	var h slog.Handler
	var opts = &slog.HandlerOptions{
		Level:     c.Level.Level,
		AddSource: c.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				label, ok := levelLabels[level]
				if ok {
					a.Value = slog.StringValue(strings.ToUpper(label))
				} else {
					a.Value = slog.StringValue(level.String())
				}
			}
			return a
		},
	}
	opts.AddSource = c.AddSource
	var w io.WriteCloser
	if c.Handler != "" {
		var err error
		w, err = logger.CreateHandler(c.Handler, c.HandlerWith)
		if err != nil {
			return nil, err
		}
	} else {
		w = os.Stdout
	}
	if c.Encoder == EncoderText {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}
	return &handler{
		handler: h,
	}, nil
}

// UnmarshalOptions unmarshal the options.
func UnmarshalOptions(options map[string]any) (*Config, error) {
	cfg := NewConfig()
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			env.ExpandStringWithEnvHookFunc(),
			env.ExpandSliceWithEnvHookFunc(),
			env.ExpandStringKeyMapWithEnvHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToBasicTypeHookFunc(),
		),
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName) || strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))
		},
		WeaklyTypedInput: true,
		Result:           cfg,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(options); err != nil {
		return nil, err
	}
	return cfg, nil
}
