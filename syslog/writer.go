// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package syslog

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// A Writer is a connection to a syslog server.
type Writer struct {
	priority Priority
	tag      string
	hostname string
	network  string
	raddr    string

	mu   sync.Mutex // guards conn
	conn serverConn
}

// This interface and the separate syslog_unix.go file exist for
// Solaris support as implemented by gccgo. On Solaris you cannot
// simply open a TCP connection to the syslog daemon. The gccgo
// sources have a syslog_solaris.go file that implements unixSyslog to
// return a type that satisfies this interface and simply calls the C
// library syslog function.
type serverConn interface {
	writeString(p Priority, hostname, tag, s, nl string) error
	close() error
}

type netConn struct {
	local bool
	conn  net.Conn
}

// New establishes a new connection to the system log daemon. Each
// write to the returned writer sends a log message with the given
// priority (a combination of the syslog facility and severity) and
// prefix tag. If tag is empty, the os.Args[0] is used.
func New(priority Priority, tag string) (*Writer, error) {
	return Dial("", "", priority, tag)
}

// Dial establishes a connection to a log daemon by connecting to
// address raddr on the specified network. Each write to the returned
// writer sends a log message with the facility and severity
// (from priority) and tag. If tag is empty, the os.Args[0] is used.
// If network is empty, Dial will connect to the local syslog server.
// Otherwise, see the documentation for net.Dial for valid values
// of network and raddr.
func Dial(network, raddr string, priority Priority, tag string) (*Writer, error) {
	if priority < 0 || priority > LOG_LOCAL7|LOG_DEBUG {
		return nil, errors.New("log/syslog: invalid priority")
	}

	if tag == "" {
		tag = os.Args[0]
	}
	hostname, _ := os.Hostname()

	w := &Writer{
		priority: priority,
		tag:      tag,
		hostname: hostname,
		network:  network,
		raddr:    raddr,
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.connect()
	if err != nil {
		return nil, err
	}
	return w, err
}

// connect makes a connection to the syslog server.
// It must be called with w.mu held.
func (w *Writer) connect() (err error) {
	if w.conn != nil {
		// ignore err from close, it makes sense to continue anyway
		w.conn.close()
		w.conn = nil
	}

	if w.network == "" {
		w.conn, err = unixSyslog()
		if w.hostname == "" {
			w.hostname = "localhost"
		}
	} else {
		var c net.Conn
		c, err = net.Dial(w.network, w.raddr)
		if err == nil {
			w.conn = &netConn{
				conn:  c,
				local: w.network == "unixgram" || w.network == "unix",
			}
			if w.hostname == "" {
				w.hostname = c.LocalAddr().String()
			}
		}
	}
	return
}

// Write sends a log message to the syslog daemon.
func (w *Writer) Write(b []byte) (int, error) {
	severity := w.priority
	if len(b) >= 3 && b[0] == '<' && b[2] == '>' {
		switch b[1] {
		case '0':
			severity = LOG_EMERG
		case '1':
			severity = LOG_ALERT
		case '2':
			severity = LOG_CRIT
		case '3':
			severity = LOG_ERR
		case '4':
			severity = LOG_WARNING
		case '5':
			severity = LOG_NOTICE
		case '6':
			severity = LOG_INFO
		case '7':
			severity = LOG_DEBUG
		}
		return w.writeAndRetry(severity, string(b[3:]))
	}
	return w.writeAndRetry(severity, string(b))
}

// Close closes a connection to the syslog daemon.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		err := w.conn.close()
		w.conn = nil
		return err
	}
	return nil
}

// Emerg logs a message with severity LOG_EMERG, ignoring the severity
// passed to New.
func (w *Writer) Emerg(m string) error {
	_, err := w.writeAndRetry(LOG_EMERG, m)
	return err
}

// Alert logs a message with severity LOG_ALERT, ignoring the severity
// passed to New.
func (w *Writer) Alert(m string) error {
	_, err := w.writeAndRetry(LOG_ALERT, m)
	return err
}

// Crit logs a message with severity LOG_CRIT, ignoring the severity
// passed to New.
func (w *Writer) Crit(m string) error {
	_, err := w.writeAndRetry(LOG_CRIT, m)
	return err
}

// Err logs a message with severity LOG_ERR, ignoring the severity
// passed to New.
func (w *Writer) Err(m string) error {
	_, err := w.writeAndRetry(LOG_ERR, m)
	return err
}

// Warning logs a message with severity LOG_WARNING, ignoring the
// severity passed to New.
func (w *Writer) Warning(m string) error {
	_, err := w.writeAndRetry(LOG_WARNING, m)
	return err
}

// Notice logs a message with severity LOG_NOTICE, ignoring the
// severity passed to New.
func (w *Writer) Notice(m string) error {
	_, err := w.writeAndRetry(LOG_NOTICE, m)
	return err
}

// Info logs a message with severity LOG_INFO, ignoring the severity
// passed to New.
func (w *Writer) Info(m string) error {
	_, err := w.writeAndRetry(LOG_INFO, m)
	return err
}

// Debug logs a message with severity LOG_DEBUG, ignoring the severity
// passed to New.
func (w *Writer) Debug(m string) error {
	_, err := w.writeAndRetry(LOG_DEBUG, m)
	return err
}

func (w *Writer) writeAndRetry(p Priority, s string) (int, error) {
	pr := (w.priority & facilityMask) | (p & severityMask)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		if n, err := w.write(pr, s); err == nil {
			return n, err
		}
	}
	if err := w.connect(); err != nil {
		return 0, err
	}
	return w.write(pr, s)
}

// write generates and writes a syslog formatted string. The
// format is as follows: <PRI>TIMESTAMP HOSTNAME TAG[PID]: MSG
func (w *Writer) write(p Priority, msg string) (int, error) {
	// ensure it ends in a \n
	nl := ""
	if !strings.HasSuffix(msg, "\n") {
		nl = "\n"
	}

	err := w.conn.writeString(p, w.hostname, w.tag, msg, nl)
	if err != nil {
		return 0, err
	}
	// Note: return the length of the input, not the number of
	// bytes printed by Fprintf, because this must behave like
	// an io.Writer.
	return len(msg), nil
}

func (n *netConn) writeString(p Priority, hostname, tag, msg, nl string) error {
	if n.local {
		// Compared to the network form below, the changes are:
		//	1. Use time.Stamp instead of time.RFC3339.
		//	2. Drop the hostname field from the Fprintf.
		timestamp := time.Now().Format(time.Stamp)
		_, err := fmt.Fprintf(n.conn, "<%d>%s %s[%d]: %s%s",
			p, timestamp,
			tag, os.Getpid(), msg, nl)
		return err
	}
	timestamp := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(n.conn, "<%d>%s %s %s[%d]: %s%s",
		p, timestamp, hostname,
		tag, os.Getpid(), msg, nl)
	return err
}

func (n *netConn) close() error {
	return n.conn.Close()
}
