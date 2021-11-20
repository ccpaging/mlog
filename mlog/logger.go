package mlog

import (
	"fmt"
	stdlog "log"
)

var (
	Ldebug string = "DEBG "
	Ltrace string = "TRAC "
	Linfo  string = "INFO "
	Lwarn  string = "WARN "
	Lerror string = "EROR "
)

type Logger struct {
	Name string
	Lmap map[string]*stdlog.Logger
}

func NewLogger(name string, root *stdlog.Logger, levelStrings []string) *Logger {
	if root == nil {
		root = stdlog.Default()
	}
	if levelStrings == nil {
		levelStrings = []string{Ldebug, Ltrace, Linfo, Lwarn, Lerror}
	}

	w, flag := root.Writer(), root.Flags()
	l := &Logger{
		Name: name,
		Lmap: make(map[string]*stdlog.Logger),
	}
	for _, level := range levelStrings {
		l.Lmap[level] = stdlog.New(w, l.Name+level, flag)
	}
	return l
}

func (l Logger) New(name string) *Logger {
	o := &Logger{
		Name: name,
		Lmap: make(map[string]*stdlog.Logger),
	}
	for level, ll := range l.Lmap {
		w, flag := ll.Writer(), ll.Flags()
		o.Lmap[level] = stdlog.New(w, l.Name+level, flag)
	}
	return o
}

func (l *Logger) Debug(v ...interface{}) {
	if ll, ok := l.Lmap[Ldebug]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	if ll, ok := l.Lmap[Ltrace]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if ll, ok := l.Lmap[Linfo]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if ll, ok := l.Lmap[Lwarn]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	if ll, ok := l.Lmap[Lerror]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

var (
	std   = NewLogger("", nil, nil)
	Debug = std.Debug
	Trace = std.Trace
	Info  = std.Info
	Warn  = std.Warn
)
