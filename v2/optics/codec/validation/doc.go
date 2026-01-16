// Package validation provides functional validation types and operations for the codec system.
//
// This package implements a validation monad that accumulates errors during validation operations,
// making it ideal for form validation, data parsing, and other scenarios where you want to collect
// all validation errors rather than failing on the first error.
//
// # Core Concepts
//
// Validation[A]: Represents the result of a validation operation as Either[Errors, A]:
//   - Left(Errors): Validation failed with one or more errors
//   - Right(A): Successfully validated value of type A
//
// ValidationError: A detailed error type that includes:
//   - Value: The actual value that failed validation
//   - Context: The path through nested structures (e.g., "user.address.zipCode")
//   - Message: Human-readable error description
//   - Cause: Optional underlying error
//
// Context: A stack of ContextEntry values that tracks the validation path through
// nested data structures, enabling precise error reporting.
//
// # Basic Usage
//
// Creating validation results:
//
//	// Success case
//	valid := validation.Success(42)
//
//	// Failure case
//	invalid := validation.Failures[int](validation.Errors{
//		&validation.ValidationError{
//			Value:    "not a number",
//			Message:  "expected integer",
//			Context:  nil,
//		},
//	})
//
// Using with context:
//
//	failWithMsg := validation.FailureWithMessage[int]("invalid", "must be positive")
//	result := failWithMsg([]validation.ContextEntry{
//		{Key: "age", Type: "int"},
//	})
//
// # Applicative Validation
//
// The validation type supports applicative operations, allowing you to combine
// multiple validations and accumulate all errors:
//
//	type User struct {
//		Name  string
//		Email string
//		Age   int
//	}
//
//	validateName := func(s string) validation.Validation[string] {
//		if len(s) > 0 {
//			return validation.Success(s)
//		}
//		return validation.Failures[string](/* error */)
//	}
//
//	// Combine validations - all errors will be collected
//	result := validation.Ap(validation.Ap(validation.Ap(
//		validation.Of(func(name string) func(email string) func(age int) User {
//			return func(email string) func(age int) User {
//				return func(age int) User {
//					return User{name, email, age}
//				}
//			}
//		}),
//	)(validateName("")))(validateEmail("")))(validateAge(-1))
//
// # Error Formatting
//
// ValidationError implements custom formatting for detailed error messages:
//
//	err := &ValidationError{
//		Value:    "abc",
//		Context:  []ContextEntry{{Key: "user"}, {Key: "age"}},
//		Message:  "expected integer",
//	}
//
//	fmt.Printf("%v", err)   // at user.age: expected integer
//	fmt.Printf("%+v", err)  // at user.age: expected integer
//	                        //   value: "abc"
//
// # Monoid Operations
//
// The package provides monoid instances for combining validations:
//
//	// Combine validation results
//	m := validation.ApplicativeMonoid(stringMonoid)
//	combined := m.Concat(validation.Success("hello"), validation.Success(" world"))
//	// Result: Success("hello world")
//
// # Integration
//
// This package integrates with:
//   - either: Validation is built on Either for error handling
//   - array: For collecting multiple errors
//   - monoid: For combining validation results
//   - reader: For context-dependent validation operations
package validation
