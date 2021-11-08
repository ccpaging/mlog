// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

package mlog

import (
	stdlog "log"
)

const (
	Ldebug int = iota
	Ltrace
	Linfo
	Lwarn
	Lerror
)

var levelStrings = []string{"DEBG ", "TRAC ", "INFO ", "WARN ", "EROR "}

// OutFunc handles the arguments in the manner of fmt.Print.
type OutFunc func(v ...interface{})

// OutlnFunc handles the arguments in the manner of fmt.Println.
type OutlnFunc func(v ...interface{})

// OutfFunc handles the arguments in the manner of fmt.Printf.
type OutfFunc func(format string, v ...interface{})

// A LogOutput represents the log output functions.
// Those functions have three types:
//	OutFunc func(v ...interface{})
//	OutlnFunc func(v ...interface{})
//	OutfFunc func(format string, v ...interface{})
type LogOutput struct {
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

// New define the real function by the loggers bundle.
func (*LogOutput) New(b *LogBundle) *LogOutput {
	return &LogOutput{
		Printf: b.Outf2, // set as Infof
		Debugf: b.Outf0,
		Tracef: b.Outf1,
		Infof:  b.Outf2,
		Warnf:  b.Outf3,
		Errorf: b.Outf4,

		Print: b.Out2, // set as Info
		Debug: b.Out0,
		Trace: b.Out1,
		Info:  b.Out2,
		Warn:  b.Out3,
		Error: b.Out4,
	}
}

// A Logger is the wrapper of LogBundle and LogOutput.
type Logger struct {
	*LogBundle
	*LogOutput
}

// New creates a new Logger.
func New(l *stdlog.Logger) *Logger {
	b := Bundle(l, levelStrings)
	o := &LogOutput{}
	return &Logger{
		LogBundle: b,
		LogOutput: o.New(b),
	}
}

// New creates a new duplicate logger with new prefix.
func (l *Logger) New(prefix string) *Logger {
	b := l.LogBundle.New(prefix)
	o := &LogOutput{}
	return &Logger{
		LogBundle: b,
		LogOutput: o.New(b),
	}
}

var std = New(stdlog.Default())

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

// These functions write to the standard logger.

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
var Printf = std.Printf
var Debugf = std.Debugf
var Tracef = std.Tracef
var Infof = std.Infof
var Warnf = std.Warnf
var Errorf = std.Errorf

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
var Print = std.Print
var Debug = std.Debug
var Trace = std.Trace
var Info = std.Info
var Warn = std.Warn
var Error = std.Error
