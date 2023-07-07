package predicate

import (
	F "github.com/ibm/fp-go/function"
)

// ContraMap creates a predicate from an existing predicate given a mapping function
func ContraMap[A, B any](f func(B) A) func(func(A) bool) func(B) bool {
	return func(pred func(A) bool) func(B) bool {
		return F.Flow2(
			f,
			pred,
		)
	}
}
