// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package file

import (
	"os"

	"github.com/IBM/fp-go/v2/ioeither/file"
)

var (
	// CreateTemp creates a temporary file with proper parametrization.
	// It is an alias for ioeither.file.CreateTemp which wraps os.CreateTemp
	// in an IOResult context for functional composition.
	//
	// This function takes a directory and pattern parameter and returns an IOResult
	// that produces a temporary file handle when executed.
	//
	// Parameters:
	//   - dir: directory where the temporary file should be created (empty string uses default temp dir)
	//   - pattern: filename pattern with optional '*' placeholder for random suffix
	//
	// Returns:
	//   IOResult[*os.File] that when executed creates and returns a temporary file handle
	//
	// Example:
	//   tempFile := CreateTemp("", "myapp-*.tmp")
	//   result := tempFile()
	//   file, err := E.UnwrapError(result)
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   defer file.Close()
	CreateTemp = file.CreateTemp
)

// WithTempFile creates a temporary file, then invokes a callback to create a resource
// based on the file, then automatically closes and removes the temp file.
//
// This function provides safe temporary file management by:
//  1. Creating a temporary file with sensible defaults
//  2. Passing the file handle to the provided callback function
//  3. Ensuring the file is closed and removed, even if the callback fails
//
// Type Parameters:
//   - A: The type of result produced by the callback function
//
// Parameters:
//   - f: A Kleisli function that takes a *os.File and returns an IOResult[A]
//
// Returns:
//
//	IOResult[A] that when executed creates a temp file, runs the callback,
//	and cleans up the file regardless of success or failure
//
// Example - Writing and reading from a temporary file:
//
//	import (
//	    "io"
//	    "os"
//	    E "github.com/IBM/fp-go/v2/either"
//	    "github.com/IBM/fp-go/v2/ioresult"
//	    "github.com/IBM/fp-go/v2/ioresult/file"
//	)
//
//	// Write data to temp file and return the number of bytes written
//	writeToTemp := func(f *os.File) ioresult.IOResult[int] {
//	    return ioresult.TryCatchError(func() (int, error) {
//	        data := []byte("Hello, temporary world!")
//	        return f.Write(data)
//	    })
//	}
//
//	result := file.WithTempFile(writeToTemp)
//	bytesWritten, err := E.UnwrapError(result())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Wrote %d bytes to temporary file\n", bytesWritten)
//
// Example - Processing data through a temporary file:
//
//	processData := func(data []byte) ioresult.IOResult[string] {
//	    return file.WithTempFile(func(f *os.File) ioresult.IOResult[string] {
//	        return ioresult.TryCatchError(func() (string, error) {
//	            // Write data to temp file
//	            if _, err := f.Write(data); err != nil {
//	                return "", err
//	            }
//
//	            // Seek back to beginning
//	            if _, err := f.Seek(0, 0); err != nil {
//	                return "", err
//	            }
//
//	            // Read and process
//	            processed, err := io.ReadAll(f)
//	            if err != nil {
//	                return "", err
//	            }
//
//	            return strings.ToUpper(string(processed)), nil
//	        })
//	    })
//	}
//
//	result := processData([]byte("hello world"))
//	output, err := E.UnwrapError(result())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output) // "HELLO WORLD"
//
// The temporary file is guaranteed to be cleaned up even if the callback function
// panics or returns an error, providing safe resource management in a functional style.
//
//go:inline
func WithTempFile[A any](f Kleisli[*os.File, A]) IOResult[A] {
	return file.WithTempFile(f)
}
