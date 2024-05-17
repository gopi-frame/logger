package logger

import (
	"io"

	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/support/maps"
)

var defaultChannel = "default"

var _ logger.Logger = (*Manager)(nil)

// NewManager new manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.log = NewLogger()
	manager.channels = maps.NewMap[string, logger.Logger]()
	manager.channels.Set(defaultChannel, manager.log)
	return manager
}

// Manager logger manager
type Manager struct {
	log      logger.Logger // default logger, channel default
	channels *maps.Map[string, logger.Logger]
}

// Channel channel
func (m *Manager) Channel(name string) logger.Logger {
	return m.channels.GetOr(name, m.log)
}

// AddChannel add channel
func (m *Manager) AddChannel(name string, logger logger.Logger) {
	m.channels.Set(name, logger)
	if name == defaultChannel {
		m.log = logger
	}
}

// Dispatcher set dispatcher for default logger
func (m *Manager) Dispatcher(d event.Dispatcher) {
	m.log.Dispatcher(d)
}

// Formatter set formatter for default logger
func (m *Manager) Formatter(formatter logger.Formatter) {
	m.log.Formatter(formatter)
}

// Hooks set hooks for default logger
func (m *Manager) Hooks(hooks ...logger.Hook) {
	m.log.Hooks(hooks...)
}

// Outputs set outputs for default logger
func (m *Manager) Outputs(outputs ...io.Writer) {
	m.log.Outputs(outputs...)
}

// Debug debug
func (m *Manager) Debug(message string, context map[string]any) {
	m.log.Debug(message, context)
}

// Info info
func (m *Manager) Info(message string, context map[string]any) {
	m.log.Info(message, context)
}

// Warn warn
func (m *Manager) Warn(message string, context map[string]any) {
	m.log.Warn(message, context)
}

// Error error
func (m *Manager) Error(message string, context map[string]any) {
	m.log.Error(message, context)
}

// Fatal fatal
func (m *Manager) Fatal(message string, context map[string]any) {
	m.log.Fatal(message, context)
}

// Panic panic
func (m *Manager) Panic(message string, context map[string]any) {
	m.log.Panic(message, context)
}
