package generic

import F "github.com/IBM/fp-go/v2/function"

// StrictUniq converts an array of arbitrary items into an array or unique items
// where uniqueness is determined by the built-in uniqueness constraint
func StrictUniq[AS ~[]A, A comparable](as AS) AS {
	return Uniq[AS](F.Identity[A])(as)
}

// uniquePredUnsafe returns a predicate on a map for uniqueness
func uniquePredUnsafe[PRED ~func(A) K, A any, K comparable](f PRED) func(int, A) bool {
	lookup := make(map[K]bool)
	return func(_ int, a A) bool {
		k := f(a)
		_, has := lookup[k]
		if has {
			return false
		}
		lookup[k] = true
		return true
	}
}

// Uniq converts an array of arbitrary items into an array or unique items
// where uniqueness is determined based on a key extractor function
func Uniq[AS ~[]A, PRED ~func(A) K, A any, K comparable](f PRED) func(as AS) AS {
	return func(as AS) AS {
		// we need to create a new predicate for each iteration
		return filterWithIndex(as, uniquePredUnsafe(f))
	}
}
