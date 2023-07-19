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

func Info(correlationID string) {
	l := queue.LogInfo{
		CorrelationID: correlationID,
		Level:         "info",
		Service:       viper.GetString("server.name"),
		Datetime:      time.Now(),
	}
	payload := l.PrepareLog()

	queue.Publish(payload)
}
