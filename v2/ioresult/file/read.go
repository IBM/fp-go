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

package file

import (
	"io"

	"github.com/IBM/fp-go/v2/ioeither/file"
)

// Read uses a generator function to create a stream, reads data from it using a provided
// reader function, and ensures the stream is properly closed after reading.
//
// This function provides safe resource management for reading operations by:
//  1. Acquiring a ReadCloser resource using the provided acquire function
//  2. Applying a reader function to extract data from the resource
//  3. Ensuring the resource is closed, even if an error occurs during reading
//
// Type Parameters:
//   - R: The type of data to be read from the stream
//   - RD: The type of the ReadCloser resource (must implement io.ReadCloser)
//
// Parameters:
//   - acquire: An IOResult that produces the ReadCloser resource
//
// Returns:
//
//	A Kleisli function that takes a reader function (which transforms RD to R)
//	and returns an IOResult that produces the read result R or an error.
//
// The key difference from ioeither.Read is that this returns IOResult[R] which is
// IO[Result[R]], representing a computation that returns a Result type (tuple of value and error)
// rather than an Either type.
//
// Example - Reading first N bytes from a file:
//
//	import (
//	    "os"
//	    "io"
//	    F "github.com/IBM/fp-go/v2/function"
//	    R "github.com/IBM/fp-go/v2/result"
//	    "github.com/IBM/fp-go/v2/ioresult"
//	    "github.com/IBM/fp-go/v2/ioresult/file"
//	)
//
//	// Read first 10 bytes from a file
//	readFirst10 := func(f *os.File) ioresult.IOResult[[]byte] {
//	    return ioresult.TryCatch(func() ([]byte, error) {
//	        buf := make([]byte, 10)
//	        n, err := f.Read(buf)
//	        return buf[:n], err
//	    })
//	}
//
//	result := F.Pipe1(
//	    file.Open("data.txt"),
//	    file.Read[[]byte, *os.File],
//	)(readFirst10)
//
//	// Execute the IO operation to get the Result
//	res := result()
//	data, err := res()  // Result is a tuple function
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Read: %s\n", data)
//
// Example - Using with Result combinators:
//
//	result := F.Pipe1(
//	    file.Open("config.json"),
//	    file.Read[[]byte, *os.File],
//	)(readFirst10)
//
//	// Chain operations using Result combinators
//	processed := F.Pipe2(
//	    result,
//	    ioresult.Map(func(data []byte) string {
//	        return string(data)
//	    }),
//	    ioresult.ChainFirst(func(s string) ioresult.IOResult[any] {
//	        return ioresult.Of[any](fmt.Printf("Read: %s\n", s))
//	    }),
//	)
//
//	res := processed()
//	str, err := res()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The Read function ensures that the file is closed even if the reading operation fails,
// providing safe and composable resource management in a functional style.
//
//go:inline
func Read[R any, RD io.ReadCloser](acquire IOResult[RD]) Kleisli[Kleisli[RD, R], R] {
	return file.Read[R](acquire)
}
