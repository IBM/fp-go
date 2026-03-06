package readerio

import (
	"github.com/IBM/fp-go/v2/function"
	RIO "github.com/IBM/fp-go/v2/readerio"
)

// Bracket ensures that a resource is properly acquired, used, and released, even if an error occurs.
// This implements the bracket pattern for safe resource management with [ReaderIO].
//
// The bracket pattern guarantees that:
//   - The acquire action is executed first to obtain the resource
//   - The use function is called with the acquired resource
//   - The release function is always called with the resource and result, regardless of success or failure
//   - The final result from the use function is returned
//
// This is particularly useful for managing resources like file handles, database connections,
// or locks that must be cleaned up properly.
//
// Type Parameters:
//   - A: The type of the acquired resource
//   - B: The type of the result produced by the use function
//   - ANY: The type returned by the release function (typically ignored)
//
// Parameters:
//   - acquire: A ReaderIO that acquires the resource
//   - use: A Kleisli arrow that uses the resource and produces a result
//   - release: A function that releases the resource, receiving both the resource and the result
//
// Returns:
//   - A ReaderIO[B] that safely manages the resource lifecycle
//
// Example:
//
//	// Acquire a file handle
//	acquireFile := func(ctx context.Context) IO[*os.File] {
//	    return func() *os.File {
//	        f, _ := os.Open("data.txt")
//	        return f
//	    }
//	}
//
//	// Use the file
//	readFile := func(f *os.File) ReaderIO[string] {
//	    return func(ctx context.Context) IO[string] {
//	        return func() string {
//	            data, _ := io.ReadAll(f)
//	            return string(data)
//	        }
//	    }
//	}
//
//	// Release the file
//	closeFile := func(f *os.File, result string) ReaderIO[any] {
//	    return func(ctx context.Context) IO[any] {
//	        return func() any {
//	            f.Close()
//	            return nil
//	        }
//	    }
//	}
//
//	// Safely read file with automatic cleanup
//	safeRead := Bracket(acquireFile, readFile, closeFile)
//	result := safeRead(t.Context())()
//
//go:inline
func Bracket[
	A, B, ANY any](

	acquire ReaderIO[A],
	use Kleisli[A, B],
	release func(A, B) ReaderIO[ANY],
) ReaderIO[B] {
	return RIO.Bracket(acquire, use, release)
}

// WithResource creates a higher-order function that manages a resource lifecycle for any operation.
// It returns a Kleisli arrow that takes a use function and automatically handles resource
// acquisition and cleanup using the bracket pattern.
//
// This is a more composable alternative to Bracket, allowing you to define resource management
// once and reuse it with different use functions. The resource is acquired when the returned
// Kleisli arrow is invoked, used by the provided function, and then released regardless of
// success or failure.
//
// Type Parameters:
//   - A: The type of the resource to be managed
//   - B: The type of the result produced by the use function
//   - ANY: The type returned by the release function (typically ignored)
//
// Parameters:
//   - onCreate: A ReaderIO that acquires/creates the resource
//   - onRelease: A Kleisli arrow that releases/cleans up the resource
//
// Returns:
//   - A Kleisli arrow that takes a use function and returns a ReaderIO managing the full lifecycle
//
// Example with database connection:
//
//	// Define resource management once
//	withDB := WithResource(
//	    // Acquire connection
//	    func(ctx context.Context) IO[*sql.DB] {
//	        return func() *sql.DB {
//	            db, _ := sql.Open("postgres", "connection-string")
//	            return db
//	        }
//	    },
//	    // Release connection
//	    func(db *sql.DB) ReaderIO[any] {
//	        return func(ctx context.Context) IO[any] {
//	            return func() any {
//	                db.Close()
//	                return nil
//	            }
//	        }
//	    },
//	)
//
//	// Reuse with different operations
//	queryUsers := withDB(func(db *sql.DB) ReaderIO[[]User] {
//	    return func(ctx context.Context) IO[[]User] {
//	        return func() []User {
//	            // Query users from db
//	            return users
//	        }
//	    }
//	})
//
//	insertUser := withDB(func(db *sql.DB) ReaderIO[int64] {
//	    return func(ctx context.Context) IO[int64] {
//	        return func() int64 {
//	            // Insert user into db
//	            return userID
//	        }
//	    }
//	})
//
// Example with file handling:
//
//	withFile := WithResource(
//	    func(ctx context.Context) IO[*os.File] {
//	        return func() *os.File {
//	            f, _ := os.Open("data.txt")
//	            return f
//	        }
//	    },
//	    func(f *os.File) ReaderIO[any] {
//	        return func(ctx context.Context) IO[any] {
//	            return func() any {
//	                f.Close()
//	                return nil
//	            }
//	        }
//	    },
//	)
//
//	// Use for reading
//	readContent := withFile(func(f *os.File) ReaderIO[string] {
//	    return func(ctx context.Context) IO[string] {
//	        return func() string {
//	            data, _ := io.ReadAll(f)
//	            return string(data)
//	        }
//	    }
//	})
//
//	// Use for getting file info
//	getSize := withFile(func(f *os.File) ReaderIO[int64] {
//	    return func(ctx context.Context) IO[int64] {
//	        return func() int64 {
//	            info, _ := f.Stat()
//	            return info.Size()
//	        }
//	    }
//	})
//
// Use Cases:
//   - Database connections: Acquire connection, execute queries, close connection
//   - File handles: Open file, read/write, close file
//   - Network connections: Establish connection, transfer data, close connection
//   - Locks: Acquire lock, perform critical section, release lock
//   - Temporary resources: Create temp file/directory, use it, clean up
//
//go:inline
func WithResource[A, B, ANY any](
	onCreate ReaderIO[A], onRelease Kleisli[A, ANY]) Kleisli[Kleisli[A, B], B] {
	return function.Bind13of3(Bracket[A, B, ANY])(onCreate, function.Ignore2of2[B](onRelease))
}
