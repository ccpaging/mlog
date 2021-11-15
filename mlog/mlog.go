package mlog

import (
	stdlog "log"
)

type MultiLogger struct {
	pre string
	lm  map[string]*stdlog.Logger
}

// New creates a new Logger wrapper with new prefix,
// and copy the root go stdlib logger to internal map.
func New(prefix string) *MultiLogger {
	return &MultiLogger{
		pre: prefix,
		lm:  make(map[string]*stdlog.Logger),
	}
}

// CreateByMap creates a new logger wrapper with new prefix,
// and copy the map go stdlib logger to internal. The map keys
// should be Ldebug, Ltrace, Linfo, Lwarn.
func CreateByMap(prefix string, lm map[string]*stdlog.Logger) *MultiLogger {
	l := New(prefix)
	for s, ll := range lm {
		l.Set(s, ll)
	}
	return l
}

func (l *MultiLogger) Set(level string, root *stdlog.Logger) *MultiLogger {
	w, flag := root.Writer(), root.Flags()
	l.lm[level] = stdlog.New(w, level+l.pre, flag)
	return l
}

func (l *MultiLogger) remove(level string) {
	delete(l.lm, level)
}

func (l *MultiLogger) Get(level string) (ll *stdlog.Logger, ok bool) {
	ll, ok = l.lm[level]
	return
}

// New creates a new logger wrapper with new prefix.
func (l *MultiLogger) New(prefix string) *MultiLogger {
	return CreateByMap(prefix, l.lm)
}
