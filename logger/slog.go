package logger

import (
	"log/slog"
	"os"
	"runtime"
)

var _ LoggerInterface = (*ZapSugaredLogger)(nil)

type SlogLogger struct {
	logger *slog.Logger
	level  *slog.LevelVar

	// true 代表使用slog打印文件路径,false 会使用自定的方法给日志 增加字段 file line
	addSource bool
}

// newSlog
func newSlog(level Level, addSource bool) LoggerInterface {
	levelVar := &slog.LevelVar{}
	levelVar.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: addSource, Level: levelVar}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return &SlogLogger{
		logger: logger,
		level:  levelVar,
	}
}

func (l *SlogLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l *SlogLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *SlogLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *SlogLogger) Infow(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.ApppendFileLine(keysAndValues...)
	l.logger.Info(msg, keysAndValues...)
}

func (l *SlogLogger) SetLevel(level Level) {
	zapLevelToSlogLevel(level)
	l.level.Set(slog.Level(zapLevelToSlogLevel(level)))
}

func (l *SlogLogger) With(keyValues ...interface{}) LoggerInterface {
	newLog := l.logger.With(keyValues...)
	return &SlogLogger{
		logger: newLog,
		level:  l.level,
	}
}

// ApppendFileLine 获取调用方的文件和文件号
// slog 原生 暂不支持 callerSkip,使用此函数啃根会有性能问题,最好等slog提供 CallerSkip 的参数
func (l *SlogLogger) ApppendFileLine(keyValues ...interface{}) []interface{} {
	l.addSource = false
	if !l.addSource {
		var pc uintptr
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(4, pcs[:])
		pc = pcs[0]
		fs := runtime.CallersFrames([]uintptr{pc})
		f, _ := fs.Next()
		keyValues = append(keyValues, "file", f.File, "line", f.Line)
		return keyValues

	}
	return keyValues
}
