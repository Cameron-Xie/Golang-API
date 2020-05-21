package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewJsonLogger(t *testing.T) {
	a := assert.New(t)
	l := NewJSONLogger().(*logrus.Logger)
	f, ok := l.Formatter.(*logrus.JSONFormatter)

	a.True(ok)
	a.True(f.DisableTimestamp)
}
