package logger

import (
	"github.com/gopi-frame/container"
	"github.com/gopi-frame/contract/logger"
)

// Manager is a logger manager
type Manager struct {
	defaultLogger string
	container.Container[logger.Logger]
}

// New creates a new logger manager
func New(defaultLogger string) *Manager {
	return &Manager{
		defaultLogger: defaultLogger,
		Container:     *container.New[logger.Logger](),
	}
}

// Channel gets the channel by name or default logger
func (m *Manager) Channel(name ...string) logger.Logger {
	if len(name) == 0 {
		return m.Get(m.defaultLogger)
	}
	return m.Get(name[0])
}

// Debug logs a debug message
func (m *Manager) Debug(message string) {
	m.Channel().Debug(message)
}

// Info logs an info message
func (m *Manager) Info(message string) {
	m.Channel().Info(message)
}

// Warn logs a warning message
func (m *Manager) Warn(message string) {
	m.Channel().Warn(message)
}

// Error logs an exception message
func (m *Manager) Error(message string) {
	m.Channel().Error(message)
}

// Fatal logs a fatal message
func (m *Manager) Fatal(message string) {
	m.Channel().Fatal(message)
}

// Panic logs a panic message
func (m *Manager) Panic(message string) {
	m.Channel().Panic(message)
}
