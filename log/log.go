package log

import (
	"github.com/sirupsen/logrus"
)

type Logger logrus.Entry

//type Logger struct {
//	Name string
//}

func NewLogger(name string) *Logger {
	return (*Logger)(logrus.New().WithField("name", name))
}

func (l *Logger) Error(str string) { (*logrus.Entry)(l).Error(str) }
func (l *Logger) Errorf(format string, args ...interface{}) {
	(*logrus.Entry)(l).Errorf(format, args...)
}

func (l *Logger) Warning(str string) { (*logrus.Entry)(l).Warning(str) }
func (l *Logger) Warningf(format string, args ...interface{}) {
	(*logrus.Entry)(l).Warningf(format, args...)
}

func (l *Logger) Info(str string)                          { (*logrus.Entry)(l).Info(str) }
func (l *Logger) Infof(format string, args ...interface{}) { (*logrus.Entry)(l).Infof(format, args...) }

func (l *Logger) Debug(str string) { (*logrus.Entry)(l).Debug(str) }
func (l *Logger) Debugf(format string, args ...interface{}) {
	(*logrus.Entry)(l).Debugf(format, args...)
}
