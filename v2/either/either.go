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

// package either implements the Either monad
//
// A data type that can be of either of two types but not both. This is
// typically used to carry an error or a return value
package either

import (
	"github.com/IBM/fp-go/v2/internal/common"
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
func Of[E, A any](value A) Either[E, A] {
	return common.EitherOf[E](value)
}

// FromIO executes an IO operation and wraps the result in a Right value.
// This is useful for lifting pure IO operations into the Either context.
//
// Example:
//
//	getValue := lazy.Of(42)
//	result := either.FromIO[error](getValue) // Right(42)
//
// go: inline
func FromIO[E any, IO ~func() A, A any](f IO) Either[E, A] {
	return common.EitherFromIO[E, IO](f)
}

// MonadAp applies a function wrapped in Either to a value wrapped in Either.
// If either the function or the value is Left, returns Left.
// This is the applicative apply operation.
//
// Example:
//
//	fab := either.Right[error](N.Mul(2))
//	fa := either.Right[error](21)
//	result := either.MonadAp(fab, fa) // Right(42)
func MonadAp[B, E, A any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	return common.EitherMonadAp[B](fab, fa)
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
//
//go:inline
func Ap[B, E, A any](fa Either[E, A]) Operator[E, func(A) B, B] {
	return common.EitherAp[B](fa)
}

// MonadMap transforms the Right value using the provided function.
// If the Either is Left, returns Left unchanged.
// This is the functor map operation.
//
// Example:
//
//	result := either.MonadMap(
//	    either.Right[error](21),
//	    N.Mul(2),
//	) // Right(42)
//
//go:inline
func MonadMap[E, A, B any](fa Either[E, A], f func(A) B) Either[E, B] {
	return common.EitherMonadMap(fa, f)
}

// MonadBiMap applies two functions: one to transform a Left value, another to transform a Right value.
// This allows transforming both channels of the Either simultaneously.
//
// Example:
//
//	result := either.MonadBiMap(
//	    either.Left[int](errors.New("error")),
//	    error.Error,
//	    func(n int) string { return fmt.Sprint(n) },
//	) // Left("error")
func MonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B] {
	return common.EitherMonadBiMap(fa, f, g)
}

// BiMap is the curried version of [MonadBiMap].
// Maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B] {
	return common.EitherBiMap(f, g)
}

// MonadMapTo replaces the Right value with a constant value.
// If the Either is Left, returns Left unchanged.
//
// Example:
//
//	result := either.MonadMapTo(either.Right[error](21), "success") // Right("success")
func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	return common.EitherMonadMapTo(fa, b)
}

// MapTo is the curried version of [MonadMapTo].
func MapTo[E, A, B any](b B) Operator[E, A, B] {
	return common.EitherMapTo[E, A](b)
}

// MonadMapLeft applies a transformation function to the Left (error) value.
// If the Either is Right, returns Right unchanged.
//
// Example:
//
//	result := either.MonadMapLeft(
//	    either.Left[int](errors.New("error")),
//	    error.Error,
//	) // Left("error")
func MonadMapLeft[E1, A, E2 any](fa Either[E1, A], f func(E1) E2) Either[E2, A] {
	return common.EitherMonadMapLeft(fa, f)
}

// Map is the curried version of [MonadMap].
// Transforms the Right value using the provided function.
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return common.EitherMap[E](f)
}

// MapLeft is the curried version of [MonadMapLeft].
// Applies a mapping function to the Left (error) channel.
func MapLeft[A, E1, E2 any](f func(E1) E2) func(fa Either[E1, A]) Either[E2, A] {
	return common.EitherMapLeft[A](f)
}

// MonadChain sequences two computations, where the second depends on the result of the first.
// If the first Either is Left, returns Left without executing the second computation.
// This is the monadic bind operation (also known as flatMap).
//
// Example:
//
//	result := either.MonadChain(
//	    either.Right[error](21),
//	    func(x int) either.Either[error, int] {
//	        return either.Right[error](x * 2)
//	    },
//	) // Right(42)
//
//go:inline
func MonadChain[E, A, B any](fa Either[E, A], f Kleisli[E, A, B]) Either[E, B] {
	return common.EitherMonadChain(fa, f)
}

// MonadChainLeft sequences a computation on the Left (error) value, allowing error recovery or transformation.
// If the Either is Left, applies the provided function to the error value, which returns a new Either.
// If the Either is Right, returns the Right value unchanged with the new error type.
//
// This is the dual of [MonadChain] - while MonadChain operates on Right values (success),
// MonadChainLeft operates on Left values (errors). It's useful for error recovery, error transformation,
// or chaining alternative computations when an error occurs.
//
// Note: MonadChainLeft is identical to [OrElse] - both provide the same functionality for error recovery.
//
// The error type can be transformed from EA to EB, allowing flexible error type conversions.
//
// Example:
//
//	// Error recovery: convert specific errors to success
//	result := either.MonadChainLeft(
//	    either.Left[int](errors.New("not found")),
//	    func(err error) either.Either[string, int] {
//	        if err.Error() == "not found" {
//	            return either.Right[string](0) // default value
//	        }
//	        return either.Left[int](err.Error()) // transform error
//	    },
//	) // Right(0)
//
//	// Error transformation: change error type
//	result := either.MonadChainLeft(
//	    either.Left[int](404),
//	    func(code int) either.Either[string, int] {
//	        return either.Left[int](fmt.Sprintf("Error code: %d", code))
//	    },
//	) // Left("Error code: 404")
//
//	// Right values pass through unchanged
//	result := either.MonadChainLeft(
//	    either.Right[error](42),
//	    func(err error) either.Either[string, int] {
//	        return either.Left[int]("error")
//	    },
//	) // Right(42)
//
//go:inline
func MonadChainLeft[EA, EB, A any](fa Either[EA, A], f Kleisli[EB, EA, A]) Either[EB, A] {
	return common.EitherMonadChainLeft(fa, f)
}

// ChainLeft is the curried version of [MonadChainLeft].
// Returns a function that sequences a computation on the Left (error) value.
//
// Note: ChainLeft is identical to [OrElse] - both provide the same functionality for error recovery.
//
// This is useful for creating reusable error handlers or transformers that can be
// composed with other Either operations using pipes or function composition.
//
// Example:
//
//	// Create a reusable error handler
//	handleNotFound := either.ChainLeft[error, string](func(err error) either.Either[string, int] {
//	    if err.Error() == "not found" {
//	        return either.Right[string](0)
//	    }
//	    return either.Left[int](err.Error())
//	})
//
//	// Use in a pipeline
//	result := F.Pipe1(
//	    either.Left[int](errors.New("not found")),
//	    handleNotFound,
//	) // Right(0)
//
//go:inline
func ChainLeft[EA, EB, A any](f Kleisli[EB, EA, A]) Kleisli[EB, Either[EA, A], A] {
	return common.EitherChainLeft(f)
}

// MonadChainFirst executes a side-effect computation but returns the original value.
// Useful for performing actions (like logging) without changing the value.
//
// Example:
//
//	result := either.MonadChainFirst(
//	    either.Right[error](42),
//	    func(x int) either.Either[error, string] {
//	        fmt.Println(x) // side effect
//	        return either.Right[error]("logged")
//	    },
//	) // Right(42) - original value preserved
func MonadChainFirst[E, A, B any](ma Either[E, A], f Kleisli[E, A, B]) Either[E, A] {
	return common.EitherMonadChainFirst(ma, f)
}

// MonadChainTo ignores the first Either and returns the second.
// Useful for sequencing operations where you don't need the first result.
func MonadChainTo[A, E, B any](ma Either[E, A], mb Either[E, B]) Either[E, B] {
	return common.EitherMonadChainTo[A](ma, mb)
}

// MonadChainOptionK chains a function that returns an Option, converting None to Left.
//
// Example:
//
//	result := either.MonadChainOptionK(
//	    func() error { return errors.New("not found") },
//	    either.Right[error](42),
//	    func(x int) option.Option[string] {
//	        if x > 0 { return option.Some("positive") }
//	        return option.None[string]()
//	    },
//	) // Right("positive")
func MonadChainOptionK[A, B, E any](onNone func() E, ma Either[E, A], f func(A) Option[B]) Either[E, B] {
	return common.EitherMonadChainOptionK[A, B](onNone, ma, f)
}

// ChainOptionK is the curried version of [MonadChainOptionK].
func ChainOptionK[A, B, E any](onNone func() E) func(func(A) Option[B]) Operator[E, A, B] {
	return common.EitherChainOptionK[A, B](onNone)
}

// ChainTo is the curried version of [MonadChainTo].
func ChainTo[A, E, B any](mb Either[E, B]) Operator[E, A, B] {
	return common.EitherChainTo[A](mb)
}

// Chain is the curried version of [MonadChain].
// Sequences two computations where the second depends on the first.
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return common.EitherChain(f)
}

// ChainFirst is the curried version of [MonadChainFirst].
func ChainFirst[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A] {
	return common.EitherChainFirst(f)
}

// Flatten removes one level of nesting from a nested Either.
//
// Example:
//
//	nested := either.Right[error](either.Right[error](42))
//	result := either.Flatten(nested) // Right(42)
func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A] {
	return common.EitherFlatten(mma)
}

// TryCatch converts a (value, error) tuple into an Either, applying a transformation to the error.
//
// Example:
//
//	result := either.TryCatch(
//	    42, nil,
//	    func(err error) string { return err.Error() },
//	) // Right(42)
func TryCatch[FE func(error) E, E, A any](val A, err error, onThrow FE) Either[E, A] {
	return common.EitherTryCatch[FE](val, err, onThrow)
}

// TryCatchError is a specialized version of [TryCatch] for error types.
// Converts a (value, error) tuple into Either[error, A].
//
// Example:
//
//	result := either.TryCatchError(42, nil) // Right(42)
//	result := either.TryCatchError(0, errors.New("fail")) // Left(error)
func TryCatchError[A any](val A, err error) Either[error, A] {
	return common.EitherTryCatchError[A](val, err)
}

// Sequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R] {
	return common.EitherSequence2(f)
}

// Sequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R] {
	return common.EitherSequence3(f)
}

// FromOption converts an Option to an Either, using the provided function to generate a Left value for None.
//
// Example:
//
//	opt := option.Some(42)
//	result := either.FromOption[int](func() error { return errors.New("none") })(opt) // Right(42)
func FromOption[A, E any](onNone func() E) func(Option[A]) Either[E, A] {
	return common.EitherFromOption[A](onNone)
}

// ToOption converts an Either to an Option, discarding the Left value.
//
// Example:
//
//	result := either.ToOption(either.Right[error](42)) // Some(42)
//	result := either.ToOption(either.Left[int](errors.New("err"))) // None
//
//go:inline
func ToOption[E, A any](ma Either[E, A]) Option[A] {
	return common.EitherToOption(ma)
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
func FromError[A any](f func(a A) error) func(A) Either[error, A] {
	return common.EitherFromError(f)
}

// ToError converts an Either[error, A] to an error, returning nil for Right values.
//
// Example:
//
//	err := either.ToError(either.Left[int](errors.New("fail"))) // error
//	err := either.ToError(either.Right[error](42)) // nil
func ToError[A any](e Either[error, A]) error {
	return common.EitherToError(e)
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
//
//go:inline
func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B {
	return common.EitherFold(onLeft, onRight)
}

// UnwrapError converts an Either[error, A] into the idiomatic Go tuple (A, error).
//
// Example:
//
//	val, err := either.UnwrapError(either.Right[error](42)) // 42, nil
//	val, err := either.UnwrapError(either.Left[int](errors.New("fail"))) // zero, error
//
//go:inline
func UnwrapError[A any](ma Either[error, A]) (A, error) {
	return common.EitherUnwrapError(ma)
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
func FromPredicate[E, A any](pred Predicate[A], onFalse func(A) E) Kleisli[E, A, A] {
	return common.EitherFromPredicate(pred, onFalse)
}

// FromNillable creates an Either from a pointer, using the provided error for nil pointers.
//
// Example:
//
//	var ptr *int = nil
//	result := either.FromNillable[int](errors.New("nil"))(ptr) // Left(error)
//	val := 42
//	result := either.FromNillable[int](errors.New("nil"))(&val) // Right(&42)
func FromNillable[A, E any](e E) Kleisli[E, *A, *A] {
	return common.EitherFromNillable[A](e)
}

// GetOrElse extracts the Right value or computes a default from the Left value.
//
// Example:
//
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Right[error](42)) // 42
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Left[int](err)) // 0
func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A {
	return common.EitherGetOrElse(onLeft)
}

// Reduce folds an Either into a single value using a reducer function.
// Returns the initial value for Left, or applies the reducer to the Right value.
func Reduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B {
	return common.EitherReduce[E](f, initial)
}

// AltW provides an alternative Either if the first is Left, allowing different error types.
// The 'W' suffix indicates "widening" of the error type.
//
// Example:
//
//	alternative := either.AltW[error, string](func() either.Either[string, int] {
//	    return either.Right[string](99)
//	})
//	result := alternative(either.Left[int](errors.New("fail"))) // Right(99)
func AltW[E, E1, A any](that Lazy[Either[E1, A]]) Kleisli[E1, Either[E, A], A] {
	return common.EitherAltW[E](that)
}

// Alt provides an alternative Either if the first is Left.
//
// Example:
//
//	alternative := either.Alt[error](func() either.Either[error, int] {
//	    return either.Right[error](99)
//	})
//	result := alternative(either.Left[int](errors.New("fail"))) // Right(99)
func Alt[E, A any](that Lazy[Either[E, A]]) Operator[E, A, A] {
	return common.EitherAlt(that)
}

// OrElse recovers from a Left (error) by providing an alternative computation.
// If the Either is Right, it returns the value unchanged.
// If the Either is Left, it applies the provided function to the error value,
// which returns a new Either that replaces the original.
//
// Note: OrElse is identical to [ChainLeft] - both provide the same functionality for error recovery.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations.
// The error type can be widened from E1 to E2, allowing transformation of error types.
//
// Example:
//
//	// Recover from specific errors with fallback values
//	recover := either.OrElse(func(err error) either.Either[error, int] {
//	    if err.Error() == "not found" {
//	        return either.Right[error](0) // default value
//	    }
//	    return either.Left[int](err) // propagate other errors
//	})
//	result := recover(either.Left[int](errors.New("not found"))) // Right(0)
//	result := recover(either.Right[error](42)) // Right(42) - unchanged
//
//go:inline
func OrElse[E1, E2, A any](onLeft Kleisli[E2, E1, A]) Kleisli[E2, Either[E1, A], A] {
	return common.EitherOrElse(onLeft)
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
func ToType[A, E any](onError func(any) E) func(any) Either[E, A] {
	return common.EitherToType[A](onError)
}

// Memoize returns the Either unchanged (Either values are already memoized).
func Memoize[E, A any](val Either[E, A]) Either[E, A] {
	return common.EitherMemoize(val)
}

// MonadSequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	return common.EitherMonadSequence2(e1, e2, f)
}

// MonadSequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	return common.EitherMonadSequence3(e1, e2, e3, f)
}

// Swap exchanges the Left and Right type parameters.
//
// Example:
//
//	result := either.Swap(either.Right[error](42)) // Left(42)
//	result := either.Swap(either.Left[int](errors.New("err"))) // Right(error)
//
//go:inline
func Swap[E, A any](val Either[E, A]) Either[A, E] {
	return common.EitherSwap(val)
}

// MonadFlap applies a value to a function wrapped in Either.
// This is the reverse of [MonadAp].
func MonadFlap[E, B, A any](fab Either[E, func(A) B], a A) Either[E, B] {
	return common.EitherMonadFlap[E, B](fab, a)
}

// Flap is the curried version of [MonadFlap].
func Flap[E, B, A any](a A) Operator[E, func(A) B, B] {
	return common.EitherFlap[E, B](a)
}

// MonadAlt provides an alternative Either if the first is Left.
// This is the monadic version of [Alt].
func MonadAlt[E, A any](fa Either[E, A], that Lazy[Either[E, A]]) Either[E, A] {
	return common.EitherMonadAlt(fa, that)
}

// Zero returns the zero value of an [Either], which is a Right containing the zero value of type A.
// This function is useful as an identity element in monoid operations or for creating an empty Either
// in a Right state.
//
// The returned Either is always a Right value containing the zero value of type A. For reference types
// (pointers, slices, maps, channels, functions, interfaces), the zero value is nil. For value types
// (numbers, booleans, structs), it's the type's zero value.
//
// Important: Zero() returns the same value as the default initialization of Either[E, A].
// When you declare `var e Either[E, A]` without initialization, it has the same value as Zero[E, A]().
//
// Note: This differs from creating a Left value, which would represent an error or failure state.
// Zero always produces a successful (Right) state with a zero value.
//
// Example:
//
//	// Zero Either with int value
//	e1 := either.Zero[error, int]()  // Right(0)
//
//	// Zero Either with string value
//	e2 := either.Zero[error, string]()  // Right("")
//
//	// Zero Either with pointer type
//	e3 := either.Zero[error, *int]()  // Right(nil)
//
//	// Zero equals default initialization
//	var defaultInit Either[error, int]
//	zero := either.Zero[error, int]()
//	assert.Equal(t, defaultInit, zero) // true
//
//	// Verify it's a Right value
//	e := either.Zero[error, int]()
//	assert.True(t, either.IsRight(e))  // true
//	assert.False(t, either.IsLeft(e))  // false
func Zero[E, A any]() Either[E, A] {
	return common.EitherZero[E, A]()
}
