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
		logger = New("logger: ").SetOutput(log.New(&buf, "logger: ", log.Lshortfile))
	)

	logger.Info("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// INFO logger: example_test.go:19: Hello, log file!
}

func ExampleLogger_Debug() {
	var (
		buf    bytes.Buffer
		logger = New("main: ").SetOutput(log.New(&buf, "", log.Lshortfile))

		debugln = func(info string) {
			logger.Debug(info)
		}
	)

	debugln("Hello world")

	fmt.Print(&buf)
	// Output:
	// DEBG main: example_test.go:32: Hello world
}

func ExampleLogger_Info() {
	var (
		buf    bytes.Buffer
		logger = New("main: ").SetOutput(log.New(&buf, "", log.Lmsgprefix))
	)

	logger.Info("Hello world")

	fmt.Print(&buf)
	// Output:
	// INFO main: Hello world
}
