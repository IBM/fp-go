package iter

import (
	"github.com/IBM/fp-go/v2/option"
)

// Last returns the last element from an [Iterator] wrapped in an [Option].
//
// This function retrieves the last element from the iterator by consuming the entire
// sequence. If the iterator contains at least one element, it returns Some(element).
// If the iterator is empty, it returns None.
//
// RxJS Equivalent: [last] - https://rxjs.dev/api/operators/last
//
// Type Parameters:
//   - U: The type of elements in the iterator
//
// Parameters:
//   - it: The input iterator to get the last element from
//
// Returns:
//   - Option[U]: Some(last element) if the iterator is non-empty, None otherwise
//
// Example with non-empty sequence:
//
//	seq := iter.From(1, 2, 3, 4, 5)
//	last := iter.Last(seq)
//	// Returns: Some(5)
//
// Example with empty sequence:
//
//	seq := iter.Empty[int]()
//	last := iter.Last(seq)
//	// Returns: None
//
// Example with filtered sequence:
//
//	seq := iter.From(1, 2, 3, 4, 5)
//	filtered := iter.Filter(func(x int) bool { return x < 4 })(seq)
//	last := iter.Last(filtered)
//	// Returns: Some(3)
func Last[U any](it Seq[U]) Option[U] {
	var last U
	found := false

	for last = range it {
		found = true
	}

	if !found {
		return option.None[U]()
	}
	return option.Some(last)
}
