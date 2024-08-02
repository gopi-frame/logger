package zap

import (
	"bytes"
	"fmt"
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

// DefaultWriter is the default writer for zap logger.
var DefaultWriter zapcore.WriteSyncer = os.Stdout

// Config is the configuration for zap logger.
type Config struct {
	Level         zapcore.Level             `json:"level" yaml:"level" toml:"level"`
	Development   bool                      `json:"development" yaml:"development" toml:"development"`
	Fields        map[string]any            `json:"fields" yaml:"fields" toml:"fields"`
	Caller        bool                      `json:"caller" yaml:"caller" toml:"caller"`
	CallerSkip    int                       `json:"callerSkip" yaml:"callerSkip" toml:"callerSkip"`
	Encoder       string                    `json:"encoder" yaml:"encoder" toml:"encoder"`
	EncoderConfig zapcore.EncoderConfig     `json:"encoderConfig" yaml:"encoderConfig" toml:"encoderConfig"`
	Writers       map[string]map[string]any `json:"writers" yaml:"writers" toml:"writers"`

	Hooks         []func(zapcore.Entry) error `json:"-" yaml:"-" toml:"-"`
	Stacktrace    zapcore.LevelEnabler        `json:"-" yaml:"-" toml:"-"`
	IncreaseLevel zapcore.LevelEnabler        `json:"-" yaml:"-" toml:"-"`
	PanicHook     zapcore.CheckWriteHook      `json:"-" yaml:"-" toml:"-"`
	FatalHook     zapcore.CheckWriteHook      `json:"-" yaml:"-" toml:"-"`
	Clock         zapcore.Clock               `json:"-" yaml:"-" toml:"-"`
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{
		Level:   DefaultLevel,
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
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			LineEnding:     zapcore.DefaultLineEnding,
			SkipLineEnding: false,
		},
		Stacktrace: zapcore.ErrorLevel,
	}
}

// Option is a function that can be used to configure zap logger.
type Option func(cfg *Config) error

// Apply applies the given options to the [Config].
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}

// Level sets the log level.
func Level(level zapcore.Level) Option {
	return func(cfg *Config) error {
		cfg.Level = level
		return nil
	}
}

// Encoder sets the log encoder.
// AvailableAt encoders: [EncoderJSON], [EncoderText].
// If empty string is given, it does nothing.
func Encoder(level string) Option {
	return func(cfg *Config) error {
		if level == "" {
			return nil
		}
		cfg.Encoder = level
		return nil
	}
}

// MessageKey sets the message key.
// If empty string is given, it does nothing.
func MessageKey(messageKey string) Option {
	return func(cfg *Config) error {
		if messageKey == "" {
			return nil
		}
		cfg.EncoderConfig.MessageKey = messageKey
		return nil
	}
}

// LevelKey sets the level key.
// If empty string is given, it does nothing.
func LevelKey(levelKey string) Option {
	return func(cfg *Config) error {
		if levelKey == "" {
			return nil
		}
		cfg.EncoderConfig.LevelKey = levelKey
		return nil
	}
}

// TimeKey sets the time key.
// If empty string is given, it does nothing.
func TimeKey(timeKey string) Option {
	return func(cfg *Config) error {
		if timeKey == "" {
			return nil
		}
		cfg.EncoderConfig.TimeKey = timeKey
		return nil
	}
}

// NameKey sets the name key.
// If empty string is given, it does nothing.
func NameKey(nameKey string) Option {
	return func(cfg *Config) error {
		if nameKey == "" {
			return nil
		}
		cfg.EncoderConfig.NameKey = nameKey
		return nil
	}
}

// CallerKey sets the caller key.
// If empty string is given, it does nothing.
func CallerKey(callerKey string) Option {
	return func(cfg *Config) error {
		if callerKey == "" {
			return nil
		}
		cfg.EncoderConfig.CallerKey = callerKey
		return nil
	}
}

// FunctionKey sets the function key.
// If empty string is given, it does nothing.
func FunctionKey(functionKey string) Option {
	return func(cfg *Config) error {
		if functionKey == "" {
			return nil
		}
		cfg.EncoderConfig.FunctionKey = functionKey
		return nil
	}
}

// StacktraceKey sets the stacktrace key.
// If empty string is given, it does nothing.
func StacktraceKey(stacktraceKey string) Option {
	return func(cfg *Config) error {
		if stacktraceKey == "" {
			return nil
		}
		cfg.EncoderConfig.StacktraceKey = stacktraceKey
		return nil
	}
}

// SkipLineEnding sets whether to skip line ending in the log.
func SkipLineEnding(skipLineEnding bool) Option {
	return func(cfg *Config) error {
		cfg.EncoderConfig.SkipLineEnding = skipLineEnding
		return nil
	}
}

// LineEnding sets the line ending in the log.
// If empty string is given, it does nothing.
func LineEnding(lineEnding string) Option {
	return func(cfg *Config) error {
		if lineEnding == "" {
			return nil
		}
		cfg.EncoderConfig.LineEnding = lineEnding
		return nil
	}
}

// LevelEncoder sets the level encoder.
// If nil is given, it does nothing.
func LevelEncoder(levelEncoder zapcore.LevelEncoder) Option {
	return func(cfg *Config) error {
		if levelEncoder == nil {
			return nil
		}
		cfg.EncoderConfig.EncodeLevel = levelEncoder
		return nil
	}
}

// TimeEncoder sets the time encoder.
// If nil is given, it does nothing.
func TimeEncoder(timeEncoder zapcore.TimeEncoder) Option {
	return func(cfg *Config) error {
		if timeEncoder == nil {
			return nil
		}
		cfg.EncoderConfig.EncodeTime = timeEncoder
		return nil
	}
}

// DurationEncoder sets the duration encoder.
// If nil is given, it does nothing.
func DurationEncoder(durationEncoder zapcore.DurationEncoder) Option {
	return func(cfg *Config) error {
		if durationEncoder == nil {
			return nil
		}
		cfg.EncoderConfig.EncodeDuration = durationEncoder
		return nil
	}
}

// CallerEncoder sets the caller encoder.
// If nil is given, it does nothing.
func CallerEncoder(callerEncoder zapcore.CallerEncoder) Option {
	return func(cfg *Config) error {
		if callerEncoder == nil {
			return nil
		}
		cfg.EncoderConfig.EncodeCaller = callerEncoder
		return nil
	}
}

// NameEncoder sets the name encoder.
// If nil is given, it does nothing.
func NameEncoder(nameEncoder zapcore.NameEncoder) Option {
	return func(cfg *Config) error {
		if nameEncoder == nil {
			return nil
		}
		cfg.EncoderConfig.EncodeName = nameEncoder
		return nil
	}
}

// Development sets the development mode.
func Development() Option {
	return func(cfg *Config) error {
		cfg.Development = true
		return nil
	}
}

// Fields sets the extra fields to the log.
func Fields(fields map[string]any) Option {
	return func(cfg *Config) error {
		cfg.Fields = fields
		return nil
	}
}

// AddCaller enables caller.
func AddCaller() Option {
	return WithCaller(true)
}

// WithCaller enables or disables caller.
func WithCaller(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.Caller = enabled
		return nil
	}
}

// Hooks adds hooks.
func Hooks(hooks ...func(zapcore.Entry) error) Option {
	return func(cfg *Config) error {
		cfg.Hooks = append(cfg.Hooks, hooks...)
		return nil
	}
}

// AddCallerSkip adds caller skip.
func AddCallerSkip(skip int) Option {
	return func(cfg *Config) error {
		cfg.CallerSkip = skip
		return nil
	}
}

// AddStacktrace adds stacktrace to the log if the level is enabled.
// If nil is given, it does nothing.
// For more details, see https://pkg.go.dev/go.uber.org/zap#AddStacktrace.
func AddStacktrace(level zapcore.LevelEnabler) Option {
	return func(cfg *Config) error {
		if cfg.Stacktrace == nil {
			cfg.Stacktrace = zap.ErrorLevel
		}
		cfg.Stacktrace = level
		return nil
	}
}

// IncreaseLevel increases the level of the log.
// For more details, see https://pkg.go.dev/go.uber.org/zap#IncreaseLevel.
func IncreaseLevel(level zapcore.LevelEnabler) Option {
	return func(cfg *Config) error {
		cfg.IncreaseLevel = level
		return nil
	}
}

// WithPanicHook sets the panic hook.
// For more details, see https://pkg.go.dev/go.uber.org/zap#WithPanicHook.
func WithPanicHook(hook zapcore.CheckWriteHook) Option {
	return func(cfg *Config) error {
		cfg.PanicHook = hook
		return nil
	}
}

// WithFatalHook sets the fatal hook.
// For more details, see https://pkg.go.dev/go.uber.org/zap#WithFatalHook.
func WithFatalHook(hook zapcore.CheckWriteHook) Option {
	return func(cfg *Config) error {
		cfg.FatalHook = hook
		return nil
	}
}

// WithClock sets the clock.
// For more details, see https://pkg.go.dev/go.uber.org/zap#WithClock.
func WithClock(clock zapcore.Clock) Option {
	return func(cfg *Config) error {
		cfg.Clock = clock
		return nil
	}
}

// ZapEncoder returns the zap encoder.
func (cfg *Config) ZapEncoder() (zapcore.Encoder, error) {
	var encoder zapcore.Encoder
	if cfg.Encoder == EncoderJSON {
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}
	return encoder, nil
}

// ZapWriters returns the zap writers.
func (cfg *Config) ZapWriters() (zapcore.WriteSyncer, error) {
	if len(cfg.Writers) == 0 {
		return DefaultWriter, nil
	}
	var ws []zapcore.WriteSyncer
	for driverName, options := range cfg.Writers {
		w, err := writer.Open(driverName, options)
		if err != nil {
			return nil, err
		}
		ws = append(ws, zapcore.AddSync(w))
	}
	return zapcore.Lock(zapcore.NewMultiWriteSyncer(ws...)), nil
}

// ZapOptions returns the zap options.
func (cfg *Config) ZapOptions() []zap.Option {
	var opts []zap.Option
	if cfg.Development {
		opts = append(opts, zap.Development())
	}
	var fields []zapcore.Field
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

// UnmarshalOptions unmarshal the options.
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
			DecodeLevelEnablerHook,
			DecodeCheckWriteHook,
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

// DecodeLevelHook decodes the level.
//
//   - For type string or []byte, [zapcore.Level.UnmarshalText] will be called to parse the value.
//   - For type which is convertible to int8, it will convert to the [zapcore.Level] by [reflect.Type.Convert].
//   - For type [zapcore.Level], it will return the value directly.
//   - For type [zapcore.LevelEnabler], it will call [zapcore.LevelOf] to get the level.
func DecodeLevelHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.Level]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var l = new(zapcore.Level)
		if err := l.UnmarshalText([]byte(strings.ToUpper(data.(string)))); err != nil {
			return nil, err
		}
		return *l, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var l = new(zapcore.Level)
		if err := l.UnmarshalText(bytes.ToUpper(data.([]byte))); err != nil {
			return nil, err
		}
		return *l, nil
	}
	if f.ConvertibleTo(reflect.TypeFor[int8]()) {
		l := zapcore.Level(reflect.ValueOf(data).Convert(reflect.TypeFor[int8]()).Interface().(int8))
		if l >= zapcore.DebugLevel && l <= zapcore.FatalLevel {
			return l, nil
		}
		return nil, exception.NewArgumentException("data", data, "invalid level value")
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

// DecodeDurationEncoderHook decodes the duration encoder.
//
//   - For type string or []byte, [zapcore.DurationEncoder.UnmarshalText] will be called to parse the value.
//   - For type [zapcore.DurationEncoder], it will return the value directly.
//   - For type func([time.Duration], [zapcore.PrimitiveArrayEncoder]),
//     it will convert the value to [zapcore.DurationEncoder] and return it.
func DecodeDurationEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.DurationEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var de = new(zapcore.DurationEncoder)
		_ = de.UnmarshalText([]byte(data.(string)))
		return *de, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var de = new(zapcore.DurationEncoder)
		_ = de.UnmarshalText(data.([]byte))
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

// DecodeTimeEncoderHook decodes the time encoder.
//
//   - For type string or []byte, [zapcore.TimeEncoder.UnmarshalText] will be called to parse the value.
//   - For type [zapcore.TimeEncoder], it will return the value directly.
//   - For type func([time.Time], [zapcore.PrimitiveArrayEncoder]),
//     it will convert the value to [zapcore.TimeEncoder] and return it.
//   - For type map which contains a key named "layout" (case-insensitive),
//     [zapcore.TimeEncoderOfLayout] will be called to create a time encoder with given layout.
func DecodeTimeEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.TimeEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var te = new(zapcore.TimeEncoder)
		_ = te.UnmarshalText([]byte(strings.ToLower(data.(string))))
		return *te, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var te = new(zapcore.TimeEncoder)
		_ = te.UnmarshalText(bytes.ToLower(data.([]byte)))
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

// DecodeLevelEncoderHook decodes the level encoder.
//
//   - For type string or []byte, [zapcore.LevelEncoder.UnmarshalText] will be called to parse the value.
//   - For type [zapcore.LevelEncoder], it will return the value directly.
//   - For type func([zapcore.Level], [zapcore.PrimitiveArrayEncoder]), it will convert the value to [zapcore.LevelEncoder] and return it.
func DecodeLevelEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.LevelEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var le = new(zapcore.LevelEncoder)
		_ = le.UnmarshalText([]byte(data.(string)))
		return *le, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var le = new(zapcore.LevelEncoder)
		_ = le.UnmarshalText(data.([]byte))
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

// DecodeCallerEncoderHook decodes the caller encoder.
//
//   - For type string or []byte, [zapcore.CallerEncoder.UnmarshalText] will be called to parse the value.
//   - For type [zapcore.CallerEncoder], it will return the value directly.
//   - For type func([zapcore.EntryCaller], [zapcore.PrimitiveArrayEncoder]), it will convert the value to [zapcore.CallerEncoder] and return it.
func DecodeCallerEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.CallerEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var ce = new(zapcore.CallerEncoder)
		_ = ce.UnmarshalText([]byte(data.(string)))
		return *ce, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var ce = new(zapcore.CallerEncoder)
		_ = ce.UnmarshalText(data.([]byte))
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

// DecodeNameEncoderHook decodes the name encoder.
//
//   - For type string or []byte, [zapcore.NameEncoder.UnmarshalText] will be called to parse the value.
//   - For type [zapcore.NameEncoder], it will return the value directly.
//   - For type func(string, [zapcore.PrimitiveArrayEncoder]), it will convert the value to [zapcore.NameEncoder] and return it.
func DecodeNameEncoderHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.NameEncoder]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var ne = new(zapcore.NameEncoder)
		_ = ne.UnmarshalText([]byte(data.(string)))
		return *ne, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var ne = new(zapcore.NameEncoder)
		_ = ne.UnmarshalText(data.([]byte))
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

// DecodeLevelEnablerHook decodes the level enabler.
//
//   - For type string or []byte, [zapcore.Level.UnmarshalText] will be called to parse the value.
//   - For type which can be converted to int8, it will convert the value to [zapcore.Level] and return it.
//   - For type which implements [zapcore.LevelEnabler], it will return the value directly.
//   - For type func([zapcore.Level]) bool, it will convert the value to [zapcore.LevelEnabler] and return it.
func DecodeLevelEnablerHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.LevelEnabler]() {
		return data, nil
	}
	if f == reflect.TypeFor[string]() {
		var le = new(zapcore.Level)
		if err := le.UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return *le, nil
	}
	if f == reflect.TypeFor[[]byte]() {
		var le = new(zapcore.Level)
		if err := le.UnmarshalText(data.([]byte)); err != nil {
			return nil, err
		}
		return *le, nil
	}
	if f.ConvertibleTo(reflect.TypeFor[int8]()) {
		le := zapcore.Level(reflect.ValueOf(data).Convert(reflect.TypeFor[int8]()).Interface().(int8))
		if le >= zapcore.DebugLevel && le <= zapcore.PanicLevel {
			return le, nil
		}
		return nil, exception.NewArgumentException("data", data, "invalid value")
	}
	if f.Implements(reflect.TypeFor[zapcore.LevelEnabler]()) {
		return data.(zapcore.LevelEnabler), nil
	}
	if f == reflect.TypeFor[zap.LevelEnablerFunc]() {
		if value, ok := data.(zap.LevelEnablerFunc); ok {
			return value, nil
		} else {
			return zap.LevelEnablerFunc(data.(func(zapcore.Level) bool)), nil
		}
	}
	if f == reflect.TypeFor[func(zapcore.Level) bool]() {
		return zap.LevelEnablerFunc(data.(func(zapcore.Level) bool)), nil
	}
	return data, nil
}

// CheckWriteHookFunc returns a [zapcore.CheckWriteHook] by the given function.
func CheckWriteHookFunc(fn func(entry *zapcore.CheckedEntry, fields []zapcore.Field)) zapcore.CheckWriteHook {
	return &checkWriteAction{
		callback: fn,
	}
}

type checkWriteAction struct {
	callback func(entry *zapcore.CheckedEntry, fields []zapcore.Field)
}

func (c *checkWriteAction) OnWrite(entry *zapcore.CheckedEntry, fields []zapcore.Field) {
	c.callback(entry, fields)
}

// DecodeCheckWriteHook decodes the check write hook.
//
//	For type string, it will be converted to [zapcore.CheckWriteHook] by the following rules:
//	  1. "0" or "noop" will be converted to [zapcore.WriteThenNoop].
//	  2. "1" or "goexit" will be converted to [zapcore.WriteThenGoexit].
//	  3. "2" or "panic" will be converted to [zapcore.WriteThenPanic].
//	  4. "3" or "fatal" will be converted to [zapcore.WriteThenFatal].
//	For type [zapcore.CheckWriteHook], it will return the value directly.
//	For type func([zapcore.CheckedEntry], []zapcore.Field), it will convert the value to [zapcore.CheckWriteHook] and return it.
//	For type which can be converted to uint8, it will be converted to [zapcore.CheckWriteHook] by the following rules:
//	  1. 0 will be converted to [zapcore.WriteThenNoop].
//	  2. 1 will be converted to [zapcore.WriteThenGoexit].
//	  3. 2 will be converted to [zapcore.WriteThenPanic].
//	  4. 3 will be converted to [zapcore.WriteThenFatal].
func DecodeCheckWriteHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeFor[zapcore.CheckWriteHook]() {
		return data, nil
	}
	if f.Implements(reflect.TypeFor[zapcore.CheckWriteHook]()) {
		return data.(zapcore.CheckWriteHook), nil
	}
	if f == reflect.TypeFor[string]() {
		switch strings.ToLower(data.(string)) {
		case "0", "noop":
			return zapcore.WriteThenNoop, nil
		case "1", "goexit":
			return zapcore.WriteThenGoexit, nil
		case "2", "panic":
			return zapcore.WriteThenPanic, nil
		case "3", "fatal":
			return zapcore.WriteThenFatal, nil
		default:
			return nil, exception.NewArgumentException("data", data, "invalid value for CheckWriteHook: "+data.(string))
		}
	}
	if f == reflect.TypeFor[[]byte]() {
		switch strings.ToLower(string(data.([]byte))) {
		case "0", "noop":
			return zapcore.WriteThenNoop, nil
		case "1", "goexit":
			return zapcore.WriteThenGoexit, nil
		case "2", "panic":
			return zapcore.WriteThenPanic, nil
		case "3", "fatal":
			return zapcore.WriteThenFatal, nil
		default:
			return nil, exception.NewArgumentException("data", data, "invalid value for CheckWriteHook: "+string(data.([]byte)))
		}
	}
	if f.ConvertibleTo(reflect.TypeFor[uint8]()) {
		switch reflect.ValueOf(data).Convert(reflect.TypeFor[uint8]()).Interface().(uint8) {
		case 0:
			return zapcore.WriteThenNoop, nil
		case 1:
			return zapcore.WriteThenGoexit, nil
		case 2:
			return zapcore.WriteThenPanic, nil
		case 3:
			return zapcore.WriteThenFatal, nil
		default:
			return nil, exception.NewArgumentException("data", data, fmt.Sprintf("invalid value for CheckWriteHook: %v", data))
		}
	}
	if f == reflect.TypeFor[func(*zapcore.CheckedEntry, []zapcore.Field)]() {
		return CheckWriteHookFunc(data.(func(*zapcore.CheckedEntry, []zapcore.Field))), nil
	}
	return data, nil
}
