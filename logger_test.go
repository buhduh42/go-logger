package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("level=SILENT", func(t *testing.T) {
		testLevel(t, SILENT)
	})
	t.Run("level=FATAL", func(t *testing.T) {
		testLevel(t, FATAL)
		testHelperFunc(t, FATAL, int(FATAL), "Fatal", (*logger).Fatal, (*logger).Fatalf)
	})
	t.Run("level=ERROR", func(t *testing.T) {
		testLevel(t, ERROR)
		testHelperFunc(t, ERROR, int(ERROR), "Error", (*logger).Error, (*logger).Errorf)
	})
	t.Run("level=WARN", func(t *testing.T) {
		testLevel(t, WARN)
		testHelperFunc(t, WARN, int(WARN), "Warn", (*logger).Warn, (*logger).Warnf)
	})
	t.Run("level=INFO", func(t *testing.T) {
		testLevel(t, INFO)
		testHelperFunc(t, INFO, int(INFO), "Info", (*logger).Info, (*logger).Infof)
	})
	t.Run("level=DEBUG", func(t *testing.T) {
		testLevel(t, DEBUG)
		testHelperFunc(t, DEBUG, int(DEBUG), "Debug", (*logger).Debug, (*logger).Debugf)
	})
	t.Run("level=TRACE", func(t *testing.T) {
		testLevel(t, TRACE)
		testHelperFunc(t, TRACE, int(TRACE), "Trace", (*logger).Trace, (*logger).Tracef)
	})
	t.Run("multilogger", testMultiLogger)
	t.Run("multiple writers", testMultipleWriters)
}

type errorReadWriter struct {
	buff *bytes.Buffer
}

func (e *errorReadWriter) Read(p []byte) (int, error) {
	return e.buff.Read(p)
}

func (e *errorReadWriter) Write(p []byte) (int, error) {
	temp := make([]byte, len(p))
	copy(temp, p)
	e.buff = bytes.NewBuffer(temp)
	return 0, fmt.Errorf("writer failed")
}

//all the Logger methods were tested with the lovel=* tests
//no need to test them here
//TODO this test is all messed up
var multipleWriterTestCases = []*struct {
	readWriters []io.ReadWriter
	message     string
	name        string
	loggerLevel LogLevel
	toLogLevel  LogLevel
	logged      []bool
	expError    bool
}{
	{
		[]io.ReadWriter{
			new(bytes.Buffer),
			new(bytes.Buffer),
			new(bytes.Buffer),
		},
		"this is a logged message",
		"standard length 3 log writers",
		INFO,
		INFO,
		[]bool{true, true, true},
		false,
	},
	/*
		{
			[]io.ReadWriter{
				new(bytes.Buffer),
				&errorReadWriter{},
				new(bytes.Buffer),
			},
			"this is a logged message",
			"2 standard writers and one error writer",
			INFO,
			INFO,
			[]bool{true, false, true},
			true,
		},
	*/
}

func testMultipleWriters(t *testing.T) {
	for _, tCase := range multipleWriterTestCases {
		dumb := make([]io.Writer, len(tCase.readWriters))
		for i, stupid := range tCase.readWriters {
			dumb[i] = stupid.(io.Writer)
		}
		logger := NewLogger(tCase.loggerLevel, tCase.name, dumb...)
		err := logger.Log(tCase.toLogLevel, tCase.message)
		if tCase.expError && err == nil {
			t.Errorf("Logger.Log() expected an error")
		}
		if !tCase.expError && err != nil {
			t.Errorf("logger.Log() did not expect an error, got '%s'", err)
		}
		if len(tCase.logged) != len(tCase.readWriters) {
			t.Errorf(
				"test case '%s' requires equal length logged and readWriter slices",
				tCase.name,
			)
			continue
		}
		for i, w := range tCase.readWriters {
			read, _ := ioutil.ReadAll(w)
			if len(read) > 0 && !tCase.logged[i] {
				t.Logf("idx: %d, len: %d", i, len(read))
				t.Errorf(
					"test case: '%s' got a log message, should not have, got: '%s'",
					tCase.name, read,
				)
			}
			if len(read) == 0 && tCase.logged[i] {
				t.Logf("idx: %d, len: %d", i, len(read))
				t.Errorf(
					"test case '%s' should not have gotten a log messge",
					tCase.name,
				)
			}
		}
	}
}

var levels []LogLevel = []LogLevel{
	SILENT,
	FATAL,
	ERROR,
	WARN,
	INFO,
	DEBUG,
	TRACE,
}

var multiLoggerTestCases = []*struct {
	//loggers
}{}

func testMultiLogger(t *testing.T) {
	//for _, tCase := range multiLoggerTestCases {

	//}
}

func testHelperFunc(
	t *testing.T,
	level LogLevel,
	cbLevel int,
	cbName string,
	cb func(*logger, string) error,
	cbf func(*logger, string, ...interface{}) error,
) {
	buff := new(bytes.Buffer)
	temp := NewLogger(level, fmt.Sprintf("logger %s", level), buff)
	myLogger := temp.(*logger)
	err := cb(myLogger, "msg")
	if err != nil {
		t.Errorf("logger.%s() should not have returned an error, got '%s'", cbName, err)
	}
	if cbLevel <= int(level) && buff.String() == "" {
		t.Errorf("logger.%s() should have logged for logger level %s", cbName, level)
	}
	if cbLevel > int(level) && buff.String() != "" {
		t.Errorf("logger.%s() should NOT have logged for logger level %s", cbName, level)
	}
	buff.Reset()
	err = cbf(myLogger, "%s", "msg")
	if err != nil {
		t.Errorf("logger.%sf() should not have returned an error, got '%s'", cbName, err)
	}
	if cbLevel <= int(level) && buff.String() == "" {
		t.Errorf("logger.%sf() should have logged for logger level %s", cbName, level)
	}
	if cbLevel > int(level) && buff.String() != "" {
		t.Errorf("logger.%sf() should NOT have logged for logger level %s", cbName, level)
	}
}

func testLevel(t *testing.T, level LogLevel) {
	buff := new(bytes.Buffer)
	bufff := new(bytes.Buffer)
	myLogger := NewLogger(level, fmt.Sprintf("%s logger", level), buff)
	myLoggerf := NewLogger(level, fmt.Sprintf("%s loggerf", level), bufff)
	for _, l := range levels {
		buff.Reset()
		bufff.Reset()
		myLogger.Log(l, "msg")
		myLoggerf.Logf(l, "%s", "msg")
		if l == SILENT || level == SILENT {
			if buff.Len() > 0 {
				t.Errorf(
					"logger level SILENT or SILEN log should not log a message\n\tgenerated message: '%s'",
					buff.String(),
				)
			}
			if bufff.Len() > 0 {
				t.Errorf(
					"logger level SILENT or SILEN log should not logf a message\n\tgenerated message: '%s'",
					bufff.String(),
				)
			}
			continue
		}
		if int(l) <= int(level) {
			if buff.String() == "" {
				t.Errorf(
					"logger level '%s' should log for level '%s'\n\tgenerated message: '%s'",
					level, l, buff.String(),
				)
			}
			if bufff.String() == "" {
				t.Errorf(
					"logger level '%s' should logf for level '%s'\n\tgenerated message: '%s'",
					level, l, bufff.String(),
				)
			}
		}
		if int(l) > int(level) {
			if buff.Len() > 0 {
				t.Errorf(
					"logger level '%s' should NOT log for level '%s'\n\tgenerated message: '%s'",
					level, l, buff.String(),
				)
			}
			if bufff.Len() > 0 {
				t.Errorf(
					"logger level '%s' should NOT logf for level '%s'\n\tgenerated message: '%s'",
					level, l, bufff.String(),
				)
			}
		}
	}
}
