package mlog

import (
	"fmt"
	stdlog "log"
)

var (
	Ldebug string = "DEBG "
	Linfo  string = "INFO "
	Lwarn  string = "WARN "
)

type Logger map[string]*stdlog.Logger

func New(name string) Logger {
	return NewLogger(name, nil, nil)
}

func NewLogger(name string, root *stdlog.Logger, levelStrings []string) Logger {
	if root == nil {
		root = stdlog.Default()
	}
	if len(levelStrings) == 0 {
		levelStrings = []string{Ldebug, Linfo, Lwarn}
	}

	w, flag := root.Writer(), root.Flags()
	l := make(Logger)
	for _, level := range levelStrings {
		l[level] = stdlog.New(w, level+name, flag)
	}
	return l
}

func (l Logger) New(name string) Logger {
	out := make(Logger)
	for level, ll := range l {
		w, flag := ll.Writer(), ll.Flags()
		out[level] = stdlog.New(w, level+name, flag)
	}
	return out
}

func (l Logger) Debug(v ...interface{}) {
	if ll, ok := l[Ldebug]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l Logger) Info(v ...interface{}) {
	if ll, ok := l[Linfo]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l Logger) Warn(v ...interface{}) {
	if ll, ok := l[Lwarn]; ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

var (
	std   = NewLogger("", nil, nil)
	Debug = std.Debug
	Info  = std.Info
	Warn  = std.Warn
)

// Default returns the standard logger used by the package-level output functions.
func Default() Logger { return std }
