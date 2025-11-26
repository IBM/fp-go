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

	"github.com/IBM/fp-go/v2/ioeither"
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
//   - acquire: An IOEither that produces the ReadCloser resource
//
// Returns:
//
//	A Kleisli function that takes a reader function (which transforms RD to R)
//	and returns an IOEither that produces the read result R or an error.
//
// Example:
//
//	import (
//	    "os"
//	    "io"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/ioeither"
//	    "github.com/IBM/fp-go/v2/ioeither/file"
//	)
//
//	// Read first 10 bytes from a file
//	readFirst10 := func(f *os.File) ioeither.IOEither[error, []byte] {
//	    return ioeither.TryCatchError(func() ([]byte, error) {
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
//	data, err := result()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Read: %s\n", data)
//
// The Read function ensures that the file is closed even if the reading operation fails,
// providing safe and composable resource management in a functional style.
func Read[R any, RD io.ReadCloser](acquire IOEither[error, RD]) Kleisli[error, Kleisli[error, RD, R], R] {
	return ioeither.WithResource[R](
		acquire,
		Close[RD])
}
