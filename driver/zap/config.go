package zap

import (
	"os"

	"github.com/gopi-frame/logger"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level         zapcore.Level             `json:"level" yaml:"level" toml:"level"`
	Development   bool                      `json:"development" yaml:"development" toml:"development"`
	Fields        map[string]any            `json:"fields" yaml:"fields" toml:"fields"`
	Caller        bool                      `json:"caller" yaml:"caller" toml:"caller"`
	CallerSkip    int                       `json:"caller_skip" yaml:"caller_skip" toml:"caller_skip"`
	Encoder       string                    `json:"encoder" yaml:"encoder" toml:"encoder"`
	EncoderConfig zapcore.EncoderConfig     `json:"encoder_config" yaml:"encoder_config" toml:"encoder_config"`
	Outputs       map[string]map[string]any `json:"outputs" yaml:"outputs" toml:"outputs"`

	Hooks         []func(zapcore.Entry) error `json:"-" yaml:"-" toml:"-"`
	Stacktrace    zapcore.LevelEnabler        `json:"-" yaml:"-" toml:"-"`
	IncreaseLevel zapcore.LevelEnabler        `json:"-" yaml:"-" toml:"-"`
	PanicHook     zapcore.CheckWriteHook      `json:"-" yaml:"-" toml:"-"`
	FatalHook     zapcore.CheckWriteHook      `json:"-" yaml:"-" toml:"-"`
	Clock         zapcore.Clock               `json:"-" yaml:"-" toml:"-"`
}

func NewConfig() *Config {
	return &Config{
		Level:         zapcore.WarnLevel,
		Encoder:       EncoderJSON,
		EncoderConfig: zapcore.EncoderConfig{},
	}
}

type Option func(cfg *Config) error

func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}

func Level(level zapcore.Level) Option {
	return func(cfg *Config) error {
		cfg.Level = level
		return nil
	}
}

func Encoder(level string) Option {
	return func(cfg *Config) error {
		if level == "" {
			return nil
		}
		cfg.Encoder = level
		return nil
	}
}

func MessageKey(messageKey string) Option {
	return func(cfg *Config) error {
		if messageKey == "" {
			return nil
		}
		cfg.EncoderConfig.MessageKey = messageKey
		return nil
	}
}

func Levelkey(levelKey string) Option {
	return func(cfg *Config) error {
		if levelKey == "" {
			return nil
		}
		cfg.EncoderConfig.LevelKey = levelKey
		return nil
	}
}

func TimeKey(timeKey string) Option {
	return func(cfg *Config) error {
		if timeKey == "" {
			return nil
		}
		cfg.EncoderConfig.TimeKey = timeKey
		return nil
	}
}

func NameKey(nameKey string) Option {
	return func(cfg *Config) error {
		if nameKey == "" {
			return nil
		}
		cfg.EncoderConfig.NameKey = nameKey
		return nil
	}
}

func CallerKey(callerKey string) Option {
	return func(cfg *Config) error {
		if callerKey == "" {
			return nil
		}
		cfg.EncoderConfig.CallerKey = callerKey
		return nil
	}
}

func FunctionKey(functionKey string) Option {
	return func(cfg *Config) error {
		if functionKey == "" {
			return nil
		}
		cfg.EncoderConfig.FunctionKey = functionKey
		return nil
	}
}

func StacktraceKey(stacktraceKey string) Option {
	return func(cfg *Config) error {
		if stacktraceKey == "" {
			return nil
		}
		cfg.EncoderConfig.StacktraceKey = stacktraceKey
		return nil
	}
}

func SkipLineEnding(skipLineEnding bool) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.SkipLineEnding = skipLineEnding
		return nil
	}
}

func LineEnding(lineEnding string) Option {
	return func(cfg *Config) error {
		if lineEnding == "" {
			return nil
		}
		cfg.EncoderConfig.LineEnding = lineEnding
		return nil
	}
}

func LevelEncoder(levelEncoder zapcore.LevelEncoder) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.EncodeLevel = levelEncoder
		return nil
	}
}

func TimeEncoder(timeEncoder zapcore.TimeEncoder) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.EncodeTime = timeEncoder
		return nil
	}
}

func DurationEncoder(durationEncoder zapcore.DurationEncoder) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.EncodeDuration = durationEncoder
		return nil
	}
}

func CallerEncoder(callerEncoder zapcore.CallerEncoder) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.EncodeCaller = callerEncoder
		return nil
	}
}

func NameEncoder(nameEncoder zapcore.NameEncoder) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.EncodeName = nameEncoder
		return nil
	}
}

func Development() Option {
	return func(cfg *Config) error {
		cfg.Development = true
		return nil
	}
}

func Fields(fields map[string]any) Option {
	return func(cfg *Config) error {
		cfg.Fields = fields
		return nil
	}
}

func AddCaller() Option {
	return func(cfg *Config) error {
		cfg.Caller = true
		return nil
	}
}

func WithCaller(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.Caller = enabled
		return nil
	}
}

func Hooks(hooks ...func(zapcore.Entry) error) Option {
	return func(cfg *Config) error {
		cfg.Hooks = append(cfg.Hooks, hooks...)
		return nil
	}
}

func AddCallerSkip(skip int) Option {
	return func(cfg *Config) error {
		cfg.CallerSkip = skip
		return nil
	}
}

func AddStacktrace(level zapcore.LevelEnabler) Option {
	return func(cfg *Config) error {
		cfg.Stacktrace = level
		return nil
	}
}

func IncreaseLevel(level zapcore.LevelEnabler) Option {
	return func(cfg *Config) error {
		cfg.IncreaseLevel = level
		return nil
	}
}

func WithPanicHook(hook zapcore.CheckWriteHook) Option {
	return func(cfg *Config) error {
		cfg.PanicHook = hook
		return nil
	}
}

func WithFatalHook(hook zapcore.CheckWriteHook) Option {
	return func(cfg *Config) error {
		cfg.FatalHook = hook
		return nil
	}
}

func WithClock(clock zapcore.Clock) Option {
	return func(cfg *Config) error {
		cfg.Clock = clock
		return nil
	}
}

func (cfg *Config) ZapEncoder() (zapcore.Encoder, error) {
	var encoder zapcore.Encoder
	if cfg.Encoder == EncoderJSON {
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}
	return encoder, nil
}

func (cfg *Config) ZapWriters() (zapcore.WriteSyncer, error) {
	if len(cfg.Outputs) == 0 {
		return zapcore.AddSync(os.Stdout), nil
	}
	ws := []zapcore.WriteSyncer{}
	for name, options := range cfg.Outputs {
		w, err := logger.Writer(name, options)
		if err != nil {
			return nil, err
		}
		ws = append(ws, zapcore.AddSync(w))
	}
	return zapcore.NewMultiWriteSyncer(ws...), nil
}

func (cfg *Config) ZapOptions() []zap.Option {
	opts := []zap.Option{}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}
	fields := []zapcore.Field{}
	for key, value := range cfg.Fields {
		fields = append(fields, zap.Any(key, value))
	}
	opts = append(opts, zap.Fields(fields...))
	opts = append(opts, zap.WithCaller(cfg.Caller))
	opts = append(opts, zap.AddCallerSkip(cfg.CallerSkip))
	opts = append(opts, zap.Hooks(cfg.Hooks...))
	if cfg.Stacktrace != nil {
		opts = append(opts, zap.AddStacktrace(cfg.Stacktrace))
	}
	if cfg.IncreaseLevel != nil {
		opts = append(opts, zap.IncreaseLevel(cfg.IncreaseLevel))
	}
	if cfg.PanicHook != nil {
		opts = append(opts, zap.WithPanicHook(cfg.PanicHook))
	}
	if cfg.FatalHook != nil {
		opts = append(opts, zap.WithFatalHook(cfg.FatalHook))
	}
	if cfg.Clock != nil {
		opts = append(opts, zap.WithClock(cfg.Clock))
	}
	return opts
}

func UnmarshalOptions(options map[string]any) (*Config, error) {
	cfg := NewConfig()
	err := mapstructure.WeakDecode(options, cfg)
	return cfg, err
}
