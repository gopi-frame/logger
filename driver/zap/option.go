package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
