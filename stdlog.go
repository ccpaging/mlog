package mlog

import (
	"io"
	"log"
)

func (l *Logger) StdLog(name string) *log.Logger {
	return l.StdLogAt(Ldebug, name)
}

// StdLogAt returns *log.Logger which writes to supplied zap logger at required level.
func (l *Logger) StdLogAt(level, name string) *log.Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	n := ltoi(level)
	prefix := levelStrings[n] + name
	if w := l.levelWriter(n); w != nil {
		return log.New(w, prefix, LstdFlags)
	}

	return log.New(io.Discard, prefix, LstdFlags)
}

// NewStdLog returns a *log.Logger which writes to the supplied zap Logger at
// InfoLevel. To redirect the standard library's package-global logging
// functions, use RedirectStdLog instead.
func NewStdLog(name string) *log.Logger {
	return global.StdLog(name)
}

// NewStdLogAt returns *log.Logger which writes to supplied zap logger at
// required level.
func NewStdLogAt(level, name string) *log.Logger {
	return global.StdLogAt(level, name)
}
