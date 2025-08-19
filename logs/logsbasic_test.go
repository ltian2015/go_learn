package main

import (
	"log/slog"
	"os"
	"testing"
)

func TestJsonlog(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 为 logger 添加固定属性
	loggerWithReqID := logger.With("request_id", "req-123abc")
	loggerWithReqID.Info("Request started") // 所有通过此 logger 记录的日志都会包含 request_id

	// 使用 Group
	logger.Info("database query",
		slog.Group("query_info",
			slog.String("table", "users"),
			slog.Int("duration_ms", 12),
		),
		slog.Group("connection",
			slog.String("host", "db01.example.com"),
			slog.Int("port", 5432),
		),
	)
	// 输出 JSON 中会包含: "query_info": {"table": "users", "duration_ms": 12}, "connection": {"host": "db01.example.com", "port": 5432}
}
