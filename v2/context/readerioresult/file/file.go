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
	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/file"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
)

// STDIO is a special constant representing standard input/output streams.
// When used as a filename with ReadFile or WriteFile, it causes the operation
// to use os.Stdin or os.Stdout respectively, instead of opening a file.
//
// This convention is commonly used in Unix command-line tools to allow
// reading from stdin or writing to stdout by specifying "-" as the filename.
//
// Example:
//
//	// Read from stdin
//	data := ReadFile(STDIO)(ctx)()
//
//	// Write to stdout
//	result := WriteFile([]byte("Hello"))(STDIO)(ctx)()
const (
	STDIO = "-"
)

var (
	isNotStdIO = O.FromPredicate(P.Not(P.IsStrictEqual[string]()(STDIO)))

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

	// Create creates or truncates a file for writing within the given context.
	// If the file already exists, it is truncated. If it doesn't exist, it is created
	// with mode 0666 (before umask).
	//
	// The operation respects context cancellation and returns a ReaderIOResult
	// that produces an os.File handle on success.
	//
	// The returned file handle should be closed using the Close function when no longer needed,
	// or managed automatically using WithResource or WriteFile.
	//
	// Parameters:
	//   - path: The path to the file to create or truncate
	//
	// Returns:
	//   - ReaderIOResult[*os.File]: A context-aware computation that creates the file
	//
	// Example:
	//
	//	createFile := Create("output.txt")
	//	result := createFile(ctx)()
	//	either.Fold(
	//	    result,
	//	    func(err error) { log.Printf("Error: %v", err) },
	//	    func(f *os.File) {
	//	        defer f.Close()
	//	        f.WriteString("Hello, World!")
	//	    },
	//	)
	//
	// See Also:
	//   - WriteFile: For writing data to a file with automatic resource management
	//   - Open: For opening files for reading
	//   - Close: For closing file handles
	Create = F.Flow3(
		IOEF.Create,
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
// If the filename is "-", the data is read from os.Stdin instead.
//
// The operation:
//   - Opens the file for reading (or uses stdin if filename is "-")
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
	return RIOE.WithResource[[]byte](OpenOrStdIn()(path), Close[io.ReadCloser])(func(r io.ReadCloser) ReaderIOResult[[]byte] {
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return func() ET.Either[error, []byte] {
				return file.ReadAll(ctx, r)
			}
		}
	})
}

// WriteFile writes data to a file in a context-aware manner.
// This function automatically manages the file resource using the RAII pattern,
// ensuring the file is properly closed even if an error occurs or the context is canceled.
//
// If the file doesn't exist, it is created with mode 0666 (before umask).
// If the file already exists, it is truncated before writing.
// If the filename is "-", the data is written to os.Stdout instead.
//
// The operation:
//   - Creates or truncates the file for writing (or uses stdout if filename is "-")
//   - Writes all data to the file
//   - Automatically closes the file when done
//   - Respects context cancellation during the write operation
//
// Parameters:
//   - data: The byte slice to write to the file
//
// Returns:
//   - Kleisli[string, []byte]: A function that takes a file path and returns a computation
//     that writes the data and returns the written bytes on success
//
// Example:
//
//	writeOp := WriteFile([]byte("Hello, World!"))
//	result := writeOp("output.txt")(ctx)()
//	either.Fold(
//	    result,
//	    func(err error) { log.Printf("Write error: %v", err) },
//	    func(data []byte) { log.Printf("Wrote %d bytes", len(data)) },
//	)
//
// The function uses WithResource internally to ensure proper cleanup:
//
//	WriteFile(data) = Create >> WriteAll(data) >> Close
//
// See Also:
//   - ReadFile: For reading file contents with automatic resource management
//   - Create: For creating files without automatic writing
//   - WriteAll: For writing to an already-open file handle
func WriteFile(data []byte) Kleisli[string, []byte] {
	return F.Flow2(
		CreateOrStdOut(),
		WriteAll[io.WriteCloser](data),
	)
}

type noCloseable struct {
	delegate *os.File
}

func (_ *noCloseable) Close() error {
	return nil
}

func (nc *noCloseable) Write(p []byte) (n int, err error) {
	return nc.delegate.Write(p)
}

func (nc *noCloseable) Read(p []byte) (n int, err error) {
	return nc.delegate.Read(p)
}

// CreateOrStdOut creates a file for writing or returns stdout if the path is "-".
// This function is useful for CLI applications that need to support writing to either
// a file or stdout based on user input.
//
// The function uses a special convention where the path "-" (STDIO constant) represents
// stdout. For any other path, it creates or truncates the file normally.
//
// The returned io.WriteCloser is safe to close in both cases:
//   - For regular files, Close closes the file handle
//   - For stdout, Close is a no-op (does not close stdout)
//
// Parameters:
//   - path: The file path to create, or "-" for stdout
//
// Returns:
//   - Kleisli[string, io.WriteCloser]: A function that takes a path and returns
//     a computation producing a WriteCloser
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Write to a file
//	writeOp := F.Pipe1("output.txt", CreateOrStdOut())
//	result := writeOp(ctx)()
//	either.Fold(
//	    result,
//	    func(err error) { log.Printf("Error: %v", err) },
//	    func(w io.WriteCloser) {
//	        defer w.Close()
//	        w.Write([]byte("Hello, World!"))
//	    },
//	)
//
//	// Write to stdout
//	writeOp := F.Pipe1("-", CreateOrStdOut())
//	result := writeOp(ctx)()
//	// Writes to stdout, Close is safe but does nothing
//
// See Also:
//   - OpenOrStdIn: For reading from files or stdin
//   - Create: For creating files without stdout fallback
func CreateOrStdOut() Kleisli[string, io.WriteCloser] {
	return F.Flow3(
		isNotStdIO,
		O.Map(F.Flow2(
			Create,
			RIOE.Map[*os.File](FL.ToWriteCloser),
		)),
		O.GetOrElse(lazy.Of(RIOE.Of[io.WriteCloser](&noCloseable{os.Stdout}))),
	)
}

// OpenOrStdIn opens a file for reading or returns stdin if the path is "-".
// This function is useful for CLI applications that need to support reading from either
// a file or stdin based on user input.
//
// The function uses a special convention where the path "-" (STDIO constant) represents
// stdin. For any other path, it opens the file normally.
//
// The returned io.ReadCloser is safe to close in both cases:
//   - For regular files, Close closes the file handle
//   - For stdin, Close is a no-op (does not close stdin)
//
// Parameters:
//   - path: The file path to open, or "-" for stdin
//
// Returns:
//   - Kleisli[string, io.ReadCloser]: A function that takes a path and returns
//     a computation producing a ReadCloser
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Read from a file
//	readOp := F.Pipe1("input.txt", OpenOrStdIn())
//	result := readOp(ctx)()
//	either.Fold(
//	    result,
//	    func(err error) { log.Printf("Error: %v", err) },
//	    func(r io.ReadCloser) {
//	        defer r.Close()
//	        data, _ := io.ReadAll(r)
//	        fmt.Println(string(data))
//	    },
//	)
//
//	// Read from stdin
//	readOp := F.Pipe1("-", OpenOrStdIn())
//	result := readOp(ctx)()
//	// Reads from stdin, Close is safe but does nothing
//
// See Also:
//   - CreateOrStdOut: For writing to files or stdout
//   - Open: For opening files without stdin fallback
func OpenOrStdIn() Kleisli[string, io.ReadCloser] {
	return F.Flow3(
		isNotStdIO,
		O.Map(F.Flow2(
			Open,
			RIOE.Map[*os.File](FL.ToReadCloser),
		)),
		O.GetOrElse(lazy.Of(RIOE.Of[io.ReadCloser](&noCloseable{os.Stdin}))),
	)
}
