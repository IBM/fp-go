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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

// FromIOResult lifts an IOResult into a ReaderIOResult context.
// The resulting computation ignores the environment parameter and directly executes the IOResult.
//
// Type Parameters:
//   - R: The type of the environment (ignored by the computation)
//   - A: The type of the success value
//
// Parameters:
//   - ma: The IOResult to lift
//
// Returns:
//   - A ReaderIOResult that executes the IOResult regardless of the environment
//
// Example:
//
//	ioResult := func() (int, error) { return 42, nil }
//	readerIOResult := FromIOResult[Config](ioResult)
//	result, err := readerIOResult(cfg)() // Returns 42, nil
//
//go:inline
func FromIOResult[R, A any](ma IOResult[A]) ReaderIOResult[R, A] {
	return reader.Of[R](ma)
}

// RightIO lifts an IO computation into a ReaderIOResult as a successful value.
// The IO computation always succeeds, so it's wrapped in the Right (success) side.
//
// Type Parameters:
//   - R: The type of the environment (ignored by the computation)
//   - A: The type of the value produced by the IO
//
// Parameters:
//   - ma: The IO computation to lift
//
// Returns:
//   - A ReaderIOResult that executes the IO and wraps the result as a success
//
// Example:
//
//	getCurrentTime := func() time.Time { return time.Now() }
//	readerIOResult := RightIO[Config](getCurrentTime)
//	result, err := readerIOResult(cfg)() // Returns current time, nil
func RightIO[R, A any](ma IO[A]) ReaderIOResult[R, A] {
	return function.Pipe2(ma, ioresult.RightIO[A], FromIOResult[R, A])
}

// LeftIO lifts an IO computation that produces an error into a ReaderIOResult as a failure.
// The IO computation produces an error, which is wrapped in the Left (error) side.
//
// Type Parameters:
//   - R: The type of the environment (ignored by the computation)
//   - A: The type of the success value (never produced)
//
// Parameters:
//   - ma: The IO computation that produces an error
//
// Returns:
//   - A ReaderIOResult that executes the IO and wraps the error as a failure
//
// Example:
//
//	getError := func() error { return errors.New("something went wrong") }
//	readerIOResult := LeftIO[Config, int](getError)
//	_, err := readerIOResult(cfg)() // Returns error
func LeftIO[R, A any](ma IO[error]) ReaderIOResult[R, A] {
	return function.Pipe2(ma, ioresult.LeftIO[A], FromIOResult[R, A])
}

// FromIO lifts an IO computation into a ReaderIOResult context.
// This is an alias for RightIO - the IO computation always succeeds.
//
// Type Parameters:
//   - R: The type of the environment (ignored by the computation)
//   - E: Unused type parameter (kept for compatibility)
//   - A: The type of the value produced by the IO
//
// Parameters:
//   - ma: The IO computation to lift
//
// Returns:
//   - A ReaderIOResult that executes the IO and wraps the result as a success
//
//go:inline
func FromIO[R, E, A any](ma IO[A]) ReaderIOResult[R, A] {
	return RightIO[R](ma)
}

// FromReaderIO lifts a ReaderIO into a ReaderIOResult context.
// The ReaderIO computation always succeeds, so it's wrapped in the Right (success) side.
// This is an alias for RightReaderIO.
//
// Type Parameters:
//   - R: The type of the environment
//   - A: The type of the value produced
//
// Parameters:
//   - ma: The ReaderIO to lift
//
// Returns:
//   - A ReaderIOResult that executes the ReaderIO and wraps the result as a success
//
// Example:
//
//	getConfigValue := func(cfg Config) func() int {
//	    return func() int { return cfg.Timeout }
//	}
//	readerIOResult := FromReaderIO(getConfigValue)
//	result, err := readerIOResult(cfg)() // Returns cfg.Timeout, nil
//
//go:inline
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A] {
	return RightReaderIO(ma)
}

// RightReaderIO lifts a ReaderIO into a ReaderIOResult as a successful value.
// The ReaderIO computation always succeeds, so it's wrapped in the Right (success) side.
//
// Type Parameters:
//   - R: The type of the environment
//   - A: The type of the value produced
//
// Parameters:
//   - ma: The ReaderIO to lift
//
// Returns:
//   - A ReaderIOResult that executes the ReaderIO and wraps the result as a success
//
// Example:
//
//	logMessage := func(cfg Config) func() string {
//	    return func() string {
//	        log.Printf("Processing with timeout: %d", cfg.Timeout)
//	        return "logged"
//	    }
//	}
//	readerIOResult := RightReaderIO(logMessage)
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A] {
	return function.Flow2(
		ma,
		ioresult.FromIO,
	)
}

// MonadMap transforms the success value of a ReaderIOResult using the provided function.
// If the computation fails, the error is propagated unchanged.
//
// Type Parameters:
//   - R: The type of the environment
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - fa: The ReaderIOResult to transform
//   - f: The transformation function
//
// Returns:
//   - A new ReaderIOResult with the transformed value
//
// Example:
//
//	getValue := Right[Config](10)
//	doubled := MonadMap(getValue, func(x int) int { return x * 2 })
//	result, err := doubled(cfg)() // Returns 20, nil
func MonadMap[R, A, B any](fa ReaderIOResult[R, A], f func(A) B) ReaderIOResult[R, B] {
	return function.Flow2(
		fa,
		ioresult.Map(f),
	)
}

// Map transforms the success value of a ReaderIOResult using the provided function.
// This is the curried version of MonadMap, useful for composition in pipelines.
//
// Type Parameters:
//   - R: The type of the environment
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - f: The transformation function
//
// Returns:
//   - A function that transforms a ReaderIOResult
//
// Example:
//
//	double := Map[Config](func(x int) int { return x * 2 })
//	getValue := Right[Config](10)
//	result := F.Pipe1(getValue, double)
//	value, err := result(cfg)() // Returns 20, nil
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	mp := ioresult.Map(f)
	return func(ri ReaderIOResult[R, A]) ReaderIOResult[R, B] {
		return function.Flow2(
			ri,
			mp,
		)
	}
}

// MonadMapTo replaces the success value with a constant value.
// Useful when you want to discard the result but keep the effect.
func MonadMapTo[R, A, B any](fa ReaderIOResult[R, A], b B) ReaderIOResult[R, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

// MapTo returns a function that replaces the success value with a constant.
// This is the curried version of MonadMapTo.
func MapTo[R, A, B any](b B) Operator[R, A, B] {
	return Map[R](function.Constant1[A](b))
}

//go:inline
func MonadChain[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return func(r R) IOResult[B] {
		return function.Pipe1(
			fa(r),
			ioresult.Chain(func(a A) IOResult[B] {
				return f(a)(r)
			}),
		)
	}
}

//go:inline
func MonadChainFirst[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return chain.MonadChainFirst(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		fa,
		f)
}

//go:inline
func MonadTap[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirst(fa, f)
}

// MonadChainEitherK chains a computation that returns an Either into a ReaderIOResult.
// The Either is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, A, B],
		FromEither[R, B],
		ma,
		f,
	)
}

// ChainEitherK returns a function that chains an Either-returning function into ReaderIOResult.
// This is the curried version of MonadChainEitherK.
//
//go:inline
func ChainEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, A, B],
		FromEither[R, B],
		f,
	)
}

// MonadChainFirstEitherK chains an Either-returning computation but keeps the original value.
// Useful for validation or side effects that return Either.
//
//go:inline
func MonadChainFirstEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromEither[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapEitherK[R, A, B any](ma ReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, f)
}

// ChainFirstEitherK returns a function that chains an Either computation while preserving the original value.
// This is the curried version of MonadChainFirstEitherK.
//
//go:inline
func ChainFirstEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, A, A],
		Map[R, B, A],
		FromEither[R, B],
		f,
	)
}

//go:inline
func TapEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return ChainFirstEitherK[R](f)
}

// MonadChainReaderK chains a Reader-returning computation into a ReaderIOResult.
// The Reader is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReader[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstReaderK(ma, f)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReader[R, B],
		f,
	)
}

//go:inline
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderK(f)
}

//go:inline
func MonadChainReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

//go:inline
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderIOK(f)
}

// //go:inline
// func MonadChainReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, A, B]) ReaderIOResult[R, B] {
// 	return fromreader.MonadChainReaderK(
// 		MonadChain[R, A, B],
// 		FromReaderEither[R, B],
// 		ma,
// 		f,
// 	)
// }

// // ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// // This is the curried version of MonadChainReaderK.
// //
// //go:inline
// func ChainReaderEitherK[R, A, B any](f RE.Kleisli[R, A, B]) Operator[R, A, B] {
// 	return fromreader.ChainReaderK(
// 		Chain[R, A, B],
// 		FromReaderEither[R, B],
// 		f,
// 	)
// }

// //go:inline
// func MonadChainFirstReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, A, B]) ReaderIOResult[R, A] {
// 	return fromreader.MonadChainFirstReaderK(
// 		MonadChainFirst[R, A, B],
// 		FromReaderEither[R, B],
// 		ma,
// 		f,
// 	)
// }

// //go:inline
// func MonadTapReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, A, B]) ReaderIOResult[R, A] {
// 	return MonadChainFirstReaderEitherK(ma, f)
// }

// // ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// // This is the curried version of MonadChainReaderK.
// //
// //go:inline
// func ChainFirstReaderEitherK[R, A, B any](f RE.Kleisli[R, A, B]) Operator[R, A, A] {
// 	return fromreader.ChainFirstReaderK(
// 		ChainFirst[R, A, B],
// 		FromReaderEither[R, B],
// 		f,
// 	)
// }

// //go:inline
// func TapReaderEitherK[R, A, B any](f RE.Kleisli[R, A, B]) Operator[R, A, A] {
// 	return ChainFirstReaderEitherK(f)
// }

// //go:inline
// func ChainReaderOptionK[R, A, B any](onNone func() E) func(readeroption.Kleisli[R, A, B]) Operator[R, A, B] {
// 	fro := FromReaderOption[R, B](onNone)
// 	return func(f readeroption.Kleisli[R, A, B]) Operator[R, A, B] {
// 		return fromreader.ChainReaderK(
// 			Chain[R, A, B],
// 			fro,
// 			f,
// 		)
// 	}
// }

// //go:inline
// func ChainFirstReaderOptionK[R, A, B any](onNone func() E) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
// 	fro := FromReaderOption[R, B](onNone)
// 	return func(f readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
// 		return fromreader.ChainFirstReaderK(
// 			ChainFirst[R, A, B],
// 			fro,
// 			f,
// 		)
// 	}
// }

// //go:inline
// func TapReaderOptionK[R, A, B any](onNone func() E) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
// 	return ChainFirstReaderOptionK[R, A, B](onNone)
// }

// // MonadChainIOEitherK chains an IOEither-returning computation into a ReaderIOResult.
// // The IOEither is automatically lifted into the ReaderIOResult context.
// //
// //go:inline
// func MonadChainIOEitherK[R, A, B any](ma ReaderIOResult[R, A], f IOE.Kleisli[A, B]) ReaderIOResult[R, B] {
// 	return fromioeither.MonadChainIOEitherK(
// 		MonadChain[R, A, B],
// 		FromIOEither[R, B],
// 		ma,
// 		f,
// 	)
// }

// // ChainIOEitherK returns a function that chains an IOEither-returning function into ReaderIOResult.
// // This is the curried version of MonadChainIOEitherK.
// //
// //go:inline
// func ChainIOEitherK[R, A, B any](f IOE.Kleisli[A, B]) Operator[R, A, B] {
// 	return fromioeither.ChainIOEitherK(
// 		Chain[R, A, B],
// 		FromIOEither[R, B],
// 		f,
// 	)
// }

// MonadChainIOK chains an IO-returning computation into a ReaderIOResult.
// The IO is automatically lifted into the ReaderIOResult context (always succeeds).
//
//go:inline
func MonadChainIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma,
		f,
	)
}

// ChainIOK returns a function that chains an IO-returning function into ReaderIOResult.
// This is the curried version of MonadChainIOK.
//
//go:inline
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
}

// MonadChainFirstIOK chains an IO computation but keeps the original value.
// Useful for performing IO side effects while preserving the original value.
//
//go:inline
func MonadChainFirstIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromIO[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapIOK[R, A, B any](ma ReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstIOK(ma, f)
}

// ChainFirstIOK returns a function that chains an IO computation while preserving the original value.
// This is the curried version of MonadChainFirstIOK.
//
//go:inline
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, A, A],
		Map[R, B, A],
		FromIO[R, B],
		f,
	)
}

//go:inline
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOK[R](f)
}

// // ChainOptionK returns a function that chains an Option-returning function into ReaderIOResult.
// // If the Option is None, the provided error function is called to produce the error value.
// //
// //go:inline
// func ChainOptionK[R, A, B any](onNone func() E) func(func(A) Option[B]) Operator[R, A, B] {
// 	return fromeither.ChainOptionK(
// 		MonadChain[R, A, B],
// 		FromEither[R, B],
// 		onNone,
// 	)
// }

// MonadAp applies a function wrapped in a context to a value wrapped in a context.
// Both computations are executed (default behavior may be sequential or parallel depending on implementation).
//
//go:inline
func MonadAp[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return func(r R) IOResult[B] {
		return ioresult.MonadAp(fab(r), fa(r))
	}
}

// MonadApSeq applies a function in a context to a value in a context, executing them sequentially.
//
//go:inline
func MonadApSeq[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return func(r R) IOResult[B] {
		return ioresult.MonadApSeq(fab(r), fa(r))
	}
}

// MonadApPar applies a function in a context to a value in a context, executing them in parallel.
//
//go:inline
func MonadApPar[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return func(r R) IOResult[B] {
		return ioresult.MonadApPar(fab(r), fa(r))
	}
}

// Ap returns a function that applies a function in a context to a value in a context.
// This is the curried version of MonadAp.
func Ap[B, R, A any](fa ReaderIOResult[R, A]) func(fab ReaderIOResult[R, func(A) B]) ReaderIOResult[R, B] {
	return function.Bind2nd(MonadAp[R, A, B], fa)
}

// Chain returns a function that sequences computations where the second depends on the first.
// This is the curried version of MonadChain.
//
//go:inline
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return function.Bind2nd(MonadChain, f)
}

// ChainFirst returns a function that sequences computations but keeps the first result.
// This is the curried version of MonadChainFirst.
//
//go:inline
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return chain.ChainFirst(
		Chain[R, A, A],
		Map[R, B, A],
		f)
}

//go:inline
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirst(f)
}

// Right creates a successful ReaderIOResult with the given value.
//
//go:inline
func Right[R, A any](a A) ReaderIOResult[R, A] {
	return reader.Of[R](ioresult.Right(a))
}

// Left creates a failed ReaderIOResult with the given error.
//
//go:inline
func Left[R, A any](e error) ReaderIOResult[R, A] {
	return reader.Of[R](ioresult.Left[A](e))
}

// Of creates a successful ReaderIOResult with the given value.
// This is the pointed functor operation, lifting a pure value into the ReaderIOResult context.
func Of[R, A any](a A) ReaderIOResult[R, A] {
	return Right[R](a)
}

// Flatten removes one level of nesting from a nested ReaderIOResult.
// Converts ReaderIOResult[R, ReaderIOResult[R, A]] to ReaderIOResult[R, A].
func Flatten[R, A any](mma ReaderIOResult[R, ReaderIOResult[R, A]]) ReaderIOResult[R, A] {
	return MonadChain(mma, function.Identity[ReaderIOResult[R, A]])
}

// FromEither lifts an Either into a ReaderIOResult context.
// The Either value is independent of any context or IO effects.
func FromEither[R, A any](t either.Either[error, A]) ReaderIOResult[R, A] {
	return func(r R) IOResult[A] {
		return func() (A, error) {
			return either.Unwrap(t)
		}
	}
}

// RightReader lifts a Reader into a ReaderIOResult, placing the result in the Right side.
func RightReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A] {
	return function.Flow2(ma, ioresult.Right[A])
}

// LeftReader lifts a Reader into a ReaderIOResult, placing the result in the Left (error) side.
func LeftReader[A, R any](ma Reader[R, error]) ReaderIOResult[R, A] {
	return function.Flow2(ma, ioresult.Left[A])
}

// FromReader lifts a Reader into a ReaderIOResult context.
// The Reader result is placed in the Right side (success).
func FromReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A] {
	return RightReader(ma)
}

// // FromIOEither lifts an IOEither into a ReaderIOResult context.
// // The computation becomes independent of any reader context.
// //
// //go:inline
// func FromIOEither[R, A any](ma IOEither[A]) ReaderIOResult[R, A] {
// 	return reader.Of[R](ma)
// }

// // FromReaderEither lifts a ReaderEither into a ReaderIOResult context.
// // The Either result is lifted into an IO effect.
// func FromReaderEither[R, A any](ma RE.ReaderEither[R, A]) ReaderIOResult[R, A] {
// 	return function.Flow2(ma, IOE.FromEither[A])
// }

// Ask returns a ReaderIOResult that retrieves the current context.
// Useful for accessing configuration or dependencies.
//
//go:inline
func Ask[R any]() ReaderIOResult[R, R] {
	return fromreader.Ask(FromReader[R, R])()
}

// Asks returns a ReaderIOResult that retrieves a value derived from the context.
// This is useful for extracting specific fields from a configuration object.
//
//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderIOResult[R, A] {
	return fromreader.Asks(FromReader[R, A])(r)
}

// // FromOption converts an Option to a ReaderIOResult.
// // If the Option is None, the provided function is called to produce the error.
// //
// //go:inline
// func FromOption[R, A any](onNone func() E) func(Option[A]) ReaderIOResult[R, A] {
// 	return fromeither.FromOption(FromEither[R, A], onNone)
// }

// // FromPredicate creates a ReaderIOResult from a predicate.
// // If the predicate returns false, the onFalse function is called to produce the error.
// //
// //go:inline
// func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOResult[R, A] {
// 	return fromeither.FromPredicate(FromEither[R, A], pred, onFalse)
// }

// // Fold handles both success and error cases, producing a ReaderIO.
// // This is useful for converting a ReaderIOResult into a ReaderIO by handling all cases.
// //
// //go:inline
// func Fold[R, A, B any](onLeft func(E) ReaderIO[R, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOResult[R, A]) ReaderIO[R, B] {
// 	return eithert.MatchE(readerio.MonadChain[R, either.Either[A], B], onLeft, onRight)
// }

// //go:inline
// func MonadFold[R, A, B any](ma ReaderIOResult[R, A], onLeft func(E) ReaderIO[R, B], onRight func(A) ReaderIO[R, B]) ReaderIO[R, B] {
// 	return eithert.FoldE(readerio.MonadChain[R, either.Either[A], B], ma, onLeft, onRight)
// }

// // GetOrElse provides a default value in case of error.
// // The default is computed lazily via a ReaderIO.
// //
// //go:inline
// func GetOrElse[R, A any](onLeft func(E) ReaderIO[R, A]) func(ReaderIOResult[R, A]) ReaderIO[R, A] {
// 	return eithert.GetOrElse(readerio.MonadChain[R, either.Either[A], A], readerio.Of[R, A], onLeft)
// }

// // OrElse tries an alternative computation if the first one fails.
// // The alternative can produce a different error type.
// //
// //go:inline
// func OrElse[R1, A2 any](onLeft func(E1) ReaderIOResult[R2, A]) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A] {
// 	return eithert.OrElse(readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]], readerio.Of[R, either.Either[E2, A]], onLeft)
// }

// // OrLeft transforms the error using a ReaderIO if the computation fails.
// // The success value is preserved unchanged.
// //
// //go:inline
// func OrLeft[A1, R2 any](onLeft func(E1) ReaderIO[R2]) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A] {
// 	return eithert.OrLeft(
// 		readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]],
// 		readerio.MonadMap[R2, either.Either[E2, A]],
// 		readerio.Of[R, either.Either[E2, A]],
// 		onLeft,
// 	)
// }

// // MonadBiMap applies two functions: one to the error, one to the success value.
// // This allows transforming both channels simultaneously.
// //
// //go:inline
// func MonadBiMap[R12, A, B any](fa ReaderIOResult[R1, A], f func(E1) E2, g func(A) B) ReaderIOResult[R2, B] {
// 	return eithert.MonadBiMap(
// 		readerio.MonadMap[R, either.Either[E1, A], either.Either[E2, B]],
// 		fa, f, g,
// 	)
// }

// // BiMap returns a function that maps over both the error and success channels.
// // This is the curried version of MonadBiMap.
// //
// //go:inline
// func BiMap[R12, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, B] {
// 	return eithert.BiMap(readerio.Map[R, either.Either[E1, A], either.Either[E2, B]], f, g)
// }

// // TryCatch wraps a function that returns (value, error) into a ReaderIOResult.
// // The onThrow function converts the error into the desired error type.
// func TryCatch[R, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOResult[R, A] {
// 	return func(r R) IOEither[A] {
// 		return IOE.TryCatch(f(r), onThrow)
// 	}
// }

// // MonadAlt tries the first computation, and if it fails, tries the second.
// // This implements the Alternative pattern for error recovery.
// //
// //go:inline
// func MonadAlt[R, A any](first ReaderIOResult[R, A], second L.Lazy[ReaderIOResult[R, A]]) ReaderIOResult[R, A] {
// 	return eithert.MonadAlt(
// 		readerio.Of[Rither[A]],
// 		readerio.MonadChain[Rither[A]ither[A]],

// 		first,
// 		second,
// 	)
// }

// // Alt returns a function that tries an alternative computation if the first fails.
// // This is the curried version of MonadAlt.
// //
// //go:inline
// func Alt[R, A any](second L.Lazy[ReaderIOResult[R, A]]) Operator[R, A, A] {
// 	return eithert.Alt(
// 		readerio.Of[Rither[A]],
// 		readerio.Chain[Rither[A]ither[A]],

// 		second,
// 	)
// }

// // Memoize computes the value of the ReaderIOResult lazily but exactly once.
// // The context used is from the first call. Do not use if the value depends on the context.
// //
// //go:inline
// func Memoize[
// 	R, A any](rdr ReaderIOResult[R, A]) ReaderIOResult[R, A] {
// 	return readerio.Memoize(rdr)
// }

// MonadFlap applies a value to a function wrapped in a context.
// This is the reverse of Ap - the value is fixed and the function varies.
//
//go:inline
func MonadFlap[R, B, A any](fab ReaderIOResult[R, func(A) B], a A) ReaderIOResult[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap returns a function that applies a fixed value to a function in a context.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[R, B, A any](a A) func(ReaderIOResult[R, func(A) B]) ReaderIOResult[R, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

// // MonadMapLeft applies a function to the error value, leaving success unchanged.
// //
// //go:inline
// func MonadMapLeft[R12, A any](fa ReaderIOResult[R1, A], f func(E1) E2) ReaderIOResult[R2, A] {
// 	return eithert.MonadMapLeft(readerio.MonadMap[Rither[E1, A]ither[E2, A]], fa, f)
// }

// // MapLeft returns a function that transforms the error channel.
// // This is the curried version of MonadMapLeft.
// //
// //go:inline
// func MapLeft[R, A12 any](f func(E1) E2) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A] {
// 	return eithert.MapLeft(readerio.Map[Rither[E1, A]ither[E2, A]], f)
// }

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
// This is similar to Contravariant's contramap operation.
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A] {
	return reader.Local[IOResult[A]](f)
}

// Read executes a ReaderIOResult by providing it with a concrete environment value.
// This function "runs" the reader computation by supplying the required environment,
// converting a ReaderIOResult into an IOResult that can be executed.
//
// This is the fundamental way to execute a ReaderIOResult computation - you provide
// the environment it needs, and get back an IOResult that can be run.
//
// Type Parameters:
//   - A: The type of the success value
//   - R: The type of the environment/context
//
// Parameters:
//   - r: The environment value to provide to the computation
//
// Returns:
//   - A function that takes a ReaderIOResult and returns an IOResult
//
// Example:
//
//	// Define a computation that needs configuration
//	computation := func(cfg Config) IOResult[string] {
//	    return func() (string, error) {
//	        return fmt.Sprintf("Value: %d", cfg.Value), nil
//	    }
//	}
//
//	// Provide the configuration and execute
//	cfg := Config{Value: 42}
//	result := Read[string](cfg)(computation)
//	value, err := result() // Returns "Value: 42", nil
//
//go:inline
func Read[A, R any](r R) func(ReaderIOResult[R, A]) IOResult[A] {
	return reader.Read[IOResult[A]](r)
}

// ReadIO executes a ReaderIOResult by providing it with an environment value wrapped in IO.
// This is useful when the environment itself needs to be computed or retrieved through an IO operation.
// The IO effect is executed first to obtain the environment, then that environment is provided
// to the ReaderIOResult computation.
//
// This allows for dynamic environment resolution where the configuration or context is not
// immediately available but must be computed or fetched.
//
// Type Parameters:
//   - A: The type of the success value
//   - R: The type of the environment/context
//
// Parameters:
//   - r: An IO operation that produces the environment value
//
// Returns:
//   - A function that takes a ReaderIOResult and returns an IOResult
//
// Example:
//
//	// Environment that needs to be loaded
//	loadConfig := func() Config {
//	    // Simulate loading config from file or environment
//	    return Config{Value: 42}
//	}
//
//	// Computation that needs the config
//	computation := func(cfg Config) IOResult[string] {
//	    return func() (string, error) {
//	        return fmt.Sprintf("Loaded: %d", cfg.Value), nil
//	    }
//	}
//
//	// Load config and execute computation
//	result := ReadIO[string](loadConfig)(computation)
//	value, err := result() // Loads config, then returns "Loaded: 42", nil
func ReadIO[A, R any](r IO[R]) func(ReaderIOResult[R, A]) IOResult[A] {
	return func(ri ReaderIOResult[R, A]) IOResult[A] {
		return func() (A, error) {
			return ri(r())()
		}
	}
}

// ReadIOResult executes a ReaderIOResult by providing it with an environment value wrapped in IOResult.
// This is the most flexible variant, allowing the environment itself to be the result of a computation
// that may fail. If the environment computation fails, the entire computation fails without executing
// the ReaderIOResult.
//
// This is useful when the environment must be validated, loaded from external sources, or computed
// in a way that might fail. The error from environment resolution is propagated as the final error.
//
// Type Parameters:
//   - A: The type of the success value
//   - R: The type of the environment/context
//
// Parameters:
//   - r: An IOResult operation that produces the environment value or an error
//
// Returns:
//   - A function that takes a ReaderIOResult and returns an IOResult
//
// Example:
//
//	// Environment that might fail to load
//	loadConfig := func() (Config, error) {
//	    cfg, err := os.ReadFile("config.json")
//	    if err != nil {
//	        return Config{}, fmt.Errorf("failed to load config: %w", err)
//	    }
//	    return parseConfig(cfg)
//	}
//
//	// Computation that needs the config
//	computation := func(cfg Config) IOResult[string] {
//	    return func() (string, error) {
//	        return fmt.Sprintf("Using: %d", cfg.Value), nil
//	    }
//	}
//
//	// Try to load config and execute computation
//	result := ReadIOResult[string](loadConfig)(computation)
//	value, err := result() // Returns error if config loading fails
//
//go:inline
func ReadIOResult[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A] {
	return function.Flow2(
		ioresult.Chain[R, A],
		Read[A](r),
	)
}

// //go:inline
// func MonadChainLeft[RAB, A any](fa ReaderIOResult[RA, A], f Kleisli[RBA, A]) ReaderIOResult[RB, A] {
// 	return readert.MonadChain(
// 		IOE.MonadChainLeft[EAB, A],
// 		fa,
// 		f,
// 	)
// }

// //go:inline
// func ChainLeft[RAB, A any](f Kleisli[RBA, A]) func(ReaderIOResult[RA, A]) ReaderIOResult[RB, A] {
// 	return readert.Chain[ReaderIOResult[RA, A]](
// 		IOE.ChainLeft[EAB, A],
// 		f,
// 	)
// }

// // MonadChainFirstLeft chains a computation on the left (error) side but always returns the original error.
// // If the input is a Left value, it applies the function f to the error and executes the resulting computation,
// // but always returns the original Left error regardless of what f returns (Left or Right).
// // If the input is a Right value, it passes through unchanged without calling f.
// //
// // This is useful for side effects on errors (like logging or metrics) where you want to perform an action
// // when an error occurs but always propagate the original error, ensuring the error path is preserved.
// //
// // Parameters:
// //   - ma: The input ReaderIOResult that may contain an error of type EA
// //   - f: A function that takes an error of type EA and returns a ReaderIOResult (typically for side effects)
// //
// // Returns:
// //   - A ReaderIOResult with the original error preserved if input was Left, or the original Right value
// //
// //go:inline
// func MonadChainFirstLeft[A, RAB, B any](ma ReaderIOResult[RA, A], f Kleisli[RBA, B]) ReaderIOResult[RA, A] {
// 	return MonadChainLeft(ma, function.Flow2(f, Fold(function.Constant1[EB](ma), function.Constant1[B](ma))))
// }

// //go:inline
// func MonadTapLeft[A, RAB, B any](ma ReaderIOResult[RA, A], f Kleisli[RBA, B]) ReaderIOResult[RA, A] {
// 	return MonadChainFirstLeft(ma, f)
// }

// // ChainFirstLeft is the curried version of [MonadChainFirstLeft].
// // It returns a function that chains a computation on the left (error) side while always preserving the original error.
// //
// // This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// // in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// // ensuring the error path is preserved.
// //
// // Parameters:
// //   - f: A function that takes an error of type EA and returns a ReaderIOResult (typically for side effects)
// //
// // Returns:
// //   - An Operator that performs the side effect but always returns the original error if input was Left
// //
// //go:inline
// func ChainFirstLeft[A, RAB, B any](f Kleisli[RBA, B]) Operator[RA, A, A] {
// 	return ChainLeft(func(e EA) ReaderIOResult[RA, A] {
// 		ma := Left[R, A](e)
// 		return MonadFold(f(e), function.Constant1[EB](ma), function.Constant1[B](ma))
// 	})
// }

// //go:inline
// func TapLeft[A, RAB, B any](f Kleisli[RBA, B]) Operator[RA, A, A] {
// 	return ChainFirstLeft[A](f)
// }
