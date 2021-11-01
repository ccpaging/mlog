// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package syslog

import (
	"fmt"
	stdlog "log"
	"strconv"
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
}

// NewLogger creates a log.Logger whose output is written to the
// system log service with the specified priority, a combination of
// the syslog facility and severity. The logFlag argument is the flag
// set passed through to log.New to create the Logger.
func NewLogger(p Priority, logFlag int) (*Logger, error) {
	s, err := New(p, "")
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: stdlog.New(s, "", logFlag|stdlog.Llevel),
	}, nil
}

func (l *Logger) pr(severity Priority) string {
	return "<" + strconv.Itoa(int(severity)) + ">"
}

// Emerg logs a message with severity LOG_EMERG, ignoring the severity
// passed to New.
func (l *Logger) Emerg(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_EMERG), fmt.Sprint(v...))
}

// Alert logs a message with severity LOG_ALERT, ignoring the severity
// passed to New.
func (l *Logger) Alert(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_ALERT), fmt.Sprint(v...))
}

// Crit logs a message with severity LOG_CRIT, ignoring the severity
// passed to New.
func (l *Logger) Crit(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_CRIT), fmt.Sprint(v...))
}

// Err logs a message with severity LOG_ERR, ignoring the severity
// passed to New.
func (l *Logger) Err(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_ERR), fmt.Sprint(v...))
}

// Warning logs a message with severity LOG_WARNING, ignoring the
// severity passed to New.
func (l *Logger) Warning(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_WARNING), fmt.Sprint(v...))
}

// Notice logs a message with severity LOG_NOTICE, ignoring the
// severity passed to New.
func (l *Logger) Notice(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_NOTICE), fmt.Sprint(v...))
}

// Info logs a message with severity LOG_INFO, ignoring the severity
// passed to New.
func (l *Logger) Info(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_INFO), fmt.Sprint(v...))
}

// Debug logs a message with severity LOG_DEBUG, ignoring the severity
// passed to New.
func (l *Logger) Debug(v ...interface{}) error {
	return l.OutputL(2, l.pr(LOG_DEBUG), fmt.Sprint(v...))
}
