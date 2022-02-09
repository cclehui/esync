package esyncutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

type LogData struct {
	Source  string `json:"source"`
	Level   string `json:"level"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

type Logger interface {
	Errorf(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
}

var logger Logger = &DefaultLogger{}

func SetLogger(newLogger Logger) {
	logger = newLogger
}

func GetLogger() Logger {
	return logger
}

type DefaultLogger struct{}

func (l *DefaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	l.log(ctx, &LogData{Level: LevelError, Content: content})
}

func (l *DefaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	l.log(ctx, &LogData{Level: LevelWarn, Content: content})
}

func (l *DefaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	l.log(ctx, &LogData{Level: LevelInfo, Content: content})
}

func (l *DefaultLogger) log(ctx context.Context, logData *LogData) {
	logData.Time = time.Now().Format("2006-01-02 15:04:05")
	logData.Source = "ESYNC"
	logStr, _ := json.Marshal(logData)
	fmt.Println(string(logStr))
}
