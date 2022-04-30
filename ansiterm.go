package mlog

import (
	"bytes"
	"io"
)

// 0, Black; 1, Red; 2, Green; 3, Yellow; 4, Blue; 5, Purple; 6, Cyan; 7, White
var (
	colorDebug = []byte("\033[32m")
	colorTrace = []byte("\033[35m")
	colorWarn  = []byte("\033[33m")
	colorError = []byte("\033[31m")
	colorReset = []byte("\033[0m")
)

type ansiTerm struct {
	w io.Writer
}

func (t *ansiTerm) Write(b []byte) (n int, err error) {
	n = len(b)
	var cb []byte
	if n > 20 {
		switch b[20] {
		case 'D':
			cb = colorDebug
		case 'T':
			cb = colorTrace
		case 'W':
			cb = colorWarn
		case 'E':
			cb = colorError
		case 'F':
			cb = colorError
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
