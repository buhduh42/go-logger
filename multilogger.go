package logger

import (
	"fmt"
	"strings"
)

type multiLogger []Logger

func MultiLogger(loggers ...Logger) Logger {
	toRet := make(multiLogger, len(loggers))
	for i, l := range loggers {
		toRet[i] = l
	}
	return toRet
}

func (m multiLogger) Log(level LogLevel, msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		if err := l.Log(level, msg); err != nil {
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("multilogger Log() failed: index %d: error: '%s'", i, err),
			)
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Logf(level LogLevel, msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		if err := l.Logf(level, msg, data...); err != nil {
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("multilogger Logf() failed: index %d: error: '%s'", i, err),
			)
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Write(p []byte) (int, error) {
	numWrote := make([]int, len(m))
	errMsgs := make([]string, 0)
	for i, l := range m {
		n, err := l.Write(p)
		numWrote = append(numWrote, n)
		if err != nil {
			errMsgs = append(
				errMsgs,
				fmt.Sprintf("multilogger Write() failed: index %d: error: '%s'", i, err),
			)
		}
	}
	min := numWrote[0]
	for _, n := range numWrote {
		if n < min {
			min = n
		}
	}
	if len(errMsgs) > 0 {
		return min, fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
	}
	return min, nil
}

func multilogMsgHelper(
	cb func(string) error,
	cbName, msg string,
	errMsgs []string,
	i int, pLogger Logger,
) {
	err := cb(msg)
	if err != nil {
		var msg string
		if l, ok := pLogger.(*logger); ok {
			msg = fmt.Sprintf(
				"multilogger %s() failed, logger name: '%s' error: '%s'",
				cbName, l.name, err,
			)
		} else {
			msg = fmt.Sprintf(
				"multilogger %s() failed, index %d: error: '%s'",
				cbName, i, err,
			)
		}
		errMsgs = append(errMsgs, msg)
	}
}

func multilogfMsgHelper(
	cb func(string, ...interface{}) error,
	cbName, msg string,
	errMsgs []string,
	i int, pLogger Logger,
	data ...interface{},
) {
	err := cb(msg, data...)
	if err != nil {
		var msg string
		if l, ok := pLogger.(*logger); ok {
			msg = fmt.Sprintf(
				"multilogger %s() failed, logger name: '%s' error: '%s'",
				cbName, l.name, err,
			)
		} else {
			msg = fmt.Sprintf(
				"multilogger %s() failed, index %d: error: '%s'",
				cbName, i, err,
			)
		}
		errMsgs = append(errMsgs, msg)
	}
}

func (m multiLogger) Fatal(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Fatal, "Fatal", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Fatalf(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Fatalf, "Fatalf", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Error(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Error, "Error", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Errorf(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Errorf, "Errorf", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Warn(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Warn, "Warn", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Warnf(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Warnf, "Warnf", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Info(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Info, "Info", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Infof(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Infof, "Infof", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Debug(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Debug, "Debug", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Debugf(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Debugf, "Debugf", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Trace(msg string) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogMsgHelper(l.Trace, "Trace", msg, errMsgs, i, l)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func (m multiLogger) Tracef(msg string, data ...interface{}) error {
	errMsgs := make([]string, 0)
	for i, l := range m {
		multilogfMsgHelper(l.Tracef, "Tracef", msg, errMsgs, i, l, data...)
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}
