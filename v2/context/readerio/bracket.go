package readerio

import (
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
