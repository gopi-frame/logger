package slog

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/writer"
	"io"
	"log/slog"
	"os"
	"strings"
)

var DefaultWriter io.Writer = os.Stdout

// Config is the configuration for the [Logger].
type Config struct {
	Level        Level                     `json:"level" yaml:"level" toml:"level"`
	Fields       map[string]any            `json:"fields" yaml:"fields" toml:"fields"`
	Encoder      string                    `json:"encoder" yaml:"encoder" toml:"encoder"`
	AddSource    bool                      `json:"addSource" yaml:"addSource" toml:"addSource"`
	Writers      map[string]map[string]any `json:"writers" yaml:"writers" toml:"writers"`
	PanicOnFatal bool                      `json:"panicOnFatal" yaml:"panicOnFatal" toml:"panicOnFatal"`
}

// NewConfig creates a new [Config] instance with default values.
func NewConfig() *Config {
	return &Config{
		Level:   Level{slog.LevelDebug},
		Fields:  make(map[string]any),
		Encoder: EncoderJSON,
	}
}

// Option is a function that configures the [Config].
type Option func(*Config) error

// Apply applies the options to the config.
func (c *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

// WithLevel sets the log level.
func WithLevel(level slog.Level) Option {
	return func(cfg *Config) error {
		cfg.Level = Level{level}
		return nil
	}
}

// AddSource adds source to the log message.
func AddSource() Option {
	return func(cfg *Config) error {
		cfg.AddSource = true
		return nil
	}
}

// WithFields sets fields for the log message.
func WithFields(fields map[string]any) Option {
	return func(cfg *Config) error {
		cfg.Fields = fields
		return nil
	}
}

// PanicOnFatal replaces calling os.Exit(1) on fatal level with panic.
func PanicOnFatal() Option {
	return func(cfg *Config) error {
		cfg.PanicOnFatal = true
		return nil
	}
}

// SlogHandler creates a new slog handler.
func (c *Config) SlogHandler() (slog.Handler, error) {
	var w = DefaultWriter
	if len(c.Writers) > 0 {
		var ws []io.Writer
		for driverName, options := range c.Writers {
			w, err := writer.Open(driverName, options)
			if err != nil {
				return nil, err
			}
			ws = append(ws, w)
		}
		w = io.MultiWriter(ws...)
	}
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
			mapstructure.TextUnmarshallerHookFunc(),
		),
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
