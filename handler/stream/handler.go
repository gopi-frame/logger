package stream

import (
	"github.com/gopi-frame/logger"
	"github.com/spf13/cast"
	"io"
	"os"
	"strings"
)

var handlerName = "stream"

//goland:noinspection GoBoolExpressions
func init() {
	if handlerName != "" {
		logger.RegisterHandler(handlerName, func(config map[string]any) (io.WriteCloser, error) {
			return NewStreamHandlerFromConfig(config)
		})
	}
}

type StreamHandler struct {
	w io.Writer
}

func NewStreamHandler(w io.Writer) *StreamHandler {
	return &StreamHandler{
		w: w,
	}
}

func NewStreamHandlerFromConfig(config map[string]any) (*StreamHandler, error) {
	stream, ok := config["stream"]
	if !ok || stream == nil {
		return NewStreamHandler(io.Discard), nil
	}
	if stream, ok := stream.(io.Writer); ok {
		return NewStreamHandler(stream), nil
	}
	if stream, ok := stream.(string); ok {
		switch stream {
		case "stdout":
			return NewStreamHandler(os.Stdout), nil
		case "stderr":
			return NewStreamHandler(os.Stderr), nil
		case "discard", "null":
			return NewStreamHandler(io.Discard), nil
		default:
			if strings.HasPrefix(stream, "file://") {
				var mode uint32 = 0644
				if m, ok := config["mode"]; ok {
					mode = cast.ToUint32(m)
				}
				stream, err := os.OpenFile(strings.TrimLeft(stream, "file://"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(mode))
				if err != nil {
					panic(err)
				}
				return NewStreamHandler(stream), nil
			}
		}
	}
	return nil, NewInvalidStreamException()
}

func (h *StreamHandler) Write(p []byte) (n int, err error) {
	return h.w.Write(p)
}

func (h *StreamHandler) Close() error {
	if h.w == nil {
		return nil
	}
	if closer, ok := h.w.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return err
		}
		h.w = nil
	}
	return nil
}
