package driver

// Logger logger
type Logger interface {
	Debug(message string, fields map[string]any)
	Info(message string, fields map[string]any)
	Warn(message string, fields map[string]any)
	Error(message string, fields map[string]any)
	Fatal(message string, fields map[string]any)
	Panic(message string, fields map[string]any)
}
