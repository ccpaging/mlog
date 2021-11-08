package mlog

import (
	"fmt"
	stdlog "log"
)

var (
	Ldebug = "DEBG "
	Ltrace = "TRAC "
	Linfo  = "INFO "
	Lwarn  = "WARN "
)

type Logger struct {
	Prefix string
	Lmap   map[string]*stdlog.Logger
}

func New(prefix string) *Logger {
	return &Logger{Prefix: prefix}
}

func (l *Logger) SetOutput(root *stdlog.Logger) *Logger {
	if l.Lmap == nil {
		l.Lmap = make(map[string]*stdlog.Logger)
	}
	w, flag := root.Writer(), root.Flags()
	l.Lmap[Ldebug] = stdlog.New(w, Ldebug+l.Prefix, flag)
	l.Lmap[Ltrace] = stdlog.New(w, Ltrace+l.Prefix, flag)
	l.Lmap[Linfo] = stdlog.New(w, Linfo+l.Prefix, flag)
	l.Lmap[Lwarn] = stdlog.New(w, Lwarn+l.Prefix, flag)
	return l
}

func (l *Logger) Add(name string, root *stdlog.Logger) *Logger {
	if l.Lmap == nil {
		l.Lmap = make(map[string]*stdlog.Logger)
	}
	w, flag := root.Writer(), root.Flags()
	l.Lmap[name] = stdlog.New(w, name+l.Prefix, flag)
	return l
}

func (l *Logger) Remove(name string) {
	delete(l.Lmap, name)
}

func (l *Logger) New(prefix string) *Logger {
	nl := &Logger{Prefix: prefix}
	if l.Lmap == nil {
		return nl
	}
	for name, ll := range l.Lmap {
		nl.Add(name, ll)
	}
	return nl
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

func (l *Logger) Debugf(format string, v ...interface{}) {
	if ll, ok := l.Lmap[Ldebug]; ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if ll, ok := l.Lmap[Ltrace]; ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if ll, ok := l.Lmap[Linfo]; ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if ll, ok := l.Lmap[Lwarn]; ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

var std = New("").SetOutput(stdlog.Default())

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

var (
	Debug = std.Debug
	Trace = std.Trace
	Info  = std.Info
	Warn  = std.Warn

	Debugf = std.Debugf
	Tracef = std.Tracef
	Infof  = std.Infof
	Warnf  = std.Warnf
)
