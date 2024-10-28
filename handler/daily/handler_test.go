package daily

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewDailyHandlerFromConfig(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		_ = os.Setenv("TEST_DATA_DIR", "testdata")
		fh, err := NewDailyHandlerFromConfig(map[string]any{
			"filename": "${TEST_DATA_DIR}/test.log",
			"mode":     0644,
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		now := time.Now()
		fh.now = func() time.Time {
			return now
		}
		for i := 0; i < 10; i++ {
			_, err := fh.Write([]byte("test"))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			content, err := os.ReadFile(fmt.Sprintf("testdata/test.%s.log", now.Format("2006-01-02")))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, "test", string(content))
			now = now.Add(time.Hour * 24)
		}
		entries, err := os.ReadDir("testdata")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, 10, len(entries))
		_ = os.RemoveAll("testdata")
	})

	t.Run("with compress", func(t *testing.T) {
		_ = os.Setenv("TEST_DATA_DIR", "testdata")
		fh, err := NewDailyHandlerFromConfig(map[string]any{
			"filename": "${TEST_DATA_DIR}/test.log",
			"mode":     0644,
			"compress": true,
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		now := time.Now()
		fh.now = func() time.Time {
			return now
		}
		for i := 0; i < 10; i++ {
			_, err := fh.Write([]byte("test"))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			content, err := os.ReadFile(fmt.Sprintf("testdata/test.%s.log", now.Format("2006-01-02")))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, "test", string(content))
			if i > 0 {
				_, err := os.Stat(fmt.Sprintf("testdata/test.%s.log.gz", now.Add(-time.Hour*24).Format("2006-01-02")))
				if !assert.NoError(t, err) {
					assert.FailNow(t, err.Error())
				}
			}
			now = now.Add(time.Hour * 24)
		}
		entries, err := os.ReadDir("testdata")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, 10, len(entries))
		count := 0
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(entry.Name(), ".gz") {
				count++
			}
		}
		assert.Equal(t, 9, count)
		_ = os.RemoveAll("testdata")
	})

	t.Run("with maxAge", func(t *testing.T) {
		_ = os.Setenv("TEST_DATA_DIR", "testdata")
		fh, err := NewDailyHandlerFromConfig(map[string]any{
			"filename": "${TEST_DATA_DIR}/test.log",
			"mode":     0644,
			"maxAge":   7,
			"compress": true,
		})
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		now := time.Now()
		fh.now = func() time.Time {
			return now
		}
		for i := 0; i < 10; i++ {
			_, err := fh.Write([]byte("test"))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			content, err := os.ReadFile(fmt.Sprintf("testdata/test.%s.log", now.Format("2006-01-02")))
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, "test", string(content))
			if i > 0 {
				_, err := os.Stat(fmt.Sprintf("testdata/test.%s.log.gz", now.Add(-time.Hour*24).Format("2006-01-02")))
				if !assert.NoError(t, err) {
					assert.FailNow(t, err.Error())
				}
			}
			now = now.Add(time.Hour * 24)
		}
		entries, err := os.ReadDir("testdata")
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, 8, len(entries))
		count := 0
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.HasSuffix(entry.Name(), ".gz") {
				count++
			}
		}
		assert.Equal(t, 7, count)
		_ = os.RemoveAll("testdata")
	})
}
