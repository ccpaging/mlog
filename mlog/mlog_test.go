// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

// These tests are too simple.

import (
	"bytes"
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

var tests = []tester{
	{Debug, Debugf, "Logger: ", "DEBG Logger: hello 23 world"},
	{Trace, Tracef, "Logger: ", "TRAC Logger: hello 23 world"},
	{Info, Infof, "Logger: ", "INFO Logger: hello 23 world"},
	{Warn, Warnf, "Logger: ", "WARN Logger: hello 23 world"},
	{Error, Errorf, "Logger: ", "EROR Logger: hello 23 world"},
}

// Test using Debug("hello", 23, "world") or using Debugf("hello %d world", 23)
func testExtPrint(t *testing.T, testcase *tester) {
	buf := new(bytes.Buffer)
	std.SetOutput(buf)
	std.SetFlags(0)
	std.SetPrefix(testcase.prefix)
	testcase.output("hello ", 23, " world")
	testcase.outputf("hello %d world", 23)
	line := buf.String()
	line = line[0 : len(line)-1]
	if got, want := line, testcase.expect+"\n"+testcase.expect; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
	std.SetOutput(os.Stderr)
}

func TestExtAll(t *testing.T) {
	for _, testcase := range tests {
		testExtPrint(t, &testcase)
	}
}

func TestLevelSetting(t *testing.T) {
	buf := new(bytes.Buffer)
	std.SetOutput(buf)
	std.SetPrefix("Reality: ")
	std.SetLevel(Linfo)
	// Verify a log message looks right, with our prefix and microseconds present.
	std.Debug("hello")
	std.Info("hello")
	if expect := "INFO Reality: hello\n"; buf.String() != expect {
		t.Errorf("log output should match %q is %q", expect, buf.String())
	}
	std.SetOutput(os.Stderr)
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(log.New(&b, "", 0))
	l.Print(testString)
	if expect := "INFO " + testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := New(log.New(&b, "Header:", log.LstdFlags))
	l.Print()
	l.Print("non-empty")
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
	l := New(log.New(&b, "", 0))
	for i := 0; i < 100; i++ {
		go func() {
			l.SetFlags(0)
		}()
		l.Info(0, "")
	}
}

func BenchmarkCaller(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(log.New(&buf, "", log.LstdFlags|log.Lshortfile))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Print(testString)
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(log.New(&buf, "", log.LstdFlags))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Print(testString)
	}
}

func BenchmarkPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(log.New(&buf, "", 0))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Print(testString)
	}
}

func BenchmarkInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(log.New(&buf, "", log.LstdFlags))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}
