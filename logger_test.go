package mlog

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testFiles []string = []string{"_test.log", "_test.1.log"}

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	root := log.New(&buf, "", 0)
	l := New("main: ", root, "")
	l.Info("This is info")
	if want, got := "INFO main: This is info\n", buf.String(); want != got {
		t.Errorf("logger output should match %q is %q", want, got)
	}
}

func TestNewLogger(t *testing.T) {
	testFile := testFiles[0]
	defer os.Remove(testFile)

	prev := COLOR
	COLOR = true

	l := NewLogger("test ",
		&Settings{
			EnableConsole: true, ConsoleLevel: "debug",
			EnableFile: true, FileLevel: "info", FileLocation: testFile, FileLimitSize: "5k", FileBackupCount: 1,
		})
	l.Debug("this is debug")
	l.Trace("this is trace")
	l.Info("this is info")
	l.Warn("this is warn")
	l.Error(errors.New("this is error"))
	l.Close()
	COLOR = prev
	l.Info("Omitted after Close()")

	if contents, err := ioutil.ReadFile(testFile); err != nil {
		t.Errorf("read(%q): %s", testFiles, err)
	} else if len(contents) != 130 {
		t.Errorf("malformed file: %q (%d bytes)", string(contents), len(contents))
	}
}

func TestLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	l := New("test: ", log.New(&buf, "", log.Lshortfile|log.Lmsgprefix), "")
	l.Debug("This is debug")
	if want, got := "logger_test.go:55: DEBG test: This is debug\n", buf.String(); want != got {
		t.Errorf("logger debug should match %q is %q", want, got)
	}
}

func TestLoggerNew(t *testing.T) {
	l := New("new: ", nil, "")

	var buf bytes.Buffer
	root := log.New(&buf, "", 0)
	ll := New("test: ", root, "")
	dup := ll.WithName("temp")
	if want, got := &ll.core, &dup.core; want == got {
		t.Errorf("logger new should has a new core %v is %v", want, got)
	}

	p1 := &l.core
	l.CopyFrom(dup)
	p2 := &l.core
	if p1 != p2 {
		t.Errorf("logger copyfrom should match %v is %v", p1, p2)
	}

	l.Info("This is info")
	if want, got := "INFO new: This is info\n", buf.String(); want != got {
		t.Errorf("logger output should match %q is %q", want, got)
	}
}

func BenchmarkStdlogPrint(b *testing.B) {
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

	l := New("", log.New(&buf, "", log.LstdFlags), "")

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

	l := New("", log.New(&buf, "", log.LstdFlags), "")
	info := l.Info

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		info(testString)
	}
	b.StopTimer()
}
