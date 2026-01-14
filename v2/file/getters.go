// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package file provides utility functions for working with file paths and I/O interfaces.
// It offers functional programming utilities for path manipulation and type conversions
// for common I/O interfaces.
package file

import (
	"io"
	"path/filepath"
)

// Join appends a filename to a root path using the operating system's path separator.
// Returns a curried function that takes a root path and joins it with the provided name.
//
// This function follows the "data last" principle, where the data (root path) is provided
// last, making it ideal for use in functional pipelines and partial application. The name
// parameter is fixed first, creating a reusable path builder function.
//
// This is useful for creating reusable path builders in functional pipelines.
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Data last: fix the filename first, apply root path later
//	addConfig := file.Join("config.json")
//	path := addConfig("/etc/myapp")
//	// path is "/etc/myapp/config.json" on Unix
//	// path is "\etc\myapp\config.json" on Windows
//
//	// Using with Pipe (data flows through the pipeline)
//	result := F.Pipe1("/var/log", file.Join("app.log"))
//	// result is "/var/log/app.log" on Unix
//
//	// Chain multiple joins
//	result := F.Pipe2(
//	    "/root",
//	    file.Join("subdir"),
//	    file.Join("file.txt"),
//	)
//	// result is "/root/subdir/file.txt"
func Join(name string) Endomorphism[string] {
	return func(root string) string {
		return filepath.Join(root, name)
	}
}

// ToReader converts any type that implements io.Reader to the io.Reader interface.
// This is useful for type erasure when you need to work with the interface type
// rather than a concrete implementation.
//
// Example:
//
//	import (
//	    "bytes"
//	    "io"
//	)
//
//	buf := bytes.NewBuffer([]byte("hello"))
//	var reader io.Reader = file.ToReader(buf)
//	// reader is now of type io.Reader
func ToReader[R io.Reader](r R) io.Reader {
	return r
}

// ToWriter converts any type that implements io.Writer to the io.Writer interface.
// This is useful for type erasure when you need to work with the interface type
// rather than a concrete implementation.
//
// Example:
//
//	import (
//	    "bytes"
//	    "io"
//	)
//
//	buf := &bytes.Buffer{}
//	var writer io.Writer = file.ToWriter(buf)
//	// writer is now of type io.Writer
func ToWriter[W io.Writer](w W) io.Writer {
	return w
}

// ToCloser converts any type that implements io.Closer to the io.Closer interface.
// This is useful for type erasure when you need to work with the interface type
// rather than a concrete implementation.
//
// Example:
//
//	import (
//	    "os"
//	    "io"
//	)
//
//	f, _ := os.Open("file.txt")
//	var closer io.Closer = file.ToCloser(f)
//	defer closer.Close()
//	// closer is now of type io.Closer
func ToCloser[C io.Closer](c C) io.Closer {
	return c
}
