// Copyright (c) 2024 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package either

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// eitherApplicative is the internal implementation of the Applicative type class for Either.
// It provides the basic applicative operations: Of (lift), Map (transform), and Ap (apply).
type eitherApplicative[E, A, B any] struct {
	fof  func(a A) Either[E, A]
	fmap func(func(A) B) Operator[E, A, B]
	fap  func(Either[E, A]) Operator[E, func(A) B, B]
}

// Of lifts a pure value into a Right context.
func (o *eitherApplicative[E, A, B]) Of(a A) Either[E, A] {
	return o.fof(a)
}

// Map applies a transformation function to the Right value, preserving Left values.
func (o *eitherApplicative[E, A, B]) Map(f func(A) B) Operator[E, A, B] {
	return o.fmap(f)
}

// Ap applies a wrapped function to a wrapped value.
// The behavior depends on which Ap implementation is used (fail-fast or validation).
func (o *eitherApplicative[E, A, B]) Ap(fa Either[E, A]) Operator[E, func(A) B, B] {
	return o.fap(fa)
}

// Applicative creates a standard Applicative instance for Either with fail-fast error handling.
//
// This returns a lawful Applicative that satisfies all applicative laws:
//   - Identity: Ap(Of(identity))(v) == v
//   - Homomorphism: Ap(Of(f))(Of(x)) == Of(f(x))
//   - Interchange: Ap(Of(f))(u) == Ap(Map(f => f(y))(u))(Of(y))
//   - Composition: Ap(Ap(Map(compose)(f))(g))(x) == Ap(f)(Ap(g)(x))
//
// The Applicative operations behave as follows:
//   - Of: lifts a value into Right
//   - Map: transforms Right values, preserves Left (standard functor)
//   - Ap: fails fast - if either operand is Left, returns the first Left encountered
//
// This is the standard Either applicative that stops at the first error, making it
// suitable for computations where you want to short-circuit on failure.
//
// Example - Fail-Fast Behavior:
//
//	app := either.Applicative[error, int, string]()
//
//	// Both succeed - function application works
//	value := either.Right[error](42)
//	fn := either.Right[error](strconv.Itoa)
//	result := app.Ap(value)(fn)
//	// result is Right("42")
//
//	// First error stops computation
//	err1 := either.Left[func(int) string](errors.New("error 1"))
//	err2 := either.Left[int](errors.New("error 2"))
//	result2 := app.Ap(err2)(err1)
//	// result2 is Left(error 1) - only first error is returned
//
// Type Parameters:
//   - E: The error type (Left value)
//   - A: The input value type (Right value)
//   - B: The output value type after transformation
func Applicative[E, A, B any]() applicative.Applicative[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]] {
	return &eitherApplicative[E, A, B]{
		Of[E, A],
		Map[E, A, B],
		Ap[B, E, A],
	}
}

// ApplicativeV creates an Applicative with validation-style error accumulation.
//
// This returns a lawful Applicative that accumulates errors using a Semigroup when
// combining independent computations with Ap. This is the "validation" pattern commonly
// used for form validation, configuration validation, and parallel error collection.
//
// The returned instance satisfies all applicative laws:
//   - Identity: Ap(Of(identity))(v) == v
//   - Homomorphism: Ap(Of(f))(Of(x)) == Of(f(x))
//   - Interchange: Ap(Of(f))(u) == Ap(Map(f => f(y))(u))(Of(y))
//   - Composition: Ap(Ap(Map(compose)(f))(g))(x) == Ap(f)(Ap(g)(x))
//
// Key behaviors:
//   - Of: lifts a value into Right
//   - Map: transforms Right values, preserves Left (standard functor)
//   - Ap: when both operands are Left, combines errors using the Semigroup
//
// Comparison with standard Applicative:
//   - Applicative: Ap fails fast (returns first error)
//   - ApplicativeV: Ap accumulates errors (combines all errors via Semigroup)
//
// Use cases:
//   - Form validation: collect all validation errors at once
//   - Configuration validation: report all configuration problems
//   - Parallel independent checks: accumulate all failures
//
// Example - Error Accumulation for Form Validation:
//
//	type ValidationErrors []string
//
//	// Define how to combine error lists
//	sg := semigroup.MakeSemigroup(func(a, b ValidationErrors) ValidationErrors {
//	    return append(append(ValidationErrors{}, a...), b...)
//	})
//
//	app := either.ApplicativeV[ValidationErrors, User, User](sg)
//
//	// Validate multiple fields independently
//	validateName := func(name string) Either[ValidationErrors, string] {
//	    if len(name) < 3 {
//	        return Left[string](ValidationErrors{"Name must be at least 3 characters"})
//	    }
//	    return Right[ValidationErrors](name)
//	}
//
//	validateAge := func(age int) Either[ValidationErrors, int] {
//	    if age < 18 {
//	        return Left[int](ValidationErrors{"Must be 18 or older"})
//	    }
//	    return Right[ValidationErrors](age)
//	}
//
//	validateEmail := func(email string) Either[ValidationErrors, string] {
//	    if !strings.Contains(email, "@") {
//	        return Left[string](ValidationErrors{"Invalid email format"})
//	    }
//	    return Right[ValidationErrors](email)
//	}
//
//	// Create a constructor function lifted into Either
//	makeUser := func(name string) func(int) func(string) User {
//	    return func(age int) func(string) User {
//	        return func(email string) User {
//	            return User{Name: name, Age: age, Email: email}
//	        }
//	    }
//	}
//
//	// Apply validations - all errors are collected
//	name := validateName("ab")          // Left: name too short
//	age := validateAge(16)              // Left: age too low
//	email := validateEmail("invalid")    // Left: invalid email
//
//	// Combine all validations using ApV
//	result := app.Ap(name)(
//	    app.Ap(age)(
//	        app.Ap(email)(
//	            app.Of(makeUser),
//	        ),
//	    ),
//	)
//	// result is Left(ValidationErrors{
//	//     "Name must be at least 3 characters",
//	//     "Must be 18 or older",
//	//     "Invalid email format"
//	// })
//	// All three errors are collected!
//
// Type Parameters:
//   - E: The error type that must have a Semigroup for combining errors
//   - A: The input value type (Right value)
//   - B: The output value type after transformation
//   - sg: Semigroup instance for combining Left values when both operands of Ap are Left
func ApplicativeV[E, A, B any](sg S.Semigroup[E]) applicative.Applicative[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]] {
	return &eitherApplicative[E, A, B]{
		Of[E, A],
		Map[E, A, B],
		ApV[B, A](sg),
	}
}
