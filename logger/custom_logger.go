package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type CustomLogger struct {
	SugarLogger *zap.SugaredLogger
}

func (l *CustomLogger) Info(args ...interface{}) {
	l.SugarLogger.Info(args...)
}

func (l *CustomLogger) Infof(template string, args ...interface{}) {
	l.SugarLogger.Infof(template, args...)
}

func (l *CustomLogger) Error(args ...interface{}) {
	l.SugarLogger.Error(args...)
}

func (l *CustomLogger) Errorf(template string, args ...interface{}) {
	l.SugarLogger.Errorf(template, args...)
}

func (l *CustomLogger) Fatal(args ...interface{}) {
	l.SugarLogger.Fatal(args...)
}

func (l *CustomLogger) Fatalf(template string, args ...interface{}) {
	l.SugarLogger.Fatalf(template, args...)
}
