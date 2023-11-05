package logs

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	ErrorKey   = "error"
	TraceIdKey = "trace_id"

	TraceIdCtxKey = "TraceId"
)

type loggerEntry struct {
	Logger *logger
	entry  *logrus.Entry
	caller bool
}

func (l *loggerEntry) withField(field string, value any) *loggerEntry {
	return &loggerEntry{
		Logger: l.Logger,
		entry:  l.entry.WithField(field, value),
		caller: l.caller,
	}
}

func (l *loggerEntry) withFields(fields map[string]any) *loggerEntry {
	return &loggerEntry{
		Logger: l.Logger,
		entry:  l.entry.WithFields(fields),
		caller: l.caller,
	}
}

func (l *loggerEntry) WithFields(fields map[string]any) LoggerEntry {
	return l.withFields(fields)
}

func (l *loggerEntry) WithField(field string, value interface{}) LoggerEntry {
	return l.withField(field, value)
}

func (l *loggerEntry) WithContext(ctx context.Context) LoggerEntry {
	return &loggerEntry{
		Logger: l.Logger,
		entry: l.entry.WithFields(map[string]any{
			TraceIdKey: ctx.Value(TraceIdCtxKey),
		}),
		caller: l.caller,
	}
}

func (l *loggerEntry) WithError(err error) LoggerEntry {
	return &loggerEntry{
		Logger: l.Logger,
		entry:  l.entry.WithField(ErrorKey, err),
		caller: l.caller,
	}
}

func (l *loggerEntry) WithoutCaller() LoggerEntry {
	return &loggerEntry{
		Logger: l.Logger,
		entry:  l.entry,
		caller: false,
	}
}

func (l *loggerEntry) withCaller() *loggerEntry {
	info := newFileInfo(3)
	return l.withFields(map[string]any{
		"file":     info.file,
		"line":     info.line,
		"funcName": info.funcName,
	})
}
func (l *loggerEntry) log(level logrus.Level, args ...any) {
	le := l
	if l.caller {
		le = le.withCaller()
	}
	le.entry.Log(level, args...)
}

func (l *loggerEntry) logln(level logrus.Level, args ...any) {
	le := l
	if l.caller {
		le = le.withCaller()
	}
	le.entry.Logln(level, args...)
}

func (l *loggerEntry) logf(level logrus.Level, format string, args ...interface{}) {
	le := l
	if l.caller {
		le = le.withCaller()
	}
	le.entry.Logf(level, format, args...)
}

func (l *loggerEntry) Debug(msg string) {
	l.log(logrus.DebugLevel, msg)
}

func (l *loggerEntry) Info(msg string) {
	l.log(logrus.InfoLevel, msg)
}

func (l *loggerEntry) Warn(msg string) {
	l.log(logrus.WarnLevel, msg)
}

func (l *loggerEntry) Error(msg string) {
	l.log(logrus.ErrorLevel, msg)
}

func (l *loggerEntry) Fatal(msg string) {
	l.log(logrus.FatalLevel, msg)
}

func (l *loggerEntry) Panic(msg string) {
	l.log(logrus.PanicLevel, msg)
}

func (l *loggerEntry) Debugln(msg string) {
	l.logln(logrus.DebugLevel, msg)
}

func (l *loggerEntry) Infoln(msg string) {
	l.logln(logrus.InfoLevel, msg)
}

func (l *loggerEntry) Warnln(msg string) {
	l.logln(logrus.WarnLevel, msg)
}

func (l *loggerEntry) Errorln(msg string) {
	l.logln(logrus.ErrorLevel, msg)
}

func (l *loggerEntry) Fatalln(msg string) {
	l.logln(logrus.FatalLevel, msg)
	os.Exit(1)
}

func (l *loggerEntry) Panicln(msg string) {
	l.logln(logrus.PanicLevel, msg)
	os.Exit(1)
}

func (l *loggerEntry) Debugf(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}

func (l *loggerEntry) Infof(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}

func (l *loggerEntry) Warnf(format string, args ...interface{}) {
	l.logf(logrus.WarnLevel, format, args...)
}

func (l *loggerEntry) Errorf(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}

func (l *loggerEntry) Fatalf(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}

func (l *loggerEntry) Panicf(format string, args ...interface{}) {
	l.logf(logrus.PanicLevel, format, args...)
}
