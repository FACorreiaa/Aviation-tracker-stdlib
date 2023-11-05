package logs

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LoggerEntry interface {
	WithFields(fields map[string]any) LoggerEntry
	WithField(field string, value any) LoggerEntry
	WithContext(ctx context.Context) LoggerEntry
	WithError(err error) LoggerEntry
	WithoutCaller() LoggerEntry

	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Debugln(string)
	Infoln(string)
	Warnln(string)
	Errorln(string)
	Fatalln(string)

	log(level logrus.Level, args ...any)
	logf(level logrus.Level, format string, args ...any)
	logln(level logrus.Level, args ...any)
}

type Logger interface {
	LoggerEntry
	NewEntry() LoggerEntry
	ConfigureLogger(formatter Formatter)
}
