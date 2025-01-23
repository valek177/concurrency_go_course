package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

var (
	globalLogger *zap.Logger

	defaultLoggerFilename = "log/output.log"
	loggerMaxSizeMb       = 10
	loggerMaxBackupsCount = 3
	loggerMaxAgeDays      = 7
)

// MockLogger mocks logger
func MockLogger() {
	globalLogger = zap.NewNop()
}

// InitLogger initializes logger with level
func InitLogger(logLevel, filename string) {
	Init(getCore(getAtomicLevel(logLevel), filename))
}

// Init initializes new logger
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// Debug is used for debug logging
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info is used for info logging
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn is used for warn logging
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error is used for error logging
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// ErrorWithMsg is used for error logging with error param
func ErrorWithMsg(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	globalLogger.Error(msg, fields...)
}

// Fatal is used for fatal logging
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// FatalWithMsg is used for fatal logging with error param
func FatalWithMsg(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	globalLogger.Fatal(msg, fields...)
}

// WithOptions applies options
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}

func getAtomicLevel(logLevel string) zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		FatalWithMsg("failed to set log level: ", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func getCore(level zap.AtomicLevel, filename string) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	loggerFilename := defaultLoggerFilename
	if len(filename) != 0 {
		loggerFilename = filename
	}

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   loggerFilename,
		MaxSize:    loggerMaxSizeMb,
		MaxBackups: loggerMaxBackupsCount,
		MaxAge:     loggerMaxAgeDays,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}
