package mlog

import (
	"bytes"
	stdlog "log"
	"testing"
)

func TestLoggerInfo(t *testing.T) {
	// Setup Logger
	l := NewLogger("main: ", nil, nil)
	l.Info("Hello, world!")

	// New Logger
	m := l.New("new: ")
	m.Info("Hello, world!")
}

func TestDefault(t *testing.T) {
	Debug("This is debug")
	Trace("This is trace")
	Info("This is info")
	Warn("This is warn")
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
