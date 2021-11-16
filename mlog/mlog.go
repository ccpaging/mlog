// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package mlog

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"sync/atomic"
)

const flagMask = 0x00ff
const levelMask = 0x0f00

// These are the integer logging levels used by the logger
const (
	_ int = (1 + iota) << 8
	Ldebug
	_
	Linfo
	_
	Lwarn
	Lerror
	_
)

var LevelStrings map[int]string = map[int]string{
	Ldebug: "DEBG ",
	Linfo:  "INFO ",
	Lwarn:  "WARN ",
	Lerror: "EROR ",
}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	*stdlog.Logger
	level int32
}

// New creates a new Logger.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{Logger: stdlog.New(out, prefix, flag)}
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) SetFlags(flag int) {
	l.Logger.SetFlags(flag)
	// set filter level
	atomic.StoreInt32(&l.level, int32(flag&levelMask))
}

func (l Logger) New(prefix string) *Logger {
	return &Logger{
		Logger: stdlog.New(l.Writer(), prefix, l.Flags()),
		level:  l.level,
	}
}

func (l Logger) Exit(n int) {
	os.Exit(n)
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Print(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprint(v...))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprintf(format, v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Println(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprintln(v...))
}

// Debug calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Debug(v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Ldebug]+fmt.Sprint(v...))
}

// Debugf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Debugf(format string, v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Ldebug]+fmt.Sprintf(format, v...))
}

// Debug calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Debugln(v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Ldebug]+fmt.Sprintln(v...))
}

// Info calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Info(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprint(v...))
}

// Infof calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprintf(format, v...))
}

// Infoln calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Infoln(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Linfo]+fmt.Sprintln(v...))
}

// Warn calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Warn(v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lwarn]+fmt.Sprint(v...))
}

// Warnf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Warnf(format string, v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lwarn]+fmt.Sprintf(format, v...))
}

// Warnln calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Warnln(v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lwarn]+fmt.Sprintln(v...))
}

// Error calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Error(err error, v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lerror]+err.Error()+" "+fmt.Sprint(v...))
}

// Errorf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Errorf(format string, v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lerror]+fmt.Sprintf(format, v...))
}

// Errorln calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Errorln(err error, v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, LevelStrings[Lerror]+err.Error()+" "+fmt.Sprintln(v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprint(v...))
	l.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprintf(format, v...))
	l.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprintln(v...))
	l.Exit(1)
}

// Global Logger and functions

var std = New(os.Stderr, "", stdlog.LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }
