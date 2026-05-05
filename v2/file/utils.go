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

import "io"

type nopWriteCloser struct {
	delegate io.Writer
}

func (_ *nopWriteCloser) Close() error {
	return nil
}

func (nc *nopWriteCloser) Write(p []byte) (n int, err error) {
	return nc.delegate.Write(p)
}

// NopReadCloser wraps an io.Reader with a no-op Close method, converting it to an io.ReadCloser.
// This is useful when you need to satisfy an io.ReadCloser interface but don't need actual
// close functionality, such as when working with in-memory buffers or strings.
//
// This function is a generic wrapper around io.NopCloser that preserves type information
// for better type inference in functional pipelines.
//
// Type Parameters:
//   - R: Any type that implements io.Reader
//
// Parameters:
//   - r: The io.Reader to wrap
//
// Returns:
//   - io.ReadCloser: A ReadCloser that delegates Read calls to r and has a no-op Close
//
// Example:
//
//   reader := strings.NewReader("hello world")
//   readCloser := NopReadCloser(reader)
//   defer readCloser.Close() // no-op, safe to call
//   data, _ := io.ReadAll(readCloser)
//
// See Also:
//   - NopWriteCloser: Similar wrapper for io.Writer
func NopReadCloser[R io.Reader](r R) io.ReadCloser {
	return io.NopCloser(r)
}

// NopWriteCloser wraps an io.Writer with a no-op Close method, converting it to an io.WriteCloser.
// This is useful when you need to satisfy an io.WriteCloser interface but don't need actual
// close functionality, such as when writing to in-memory buffers or when the underlying
// writer doesn't require cleanup.
//
// The returned WriteCloser delegates all Write calls to the underlying writer and returns
// nil from Close without performing any cleanup operations.
//
// Type Parameters:
//   - W: Any type that implements io.Writer
//
// Parameters:
//   - w: The io.Writer to wrap
//
// Returns:
//   - io.WriteCloser: A WriteCloser that delegates Write calls to w and has a no-op Close
//
// Example:
//
//   var buf bytes.Buffer
//   writeCloser := NopWriteCloser(&buf)
//   writeCloser.Write([]byte("hello"))
//   writeCloser.Close() // no-op, safe to call
//   fmt.Println(buf.String()) // prints: hello
//
// See Also:
//   - NopReadCloser: Similar wrapper for io.Reader
func NopWriteCloser[W io.Writer](w W) io.WriteCloser {
	return &nopWriteCloser{w}
}
