// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

import (
	"bytes"
	"fmt"
	"log"
)

func ExampleLogWrapper() {
	var (
		buf    bytes.Buffer
		logger = Wrap(log.New(&buf, "logger: ", log.Lshortfile))
	)

	logger.Print("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// INFO logger: example_test.go:19: Hello, log file!
}

func ExampleLogWrapper_Debug() {
	var (
		buf    bytes.Buffer
		logger = Wrap(log.New(&buf, "INFO: ", log.Lshortfile))

		debugln = func(info string) {
			logger.Debug(info)
		}
	)

	debugln("Hello world")

	fmt.Print(&buf)
	// Output:
	// DEBG INFO: example_test.go:32: Hello world
}
