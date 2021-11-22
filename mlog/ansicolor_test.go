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

// Write may used in io.MultiWriter. If a writer returns an error
// or n != len(b), io.ErrShortWrite, MultiWriter write operation
// stops and returns the error; it does not continue down the list.
func (t *ansiTermWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	var cb []byte
	if n > 0 {
		level := bytes.IndexByte([]byte("DTIWE"), b[0])
		if level >= 0 {
			cb = colors[level]
		}
	}
	if len(cb) == 0 {
		_, err = t.w.Write(b)
		return
	}
	var bb []byte
	bb = append(bb, cb...)
	bb = append(bb, bytes.Trim(b, "\r\n")...)
	bb = append(bb, colorReset...)
	bb = append(bb, '\n')
	_, err = t.w.Write(bb)
	return
}

func TestAnsiTerm(t *testing.T) {
	logger := NewLogger("", log.New(&ansiTermWriter{w: os.Stdout}, "", log.Lshortfile), nil)
	logger.Debug("This is a debug")
	logger.Info("This is a info")
	logger.Warn("This is a warn")
}
