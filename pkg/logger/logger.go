package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	sugaredLogger *zap.SugaredLogger
}

func (logger *Logger) InitLogger(logFilePath string) {
	encoder := logger.getEncoder()
	writerSyncer := logger.getLogWriter(logFilePath)
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)
	newLogger := zap.New(core)
	logger.sugaredLogger = newLogger.Sugar()
}

func (logger *Logger) getEncoder() (encoder zapcore.Encoder) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	return
}

func (logger *Logger) getLogWriter(logFilePath string) zapcore.WriteSyncer {
	if logFilePath != "" && logFilePath != "stdout" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    1 << 7, // 128 Mb
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   false,
		}
		return zapcore.AddSync(lumberJackLogger)
	} else {
		return zapcore.AddSync(os.Stdout)
	}
}

func (logger *Logger) Sync() (err error) {
	err = logger.sugaredLogger.Sync()
	return
}

func (logger *Logger) Info(args ...interface{}) {
	logger.sugaredLogger.Info(args...)
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.sugaredLogger.Infof(template, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.sugaredLogger.Error(args...)
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.sugaredLogger.Errorf(template, args...)
}
