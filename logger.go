package logger

import "github.com/gopi-frame/logger/contract"

type logger struct {
	logger  contract.Logger
	options Options
}

func (w logger) open() contract.Logger {
	if w.logger != nil {
		return w.logger
	}
	driver := w.options.Driver()
	logger, err := Open(driver, w.options)
	if err != nil {
		panic(err)
	}
	w.logger = logger
	return w.logger
}

func active(l contract.Logger) logger {
	return logger{
		logger: l,
	}
}

func lazy(options Options) logger {
	return logger{
		options: options,
	}
}

func (l logger) Debug(message string, fields map[string]any) {
	l.open().Debug(message, fields)
}

func (l logger) Info(message string, fields map[string]any) {
	l.open().Info(message, fields)
}

func (l logger) Warn(message string, fields map[string]any) {
	l.open().Warn(message, fields)
}

func (l logger) Error(message string, fields map[string]any) {
	l.open().Error(message, fields)
}

func (l logger) Fatal(message string, fields map[string]any) {
	l.open().Fatal(message, fields)
}

func (l logger) Panic(message string, fields map[string]any) {
	l.open().Panic(message, fields)
}
