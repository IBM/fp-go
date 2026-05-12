package iterresult

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	// Trampoline represents a tail-recursive computation that can be evaluated safely
	// without stack overflow. It's used for implementing stack-safe recursive algorithms.
	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	// Void represents the absence of a value, similar to void in other languages.
	// It's used in functions that perform side effects but don't return meaningful values.
	// The zero value is function.VOID.
	//
	// Example:
	//
	//	logMessage := func() Void {
	//	    fmt.Println("Logging...")
	//	    return function.VOID
	//	}
	Void = function.Void

	// Lazy represents a deferred computation that produces a value of type A.
	// The computation is not executed until the Lazy value is called.
	// This is an alias for IO and provides semantic clarity when working with
	// lazy evaluations that don't necessarily involve side effects.
	//
	// Type Parameters:
	//   - A: The type of value produced by the lazy computation
	//
	// Example:
	//
	//	expensiveCalc := func() int {
	//	    // Expensive computation
	//	    return 42
	//	}
	//	// expensiveCalc is a Lazy[int]
	Lazy[A any] = lazy.Lazy[A]

	// Endomorphism represents a function from a type to itself.
	// It represents transformations that preserve the type.
	//
	// Type Parameters:
	//   - T: The type being transformed
	//
	// Example:
	//
	//	increment := func(x int) int { return x + 1 }
	//	// increment is an Endomorphism[int]
	Endomorphism[T any] = endomorphism.Endomorphism[T]

	// IO represents a computation that performs side effects and returns a value of type T.
	// IO computations are lazy - they don't execute until explicitly called.
	// This allows for composing side-effecting operations in a pure, functional way.
	//
	// Type Parameters:
	//   - T: The type of value produced by the IO computation
	//
	// Example:
	//
	//	getCurrentTime := func() time.Time { return time.Now() }
	//	// getCurrentTime is an IO[time.Time]
	IO[T any] = io.IO[T]

	// IOResult represents an IO computation that can either fail with an error
	// or succeed with a value of type T. It combines the side-effect handling of IO
	// with the error handling of Result.
	//
	// This is useful for operations that perform I/O and may fail, providing a
	// type-safe alternative to returning (T, error) tuples.
	//
	// Type Parameters:
	//   - T: The success type
	//
	// Example:
	//
	//	readFile := func() result.Result[string] {
	//	    data, err := os.ReadFile("config.txt")
	//	    if err != nil {
	//	        return result.Left(err)
	//	    }
	//	    return result.Of(string(data))
	//	}
	//	// readFile is an IOResult[string]
	IOResult[T any] = ioresult.IOResult[T]
)
