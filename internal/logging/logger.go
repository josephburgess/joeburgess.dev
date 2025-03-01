package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

var Log *Logger

func NewLogger() *Logger {
	zapLogger, _ := zap.NewDevelopment()
	sugar := zapLogger.Sugar()
	Log = &Logger{
		SugaredLogger: sugar,
	}
	return Log
}

func Error(msg string, err error) {
	if Log != nil {
		Log.Errorw(msg, "error", err)
	}
}

func Info(msg string, args ...any) {
	if Log != nil {
		Log.Infof(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if Log != nil {
		Log.Warnf(msg, args...)
	}
}

func Debug(msg string, args ...any) {
	if Log != nil {
		Log.Debugf(msg, args...)
	}
}
