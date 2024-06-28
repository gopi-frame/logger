package zap

import "go.uber.org/zap"

const (
	OptKeyLevel                   = "level"
	OptKeyEncoder                 = "encoder"
	OptKeyEncoderConfig           = "encoderConfig"
	OptKeyEncoderMessageKey       = "messageKey"
	OptKeyEncoderLevelKey         = "levelKey"
	OptKeyEncoderTimeKey          = "timeKey"
	OptKeyEncoderNameKey          = "nameKey"
	OptKeyEncoderCallerKey        = "callerKey"
	OptKeyEncoderFunctionKey      = "functionKey"
	OptKeyEncoderStacktraceKey    = "stacktraceKey"
	OptKeyEncoderSkipLineEnding   = "skipLineEnding"
	OptKeyEncoderLineEnding       = "lineEnding"
	OptKeyEncoderLevelEncoder     = "levelEncoder"
	OptKeyEncoderTimeEncoder      = "timeEncoder"
	OptKeyEncoderTimeLayout       = "timeLayout"
	OptKeyEncoderDurationEncoder  = "durationEncoder"
	OptKeyEncoderCallerEncoder    = "callerEncoder"
	OptKeyEncoderNameEncoder      = "nameEncoder"
	OptKeyEncoderConsoleSeparator = "consoleSeparator"
	OptKeyWriters                 = "writers"
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

	DefaultLevelEncoder    = LevelEncoderLowerCase
	DefaultTimeEncoder     = TimeEncoderTimestamp
	DefaultTimeLayout      = "2006-01-02 15:04:05"
	DefaultDurationEncoder = DurationEncoderString
	DefaultCallerEncoder   = CallerEncoderFull
	DefaultNameEncoder     = NameEncoderFull
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
