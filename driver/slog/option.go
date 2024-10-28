package slog

import "log/slog"

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
