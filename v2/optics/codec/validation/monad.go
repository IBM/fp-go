package validation

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/internal/applicative"
)

// Of creates a successful validation result containing the given value.
// This is the pure/return operation for the Validation monad.
//
// Example:
//
//	valid := Of(42) // Validation[int] containing 42
func Of[A any](a A) Validation[A] {
	return either.Of[Errors](a)
}

// Ap applies a validation containing a function to a validation containing a value.
// This is the applicative apply operation that accumulates errors from both validations.
// If either validation fails, all errors are collected. If both succeed, the function is applied.
//
// This enables combining multiple validations while collecting all errors:
//
// Example:
//
//	// Validate multiple fields and collect all errors
//	validateUser := Ap(Ap(Of(func(name string) func(age int) User {
//		return func(age int) User { return User{name, age} }
//	}))(validateName))(validateAge)
func Ap[B, A any](fa Validation[A]) Operator[func(A) B, B] {
	return either.ApV[B, A](ErrorsMonoid())(fa)
}

// Map transforms the value inside a successful validation using the provided function.
// If the validation is a failure, the errors are preserved unchanged.
// This is the functor map operation for Validation.
//
// Example:
//
//	doubled := Map(func(x int) int { return x * 2 })(Of(21))
//	// Result: Success(42)
func Map[A, B any](f func(A) B) Operator[A, B] {
	return either.Map[Errors](f)
}

// Applicative creates an Applicative instance for Validation with error accumulation.
//
// This returns a lawful Applicative that accumulates validation errors using the Errors monoid.
// Unlike the standard Either applicative which fails fast, this validation applicative collects
// all errors when combining independent validations with Ap.
//
// The returned instance satisfies all applicative laws:
//   - Identity: Ap(Of(identity))(v) == v
//   - Homomorphism: Ap(Of(f))(Of(x)) == Of(f(x))
//   - Interchange: Ap(Of(f))(u) == Ap(Map(f => f(y))(u))(Of(y))
//   - Composition: Ap(Ap(Map(compose)(f))(g))(x) == Ap(f)(Ap(g)(x))
//
// Key behaviors:
//   - Of: lifts a value into a successful Validation (Right)
//   - Map: transforms successful values, preserves failures (standard functor)
//   - Ap: when both operands fail, combines all errors using the Errors monoid
//
// This is particularly useful for form validation, configuration validation, and any scenario
// where you want to collect all validation errors at once rather than stopping at the first failure.
//
// Example - Validating Multiple Fields:
//
//	app := Applicative[string, User]()
//
//	// Validate individual fields
//	validateName := func(name string) Validation[string] {
//		if len(name) < 3 {
//			return Failure("Name must be at least 3 characters")
//		}
//		return Success(name)
//	}
//
//	validateAge := func(age int) Validation[int] {
//		if age < 18 {
//			return Failure("Must be 18 or older")
//		}
//		return Success(age)
//	}
//
//	// Create a curried constructor
//	makeUser := func(name string) func(int) User {
//		return func(age int) User {
//			return User{Name: name, Age: age}
//		}
//	}
//
//	// Combine validations - all errors are collected
//	name := validateName("ab")  // Failure: name too short
//	age := validateAge(16)      // Failure: age too low
//
//	result := app.Ap(age)(app.Ap(name)(app.Of(makeUser)))
//	// result contains both validation errors:
//	// - "Name must be at least 3 characters"
//	// - "Must be 18 or older"
//
// Type Parameters:
//   - A: The input value type (Right value)
//   - B: The output value type after transformation
//
// Returns:
//
//	An Applicative instance with Of, Map, and Ap operations that accumulate errors
func Applicative[A, B any]() applicative.Applicative[A, B, Validation[A], Validation[B], Validation[func(A) B]] {
	return either.ApplicativeV[Errors, A, B](
		ErrorsMonoid(),
	)
}
