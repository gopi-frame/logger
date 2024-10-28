package zap

import (
	"bytes"
	"fmt"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/exception"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is the configuration for zap logger.
type Config struct {
	Level         zapcore.Level         `json:"level" yaml:"level" toml:"level"`
	Development   bool                  `json:"development" yaml:"development" toml:"development"`
	Fields        map[string]any        `json:"fields" yaml:"fields" toml:"fields"`
	Caller        bool                  `json:"caller" yaml:"caller" toml:"caller"`
	CallerSkip    int                   `json:"callerSkip" yaml:"callerSkip" toml:"callerSkip"`
	Encoder       string                `json:"encoder" yaml:"encoder" toml:"encoder"`
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig" toml:"encoderConfig"`
	Handler       string                `json:"handler" yaml:"handler" toml:"handler"`
	HandlerWith   map[string]any        `json:"handlerWith" yaml:"handlerWith" toml:"handlerWith"`

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
	if cfg.Handler != "" {
		handler, err := logger.CreateHandler(cfg.Handler, cfg.HandlerWith)
		if err != nil {
			return nil, err
		}
		return zapcore.Lock(zapcore.AddSync(handler)), nil
	}
	return zapcore.Lock(zapcore.AddSync(os.Stdout)), nil
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
			env.ExpandStringWithEnvHookFunc(),
			env.ExpandSliceWithEnvHookFunc(),
			env.ExpandStringKeyMapWithEnvHookFunc(),
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
			if strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", "")) {
				return true
			}
			if (strings.EqualFold(mapKey, "durationEncoder") ||
				strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))) && fieldName == "EncodeDuration" {
				return true
			}
			if (strings.EqualFold(mapKey, "timeEncoder") ||
				strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))) && fieldName == "EncodeTime" {
				return true
			}
			if (strings.EqualFold(mapKey, "levelEncoder") ||
				strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))) && fieldName == "EncodeLevel" {
				return true
			}
			if (strings.EqualFold(mapKey, "callerEncoder") ||
				strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))) && fieldName == "EncodeCaller" {
				return true
			}
			if (strings.EqualFold(mapKey, "nameEncoder") ||
				strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))) && fieldName == "EncodeName" {
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
