package zap

import (
	"go.uber.org/zap"
)

const (
	// DefaultLevel is the default logging level.
	DefaultLevel = zap.WarnLevel
	// DefaultEncoder is the default logging encoder.
	DefaultEncoder = EncoderJSON

	// DefaultEncoderMessageKey is the default message key for encoder.
	DefaultEncoderMessageKey = "message"
	// DefaultEncoderLevelKey is the default level key for encoder.
	DefaultEncoderLevelKey = "level"
	// DefaultEncoderTimeKey is the default time key for encoder.
	DefaultEncoderTimeKey = "time"
	// DefaultEncoderNameKey is the default name key for encoder.
	DefaultEncoderNameKey = "name"
	// DefaultEncoderCallerKey     = "caller"
	DefaultEncoderCallerKey = "caller"
	// DefaultEncoderFunctionKey   = "function"
	DefaultEncoderFunctionKey = "function"
	// DefaultEncoderStacktraceKey = "stacktrace"
	DefaultEncoderStacktraceKey = "stacktrace"
)

// Encoder type enums
const (
	EncoderJSON = "json"
	EncoderText = "text"
)

// Level encoder enums
const (
	LevelEncoderCapital        = "capital"
	LevelEncoderCapitalColor   = "capitalColor"
	LevelEncoderLowerCaseColor = "color"
	LevelEncoderLowerCase      = "lowercase"
)

// Time encoder enums
const (
	TimeEncoderRFC3339Nano = "rfc3339nano"
	TimeEncoderRFC3339     = "rfc3339"
	TimeEncoderISO8601     = "iso8601"
	TimeEncoderMillis      = "millis"
	TimeEncoderNanos       = "nanos"
	TimeEncoderTimestamp   = "timestamp"
)

// Duration encoder enums
const (
	DurationEncoderSeconds = "seconds"
	DurationEncoderString  = "string"
	DurationEncoderNanos   = "nanos"
	DurationEncoderMillis  = "ms"
)

// Caller encoder enums
const (
	CallerEncoderFull  = "full"
	CallerEncoderShort = "short"
)

// Name encoder enums
const (
	NameEncoderFull = "full"
)
