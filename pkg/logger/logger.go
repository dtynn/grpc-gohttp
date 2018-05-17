package logger

// Logger logger
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Level log level
type Level int

func (l Level) prefix() string {
	switch l {
	case LDebug:
		return "[D] "

	case LInfo:
		return "[I] "

	case LWarn:
		return "[W] "

	case LError:
		return "[E] "

	case lFatal:
		return "[F] "

	default:
		return ""
	}
}

// log levels
const (
	LDebug Level = iota
	LInfo
	LWarn
	LError
	lFatal
)
