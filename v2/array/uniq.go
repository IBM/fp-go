package array

import (
	G "github.com/IBM/fp-go/v2/array/generic"
)

// StrictUniq converts an array of arbitrary items into an array of unique items
// where uniqueness is determined by the built-in equality constraint (comparable).
// The first occurrence of each unique value is kept, subsequent duplicates are removed.
//
// Example:
//
//	numbers := []int{1, 2, 2, 3, 3, 3, 4}
//	unique := array.StrictUniq(numbers) // [1, 2, 3, 4]
//
//	strings := []string{"a", "b", "a", "c", "b"}
//	unique2 := array.StrictUniq(strings) // ["a", "b", "c"]
//
//go:inline
func StrictUniq[A comparable](as []A) []A {
	return G.StrictUniq(as)
}

// Uniq converts an array of arbitrary items into an array of unique items
// where uniqueness is determined based on a key extractor function.
// The first occurrence of each unique key is kept, subsequent duplicates are removed.
//
// This is useful for removing duplicates from arrays of complex types based on a specific field.
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	people := []Person{
//	    {"Alice", 30},
//	    {"Bob", 25},
//	    {"Alice", 35}, // duplicate name
//	    {"Charlie", 30},
//	}
//
//	uniqueByName := array.Uniq(func(p Person) string { return p.Name })
//	result := uniqueByName(people)
//	// Result: [{"Alice", 30}, {"Bob", 25}, {"Charlie", 30}]
//
//go:inline
func Uniq[A any, K comparable](f func(A) K) Operator[A, A] {
	return G.Uniq[[]A](f)
}
