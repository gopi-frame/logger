package zap

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	_ "github.com/gopi-frame/writer/file"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var file = filepath.ToSlash(filepath.Join(os.TempDir(), "gopi-frame-test", "logger", "zap", "log.txt."+strconv.FormatInt(time.Now().Unix(), 10)))
	var options = map[string]any{
		"level": "info",
		"fields": map[string]any{
			"driver": "zap",
		},
		"caller":  true,
		"encoder": "json",
		"encoderConfig": map[string]any{
			"messageKey": "message",
			"timeEncoder": map[string]string{
				"layout": "2006-01-02 15:04:05",
			},
		},
		"outputs": map[string]any{
			"file": map[string]any{
				"file": file,
			},
		},
	}
	logger, err := new(Driver).Open(options)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	logger.Debug("debug", nil)
	logger.Info("info", nil)
	logger.Warn("warn", map[string]any{"extra": "error"})
	logger.Error("error", nil)
}
