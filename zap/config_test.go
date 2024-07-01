package zap

import (
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

	t.Run("decode bytes", func(t *testing.T) {
		var level = []byte("warn")
		decodedLevel, err := DecodeLevelHook(reflect.TypeFor[[]byte](), reflect.TypeFor[zapcore.Level](), level)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, zapcore.WarnLevel, decodedLevel)
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
