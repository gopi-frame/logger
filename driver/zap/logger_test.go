package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/gopi-frame/logger/writer/file"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	var file = filepath.Join(os.TempDir(), "gopi-frame-logger-test-1", fmt.Sprintf("log.%d.txt", time.Now().Unix()))
	var testOptions = map[string]any{
		OptKeyLevel:   "warn",
		OptKeyEncoder: EncoderJSON,
		OptKeyEncoderConfig: map[string]any{
			OptKeyEncoderMessageKey:       DefaultEncoderMessageKey,
			OptKeyEncoderLevelKey:         DefaultEncoderLevelKey,
			OptKeyEncoderTimeKey:          DefaultEncoderTimeKey,
			OptKeyEncoderNameKey:          DefaultEncoderNameKey,
			OptKeyEncoderCallerKey:        DefaultEncoderCallerKey,
			OptKeyEncoderFunctionKey:      DefaultEncoderFunctionKey,
			OptKeyEncoderStacktraceKey:    DefaultEncoderStacktraceKey,
			OptKeyEncoderSkipLineEnding:   false,
			OptKeyEncoderLevelEncoder:     DefaultLevelEncoder,
			OptKeyEncoderTimeLayout:       "2006-01-02",
			OptKeyEncoderDurationEncoder:  DefaultDurationEncoder,
			OptKeyEncoderCallerEncoder:    DefaultCallerEncoder,
			OptKeyEncoderNameEncoder:      DefaultNameEncoder,
			OptKeyEncoderConsoleSeparator: zapcore.DefaultLineEnding,
		},
		OptKeyWriters: map[string]any{
			"file": map[string]any{
				"file": file,
			},
		},
	}
	logger, err := new(Driver).Open(testOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", logger)
	// assert.Nil(t, err)
	// logger.Debug("this is a debug message", map[string]any{})
	// logger.Info("this is a info message", nil)
	// logger.Warn("this is a warn message", nil)
	// logger.Error("this is an error message", nil)
	// logger.Error("this is an error message with extra fields", map[string]any{
	// 	"err": "error",
	// })
	// _, err = os.Open(file)
	// fmt.Println(err.Error())
}
