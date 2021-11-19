// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mlog

import (
	"bytes"
	"errors"
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
	// logger: example_test.go:20: INFO Hello, log file!
}

func ExampleLogger_Info() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "", log.Lshortfile)
	)

	logger.Info("Hello world")

	fmt.Print(&buf)
	// Output:
	// example_test.go:33: INFO Hello world
}

func ExampleLogger_Error() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "", log.Lshortfile)
	)

	logger.Error(errors.New("Error message."), " Hello world")

	fmt.Print(&buf)
	// Output:
	// example_test.go:46: EROR Error message. Hello world
}

func ExampleLogger_New() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "logger: ", log.Lshortfile).New("new: ")
	)

	logger.Print("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// new: example_test.go:59: INFO Hello, log file!
}
