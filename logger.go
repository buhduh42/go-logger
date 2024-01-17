package logger

/* a message to view a commit */

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel uint8

//go:generate stringer -type=LogLevel
const (
	SILENT LogLevel = 1 << iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

type Logger interface {
	//primary methods
	Log(LogLevel, string) error
	Logf(LogLevel, string, ...interface{}) error

	//satisfies the writer interface
	Write([]byte) (int, error)

	//helpers
	Fatal(string) error
	Fatalf(string, ...interface{}) error
	Error(string) error
	Errorf(string, ...interface{}) error
	Warn(string) error
	Warnf(string, ...interface{}) error
	Info(string) error
	Infof(string, ...interface{}) error
	Debug(string) error
	Debugf(string, ...interface{}) error
	Trace(string) error
	Tracef(string, ...interface{}) error
}

type FormattedCallback func(string, ...interface{}) error

type logger struct {
	writer io.Writer
	level  LogLevel
	name   string
	_      struct{}
}

func NopLogger() Logger {
	return NewLogger(SILENT, "nop logger", io.Discard)
}

func NewLogLevel(level string) (LogLevel, error) {
	switch level {
	case SILENT.String():
		return SILENT, nil
	case FATAL.String():
		return FATAL, nil
	case ERROR.String():
		return ERROR, nil
	case WARN.String():
		return WARN, nil
	case INFO.String():
		return INFO, nil
	case DEBUG.String():
		return DEBUG, nil
	case TRACE.String():
		return TRACE, nil
	}
	return WARN, fmt.Errorf("unrecognized level string: '%s'", level)
}

func NewLogger(level LogLevel, name string, locations ...io.Writer) Logger {
	var writers io.Writer
	if level == SILENT {
		writers = io.Discard
	} else if len(locations) == 0 {
		writers = os.Stdout
	} else {
		writers = io.MultiWriter(locations...)
	}
	return &logger{
		writer: writers,
		level:  level,
		name:   name,
	}
}

func (l *logger) Fatal(msg string) error {
	return l.Log(FATAL, msg)
}

func (l *logger) Fatalf(format string, data ...interface{}) error {
	return l.Logf(FATAL, format, data...)
}

func (l *logger) Error(msg string) error {
	return l.Log(ERROR, msg)
}

func (l *logger) Errorf(format string, data ...interface{}) error {
	return l.Logf(ERROR, format, data...)
}

func (l *logger) Warn(msg string) error {
	return l.Log(WARN, msg)
}

func (l *logger) Warnf(format string, data ...interface{}) error {
	return l.Logf(WARN, format, data...)
}

func (l *logger) Info(msg string) error {
	return l.Log(INFO, msg)
}

func (l *logger) Infof(format string, data ...interface{}) error {
	return l.Logf(INFO, format, data...)
}

func (l *logger) Debug(msg string) error {
	return l.Log(DEBUG, msg)
}

func (l *logger) Debugf(format string, data ...interface{}) error {
	return l.Logf(DEBUG, format, data...)
}

func (l *logger) Trace(msg string) error {
	return l.Log(TRACE, msg)
}

func (l *logger) Tracef(format string, data ...interface{}) error {
	return l.Logf(TRACE, format, data...)
}

func (l *logger) Log(level LogLevel, msg string) error {
	if level == SILENT || l.level < level {
		return nil
	}
	_, err := l.Write(l.getMsg(level, msg))
	return err
}

func (l *logger) Logf(level LogLevel, format string, data ...interface{}) error {
	return l.Log(level, fmt.Sprintf(format, data...))
}

func (l *logger) Write(p []byte) (int, error) {
	n, err := l.writer.Write(p)
	if n < len(p) {
		//this could cause an infinite recursion, do it sepearately
		l.warnWrite(n, len(p))
	}
	return n, err
}

func (l *logger) warnWrite(numBytes, total int) {
	//println("this is the warn logger")
	if l.level >= WARN {
		l.writer.Write(
			l.getMsg(
				WARN,
				fmt.Sprintf("only %d bytes written of a %d bytes length message", numBytes, total),
			),
		)
	}
}

func (l *logger) getMsg(level LogLevel, msg string) []byte {
	return []byte(
		fmt.Sprintf("logger[%s][%s] - %s - %s\n", level, l.name, time.Now().Format(time.RFC1123Z), msg),
	)
}
