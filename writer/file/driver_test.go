package file

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
	var file = filepath.ToSlash(filepath.Join(os.TempDir(), "gopi-frame-test", "writer", "driver", "file", fmt.Sprintf("log.%d.txt", time.Now().Unix())))
	var options = map[string]any{
		"file": file,
	}
	w, err := logger.Writer(driverName, options)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	_, err = w.Write([]byte("hello world"))
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	w.Close()
	content, err := os.ReadFile(file)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "hello world", string(content))
}
