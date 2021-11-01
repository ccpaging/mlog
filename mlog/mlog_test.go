// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

// These tests are too simple.

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	if got := Default(); got != std {
		t.Errorf("Default [%p] should be std [%p]", got, std)
	}
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(&b, "", 0)
	l.Println(testString)
	if expect := "INFO  " + testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestOutputRace(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "", 0)
	for i := 0; i < 100; i++ {
		go func() {
			l.SetFlags(0)
		}()
		l.Output(0, "")
	}
}

func TestFlagAndPrefixSetting(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "Test:", log.LstdFlags)
	f := l.Flags()
	if f != log.LstdFlags {
		t.Errorf("Flags 1: expected %x got %x", log.LstdFlags, f)
	}
	l.SetFlags(f | log.Lmicroseconds)
	f = l.Flags()
	if f != log.LstdFlags|log.Lmicroseconds {
		t.Errorf("Flags 2: expected %x got %x", log.LstdFlags|log.Lmicroseconds, f)
	}
	p := l.Prefix()
	if p != "Test:" {
		t.Errorf(`Prefix: expected "Test:" got %q`, p)
	}
	l.SetPrefix("Reality:")
	p = l.Prefix()
	if p != "Reality:" {
		t.Errorf(`Prefix: expected "Reality:" got %q`, p)
	}
	// Verify a log message looks right, with our prefix and microseconds present.
	l.Print("hello")
	if expect := "Reality:"; b.String()[0:len(expect)] != expect {
		t.Errorf("log output should match %q is %q", expect, b.String()[0:len(expect)])
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
