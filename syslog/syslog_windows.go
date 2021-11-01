// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslog

import (
	"errors"
)

// unixSyslog opens a connection to the syslog daemon running on the
// local machine using a Unix domain socket.

func unixSyslog() (conn serverConn, err error) {
	return nil, errors.New("Unix syslog delivery error")
}
