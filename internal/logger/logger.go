package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	core   *zap.Logger
	level  string
	writer io.Writer
}

var custLogger *Logger

const (
	ErrorLevel = "ERROR"
	WarnLevel  = "WARN"
	InfoLevel  = "INFO"
	DebugLevel = "DEBUG"
)

func init() {
	custLogger = &Logger{core: zap.Must(zap.NewDevelopment()), level: InfoLevel, writer: os.Stdout}
	initCore()
}

func SetLogLevel(level string) {
	custLogger.level = level
	initCore()
}

func SetWriter(writer io.Writer) {
	custLogger.writer = writer
	initCore()
}

func initCore() {
	var zapLevel zapcore.Level

	switch custLogger.level {
	case "WARN":
		zapLevel = zap.WarnLevel
	case "INFO":
		zapLevel = zap.InfoLevel
	case "DEBUG":
		zapLevel = zap.DebugLevel
	default:
		zapLevel = zap.ErrorLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(custLogger.writer),
		zap.NewAtomicLevelAt(zapLevel))

	custLogger.core = zap.New(core)
}

func GetLogger() *Logger {
	return custLogger
}

func Error(msg string) {
	custLogger.core.Error(msg)
}

func Warn(msg string) {
	custLogger.core.Warn(msg)
}

func Info(msg string) {
	custLogger.core.Info(msg)
}

func Debug(msg string) {
	custLogger.core.Debug(msg)
}

func Fatal(msg string) {
	custLogger.core.Fatal(msg)
}
