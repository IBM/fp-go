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

// package result implements the Either monad
//
// A data type that can be of either of two types but not both. This is
// typically used to carry an error or a return value
package result

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/option"
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
func Of[A any](value A) Result[A] {
	return either.Of[error](value)
}

// FromIO executes an IO operation and wraps the result in a Right value.
// This is useful for lifting pure IO operations into the Either context.
//
// Example:
//
//	getValue := func() int { return 42 }
//	result := either.FromIO[error](getValue) // Right(42)
//
//go:inline
func FromIO[IO ~func() A, A any](f IO) Result[A] {
	return either.FromIO[error](f)
}

// MonadAp applies a function wrapped in Either to a value wrapped in Either.
// If either the function or the value is Left, returns Left.
// This is the applicative apply operation.
//
// Example:
//
//	fab := either.Right[error](func(x int) int { return x * 2 })
//	fa := either.Right[error](21)
//	result := either.MonadAp(fab, fa) // Right(42)
//
//go:inline
func MonadAp[B, A any](fab Result[func(a A) B], fa Result[A]) Result[B] {
	return either.MonadAp(fab, fa)
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
//
//go:inline
func Ap[B, A any](fa Result[A]) Operator[func(A) B, B] {
	return either.Ap[B](fa)
}

// MonadMap transforms the Right value using the provided function.
// If the Either is Left, returns Left unchanged.
// This is the functor map operation.
//
// Example:
//
//	result := either.MonadMap(
//	    either.Right[error](21),
//	    func(x int) int { return x * 2 },
//	) // Right(42)
//
//go:inline
func MonadMap[A, B any](fa Result[A], f func(a A) B) Result[B] {
	return either.MonadMap(fa, f)
}

// MonadBiMap applies two functions: one to transform a Left value, another to transform a Right value.
// This allows transforming both channels of the Either simultaneously.
//
// Example:
//
//	result := either.MonadBiMap(
//	    either.Left[int](errors.New("error")),
//	    func(e error) string { return e.Error() },
//	    func(n int) string { return fmt.Sprint(n) },
//	) // Left("error")
//
//go:inline
func MonadBiMap[E, A, B any](fa Result[A], f func(error) E, g func(a A) B) Either[E, B] {
	return either.MonadBiMap(fa, f, g)
}

// BiMap is the curried version of [MonadBiMap].
// Maps a pair of functions over the two type arguments of the bifunctor.
//
//go:inline
func BiMap[E, A, B any](f func(error) E, g func(a A) B) func(Result[A]) Either[E, B] {
	return either.BiMap(f, g)
}

// MonadMapTo replaces the Right value with a constant value.
// If the Either is Left, returns Left unchanged.
//
// Example:
//
//	result := either.MonadMapTo(either.Right[error](21), "success") // Right("success")
//
//go:inline
func MonadMapTo[A, B any](fa Result[A], b B) Result[B] {
	return either.MonadMapTo(fa, b)
}

// MapTo is the curried version of [MonadMapTo].
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return either.MapTo[error, A](b)
}

// MonadMapLeft applies a transformation function to the Left (error) value.
// If the Either is Right, returns Right unchanged.
//
// Example:
//
//	result := either.MonadMapLeft(
//	    either.Left[int](errors.New("error")),
//	    func(e error) string { return e.Error() },
//	) // Left("error")
//
//go:inline
func MonadMapLeft[A, E any](fa Result[A], f func(error) E) Either[E, A] {
	return either.MonadMapLeft(fa, f)
}

// Map is the curried version of [MonadMap].
// Transforms the Right value using the provided function.
//
//go:inline
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return either.Map[error](f)
}

// MapLeft is the curried version of [MonadMapLeft].
// Applies a mapping function to the Left (error) channel.
//
//go:inline
func MapLeft[A, E any](f func(error) E) func(fa Result[A]) Either[E, A] {
	return either.MapLeft[A](f)
}

// MonadChain sequences two computations, where the second depends on the result of the first.
// If the first Either is Left, returns Left without executing the second computation.
// This is the monadic bind operation (also known as flatMap).
//
// Example:
//
//	result := either.MonadChain(
//	    either.Right[error](21),
//	    func(x int) either.Result[int] {
//	        return either.Right[error](x * 2)
//	    },
//	) // Right(42)
//
//go:inline
func MonadChain[A, B any](fa Result[A], f Kleisli[A, B]) Result[B] {
	return either.MonadChain(fa, f)
}

// MonadChainFirst executes a side-effect computation but returns the original value.
// Useful for performing actions (like logging) without changing the value.
//
// Example:
//
//	result := either.MonadChainFirst(
//	    either.Right[error](42),
//	    func(x int) either.Result[string] {
//	        fmt.Println(x) // side effect
//	        return either.Right[error]("logged")
//	    },
//	) // Right(42) - original value preserved
//
//go:inline
func MonadChainFirst[A, B any](ma Result[A], f Kleisli[A, B]) Result[A] {
	return either.MonadChainFirst(ma, f)
}

// MonadChainTo ignores the first Either and returns the second.
// Useful for sequencing operations where you don't need the first result.
//
//go:inline
func MonadChainTo[A, B any](ma Result[A], mb Result[B]) Result[B] {
	return either.MonadChainTo(ma, mb)
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
//
//go:inline
func MonadChainOptionK[A, B any](onNone func() error, ma Result[A], f option.Kleisli[A, B]) Result[B] {
	return either.MonadChainOptionK(onNone, ma, f)
}

// ChainOptionK is the curried version of [MonadChainOptionK].
//
//go:inline
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B] {
	return either.ChainOptionK[A, B](onNone)
}

// ChainTo is the curried version of [MonadChainTo].
//
//go:inline
func ChainTo[A, B any](mb Result[B]) Operator[A, B] {
	return either.ChainTo[A](mb)
}

// Chain is the curried version of [MonadChain].
// Sequences two computations where the second depends on the first.
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return either.Chain(f)
}

// ChainFirst is the curried version of [MonadChainFirst].
//
//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return either.ChainFirst(f)
}

// Flatten removes one level of nesting from a nested Either.
//
// Example:
//
//	nested := either.Right[error](either.Right[error](42))
//	result := either.Flatten(nested) // Right(42)
//
//go:inline
func Flatten[A any](mma Result[Result[A]]) Result[A] {
	return either.Flatten(mma)
}

// TryCatch converts a (value, error) tuple into an Either, applying a transformation to the error.
//
// Example:
//
//	result := either.TryCatch(
//	    42, nil,
//	    func(err error) string { return err.Error() },
//	) // Right(42)
//
//go:inline
func TryCatch[FE Endomorphism[error], A any](val A, err error, onThrow FE) Result[A] {
	return either.TryCatch(val, err, onThrow)
}

// TryCatchError is a specialized version of [TryCatch] for error types.
// Converts a (value, error) tuple into Result[A].
//
// Example:
//
//	result := either.TryCatchError(42, nil) // Right(42)
//	result := either.TryCatchError(0, errors.New("fail")) // Left(error)
//
//go:inline
func TryCatchError[A any](val A, err error) Result[A] {
	return either.TryCatchError(val, err)
}

// Sequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
//
//go:inline
func Sequence2[T1, T2, R any](f func(T1, T2) Result[R]) func(Result[T1], Result[T2]) Result[R] {
	return either.Sequence2(f)
}

// Sequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
//
//go:inline
func Sequence3[T1, T2, T3, R any](f func(T1, T2, T3) Result[R]) func(Result[T1], Result[T2], Result[T3]) Result[R] {
	return either.Sequence3(f)
}

// FromOption converts an Option to an Either, using the provided function to generate a Left value for None.
//
// Example:
//
//	opt := option.Some(42)
//	result := either.FromOption[int](func() error { return errors.New("none") })(opt) // Right(42)
//
//go:inline
func FromOption[A any](onNone func() error) func(Option[A]) Result[A] {
	return either.FromOption[A](onNone)
}

// ToOption converts an Either to an Option, discarding the Left value.
//
// Example:
//
//	result := either.ToOption(either.Right[error](42)) // Some(42)
//	result := either.ToOption(either.Left[int](errors.New("err"))) // None
//
//go:inline
func ToOption[A any](ma Result[A]) Option[A] {
	return either.ToOption(ma)
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
//
//go:inline
func FromError[A any](f func(a A) error) Kleisli[A, A] {
	return either.FromError(f)
}

// ToError converts an Result[A] to an error, returning nil for Right values.
//
// Example:
//
//	err := either.ToError(either.Left[int](errors.New("fail"))) // error
//	err := either.ToError(either.Right[error](42)) // nil
//
//go:inline
func ToError[A any](e Result[A]) error {
	return either.ToError(e)
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
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(Result[A]) B {
	return either.Fold(onLeft, onRight)
}

// UnwrapError converts an Result[A] into the idiomatic Go tuple (A, error).
//
// Example:
//
//	val, err := either.UnwrapError(either.Right[error](42)) // 42, nil
//	val, err := either.UnwrapError(either.Left[int](errors.New("fail"))) // zero, error
//
//go:inline
func UnwrapError[A any](ma Result[A]) (A, error) {
	return either.UnwrapError(ma)
}

// FromPredicate creates an Either based on a predicate.
// If the predicate returns true, creates a Right; otherwise creates a Left using onFalse.
//
// Example:
//
//	isPositive := either.FromPredicate(
//	    func(x int) bool { return x > 0 },
//	    func(x int) error { return errors.New("not positive") },
//	)
//	result := isPositive(42) // Right(42)
//	result := isPositive(-1) // Left(error)
//
//go:inline
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return either.FromPredicate(pred, onFalse)
}

// FromNillable creates an Either from a pointer, using the provided error for nil pointers.
//
// Example:
//
//	var ptr *int = nil
//	result := either.FromNillable[int](errors.New("nil"))(ptr) // Left(error)
//	val := 42
//	result := either.FromNillable[int](errors.New("nil"))(&val) // Right(&42)
//
//go:inline
func FromNillable[A any](e error) func(*A) Result[*A] {
	return either.FromNillable[A](e)
}

// GetOrElse extracts the Right value or computes a default from the Left value.
//
// Example:
//
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Right[error](42)) // 42
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Left[int](err)) // 0
//
//go:inline
func GetOrElse[A any](onLeft func(error) A) func(Result[A]) A {
	return either.GetOrElse(onLeft)
}

// Reduce folds an Either into a single value using a reducer function.
// Returns the initial value for Left, or applies the reducer to the Right value.
//
//go:inline
func Reduce[A, B any](f func(B, A) B, initial B) func(Result[A]) B {
	return either.Reduce[error](f, initial)
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
//
//go:inline
func AltW[E1, A any](that Lazy[Either[E1, A]]) func(Result[A]) Either[E1, A] {
	return either.AltW[error](that)
}

// Alt provides an alternative Either if the first is Left.
//
// Example:
//
//	alternative := either.Alt[error](func() either.Result[int] {
//	    return either.Right[error](99)
//	})
//	result := alternative(either.Left[int](errors.New("fail"))) // Right(99)
//
//go:inline
func Alt[A any](that Lazy[Result[A]]) Operator[A, A] {
	return either.Alt(that)
}

// OrElse recovers from a Left by providing an alternative computation.
//
// Example:
//
//	recover := either.OrElse(func(err error) either.Result[int] {
//	    return either.Right[error](0) // default value
//	})
//	result := recover(either.Left[int](errors.New("fail"))) // Right(0)
//
//go:inline
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A] {
	return either.OrElse(onLeft)
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
//
//go:inline
func ToType[A any](onError func(any) error) Kleisli[any, A] {
	return either.ToType[A](onError)
}

// Memoize returns the Either unchanged (Either values are already memoized).
//
//go:inline
func Memoize[A any](val Result[A]) Result[A] {
	return either.Memoize(val)
}

// MonadSequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
//
//go:inline
func MonadSequence2[T1, T2, R any](e1 Result[T1], e2 Result[T2], f func(T1, T2) Result[R]) Result[R] {
	return either.MonadSequence2(e1, e2, f)
}

// MonadSequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
//
//go:inline
func MonadSequence3[T1, T2, T3, R any](e1 Result[T1], e2 Result[T2], e3 Result[T3], f func(T1, T2, T3) Result[R]) Result[R] {
	return either.MonadSequence3(e1, e2, e3, f)
}

// Swap exchanges the Left and Right type parameters.
//
// Example:
//
//	result := either.Swap(either.Right[error](42)) // Left(42)
//	result := either.Swap(either.Left[int](errors.New("err"))) // Right(error)
//
//go:inline
func Swap[A any](val Result[A]) Either[A, error] {
	return either.Swap(val)
}

// MonadFlap applies a value to a function wrapped in Either.
// This is the reverse of [MonadAp].
//
//go:inline
func MonadFlap[B, A any](fab Result[func(A) B], a A) Result[B] {
	return either.MonadFlap(fab, a)
}

// Flap is the curried version of [MonadFlap].
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return either.Flap[error, B](a)
}

// MonadAlt provides an alternative Either if the first is Left.
// This is the monadic version of [Alt].
//
//go:inline
func MonadAlt[A any](fa Result[A], that Lazy[Result[A]]) Result[A] {
	return either.MonadAlt(fa, that)
}
