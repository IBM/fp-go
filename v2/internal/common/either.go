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
package common

import (
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	FC "github.com/IBM/fp-go/v2/internal/functor"
)

// EitherOf constructs a Right value containing the given value.
// This is the monadic return/pure operation for Either.
// Equivalent to [EitherRight].
//
// Example:
//
//	result := either.EitherOf[error](42) // Right(42)
//
//go:inline
func EitherOf[E, A any](value A) Either[E, A] {
	return EitherRight[E](value)
}

// EitherFromIO executes an IO operation and wraps the result in a Right value.
// This is useful for lifting pure IO operations into the Either context.
//
// Example:
//
//	getValue := lazy.Of(42)
//	result := either.EitherFromIO[error](getValue) // Right(42)
//
// go: inline
func EitherFromIO[E any, IO ~func() A, A any](f IO) Either[E, A] {
	return EitherOf[E](f())
}

// EitherMonadAp applies a function wrapped in Either to a value wrapped in Either.
// If either the function or the value is Left, returns Left.
// This is the applicative apply operation.
//
// Example:
//
//	fab := either.EitherRight[error](N.Mul(2))
//	fa := either.EitherRight[error](21)
//	result := either.EitherMonadAp(fab, fa) // Right(42)
func EitherMonadAp[B, E, A any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	if fab.l {
		return EitherLeft[B](fab.e)
	}
	if fa.l {
		return EitherLeft[B](fa.e)
	}
	return EitherOf[E](fab.a(fa.a))
}

// EitherAp is the curried version of [EitherMonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
//
//go:inline
func EitherAp[B, E, A any](fa Either[E, A]) EitherOperator[E, func(A) B, B] {
	return F.Bind2nd(EitherMonadAp[B, E, A], fa)
}

// EitherMonadMap transforms the Right value using the provided function.
// If the Either is Left, returns Left unchanged.
// This is the functor map operation.
//
// Example:
//
//	result := either.EitherMonadMap(
//	    either.EitherRight[error](21),
//	    N.Mul(2),
//	) // Right(42)
//
//go:inline
func EitherMonadMap[E, A, B any](fa Either[E, A], f func(A) B) Either[E, B] {
	if fa.l {
		return EitherLeft[B](fa.e)
	}
	return EitherOf[E](f(fa.a))
}

// EitherMonadBiMap applies two functions: one to transform a Left value, another to transform a Right value.
// This allows transforming both channels of the Either simultaneously.
//
// Example:
//
//	result := either.EitherMonadBiMap(
//	    either.EitherLeft[int](errors.New("error")),
//	    error.Error,
//	    func(n int) string { return fmt.Sprint(n) },
//	) // Left("error")
func EitherMonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B] {
	if fa.l {
		return EitherLeft[B](f(fa.e))
	}
	return EitherOf[E2](g(fa.a))
}

// EitherBiMap is the curried version of [EitherMonadBiMap].
// Maps a pair of functions over the two type arguments of the bifunctor.
func EitherBiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B] {
	return func(fa Either[E1, A]) Either[E2, B] {
		return EitherMonadBiMap(fa, f, g)
	}
}

// EitherMonadMapTo replaces the Right value with a constant value.
// If the Either is Left, returns Left unchanged.
//
// Example:
//
//	result := either.EitherMonadMapTo(either.EitherRight[error](21), "success") // Right("success")
func EitherMonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	if fa.l {
		return EitherLeft[B](fa.e)
	}
	return EitherOf[E](b)
}

// EitherMapTo is the curried version of [EitherMonadMapTo].
func EitherMapTo[E, A, B any](b B) EitherOperator[E, A, B] {
	return F.Bind2nd(EitherMonadMapTo[E, A], b)
}

// EitherMonadMapLeft applies a transformation function to the Left (error) value.
// If the Either is Right, returns Right unchanged.
//
// Example:
//
//	result := either.EitherMonadMapLeft(
//	    either.EitherLeft[int](errors.New("error")),
//	    error.Error,
//	) // Left("error")
func EitherMonadMapLeft[E1, A, E2 any](fa Either[E1, A], f func(E1) E2) Either[E2, A] {
	return EitherMonadFold(fa, F.Flow2(f, EitherLeft[A, E2]), EitherRight[E2, A])
}

// EitherMap is the curried version of [EitherMonadMap].
// Transforms the Right value using the provided function.
func EitherMap[E, A, B any](f func(A) B) EitherOperator[E, A, B] {
	return F.Bind2nd(EitherMonadMap[E], f)
}

// EitherMapLeft is the curried version of [EitherMonadMapLeft].
// Applies a mapping function to the Left (error) channel.
func EitherMapLeft[A, E1, E2 any](f func(E1) E2) func(fa Either[E1, A]) Either[E2, A] {
	return EitherFold(F.Flow2(f, EitherLeft[A, E2]), EitherRight[E2, A])
}

// EitherMonadChain sequences two computations, where the second depends on the result of the first.
// If the first Either is Left, returns Left without executing the second computation.
// This is the monadic bind operation (also known as flatMap).
//
// Example:
//
//	result := either.EitherMonadChain(
//	    either.EitherRight[error](21),
//	    func(x int) either.Either[error, int] {
//	        return either.EitherRight[error](x * 2)
//	    },
//	) // Right(42)
//
//go:inline
func EitherMonadChain[E, A, B any](fa Either[E, A], f EitherKleisli[E, A, B]) Either[E, B] {
	if fa.l {
		return EitherLeft[B](fa.e)
	}
	return f(fa.a)
}

// EitherMonadChainLeft sequences a computation on the Left (error) value, allowing error recovery or transformation.
// If the Either is Left, applies the provided function to the error value, which returns a new Either.
// If the Either is Right, returns the Right value unchanged with the new error type.
//
// This is the dual of [EitherMonadChain] - while EitherMonadChain operates on Right values (success),
// EitherMonadChainLeft operates on Left values (errors). It's useful for error recovery, error transformation,
// or chaining alternative computations when an error occurs.
//
// Note: EitherMonadChainLeft is identical to [EitherOrElse] - both provide the same functionality for error recovery.
//
// The error type can be transformed from EA to EB, allowing flexible error type conversions.
//
// Example:
//
//	// Error recovery: convert specific errors to success
//	result := either.EitherMonadChainLeft(
//	    either.EitherLeft[int](errors.New("not found")),
//	    func(err error) either.Either[string, int] {
//	        if err.Error() == "not found" {
//	            return either.EitherRight[string](0) // default value
//	        }
//	        return either.EitherLeft[int](err.Error()) // transform error
//	    },
//	) // Right(0)
//
//	// Error transformation: change error type
//	result := either.EitherMonadChainLeft(
//	    either.EitherLeft[int](404),
//	    func(code int) either.Either[string, int] {
//	        return either.EitherLeft[int](fmt.Sprintf("Error code: %d", code))
//	    },
//	) // Left("Error code: 404")
//
//	// Right values pass through unchanged
//	result := either.EitherMonadChainLeft(
//	    either.EitherRight[error](42),
//	    func(err error) either.Either[string, int] {
//	        return either.EitherLeft[int]("error")
//	    },
//	) // Right(42)
//
//go:inline
func EitherMonadChainLeft[EA, EB, A any](fa Either[EA, A], f EitherKleisli[EB, EA, A]) Either[EB, A] {
	return EitherMonadFold(fa, f, EitherOf[EB])
}

// EitherChainLeft is the curried version of [EitherMonadChainLeft].
// Returns a function that sequences a computation on the Left (error) value.
//
// Note: EitherChainLeft is identical to [EitherOrElse] - both provide the same functionality for error recovery.
//
// This is useful for creating reusable error handlers or transformers that can be
// composed with other Either operations using pipes or function composition.
//
// Example:
//
//	// Create a reusable error handler
//	handleNotFound := either.EitherChainLeft[error, string](func(err error) either.Either[string, int] {
//	    if err.Error() == "not found" {
//	        return either.EitherRight[string](0)
//	    }
//	    return either.EitherLeft[int](err.Error())
//	})
//
//	// Use in a pipeline
//	result := F.Pipe1(
//	    either.EitherLeft[int](errors.New("not found")),
//	    handleNotFound,
//	) // Right(0)
//
//go:inline
func EitherChainLeft[EA, EB, A any](f EitherKleisli[EB, EA, A]) EitherKleisli[EB, Either[EA, A], A] {
	return EitherFold(f, EitherOf[EB])
}

// EitherMonadChainFirst executes a side-effect computation but returns the original value.
// Useful for performing actions (like logging) without changing the value.
//
// Example:
//
//	result := either.EitherMonadChainFirst(
//	    either.EitherRight[error](42),
//	    func(x int) either.Either[error, string] {
//	        fmt.Println(x) // side effect
//	        return either.EitherRight[error]("logged")
//	    },
//	) // Right(42) - original value preserved
func EitherMonadChainFirst[E, A, B any](ma Either[E, A], f EitherKleisli[E, A, B]) Either[E, A] {
	return C.MonadChainFirst(
		EitherMonadChain[E, A, A],
		EitherMonadMap[E, B, A],
		ma,
		f,
	)
}

// EitherMonadChainTo ignores the first Either and returns the second.
// Useful for sequencing operations where you don't need the first result.
func EitherMonadChainTo[A, E, B any](_ Either[E, A], mb Either[E, B]) Either[E, B] {
	return mb
}

// EitherMonadChainOptionK chains a function that returns an Option, converting None to Left.
//
// Example:
//
//	result := either.EitherMonadChainOptionK(
//	    func() error { return errors.New("not found") },
//	    either.EitherRight[error](42),
//	    func(x int) option.Option[string] {
//	        if x > 0 { return option.Some("positive") }
//	        return option.None[string]()
//	    },
//	) // Right("positive")
func EitherMonadChainOptionK[A, B, E any](onNone func() E, ma Either[E, A], f func(A) Option[B]) Either[E, B] {
	return EitherMonadChain(ma, F.Flow2(f, EitherFromOption[B](onNone)))
}

// EitherChainOptionK is the curried version of [EitherMonadChainOptionK].
func EitherChainOptionK[A, B, E any](onNone func() E) func(func(A) Option[B]) EitherOperator[E, A, B] {
	from := EitherFromOption[B](onNone)
	return func(f func(A) Option[B]) EitherOperator[E, A, B] {
		return EitherChain(F.Flow2(f, from))
	}
}

// EitherChainTo is the curried version of [EitherMonadChainTo].
func EitherChainTo[A, E, B any](mb Either[E, B]) EitherOperator[E, A, B] {
	return F.Constant1[Either[E, A]](mb)
}

// EitherChain is the curried version of [EitherMonadChain].
// Sequences two computations where the second depends on the first.
func EitherChain[E, A, B any](f EitherKleisli[E, A, B]) EitherOperator[E, A, B] {
	return F.Bind2nd(EitherMonadChain[E], f)
}

// EitherChainFirst is the curried version of [EitherMonadChainFirst].
func EitherChainFirst[E, A, B any](f EitherKleisli[E, A, B]) EitherOperator[E, A, A] {
	return C.ChainFirst(
		EitherChain[E, A, A],
		EitherMap[E, B, A],
		f,
	)
}

// EitherFlatten removes one level of nesting from a nested Either.
//
// Example:
//
//	nested := either.EitherRight[error](either.EitherRight[error](42))
//	result := either.EitherFlatten(nested) // Right(42)
func EitherFlatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A] {
	return EitherMonadChain(mma, F.Identity[Either[E, A]])
}

// EitherTryCatch converts a (value, error) tuple into an Either, applying a transformation to the error.
//
// Example:
//
//	result := either.EitherTryCatch(
//	    42, nil,
//	    func(err error) string { return err.Error() },
//	) // Right(42)
func EitherTryCatch[FE func(error) E, E, A any](val A, err error, onThrow FE) Either[E, A] {
	if err != nil {
		return F.Pipe2(err, onThrow, EitherLeft[A, E])
	}
	return F.Pipe1(val, EitherRight[E, A])
}

// EitherTryCatchError is a specialized version of [EitherTryCatch] for error types.
// Converts a (value, error) tuple into Either[error, A].
//
// Example:
//
//	result := either.EitherTryCatchError(42, nil) // Right(42)
//	result := either.EitherTryCatchError(0, errors.New("fail")) // Left(error)
func EitherTryCatchError[A any](val A, err error) Either[error, A] {
	return EitherTryCatch(val, err, F.Identity[error])
}

// EitherSequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func EitherSequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2]) Either[E, R] {
		return EitherMonadSequence2(e1, e2, f)
	}
}

// EitherSequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func EitherSequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3]) Either[E, R] {
		return EitherMonadSequence3(e1, e2, e3, f)
	}
}

// EitherFromOption converts an Option to an Either, using the provided function to generate a Left value for None.
//
// Example:
//
//	opt := option.Some(42)
//	result := either.EitherFromOption[int](func() error { return errors.New("none") })(opt) // Right(42)
func EitherFromOption[A, E any](onNone func() E) func(Option[A]) Either[E, A] {
	return OptionFold(F.Nullary2(onNone, EitherLeft[A, E]), EitherRight[E, A])
}

// EitherToOption converts an Either to an Option, discarding the Left value.
//
// Example:
//
//	result := either.EitherToOption(either.EitherRight[error](42)) // Some(42)
//	result := either.EitherToOption(either.EitherLeft[int](errors.New("err"))) // None
//
//go:inline
func EitherToOption[E, A any](ma Either[E, A]) Option[A] {
	return EitherMonadFold(ma, F.Ignore1of1[E](OptionNone[A]), OptionSome[A])
}

// EitherFromError creates an Either from a function that may return an error.
//
// Example:
//
//	validate := func(x int) error {
//	    if x < 0 { return errors.New("negative") }
//	    return nil
//	}
//	toEither := either.EitherFromError(validate)
//	result := toEither(42) // Right(42)
func EitherFromError[A any](f func(a A) error) func(A) Either[error, A] {
	return func(a A) Either[error, A] {
		return EitherTryCatchError(a, f(a))
	}
}

// EitherToError converts an Either[error, A] to an error, returning nil for Right values.
//
// Example:
//
//	err := either.EitherToError(either.EitherLeft[int](errors.New("fail"))) // error
//	err := either.EitherToError(either.EitherRight[error](42)) // nil
func EitherToError[A any](e Either[error, A]) error {
	return EitherMonadFold(e, F.Identity[error], F.Constant1[A, error](nil))
}

// EitherFold is the curried version of [EitherMonadFold].
// Extracts the value from an Either by providing handlers for both cases.
//
// Example:
//
//	result := either.EitherFold(
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Value: %d", n) },
//	)(either.EitherRight[error](42)) // "Value: 42"
//
//go:inline
func EitherFold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B {
	return func(ma Either[E, A]) B {
		return EitherMonadFold(ma, onLeft, onRight)
	}
}

// EitherUnwrapError converts an Either[error, A] into the idiomatic Go tuple (A, error).
//
// Example:
//
//	val, err := either.EitherUnwrapError(either.EitherRight[error](42)) // 42, nil
//	val, err := either.EitherUnwrapError(either.EitherLeft[int](errors.New("fail"))) // zero, error
//
//go:inline
func EitherUnwrapError[A any](ma Either[error, A]) (A, error) {
	return EitherUnwrap(ma)
}

// EitherFromPredicate creates an Either based on a predicate.
// If the predicate returns true, creates a Right; otherwise creates a Left using onFalse.
//
// Example:
//
//	isPositive := either.EitherFromPredicate(
//	    N.MoreThan(0),
//	    func(x int) error { return errors.New("not positive") },
//	)
//	result := isPositive(42) // Right(42)
//	result := isPositive(-1) // Left(error)
func EitherFromPredicate[E, A any](pred Predicate[A], onFalse func(A) E) EitherKleisli[E, A, A] {
	return func(a A) Either[E, A] {
		if pred(a) {
			return EitherRight[E](a)
		}
		return EitherLeft[A](onFalse(a))
	}
}

// EitherFromNillable creates an Either from a pointer, using the provided error for nil pointers.
//
// Example:
//
//	var ptr *int = nil
//	result := either.EitherFromNillable[int](errors.New("nil"))(ptr) // Left(error)
//	val := 42
//	result := either.EitherFromNillable[int](errors.New("nil"))(&val) // Right(&42)
func EitherFromNillable[A, E any](e E) EitherKleisli[E, *A, *A] {
	return EitherFromPredicate(F.IsNonNil[A], F.Constant1[*A](e))
}

// EitherGetOrElse extracts the Right value or computes a default from the Left value.
//
// Example:
//
//	result := either.EitherGetOrElse(func(err error) int { return 0 })(either.EitherRight[error](42)) // 42
//	result := either.EitherGetOrElse(func(err error) int { return 0 })(either.EitherLeft[int](err)) // 0
func EitherGetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A {
	return EitherFold(onLeft, F.Identity[A])
}

// EitherReduce folds an Either into a single value using a reducer function.
// Returns the initial value for Left, or applies the reducer to the Right value.
func EitherReduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B {
	return func(fa Either[E, A]) B {
		if fa.l {
			return initial
		}
		return f(initial, fa.a)
	}
}

// EitherAltW provides an alternative Either if the first is Left, allowing different error types.
// The 'W' suffix indicates "widening" of the error type.
//
// Example:
//
//	alternative := either.EitherAltW[error, string](func() either.Either[string, int] {
//	    return either.EitherRight[string](99)
//	})
//	result := alternative(either.EitherLeft[int](errors.New("fail"))) // Right(99)
func EitherAltW[E, E1, A any](that func() Either[E1, A]) EitherKleisli[E1, Either[E, A], A] {
	return EitherFold(F.Ignore1of1[E](that), EitherRight[E1, A])
}

// EitherAlt provides an alternative Either if the first is Left.
//
// Example:
//
//	alternative := either.EitherAlt[error](func() either.Either[error, int] {
//	    return either.EitherRight[error](99)
//	})
//	result := alternative(either.EitherLeft[int](errors.New("fail"))) // Right(99)
func EitherAlt[E, A any](that func() Either[E, A]) EitherOperator[E, A, A] {
	return EitherAltW[E](that)
}

// EitherOrElse recovers from a Left (error) by providing an alternative computation.
// If the Either is Right, it returns the value unchanged.
// If the Either is Left, it applies the provided function to the error value,
// which returns a new Either that replaces the original.
//
// Note: EitherOrElse is identical to [EitherChainLeft] - both provide the same functionality for error recovery.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations.
// The error type can be widened from E1 to E2, allowing transformation of error types.
//
// Example:
//
//	// Recover from specific errors with fallback values
//	recover := either.EitherOrElse(func(err error) either.Either[error, int] {
//	    if err.Error() == "not found" {
//	        return either.EitherRight[error](0) // default value
//	    }
//	    return either.EitherLeft[int](err) // propagate other errors
//	})
//	result := recover(either.EitherLeft[int](errors.New("not found"))) // Right(0)
//	result := recover(either.EitherRight[error](42)) // Right(42) - unchanged
//
//go:inline
func EitherOrElse[E1, E2, A any](onLeft EitherKleisli[E2, E1, A]) EitherKleisli[E2, Either[E1, A], A] {
	return EitherFold(onLeft, EitherOf[E2, A])
}

// EitherToType attempts to convert an any value to a specific type, returning Either.
//
// Example:
//
//	convert := either.EitherToType[int](func(v any) error {
//	    return fmt.Errorf("cannot convert %v to int", v)
//	})
//	result := convert(42) // Right(42)
//	result := convert("string") // Left(error)
func EitherToType[A, E any](onError func(any) E) func(any) Either[E, A] {
	return func(value any) Either[E, A] {
		return F.Pipe2(
			value,
			OptionInstanceOf[A],
			OptionFold(F.Nullary3(F.Constant(value), onError, EitherLeft[A, E]), EitherRight[E, A]),
		)
	}
}

// EitherMemoize returns the Either unchanged (Either values are already memoized).
func EitherMemoize[E, A any](val Either[E, A]) Either[E, A] {
	return val
}

// EitherMonadSequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func EitherMonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	if e1.l {
		return EitherLeft[R](e1.e)
	}
	if e2.l {
		return EitherLeft[R](e2.e)
	}
	return f(e1.a, e2.a)
}

// EitherMonadSequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func EitherMonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	if e1.l {
		return EitherLeft[R](e1.e)
	}
	if e2.l {
		return EitherLeft[R](e2.e)
	}
	if e3.l {
		return EitherLeft[R](e3.e)
	}
	return f(e1.a, e2.a, e3.a)
}

// EitherSwap exchanges the Left and Right type parameters.
//
// Example:
//
//	result := either.EitherSwap(either.EitherRight[error](42)) // Left(42)
//	result := either.EitherSwap(either.EitherLeft[int](errors.New("err"))) // Right(error)
//
//go:inline
func EitherSwap[E, A any](val Either[E, A]) Either[A, E] {
	return EitherMonadFold(val, EitherRight[A, E], EitherLeft[E, A])
}

// EitherMonadFlap applies a value to a function wrapped in Either.
// This is the reverse of [EitherMonadAp].
func EitherMonadFlap[E, B, A any](fab Either[E, func(A) B], a A) Either[E, B] {
	return FC.MonadFlap(EitherMonadMap[E, func(A) B, B], fab, a)
}

// EitherFlap is the curried version of [EitherMonadFlap].
func EitherFlap[E, B, A any](a A) EitherOperator[E, func(A) B, B] {
	return FC.Flap(EitherMap[E, func(A) B, B], a)
}

// EitherMonadAlt provides an alternative Either if the first is Left.
// This is the monadic version of [EitherAlt].
func EitherMonadAlt[E, A any](fa Either[E, A], that func() Either[E, A]) Either[E, A] {
	return EitherMonadFold(fa, F.Ignore1of1[E](that), EitherOf[E, A])
}

// EitherZero returns the zero value of an [Either], which is a Right containing the zero value of type A.
// This function is useful as an identity element in monoid operations or for creating an empty Either
// in a Right state.
//
// The returned Either is always a Right value containing the zero value of type A. For reference types
// (pointers, slices, maps, channels, functions, interfaces), the zero value is nil. For value types
// (numbers, booleans, structs), it's the type's zero value.
//
// Important: EitherZero() returns the same value as the default initialization of Either[E, A].
// When you declare `var e Either[E, A]` without initialization, it has the same value as EitherZero[E, A]().
//
// Note: This differs from creating a Left value, which would represent an error or failure state.
// EitherZero always produces a successful (Right) state with a zero value.
//
// Example:
//
//	// Zero Either with int value
//	e1 := either.EitherZero[error, int]()  // Right(0)
//
//	// Zero Either with string value
//	e2 := either.EitherZero[error, string]()  // Right("")
//
//	// Zero Either with pointer type
//	e3 := either.EitherZero[error, *int]()  // Right(nil)
//
//	// Zero equals default initialization
//	var defaultInit Either[error, int]
//	zero := either.EitherZero[error, int]()
//	assert.Equal(t, defaultInit, zero) // true
//
//	// Verify it's a Right value
//	e := either.EitherZero[error, int]()
//	assert.True(t, either.IsRight(e))  // true
//	assert.False(t, either.IsLeft(e))  // false
func EitherZero[E, A any]() Either[E, A] {
	return Either[E, A]{}
}
