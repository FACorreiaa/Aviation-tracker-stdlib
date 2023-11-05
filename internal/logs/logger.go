package logs

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

type Formatter int

const (
	DefaultFormatter Formatter = iota
	JSONFormatter
)

type logger struct {
	logger *logrus.Logger
}

func NewLogger() Logger {
	return &logger{
		logrus.New(),
	}
}

func InitDefaultLogger() {
	DefaultLogger = NewLogger()
}

func (l *logger) ConfigureLogger(formatter Formatter) {
	l.logger.SetFormatter(getFormatter(formatter))
}

func (l *logger) NewEntry() LoggerEntry {
	return &loggerEntry{
		Logger: l,
		entry:  logrus.NewEntry(l.logger),
		caller: true,
	}
}

func (l *logger) WithFields(fields map[string]any) LoggerEntry {
	return &loggerEntry{
		Logger: l,
		entry:  l.logger.WithFields(fields),
		caller: true,
	}
}

func (l *logger) WithField(field string, value any) LoggerEntry {
	return &loggerEntry{
		Logger: l,
		entry:  l.logger.WithField(field, value),
		caller: true,
	}
}

func (l *logger) WithContext(ctx context.Context) LoggerEntry {
	return l.NewEntry().WithContext(ctx)
}

func (l *logger) WithError(err error) LoggerEntry {
	return l.NewEntry().WithError(err)
}

func (l *logger) WithoutCaller() LoggerEntry {
	return &loggerEntry{
		Logger: l,
		entry:  logrus.NewEntry(l.logger),
		caller: false,
	}
}

func (l *logger) log(level logrus.Level, args ...any) {
	l.NewEntry().log(level, args...)
}

func (l *logger) logln(level logrus.Level, args ...any) {
	l.NewEntry().logln(level, args...)
}

func (l *logger) logf(level logrus.Level, format string, args ...any) {
	l.NewEntry().logf(level, format, args...)
}

func (l *logger) Debug(msg string) {
	l.log(logrus.DebugLevel, msg)
}

func (l *logger) Info(msg string) {
	l.log(logrus.InfoLevel, msg)
}

func (l *logger) Warn(msg string) {
	l.log(logrus.WarnLevel, msg)
}

func (l *logger) Error(msg string) {
	l.log(logrus.ErrorLevel, msg)
}

func (l *logger) Fatal(msg string) {
	l.log(logrus.FatalLevel, msg)
	os.Exit(1)
}

func (l *logger) Panic(msg string) {
	l.log(logrus.PanicLevel, msg)
	os.Exit(1)
}

func (l *logger) Debugln(msg string) {
	l.logln(logrus.DebugLevel, msg)
}

func (l *logger) Infoln(msg string) {
	l.logln(logrus.InfoLevel, msg)
}

func (l *logger) Warnln(msg string) {
	l.logln(logrus.WarnLevel, msg)
}

func (l *logger) Errorln(msg string) {
	l.logln(logrus.ErrorLevel, msg)
}

func (l *logger) Fatalln(msg string) {
	l.logln(logrus.FatalLevel, msg)
}

func (l *logger) Panicln(msg string) {
	l.logln(logrus.PanicLevel, msg)
}

// logger Printf family functions

func (l *logger) Debugf(format string, args ...any) {
	l.logf(logrus.DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.logf(logrus.InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.logf(logrus.WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.logf(logrus.ErrorLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.logf(logrus.FatalLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...any) {
	l.logf(logrus.PanicLevel, format, args...)
}
