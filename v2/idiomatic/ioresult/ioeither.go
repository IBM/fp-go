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

package ioresult

import (
	"sync"
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/result"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	RES "github.com/IBM/fp-go/v2/result"
)

// Left creates an IOResult that represents a failed computation with the given error.
// When executed, it returns the zero value for type A and the provided error.
func Left[A any](l error) IOResult[A] {
	return func() (A, error) {
		return result.Left[A](l)
	}
}

// Right creates an IOResult that represents a successful computation with the given value.
// When executed, it returns the provided value and nil error.
func Right[A any](r A) IOResult[A] {
	return func() (A, error) {
		return result.Of(r)
	}
}

//go:inline
func Of[A any](r A) IOResult[A] {
	return Right(r)
}

//go:inline
func MonadOf[A any](r A) IOResult[A] {
	return Of(r)
}

// LeftIO creates an IOResult from an IO computation that produces an error.
// The error from the IO is used as the Left value.
func LeftIO[A any](ml IO[error]) IOResult[A] {
	return func() (a A, e error) {
		e = ml()
		return
	}
}

// RightIO creates an IOResult from an IO computation that produces a value.
// The IO is executed and its result is wrapped in a successful IOResult.
func RightIO[A any](mr IO[A]) IOResult[A] {
	return func() (A, error) {
		return result.Of(mr())
	}
}

// FromEither converts an Either (Result[A]) to an IOResult.
// Either's Left becomes an error, Either's Right becomes a successful value.
func FromEither[A any](e Result[A]) IOResult[A] {
	return func() (A, error) {
		return RES.Unwrap(e)
	}
}

// FromResult converts a (value, error) tuple to an IOResult.
// This is the primary way to convert Go's standard error handling pattern to IOResult.
func FromResult[A any](a A, err error) IOResult[A] {
	return func() (A, error) {
		return a, err
	}
}

// FromOption converts an Option (represented as value, bool) to an IOResult.
// If the bool is true, the value is wrapped in a successful IOResult.
// If the bool is false, onNone is called to generate the error.
func FromOption[A any](onNone Lazy[error]) func(A, bool) IOResult[A] {
	return func(a A, ok bool) IOResult[A] {
		return func() (A, error) {
			if ok {
				return result.Of(a)
			}
			return result.Left[A](onNone())
		}
	}
}

// ChainOptionK chains a function that returns an Option (value, bool).
// The None case (false) is converted to an error using onNone.
func ChainOptionK[A, B any](onNone Lazy[error]) func(func(A) (B, bool)) Operator[A, B] {
	return func(f func(A) (B, bool)) Operator[A, B] {
		return func(i IOResult[A]) IOResult[B] {
			return func() (B, error) {
				a, err := i()
				if err != nil {
					return result.Left[B](err)
				}
				b, ok := f(a)
				if ok {
					return result.Of(b)
				}
				return result.Left[B](onNone())
			}
		}
	}
}

// MonadChainIOK chains an IO kleisli function to an IOResult.
// If the IOResult fails, the function is not executed. Otherwise, the IO is executed and wrapped.
func MonadChainIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[B] {
	return fromio.MonadChainIOK(
		MonadChain[A, B],
		FromIO[B],
		ma,
		f,
	)
}

//go:inline
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B] {
	return fromio.ChainIOK(
		Chain[A, B],
		FromIO[B],
		f,
	)
}

//go:inline
func ChainLazyK[A, B any](f func(A) Lazy[B]) Operator[A, B] {
	return ChainIOK(f)
}

//go:inline
func FromIO[A any](mr IO[A]) IOResult[A] {
	return RightIO(mr)
}

//go:inline
func FromLazy[A any](mr Lazy[A]) IOResult[A] {
	return FromIO(mr)
}

// MonadMap transforms the value inside an IOResult using the given function.
// If the IOResult is a Left (error), the function is not applied.
func MonadMap[A, B any](fa IOResult[A], f func(A) B) IOResult[B] {
	return func() (B, error) {
		a, err := fa()
		if err != nil {
			return result.Left[B](err)
		}
		return result.Of(f(a))
	}
}

// Map returns an operator that transforms values using the given function.
// This is the Functor map operation for IOResult.
func Map[A, B any](f func(A) B) Operator[A, B] {
	return function.Bind2nd(MonadMap[A, B], f)
}

//go:inline
func MonadMapTo[A, B any](fa IOResult[A], b B) IOResult[B] {
	return MonadMap(fa, function.Constant1[A](b))
}

//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return function.Bind2nd(MonadMapTo[A, B], b)
}

// MonadChain chains a kleisli function that depends on the current value.
// This is the Monad bind operation for IOResult.
func MonadChain[A, B any](fa IOResult[A], f Kleisli[A, B]) IOResult[B] {
	return func() (B, error) {
		a, err := fa()
		if err != nil {
			return result.Left[B](err)
		}
		return f(a)()
	}
}

// Chain returns an operator that chains a kleisli function.
// This enables dependent computations where the next step depends on the previous result.
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return function.Bind2nd(MonadChain[A, B], f)
}

// MonadChainEitherK chains a function that returns an Either.
// The Either is converted to IOResult: Left becomes error, Right becomes success.
func MonadChainEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[B] {
	return func() (B, error) {
		a, err := ma()
		if err != nil {
			return result.Left[B](err)
		}
		return either.Unwrap(f(a))
	}
}

//go:inline
func ChainEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B] {
	return function.Bind2nd(MonadChainEitherK[A, B], f)
}

// MonadChainResultK chains a function that returns a (value, error) tuple.
// This allows chaining standard Go functions that return errors.
func MonadChainResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[B] {
	return func() (B, error) {
		a, err := ma()
		if err != nil {
			return result.Left[B](err)
		}
		return f(a)
	}
}

//go:inline
func ChainResultK[A, B any](f result.Kleisli[A, B]) Operator[A, B] {
	return function.Bind2nd(MonadChainResultK[A, B], f)
}

//go:inline
func MonadAp[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B] {
	return MonadApPar(mab, ma)
}

//go:inline
func Ap[B, A any](ma IOResult[A]) Operator[func(A) B, B] {
	return ApPar[B](ma)
}

// MonadApPar applies a function to a value, executing both in parallel.
// Both IOResults are executed concurrently for better performance.
func MonadApPar[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B] {
	return func() (B, error) {
		var wg sync.WaitGroup
		wg.Add(1)

		var fab func(A) B
		var faberr error

		go func() {
			defer wg.Done()
			fab, faberr = mab()
		}()

		fa, faerr := ma()
		wg.Wait()

		if faberr != nil {
			return result.Left[B](faberr)
		}
		if faerr != nil {
			return result.Left[B](faerr)
		}

		return result.Of(fab(fa))
	}
}

//go:inline
func ApPar[B, A any](ma IOResult[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApPar[B, A], ma)
}

// MonadApSeq applies a function to a value sequentially.
// The function IOResult is executed first, then the value IOResult.
func MonadApSeq[B, A any](mab IOResult[func(A) B], ma IOResult[A]) IOResult[B] {
	return func() (B, error) {
		fab, err := mab()
		if err != nil {
			return result.Left[B](err)
		}

		fa, err := ma()
		if err != nil {
			return result.Left[B](err)
		}

		return result.Of(fab(fa))
	}
}

//go:inline
func ApSeq[B, A any](ma IOResult[A]) func(IOResult[func(A) B]) IOResult[B] {
	return function.Bind2nd(MonadApSeq[B, A], ma)
}

//go:inline
func Flatten[A any](mma IOResult[IOResult[A]]) IOResult[A] {
	return MonadChain(mma, function.Identity[IOResult[A]])
}

// Memoize caches the result of an IOResult so it only executes once.
// Subsequent calls return the cached result without re-executing the computation.
func Memoize[A any](ma IOResult[A]) IOResult[A] {
	// synchronization primitives
	var once sync.Once
	var fa A
	var faerr error
	// callback
	gen := func() {
		fa, faerr = ma()
	}
	// returns our memoized wrapper
	return func() (A, error) {
		once.Do(gen)
		return fa, faerr
	}
}

// MonadMapLeft transforms the error value using the given function.
// The success value is left unchanged.
func MonadMapLeft[A any](fa IOResult[A], f Endomorphism[error]) IOResult[A] {
	return func() (A, error) {
		a, err := fa()
		if err != nil {
			return result.Left[A](f(err))
		}
		return result.Of(a)
	}
}

//go:inline
func MapLeft[A any](f Endomorphism[error]) Operator[A, A] {
	return function.Bind2nd(MonadMapLeft[A], f)
}

// MonadBiMap transforms both the error (left) and success (right) values.
func MonadBiMap[A, B any](fa IOResult[A], f Endomorphism[error], g func(A) B) IOResult[B] {
	return func() (B, error) {
		a, err := fa()
		if err != nil {
			return result.Left[B](f(err))
		}
		return result.Of(g(a))
	}
}

//go:inline
func BiMap[A, B any](f Endomorphism[error], g func(A) B) Operator[A, B] {
	return function.Bind23of3(MonadBiMap[A, B])(f, g)
}

//go:inline
func Fold[A, B any](onLeft func(error) IO[B], onRight io.Kleisli[A, B]) func(IOResult[A]) IO[B] {
	return function.Bind23of3(MonadFold[A, B])(onLeft, onRight)
}

// GetOrElse extracts the value from an IOResult, using a default IO for error cases.
// This converts an IOResult to an IO that cannot fail.
func GetOrElse[A any](onLeft func(error) IO[A]) func(IOResult[A]) IO[A] {
	return func(fa IOResult[A]) IO[A] {
		return func() A {
			a, err := fa()
			if err != nil {
				return onLeft(err)()
			}
			return a
		}
	}
}

//go:inline
func MonadChainTo[A, B any](fa IOResult[A], fb IOResult[B]) IOResult[B] {
	return MonadChain(fa, function.Constant1[A](fb))
}

//go:inline
func ChainTo[A, B any](fb IOResult[B]) Operator[A, B] {
	return function.Bind2nd(MonadChainTo[A, B], fb)
}

// MonadChainFirst chains a computation but returns the original value if both succeed.
// If either computation fails, the error is returned.
func MonadChainFirst[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A] {
	return chain.MonadChainFirst(
		MonadChain[A, A],
		MonadMap[B, A],
		ma,
		f,
	)
}

//go:inline
func MonadTap[A, B any](ma IOResult[A], f Kleisli[A, B]) IOResult[A] {
	return MonadChainFirst(ma, f)
}

//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(
		Chain[A, A],
		Map[B, A],
		f,
	)
}

//go:inline
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return ChainFirst(f)
}

func MonadChainFirstEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[A] {
	return func() (A, error) {
		a, err := ma()
		if err != nil {
			return result.Left[A](err)
		}
		_, err = either.Unwrap(f(a))
		if err != nil {
			return result.Left[A](err)
		}
		return result.Of(a)
	}
}

//go:inline
func ChainFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return function.Bind2nd(MonadChainFirstEitherK[A, B], f)
}

func MonadChainFirstResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A] {
	return func() (A, error) {
		a, err := ma()
		if err != nil {
			return result.Left[A](err)
		}
		_, err = f(a)
		if err != nil {
			return result.Left[A](err)
		}
		return result.Of(a)
	}
}

//go:inline
func ChainFirstResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A] {
	return function.Bind2nd(MonadChainFirstResultK[A, B], f)
}

// MonadChainFirstIOK runs [IO] the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[A, A],
		MonadMap[B, A],
		FromIO[B],
		ma,
		f,
	)
}

//go:inline
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return fromio.ChainFirstIOK(
		Chain[A, A],
		Map[B, A],
		FromIO[B],
		f,
	)
}

//go:inline
func MonadTapEitherK[A, B any](ma IOResult[A], f either.Kleisli[error, A, B]) IOResult[A] {
	return MonadChainFirstEitherK(ma, f)
}

//go:inline
func TapEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return ChainFirstEitherK(f)
}

//go:inline
func MonadTapResultK[A, B any](ma IOResult[A], f result.Kleisli[A, B]) IOResult[A] {
	return MonadChainFirstResultK(ma, f)
}

//go:inline
func TapResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A] {
	return ChainFirstResultK(f)
}

//go:inline
func MonadTapIOK[A, B any](ma IOResult[A], f io.Kleisli[A, B]) IOResult[A] {
	return MonadChainFirstIOK(ma, f)
}

//go:inline
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return ChainFirstIOK(f)
}

// MonadFold handles both error and success cases explicitly, converting to an IO.
// This is useful for pattern matching on the IOResult.
func MonadFold[A, B any](ma IOResult[A], onLeft func(error) IO[B], onRight io.Kleisli[A, B]) IO[B] {
	return func() B {
		a, err := ma()
		if err != nil {
			return onLeft(err)()
		}
		return onRight(a)()
	}
}

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
// WithResource constructs a bracket pattern for resource management.
// It ensures resources are properly acquired, used, and released even if errors occur.
// The release function is always called, similar to defer.
func WithResource[A, R, ANY any](
	onCreate IOResult[R],
	onRelease Kleisli[R, ANY],
) Kleisli[Kleisli[R, A], A] {
	return func(k Kleisli[R, A]) IOResult[A] {
		return func() (A, error) {
			r, rerr := onCreate()
			if rerr != nil {
				return result.Left[A](rerr)
			}
			a, aerr := k(r)()
			_, nerr := onRelease(r)()
			if aerr != nil {
				return result.Left[A](aerr)
			}
			if nerr != nil {
				return result.Left[A](nerr)
			}
			return result.Of(a)
		}
	}
}

// FromImpure converts an impure side-effecting function into an IOResult.
// The function is executed when the IOResult runs, and always succeeds with nil.
func FromImpure(f func()) IOResult[any] {
	return function.Pipe2(
		f,
		io.FromImpure,
		FromIO[any],
	)
}

// Defer defers the creation of an IOResult until it is executed.
// This allows lazy evaluation of the IOResult itself.
func Defer[A any](gen Lazy[IOResult[A]]) IOResult[A] {
	return func() (A, error) {
		return gen()()
	}
}

// MonadAlt tries the first IOResult, and if it fails, tries the second.
// This provides a fallback mechanism for error recovery.
func MonadAlt[A any](first IOResult[A], second Lazy[IOResult[A]]) IOResult[A] {
	return func() (A, error) {
		a, err := first()
		if err != nil {
			return second()()
		}
		return result.Of(a)
	}
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second Lazy[IOResult[A]]) Operator[A, A] {
	return function.Bind2nd(MonadAlt[A], second)
}

//go:inline
func MonadFlap[B, A any](fab IOResult[func(A) B], a A) IOResult[B] {
	return functor.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return functor.Flap(Map[func(A) B, B], a)
}

// Delay creates an operation that passes in the value after some delay
// Delay creates an operator that delays execution by the specified duration.
// The IOResult is executed after waiting for the given duration.
func Delay[A any](delay time.Duration) Operator[A, A] {
	return func(fa IOResult[A]) IOResult[A] {
		return func() (A, error) {
			time.Sleep(delay)
			return fa()
		}
	}
}

// After creates an operation that passes after the given [time.Time]
// After creates an operator that delays execution until the specified timestamp.
// If the timestamp is in the past, the IOResult executes immediately.
func After[A any](timestamp time.Time) Operator[A, A] {
	return func(fa IOResult[A]) IOResult[A] {
		return func() (A, error) {
			// check if we need to wait
			current := time.Now()
			if current.Before(timestamp) {
				time.Sleep(timestamp.Sub(current))
			}
			return fa()
		}
	}
}

// MonadChainLeft handles the error case by chaining to a new computation.
// If the IOResult succeeds, it passes through unchanged.
func MonadChainLeft[A any](fa IOResult[A], f Kleisli[error, A]) IOResult[A] {
	return func() (A, error) {
		a, err := fa()
		if err != nil {
			return f(err)()
		}
		return result.Of(a)
	}
}

//go:inline
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A] {
	return function.Bind2nd(MonadChainLeft[A], f)
}

// MonadChainFirstLeft runs a computation on the error but always returns the original error.
// This is useful for side effects like logging errors without recovery.
func MonadChainFirstLeft[A, B any](ma IOResult[A], f Kleisli[error, B]) IOResult[A] {
	return func() (A, error) {
		a, err := ma()
		if err != nil {
			_, _ = f(err)()
			return result.Left[A](err)
		}
		return result.Of(a)
	}
}

//go:inline
func MonadTapLeft[A, B any](ma IOResult[A], f Kleisli[error, B]) IOResult[A] {
	return MonadChainFirstLeft(ma, f)
}

func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return function.Bind2nd(MonadChainFirstLeft[A, B], f)
}

//go:inline
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return ChainFirstLeft[A](f)
}
