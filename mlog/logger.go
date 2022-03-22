package mlog

import (
	"fmt"
	"log"
)

// There strings define log levels. You may use own strings
// to replace. Or add some new strings to define new level,
// such as Ltrace, Lerror, Lfatal, Lpanic
var (
	Ldebug string = "DEBG "
	Linfo  string = "INFO "
)

type Logger interface {
	Output(calldepth int, s string) error
}

type MultiLogger map[string]Logger

func New() MultiLogger {
	return NewMultiLogger(nil, nil)
}

// NewLogger creates a Logger with module name and level string.
func NewMultiLogger(root Logger, keys []string) MultiLogger {
	if root == nil {
		root = log.Default()
	}
	if len(keys) == 0 {
		keys = []string{Ldebug, Linfo}
	}

	ml := make(MultiLogger)
	for _, level := range keys {
		ml[level] = root
	}
	return ml
}

func (ml MultiLogger) Set(key string, value Logger) Logger {
	old := ml[key]
	ml[key] = value
	return old
}

func (ml MultiLogger) Get(key string) (l Logger, ok bool) {
	l, ok = ml[key]
	return
}

func (ml MultiLogger) Clear() {
	for key := range ml {
		delete(ml, key)
	}
}

func (ml MultiLogger) CopyFrom(in MultiLogger) {
	ml.Clear()
	for key, value := range in {
		ml[key] = value
	}
}

func (ml MultiLogger) Debug(v ...any) {
	if l, ok := ml[Ldebug]; ok {
		l.Output(2, Ldebug+fmt.Sprint(v...))
	}
}

func (ml MultiLogger) Info(v ...any) {
	if l, ok := ml[Linfo]; ok {
		l.Output(2, Linfo+fmt.Sprint(v...))
	}
}
