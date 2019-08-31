package data

var Logger interface {
	Error(str string)
	Errorf(format string, args ...interface{})

	Waring(str string)
	Waringf(format string, args ...interface{})

	Info(str string)
	Infof(format string, args ...interface{})

	Debug(str string)
	Debugf(format string, args ...interface{})
}
