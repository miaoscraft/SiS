package log

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

type Logger struct {
	Name string
}

func NewLogger(name string) *Logger { return &Logger{name} }

func (l *Logger) Error(str string)                          { cqp.AddLog(cqp.Error, l.Name, str) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.Error(fmt.Sprintf(format, args...)) }

func (l *Logger) Waring(str string)                          { cqp.AddLog(cqp.Warning, l.Name, str) }
func (l *Logger) Waringf(format string, args ...interface{}) { l.Error(fmt.Sprintf(format, args...)) }

func (l *Logger) Info(str string)                          { cqp.AddLog(cqp.Info, l.Name, str) }
func (l *Logger) Infof(format string, args ...interface{}) { l.Error(fmt.Sprintf(format, args...)) }

func (l *Logger) Debug(str string)                          { cqp.AddLog(cqp.Debug, l.Name, str) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.Error(fmt.Sprintf(format, args...)) }
