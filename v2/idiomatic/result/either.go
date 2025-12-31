// Copyright (c) 2023 - 2025 IBM Corp.
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

package result

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/option"
)

// Of constructs a Right value containing the given value.
// This is the monadic return/pure operation for Either.
// Equivalent to [Right].
//
// Example:
//
//	result := either.Of[error](42) // Right(42)
//
//go:inline
func Of[A any](a A) (A, error) {
	return a, nil
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
func Ap[B, A any](fa A, faerr error) Operator[func(A) B, B] {
	return func(fab func(A) B, faberr error) (B, error) {
		if faberr != nil {
			return Left[B](faberr)
		}
		if faerr != nil {
			return Left[B](faerr)
		}
		return Of(fab(fa))
	}
}

// BiMap is the curried version of [MonadBiMap].
// Maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[A, B any](f Endomorphism[error], g func(a A) B) Operator[A, B] {
	return func(a A, err error) (B, error) {
		if err != nil {
			return Left[B](f(err))
		}
		return Of(g(a))
	}
}

// MapTo is the curried version of [MonadMapTo].
func MapTo[A, B any](b B) Operator[A, B] {
	return func(_ A, err error) (B, error) {
		return b, err
	}
}

// Map is the curried version of [MonadMap].
// Transforms the Right value using the provided function.
func Map[A, B any](f func(A) B) Operator[A, B] {
	return func(a A, err error) (B, error) {
		if err != nil {
			return Left[B](err)
		}
		return Of(f(a))
	}
}

// MapLeft is the curried version of [MonadMapLeft].
// Applies a mapping function to the Left (error) channel.
func MapLeft[A any](f Endomorphism[error]) Operator[A, A] {
	return func(a A, err error) (A, error) {
		if err != nil {
			return Left[A](f(err))
		}
		return Of(a)
	}
}

// ChainOptionK is the curried version of [MonadChainOptionK].
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B] {
	return func(f func(A) (B, bool)) Operator[A, B] {
		return func(a A, err error) (B, error) {
			if err != nil {
				return Left[B](err)
			}
			if b, ok := f(a); ok {
				return Of(b)
			}
			return Left[B](onNone())
		}
	}
}

// ChainTo is the curried version of [MonadChainTo].
func ChainTo[A, B any](b B, berr error) Operator[A, B] {
	if berr != nil {
		return func(_ A, _ error) (B, error) {
			return Left[B](berr)
		}
	}
	return func(a A, aerr error) (B, error) {
		if aerr != nil {
			return Left[B](aerr)
		}
		return Of(b)
	}
}

// Chain is the curried version of [MonadChain].
// Sequences two computations where the second depends on the first.
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return func(a A, aerr error) (B, error) {
		if aerr != nil {
			return Left[B](aerr)
		}
		return f(a)
	}
}

// ChainFirst is the curried version of [MonadChainFirst].
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return func(a A, aerr error) (A, error) {
		if aerr != nil {
			return Left[A](aerr)
		}
		_, berr := f(a)
		return a, berr
	}
}

// FromOption converts an Option to an Either, using the provided function to generate a Left value for None.
//
// Example:
//
//	opt := option.Some(42)
//	result := either.FromOption[int](func() error { return errors.New("none") })(opt) // Right(42)
func FromOption[A any](onNone func() error) func(A, bool) (A, error) {
	return func(a A, aok bool) (A, error) {
		if !aok {
			return Left[A](onNone())
		}
		return Of(a)
	}
}

// ToOption converts an Either to an Option, discarding the Left value.
//
// Example:
//
//	result := either.ToOption(either.Right[error](42)) // Some(42)
//	result := either.ToOption(either.Left[int](errors.New("err"))) // None
//
//go:inline
func ToOption[A any](a A, aerr error) (A, bool) {
	return a, aerr == nil
}

// FromError creates an Either from a function that may return an error.
//
// Example:
//
//	validate := func(x int) error {
//	    if x < 0 { return errors.New("negative") }
//	    return nil
//	}
//	toEither := either.FromError(validate)
//	result := toEither(42) // Right(42)
func FromError[A any](f func(a A) error) Kleisli[A, A] {
	return func(a A) (A, error) {
		return a, f(a)
	}
}

// ToError converts an Either[error, A] to an error, returning nil for Right values.
//
// Example:
//
//	err := either.ToError(either.Left[int](errors.New("fail"))) // error
//	err := either.ToError(either.Right[error](42)) // nil
func ToError[A any](_ A, err error) error {
	return err
}

// Fold is the curried version of [MonadFold].
// Extracts the value from an Either by providing handlers for both cases.
//
// Example:
//
//	result := either.Fold(
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Value: %d", n) },
//	)(either.Right[error](42)) // "Value: 42"
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(A, error) B {
	return func(a A, aerr error) B {
		if aerr != nil {
			return onLeft(aerr)
		}
		return onRight(a)
	}
}

// FromPredicate creates an Either based on a predicate.
// If the predicate returns true, creates a Right; otherwise creates a Left using onFalse.
//
// Example:
//
//	isPositive := either.FromPredicate(
//	    N.MoreThan(0),
//	    func(x int) error { return errors.New("not positive") },
//	)
//	result := isPositive(42) // Right(42)
//	result := isPositive(-1) // Left(error)
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return func(a A) (A, error) {
		if pred(a) {
			return Right(a)
		}
		return Left[A](onFalse(a))
	}
}

// FromNillable creates an Either from a pointer, using the provided error for nil pointers.
//
// Example:
//
//	var ptr *int = nil
//	result := either.FromNillable[int](errors.New("nil"))(ptr) // Left(error)
//	val := 42
//	result := either.FromNillable[int](errors.New("nil"))(&val) // Right(&42)
func FromNillable[A any](e error) Kleisli[*A, *A] {
	return func(a *A) (*A, error) {
		if F.IsNil(a) {
			return Left[*A](e)
		}
		return Of(a)
	}
}

// GetOrElse extracts the Right value or computes a default from the Left value.
//
// Example:
//
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Right[error](42)) // 42
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Left[int](err)) // 0
func GetOrElse[A any](onLeft func(error) A) func(A, error) A {
	return func(a A, err error) A {
		if err != nil {
			return onLeft(err)
		}
		return a
	}
}

// Reduce folds an Either into a single value using a reducer function.
// Returns the initial value for Left, or applies the reducer to the Right value.
func Reduce[A, B any](f func(B, A) B, initial B) func(A, error) B {
	return func(a A, err error) B {
		if err != nil {
			return initial
		}
		return f(initial, a)
	}
}

// Alt provides an alternative Either if the first is Left.
//
// Example:
//
//	alternative := either.Alt[error](func() either.Either[error, int] {
//	    return either.Right[error](99)
//	})
//	result := alternative(either.Left[int](errors.New("fail"))) // Right(99)
func Alt[A any](that func() (A, error)) Operator[A, A] {
	return func(a A, err error) (A, error) {
		if err != nil {
			return that()
		}
		return Of(a)
	}
}

// OrElse recovers from a Left by providing an alternative computation.
//
// Example:
//
//	recover := either.OrElse(func(err error) either.Either[error, int] {
//	    return either.Right[error](0) // default value
//	})
//	result := recover(either.Left[int](errors.New("fail"))) // Right(0)
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A] {
	return func(a A, err error) (A, error) {
		if err != nil {
			return onLeft(err)
		}
		return Of(a)
	}
}

// ToType attempts to convert an any value to a specific type, returning Either.
//
// Example:
//
//	convert := either.ToType[int](func(v any) error {
//	    return fmt.Errorf("cannot convert %v to int", v)
//	})
//	result := convert(42) // Right(42)
//	result := convert("string") // Left(error)
func ToType[A any](onError func(any) error) Kleisli[any, A] {
	return func(x any) (A, error) {
		if a, ok := x.(A); ok {
			return Of(a)
		}
		return Left[A](onError(x))
	}
}

// Memoize returns the Either unchanged (Either values are already memoized).
func Memoize[A any](a A, err error) (A, error) {
	return a, err
}

// Flap is the curried version of [MonadFlap].
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return func(fab func(A) B, faberr error) (B, error) {
		if faberr != nil {
			return Left[B](faberr)
		}
		return Of(fab(a))
	}
}
