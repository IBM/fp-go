// Copyright (c) 2023 IBM Corp.
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
	E "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

// Of constructs a Right value containing the given value.
// This is the monadic return/pure operation for Either.
// Equivalent to [Right].
//
// Example:
//
//	result := either.Of[error](42) // Right(42)
func Of[E, A any](value A) Either[E, A] {
	return F.Pipe1(value, Right[E, A])
}

// FromIO executes an IO operation and wraps the result in a Right value.
// This is useful for lifting pure IO operations into the Either context.
//
// Example:
//
//	getValue := func() int { return 42 }
//	result := either.FromIO[error](getValue) // Right(42)
func FromIO[E any, IO ~func() A, A any](f IO) Either[E, A] {
	return F.Pipe1(f(), Right[E, A])
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
func MonadAp[B, E, A any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	return MonadFold(fab, Left[B, E], func(ab func(A) B) Either[E, B] {
		return MonadFold(fa, Left[B, E], F.Flow2(ab, Right[E, B]))
	})
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
func Ap[B, E, A any](fa Either[E, A]) func(fab Either[E, func(a A) B]) Either[E, B] {
	return F.Bind2nd(MonadAp[B, E, A], fa)
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
func MonadMap[E, A, B any](fa Either[E, A], f func(a A) B) Either[E, B] {
	return MonadChain(fa, F.Flow2(f, Right[E, B]))
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
func MonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B] {
	return MonadFold(fa, F.Flow2(f, Left[B, E2]), F.Flow2(g, Right[E2, B]))
}

// BiMap is the curried version of [MonadBiMap].
// Maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B] {
	return Fold(F.Flow2(f, Left[B, E2]), F.Flow2(g, Right[E2, B]))
}

// MonadMapTo replaces the Right value with a constant value.
// If the Either is Left, returns Left unchanged.
//
// Example:
//
//	result := either.MonadMapTo(either.Right[error](21), "success") // Right("success")
func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	return MonadMap(fa, F.Constant1[A](b))
}

// MapTo is the curried version of [MonadMapTo].
func MapTo[E, A, B any](b B) func(Either[E, A]) Either[E, B] {
	return Map[E](F.Constant1[A](b))
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
func MonadMapLeft[E1, A, E2 any](fa Either[E1, A], f func(E1) E2) Either[E2, A] {
	return MonadFold(fa, F.Flow2(f, Left[A, E2]), Right[E2, A])
}

// Map is the curried version of [MonadMap].
// Transforms the Right value using the provided function.
func Map[E, A, B any](f func(a A) B) func(fa Either[E, A]) Either[E, B] {
	return Chain(F.Flow2(f, Right[E, B]))
}

// MapLeft is the curried version of [MonadMapLeft].
// Applies a mapping function to the Left (error) channel.
func MapLeft[A, E1, E2 any](f func(E1) E2) func(fa Either[E1, A]) Either[E2, A] {
	return Fold(F.Flow2(f, Left[A, E2]), Right[E2, A])
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
func MonadChain[E, A, B any](fa Either[E, A], f func(a A) Either[E, B]) Either[E, B] {
	return MonadFold(fa, Left[B, E], f)
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
func MonadChainFirst[E, A, B any](ma Either[E, A], f func(a A) Either[E, B]) Either[E, A] {
	return C.MonadChainFirst(
		MonadChain[E, A, A],
		MonadMap[E, B, A],
		ma,
		f,
	)
}

// MonadChainTo ignores the first Either and returns the second.
// Useful for sequencing operations where you don't need the first result.
func MonadChainTo[A, E, B any](_ Either[E, A], mb Either[E, B]) Either[E, B] {
	return mb
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
	return MonadChain(ma, F.Flow2(f, FromOption[B](onNone)))
}

// ChainOptionK is the curried version of [MonadChainOptionK].
func ChainOptionK[A, B, E any](onNone func() E) func(func(A) Option[B]) func(Either[E, A]) Either[E, B] {
	from := FromOption[B](onNone)
	return func(f func(A) Option[B]) func(Either[E, A]) Either[E, B] {
		return Chain(F.Flow2(f, from))
	}
}

// ChainTo is the curried version of [MonadChainTo].
func ChainTo[A, E, B any](mb Either[E, B]) func(Either[E, A]) Either[E, B] {
	return F.Constant1[Either[E, A]](mb)
}

// Chain is the curried version of [MonadChain].
// Sequences two computations where the second depends on the first.
func Chain[E, A, B any](f func(a A) Either[E, B]) func(Either[E, A]) Either[E, B] {
	return Fold(Left[B, E], f)
}

// ChainFirst is the curried version of [MonadChainFirst].
func ChainFirst[E, A, B any](f func(a A) Either[E, B]) func(Either[E, A]) Either[E, A] {
	return C.ChainFirst(
		Chain[E, A, A],
		Map[E, B, A],
		f,
	)
}

// Flatten removes one level of nesting from a nested Either.
//
// Example:
//
//	nested := either.Right[error](either.Right[error](42))
//	result := either.Flatten(nested) // Right(42)
func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A] {
	return MonadChain(mma, F.Identity[Either[E, A]])
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
	if err != nil {
		return F.Pipe2(err, onThrow, Left[A, E])
	}
	return F.Pipe1(val, Right[E, A])
}

// TryCatchError is a specialized version of [TryCatch] for error types.
// Converts a (value, error) tuple into Either[error, A].
//
// Example:
//
//	result := either.TryCatchError(42, nil) // Right(42)
//	result := either.TryCatchError(0, errors.New("fail")) // Left(error)
func TryCatchError[A any](val A, err error) Either[error, A] {
	return TryCatch(val, err, E.IdentityError)
}

// Sequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2]) Either[E, R] {
		return MonadSequence2(e1, e2, f)
	}
}

// Sequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3]) Either[E, R] {
		return MonadSequence3(e1, e2, e3, f)
	}
}

// FromOption converts an Option to an Either, using the provided function to generate a Left value for None.
//
// Example:
//
//	opt := option.Some(42)
//	result := either.FromOption[int](func() error { return errors.New("none") })(opt) // Right(42)
func FromOption[A, E any](onNone func() E) func(Option[A]) Either[E, A] {
	return O.Fold(F.Nullary2(onNone, Left[A, E]), Right[E, A])
}

// ToOption converts an Either to an Option, discarding the Left value.
//
// Example:
//
//	result := either.ToOption(either.Right[error](42)) // Some(42)
//	result := either.ToOption(either.Left[int](errors.New("err"))) // None
func ToOption[E, A any](ma Either[E, A]) Option[A] {
	return MonadFold(ma, F.Ignore1of1[E](O.None[A]), O.Some[A])
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
	return func(a A) Either[error, A] {
		return TryCatchError(a, f(a))
	}
}

// ToError converts an Either[error, A] to an error, returning nil for Right values.
//
// Example:
//
//	err := either.ToError(either.Left[int](errors.New("fail"))) // error
//	err := either.ToError(either.Right[error](42)) // nil
func ToError[A any](e Either[error, A]) error {
	return MonadFold(e, E.IdentityError, F.Constant1[A, error](nil))
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
func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B {
	return func(ma Either[E, A]) B {
		return MonadFold(ma, onLeft, onRight)
	}
}

// UnwrapError converts an Either[error, A] into the idiomatic Go tuple (A, error).
//
// Example:
//
//	val, err := either.UnwrapError(either.Right[error](42)) // 42, nil
//	val, err := either.UnwrapError(either.Left[int](errors.New("fail"))) // zero, error
func UnwrapError[A any](ma Either[error, A]) (A, error) {
	return Unwrap[error](ma)
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
func FromPredicate[E, A any](pred func(A) bool, onFalse func(A) E) func(A) Either[E, A] {
	return func(a A) Either[E, A] {
		if pred(a) {
			return Right[E](a)
		}
		return Left[A, E](onFalse(a))
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
func FromNillable[A, E any](e E) func(*A) Either[E, *A] {
	return FromPredicate(F.IsNonNil[A], F.Constant1[*A](e))
}

// GetOrElse extracts the Right value or computes a default from the Left value.
//
// Example:
//
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Right[error](42)) // 42
//	result := either.GetOrElse(func(err error) int { return 0 })(either.Left[int](err)) // 0
func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A {
	return Fold(onLeft, F.Identity[A])
}

// Reduce folds an Either into a single value using a reducer function.
// Returns the initial value for Left, or applies the reducer to the Right value.
func Reduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B {
	return Fold(
		F.Constant1[E](initial),
		F.Bind1st(f, initial),
	)
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
func AltW[E, E1, A any](that L.Lazy[Either[E1, A]]) func(Either[E, A]) Either[E1, A] {
	return Fold(F.Ignore1of1[E](that), Right[E1, A])
}

// Alt provides an alternative Either if the first is Left.
//
// Example:
//
//	alternative := either.Alt[error](func() either.Either[error, int] {
//	    return either.Right[error](99)
//	})
//	result := alternative(either.Left[int](errors.New("fail"))) // Right(99)
func Alt[E, A any](that L.Lazy[Either[E, A]]) func(Either[E, A]) Either[E, A] {
	return AltW[E](that)
}

// OrElse recovers from a Left by providing an alternative computation.
//
// Example:
//
//	recover := either.OrElse(func(err error) either.Either[error, int] {
//	    return either.Right[error](0) // default value
//	})
//	result := recover(either.Left[int](errors.New("fail"))) // Right(0)
func OrElse[E, A any](onLeft func(e E) Either[E, A]) func(Either[E, A]) Either[E, A] {
	return Fold(onLeft, Of[E, A])
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
	return func(value any) Either[E, A] {
		return F.Pipe2(
			value,
			O.ToType[A],
			O.Fold(F.Nullary3(F.Constant(value), onError, Left[A, E]), Right[E, A]),
		)
	}
}

// Memoize returns the Either unchanged (Either values are already memoized).
func Memoize[E, A any](val Either[E, A]) Either[E, A] {
	return val
}

// MonadSequence2 sequences two Either values using a combining function.
// Short-circuits on the first Left encountered.
func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[R, E], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[R, E], func(t2 T2) Either[E, R] {
			return f(t1, t2)
		})
	})
}

// MonadSequence3 sequences three Either values using a combining function.
// Short-circuits on the first Left encountered.
func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[R, E], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[R, E], func(t2 T2) Either[E, R] {
			return MonadFold(e3, Left[R, E], func(t3 T3) Either[E, R] {
				return f(t1, t2, t3)
			})
		})
	})
}

// Swap exchanges the Left and Right type parameters.
//
// Example:
//
//	result := either.Swap(either.Right[error](42)) // Left(42)
//	result := either.Swap(either.Left[int](errors.New("err"))) // Right(error)
func Swap[E, A any](val Either[E, A]) Either[A, E] {
	return MonadFold(val, Right[A, E], Left[E, A])
}

// MonadFlap applies a value to a function wrapped in Either.
// This is the reverse of [MonadAp].
func MonadFlap[E, B, A any](fab Either[E, func(A) B], a A) Either[E, B] {
	return FC.MonadFlap(MonadMap[E, func(A) B, B], fab, a)
}

// Flap is the curried version of [MonadFlap].
func Flap[E, B, A any](a A) func(Either[E, func(A) B]) Either[E, B] {
	return FC.Flap(Map[E, func(A) B, B], a)
}

// MonadAlt provides an alternative Either if the first is Left.
// This is the monadic version of [Alt].
func MonadAlt[E, A any](fa Either[E, A], that L.Lazy[Either[E, A]]) Either[E, A] {
	return MonadFold(fa, F.Ignore1of1[E](that), Of[E, A])
}
