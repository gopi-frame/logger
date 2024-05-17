package logger

import (
	"fmt"
	"strings"
	"testing"

	eventcontract "github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/event"
	events "github.com/gopi-frame/logger/event"
	"github.com/stretchr/testify/assert"
)

type _formatter struct{}

func (f *_formatter) Format(entry logger.Entry) string {
	str := ""
	for k, v := range entry.Fields() {
		str += k + "=" + fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("message=%s %v", entry.Message(), str)
}

type _infoHook struct {
	message string
}

func (h *_infoHook) Enable(level string) bool {
	return level == "info"
}
func (h *_infoHook) Handle(e logger.Entry) error {
	h.message = e.Message()
	return nil
}

func TestManager(t *testing.T) {
	manager := NewManager()
	assert.Equal(t, int64(1), manager.channels.Count())
	assert.True(t, manager.channels.ContainsKey("default"))

	eventMessage := ""
	hook1 := new(_infoHook)
	dispatcher := event.NewDispatcher()
	dispatcher.Listen([]eventcontract.Event{new(events.MessageLogged)}, event.Listener(func(event eventcontract.Event) bool {
		e := event.(*events.MessageLogged)
		eventMessage = fmt.Sprintf("level=%s,message=%s,fields=%v", e.Level, e.Message, e.Fields)
		return true
	}))
	message := new(strings.Builder)
	testLogger := NewLogger(
		Outputs(message),
		Formatter(new(_formatter)),
		Hooks(hook1),
		Dispatcher(dispatcher),
	)
	manager.AddChannel("test", testLogger)
	assert.Equal(t, int64(2), manager.channels.Count())
	assert.True(t, manager.channels.ContainsKey("test"))

	logger := manager.Channel("test")
	logger.Info("info", map[string]any{"key": "value"})
	assert.Equal(t, "message=info context=map[key:value]", message.String())
	assert.Equal(t, "info", hook1.message)
	assert.Equal(t, "level=info,message=info,fields=map[key:value]", eventMessage)
	message.Reset()

	logger.Warn("warn", nil)
	assert.Equal(t, "message=warn context=map[]", message.String())
	assert.Equal(t, "info", hook1.message)
	assert.Equal(t, "level=warn,message=warn,fields=map[]", eventMessage)

}
