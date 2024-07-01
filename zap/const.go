package zap

import (
	"go.uber.org/zap"
)

const (
	DefaultLevel   = zap.WarnLevel
	DefaultEncoder = EncoderJSON

	DefaultEncoderMessageKey    = "message"
	DefaultEncoderLevelKey      = "level"
	DefaultEncoderTimeKey       = "time"
	DefaultEncoderNameKey       = "name"
	DefaultEncoderCallerKey     = "caller"
	DefaultEncoderFunctionKey   = "function"
	DefaultEncoderStacktraceKey = "stacktrace"
)

const (
	EncoderJSON = "json"
	EncoderText = "text"
)

const (
	LevelEncoderCapital        = "capital"
	LevelEncoderCapitalColor   = "capitalColor"
	LevelEncoderLowerCaseColor = "color"
	LevelEncoderLowerCase      = "lowercase"
)

const (
	TimeEncoderRFC3339Nano  = "rfc3339nano"
	TimeEncoderRFC3339      = "rfc3339"
	TimeEncoderISO8601      = "iso8601"
	TimeEncoderMullis       = "millis"
	TimeEncoderNanos        = "nanos"
	TimeEncoderTimestamp    = "timestamp"
	TimeEncoderCustomLayout = "custom"
)

const (
	DurationEncoderSeconds = "seconds"
	DurationEncoderString  = "string"
	DurationEncoderNanos   = "nanos"
	DurationEncoderMillis  = "ms"
)

const (
	CallerEncoderFull  = "full"
	CallerEncoderShort = "short"
)

const (
	NameEncoderFull = "full"
)

const (
	WriterStdout  = "stdout"
	WriterStderr  = "stderr"
	WriterDiscard = "discard"
)
