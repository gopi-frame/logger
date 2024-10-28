package daily

import (
	"compress/gzip"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gopi-frame/env"
	"github.com/gopi-frame/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var handlerName = "daily"

//goland:noinspection GoBoolExpressions
func init() {
	if handlerName != "" {
		logger.RegisterHandler(handlerName, func(config map[string]any) (io.WriteCloser, error) {
			return NewDailyHandlerFromConfig(config)
		})
	}
}

type DailyHandler struct {
	filename string
	mode     os.FileMode
	dir      string
	current  string
	maxAge   int
	compress bool
	now      func() time.Time // for testing
}

// NewDailyHandler creates a new daily log handler.
func NewDailyHandler(filename string, opts ...Option) (*DailyHandler, error) {
	handler := &DailyHandler{
		filename: filepath.Base(filename),
		dir:      filepath.Dir(filename),
		now:      time.Now,
	}
	if err := os.MkdirAll(handler.dir, 0755); err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler, nil
}

func NewDailyHandlerFromConfig(config map[string]any) (*DailyHandler, error) {
	var cfg struct {
		Filename string
		Mode     uint32
		MaxAge   int
		Compress bool
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
		return nil, err
	}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return NewDailyHandler(cfg.Filename, WithFileMode(os.FileMode(cfg.Mode)), WithMaxAge(cfg.MaxAge), WithCompress(cfg.Compress))
}

func (h *DailyHandler) openExistingOrNew() (*os.File, error) {
	filename := filepath.Join(h.dir,
		fmt.Sprintf("%s.%s%s",
			filepath.Base(h.filename[:len(h.filename)-len(filepath.Ext(h.filename))]),
			h.now().Format("2006-01-02"),
			filepath.Ext(h.filename)))

	if h.current != filename {
		if h.current != "" && h.compress {
			_ = h.compressFile(h.current)
		}
		h.current = filename
		if h.maxAge > 0 {
			go func() {
				_ = h.cleanOldFiles()
			}()
		}
	}
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, h.mode)
}

func (h *DailyHandler) cleanOldFiles() error {
	if h.maxAge <= 0 {
		return nil
	}
	entries, err := os.ReadDir(h.dir)
	if err != nil {
		return err
	}
	cutoff := h.now().AddDate(0, 0, -h.maxAge)
	base := filepath.Base(h.filename)
	ext := filepath.Ext(base)
	prefix := base[:len(base)-len(ext)]
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		// Parse date from filename
		dateStr := strings.TrimPrefix(name, prefix+".")
		dateStr = strings.TrimSuffix(dateStr, ext)
		dateStr = strings.TrimSuffix(dateStr, ext+".gz")
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if fileDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(h.dir, name))
		}
	}
	return nil
}

func (h *DailyHandler) compressFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	gzFile := file + ".gz"
	gzf, err := os.Create(gzFile)
	if err != nil {
		return err
	}
	defer func(gzf *os.File) {
		_ = gzf.Close()
	}(gzf)
	gzw := gzip.NewWriter(gzf)
	defer func(gzw *gzip.Writer) {
		_ = gzw.Close()
	}(gzw)
	if _, err := io.Copy(gzw, f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Remove(file)
}

func (h *DailyHandler) Write(p []byte) (int, error) {
	file, err := h.openExistingOrNew()
	if err != nil {
		return 0, err
	}
	n, err := file.Write(p)
	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	return n, err
}

func (h *DailyHandler) Close() error {
	return nil
}
