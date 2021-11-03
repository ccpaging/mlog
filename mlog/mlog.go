// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package mlog

import (
	"io"
	stdlog "log"
	"os"
)

const (
	Ldebug int = iota
	Ltrace
	Linfo
	Lwarn
	Lerror
)

var levelStrings = []string{"DEBG", "TRAC", "INFO", "WARN", "EROR"}

type Logger struct {
	*LogBase
	Printf OutfFunc
	Debugf OutfFunc
	Tracef OutfFunc
	Infof  OutfFunc
	Warnf  OutfFunc
	Errorf OutfFunc

	Print OutFunc
	Debug OutFunc
	Trace OutFunc
	Info  OutFunc
	Warn  OutFunc
	Error OutFunc
}

// New creates a new Logger.
func New(out io.Writer, prefix string, flag int) *Logger {
	b := NewLogBase(out, prefix, flag, levelStrings)
	return &Logger{
		LogBase: b,

		Printf: b.Outf2, // defult level function
		Debugf: b.Outf0,
		Tracef: b.Outf1,
		Infof:  b.Outf2,
		Warnf:  b.Outf3,
		Errorf: b.Outf4,

		Print: b.Out2, // defult level function
		Debug: b.Out0,
		Trace: b.Out1,
		Info:  b.Out2,
		Warn:  b.Out3,
		Error: b.Out4,
	}
}

// New creates a new Logger.
func (l *Logger) New(prefix string) *Logger {
	return &Logger{LogBase: l.LogBase.New(prefix)}
}

var std = New(os.Stderr, "", stdlog.LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

var Printf = std.Outf2 // defult level function
var Debugf = std.Outf0
var Tracef = std.Outf1
var Infof = std.Outf2
var Warnf = std.Outf3
var Errorf = std.Outf4

var Print = std.Out2 // defult level function
var Debug = std.Out0
var Trace = std.Out1
var Info = std.Out2
var Warn = std.Out3
var Error = std.Out4
