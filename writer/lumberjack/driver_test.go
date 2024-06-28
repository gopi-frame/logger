package lumberjack

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gopi-frame/logger"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	var filename = filepath.ToSlash(filepath.Join(os.TempDir(), "gopi-frame-test", "logger", "writer", "driver", "lumberjack", fmt.Sprintf("log.%d.txt", time.Now().Unix())))
	var options = map[string]any{
		"filename": filename,
	}
	w, err := logger.Writer(driverName, options)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	w.Close()
	_, err = w.Write([]byte("hello world"))
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "hello world", string(content))
}
