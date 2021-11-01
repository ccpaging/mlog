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
		logger = New(&buf, "logger: ", log.Lshortfile)
	)

	logger.Print("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// logger: example_test.go:19: INFO  Hello, log file!
}

func ExampleLogger_Output() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "INFO: ", log.Lshortfile)

		info = func(info string) {
			logger.Output(2, info)
		}
	)

	info("Hello world")

	fmt.Print(&buf)
	// Output:
	// INFO: example_test.go:36: Hello world
}
