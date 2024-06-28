package logger

import (
	"io"
	"os"
)

func init() {
	RegisterWriter(WriterDiscard, new(discard))
	RegisterWriter(WriterStdout, new(stdout))
	RegisterWriter(WriterStderr, new(stderr))
}

type writeCloser struct {
	io.Writer
}

func (w writeCloser) Close() error {
	return nil
}

func AddClose(writer io.Writer) io.WriteCloser {
	return writeCloser{Writer: writer}
}

type discard struct{}

func (discard) Open(map[string]any) (io.WriteCloser, error) {
	return AddClose(io.Discard), nil
}

type stdout struct{}

func (stdout) Open(map[string]any) (io.WriteCloser, error) {
	return os.Stdout, nil
}

type stderr struct{}

func (stderr) Open(map[string]any) (io.WriteCloser, error) {
	return os.Stderr, nil
}
