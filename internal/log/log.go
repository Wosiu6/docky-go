package log

import "log"

type Logger interface {
	Info(msg string, kv ...any)
	Error(msg string, kv ...any)
}

type stdLogger struct{}

func New() Logger { return &stdLogger{} }

func (l *stdLogger) Info(msg string, kv ...any)  { log.Println(append([]any{"INFO", msg}, kv...)...) }
func (l *stdLogger) Error(msg string, kv ...any) { log.Println(append([]any{"ERROR", msg}, kv...)...) }

