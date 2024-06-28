package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOptions(t *testing.T) {
	t.Run("unmarshal without perm", func(t *testing.T) {
		var file = filepath.Join(os.TempDir(), "gopi-frame-test", "writer", "driver", "file", fmt.Sprintf("log.%d.txt", time.Now().Unix()))
		var options = map[string]any{
			"file": file,
		}
		config, err := UnmarshalOptions(options)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, file, config.File)
	})

	t.Run("unmarshal with perm", func(t *testing.T) {
		var file = filepath.Join(os.TempDir(), "gopi-frame-test", "writer", "driver", "file", fmt.Sprintf("log.%d.txt", time.Now().Unix()))
		var perm = os.FileMode(0666)
		var options = map[string]any{
			"file": file,
			"perm": perm,
		}
		config, err := UnmarshalOptions(options)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, file, config.File)
		assert.Equal(t, perm, config.Perm)
	})
}

func TestConfig_UnmarshalJSON(t *testing.T) {
	var file = filepath.ToSlash(filepath.Join(os.TempDir(), "gopi-frame-test", "writer", "driver", "file", fmt.Sprintf("log.%d.txt", time.Now().Unix())))
	var config = new(Config)
	jsonbytes, err := json.Marshal(file)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	err = json.Unmarshal(jsonbytes, config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, file, config.File)
}
