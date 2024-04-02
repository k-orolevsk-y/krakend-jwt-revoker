package log

type ILogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

type NoopLogger struct{}

func (n NoopLogger) Debug(_ ...interface{})    {}
func (n NoopLogger) Info(_ ...interface{})     {}
func (n NoopLogger) Warning(_ ...interface{})  {}
func (n NoopLogger) Error(_ ...interface{})    {}
func (n NoopLogger) Critical(_ ...interface{}) {}
func (n NoopLogger) Fatal(_ ...interface{})    {}
