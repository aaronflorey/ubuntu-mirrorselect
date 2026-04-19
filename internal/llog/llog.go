package llog

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var (
	levelVar = &slog.LevelVar{}
	logger   = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: levelVar}))
)

func init() {
	levelVar.Set(slog.LevelWarn)
}

func SetLogLevel(level string) error {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		levelVar.Set(slog.LevelDebug)
	case "INFO":
		levelVar.Set(slog.LevelInfo)
	case "WARN":
		levelVar.Set(slog.LevelWarn)
	case "ERROR":
		levelVar.Set(slog.LevelError)
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	return nil
}

func Debug(v ...interface{}) {
	logger.Debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Debugln(v ...interface{}) {
	logger.Debug(strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}

func Info(v ...interface{}) {
	logger.Info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Infoln(v ...interface{}) {
	logger.Info(strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}

func Warn(v ...interface{}) {
	logger.Warn(fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Warnln(v ...interface{}) {
	logger.Warn(strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}

func Error(v ...interface{}) {
	logger.Error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Errorln(v ...interface{}) {
	logger.Error(strings.TrimSuffix(fmt.Sprintln(v...), "\n"))
}
