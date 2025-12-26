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

package readerresult

import (
	"context"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
)

func FromReader[A any](r Reader[context.Context, A]) ReaderResult[A] {
	return readereither.FromReader[error](r)
}

func FromEither[A any](e Either[A]) ReaderResult[A] {
	return readereither.FromEither[context.Context](e)
}

func Left[A any](l error) ReaderResult[A] {
	return readereither.Left[context.Context, A](l)
}

func Right[A any](r A) ReaderResult[A] {
	return readereither.Right[context.Context, error](r)
}

func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B] {
	return readereither.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) Operator[A, B] {
	return readereither.Map[context.Context, error](f)
}

func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B] {
	return readereither.MonadChain(ma, F.Flow2(f, WithContext))
}

func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return readereither.Chain(F.Flow2(f, WithContext))
}

func Of[A any](a A) ReaderResult[A] {
	return readereither.Of[context.Context, error](a)
}

func MonadAp[A, B any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B] {
	return readereither.MonadAp(fab, fa)
}

func Ap[A, B any](fa ReaderResult[A]) func(ReaderResult[func(A) B]) ReaderResult[B] {
	return readereither.Ap[B](fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return readereither.FromPredicate[context.Context](pred, onFalse)
}

// OrElse recovers from a Left (error) by providing an alternative computation with access to context.Context.
// If the ReaderResult is Right, it returns the value unchanged.
// If the ReaderResult is Left, it applies the provided function to the error value,
// which returns a new ReaderResult that replaces the original.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations
// that need access to the context (for cancellation, deadlines, or values).
//
// Example:
//
//	// Recover with context-aware fallback
//	recover := readerresult.OrElse(func(err error) readerresult.ReaderResult[int] {
//	    if err.Error() == "not found" {
//	        return func(ctx context.Context) result.Result[int] {
//	            // Could check ctx.Err() here
//	            return result.Of(42)
//	        }
//	    }
//	    return readerresult.Left[int](err)
//	})
//
//go:inline
func OrElse[A any](onLeft Kleisli[error, A]) Kleisli[ReaderResult[A], A] {
	return readereither.OrElse(F.Flow2(onLeft, WithContext))
}

func Ask() ReaderResult[context.Context] {
	return readereither.Ask[context.Context, error]()
}

func MonadChainEitherK[A, B any](ma ReaderResult[A], f func(A) Either[B]) ReaderResult[B] {
	return readereither.MonadChainEitherK(ma, f)
}

func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderResult[A]) ReaderResult[B] {
	return readereither.ChainEitherK[context.Context](f)
}

func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B] {
	return readereither.ChainOptionK[context.Context, A, B](onNone)
}

func MonadFlap[B, A any](fab ReaderResult[func(A) B], a A) ReaderResult[B] {
	return readereither.MonadFlap(fab, a)
}

func Flap[B, A any](a A) Operator[func(A) B, B] {
	return readereither.Flap[context.Context, error, B](a)
}

//go:inline
func Read[A any](r context.Context) func(ReaderResult[A]) Result[A] {
	return readereither.Read[error, A](r)
}

// MonadMapTo executes a ReaderResult computation, discards its success value, and returns a constant value.
// This is the monadic version that takes both the ReaderResult and the constant value as parameters.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MonadMapTo WILL
// execute the original ReaderResult to allow any side effects to occur, then discard the success result
// and return the constant value. If the original computation fails, the error is preserved.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The type of the constant value to return on success
//
// Parameters:
//   - ma: The ReaderResult to execute (side effects will occur, success value discarded)
//   - b: The constant value to return if ma succeeds
//
// Returns:
//   - A ReaderResult that executes ma, preserves errors, but replaces success values with b
//
// Example:
//
//	type Config struct { Counter int }
//	increment := func(ctx context.Context) result.Result[int] {
//	    // Side effect: log the operation
//	    fmt.Println("incrementing")
//	    return result.Of(5)
//	}
//	r := readerresult.MonadMapTo(increment, "done")
//	result := r(context.Background()) // Prints "incrementing", returns Right("done")
//
//go:inline
func MonadMapTo[A, B any](ma ReaderResult[A], b B) ReaderResult[B] {
	return MonadMap(ma, reader.Of[A](b))
}

// MapTo creates an operator that executes a ReaderResult computation, discards its success value,
// and returns a constant value. This is the curried version where the constant value is provided first,
// returning a function that can be applied to any ReaderResult.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MapTo WILL
// execute the input ReaderResult to allow any side effects to occur, then discard the success result
// and return the constant value. If the computation fails, the error is preserved.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult (will be discarded if successful)
//   - B: The type of the constant value to return on success
//
// Parameters:
//   - b: The constant value to return on success
//
// Returns:
//   - An Operator that executes a ReaderResult[A], preserves errors, but replaces success with b
//
// Example:
//
//	logStep := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("step executed")
//	    return result.Of(42)
//	}
//	toDone := readerresult.MapTo[int, string]("done")
//	pipeline := toDone(logStep)
//	result := pipeline(context.Background()) // Prints "step executed", returns Right("done")
//
// Example - In a functional pipeline:
//
//	step1 := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("processing")
//	    return result.Of(1)
//	}
//	pipeline := F.Pipe1(
//	    step1,
//	    readerresult.MapTo[int, string]("complete"),
//	)
//	output := pipeline(context.Background()) // Prints "processing", returns Right("complete")
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return Map(reader.Of[A](b))
}

// MonadChainTo sequences two ReaderResult computations where the second ignores the first's success value.
// This is the monadic version that takes both ReaderResults as parameters.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MonadChainTo WILL
// execute the first ReaderResult to allow any side effects to occur, then discard the success result and
// execute the second ReaderResult with the same context. If the first computation fails, the error is
// returned immediately without executing the second computation.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - ma: The first ReaderResult to execute (side effects will occur, success value discarded)
//   - b: The second ReaderResult to execute if ma succeeds
//
// Returns:
//   - A ReaderResult that executes ma, then b if ma succeeds, returning b's result
//
// Example:
//
//	logStart := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("starting")
//	    return result.Of(1)
//	}
//	logEnd := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("ending")
//	    return result.Of("done")
//	}
//	r := readerresult.MonadChainTo(logStart, logEnd)
//	result := r(context.Background()) // Prints "starting" then "ending", returns Right("done")
//
//go:inline
func MonadChainTo[A, B any](ma ReaderResult[A], b ReaderResult[B]) ReaderResult[B] {
	return MonadChain(ma, reader.Of[A](b))
}

// ChainTo creates an operator that sequences two ReaderResult computations where the second ignores
// the first's success value. This is the curried version where the second ReaderResult is provided first,
// returning a function that can be applied to any first ReaderResult.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, ChainTo WILL
// execute the first ReaderResult to allow any side effects to occur, then discard the success result and
// execute the second ReaderResult with the same context. If the first computation fails, the error is
// returned immediately without executing the second computation.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - b: The second ReaderResult to execute after the first succeeds
//
// Returns:
//   - An Operator that executes the first ReaderResult, then b if successful
//
// Example:
//
//	logEnd := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("ending")
//	    return result.Of("done")
//	}
//	thenLogEnd := readerresult.ChainTo[int, string](logEnd)
//
//	logStart := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("starting")
//	    return result.Of(1)
//	}
//	pipeline := thenLogEnd(logStart)
//	result := pipeline(context.Background()) // Prints "starting" then "ending", returns Right("done")
//
// Example - In a functional pipeline:
//
//	step1 := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("step 1")
//	    return result.Of(1)
//	}
//	step2 := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("step 2")
//	    return result.Of("complete")
//	}
//	pipeline := F.Pipe1(
//	    step1,
//	    readerresult.ChainTo[int, string](step2),
//	)
//	output := pipeline(context.Background()) // Prints "step 1" then "step 2", returns Right("complete")
//
//go:inline
func ChainTo[A, B any](b ReaderResult[B]) Operator[A, B] {
	return Chain(reader.Of[A](b))
}

//go:inline
func MonadChainFirst[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[A] {
	return chain.MonadChainFirst(
		MonadChain,
		MonadMap,
		ma,
		F.Flow2(f, WithContext),
	)
}

//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(
		Chain,
		Map,
		F.Flow2(f, WithContext),
	)
}
