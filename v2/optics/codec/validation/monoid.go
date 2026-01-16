package validation

import (
	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	M "github.com/IBM/fp-go/v2/monoid"
)

// ErrorsMonoid returns a Monoid instance for Errors (array of ValidationError pointers).
// The monoid concatenates error arrays, with an empty array as the identity element.
// This is used internally by the applicative operations to accumulate validation errors.
//
// Example:
//
//	m := ErrorsMonoid()
//	combined := m.Concat(errors1, errors2) // Concatenates both error arrays
//	empty := m.Empty()                      // Returns empty error array
func ErrorsMonoid() Monoid[Errors] {
	return A.Monoid[*ValidationError]()
}

// ApplicativeMonoid creates a Monoid instance for Validation[A] given a Monoid for A.
// This allows combining validation results where the success values are also combined
// using the provided monoid. If any validation fails, all errors are accumulated.
//
// The resulting monoid:
//   - Empty: Returns a successful validation with the empty value from the inner monoid
//   - Concat: Combines two validations:
//   - Both success: Combines values using the inner monoid
//   - Any failure: Accumulates all errors
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//
//	// Create a monoid for validations of strings
//	m := ApplicativeMonoid(string.Monoid)
//
//	v1 := Success("Hello")
//	v2 := Success(" World")
//	combined := m.Concat(v1, v2) // Success("Hello World")
//
//	v3 := Failures[string](someErrors)
//	failed := m.Concat(v1, v3) // Failures with accumulated errors
func ApplicativeMonoid[A any](m Monoid[A]) Monoid[Validation[A]] {

	return M.ApplicativeMonoid(
		Of,
		either.MonadMap,
		either.MonadApV[A, A](ErrorsMonoid()),

		m,
	)
}
