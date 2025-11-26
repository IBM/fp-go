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
	"github.com/IBM/fp-go/v2/internal/monad"
)

// eitherMonad is the internal implementation of the Monad type class for Either.
// It extends eitherApplicative by adding the Chain operation for sequential composition.
type eitherMonad[E, A, B any] struct {
	eitherApplicative[E, A, B]
	fchain func(Kleisli[E, A, B]) Operator[E, A, B]
}

// Chain sequences dependent computations, failing fast on the first Left.
func (o *eitherMonad[E, A, B]) Chain(f Kleisli[E, A, B]) Operator[E, A, B] {
	return o.fchain(f)
}

// Monad creates a lawful Monad instance for Either with fail-fast error handling.
//
// A monad combines the capabilities of four type classes:
//   - Functor (Map): transform the Right value
//   - Pointed (Of): lift a pure value into a Right
//   - Applicative (Ap): apply wrapped functions (fails fast on first Left)
//   - Chainable (Chain): sequence dependent computations (fails fast on first Left)
//
// The Either monad is left-biased and fails fast: once a Left is encountered,
// no further computations are performed and the Left is propagated immediately.
// This makes it ideal for error handling where you want to stop at the first error.
//
// This implementation satisfies all monad laws:
//
// Monad Laws:
//   - Left Identity: Chain(f)(Of(a)) == f(a)
//   - Right Identity: Chain(Of)(m) == m
//   - Associativity: Chain(g)(Chain(f)(m)) == Chain(x => Chain(g)(f(x)))(m)
//
// Additionally, it satisfies all prerequisite laws from Functor, Apply, and Applicative.
//
// Relationship to Applicative:
//
// This Monad uses the standard fail-fast Applicative (see Applicative function).
// In a lawful monad, Ap can be derived from Chain and Of:
//
//	Ap(fa)(ff) == Chain(f => Chain(a => Of(f(a)))(fa))(ff)
//
// The Either monad satisfies this property, making it a true lawful monad.
//
// When to use Monad vs Applicative:
//   - Use Monad when you need sequential dependent operations (Chain)
//   - Use Applicative when you only need independent operations (Ap, Map)
//   - Both fail fast on the first error
//
// When to use Monad vs ApplicativeV:
//   - Use Monad for sequential error handling (fail-fast)
//   - Use ApplicativeV for parallel validation (error accumulation)
//   - Note: There is no "MonadV" because Chain inherently fails fast
//
// Example - Sequential Dependent Operations:
//
//	m := either.Monad[error, int, string]()
//
//	// Chain allows each step to depend on the previous result
//	result := m.Chain(func(x int) either.Either[error, string] {
//	    if x > 0 {
//	        return either.Right[error](strconv.Itoa(x))
//	    }
//	    return either.Left[string](errors.New("value must be positive"))
//	})(either.Right[error](42))
//	// result is Right("42")
//
//	// Fails fast on first error
//	result2 := m.Chain(func(x int) either.Either[error, string] {
//	    return either.Right[error](strconv.Itoa(x))
//	})(either.Left[int](errors.New("initial error")))
//	// result2 is Left("initial error") - Chain never executes
//
// Example - Combining with Applicative operations:
//
//	m := either.Monad[error, int, int]()
//
//	// Map transforms the value
//	value := m.Map(N.Mul(2))(either.Right[error](21))
//	// value is Right(42)
//
//	// Ap applies wrapped functions (also fails fast)
//	fn := either.Right[error](N.Add(1))
//	result := m.Ap(value)(fn)
//	// result is Right(43)
//
// Example - Real-world usage with error handling:
//
//	m := either.Monad[error, User, SavedUser]()
//
//	// Pipeline of operations that can fail
//	result := m.Chain(func(user User) either.Either[error, SavedUser] {
//	    // Save to database
//	    return saveToDatabase(user)
//	})(m.Chain(func(user User) either.Either[error, User] {
//	    // Validate user
//	    return validateUser(user)
//	})(either.Right[error](inputUser)))
//
//	// If any step fails, the error propagates immediately
//
// Type Parameters:
//   - E: The error type (Left value)
//   - A: The input value type (Right value)
//   - B: The output value type after transformation
func Monad[E, A, B any]() monad.Monad[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]] {
	return &eitherMonad[E, A, B]{
		eitherApplicative[E, A, B]{
			Of[E, A],
			Map[E, A, B],
			Ap[B, E, A],
		},
		Chain[E, A, B],
	}
}
