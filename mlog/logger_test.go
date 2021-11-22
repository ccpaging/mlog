package mlog

import (
	"bytes"
	stdlog "log"
	"testing"
)

func TestAll(t *testing.T) {
	type tester struct {
		output func(v ...interface{})
		prefix string
		expect string
	}

	var buf bytes.Buffer
	testlog := NewLogger("Logger: ", stdlog.New(&buf, "", 0), nil)

	var tests = []tester{
		{testlog.Debug, "Logger: ", "DEBG Logger: hello 23 world\n"},
		{testlog.Info, "Logger: ", "INFO Logger: hello 23 world\n"},
		{testlog.Warn, "Logger: ", "WARN Logger: hello 23 world\n"},
	}

	for _, testcase := range tests {
		buf.Reset()
		testcase.output("hello ", 23, " world")
		if got, want := buf.String(), testcase.expect; got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	}
}

func TestLoggerInfo(t *testing.T) {
	// Setup Logger
	var buf bytes.Buffer
	l := NewLogger("main: ", stdlog.New(&buf, "", 0), nil)
	l.Info("Hello, world!")
	if got, want := buf.String(), "INFO main: Hello, world!\n"; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestLoggerNew(t *testing.T) {
	// Setup Logger
	var buf bytes.Buffer
	l := NewLogger("main: ", stdlog.New(&buf, "", 0), nil)

	// New Logger
	m := l.New("new: ")
	m.Info("Hello, world!")
	if got, want := buf.String(), "INFO new: Hello, world!\n"; got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestDefault(t *testing.T) {
	if got := Default(); &got == &std {
		t.Errorf("Default [%p] should be std [%p]", got, std)
	}
}

func BenchmarkStdlogPrint(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := stdlog.New(&buf, "INFO ", stdlog.LstdFlags)

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

	l := NewLogger("", stdlog.New(&buf, "", stdlog.LstdFlags), nil)

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

	l := NewLogger("", stdlog.New(&buf, "", stdlog.LstdFlags), nil)
	info := l.Info

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		info(testString)
	}
	b.StopTimer()
}
