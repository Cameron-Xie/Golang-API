package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}

func NewJSONLogger() Logger {
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	return l
}
