// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package syslog

import (
	"fmt"
	stdlog "log"
)

// The Priority is a combination of the syslog facility and
// severity. For example, LOG_ALERT | LOG_FTP sends an alert severity
// message from the FTP facility. The default severity is LOG_EMERG;
// the default facility is LOG_KERN.
type Priority int

const severityMask = 0x07
const facilityMask = 0xf8

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to LOG_FTP on Linux, BSD, and OS X.
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

type Logger struct {
	*stdlog.Logger
	out *Writer
}

// New establishes a new connection to the system log daemon. Each
// write to the returned writer sends a log message with the given
// priority (a combination of the syslog facility and severity) and
// prefix tag. If tag is empty, the os.Args[0] is used.
func New(priority Priority, tag string) (*Logger, error) {
	out, err := Dial("", "", priority, tag)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: stdlog.New(out, "", 0),
		out:    out,
	}, nil
}

// NewLogger creates a log.Logger whose output is written to the
// system log service with the specified priority, a combination of
// the syslog facility and severity. The logFlag argument is the flag
// set passed through to log.New to create the Logger.
func NewLogger(out *Writer, logFlag int) *Logger {
	return &Logger{
		Logger: stdlog.New(out, "", logFlag),
		out:    out,
	}
}

func (l *Logger) Close() error {
	if l.out != nil {
		return l.out.Close()
	}
	return nil
}

// Emerg logs a message with severity LOG_EMERG, ignoring the severity
// passed to New.
func (l *Logger) Emerg(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_EMERG, fmt.Sprintln(v...))
	return err
}

// Alert logs a message with severity LOG_ALERT, ignoring the severity
// passed to New.
func (l *Logger) Alert(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_ALERT, fmt.Sprintln(v...))
	return err
}

// Crit logs a message with severity LOG_CRIT, ignoring the severity
// passed to New.
func (l *Logger) Crit(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_CRIT, fmt.Sprintln(v...))
	return err
}

// Err logs a message with severity LOG_ERR, ignoring the severity
// passed to New.
func (l *Logger) Err(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_ERR, fmt.Sprintln(v...))
	return err
}

// Warning logs a message with severity LOG_WARNING, ignoring the
// severity passed to New.
func (l *Logger) Warning(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_WARNING, fmt.Sprintln(v...))
	return err
}

// Notice logs a message with severity LOG_NOTICE, ignoring the
// severity passed to New.
func (l *Logger) Notice(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_NOTICE, fmt.Sprintln(v...))
	return err
}

// Info logs a message with severity LOG_INFO, ignoring the severity
// passed to New.
func (l *Logger) Info(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_INFO, fmt.Sprintln(v...))
	return err
}

// Debug logs a message with severity LOG_DEBUG, ignoring the severity
// passed to New.
func (l *Logger) Debug(v ...interface{}) error {
	_, err := l.out.writeAndRetry(LOG_DEBUG, fmt.Sprintln(v...))
	return err
}
