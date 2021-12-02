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
func NewMultiLogger(root Logger, levelStrings []string) MultiLogger {
	if root == nil {
		root = log.Default()
	}
	if len(levelStrings) == 0 {
		levelStrings = []string{Ldebug, Linfo}
	}

	ml := make(MultiLogger)
	for _, level := range levelStrings {
		ml[level] = root
	}
	return ml
}

func (ml MultiLogger) Set(level string, l Logger) Logger {
	old := ml[level]
	ml[level] = l
	return old
}

func (ml MultiLogger) Get(level string) (l Logger, ok bool) {
	l, ok = ml[level]
	return
}

func (ml MultiLogger) Clear() {
	for key := range ml {
		delete(ml, key)
	}
}

func (ml MultiLogger) CopyFrom(in MultiLogger) {
	ml.Clear()
	for key, l := range in {
		ml[key] = l
	}
}

func (ml MultiLogger) Debug(v ...interface{}) {
	if l, ok := ml[Ldebug]; ok {
		l.Output(2, Ldebug+fmt.Sprint(v...))
	}
}

func (ml MultiLogger) Info(v ...interface{}) {
	if l, ok := ml[Linfo]; ok {
		l.Output(2, Linfo+fmt.Sprint(v...))
	}
}
