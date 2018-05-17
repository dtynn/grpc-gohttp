package logger

import "log"

var nl Logger = (*nop)(nil)

// Nop return nop logger
func Nop() Logger {
	return nl
}

type nop struct {
}

func (*nop) Debug(v ...interface{}) {}

func (*nop) Debugf(format string, v ...interface{}) {}

func (*nop) Info(v ...interface{}) {}

func (*nop) Infof(format string, v ...interface{}) {}

func (*nop) Warn(v ...interface{}) {}

func (*nop) Warnf(format string, v ...interface{}) {}

func (*nop) Error(v ...interface{}) {}

func (*nop) Errorf(format string, v ...interface{}) {}

func (*nop) Fatal(v ...interface{}) {
	log.Fatalln(v...)
}

func (*nop) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
