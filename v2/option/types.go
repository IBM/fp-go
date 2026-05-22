package option

import (
	"iter"

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Seq represents an iterator sequence over values of type T.
	// It's an alias for Go's standard iter.Seq[T] type.
	Seq[T any] = iter.Seq[T]

	// Endomorphism represents a function from a type to itself (T -> T).
	Endomorphism[T any] = endomorphism.Endomorphism[T]

	Predicate[T any] = predicate.Predicate[T]

	// Traversable represents a data structure that can be traversed from left to right,
	// applying an effectful function to each element and collecting the results.
	//
	// A Traversable takes a Kleisli arrow (a function that returns an Option) and
	// produces another Kleisli arrow that operates on a container of values.
	//
	// Type Parameters:
	//   - A: The input element type
	//   - B: The output element type
	//   - GA: The input container type (e.g., []A, map[K]A)
	//   - GB: The output container type (e.g., []B, map[K]B)
	//
	// The Traversable signature:
	//
	//	func(Kleisli[A, B]) Kleisli[GA, GB]
	//
	// expands to:
	//
	//	func(func(A) Option[B]) func(GA) Option[GB]
	//
	// This means: given a function that transforms A to Option[B], produce a function
	// that transforms a container of A values into an Option of a container of B values.
	//
	// Behavior:
	//   - If all transformations succeed (return Some), the result is Some containing
	//     the container of all transformed values
	//   - If any transformation fails (returns None), the entire result is None
	//
	// Common Use Cases:
	//   - Validating and transforming collections where any failure should fail the whole operation
	//   - Parsing collections of strings where all must parse successfully
	//   - Applying optional transformations across data structures
	//
	// Example:
	//
	//	// Array traversable
	//	traversable := TraversableArray[string, int]()
	//	parse := func(s string) Option[int] {
	//	    n, err := strconv.Atoi(s)
	//	    if err != nil { return None[int]() }
	//	    return Some(n)
	//	}
	//	result := traversable(parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
	//	result := traversable(parse)([]string{"1", "x", "3"}) // None
	//
	// See Also:
	//   - TraversableArray: Traversable instance for arrays
	//   - TraverseArray: Direct array traversal function
	//   - TraverseRecord: Traversal for maps/records
	Traversable[A, B, GA, GB any] = func(Kleisli[A, B]) Kleisli[GA, GB]
)
