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

// AlternativeMonoid creates a Monoid instance for Validation[A] using the Alternative pattern.
// This combines the applicative error-accumulation behavior with the alternative fallback behavior,
// allowing you to both accumulate errors and provide fallback alternatives.
//
// The Alternative pattern provides two key operations:
//   - Applicative operations (Of, Map, Ap): accumulate errors when combining validations
//   - Alternative operation (Alt): provide fallback when a validation fails
//
// This monoid is particularly useful when you want to:
//   - Try multiple validation strategies and fall back to alternatives
//   - Combine successful values using the provided monoid
//   - Accumulate all errors from failed attempts
//   - Build validation pipelines with fallback logic
//
// The resulting monoid:
//   - Empty: Returns a successful validation with the empty value from the inner monoid
//   - Concat: Combines two validations using both applicative and alternative semantics:
//   - If first succeeds and second succeeds: combines values using inner monoid
//   - If first fails: tries second as fallback (alternative behavior)
//   - If both fail: accumulates all errors
//
// Type Parameters:
//   - A: The type of the successful value
//
// Parameters:
//   - m: The monoid for combining successful values of type A
//
// Returns:
//
//	A Monoid[Validation[A]] that combines applicative and alternative behaviors
//
// Example - Combining successful validations:
//
//	import "github.com/IBM/fp-go/v2/string"
//
//	m := AlternativeMonoid(string.Monoid)
//	v1 := Success("Hello")
//	v2 := Success(" World")
//	result := m.Concat(v1, v2)
//	// Result: Success("Hello World")
//
// Example - Fallback behavior:
//
//	m := AlternativeMonoid(string.Monoid)
//	v1 := Failures[string](Errors{&ValidationError{Messsage: "first failed"}})
//	v2 := Success("fallback value")
//	result := m.Concat(v1, v2)
//	// Result: Success("fallback value") - second validation used as fallback
//
// Example - Error accumulation when both fail:
//
//	m := AlternativeMonoid(string.Monoid)
//	v1 := Failures[string](Errors{&ValidationError{Messsage: "error 1"}})
//	v2 := Failures[string](Errors{&ValidationError{Messsage: "error 2"}})
//	result := m.Concat(v1, v2)
//	// Result: Failures with accumulated errors: ["error 1", "error 2"]
//
// Example - Building validation with fallbacks:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	m := AlternativeMonoid(N.MonoidSum[int]())
//
//	// Try to parse from different sources
//	fromEnv := parseFromEnv()      // Fails
//	fromConfig := parseFromConfig() // Succeeds with 42
//	fromDefault := Success(0)       // Default fallback
//
//	result := m.Concat(m.Concat(fromEnv, fromConfig), fromDefault)
//	// Result: Success(42) - uses first successful validation
func AlternativeMonoid[A any](m Monoid[A]) Monoid[Validation[A]] {
	return M.AlternativeMonoid(
		Of[A],
		MonadMap[A, func(A) A],
		MonadAp[A, A],
		MonadAlt[A],
		m,
	)
}

// AltMonoid creates a Monoid instance for Validation[A] using the Alt (alternative) operation.
// This monoid provides a way to combine validations with fallback behavior, where the second
// validation is used as an alternative if the first one fails.
//
// The Alt operation implements the "try first, fallback to second" pattern, which is useful
// for validation scenarios where you want to attempt multiple validation strategies in sequence
// and use the first one that succeeds.
//
// The resulting monoid:
//   - Empty: Returns the provided zero value (a lazy computation that produces a Validation[A])
//   - Concat: Combines two validations using Alt semantics:
//   - If first succeeds: returns the first validation (ignores second)
//   - If first fails: returns the second validation as fallback
//
// This is different from [AlternativeMonoid] in that:
//   - AltMonoid uses a custom zero value (provided by the user)
//   - AlternativeMonoid derives the zero from an inner monoid
//   - AltMonoid is simpler and only provides fallback behavior
//   - AlternativeMonoid combines applicative and alternative behaviors
//
// Type Parameters:
//   - A: The type of the successful value
//
// Parameters:
//   - zero: A lazy computation that produces the identity/empty Validation[A].
//     This is typically a successful validation with a default value, or could be
//     a failure representing "no validation attempted"
//
// Returns:
//
//	A Monoid[Validation[A]] that combines validations with fallback behavior
//
// Example - Using default value as zero:
//
//	m := AltMonoid(func() Validation[int] { return Success(0) })
//
//	v1 := Failures[int](Errors{&ValidationError{Messsage: "failed"}})
//	v2 := Success(42)
//
//	result := m.Concat(v1, v2)
//	// Result: Success(42) - falls back to second validation
//
//	empty := m.Empty()
//	// Result: Success(0) - the provided zero value
//
// Example - Chaining multiple fallbacks:
//
//	m := AltMonoid(func() Validation[string] {
//	    return Success("default")
//	})
//
//	primary := parseFromPrimarySource()   // Fails
//	secondary := parseFromSecondary()     // Fails
//	tertiary := parseFromTertiary()       // Succeeds with "value"
//
//	result := m.Concat(m.Concat(primary, secondary), tertiary)
//	// Result: Success("value") - uses first successful validation
//
// Example - All validations fail:
//
//	m := AltMonoid(func() Validation[int] {
//	    return Failures[int](Errors{&ValidationError{Messsage: "no default"}})
//	})
//
//	v1 := Failures[int](Errors{&ValidationError{Messsage: "error 1"}})
//	v2 := Failures[int](Errors{&ValidationError{Messsage: "error 2"}})
//
//	result := m.Concat(v1, v2)
//	// Result: Failures with errors from v2: ["error 2"]
//	// Note: Unlike AlternativeMonoid, errors are NOT accumulated
//
// Example - Building a validation pipeline with fallbacks:
//
//	m := AltMonoid(func() Validation[Config] {
//	    return Success(defaultConfig)
//	})
//
//	// Try multiple configuration sources in order
//	configs := []Validation[Config]{
//	    loadFromFile("config.json"),      // Try file first
//	    loadFromEnv(),                     // Then environment
//	    loadFromRemote("api.example.com"), // Then remote API
//	}
//
//	// Fold using the monoid to get first successful config
//	result := A.MonoidFold(m)(configs)
//	// Result: First successful config, or defaultConfig if all fail
func AltMonoid[A any](zero Lazy[Validation[A]]) Monoid[Validation[A]] {
	return M.AltMonoid(
		zero,
		MonadAlt[A],
	)
}
