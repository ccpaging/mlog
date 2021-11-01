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
	Ldebug int = (1 + iota) << 8
	Ltrace
	Linfo
	Lnotice
	Lwarn
	Lerror
	Lfatal
	Lpanic
)

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

var std = New(os.Stderr, "", stdlog.LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) SetFlags(flag int) {
	l.Logger.SetFlags(flag)
	// set filter level
	atomic.StoreInt32(&l.level, int32(flag&levelMask))
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "INFO  "+fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "INFO  "+fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "INFO  "+fmt.Sprintln(v...))
}

// Fatal is equivalent to l.Print().
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprint(v...))
}

// Fatalf is equivalent to l.Printf().
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprintf(format, v...))
}

// Fatalln is equivalent to l.Println().
func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(2, "FATAL "+fmt.Sprintln(v...))
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	l.Output(2, "PANIC "+fmt.Sprint(v...))
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Output(2, "PANIC "+fmt.Sprintf(format, v...))
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	l.Output(2, "PANIC "+fmt.Sprintln(v...))
}

// Debugf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "DEBUG "+fmt.Sprintf(format, v...))
}

// Debug calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "DEBUG "+fmt.Sprintln(v...))
}

// Infof calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "INFO  "+fmt.Sprintf(format, v...))
}

// Info calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "INFO  "+fmt.Sprintln(v...))
}

// Warnf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "WARN  "+fmt.Sprintf(format, v...))
}

// Warn calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warn(v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "WARN  "+fmt.Sprintln(v...))
}

// Errorf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "ERROR "+fmt.Sprintf(format, v...))
}

// Error calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&l.level) {
		return
	}
	l.Output(2, "ERROR "+fmt.Sprintln(v...))
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// Flags returns the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// Writer returns the output destination for the standard logger.
func Writer() io.Writer {
	return std.Writer()
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Output(2, "INFO  "+fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Output(2, "INFO  "+fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Output(2, "INFO  "+fmt.Sprintln(v...))
}

// Fatal is equivalent to Print().
func Fatal(v ...interface{}) {
	std.Output(2, "FATAL "+fmt.Sprint(v...))
}

// Fatalf is equivalent to Printf().
func Fatalf(format string, v ...interface{}) {
	std.Output(2, "FATAL "+fmt.Sprintf(format, v...))
}

// Fatalln is equivalent to Println().
func Fatalln(v ...interface{}) {
	std.Output(2, "FATAL "+fmt.Sprintln(v...))
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	std.Output(2, "PANIC "+fmt.Sprint(v...))
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	std.Output(2, "PANIC "+fmt.Sprintf(format, v...))
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	std.Output(2, "PANIC "+fmt.Sprintln(v...))
}

// These functions write to the standard logger.

// Debug calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "DEBUG "+fmt.Sprintln(v...))
}

// Debugf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	if int32(Ldebug) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "DEBUG "+fmt.Sprintf(format, v...))
}

// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "INFO  "+fmt.Sprintln(v...))
}

// Infof calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	if int32(Linfo) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "INFO  "+fmt.Sprintf(format, v...))
}

// Warn calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "WARN  "+fmt.Sprintln(v...))
}

// Warnf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	if int32(Lwarn) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "WARN  "+fmt.Sprintf(format, v...))
}

// Error calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "ERROR "+fmt.Sprintln(v...))
}

// Errorf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	if int32(Lerror) < atomic.LoadInt32(&std.level) {
		return
	}
	std.Output(2, "ERROR "+fmt.Sprintf(format, v...))
}
