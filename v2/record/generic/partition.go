package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
)

// MonadPartition splits a map into two maps based on a predicate applied to
// each value. Entries for which the predicate returns true are placed in the
// right (tail) map; entries for which it returns false are placed in the left
// (head) map.
//
// Type Parameters:
//   - M: map type, must be ~map[K]V
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - kvs: the source map to partition
//   - pred: a function that tests each value
//
// Returns:
//   - pair.Pair[M, M]: a pair where Head holds the non-matching entries and
//     Tail holds the matching entries
//
// See Also:
//   - MonadPartitionWithIndex: index-aware variant
//   - Partition: curried variant
func MonadPartition[M ~map[K]V, K comparable, V any](kvs M, pred func(V) bool) pair.Pair[M, M] {
	left := make(M)
	right := make(M)
	for k, v := range kvs {
		if pred(v) {
			right[k] = v
		} else {
			left[k] = v
		}
	}
	// returns the partition
	return pair.MakePair(left, right)
}

// MonadPartitionWithIndex splits a map into two maps based on a predicate
// applied to each key-value pair. Entries for which the predicate returns true
// are placed in the right (tail) map; entries for which it returns false are
// placed in the left (head) map.
//
// Type Parameters:
//   - M: map type, must be ~map[K]V
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - kvs: the source map to partition
//   - pred: a function that tests each key-value pair; returns true to route
//     the entry to the right (tail) map
//
// Returns:
//   - pair.Pair[M, M]: a pair where Head holds the non-matching entries and
//     Tail holds the matching entries
//
// See Also:
//   - MonadPartition: value-only variant
//   - PartitionWithIndex: curried variant
func MonadPartitionWithIndex[M ~map[K]V, K comparable, V any](kvs M, pred func(K, V) bool) pair.Pair[M, M] {
	left := make(M)
	right := make(M)
	for k, v := range kvs {
		if pred(k, v) {
			right[k] = v
		} else {
			left[k] = v
		}
	}
	// returns the partition
	return pair.MakePair(left, right)
}

// Partition returns a curried function that splits a map into two maps based
// on a predicate applied to each value. The returned function accepts a map and
// returns a pair where Head holds the non-matching entries and Tail holds the
// matching entries.
//
// Type Parameters:
//   - M: map type, must be ~map[K]V
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - pred: a function that tests each value
//
// Returns:
//   - pair.Kleisli[M, M, M]: a reusable function from M to Pair[M, M]
//
// See Also:
//   - PartitionWithIndex: index-aware variant
//   - MonadPartition: uncurried variant
func Partition[M ~map[K]V, K comparable, V any](pred func(V) bool) pair.Kleisli[M, M, M] {
	return F.Bind2nd(MonadPartition[M], pred)
}

// PartitionWithIndex returns a curried function that splits a map into two maps
// based on a predicate applied to each key-value pair. The returned function
// accepts a map and returns a pair where Head holds the non-matching entries
// and Tail holds the matching entries.
//
// Type Parameters:
//   - M: map type, must be ~map[K]V
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - pred: a function that tests each key-value pair; returns true to route
//     the entry to the right (tail) map
//
// Returns:
//   - pair.Kleisli[M, M, M]: a reusable function from M to Pair[M, M]
//
// See Also:
//   - Partition: value-only variant
//   - MonadPartitionWithIndex: uncurried variant
func PartitionWithIndex[M ~map[K]V, K comparable, V any](pred func(K, V) bool) pair.Kleisli[M, M, M] {
	return F.Bind2nd(MonadPartitionWithIndex[M], pred)
}
