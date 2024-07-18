package zap

import (
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level         zapcore.Level             `json:"level" yaml:"level" toml:"level"`
	Development   bool                      `json:"development" yaml:"development" toml:"development"`
	Fields        map[string]any            `json:"fields" yaml:"fields" toml:"fields"`
	Caller        bool                      `json:"caller" yaml:"caller" toml:"caller"`
	CallerSkip    int                       `json:"callerSkip" yaml:"callerSkip" toml:"callerSkip"`
	Encoder       string                    `json:"encoder" yaml:"encoder" toml:"encoder"`
	EncoderConfig zapcore.EncoderConfig     `json:"encoderConfig" yaml:"encoderConfig" toml:"encoderConfig"`
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
		Level:   zapcore.WarnLevel,
		Encoder: EncoderJSON,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     DefaultEncoderMessageKey,
			LevelKey:       DefaultEncoderLevelKey,
			TimeKey:        DefaultEncoderTimeKey,
			NameKey:        DefaultEncoderNameKey,
			CallerKey:      DefaultEncoderCallerKey,
			FunctionKey:    DefaultEncoderFunctionKey,
			StacktraceKey:  DefaultEncoderStacktraceKey,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
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
	for driverName, options := range cfg.Outputs {
		w, err := writer.Open(driverName, options)
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
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			DecodeLevelHook,
			DecodeLevelEncoderHook,
			DecodeDurationEncoderHook,
			DecodeTimeEncoderHook,
			DecodeCallerEncoderHook,
			DecodeNameEncoderHook,
		),
		WeaklyTypedInput: true,
		Result:           cfg,
		MatchName: func(mapKey, fieldName string) bool {
			if strings.EqualFold(mapKey, fieldName) {
				return true
			}
			if strings.EqualFold(mapKey, "durationEncoder") && fieldName == "EncodeDuration" {
				return true
			}
			if strings.EqualFold(mapKey, "timeEncoder") && fieldName == "EncodeTime" {
				return true
			}
			if strings.EqualFold(mapKey, "levelEncoder") && fieldName == "EncodeLevel" {
				return true
			}
			if strings.EqualFold(mapKey, "callerEncoder") && fieldName == "EncodeCaller" {
				return true
			}
			if strings.EqualFold(mapKey, "nameEncoder") && fieldName == "EncodeName" {
				return true
			}
			return false
		},
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(options); err != nil {
		return nil, err
	}
	return cfg, nil
}

func DecodeLevelHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.Level]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var l = new(zapcore.Level)
		if err := l.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *l, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var l = new(zapcore.Level)
		if err := l.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *l, nil
	}
	if f == reflect.TypeFor[zapcore.Level]() {
		return data.(zapcore.Level), nil
	}
	if f == reflect.TypeFor[zapcore.LevelEnabler]() {
		l := zapcore.LevelOf(data.(zapcore.LevelEnabler))
		if l == zapcore.InvalidLevel {
			return nil, exception.New("failed to unmarshal level")
		}
		return l, nil
	}
	return data, nil
}

func DecodeDurationEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.DurationEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var de = new(zapcore.DurationEncoder)
		if err := de.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *de, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var de = new(zapcore.DurationEncoder)
		if err := de.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *de, nil
	}
	if f == reflect.TypeFor[zapcore.DurationEncoder]() {
		if value, ok := data.(zapcore.DurationEncoder); ok {
			return value, nil
		} else {
			return zapcore.DurationEncoder(data.(func(time.Duration, zapcore.PrimitiveArrayEncoder))), nil
		}
	}
	if f == reflect.TypeFor[func(time.Duration, zapcore.PrimitiveArrayEncoder)]() {
		return zapcore.DurationEncoder(data.(func(time.Duration, zapcore.PrimitiveArrayEncoder))), nil
	}
	return data, nil
}

func DecodeTimeEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.TimeEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var te = new(zapcore.TimeEncoder)
		if err := te.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *te, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var te = new(zapcore.TimeEncoder)
		if err := te.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *te, nil
	}
	if f == reflect.TypeFor[zapcore.TimeEncoder]() {
		if value, ok := data.(zapcore.TimeEncoder); ok {
			return value, nil
		} else {
			return zapcore.TimeEncoder(data.(func(time.Time, zapcore.PrimitiveArrayEncoder))), nil
		}
	}
	if f == reflect.TypeFor[func(time.Time, zapcore.PrimitiveArrayEncoder)]() {
		return zapcore.TimeEncoder(data.(func(time.Time, zapcore.PrimitiveArrayEncoder))), nil
	}
	if f.Kind() == reflect.Map {
		var o struct{ Layout string }
		if err := mapstructure.Decode(data, &o); err != nil {
			return data, err
		}
		return zapcore.TimeEncoderOfLayout(o.Layout), nil
	}
	return data, nil
}

func DecodeLevelEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.LevelEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var le = new(zapcore.LevelEncoder)
		if err := le.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *le, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var le = new(zapcore.LevelEncoder)
		if err := le.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *le, nil
	}
	if f == reflect.TypeFor[zapcore.LevelEncoder]() {
		if value, ok := data.(zapcore.LevelEncoder); ok {
			return value, nil
		} else {
			return zapcore.LevelEncoder(data.(func(zapcore.Level, zapcore.PrimitiveArrayEncoder))), nil
		}
	}
	if f == reflect.TypeFor[func(zapcore.Level, zapcore.PrimitiveArrayEncoder)]() {
		return zapcore.LevelEncoder(data.(func(zapcore.Level, zapcore.PrimitiveArrayEncoder))), nil
	}
	return data, nil
}

func DecodeCallerEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.CallerEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var ce = new(zapcore.CallerEncoder)
		if err := ce.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *ce, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var ce = new(zapcore.CallerEncoder)
		if err := ce.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *ce, nil
	}
	if f == reflect.TypeFor[zapcore.CallerEncoder]() {
		if value, ok := data.(zapcore.CallerEncoder); ok {
			return value, nil
		} else {
			return zapcore.CallerEncoder(data.(func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder))), nil
		}
	}
	if f == reflect.TypeFor[func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder)]() {
		return zapcore.CallerEncoder(data.(func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder))), nil
	}
	return data, nil
}

func DecodeNameEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.NameEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var ne = new(zapcore.NameEncoder)
		if err := ne.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *ne, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var ne = new(zapcore.NameEncoder)
		if err := ne.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *ne, nil
	}
	if f == reflect.TypeFor[zapcore.NameEncoder]() {
		if value, ok := data.(zapcore.NameEncoder); ok {
			return value, nil
		} else {
			return zapcore.NameEncoder(data.(func(string, zapcore.PrimitiveArrayEncoder))), nil
		}
	}
	if f == reflect.TypeFor[func(string, zapcore.PrimitiveArrayEncoder)]() {
		return zapcore.NameEncoder(data.(func(string, zapcore.PrimitiveArrayEncoder))), nil
	}
	return data, nil
}
