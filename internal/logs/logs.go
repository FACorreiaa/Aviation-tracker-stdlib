package logs

import (
	"github.com/sirupsen/logrus"
)

var DefaultLogger = NewLogger()

func getFormatter(formatter Formatter) logrus.Formatter {
	switch formatter {
	case JSONFormatter:
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}
