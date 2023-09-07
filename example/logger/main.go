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
