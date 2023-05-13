package logger

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_ZapLog(t *testing.T) {
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
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          &zap.SamplingConfig{Initial: 100, Thereafter: 100},
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	l, err := zapConfig.Build(zap.AddCallerSkip(-1))
	if err != nil {
		fmt.Printf("zap build logger fail err=%v", err)
		t.Errorf(err.Error())
	}

	l.With()
	l.Sugar().Debugw("debug_log")
	l.Sugar().With()
	fmt.Println(l.Level())
	zapConfig.Level.SetLevel(zap.InfoLevel)
	l.Sugar().Debugw("debug_log2")
	fmt.Println(l.Level())
}

func TestLogger(t *testing.T) {
	// global
	Debugw("msg1", "k1", "v1") // print
	Infow("msg2", "k2", "v2")  // print
	// global set level
	SetLevel(DebugLevel)
	Debugw("msg3", "k3", "v3") // not print
	Infow("msg4", "k4", "v4")  // print
	// new instance
	l := NewLogger(DebugLevel)
	l.Debugw("msg5", "k5", "v5") // new instance print
	Debugw("msg6", "k6", "v6")   // global not print
	// with
	ll := l.With("name", "song")
	ll.Debugw("======") // 新实例包含 key "name"
	l.Debugw("=======") // 老的不包含 key "name"
	// global with
	// 层级打印正确:新实例老例例都正常打印在 此文件中:log_test
	// zap.AddCallerSkip(0)
	a := With("name", "song")
	a.Infow("======") // 新实例包含 key "name",
	Infow("======")   // 老的不包含 key "name"
}
