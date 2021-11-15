package mlog

import (
	"fmt"
	"log"
	"sync"
)

var (
	Ldebug = "DEBG "
	Ltrace = "TRAC "
	Linfo  = "INFO "
	Lwarn  = "WARN "
)

type Logger struct {
	mu sync.RWMutex
	ml *MultiLogger
}

func NewLogger(prefix string, root *log.Logger) *Logger {
	ml := New(prefix)
	ml.Set(Ldebug, root)
	ml.Set(Ltrace, root)
	ml.Set(Linfo, root)
	ml.Set(Lwarn, root)
	return &Logger{ml: ml}
}

func (l *Logger) New(prefix string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	return &Logger{ml: l.ml.New(prefix)}
}

func (l *Logger) Debug(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if ll, ok := l.ml.Get(Ldebug); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if ll, ok := l.ml.Get(Ltrace); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if ll, ok := l.ml.Get(Linfo); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if ll, ok := l.ml.Get(Lwarn); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

var (
	std   = NewLogger("", log.Default())
	Debug = std.Debug
	Trace = std.Trace
	Info  = std.Info
	Warn  = std.Warn
)
