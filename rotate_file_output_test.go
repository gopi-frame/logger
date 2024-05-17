package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRotateFileOutput(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "logger_test")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tempDir)
	}()
	rotateFile := NewRotateFileOutput(
		filepath.Join(tempDir, "test.log"),
		1,
		1,
		10,
		true,
		false,
		"* * * * * *",
	)
	for i := 0; i < 5; i++ {
		rotateFile.Write([]byte("testlogger"))
		time.Sleep(time.Second)
	}
	entries, err := os.ReadDir(tempDir)
	assert.Nil(t, err)
	assert.Equal(t, 6, len(entries))
}
