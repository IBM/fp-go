package generic

import (
	"sort"

	O "github.com/ibm/fp-go/ord"
)

// Sort implements a stable sort on the array given the provided ordering
func Sort[GA ~[]T, T any](ord O.Ord[T]) func(ma GA) GA {

	return func(ma GA) GA {
		// nothing to sort
		l := len(ma)
		if l < 2 {
			return ma
		}
		// copy
		cpy := make(GA, l)
		copy(cpy, ma)
		sort.Slice(cpy, func(i, j int) bool {
			return ord.Compare(cpy[i], cpy[j]) < 0
		})
		return cpy
	}
}
