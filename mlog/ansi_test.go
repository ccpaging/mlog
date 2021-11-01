// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package mlog

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

// 0, Black; 1, Red; 2, Green; 3, Yellow; 4, Blue; 5, Purple; 6, Cyan; 7, White
var (
	colorDebug = []byte("\033[32m")
	colorWarn  = []byte("\033[33m")
	colorError = []byte("\033[31m")
	colorReset = []byte("\033[0m")
)

type ansiTermWriter struct {
	w io.Writer
}

func (t *ansiTermWriter) Write(b []byte) (n int, err error) {
	var cb []byte
	switch {
	case bytes.Contains(b, []byte("DEBUG")):
		cb = colorDebug
	case bytes.Contains(b, []byte("WARN ")):
		cb = colorWarn
	case bytes.Contains(b, []byte("ERROR")):
		cb = colorError
	}
	if len(cb) == 0 {
		return t.w.Write(b)
	}
	var bb []byte
	bb = append(bb, cb...)
	bb = append(bb, bytes.Trim(b, "\r\n")...)
	bb = append(bb, colorReset...)
	bb = append(bb, '\n')
	return t.w.Write(bb)
}

func TestAnsiTerm(t *testing.T) {
	w := &ansiTermWriter{w: os.Stdout}
	logger := New(w, "", log.Lshortfile)

	logger.Debug("This is a debug")
	logger.Info("This is a info")
	logger.Warn("This is a warn")
	logger.Error("This is a error")
}
