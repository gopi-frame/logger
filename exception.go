package logger

import (
	. "github.com/gopi-frame/contract/exception"
	"github.com/gopi-frame/exception"
)

type UnknownDriverException struct {
	Throwable
}

func NewUnknownDriverException(driverName string) *UnknownDriverException {
	return &UnknownDriverException{
		Throwable: exception.New("unknown driver [%s]", driverName),
	}
}

type UnknownHandlerException struct {
	Throwable
}

func NewUnknownHandlerException(handlerName string) *UnknownHandlerException {
	return &UnknownHandlerException{
		Throwable: exception.New("unknown handler [%s]", handlerName),
	}
}

type NotConfiguredChannelException struct {
	Throwable
}

func NewNotConfiguredChannelException(channel string) *NotConfiguredChannelException {
	return &NotConfiguredChannelException{
		Throwable: exception.New("channel [%s] not configured", channel),
	}
}
