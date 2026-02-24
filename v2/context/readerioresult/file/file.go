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

// Package file provides context-aware file operations that integrate with the ReaderIOResult monad.
// It offers safe, composable file I/O operations that respect context cancellation and properly
// manage resources using the RAII pattern.
//
// All operations in this package:
//   - Respect context.Context for cancellation and timeouts
//   - Return ReaderIOResult for composable error handling
//   - Automatically manage resource cleanup
//   - Are safe to use in concurrent environments
//
// # Example Usage
//
//	// Read a file with automatic resource management
//	readOp := ReadFile("data.txt")
//	result := readOp(ctx)()
//
//	// Open and manually manage a file
//	fileOp := Open("config.json")
//	fileResult := fileOp(ctx)()
package file

import (
	"context"
	"io"
	"os"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/file"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
)

var (
	// Open opens a file for reading within the given context.
	// The operation respects context cancellation and returns a ReaderIOResult
	// that produces an os.File handle on success.
	//
	// The returned file handle should be closed using the Close function when no longer needed,
	// or managed automatically using WithResource or ReadFile.
	//
	// Parameters:
	//   - path: The path to the file to open
	//
	// Returns:
	//   - ReaderIOResult[*os.File]: A context-aware computation that opens the file
	//
	// Example:
	//
	//	openFile := Open("data.txt")
	//	result := openFile(ctx)()
	//	either.Fold(
	//	    result,
	//	    func(err error) { log.Printf("Error: %v", err) },
	//	    func(f *os.File) {
	//	        defer f.Close()
	//	        // Use file...
	//	    },
	//	)
	//
	// See Also:
	//   - ReadFile: For reading entire file contents with automatic resource management
	//   - Close: For closing file handles
	Open = F.Flow3(
		IOEF.Open,
		RIOE.FromIOEither[*os.File],
		RIOE.WithContext[*os.File],
	)

	// Remove removes a file by name.
	// The operation returns the filename on success, allowing for easy composition
	// with other file operations.
	//
	// Parameters:
	//   - name: The path to the file to remove
	//
	// Returns:
	//   - ReaderIOResult[string]: A computation that removes the file and returns its name
	//
	// Example:
	//
	//	removeOp := Remove("temp.txt")
	//	result := removeOp(ctx)()
	//	either.Fold(
	//	    result,
	//	    func(err error) { log.Printf("Failed to remove: %v", err) },
	//	    func(name string) { log.Printf("Removed: %s", name) },
	//	)
	//
	// See Also:
	//   - Open: For opening files
	//   - ReadFile: For reading file contents
	Remove = F.Flow2(
		IOEF.Remove,
		RIOE.FromIOEither[string],
	)
)

// Close closes an io.Closer resource and returns a ReaderIOResult.
// This function is generic and works with any type that implements io.Closer,
// including os.File, network connections, and other closeable resources.
//
// The function captures any error that occurs during closing and returns it
// as part of the ReaderIOResult. On success, it returns Void (empty struct).
//
// Type Parameters:
//   - C: Any type that implements io.Closer
//
// Parameters:
//   - c: The resource to close
//
// Returns:
//   - ReaderIOResult[Void]: A computation that closes the resource
//
// Example:
//
//	file, _ := os.Open("data.txt")
//	closeOp := Close(file)
//	result := closeOp(ctx)()
//
// Note: This function is typically used with WithResource for automatic resource management
// rather than being called directly.
//
// See Also:
//   - Open: For opening files
//   - ReadFile: For reading files with automatic closing
func Close[C io.Closer](c C) ReaderIOResult[Void] {
	return F.Pipe2(
		c,
		IOEF.Close[C],
		RIOE.FromIOEither[Void],
	)
}

// ReadFile reads the entire contents of a file in a context-aware manner.
// This function automatically manages the file resource using the RAII pattern,
// ensuring the file is properly closed even if an error occurs or the context is canceled.
//
// The operation:
//   - Opens the file for reading
//   - Reads all contents into a byte slice
//   - Automatically closes the file when done
//   - Respects context cancellation during the read operation
//
// Parameters:
//   - path: The path to the file to read
//
// Returns:
//   - ReaderIOResult[[]byte]: A computation that reads the file contents
//
// Example:
//
//	readOp := ReadFile("config.json")
//	result := readOp(ctx)()
//	either.Fold(
//	    result,
//	    func(err error) { log.Printf("Read error: %v", err) },
//	    func(data []byte) { log.Printf("Read %d bytes", len(data)) },
//	)
//
// The function uses WithResource internally to ensure proper cleanup:
//
//	ReadFile(path) = WithResource(Open(path), Close)(readAllBytes)
//
// See Also:
//   - Open: For opening files without automatic reading
//   - Close: For closing file handles
//   - WithResource: For custom resource management patterns
func ReadFile(path string) ReaderIOResult[[]byte] {
	return RIOE.WithResource[[]byte](Open(path), Close[*os.File])(func(r *os.File) ReaderIOResult[[]byte] {
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return func() ET.Either[error, []byte] {
				return file.ReadAll(ctx, r)
			}
		}
	})
}
