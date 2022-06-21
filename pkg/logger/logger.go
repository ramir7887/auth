package logger

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger -.
type Logger struct {
	entry *logrus.Entry
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	var (
		l logrus.Level
		f logrus.Formatter
	)

	f = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}

	switch strings.ToLower(level) {
	case "fatal":
		l = logrus.FatalLevel
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		log.Fatalf("Unknown log level: %s", level)
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetFormatter(f)
	logger.SetOutput(os.Stdout)

	entry := logger.WithFields(logrus.Fields{
		"pid":     os.Getpid(),
		"service": "auth",
	})

	return &Logger{
		entry: entry,
	}
}

func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

// Debug -.
func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf -.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info -.
func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof -.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn -.
func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf -.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error -.
func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf -.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal -.
func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalf -.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}
