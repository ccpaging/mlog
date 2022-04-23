package log_test

import (
	"bytes"
	stdlog "log"
	"testing"

	"github.com/ccpaging/log"
)

func TestStdLogAt(t *testing.T) {
	var buf bytes.Buffer
	l := log.New("StdLogAt: ", stdlog.New(&buf, "", 0), log.Linfo)

	logAtInfo := l.StdLogAt("info", "test: ")
	logAtInfo.SetFlags(stdlog.Lshortfile)

	logAtInfo.Println("This is stdlog's Println")
	if want, got := "INFO test: stdlog_test.go:18: This is stdlog's Println\n", buf.String(); want != got {
		t.Errorf("\nwant: %q\ngot:  %q", want, got)
	}

	buf.Reset()
	logAtDebug := l.StdLogAt("debug", "test: ")
	logAtDebug.Println("This is stdlog's debug")
	if buf.String() != "" {
		t.Errorf("\nwant empty\ngot: %q", buf.String())
	}
}
