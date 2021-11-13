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
	colors = [][]byte{
		[]byte("\033[32m"), // Debug
		[]byte("\033[36m"), // Trace
		nil,                // Info
		[]byte("\033[33m"), // Warn
		[]byte("\033[31m"), // Error
	}
	colorReset = []byte("\033[0m")
)

type ansiTermWriter struct {
	w io.Writer
}

func (t *ansiTermWriter) Write(b []byte) (n int, err error) {
	var cb []byte
	if len(b) > 0 {
		level := bytes.IndexByte([]byte("DTIWE"), b[0])
		if level >= 0 {
			cb = colors[level]
		}
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
	logger := New("", log.New(&ansiTermWriter{w: os.Stdout}, "", log.Lshortfile))
	logger.Debug("This is a debug")
	logger.Trace("This is a trace")
	logger.Info("This is a info")
	logger.Warn("This is a warn")
}
