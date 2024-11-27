package logger

import (
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/logger"
)

// LoggerManager is a logger manager and a proxy for the default logger.
type LoggerManager struct {
	logger.Logger
	defaultChannel string
	channels       *kv.Map[string, logger.Logger]
}

// NewLoggerManager creates a new logger manager.
func NewLoggerManager() *LoggerManager {
	return &LoggerManager{
		channels: kv.NewMap[string, logger.Logger](),
	}
}

func (m *LoggerManager) init() {
	if m.Logger == nil {
		m.Logger = m.getConnectedChannel(m.defaultChannel)
	}
}

func (m *LoggerManager) getConnectedChannel(name string) logger.Logger {
	m.channels.RLock()
	defer m.channels.RUnlock()
	if channel, ok := m.channels.Get(name); ok {
		return channel
	}
	return nil
}

// SetDefault sets the default channel name.
func (m *LoggerManager) SetDefault(name string) {
	m.defaultChannel = name
}

// HasChannel returns true if the channel with the given name exists.
func (m *LoggerManager) HasChannel(name string) bool {
	m.channels.RLock()
	defer m.channels.RUnlock()
	return m.channels.ContainsKey(name)
}

// SetChannel sets the channel with the given name.
func (m *LoggerManager) SetChannel(name string, channel logger.Logger) {
	m.channels.Lock()
	defer m.channels.Unlock()
	m.channels.Set(name, channel)
}

// GetChannel returns the channel with the given name.
func (m *LoggerManager) GetChannel(name string) logger.Logger {
	if channel := m.getConnectedChannel(name); channel != nil {
		return channel
	}
	return m
}

// RemoveChannel removes the channel with the given name.
func (m *LoggerManager) RemoveChannel(name string) {
	m.channels.Lock()
	defer m.channels.Unlock()
	m.channels.Remove(name)
}

func (m *LoggerManager) GetChannels() map[string]logger.Logger {
	m.channels.RLock()
	defer m.channels.RUnlock()
	return m.channels.ToMap()
}

// GetChannelBundle returns a stack of the given channels.
// If the channel is not found, it skips the channel.
func (m *LoggerManager) GetChannelBundle(names ...string) logger.Logger {
	var channels []logger.Logger
	for _, name := range names {
		channel := m.getConnectedChannel(name)
		if channel == nil {
			continue
		}
		channels = append(channels, channel)
	}
	return NewStackLogger(channels...)
}

func (m *LoggerManager) Debug(message string) {
	m.init()
	m.Logger.Debug(message)
}

func (m *LoggerManager) Debugf(format string, args ...any) {
	m.init()
	m.Logger.Debugf(format, args...)
}

func (m *LoggerManager) Info(message string) {
	m.init()
	m.Logger.Info(message)
}

func (m *LoggerManager) Infof(format string, args ...any) {
	m.init()
	m.Logger.Infof(format, args...)
}

func (m *LoggerManager) Warn(message string) {
	m.init()
	m.Logger.Warn(message)
}

func (m *LoggerManager) Warnf(format string, args ...any) {
	m.init()
	m.Logger.Warnf(format, args...)
}

func (m *LoggerManager) Error(message string) {
	m.init()
	m.Logger.Error(message)
}

func (m *LoggerManager) Errorf(format string, args ...any) {
	m.init()
	m.Logger.Errorf(format, args...)
}

func (m *LoggerManager) Panic(message string) {
	m.init()
	m.Logger.Panic(message)
}

func (m *LoggerManager) Panicf(format string, args ...any) {
	m.init()
	m.Logger.Panicf(format, args...)
}

func (m *LoggerManager) Fatal(message string) {
	m.init()
	m.Logger.Fatal(message)
}

func (m *LoggerManager) Fatalf(format string, args ...any) {
	m.init()
	m.Logger.Fatalf(format, args...)
}
