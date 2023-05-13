package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ LoggerInterface = (*ZapSugaredLogger)(nil)

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
}

/*
	func (l *ZapSugaredLogger) Debug(args ...interface{}) {
		l.logger.Debug(args...)
	}

	func (l *ZapSugaredLogger) Info(args ...interface{}) {
		l.logger.Info(args...)
	}

	func (l *ZapSugaredLogger) Warn(args ...interface{}) {
		l.logger.Warn(args...)
	}

	func (l *ZapSugaredLogger) Error(args ...interface{}) {
		l.logger.Error(args...)
	}

	func (l *ZapSugaredLogger) DPanic(args ...interface{}) {
		l.logger.DPanic(args...)
	}

	func (l *ZapSugaredLogger) Panic(args ...interface{}) {
		l.logger.Panic(args...)
	}

	func (l *ZapSugaredLogger) Fatal(args ...interface{}) {
		l.logger.Fatal(args...)
	}

	func (l *ZapSugaredLogger) Debugf(template string, args ...interface{}) {
		l.logger.Debugf(template, args...)
	}

	func (l *ZapSugaredLogger) Infof(template string, args ...interface{}) {
		l.logger.Infof(template, args...)
	}

	func (l *ZapSugaredLogger) Warnf(template string, args ...interface{}) {
		l.logger.Warnf(template, args...)
	}

	func (l *ZapSugaredLogger) Errorf(template string, args ...interface{}) {
		l.logger.Errorf(template, args...)
	}

	func (l *ZapSugaredLogger) Fatalf(template string, args ...interface{}) {
		l.logger.Fatalf(template, args...)
	}
*/
func (l *ZapSugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *ZapSugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *ZapSugaredLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *ZapSugaredLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *ZapSugaredLogger) SetLevel(level Level) {
	l.zapConfig.Level.SetLevel(zapcore.Level(level))
}

func (l *ZapSugaredLogger) With(keyValues ...interface{}) LoggerInterface {
	ll := l.logger.With(keyValues...).WithOptions(zap.AddCallerSkip(0))
	return &ZapSugaredLogger{
		zapConfig: l.zapConfig,
		logger:    ll,
	}
}
