package file

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var handlerName = "file"

//goland:noinspection GoBoolExpressions
func init() {
	if handlerName != "" {
		logger.RegisterHandler(handlerName, func(config map[string]any) (io.WriteCloser, error) {
			return NewFileHandlerFromConfig(config)
		})
	}
}

type FileHandler struct {
	filename string
	mode     os.FileMode
}

func NewFileHandler(filename string, mode os.FileMode) (*FileHandler, error) {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return nil, err
	}
	return &FileHandler{
		filename: filename,
		mode:     mode,
	}, nil
}

func NewFileHandlerFromConfig(config map[string]any) (*FileHandler, error) {
	var cfg struct {
		Filename string
		Mode     uint32
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &cfg,
		WeaklyTypedInput: true,
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName) || strings.EqualFold(fieldName, strings.ReplaceAll(mapKey, "_", ""))
		},
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			env.ExpandStringWithEnvHookFunc(),
			mapstructure.StringToBasicTypeHookFunc(),
		),
	})
	if err != nil {
		panic(err)
	}
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}
	return NewFileHandler(cfg.Filename, os.FileMode(cfg.Mode))
}

func (h *FileHandler) Write(p []byte) (int, error) {
	file, err := os.OpenFile(h.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, h.mode)
	if err != nil {
		return 0, err
	}
	n, err := file.Write(p)
	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	return n, err
}

func (h *FileHandler) Close() error {
	return nil
}
