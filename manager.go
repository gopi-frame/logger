package logger

import (
	"fmt"

	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/logger/contract"
	"github.com/gopi-frame/support/maps"
)

type Options map[string]any

func (o Options) Driver() string {
	return o[OptKeyDriver].(string)
}

type Manager struct {
	defaultLogger string
	loggers       *maps.Map[string, logger]
}

func NewManager() *Manager {
	return &Manager{
		loggers: maps.NewMap[string, logger](),
	}
}

func (m *Manager) SetDefault(name string) {
	m.defaultLogger = name
}

func (m *Manager) Add(name string, logger contract.Logger) {
	m.loggers.Lock()
	defer m.loggers.Unlock()
	if m.loggers.ContainsKey(name) {
		panic(exception.NewArgumentException("name", name, fmt.Sprintf("duplicate logger \"%s\"", name)))
	}
	m.loggers.Set(name, active(logger))
}

func (m *Manager) AddLazy(name string, options Options) {
	m.loggers.Lock()
	defer m.loggers.Unlock()
	if m.loggers.ContainsKey(name) {
		panic(exception.NewArgumentException("name", name, fmt.Sprintf("duplicate logger \"%s\"", name)))
	}
	m.loggers.Set(name, lazy(options))
}

func (m *Manager) Logger(name string) contract.Logger {
	m.loggers.RLock()
	defer m.loggers.RUnlock()
	logger, ok := m.loggers.Get(name)
	if ok {
		return logger
	}
	logger, ok = m.loggers.Get(m.defaultLogger)
	if ok {
		return logger
	}
	panic(exception.NewArgumentException("name", name, fmt.Sprintf("unknown logger \"%s\"", name)))
}

func (m *Manager) Debug(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Debug(message, fields)
}

func (m *Manager) Info(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Info(message, fields)
}

func (m *Manager) Warn(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Warn(message, fields)
}

func (m *Manager) Error(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Error(message, fields)
}

func (m *Manager) Fatal(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Fatal(message, fields)
}

func (m *Manager) Panic(message string, fields map[string]any) {
	m.Logger(m.defaultLogger).Panic(message, fields)
}
