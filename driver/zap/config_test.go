package zap

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type enc struct {
	zapcore.PrimitiveArrayEncoder
	value string
}

func (enc *enc) AppendString(str string) {
	enc.value += str
}

func TestDecodeLevelHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var level = "warn"
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.Level](), level)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevel)
	})

	t.Run("decode invalid string", func(t *testing.T) {
		var level = "invalid"
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.Level](), level)
		assert.Error(t, err)
		assert.Equal(t, nil, decodedLevel)
	})

	t.Run("decode bytes", func(t *testing.T) {
		var level = []byte("warn")
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.Level](), level)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevel)
	})

	t.Run("decode invalid bytes", func(t *testing.T) {
		var level = []byte("invalid")
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.Level](), level)
		assert.Error(t, err)
		assert.Equal(t, nil, decodedLevel)
	})

	t.Run("decode zapcore.Level", func(t *testing.T) {
		var level = zapcore.WarnLevel
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[zapcore.Level](), reflect.TypeFor[zapcore.Level](), level)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevel)
	})

	t.Run("decode zapcore.LevelEnabler", func(t *testing.T) {
		var level = zapcore.WarnLevel
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[zapcore.LevelEnabler](), reflect.TypeFor[zapcore.Level](), level)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevel)
	})

	t.Run("decode invalid zapcore.LevelEnabler", func(t *testing.T) {
		var level = zapcore.Level(7)
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[zapcore.LevelEnabler](), reflect.TypeFor[zapcore.Level](), level)
		assert.Error(t, err)
		assert.Equal(t, nil, decodedLevel)
	})
}

func TestDecodeDurationEncoderHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var durationEncoder = DurationEncoderSeconds
		decodedDurationEncoder, err := DecodeDurationEncoderHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.DurationEncoder](), durationEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.SecondsDurationEncoder).Pointer(), reflect.ValueOf(decodedDurationEncoder).Pointer())
	})

	t.Run("decode bytes", func(t *testing.T) {
		var durationEncoder = []byte(DurationEncoderSeconds)
		decodedDurationEncoder, err := DecodeDurationEncoderHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.DurationEncoder](), durationEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.SecondsDurationEncoder).Pointer(), reflect.ValueOf(decodedDurationEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var durationEncoder = func(time.Duration, zapcore.PrimitiveArrayEncoder) {}
		decodedDurationEncoder, err := DecodeDurationEncoderHook(reflect.TypeFor[func(time.Duration, zapcore.PrimitiveArrayEncoder)](), reflect.TypeFor[zapcore.DurationEncoder](), durationEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(durationEncoder).Pointer(), reflect.ValueOf(decodedDurationEncoder).Pointer())
	})

	t.Run("decode zapcore.DurationEncoder", func(t *testing.T) {
		var durationEncoder = zapcore.SecondsDurationEncoder
		decodedDurationEncoder, err := DecodeDurationEncoderHook(reflect.TypeFor[zapcore.DurationEncoder](), reflect.TypeFor[zapcore.DurationEncoder](), durationEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.SecondsDurationEncoder).Pointer(), reflect.ValueOf(decodedDurationEncoder).Pointer())
	})
}

func TestTimeEncoderHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var timeEncoder = TimeEncoderRFC3339
		decodedTimeEncoder, err := DecodeTimeEncoderHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.TimeEncoder](), timeEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.RFC3339TimeEncoder).Pointer(), reflect.ValueOf(decodedTimeEncoder).Pointer())
	})

	t.Run("decode bytes", func(t *testing.T) {
		var timeEncoder = []byte(TimeEncoderRFC3339)
		decodedTimeEncoder, err := DecodeTimeEncoderHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.TimeEncoder](), timeEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.RFC3339TimeEncoder).Pointer(), reflect.ValueOf(decodedTimeEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var timeEncoder = func(time.Time, zapcore.PrimitiveArrayEncoder) {}
		decodedTimeEncoder, err := DecodeTimeEncoderHook(reflect.TypeFor[func(time.Time, zapcore.PrimitiveArrayEncoder)](), reflect.TypeFor[zapcore.TimeEncoder](), timeEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(timeEncoder).Pointer(), reflect.ValueOf(decodedTimeEncoder).Pointer())
	})

	t.Run("decode zapcore.TimeEncoder", func(t *testing.T) {
		var timeEncoder = zapcore.RFC3339TimeEncoder
		decodedTimeEncoder, err := DecodeTimeEncoderHook(reflect.TypeFor[zapcore.TimeEncoder](), reflect.TypeFor[zapcore.TimeEncoder](), timeEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.RFC3339TimeEncoder).Pointer(), reflect.ValueOf(decodedTimeEncoder).Pointer())
	})

	t.Run("decode layout", func(t *testing.T) {
		var enc = new(enc)
		var timeEncoder = map[string]string{"layout": "2006-01-02 15:04:05"}
		decodedTimeEncoder, err := DecodeTimeEncoderHook(reflect.TypeFor[map[string]string](), reflect.TypeFor[zapcore.TimeEncoder](), timeEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		var now = time.Now()
		(decodedTimeEncoder.(zapcore.TimeEncoder))(now, enc)
		assert.Equal(t, now.Format("2006-01-02 15:04:05"), enc.value)
	})
}

func TestLevelEncoderHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var levelEncoder = LevelEncoderCapital
		decodedLevelEncoder, err := DecodeLevelEncoderHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.LevelEncoder](), levelEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.CapitalLevelEncoder).Pointer(), reflect.ValueOf(decodedLevelEncoder).Pointer())
	})

	t.Run("decode bytes", func(t *testing.T) {
		var levelEncoder = []byte(LevelEncoderCapital)
		decodedLevelEncoder, err := DecodeLevelEncoderHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.LevelEncoder](), levelEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.CapitalLevelEncoder).Pointer(), reflect.ValueOf(decodedLevelEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var levelEncoder = func(zapcore.Level, zapcore.PrimitiveArrayEncoder) {}
		decodedLevelEncoder, err := DecodeLevelEncoderHook(reflect.TypeFor[func(zapcore.Level, zapcore.PrimitiveArrayEncoder)](), reflect.TypeFor[zapcore.LevelEncoder](), levelEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(levelEncoder).Pointer(), reflect.ValueOf(decodedLevelEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var levelEncoder = zapcore.CapitalLevelEncoder
		decodedLevelEncoder, err := DecodeLevelEncoderHook(reflect.TypeFor[zapcore.LevelEncoder](), reflect.TypeFor[zapcore.LevelEncoder](), levelEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.CapitalLevelEncoder).Pointer(), reflect.ValueOf(decodedLevelEncoder).Pointer())
	})
}

func TestDecodeCallerEncoderHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var callerEncoder = CallerEncoderShort
		decodedCallerEncoder, err := DecodeCallerEncoderHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.CallerEncoder](), callerEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.ShortCallerEncoder).Pointer(), reflect.ValueOf(decodedCallerEncoder).Pointer())
	})

	t.Run("decode bytes", func(t *testing.T) {
		var callerEncoder = []byte(CallerEncoderShort)
		decodedCallerEncoder, err := DecodeCallerEncoderHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.CallerEncoder](), callerEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.ShortCallerEncoder).Pointer(), reflect.ValueOf(decodedCallerEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var callerEncoder = func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder) {}
		decodedCallerEncoder, err := DecodeCallerEncoderHook(reflect.TypeFor[func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder)](), reflect.TypeFor[zapcore.CallerEncoder](), callerEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(callerEncoder).Pointer(), reflect.ValueOf(decodedCallerEncoder).Pointer())
	})

	t.Run("decode zapcore.CallerEncoder", func(t *testing.T) {
		var callerEncoder = zapcore.ShortCallerEncoder
		decodedCallerEncoder, err := DecodeCallerEncoderHook(reflect.TypeFor[zapcore.CallerEncoder](), reflect.TypeFor[zapcore.CallerEncoder](), callerEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.ShortCallerEncoder).Pointer(), reflect.ValueOf(decodedCallerEncoder).Pointer())
	})
}

func TestDecodeNameEncoderHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var nameEncoder = NameEncoderFull
		decodedNameEncoder, err := DecodeNameEncoderHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.NameEncoder](), nameEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.FullNameEncoder).Pointer(), reflect.ValueOf(decodedNameEncoder).Pointer())
	})

	t.Run("decode byte", func(t *testing.T) {
		var nameEncoder = []byte(NameEncoderFull)
		decodedNameEncoder, err := DecodeNameEncoderHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.NameEncoder](), nameEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.FullNameEncoder).Pointer(), reflect.ValueOf(decodedNameEncoder).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var nameEncoder = func(string, zapcore.PrimitiveArrayEncoder) {}
		decodedNameEncoder, err := DecodeNameEncoderHook(reflect.TypeFor[func(string, zapcore.PrimitiveArrayEncoder)](), reflect.TypeFor[zapcore.NameEncoder](), nameEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(nameEncoder).Pointer(), reflect.ValueOf(decodedNameEncoder).Pointer())
	})

	t.Run("decode zapcore.NameEncoder", func(t *testing.T) {
		var nameEncoder = zapcore.FullNameEncoder
		decodedNameEncoder, err := DecodeNameEncoderHook(reflect.TypeFor[zapcore.NameEncoder](), reflect.TypeFor[zapcore.NameEncoder](), nameEncoder)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.FullNameEncoder).Pointer(), reflect.ValueOf(decodedNameEncoder).Pointer())
	})
}

func TestDecodeLevelEnablerHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var levelEnabler = "warn"
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevelEnabler)
	})

	t.Run("decode invalid string", func(t *testing.T) {
		var levelEnabler = "invalid"
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		assert.Error(t, err)
		assert.Nil(t, decodedLevelEnabler)
	})

	t.Run("decode bytes", func(t *testing.T) {
		var levelEnabler = []byte("warn")
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevelEnabler)
	})

	t.Run("decode invalid bytes", func(t *testing.T) {
		var levelEnabler = []byte("invalid")
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		assert.Error(t, err)
		assert.Nil(t, decodedLevelEnabler)
	})

	t.Run("decode zapcore.LevelEnabler", func(t *testing.T) {
		var levelEnabler = zapcore.InfoLevel
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[zapcore.Level](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.InfoLevel, decodedLevelEnabler)
	})

	t.Run("decode zap.LevelEnablerFunc", func(t *testing.T) {
		var levelEnabler = zap.LevelEnablerFunc(func(zapcore.Level) bool { return true })
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[zap.LevelEnablerFunc](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(levelEnabler).Pointer(), reflect.ValueOf(decodedLevelEnabler).Pointer())
	})

	t.Run("decode func", func(t *testing.T) {
		var levelEnabler = func(zapcore.Level) bool { return true }
		decodedLevelEnabler, err := DecodeLevelEnablerHook(reflect.TypeFor[func(zapcore.Level) bool](), reflect.TypeFor[zapcore.LevelEnabler](), levelEnabler)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(levelEnabler).Pointer(), reflect.ValueOf(decodedLevelEnabler).Pointer())
	})
}

func TestDecodeCheckWriteHook(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		var checkWrite = "noop"
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenNoop, decodedCheckWrite)
	})

	t.Run("decode invalid string", func(t *testing.T) {
		var checkWrite = "invalid"
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[string](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		assert.Error(t, err)
		assert.Nil(t, decodedCheckWrite)
	})

	t.Run("decode bytes", func(t *testing.T) {
		var checkWrite = []byte("noop")
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenNoop, decodedCheckWrite)
	})

	t.Run("decode invalid bytes", func(t *testing.T) {
		var checkWrite = []byte("invalid")
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		assert.Error(t, err)
		assert.Nil(t, decodedCheckWrite)
	})

	t.Run("decode zapcore.CheckWriteHook", func(t *testing.T) {
		var checkWrite = zapcore.WriteThenNoop
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[zapcore.CheckWriteHook](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenNoop, decodedCheckWrite)
	})

	t.Run("decode number", func(t *testing.T) {
		var checkWrite = 1
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[int](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenGoexit, decodedCheckWrite)
	})

	t.Run("decode func", func(t *testing.T) {
		var checkWrite = func(*zapcore.CheckedEntry, []zapcore.Field) {}
		decodedCheckWrite, err := DecodeCheckWriteHook(reflect.TypeFor[func(*zapcore.CheckedEntry, []zapcore.Field)](), reflect.TypeFor[zapcore.CheckWriteHook](), checkWrite)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(checkWrite).Pointer(), reflect.ValueOf(decodedCheckWrite.(*checkWriteAction).callback).Pointer())
	})
}

func TestConfig_Apply(t *testing.T) {
	t.Run("Level", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, zapcore.WarnLevel, config.Level)
		if err := config.Apply(Level(zapcore.InfoLevel)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.InfoLevel, config.Level)
	})

	t.Run("Encoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, EncoderJSON, config.Encoder)

		if err := config.Apply(Encoder("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, EncoderJSON, config.Encoder)

		if err := config.Apply(Encoder(EncoderText)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, EncoderText, config.Encoder)
	})

	t.Run("MessageKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "message", config.EncoderConfig.MessageKey)
		if err := config.Apply(MessageKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "message", config.EncoderConfig.MessageKey)
		if err := config.Apply(MessageKey("msg")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "msg", config.EncoderConfig.MessageKey)
	})

	t.Run("LevelKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "level", config.EncoderConfig.LevelKey)
		if err := config.Apply(LevelKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "level", config.EncoderConfig.LevelKey)
		if err := config.Apply(LevelKey("lvl")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "lvl", config.EncoderConfig.LevelKey)
	})

	t.Run("TimeKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "time", config.EncoderConfig.TimeKey)
		if err := config.Apply(TimeKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "time", config.EncoderConfig.TimeKey)
		if err := config.Apply(TimeKey("ts")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "ts", config.EncoderConfig.TimeKey)
	})

	t.Run("NameKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "name", config.EncoderConfig.NameKey)
		if err := config.Apply(NameKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "name", config.EncoderConfig.NameKey)
		if err := config.Apply(NameKey("customName")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "customName", config.EncoderConfig.NameKey)
	})

	t.Run("CallerKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "caller", config.EncoderConfig.CallerKey)
		if err := config.Apply(CallerKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "caller", config.EncoderConfig.CallerKey)
		if err := config.Apply(CallerKey("customCaller")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "customCaller", config.EncoderConfig.CallerKey)
	})

	t.Run("FunctionKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "function", config.EncoderConfig.FunctionKey)
		if err := config.Apply(FunctionKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "function", config.EncoderConfig.FunctionKey)
		if err := config.Apply(FunctionKey("customFunction")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "customFunction", config.EncoderConfig.FunctionKey)
	})

	t.Run("StacktraceKey", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "stacktrace", config.EncoderConfig.StacktraceKey)
		if err := config.Apply(StacktraceKey("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "stacktrace", config.EncoderConfig.StacktraceKey)
		if err := config.Apply(StacktraceKey("customStacktrace")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "customStacktrace", config.EncoderConfig.StacktraceKey)
	})

	t.Run("SkipLineEnding", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, false, config.EncoderConfig.SkipLineEnding)
		if err := config.Apply(SkipLineEnding(true)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, true, config.EncoderConfig.SkipLineEnding)
	})

	t.Run("LineEnding", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, "\n", config.EncoderConfig.LineEnding)
		if err := config.Apply(LineEnding("")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "\n", config.EncoderConfig.LineEnding)
		if err := config.Apply(LineEnding("customLineEnding")); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "customLineEnding", config.EncoderConfig.LineEnding)
	})

	t.Run("LevelEncoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, reflect.ValueOf(zapcore.LowercaseLevelEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeLevel).Pointer())
		if err := config.Apply(LevelEncoder(zapcore.CapitalLevelEncoder)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.CapitalLevelEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeLevel).Pointer())
	})

	t.Run("TimeEncoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, reflect.ValueOf(zapcore.RFC3339TimeEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeTime).Pointer())
		if err := config.Apply(TimeEncoder(zapcore.EpochTimeEncoder)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.EpochTimeEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeTime).Pointer())
	})

	t.Run("DurationEncoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, reflect.ValueOf(zapcore.StringDurationEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeDuration).Pointer())
		if err := config.Apply(DurationEncoder(zapcore.NanosDurationEncoder)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.NanosDurationEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeDuration).Pointer())
	})

	t.Run("CallerEncoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, reflect.ValueOf(zapcore.ShortCallerEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeCaller).Pointer())
		if err := config.Apply(CallerEncoder(zapcore.FullCallerEncoder)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.FullCallerEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeCaller).Pointer())
	})

	t.Run("NameEncoder", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, reflect.ValueOf(zapcore.FullNameEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeName).Pointer())
		if err := config.Apply(NameEncoder(nil)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, reflect.ValueOf(zapcore.FullNameEncoder).Pointer(), reflect.ValueOf(config.EncoderConfig.EncodeName).Pointer())
	})

	t.Run("Development", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, false, config.Development)
		if err := config.Apply(Development()); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, true, config.Development)
	})

	t.Run("Fields", func(t *testing.T) {
		config := NewConfig()
		assert.Nil(t, config.Fields)
		if err := config.Apply(Fields(map[string]any{"key": "value"})); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, map[string]any{"key": "value"}, config.Fields)
	})

	t.Run("AddCaller", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, false, config.Caller)
		if err := config.Apply(AddCaller()); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, true, config.Caller)
	})

	t.Run("Hooks", func(t *testing.T) {
		config := NewConfig()
		assert.Nil(t, config.Hooks)
		var hooks []func(zapcore.Entry) error
		var hook = func(e zapcore.Entry) error { return nil }
		hooks = append(hooks, hook)
		if err := config.Apply(Hooks(hooks...)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Len(t, config.Hooks, 1)
		assert.Equal(t, reflect.ValueOf(hook).Pointer(), reflect.ValueOf(config.Hooks[0]).Pointer())
	})

	t.Run("AddCallerSkip", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, 0, config.CallerSkip)
		if err := config.Apply(AddCallerSkip(1)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, 1, config.CallerSkip)
	})

	t.Run("AddStacktrace", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, zapcore.ErrorLevel, config.Stacktrace)
		if err := config.Apply(AddStacktrace(zapcore.WarnLevel)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, config.Stacktrace)
	})

	t.Run("IncreaseLevel", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, nil, config.IncreaseLevel)
		if err := config.Apply(IncreaseLevel(zapcore.ErrorLevel)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.ErrorLevel, config.IncreaseLevel)
	})

	t.Run("WithPanicHook", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, nil, config.PanicHook)
		if err := config.Apply(WithPanicHook(zapcore.WriteThenPanic)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenPanic, config.PanicHook)
	})

	t.Run("WithFatalHook", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, nil, config.FatalHook)
		if err := config.Apply(WithFatalHook(zapcore.WriteThenFatal)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WriteThenFatal, config.FatalHook)
	})

	t.Run("WithClock", func(t *testing.T) {
		config := NewConfig()
		assert.Equal(t, nil, config.Clock)
		if err := config.Apply(WithClock(zapcore.DefaultClock)); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.DefaultClock, config.Clock)
	})
}
