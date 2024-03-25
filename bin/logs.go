package bin

import (
	"bytes"
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

var LogsBuffer bytes.Buffer
var JsonLogger log.Logger
var TextInfoLogger log.Logger
var TextErrorLogger log.Logger

func ConfigLogs() {
	logs, err := os.OpenFile("data/logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	errors, err := os.OpenFile("data/error.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	JsonLogger.SetLevel(log.TraceLevel)
	JsonLogger.SetFormatter(&log.JSONFormatter{})
	JsonLogger.SetOutput(logs)
	TextInfoLogger.SetLevel(log.InfoLevel)
	TextInfoLogger.SetFormatter(&log.TextFormatter{})
	TextInfoLogger.SetOutput(io.MultiWriter(&LogsBuffer, os.Stdout))
	TextErrorLogger.SetLevel(log.ErrorLevel)
	TextErrorLogger.SetFormatter(&log.TextFormatter{})
	TextErrorLogger.SetOutput(errors)
}

func SaveLog(args log.Fields, level log.Level, msg string) {
	JsonLogger.WithFields(args).Log(level, msg)
	TextInfoLogger.WithFields(args).Log(level, msg)
	TextErrorLogger.WithFields(args).Log(level, msg)
	if level != log.TraceLevel && level != log.DebugLevel {
		text := LogsBuffer.String()
		Client.RPush(context.Background(), "Logs", text)
		if Client.LLen(context.Background(), "Logs").Val() > 500 {
			Client.LPop(context.Background(), "Logs")
		}
	}
	LogsBuffer = bytes.Buffer{}
}
