package generic

import (
	"github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Monoid returns a Monoid instance for arrays.
// The Monoid combines arrays through concatenation, with an empty array as the identity element.
//
// Example:
//
//	m := array.Monoid[int]()
//	result := m.Concat([]int{1, 2}, []int{3, 4}) // [1, 2, 3, 4]
//	empty := m.Empty() // []
//
//go:inline
func Monoid[GT ~[]T, T any]() M.Monoid[GT] {
	return M.MakeMonoid(array.Concat[GT], Empty[GT]())
}

// Semigroup returns a Semigroup instance for arrays.
// The Semigroup combines arrays through concatenation.
//
// Example:
//
//	s := array.Semigroup[int]()
//	result := s.Concat([]int{1, 2}, []int{3, 4}) // [1, 2, 3, 4]
//
//go:inline
func Semigroup[GT ~[]T, T any]() S.Semigroup[GT] {
	return S.MakeSemigroup(array.Concat[GT])
}
