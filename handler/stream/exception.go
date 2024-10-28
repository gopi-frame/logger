package stream

import (
	. "github.com/gopi-frame/contract/exception"
	"github.com/gopi-frame/exception"
)

type InvalidStreamException struct {
	Throwable
}

func NewStreamMissingException() *InvalidStreamException {
	return &InvalidStreamException{
		Throwable: exception.New("stream is missing in configuration"),
	}
}

func NewInvalidStreamException() *InvalidStreamException {
	return &InvalidStreamException{
		exception.New("stream value is invalid in configuration"),
	}
}
