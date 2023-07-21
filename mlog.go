package mlog

import (
	"runtime/debug"
	"time"

	"github.com/praissik/mlog/pkg/queue"
	"github.com/spf13/viper"
)

func Error(correlationID string, err error) {
	l := queue.LogError{
		CorrelationID: correlationID,
		StackTrace:    string(debug.Stack()),
		Level:		   "error",
		Message:       err.Error(),
		Service:	   viper.GetString("server.name"),
		Datetime:	   time.Now(),
	}
	payload := l.PrepareLog()

	queue.Publish(payload)
}

func Info(correlationID, message string, begin int64) {
	t := time.Now()
	l := queue.LogInfo{
		CorrelationID: correlationID,
		Level:         "info",
		Message:       message,
		Service:       viper.GetString("server.name"),
		Datetime:      t,
		Duration:      t.UnixMicro() - begin,
	}
	payload := l.PrepareLog()

	queue.Publish(payload)
}
