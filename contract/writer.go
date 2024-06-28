package contract

import "io"

type Writer interface {
	Open(map[string]any) (io.WriteCloser, error)
}
