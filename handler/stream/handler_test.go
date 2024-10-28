package stream

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNewStreamHandlerFromConfig(t *testing.T) {
	t.Run("stdout", func(t *testing.T) {
		_ = os.MkdirAll("testdata", 0755)
		f, err := os.OpenFile("testdata/stdout.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		originalStdout := os.Stdout
		defer func() {
			os.Stdout = originalStdout
		}()
		os.Stdout = f
		handler, err := NewStreamHandlerFromConfig(map[string]any{
			"stream": "stdout",
		})
		defer func() {
			_ = os.Remove("testdata/stdout.log")
		}()
		_, err = handler.Write([]byte("test"))
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		_ = f.Close()
		content, err := os.ReadFile("testdata/stdout.log")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "test", string(content))
	})

	t.Run("stderr", func(t *testing.T) {
		_ = os.MkdirAll("testdata", 0755)
		f, err := os.OpenFile("testdata/stderr.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		originalStderr := os.Stderr
		defer func() {
			os.Stderr = originalStderr
		}()
		os.Stderr = f
		handler, err := NewStreamHandlerFromConfig(map[string]any{
			"stream": "stderr",
		})
		defer func() {
			_ = os.Remove("testdata/stderr.log")
		}()
		_, err = handler.Write([]byte("test"))
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		_ = f.Close()
		content, err := os.ReadFile("testdata/stderr.log")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "test", string(content))
	})

	t.Run("discard", func(t *testing.T) {
		var buf = bytes.NewBuffer(nil)
		originalDiscard := io.Discard
		defer func() {
			io.Discard = originalDiscard
		}()
		io.Discard = buf
		handler, err := NewStreamHandlerFromConfig(map[string]any{
			"stream": "discard",
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		_, err = handler.Write([]byte("test"))
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "test", buf.String())
	})

	t.Run("file", func(t *testing.T) {
		_ = os.MkdirAll("testdata", 0755)
		handler, err := NewStreamHandlerFromConfig(map[string]any{
			"stream": "file://testdata/file.log",
		})
		defer func() {
			if err := handler.Close(); err != nil {
				assert.FailNow(t, err.Error())
			} else {
				_ = os.Remove("testdata/file.log")
			}
		}()
		_, err = handler.Write([]byte("test"))
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		content, err := os.ReadFile("testdata/file.log")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "test", string(content))
	})
}
