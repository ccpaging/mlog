// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

// These tests are too simple.

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

type tester struct {
	output  func(v ...interface{})
	outputf func(format string, v ...interface{})
	prefix  string
	expect  string
}

var testLogger = New(os.Stderr, "", log.LstdFlags)

var tests = []tester{
	{testLogger.Debugln, testLogger.Debugf, "Logger: ", "DEBG Logger: hello 23 world"},
	{testLogger.Traceln, testLogger.Tracef, "Logger: ", "TRAC Logger: hello 23 world"},
	{testLogger.Infoln, testLogger.Infof, "Logger: ", "INFO Logger: hello 23 world"},
	{testLogger.Warnln, testLogger.Warnf, "Logger: ", "WARN Logger: hello 23 world"},
	{testLogger.Errorln, testLogger.Errorf, "Logger: ", "EROR Logger: hello 23 world"},
}

// Test using Println("hello", 23, "world") or using Printf("hello %d world", 23)
func testExtPrint(t *testing.T, testcase *tester) {
	buf := new(bytes.Buffer)
	testLogger.SetOutput(buf)
	testLogger.SetFlags(0)
	testLogger.SetPrefix(testcase.prefix)
	testcase.output("hello", 23, "world")
	testcase.outputf("hello %d world", 23)
	line := buf.String()
	line = line[0 : len(line)-1]
	if got, want := line, testcase.expect+"\n"+testcase.expect; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
	testLogger.SetOutput(os.Stderr)
}

func TestExtAll(t *testing.T) {
	for _, testcase := range tests {
		testExtPrint(t, &testcase)
	}
}

func TestLevelSetting(t *testing.T) {
	buf := new(bytes.Buffer)
	testLogger.SetOutput(buf)
	testLogger.SetPrefix("Reality: ")
	testLogger.SetLevelOutput(Ldebug, io.Discard)
	// Verify a log message looks right, with our prefix and microseconds present.
	testLogger.Debugln("hello")
	testLogger.Infoln("hello")
	if expect := "INFO Reality: hello\n"; buf.String() != expect {
		t.Errorf("log output should match %q is %q", expect, buf.String())
	}
	testLogger.SetOutput(os.Stderr)
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(&b, "", 0)
	l.Println(testString)
	if expect := "INFO " + testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "Header:", log.LstdFlags)
	l.Print()
	l.Println("non-empty")
	output := b.String()
	if n := strings.Count(output, "Header"); n != 2 {
		t.Errorf("expected 2 headers, got %d", n)
	}
	if n := strings.Count(output, "\n"); n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func TestLevelRace(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "", 0)
	for i := 0; i < 100; i++ {
		go func() {
			l.SetFlags(0)
		}()
		l.Infoln(0, "")
	}
}

func BenchmarkCaller(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", log.LstdFlags|log.Lshortfile)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", log.LstdFlags)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", 0)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", log.LstdFlags)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Infoln(testString)
	}
}
