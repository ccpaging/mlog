package mlog

import (
	"bytes"
	"os"
	"testing"
)

type tester struct {
	output  func(v ...interface{})
	outputf func(format string, v ...interface{})
	prefix  string
	expect  string
}

var tests = []tester{
	{Debug, Debugf, "Logger: ", "Logger: DEBUG hello 23 world"},
	{Info, Infof, "Logger: ", "Logger: INFO  hello 23 world"},
	{Warn, Warnf, "Logger: ", "Logger: WARN  hello 23 world"},
	{Error, Errorf, "Logger: ", "Logger: ERROR hello 23 world"},
}

// Test using Println("hello", 23, "world") or using Printf("hello %d world", 23)
func testExtPrint(t *testing.T, testcase *tester) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	SetFlags(0)
	SetPrefix(testcase.prefix)
	testcase.output("hello", 23, "world")
	testcase.outputf("hello %d world", 23)
	line := buf.String()
	line = line[0 : len(line)-1]
	if got, want := line, testcase.expect+"\n"+testcase.expect; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
	SetOutput(os.Stderr)
}

func TestExtAll(t *testing.T) {
	for _, testcase := range tests {
		testExtPrint(t, &testcase)
	}
}

func TestLevelSetting(t *testing.T) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	level := int(std.level)
	if level != 0 {
		t.Errorf(`Level: expected %d got %d`, 0, level)
	}
	SetFlags(Flags() | Linfo)
	level = int(std.level)
	if level != Linfo {
		t.Errorf(`Prefix: expected %d got %d`, Linfo, level)
	}
	SetPrefix("Reality:")
	// Verify a log message looks right, with our prefix and microseconds present.
	Debug("hello")
	Info("hello")
	if expect := "Reality:"; buf.String()[0:len(expect)] != expect {
		t.Errorf("log output should match %q is %q", expect, buf.String()[0:len(expect)])
	}
	SetOutput(os.Stderr)
}

func TestLevelRace(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "", 0)
	for i := 0; i < 100; i++ {
		go func() {
			l.SetFlags(Linfo)
		}()
		l.Info(0, "")
	}
}
