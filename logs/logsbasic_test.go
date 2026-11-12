package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestJsonlog(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//With（）函数在现有logger的基础上，通过修饰器模式，创建了一个新的logger，
	// 新logger在每次log操作的时候，都能够附加“固定属性”
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

func PrintFunctionCodeRuntimeAdrress() {
	var pcs [1]uintptr
	// Skip runtime.Callers and the current function (myFunc) to get the caller of myFunc if needed.
	// For this example, we want the PC of the log call itself, so 1 is appropriate if called directly,
	// or 0 if we were to get the PC from within the slog.NewRecord call.
	// However, if slog.NewRecord is called, it usually handles the PC internally.
	// For demonstration, let's manually create a record with a PC.
	// In a real slog usage, the handler would get the PC directly.

	// Get the PC of the line where slog.NewRecord would be called
	runtime.Callers(0, pcs[:])

	// Simulate creating a slog.Record (in real slog, this is done internally by the handler)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "This is a log message.", pcs[0])

	// Now, let's extract information from the record's PC
	pc := record.PC

	f := runtime.FuncForPC(pc)
	if f == nil {
		fmt.Println("Could not find function for PC")
		return
	}

	file, line := f.FileLine(pc)
	funcName := f.Name()

	fmt.Printf("Log generated from:\n")
	fmt.Printf("  File: %s\n", file)
	fmt.Printf("  Line: %d\n", line)
	fmt.Printf("  Function: %s\n", funcName)
}
func TestPrintFunctionCodeAddress(m *testing.T) {
	PrintFunctionCodeRuntimeAdrress()
}
