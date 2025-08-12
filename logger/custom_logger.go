package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ka-cache/config"
	"os"
)

type Logger interface {
	InitLogger()
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

type customLogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

func NewCustomLogger(cfg *config.Config) Logger {
	return &customLogger{cfg: cfg}
}

var loggerLevelMap = map[string]zapcore.Level{
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

func (l *customLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *customLogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)
	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	encoderCfg = zap.NewDevelopmentEncoderConfig()

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	encoder = zapcore.NewConsoleEncoder(encoderCfg)

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

func (l *customLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *customLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *customLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *customLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}
