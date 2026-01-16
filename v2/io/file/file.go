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

// Package file provides IO operations for file system interactions.
//
// This package offers functional wrappers around common file operations,
// returning IO monads that encapsulate side effects. All operations are
// lazy and only execute when the returned IO is invoked.
//
// # Core Operations
//
// The package provides two main operations:
//   - Close: Safely close io.Closer resources
//   - Remove: Remove files from the file system
//
// Both operations ignore errors and return the original input, making them
// suitable for cleanup operations where errors should not interrupt the flow.
//
// # Basic Usage
//
//	// Close a file
//	file, _ := os.Open("data.txt")
//	closeIO := file.Close(file)
//	closeIO() // Closes the file, ignoring any error
//
//	// Remove a file
//	removeIO := file.Remove("temp.txt")
//	removeIO() // Removes the file, ignoring any error
//
// # Composition with IO
//
// These operations can be composed with other IO operations:
//
//	result := pipe.Pipe2(
//	    openFile("data.txt"),
//	    io.ChainFirst(processFile),
//	    io.Chain(file.Close),
//	)
//
// # Error Handling
//
// Both Close and Remove intentionally ignore errors. This design is suitable
// for cleanup operations where:
//   - The operation is best-effort
//   - Errors should not interrupt the program flow
//   - The resource state is not critical
//
// For operations requiring error handling, use ioeither or ioresult instead.
package file

import (
	"io"
	"os"

	IO "github.com/IBM/fp-go/v2/io"
)

// Close closes a closeable resource and ignores any potential error.
// Returns an IO that, when executed, closes the resource and returns it.
//
// This function is useful for cleanup operations where errors can be safely
// ignored, such as in defer statements or resource cleanup chains.
//
// Type Parameters:
//   - R: Any type that implements io.Closer
//
// Parameters:
//   - r: The resource to close
//
// Returns:
//   - IO[R]: An IO computation that closes the resource and returns it
//
// Example:
//
//	file, _ := os.Open("data.txt")
//	defer file.Close(file)() // Close when function returns
//
// Example with IO composition:
//
//	result := pipe.Pipe3(
//	    openFile("data.txt"),
//	    io.Chain(readContent),
//	    io.ChainFirst(file.Close),
//	)
//
// Note: The #nosec comment is intentional - errors are deliberately ignored
// for cleanup operations where failure should not interrupt the flow.
func Close[R io.Closer](r R) IO.IO[R] {
	return func() R {
		r.Close() // #nosec: G104
		return r
	}
}

// Remove removes a file or directory and ignores any potential error.
// Returns an IO that, when executed, removes the named file or directory
// and returns the name.
//
// This function is useful for cleanup operations where errors can be safely
// ignored, such as removing temporary files or cache directories.
//
// Parameters:
//   - name: The path to the file or directory to remove
//
// Returns:
//   - IO[string]: An IO computation that removes the file and returns the name
//
// Example:
//
//	cleanup := file.Remove("temp.txt")
//	cleanup() // Removes temp.txt, ignoring any error
//
// Example with multiple files:
//
//	cleanup := pipe.Pipe2(
//	    file.Remove("temp1.txt"),
//	    io.ChainTo(file.Remove("temp2.txt")),
//	)
//	cleanup() // Removes both files
//
// Example in defer:
//
//	tempFile := "temp.txt"
//	defer file.Remove(tempFile)()
//	// ... use tempFile ...
//
// Note: The #nosec comment is intentional - errors are deliberately ignored
// for cleanup operations where failure should not interrupt the flow.
// This function only removes the named file or empty directory. To remove
// a directory and its contents, use os.RemoveAll wrapped in an IO.
func Remove(name string) IO.IO[string] {
	return func() string {
		os.Remove(name) // #nosec: G104
		return name
	}
}
