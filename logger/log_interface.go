package logger

type LoggerInterface interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	SetLevel(level Level)
	With(keyValues ...interface{}) LoggerInterface
}
