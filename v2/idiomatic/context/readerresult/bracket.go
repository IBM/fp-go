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

package readerresult

import (
	"context"
	"io"

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
)

// Bracket ensures safe resource management with guaranteed cleanup in the ReaderResult monad.
//
// This function implements the bracket pattern (also known as try-with-resources or RAII)
// for ReaderResult computations. It guarantees that the release action is called regardless
// of whether the use action succeeds or fails, making it ideal for managing resources like
// file handles, database connections, network sockets, or locks.
//
// The execution flow is:
//  1. Acquire the resource (lazily evaluated)
//  2. Use the resource with the provided function
//  3. Release the resource with access to: the resource, the result (if successful), and any error
//
// The release function is always called, even if:
//   - The acquire action fails (release is not called in this case)
//   - The use action fails (release receives the error)
//   - The use action succeeds (release receives nil error)
//
// Type Parameters:
//   - A: The type of the acquired resource
//   - B: The type of the result produced by using the resource
//   - ANY: The type returned by the release action (typically ignored)
//
// Parameters:
//   - acquire: Lazy computation that acquires the resource
//   - use: Function that uses the resource to produce a result
//   - release: Function that releases the resource, receiving the resource, result, and any error
//
// Returns:
//   - A ReaderResult[B] that safely manages the resource lifecycle
//
// Example - File handling:
//
//	import (
//	    "context"
//	    "os"
//	)
//
//	readFile := readerresult.Bracket(
//	    // Acquire: Open file
//	    func() readerresult.ReaderResult[*os.File] {
//	        return func(ctx context.Context) (*os.File, error) {
//	            return os.Open("data.txt")
//	        }
//	    },
//	    // Use: Read file contents
//	    func(file *os.File) readerresult.ReaderResult[string] {
//	        return func(ctx context.Context) (string, error) {
//	            data, err := io.ReadAll(file)
//	            return string(data), err
//	        }
//	    },
//	    // Release: Close file (always called)
//	    func(file *os.File, content string, err error) readerresult.ReaderResult[any] {
//	        return func(ctx context.Context) (any, error) {
//	            return nil, file.Close()
//	        }
//	    },
//	)
//
//	content, err := readFile(context.Background())
//
// Example - Database connection:
//
//	queryDB := readerresult.Bracket(
//	    // Acquire: Open connection
//	    func() readerresult.ReaderResult[*sql.DB] {
//	        return func(ctx context.Context) (*sql.DB, error) {
//	            return sql.Open("postgres", connString)
//	        }
//	    },
//	    // Use: Execute query
//	    func(db *sql.DB) readerresult.ReaderResult[[]User] {
//	        return func(ctx context.Context) ([]User, error) {
//	            return queryUsers(ctx, db)
//	        }
//	    },
//	    // Release: Close connection (always called)
//	    func(db *sql.DB, users []User, err error) readerresult.ReaderResult[any] {
//	        return func(ctx context.Context) (any, error) {
//	            return nil, db.Close()
//	        }
//	    },
//	)
//
// Example - Lock management:
//
//	withLock := readerresult.Bracket(
//	    // Acquire: Lock mutex
//	    func() readerresult.ReaderResult[*sync.Mutex] {
//	        return func(ctx context.Context) (*sync.Mutex, error) {
//	            mu.Lock()
//	            return mu, nil
//	        }
//	    },
//	    // Use: Perform critical section work
//	    func(mu *sync.Mutex) readerresult.ReaderResult[int] {
//	        return func(ctx context.Context) (int, error) {
//	            return performCriticalWork(ctx)
//	        }
//	    },
//	    // Release: Unlock mutex (always called)
//	    func(mu *sync.Mutex, result int, err error) readerresult.ReaderResult[any] {
//	        return func(ctx context.Context) (any, error) {
//	            mu.Unlock()
//	            return nil, nil
//	        }
//	    },
//	)
//
//go:inline
func Bracket[
	A, B, ANY any](

	acquire Lazy[ReaderResult[A]],
	use Kleisli[A, B],
	release func(A, B, error) ReaderResult[ANY],
) ReaderResult[B] {
	return RR.Bracket(acquire, WithContextK(use), release)
}

// WithResource creates a higher-order function for resource management with automatic cleanup.
//
// This function provides a more composable alternative to Bracket by creating a function
// that takes a resource-using function and automatically handles resource acquisition and
// release. This is particularly useful when you want to reuse the same resource management
// pattern with different operations.
//
// The pattern is:
//  1. Create a resource manager with onCreate and onRelease
//  2. Apply it to different use functions as needed
//  3. Each application ensures proper resource cleanup
//
// This is useful for:
//   - Creating reusable resource management patterns
//   - Building resource pools or factories
//   - Composing resource-dependent operations
//   - Abstracting resource lifecycle management
//
// Type Parameters:
//   - B: The type of the result produced by using the resource
//   - A: The type of the acquired resource
//   - ANY: The type returned by the release action (typically ignored)
//
// Parameters:
//   - onCreate: Lazy computation that creates/acquires the resource
//   - onRelease: Function that releases the resource (receives the resource)
//
// Returns:
//   - A Kleisli arrow that takes a resource-using function and returns a ReaderResult[B]
//     with automatic resource management
//
// Example - Reusable database connection manager:
//
//	import (
//	    "context"
//	    "database/sql"
//	)
//
//	// Create a reusable DB connection manager
//	withDB := readerresult.WithResource(
//	    // onCreate: Acquire connection
//	    func() readerresult.ReaderResult[*sql.DB] {
//	        return func(ctx context.Context) (*sql.DB, error) {
//	            return sql.Open("postgres", connString)
//	        }
//	    },
//	    // onRelease: Close connection
//	    func(db *sql.DB) readerresult.ReaderResult[any] {
//	        return func(ctx context.Context) (any, error) {
//	            return nil, db.Close()
//	        }
//	    },
//	)
//
//	// Use the manager with different operations
//	getUsers := withDB(func(db *sql.DB) readerresult.ReaderResult[[]User] {
//	    return func(ctx context.Context) ([]User, error) {
//	        return queryUsers(ctx, db)
//	    }
//	})
//
//	getOrders := withDB(func(db *sql.DB) readerresult.ReaderResult[[]Order] {
//	    return func(ctx context.Context) ([]Order, error) {
//	        return queryOrders(ctx, db)
//	    }
//	})
//
//	// Both operations automatically manage the connection
//	users, err := getUsers(context.Background())
//	orders, err := getOrders(context.Background())
//
// Example - File operations manager:
//
//	withFile := readerresult.WithResource(
//	    func() readerresult.ReaderResult[*os.File] {
//	        return func(ctx context.Context) (*os.File, error) {
//	            return os.Open("config.json")
//	        }
//	    },
//	    func(file *os.File) readerresult.ReaderResult[any] {
//	        return func(ctx context.Context) (any, error) {
//	            return nil, file.Close()
//	        }
//	    },
//	)
//
//	// Different operations on the same file
//	readConfig := withFile(func(file *os.File) readerresult.ReaderResult[Config] {
//	    return func(ctx context.Context) (Config, error) {
//	        return parseConfig(file)
//	    }
//	})
//
//	validateConfig := withFile(func(file *os.File) readerresult.ReaderResult[bool] {
//	    return func(ctx context.Context) (bool, error) {
//	        return validateConfigFile(file)
//	    }
//	})
//
// Example - Composing with other operations:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Create a pipeline with automatic resource management
//	processData := F.Pipe2(
//	    loadData,
//	    withDB(func(db *sql.DB) readerresult.ReaderResult[Result] {
//	        return saveToDatabase(db)
//	    }),
//	    readerresult.Map(formatResult),
//	)
//
//go:inline
func WithResource[B, A, ANY any](
	onCreate Lazy[ReaderResult[A]],
	onRelease Kleisli[A, ANY],
) Kleisli[Kleisli[A, B], B] {
	return WithContextK(RR.WithResource[B](onCreate, onRelease))
}

// onClose is a helper function that creates a ReaderResult that closes an io.Closer.
// This is used internally by WithCloser to provide automatic cleanup for resources
// that implement the io.Closer interface.
func onClose[A io.Closer](a A) ReaderResult[struct{}] {
	return func(_ context.Context) (struct{}, error) {
		return struct{}{}, a.Close()
	}
}

// WithCloser creates a higher-order function for managing resources that implement io.Closer.
//
// This is a specialized version of WithResource that automatically handles cleanup for any
// resource implementing the io.Closer interface (such as files, network connections, HTTP
// response bodies, etc.). It eliminates the need to manually specify the release function,
// making it more convenient for common Go resources.
//
// The function automatically calls Close() on the resource when the operation completes,
// regardless of success or failure. This ensures proper resource cleanup following Go's
// standard io.Closer pattern.
//
// Type Parameters:
//   - B: The type of the result produced by using the resource
//   - A: The type of the resource, which must implement io.Closer
//
// Parameters:
//   - onCreate: Lazy computation that creates/acquires the io.Closer resource
//
// Returns:
//   - A Kleisli arrow that takes a resource-using function and returns a ReaderResult[B]
//     with automatic Close() cleanup
//
// Example - File operations:
//
//	import (
//	    "context"
//	    "os"
//	    "io"
//	)
//
//	// Create a reusable file manager
//	withFile := readerresult.WithCloser(
//	    func() readerresult.ReaderResult[*os.File] {
//	        return func(ctx context.Context) (*os.File, error) {
//	            return os.Open("data.txt")
//	        }
//	    },
//	)
//
//	// Use with different operations - Close() is automatic
//	readContent := withFile(func(file *os.File) readerresult.ReaderResult[string] {
//	    return func(ctx context.Context) (string, error) {
//	        data, err := io.ReadAll(file)
//	        return string(data), err
//	    }
//	})
//
//	getSize := withFile(func(file *os.File) readerresult.ReaderResult[int64] {
//	    return func(ctx context.Context) (int64, error) {
//	        info, err := file.Stat()
//	        if err != nil {
//	            return 0, err
//	        }
//	        return info.Size(), nil
//	    }
//	})
//
//	content, err := readContent(context.Background())
//	size, err := getSize(context.Background())
//
// Example - HTTP response body:
//
//	import "net/http"
//
//	withResponse := readerresult.WithCloser(
//	    func() readerresult.ReaderResult[*http.Response] {
//	        return func(ctx context.Context) (*http.Response, error) {
//	            return http.Get("https://api.example.com/data")
//	        }
//	    },
//	)
//
//	// Body is automatically closed after use
//	parseJSON := withResponse(func(resp *http.Response) readerresult.ReaderResult[Data] {
//	    return func(ctx context.Context) (Data, error) {
//	        var data Data
//	        err := json.NewDecoder(resp.Body).Decode(&data)
//	        return data, err
//	    }
//	})
//
// Example - Multiple file operations:
//
//	// Read from one file, write to another
//	copyFile := func(src, dst string) readerresult.ReaderResult[int64] {
//	    withSrc := readerresult.WithCloser(
//	        func() readerresult.ReaderResult[*os.File] {
//	            return func(ctx context.Context) (*os.File, error) {
//	                return os.Open(src)
//	            }
//	        },
//	    )
//
//	    withDst := readerresult.WithCloser(
//	        func() readerresult.ReaderResult[*os.File] {
//	            return func(ctx context.Context) (*os.File, error) {
//	                return os.Create(dst)
//	            }
//	        },
//	    )
//
//	    return withSrc(func(srcFile *os.File) readerresult.ReaderResult[int64] {
//	        return withDst(func(dstFile *os.File) readerresult.ReaderResult[int64] {
//	            return func(ctx context.Context) (int64, error) {
//	                return io.Copy(dstFile, srcFile)
//	            }
//	        })
//	    })
//	}
//
// Example - Network connection:
//
//	import "net"
//
//	withConn := readerresult.WithCloser(
//	    func() readerresult.ReaderResult[net.Conn] {
//	        return func(ctx context.Context) (net.Conn, error) {
//	            return net.Dial("tcp", "localhost:8080")
//	        }
//	    },
//	)
//
//	sendData := withConn(func(conn net.Conn) readerresult.ReaderResult[int] {
//	    return func(ctx context.Context) (int, error) {
//	        return conn.Write([]byte("Hello, World!"))
//	    }
//	})
//
// Note: WithCloser is a convenience wrapper around WithResource that automatically
// provides the Close() cleanup function. For resources that don't implement io.Closer
// or require custom cleanup logic, use WithResource or Bracket instead.
//
//go:inline
func WithCloser[B any, A io.Closer](onCreate Lazy[ReaderResult[A]]) Kleisli[Kleisli[A, B], B] {
	return WithResource[B](onCreate, onClose[A])
}
