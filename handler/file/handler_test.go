package file

import (
	"github.com/gopi-frame/env"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewFileHandler(t *testing.T) {
	fh, err := NewFileHandler("testdata/test.log", 0644)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	_, err = fh.Write([]byte("test"))
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	defer func() {
		_ = os.Remove("testdata/test.log")
	}()
	content, err := os.ReadFile("testdata/test.log")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", string(content))
	}
}

func TestNewFileHandlerFromConfig(t *testing.T) {
	_ = env.Set("TEST_DATA_DIR", "testdata")
	fh, err := NewFileHandlerFromConfig(map[string]any{
		"filename": "${TEST_DATA_DIR}/test.log",
		"mode":     0644,
	})
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	_, err = fh.Write([]byte("test"))
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	defer func() {
		_ = os.Remove("testdata/test.log")
	}()
	content, err := os.ReadFile("testdata/test.log")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", string(content))
	}
}
