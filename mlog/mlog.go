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
	pre string
	lm  map[string]*stdlog.Logger
}

// New creates a new Logger wrapper with new prefix,
// and copy the root go stdlib logger to internal map.
func New(prefix string, root *stdlog.Logger) *Logger {
	l := &Logger{pre: prefix}
	l.add(Ldebug, root)
	l.add(Ltrace, root)
	l.add(Linfo, root)
	l.add(Lwarn, root)
	return l
}

// CreateByMap creates a new logger wrapper with new prefix,
// and copy the map go stdlib logger to internal. The map keys
// should be Ldebug, Ltrace, Linfo, Lwarn.
func CreateByMap(prefix string, lm map[string]*stdlog.Logger) *Logger {
	l := &Logger{pre: prefix}
	for s, ll := range lm {
		l.add(s, ll)
	}
	return l
}

func (l *Logger) add(name string, root *stdlog.Logger) *Logger {
	if l.lm == nil {
		l.lm = make(map[string]*stdlog.Logger)
	}
	w, flag := root.Writer(), root.Flags()
	l.lm[name] = stdlog.New(w, name+l.pre, flag)
	return l
}

func (l *Logger) remove(name string) {
	delete(l.lm, name)
}

func (l *Logger) get(name string) (ll *stdlog.Logger, ok bool) {
	ll, ok = l.lm[name]
	return
}

// New creates a new logger wrapper with new prefix.
func (l *Logger) New(prefix string) *Logger {
	nl := &Logger{pre: prefix}
	if l.lm == nil {
		return nl
	}
	for name, ll := range l.lm {
		nl.add(name, ll)
	}
	return nl
}

func (l *Logger) Debug(v ...interface{}) {
	if ll, ok := l.get(Ldebug); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	if ll, ok := l.get(Ltrace); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if ll, ok := l.get(Linfo); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if ll, ok := l.get(Lwarn); ok {
		ll.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if ll, ok := l.get(Ldebug); ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if ll, ok := l.get(Ltrace); ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if ll, ok := l.get(Linfo); ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if ll, ok := l.get(Lwarn); ok {
		ll.Output(2, fmt.Sprintf(format, v...))
	}
}

var std = New("", stdlog.Default())

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
