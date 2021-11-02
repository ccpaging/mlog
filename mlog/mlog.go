// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package mlog

import (
	"fmt"
	"io"
	stdlog "log"
	"sync"
)

type Level int

const (
	Ldebug Level = iota
	Ltrace
	Linfo
	Lwarn
	Lerror
)

var levelString = [Lerror + 1]string{"DEBG", "TRAC", "INFO", "WARN", "EROR"}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu    sync.Mutex // ensures atomic writes; protects the following fields
	slots [Lerror + 1]*stdlog.Logger
}

// New creates a new Logger.
func New(out io.Writer, prefix string, flag int) *Logger {
	l := &Logger{}
	for i := 0; i < int(Lerror)+1; i++ {
		l.slots[i] = stdlog.New(out, levelString[i]+" "+prefix, flag)
	}
	return l
}

// New creates a new Logger.
func (l *Logger) New(prefix string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := l.slots[0].Writer()
	flag := l.slots[0].Flags()

	ll := &Logger{}
	for i := 0; i < int(Lerror)+1; i++ {
		ll.slots[i] = stdlog.New(out, levelString[i]+" "+prefix, flag)
	}
	return ll
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, s := range l.slots {
		s.SetOutput(w)
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetLevelOutput(level Level, w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.slots[level].SetOutput(w)
}

// Writer returns the output destination for the logger.
func (l *Logger) LevelWriter(level Level) io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.slots[level].Writer()
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, s := range l.slots {
		s.SetFlags(flag)
	}
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) SetLevelFlags(level Level, flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.slots[level].SetFlags(flag)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) LevelFlags(level Level) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.slots[level].Flags()
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := 0; i < int(Lerror)+1; i++ {
		l.slots[i].SetPrefix(levelString[i] + " " + prefix)
	}
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetLevelPrefix(level Level, prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.slots[level].SetPrefix(prefix)
}

// Prefix returns the output prefix for the logger.
func (l *Logger) LevelPrefix(level Level) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.slots[level].Prefix()
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer(level Level) io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.slots[level].Writer()
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.slots[Linfo].Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.slots[Linfo].Output(2, fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	l.slots[Linfo].Output(2, fmt.Sprintln(v...))
}

// Debugf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.slots[Ldebug].Output(2, fmt.Sprintf(format, v...))
}

// Debug calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debugln(v ...interface{}) {
	l.slots[Ldebug].Output(2, fmt.Sprintln(v...))
}

// Tracef calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.slots[Ltrace].Output(2, fmt.Sprintf(format, v...))
}

// Trace calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Traceln(v ...interface{}) {
	l.slots[Ltrace].Output(2, fmt.Sprintln(v...))
}

// Infof calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.slots[Linfo].Output(2, fmt.Sprintf(format, v...))
}

// Info calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Infoln(v ...interface{}) {
	l.slots[Linfo].Output(2, fmt.Sprintln(v...))
}

// Warnf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.slots[Lwarn].Output(2, fmt.Sprintf(format, v...))
}

// Warn calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warnln(v ...interface{}) {
	l.slots[Lwarn].Output(2, fmt.Sprintln(v...))
}

// Errorf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.slots[Lerror].Output(2, fmt.Sprintf(format, v...))
}

// Error calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Errorln(v ...interface{}) {
	l.slots[Lerror].Output(2, fmt.Sprintln(v...))
}
