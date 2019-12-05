package broadcast

import (
	"github.com/sirupsen/logrus"
)

var logger Logger

func init() {
	logger = logrus.StandardLogger()
}

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Panic(i ...interface{})
	Panicf(string, ...interface{})
}

func SetLogger(log Logger) {
	logger = log
}
