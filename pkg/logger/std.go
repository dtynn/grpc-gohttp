package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Std loggers
var (
	Stdout = New(LDebug, os.Stdout)
	Stderr = New(LDebug, os.Stderr)
)

// New return a standard logger
func New(level Level, w io.Writer) Logger {
	return &std{
		level:  level,
		Logger: log.New(w, "", log.LstdFlags|log.Lshortfile),
	}
}

type std struct {
	*log.Logger
	level Level
}

func (s *std) Debug(v ...interface{}) {
	s.print(LDebug, fmt.Sprintln(v...))
}

func (s *std) Debugf(format string, v ...interface{}) {
	s.print(LDebug, fmt.Sprintf(format, v...))
}

func (s *std) Info(v ...interface{}) {
	s.print(LInfo, fmt.Sprintln(v...))
}

func (s *std) Infof(format string, v ...interface{}) {
	s.print(LInfo, fmt.Sprintf(format, v...))
}

func (s *std) Warn(v ...interface{}) {
	s.print(LWarn, fmt.Sprintln(v...))
}

func (s *std) Warnf(format string, v ...interface{}) {
	s.print(LWarn, fmt.Sprintf(format, v...))
}

func (s *std) Error(v ...interface{}) {
	s.print(LError, fmt.Sprintln(v...))
}

func (s *std) Errorf(format string, v ...interface{}) {
	s.print(LError, fmt.Sprintf(format, v...))
}

func (s *std) Fatal(v ...interface{}) {
	s.print(lFatal, fmt.Sprintln(v...))
	os.Exit(2)
}

func (s *std) Fatalf(format string, v ...interface{}) {
	s.print(lFatal, fmt.Sprintf(format, v...))
	os.Exit(2)
}

func (s *std) print(l Level, msg string) {
	if l < s.level {
		return
	}

	s.Logger.Print(l.prefix() + msg)
}
