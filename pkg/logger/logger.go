package logger

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

const (
	// We use +2 to avoid conflict with zapcore levels. So debug start with 1 and fatal end with 7

	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel = zapcore.DebugLevel + 2
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel + 2
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel = zapcore.WarnLevel + 2
	// ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel + 2
	// DPanicLevel logs are particularly important errors. In development the logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel + 2
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel + 2
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel + 2
)

func InitLogger(level int) *zap.Logger {
	var loggerLevel zapcore.Level
	switch level {
	case int(DebugLevel):
		loggerLevel = zapcore.DebugLevel
	case int(InfoLevel):
		loggerLevel = zapcore.InfoLevel
	case int(WarnLevel):
		loggerLevel = zapcore.WarnLevel
	case int(ErrorLevel):
		loggerLevel = zapcore.ErrorLevel
	case int(DPanicLevel):
		loggerLevel = zapcore.DPanicLevel
	case int(PanicLevel):
		loggerLevel = zapcore.PanicLevel
	case int(FatalLevel):
		loggerLevel = zapcore.FatalLevel
	default:
		loggerLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(loggerLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)

	return logger
}

func ForceLog(template string, args ...any) {
	template = "%v    INFO    %s " + template + "\n"
	file := getFileAndLine(2)
	args = append([]any{time.Now().Format("2006-01-02T15:04:05.000Z0700"), file}, args...)
	fmt.Printf(template, args...)
}

func getFileAndLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)

	if !ok {
		return ""
	}

	file = parseFileName(file)

	return file + ":" + strconv.Itoa(line)
}

func parseFileName(file string) string {
	if file == "" {
		return ""
	}

	file = file[len(rootDir()):]

	if file[0] == '/' {
		file = file[1:]
	}

	firstFolderRegex := regexp.MustCompile(`^[a-zA-Z0-9\-]+/`)
	file = firstFolderRegex.ReplaceAllString(file, "$1")

	return file
}

func rootDir() string {
	rootDir, _ := os.Getwd()

	return rootDir
}
