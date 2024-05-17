package event

// NewMessageLogged new message logged event
func NewMessageLogged(level string, message string, fields map[string]any) *MessageLogged {
	return &MessageLogged{
		Level:   level,
		Message: message,
		Fields:  fields,
	}
}

// MessageLogged message logged event
type MessageLogged struct {
	Level   string
	Message string
	Fields  map[string]any
}

// Topic topic
func (e *MessageLogged) Topic() string {
	return MessageLoggedTopic
}
