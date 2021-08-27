package log

import (
	"fmt"
	"log"
)

type Logger struct {
	Name string
}

func NewLogger(name string) *Logger { return &Logger{name} }

func (l *Logger) Error(str string)                          { log.Println(l.Name, "|", str) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.Error(fmt.Sprintf(format, args...)) }

func (l *Logger) Waring(str string)                          { log.Println(l.Name, "|", str) }
func (l *Logger) Waringf(format string, args ...interface{}) { l.Waring(fmt.Sprintf(format, args...)) }

func (l *Logger) Info(str string)                          { log.Println(l.Name, "|", str) }
func (l *Logger) Infof(format string, args ...interface{}) { l.Info(fmt.Sprintf(format, args...)) }

func (l *Logger) Debug(str string)                          { log.Println(l.Name, "|", str) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.Debug(fmt.Sprintf(format, args...)) }
