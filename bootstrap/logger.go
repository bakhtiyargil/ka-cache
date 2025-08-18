package bootstrap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ka-cache/config"
	"ka-cache/logger"
	"os"
)

var loggerLevelMap = map[string]zapcore.Level{
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

func initLogger() logger.Logger {
	customLogger := &logger.CustomLogger{}
	logLevel := getLoggerLevel(App.Config)
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
	lggr := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	customLogger.SugarLogger = lggr.Sugar()
	if err := customLogger.SugarLogger.Sync(); err != nil {
		customLogger.SugarLogger.Error(err)
	}
	return customLogger
}

func getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}
