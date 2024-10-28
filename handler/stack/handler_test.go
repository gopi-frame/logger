package stack

import (
	"bytes"
	"errors"
	"github.com/gopi-frame/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type mockHandler1 struct {
	bytes.Buffer
}

func (m *mockHandler1) Write(p []byte) (n int, err error) {
	return m.Buffer.Write(p)
}

func (m *mockHandler1) Close() error {
	return nil
}

type mockHandler2 struct {
	bytes.Buffer
}

func (m *mockHandler2) Write(_ []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (m *mockHandler2) Close() error {
	return nil
}

type mockHandler3 struct {
	bytes.Buffer
}

func (m *mockHandler3) Write(p []byte) (n int, err error) {
	return m.Buffer.Write(p)
}

func (m *mockHandler3) Close() error {
	return nil
}

func TestNewStackHandlerFromConfig(t *testing.T) {
	logger.RegisterHandler("mock1", func(config map[string]any) (io.WriteCloser, error) {
		return &mockHandler1{}, nil
	})
	logger.RegisterHandler("mock2", func(config map[string]any) (io.WriteCloser, error) {
		return &mockHandler2{}, nil
	})
	logger.RegisterHandler("mock3", func(config map[string]any) (io.WriteCloser, error) {
		return &mockHandler3{}, nil
	})
	t.Run("break on error", func(t *testing.T) {
		var handler *StackHandler
		handler, err := NewStackHandlerFromConfig(map[string]any{
			"break_on_error": true,
			"handlers": []map[string]any{
				{
					"driver": "mock1",
				},
				{
					"driver": "mock2",
				},
				{
					"driver": "mock3",
				},
			},
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		defer func() {
			err := handler.Close()
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
		}()
		_, err = handler.Write([]byte("test"))
		if assert.Error(t, err) {
			assert.Equal(t, "error", err.Error())
			assert.Equal(t, "test", handler.handlers[0].(*mockHandler1).String())
			assert.Equal(t, "", handler.handlers[1].(*mockHandler2).String())
			assert.Equal(t, "", handler.handlers[2].(*mockHandler3).String())
		}
	})

	t.Run("continue on error", func(t *testing.T) {
		handler, err := NewStackHandlerFromConfig(map[string]any{
			"break_on_error": false,
			"handlers": []map[string]any{
				{
					"driver": "mock1",
				},
				{
					"driver": "mock2",
				},
				{
					"driver": "mock3",
				},
			},
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		defer func() {
			err := handler.Close()
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
		}()
		_, err = handler.Write([]byte("test"))
		if assert.Error(t, err) {
			assert.Equal(t, "error", err.Error())
			assert.Equal(t, "test", handler.handlers[0].(*mockHandler1).String())
			assert.Equal(t, "", handler.handlers[1].(*mockHandler2).String())
			assert.Equal(t, "test", handler.handlers[2].(*mockHandler3).String())
		}
	})
}
