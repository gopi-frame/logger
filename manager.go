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
	lazyChannels   *kv.Map[string, func() (logger.Logger, error)]
}

// NewManager creates a new logger manager.
func NewManager() *LoggerManager {
	return &LoggerManager{
		channels:     kv.NewMap[string, logger.Logger](),
		lazyChannels: kv.NewMap[string, func() (logger.Logger, error)](),
	}
}

// SetDefaultChannel sets the default channel name.
func (m *LoggerManager) SetDefaultChannel(name string) {
	m.defaultChannel = name
}

// Use sets the default channel instance.
func (m *LoggerManager) Use(channel logger.Logger) *LoggerManager {
	m.Logger = channel
	return m
}

// AddChannel adds a channel to the manager.
func (m *LoggerManager) AddChannel(name string, logger logger.Logger) {
	m.channels.Lock()
	defer m.channels.Unlock()
	m.channels.Set(name, logger)
}

// AddLazyChannel adds a lazy channel to the manager.
func (m *LoggerManager) AddLazyChannel(name string, config map[string]any) {
	m.lazyChannels.Lock()
	defer m.lazyChannels.Unlock()
	m.lazyChannels.Set(name, func() (logger.Logger, error) {
		driver := config["driver"].(string)
		return Open(driver, config)
	})
}

// HasChannel returns true if the channel with the given name exists.
func (m *LoggerManager) HasChannel(name string) bool {
	m.channels.RLock()
	if m.channels.ContainsKey(name) {
		m.channels.RUnlock()
		return true
	}
	m.channels.RUnlock()
	m.lazyChannels.RLock()
	if m.lazyChannels.ContainsKey(name) {
		m.lazyChannels.RUnlock()
		return true
	}
	m.lazyChannels.RUnlock()
	return false
}

// TryChannel returns the channel with the given name.
// If the channel is not found or error occurred, it returns an error.
func (m *LoggerManager) TryChannel(name string) (logger.Logger, error) {
	m.channels.RLock()
	if channel, ok := m.channels.Get(name); ok {
		m.channels.RUnlock()
		return channel, nil
	}
	m.channels.RUnlock()
	m.lazyChannels.RLock()
	if lazyChannel, ok := m.lazyChannels.Get(name); ok {
		m.lazyChannels.RUnlock()
		channel, err := lazyChannel()
		if err != nil {
			return nil, err
		}
		m.channels.Lock()
		defer m.channels.Unlock()
		m.channels.Set(name, channel)
		return channel, nil
	}
	m.lazyChannels.RUnlock()
	return nil, NewNotConfiguredChannelException(name)
}

// Channel returns the channel with the given name.
// If the channel is not found or error occurred, it panics.
func (m *LoggerManager) Channel(name string) logger.Logger {
	channel, err := m.TryChannel(name)
	if err != nil {
		panic(err)
	}
	return channel
}

// ChannelOrDefault returns the channel with the given name.
// If the channel is not found or error occurred, it returns the default channel.
func (m *LoggerManager) ChannelOrDefault(name string) logger.Logger {
	channel, err := m.TryChannel(name)
	if err != nil {
		channel = m.Channel(m.defaultChannel)
	}
	return channel
}
