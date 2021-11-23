package mlog

import (
	"fmt"
	stdlog "log"
)

// There strings define log levels. You may use own strings
// to replace. Or add some new strings to define new level,
// such as Ltrace, Lerror, Lfatal, Lpanic
var (
	Ldebug string = "DEBG "
	Linfo  string = "INFO "
)

type Logger map[string]*stdlog.Logger

func New() Logger {
	return NewLogger(nil, nil)
}

// NewLogger creates a Logger with module name and level string.
func NewLogger(root *stdlog.Logger, levelStrings []string) Logger {
	if root == nil {
		root = stdlog.Default()
	}
	if len(levelStrings) == 0 {
		levelStrings = []string{Ldebug, Linfo}
	}

	w, pr, flag := root.Writer(), root.Prefix(), root.Flags()
	l := make(Logger)
	for _, level := range levelStrings {
		l[level] = stdlog.New(w, pr, flag)
	}
	return l
}

func (l Logger) Set(level string, ll *stdlog.Logger) *stdlog.Logger {
	old := l[level]
	l[level] = ll
	return old
}

func (l Logger) Get(level string) (ll *stdlog.Logger, ok bool) {
	ll, ok = l[level]
	return
}

func (l Logger) Clear() {
	for key := range l {
		delete(l, key)
	}
}

func (l Logger) CopyFrom(in Logger) {
	l.Clear()
	for key, ll := range in {
		l[key] = ll
	}
}

func (l Logger) Debug(v ...interface{}) {
	if ll, ok := l[Ldebug]; ok {
		ll.Output(2, Ldebug+fmt.Sprint(v...))
	}
}

func (l Logger) Info(v ...interface{}) {
	if ll, ok := l[Linfo]; ok {
		ll.Output(2, Linfo+fmt.Sprint(v...))
	}
}
