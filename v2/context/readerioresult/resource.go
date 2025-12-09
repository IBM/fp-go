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

package readerioresult

import (
	"context"
	"io"

	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/result"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource.
// This implements the RAII (Resource Acquisition Is Initialization) pattern, ensuring that resources are
// properly released even if the operation fails or the context is canceled.
//
// The resource is created, used, and released in a safe manner:
//   - onCreate: Creates the resource
//   - The provided function uses the resource
//   - onRelease: Releases the resource (always called, even on error)
//
// Parameters:
//   - onCreate: ReaderIOResult that creates the resource
//   - onRelease: Function to release the resource
//
// Returns a function that takes a resource-using function and returns a ReaderIOResult.
//
// Example:
//
//	file := WithResource(
//	    openFile("data.txt"),
//	    func(f *os.File) ReaderIOResult[any] {
//	        return TryCatch(func(ctx context.Context) func() (any, error) {
//	            return func() (any, error) { return nil, f.Close() }
//	        })
//	    },
//	)
//	result := file(func(f *os.File) ReaderIOResult[string] {
//	    return TryCatch(func(ctx context.Context) func() (string, error) {
//	        return func() (string, error) {
//	            data, err := io.ReadAll(f)
//	            return string(data), err
//	        }
//	    })
//	})
func WithResource[A, R, ANY any](onCreate ReaderIOResult[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A] {
	return RIOR.WithResource[A](onCreate, onRelease)
}

// onClose is a helper function that creates a ReaderIOResult for closing an io.Closer resource.
// It safely calls the Close() method and handles any errors that may occur during closing.
//
// Type Parameters:
//   - A: Must implement io.Closer interface
//
// Parameters:
//   - a: The resource to close
//
// Returns:
//   - ReaderIOResult[any]: A computation that closes the resource and returns nil on success
//
// The function ignores the context parameter since closing operations typically don't need context.
// Any error from Close() is captured and returned as a Result error.
func onClose[A io.Closer](a A) ReaderIOResult[any] {
	return func(_ context.Context) IOResult[any] {
		return func() Result[any] {
			return result.TryCatchError[any](nil, a.Close())
		}
	}
}

// WithCloser creates a resource management function specifically for io.Closer resources.
// This is a specialized version of WithResource that automatically handles closing of resources
// that implement the io.Closer interface.
//
// The function ensures that:
//   - The resource is created using the onCreate function
//   - The resource is automatically closed when the operation completes (success or failure)
//   - Any errors during closing are properly handled
//   - The resource is closed even if the main operation fails or the context is canceled
//
// Type Parameters:
//   - B: The type of value returned by the resource-using function
//   - A: The type of resource that implements io.Closer
//
// Parameters:
//   - onCreate: ReaderIOResult that creates the io.Closer resource
//
// Returns:
//   - A function that takes a resource-using function and returns a ReaderIOResult[B]
//
// Example with file operations:
//
//	openFile := func(filename string) ReaderIOResult[*os.File] {
//	    return TryCatch(func(ctx context.Context) func() (*os.File, error) {
//	        return func() (*os.File, error) {
//	            return os.Open(filename)
//	        }
//	    })
//	}
//
//	fileReader := WithCloser(openFile("data.txt"))
//	result := fileReader(func(f *os.File) ReaderIOResult[string] {
//	    return TryCatch(func(ctx context.Context) func() (string, error) {
//	        return func() (string, error) {
//	            data, err := io.ReadAll(f)
//	            return string(data), err
//	        }
//	    })
//	})
//
// Example with HTTP response:
//
//	httpGet := func(url string) ReaderIOResult[*http.Response] {
//	    return TryCatch(func(ctx context.Context) func() (*http.Response, error) {
//	        return func() (*http.Response, error) {
//	            return http.Get(url)
//	        }
//	    })
//	}
//
//	responseReader := WithCloser(httpGet("https://api.example.com/data"))
//	result := responseReader(func(resp *http.Response) ReaderIOResult[[]byte] {
//	    return TryCatch(func(ctx context.Context) func() ([]byte, error) {
//	        return func() ([]byte, error) {
//	            return io.ReadAll(resp.Body)
//	        }
//	    })
//	})
//
// Example with database connection:
//
//	openDB := func(dsn string) ReaderIOResult[*sql.DB] {
//	    return TryCatch(func(ctx context.Context) func() (*sql.DB, error) {
//	        return func() (*sql.DB, error) {
//	            return sql.Open("postgres", dsn)
//	        }
//	    })
//	}
//
//	dbQuery := WithCloser(openDB("postgres://..."))
//	result := dbQuery(func(db *sql.DB) ReaderIOResult[[]User] {
//	    return TryCatch(func(ctx context.Context) func() ([]User, error) {
//	        return func() ([]User, error) {
//	            rows, err := db.QueryContext(ctx, "SELECT * FROM users")
//	            if err != nil {
//	                return nil, err
//	            }
//	            defer rows.Close()
//	            return scanUsers(rows)
//	        }
//	    })
//	})
func WithCloser[B any, A io.Closer](onCreate ReaderIOResult[A]) Kleisli[Kleisli[A, B], B] {
	return WithResource[B](onCreate, onClose[A])
}
