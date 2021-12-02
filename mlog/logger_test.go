package mlog

import (
	"bytes"
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	type tester struct {
		output func(v ...interface{})
		expect string
	}

	var buf bytes.Buffer
	testlog := NewMultiLogger(log.New(&buf, "Logger: ", 0), nil)

	var tests = []tester{
		{testlog.Debug, "Logger: DEBG hello 23 world\n"},
		{testlog.Info, "Logger: INFO hello 23 world\n"},
	}

	for _, testcase := range tests {
		buf.Reset()
		testcase.output("hello ", 23, " world")
		if got, want := buf.String(), testcase.expect; got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	}
}

func TestMultiLoggerInfo(t *testing.T) {
	// Setup Logger
	var buf bytes.Buffer
	l := NewMultiLogger(log.New(&buf, "main: ", 0), nil)
	l.Info("Hello, world!")
	if got, want := buf.String(), "main: INFO Hello, world!\n"; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func BenchmarkPrint(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := log.New(&buf, "INFO ", log.LstdFlags)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Print(testString)
	}
	b.StopTimer()
}

func BenchmarkInfo(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer

	l := NewMultiLogger(log.New(&buf, "", log.LstdFlags), nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
	b.StopTimer()
}

func BenchmarkInfoPtr(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer

	l := NewMultiLogger(log.New(&buf, "", log.LstdFlags), nil)
	info := l.Info

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		info(testString)
	}
	b.StopTimer()
}
