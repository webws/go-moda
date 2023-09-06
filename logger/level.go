package logger

import "log/slog"

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

func zapLevelToSlogLevel(level Level) slog.Level {
	levelMap := map[Level]slog.Level{
		DebugLevel: slog.LevelDebug,
		InfoLevel:  slog.LevelInfo,
		WarnLevel:  slog.LevelWarn,
		ErrorLevel: slog.LevelError,
	}
	if slogLevel, ok := levelMap[level]; ok {
		return slogLevel
	}
	return slog.LevelError
}
