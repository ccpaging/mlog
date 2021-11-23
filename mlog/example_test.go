// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

import (
	"bytes"
	"fmt"
	"log"
)

func ExampleLogger() {
	var (
		buf    bytes.Buffer
		logger = NewLogger(log.New(&buf, "main: ", log.Lshortfile|log.Lmsgprefix), nil)
	)

	logger.Info("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// example_test.go:19: main: INFO Hello, log file!
}

func ExampleLogger_Debug() {
	var (
		buf    bytes.Buffer
		logger = NewLogger(log.New(&buf, "main: ", log.Lshortfile), nil)

		debugln = func(info string) {
			logger.Debug(info)
		}
	)

	debugln("Hello world")

	fmt.Print(&buf)
	// Output:
	// main: example_test.go:32: DEBG Hello world
}

func ExampleLogger_Info() {
	var (
		buf    bytes.Buffer
		logger = NewLogger(log.New(&buf, "main: ", log.Lmsgprefix), nil)
	)

	logger.Info("Hello world")

	fmt.Print(&buf)
	// Output:
	// main: INFO Hello world
}
