package logger

// 全局 log,也可以单独 NewLogger 获取新的实例
var globalog = newlogger(DebugLevel)

func newlogger(level Level) *Logger {
	l := &Logger{logger: buildZapLog(level)}
	return l
}

func (l *Logger) setLogger(in LoggerInterface) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logger = in
}

func Infow(msg string, keysAndValues ...interface{}) {
	globalog.logger.Infow(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	globalog.logger.Debugw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	globalog.logger.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	globalog.logger.Fatalw(msg, keysAndValues...)
}

func SetLevel(level Level) {
	globalog.logger.SetLevel(level)
}

func GetLogger() LoggerInterface {
	return globalog
}

func With(keyValues ...interface{}) LoggerInterface {
	newLog := globalog.logger.With(keyValues...)
	l := &Logger{logger: newLog}
	return l
}
