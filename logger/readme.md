### 快速体验
以下是 项目中 已经用slog替换 zap 后的 logger 使用方法,无任何感知,与之前一模一样
``` golang
package main

import "github.com/webws/go-moda/logger"

func main() {
	// 格式化打印 {"time":"2023-09-08T01:25:21.313463+08:00","level":"INFO","msg":"info hello slog","key":"value","file":"/Users/xxx/w/pro/go-moda/example/logger/main.go","line":6}
	logger.Infow("info hello slog", "key", "value")   // 打印json
	logger.Debugw("debug hello slog", "key", "value") // 不展示
	logger.SetLevel(logger.DebugLevel)                // 设置等级
	logger.Debugw("debug hello slog", "key", "value") // 设置了等级之后展示 debug
	// with
	newLog := logger.With("newkey", "newValue")
	newLog.Debugw("new hello slog") // 会打印 newkey:newValue
	logger.Debugw("old hello slog") // 不会打印 newkey:newValue
}
```
### slog 基础使用
Go 1.21版本中 将 golang.org/x/exp/slog 引入了go标准库  路径为 log/slog。 
新项目的 如果不使用第三方包，可以直接用slog当你的 logger 

slog 简单示例:
```golang
        slog.Info("finished", "key", "value")
	slog.Debug("finished", "key1", "value1")
```
以下是打印日志 默认slog 输出级别是info以上,所以debug是打印不出来.
```
2023/09/08 00:27:24 INFO finished key=value
```

json格式化,设置日志等级,并打印调用函数和文件
```golang
opts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	logger.Info("finished", "key", "value")
```
输出
```
{"time":"2023-09-08T00:34:22.035962+08:00","level":"INFO","source":{"function":"callvis/slog.TestLogJsonHandler","file":"/Users/websong/w/pro/go-note/slog/main_test.go","line":39},"msg":"finished","key":"value"}

```
### 原有 logger zap实现
原有的项目已经实现了一套logger,使用zap log 实现接口

#### 原有代码示例
logger interface LoggerInterface
``` golang
package logger

type LoggerInterface interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	SetLevel(level Level)
	With(keyValues ...interface{}) LoggerInterface
}
```
zap log 实现  LoggerInterface
``` golang
type ZapSugaredLogger struct {
	logger    *zap.SugaredLogger
	zapConfig *zap.Config
}

func buildZapLog(level Level) LoggerInterface {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	zapConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          &zap.SamplingConfig{Initial: 100, Thereafter: 100},
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	l, err := zapConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		fmt.Printf("zap build logger fail err=%v", err)
		return nil
	}
	return &ZapSugaredLogger{
		logger:    l.Sugar(),
		zapConfig: zapConfig,
	}

    func (l *ZapSugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
    }

    func (l *ZapSugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	    l.logger.Errorw(msg, keysAndValues...)
    }
    // ...省略info 之类其他实现接口的方法 
}
```
全局初始化logger,因代码量太大,以下是伪代码,主要提供思路,为下文 slog 无侵入替换zap 预热
``` golang
package logger

// 全局 log,也可以单独 NewLogger 获取新的实例
var globalog = newlogger(DebugLevel)

func newlogger(level Level) *Logger {
	l := &Logger{logger: buildZapLog(level)}
	return l
}
func Infow(msg string, keysAndValues ...interface{}) {
	globalog.logger.Infow(msg, keysAndValues...)
}
// ...省略其他全局方法,比如DebugW 之类
```
在项目里就可以通过logger 使用日志
``` golang
    logger.Debugw("msg1", "k1", "v1") // debug
	logger.SetLevel(DebugLevel)      //设置等级
	logger.Debugw("msg3", "k3", "v3") 
	newLogger := logger.With("name", "song")
	logger.Infow("msg4", "k4", "v4")  // print
```
### slog 不侵入业务 替换zap 
logger interface  接口保持不变

slog 实现 代码
```golang
package logger

import (
	"log/slog"
	"os"
	"runtime"
)

var _ LoggerInterface = (*SlogLogger)(nil)

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
func (l *SlogLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.ApppendFileLine(keysAndValues...)
	l.logger.Error(msg, keysAndValues...)
	os.Exit(1)
}

func (l *SlogLogger) Infow(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.ApppendFileLine(keysAndValues...)
	l.logger.Info(msg, keysAndValues...)
}
// 省略继承接口的其他方法 DebugW 之类的
func (l *SlogLogger) SetLevel(level Level) {
	zapLevelToSlogLevel(level)
	l.level.Set(slog.Level(zapLevelToSlogLevel(level)))
}
// 
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
```

全局初始化logger,以下伪代码
``` golang
package logger
// 全局 log,也可以单独 NewLogger 获取新的实例
var globalog = newlogger(DebugLevel)

func newlogger(level Level) *Logger {
	l := &Logger{logger: newSlog(level, false)}
	return l
}
func Infow(msg string, keysAndValues ...interface{}) {
	globalog.logger.Infow(msg, keysAndValues...)
}
// ...省略其他全局方法,比如DebugW 之类
```
### slog 实现 callerSkip 功能 

slog  的 addsource 参数 会打印文件名和行号,但 并不能像 zap 那样支持 callerSkip,也就是说  如果将 slog 封装在 logger 目录的log.go 文件下,使用logger进行打印,展示的文件会一只是log.go 

看了 slog 的源码,其实slog 使用了 runtime.Callers 在内部实现了 callerSkip 功能,但是没有对外暴露 callerSkip 参数

我就封装了 ApppendFileLine 方法,使用  runtime.Callers 获取到 文件名 和 行号,增加 file 和 line 的key value到日志

可能会有性能问题,希望slog能对外提供一个 callerSkip 参数

```golang
    var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(4, pcs[:])
	pc = pcs[0]
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	keyValues = append(keyValues, "file", f.File, "line", f.Line)
```
g