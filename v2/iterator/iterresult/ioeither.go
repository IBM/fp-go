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

package iterresult

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/iterator/itereither"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Seq[A any]    = iter.Seq[A]
	Result[A any] = result.Result[A]

	SeqResult[A any] = Seq[Result[A]]

	Kleisli[A, B any]  = R.Reader[A, SeqResult[B]]
	Operator[A, B any] = Kleisli[SeqResult[A], B]
)

// Left constructs a SeqResult that represents a failure with an error value
func Left[A any](l error) SeqResult[A] {
	return itereither.Left[A](l)
}

// Right constructs a SeqResult that represents a successful computation with a value of type A
func Right[A any](r A) SeqResult[A] {
	return itereither.Right[error](r)
}

// Of constructs a SeqResult that represents a successful computation with a value of type A.
// This is an alias for Right and is the canonical way to lift a pure value into the SeqResult context.
func Of[A any](r A) SeqResult[A] {
	return itereither.Of[error](r)
}

// MonadOf is an alias for Of, provided for consistency with monad naming conventions
func MonadOf[A any](r A) SeqResult[A] {
	return itereither.MonadOf[error](r)
}

// LeftSeq constructs a SeqResult from a Seq that produces an error value
func LeftSeq[A any](ml Seq[error]) SeqResult[A] {
	return itereither.LeftSeq[A](ml)
}

// RightSeq constructs a SeqResult from a Seq that produces a success value
func RightSeq[A any](mr Seq[A]) SeqResult[A] {
	return itereither.RightSeq[error](mr)
}

// FromEither lifts a Result value into the SeqResult context
func FromEither[A any](e Result[A]) SeqResult[A] {
	return itereither.FromEither(e)
}

func FromOption[A any](onNone Lazy[error]) Kleisli[O.Option[A], A] {
	return itereither.FromOption[A](onNone)
}

func ChainOptionK[A, B any](onNone Lazy[error]) func(O.Kleisli[A, B]) Operator[A, B] {
	return itereither.ChainOptionK[A, B](onNone)
}

func MonadChainSeqK[A, B any](ma SeqResult[A], f iter.Kleisli[A, B]) SeqResult[B] {
	return itereither.MonadChainSeqK(ma, f)
}

func ChainSeqK[A, B any](f iter.Kleisli[A, B]) Operator[A, B] {
	return itereither.ChainSeqK[error](f)
}

func MonadMergeMapSeqK[A, B any](ma SeqResult[A], f iter.Kleisli[A, B]) SeqResult[B] {
	return itereither.MonadMergeMapSeqK(ma, f)
}

func MergeMapSeqK[A, B any](f iter.Kleisli[A, B]) Operator[A, B] {
	return itereither.MergeMapSeqK[error](f)
}

// FromSeq creates a SeqResult from a Seq instance, invoking Seq for each invocation of SeqResult
func FromSeq[A any](mr Seq[A]) SeqResult[A] {
	return itereither.FromSeq[error](mr)
}

// MonadMap applies a function to the value inside a successful SeqResult, leaving errors unchanged.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---R(3)---|
//	f(x) = x * 2
//	Output: ---R(2)---R(4)---L(e)---R(6)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func MonadMap[A, B any](fa SeqResult[A], f func(A) B) SeqResult[B] {
	return itereither.MonadMap(fa, f)
}

// Map returns a function that applies a transformation to the value inside a successful SeqResult.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---R(3)---|
//	f(x) = x * 2
//	Output: ---R(2)---R(4)---L(e)---R(6)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func Map[A, B any](f func(A) B) Operator[A, B] {
	return itereither.Map[error](f)
}

// MonadMapTo replaces the value inside a successful SeqResult with a constant value
func MonadMapTo[A, B any](fa SeqResult[A], b B) SeqResult[B] {
	return itereither.MonadMapTo(fa, b)
}

// MapTo returns a function that replaces the value inside a successful SeqResult with a constant value
func MapTo[A, B any](b B) Operator[A, B] {
	return itereither.MapTo[error, A](b)
}

// MonadChain sequences two SeqResult computations, where the second depends on the result of the first.
//
// Marble diagram:
//
//	Input:  ---R(1)-------R(2)---L(e)---|
//	f(1) -> ---R(10)---R(11)---|
//	f(2) -> ---R(20)---R(21)---|
//	Output: ---R(10)---R(11)---R(20)---R(21)---L(e)---|
//
// Each Right value is transformed into a sequence, which is then flattened.
// Left values pass through unchanged and stop further processing.
func MonadChain[A, B any](fa SeqResult[A], f Kleisli[A, B]) SeqResult[B] {
	return itereither.MonadChain(fa, f)
}

// MonadMergeMap sequences two SeqResult computations, where the second depends on the result of the first.
// Unlike MonadChain, MergeMap interleaves results from concurrent sequences.
//
// Marble diagram:
//
//	Input:  ---R(1)-------R(2)---|
//	f(1) -> ---R(10)------R(11)---|
//	f(2) -> ------R(20)------R(21)---|
//	Output: ---R(10)---R(20)---R(11)---R(21)---|
//
// Results are interleaved as they become available, rather than waiting for each sequence to complete.
func MonadMergeMap[A, B any](fa SeqResult[A], f Kleisli[A, B]) SeqResult[B] {
	return itereither.MonadMergeMap(fa, f)
}

// Chain returns a function that sequences two SeqResult computations
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return itereither.Chain(f)
}

// MergeMap returns a function that sequences two SeqResult computations
func MergeMap[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return itereither.MergeMap(f)
}

func MonadChainEitherK[A, B any](ma SeqResult[A], f result.Kleisli[A, B]) SeqResult[B] {
	return itereither.MonadChainEitherK(ma, f)
}

func ChainEitherK[A, B any](f result.Kleisli[A, B]) Operator[A, B] {
	return itereither.ChainEitherK(f)
}

// MonadAp applies a function wrapped in a SeqResult to a value wrapped in a SeqResult
func MonadAp[B, A any](mab SeqResult[func(A) B], ma SeqResult[A]) SeqResult[B] {
	return itereither.MonadAp(mab, ma)
}

// Ap applies a function wrapped in a SeqResult to a value wrapped in a SeqResult.
// This is an alias of ApPar which applies the function and value in parallel.
func Ap[B, A any](ma SeqResult[A]) Operator[func(A) B, B] {
	return itereither.Ap[B](ma)
}

// Flatten removes one level of nesting from a nested SeqResult
func Flatten[A any](mma SeqResult[SeqResult[A]]) SeqResult[A] {
	return itereither.Flatten(mma)
}

// MonadMapLeft applies a function to the error value of a failed SeqResult, leaving successful values unchanged.
//
// Marble diagram:
//
//	Input:  ---L(e1)---R(1)---L(e2)---R(2)---|
//	f(e) = "error: " + e
//	Output: ---L("error: e1")---R(1)---L("error: e2")---R(2)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func MonadMapLeft[A any](fa SeqResult[A], f Endomorphism[error]) SeqResult[A] {
	return itereither.MonadMapLeft(fa, f)
}

// MapLeft returns a function that applies a transformation to the error value of a failed SeqResult
func MapLeft[A any](f Endomorphism[error]) Operator[A, A] {
	return itereither.MapLeft[A](f)
}

// MonadBiMap applies one function to the error value and another to the success value of a SeqResult.
//
// Marble diagram:
//
//	Input:  ---L(e1)---R(1)---L(e2)---R(2)---|
//	f(e) = len(e), g(x) = x * 2
//	Output: ---L(3)---R(2)---L(3)---R(4)---|
//
// Both Left and Right values are transformed according to their respective functions.
func MonadBiMap[A, B any](fa SeqResult[A], f Endomorphism[error], g func(A) B) SeqResult[B] {
	return itereither.MonadBiMap(fa, f, g)
}

// BiMap returns a function that maps a pair of functions over the two type arguments of the bifunctor
func BiMap[A, B any](f Endomorphism[error], g func(A) B) Operator[A, B] {
	return itereither.BiMap(f, g)
}

// Fold converts a SeqResult into a Seq by providing handlers for both the error and success cases
func Fold[A, B any](onLeft iter.Kleisli[error, B], onRight iter.Kleisli[A, B]) func(SeqResult[A]) Seq[B] {
	return itereither.Fold(onLeft, onRight)
}

// GetOrElse extracts the value from a successful SeqResult or computes a default value from the error
func GetOrElse[A any](onLeft iter.Kleisli[error, A]) func(SeqResult[A]) Seq[A] {
	return itereither.GetOrElse(onLeft)
}

// GetOrElseOf extracts the value from a successful SeqResult or computes a default value from the error
func GetOrElseOf[A any](onLeft func(error) A) func(SeqResult[A]) Seq[A] {
	return itereither.GetOrElseOf(onLeft)
}

// MonadChainTo sequences two SeqResult computations, discarding the result of the first
func MonadChainTo[A, B any](fa SeqResult[A], fb SeqResult[B]) SeqResult[B] {
	return itereither.MonadChainTo(fa, fb)
}

// ChainTo returns a function that sequences two SeqResult computations, discarding the result of the first
func ChainTo[A, B any](fb SeqResult[B]) Operator[A, B] {
	return itereither.ChainTo[A](fb)
}

// MonadChainToSeq sequences a SeqResult with a Seq, discarding the result of the first
func MonadChainToSeq[A, B any](fa SeqResult[A], fb Seq[B]) SeqResult[B] {
	return itereither.MonadChainToSeq(fa, fb)
}

// ChainToSeq returns a function that sequences a SeqResult with a Seq, discarding the result of the first
func ChainToSeq[A, B any](fb Seq[B]) Operator[A, B] {
	return itereither.ChainToSeq[error, A](fb)
}

// MonadChainFirst executes a side-effecting SeqResult computation but returns the original value
func MonadChainFirst[A, B any](ma SeqResult[A], f Kleisli[A, B]) SeqResult[A] {
	return itereither.MonadChainFirst(ma, f)
}

// MonadTap is an alias for MonadChainFirst, executing a side effect while preserving the original value
func MonadTap[A, B any](ma SeqResult[A], f Kleisli[A, B]) SeqResult[A] {
	return itereither.MonadTap(ma, f)
}

// ChainFirst returns a function that executes a side-effecting SeqResult computation but returns the original value
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return itereither.ChainFirst(f)
}

// Tap is an alias for ChainFirst, executing a side effect while preserving the original value
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return itereither.Tap(f)
}

// MonadChainFirstEitherK executes a side-effecting Result computation but returns the original SeqResult value
func MonadChainFirstEitherK[A, B any](ma SeqResult[A], f either.Kleisli[error, A, B]) SeqResult[A] {
	return itereither.MonadChainFirstEitherK(ma, f)
}

// ChainFirstEitherK returns a function that executes a side-effecting Result computation but returns the original value
func ChainFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return itereither.ChainFirstEitherK(f)
}

// MonadChainFirstResultK executes a side-effecting Result computation but returns the original SeqResult value
func MonadChainFirstResultK[A, B any](ma SeqResult[A], f result.Kleisli[A, B]) SeqResult[A] {
	return itereither.MonadChainFirstEitherK(ma, f)
}

// ChainFirstResultK returns a function that executes a side-effecting Result computation but returns the original value
func ChainFirstResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A] {
	return itereither.ChainFirstEitherK(f)
}

// MonadMergeMapFirstEitherK executes a side-effecting Result computation but returns the original SeqResult value
func MonadMergeMapFirstEitherK[A, B any](ma SeqResult[A], f either.Kleisli[error, A, B]) SeqResult[A] {
	return itereither.MonadMergeMapFirstEitherK(ma, f)
}

// MergeMapFirstEitherK returns a function that executes a side-effecting Result computation but returns the original value
func MergeMapFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return itereither.MergeMapFirstEitherK(f)
}

// MonadMergeMapFirstResultK executes a side-effecting Result computation but returns the original SeqResult value
func MonadMergeMapFirstResultK[A, B any](ma SeqResult[A], f result.Kleisli[A, B]) SeqResult[A] {
	return itereither.MonadMergeMapFirstEitherK(ma, f)
}

// MergeMapFirstResultK returns a function that executes a side-effecting Result computation but returns the original value
func MergeMapFirstResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A] {
	return itereither.MergeMapFirstEitherK(f)
}

// MonadTapEitherK is an alias for MonadChainFirstEitherK, executing a Result side effect while preserving the original value
func MonadTapEitherK[A, B any](ma SeqResult[A], f either.Kleisli[error, A, B]) SeqResult[A] {
	return itereither.MonadTapEitherK(ma, f)
}

// TapEitherK is an alias for ChainFirstEitherK, executing a Result side effect while preserving the original value
func TapEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return itereither.TapEitherK(f)
}

// MonadTapResultK is an alias for MonadChainFirstEitherK, executing a Result side effect while preserving the original value
func MonadTapResultK[A, B any](ma SeqResult[A], f result.Kleisli[A, B]) SeqResult[A] {
	return itereither.MonadTapEitherK(ma, f)
}

// TapResultK is an alias for ChainFirstEitherK, executing a Result side effect while preserving the original value
func TapResultK[A, B any](f result.Kleisli[A, B]) Operator[A, A] {
	return itereither.TapEitherK(f)
}

// MonadFold eliminates a SeqResult by providing handlers for both error and success cases, returning a Seq
func MonadFold[A, B any](ma SeqResult[A], onLeft iter.Kleisli[error, B], onRight iter.Kleisli[A, B]) Seq[B] {
	return itereither.MonadFold(ma, onLeft, onRight)
}

// WithResource constructs a function that safely manages a resource with automatic cleanup.
// It creates a resource, operates on it, and ensures the resource is released even if an error occurs.
func WithResource[A, R, ANY any](onCreate SeqResult[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A] {
	return itereither.WithResource[A](onCreate, onRelease)
}

// MonadAlt provides an alternative SeqResult computation if the first one fails.
//
// Marble diagram:
//
//	First:  ---L(e1)---L(e2)---|
//	Second: ---R(1)---R(2)---|
//	Output: ---R(1)---R(2)---L(e2)---|
//
// When a Left is encountered, it's replaced with values from the second sequence.
// Right values from the first sequence pass through unchanged.
func MonadAlt[A any](first SeqResult[A], second Lazy[SeqResult[A]]) SeqResult[A] {
	return itereither.MonadAlt(first, second)
}

// Alt returns a function that provides an alternative SeqResult computation if the first one fails
func Alt[A any](second Lazy[SeqResult[A]]) Operator[A, A] {
	return itereither.Alt(second)
}

// MonadFlap applies a value to a function wrapped in a SeqResult
func MonadFlap[B, A any](fab SeqResult[func(A) B], a A) SeqResult[B] {
	return itereither.MonadFlap(fab, a)
}

// Flap returns a function that applies a value to a function wrapped in a SeqResult
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return itereither.Flap[error, B](a)
}

// MonadChainLeft chains a computation on the left (error) side of a SeqResult.
// If the input is a Left value, it applies the function f to transform the error and potentially
// change the error type from EA to EB. If the input is a Right value, it passes through unchanged.
//
// Note: MonadChainLeft is identical to OrElse - both provide the same functionality for error recovery.
//
// This is useful for error recovery or error transformation scenarios where you want to handle
// errors by performing another computation that may also fail.
//
// Parameters:
//   - fa: The input SeqResult that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a SeqResult with error type EB
//
// Returns:
//   - A SeqResult with the potentially transformed error type EB
//
// Example:
//
//	// Recover from a specific error by trying an alternative computation
//	result := MonadChainLeft(
//	    Left[int](errors.New("network error")),
//	    func(err error) SeqResult[int] {
//	        if err.Error() == "network error" {
//	            return Right(42) // recover with default value
//	        }
//	        return Left[int](fmt.Errorf("unrecoverable: %w", err))
//	    },
//	)
func MonadChainLeft[A any](fa SeqResult[A], f Kleisli[error, A]) SeqResult[A] {
	return itereither.MonadChainLeft(fa, f)
}

// ChainLeft is the curried version of MonadChainLeft.
// It returns a function that chains a computation on the left (error) side of a SeqResult.
//
// Note: ChainLeft is identical to OrElse - both provide the same functionality for error recovery.
//
// This is particularly useful in functional composition pipelines where you want to handle
// errors by performing another computation that may also fail.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a SeqResult with error type EB
//
// Returns:
//   - A function that transforms a SeqResult with error type EA to one with error type EB
//
// Example:
//
//	// Create a reusable error handler
//	recoverFromNetworkError := ChainLeft(func(err error) SeqResult[int] {
//	    if strings.Contains(err.Error(), "network") {
//	        return Right(0) // return default value
//	    }
//	    return Left[int](err) // propagate other errors
//	})
//
//	result := F.Pipe1(
//	    Left[int](errors.New("network timeout")),
//	    recoverFromNetworkError,
//	)
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A] {
	return itereither.ChainLeft(f)
}

// MonadChainFirstLeft chains a computation on the left (error) side but always returns the original error.
// If the input is a Left value, it applies the function f to the error and executes the resulting computation,
// but always returns the original Left error regardless of what f returns (Left or Right).
// If the input is a Right value, it passes through unchanged without calling f.
//
// This is useful for side effects on errors (like logging or metrics) where you want to perform an action
// when an error occurs but always propagate the original error, ensuring the error path is preserved.
//
// Parameters:
//   - ma: The input SeqResult that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a SeqResult (typically for side effects)
//
// Returns:
//   - A SeqResult with the original error preserved if input was Left, or the original Right value
//
// Example:
//
//	// Log errors but always preserve the original error
//	result := MonadChainFirstLeft(
//	    Left[int](errors.New("database error")),
//	    func(err error) SeqResult[int] {
//	        log.Printf("Error occurred: %s", err)
//	        return Right(0)
//	    },
//	)
//	// result will always be Left(error: "database error"), even though f returns Right
func MonadChainFirstLeft[A, B any](ma SeqResult[A], f Kleisli[error, B]) SeqResult[A] {
	return itereither.MonadChainFirstLeft(ma, f)
}

//go:inline
func MonadTapLeft[A, B any](ma SeqResult[A], f Kleisli[error, B]) SeqResult[A] {
	return itereither.MonadTapLeft(ma, f)
}

// ChainFirstLeft is the curried version of MonadChainFirstLeft.
// It returns a function that chains a computation on the left (error) side while always preserving the original error.
//
// This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// ensuring the error path is preserved.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a SeqResult (typically for side effects)
//
// Returns:
//   - An Operator that performs the side effect but always returns the original error if input was Left
//
// Example:
//
//	// Create a reusable error logger
//	logError := ChainFirstLeft(func(err error) SeqResult[int] {
//	    log.Printf("Error: %s", err)
//	    return Right(0)
//	})
//
//	result := F.Pipe1(
//	    Left[int](errors.New("validation failed")),
//	    logError, // logs the error
//	)
//	// result is always Left(error: "validation failed"), even though f returns Right
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return itereither.ChainFirstLeft[A](f)
}

//go:inline
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return itereither.TapLeft[A](f)
}

// OrElse recovers from a Left (error) by providing an alternative computation.
// If the SeqResult is Right, it returns the value unchanged.
// If the SeqResult is Left, it applies the provided function to the error value,
// which returns a new SeqResult that replaces the original.
//
// Note: OrElse is identical to ChainLeft - both provide the same functionality for error recovery.
//
// This is useful for error recovery, fallback logic, or chaining alternative IO computations.
// The error type can be widened from E1 to E2, allowing transformation of error types.
//
// Example:
//
//	// Recover from specific errors with fallback operations
//	recover := OrElse(func(err error) SeqResult[int] {
//	    if err.Error() == "not found" {
//	        return Right(0) // default value
//	    }
//	    return Left[int](err) // propagate other errors
//	})
//	result := recover(Left[int](errors.New("not found"))) // Right(0)
//	result := recover(Right(42)) // Right(42) - unchanged
//
//go:inline
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A] {
	return itereither.OrElse(onLeft)
}

// MonadReduce reduces a SeqResult to a single Result value by applying a function to each
// Right element and an accumulator, starting with an initial value. If any Left is encountered,
// reduction stops immediately and returns that Left value.
//
// Type Parameters:
//   - E: The error type
//   - A: The element type in the sequence
//   - B: The accumulator and result type
//
// Parameters:
//   - fa: The SeqResult to reduce
//   - f: The reducer function that combines the accumulator with each element
//   - initial: The initial accumulator value
//
// Returns:
//   - Result[B]: Left with the first error encountered, or Right with the final accumulated value
//
// Marble Diagram:
//
//	Input:  --R(1)--R(2)--R(3)--R(4)--R(5)--|
//	Reduce((acc, x) => acc + x, 0)
//	Output: Right(15)
//
//	Input:  --R(1)--R(2)--L(e)--R(4)--R(5)--|
//	Reduce((acc, x) => acc + x, 0)
//	Output: Left(e)
//
// Example:
//
//	seq := func(yield func(Result[int]) bool) {
//	    yield(Right(1))
//	    yield(Right(2))
//	    yield(Right(3))
//	}
//	sum := MonadReduce(seq, func(acc, x int) int { return acc + x }, 0)
//	// returns: Right(6)
//
//	seqWithError := func(yield func(Result[int]) bool) {
//	    yield(Right(1))
//	    yield(Left[int](errors.New("error")))
//	    yield(Right(3))
//	}
//	result := MonadReduce(seqWithError, func(acc, x int) int { return acc + x }, 0)
//	// returns: Left(error: "error")
//
//go:inline
func MonadReduce[A, B any](fa SeqResult[A], f func(B, A) B, initial B) Result[B] {
	return itereither.MonadReduce(fa, f, initial)
}

// Reduce returns a function that reduces a SeqResult to a single Result value.
// This is the curried version of MonadReduce.
//
// Type Parameters:
//   - E: The error type
//   - A: The element type in the sequence
//   - B: The accumulator and result type
//
// Parameters:
//   - f: The reducer function that combines the accumulator with each element
//   - initial: The initial accumulator value
//
// Returns:
//   - A function that takes a SeqResult and returns Result[B]
//
// Example:
//
//	sum := Reduce(func(acc, x int) int { return acc + x }, 0)
//	seq := func(yield func(Result[int]) bool) {
//	    yield(Right(1))
//	    yield(Right(2))
//	    yield(Right(3))
//	}
//	result := sum(seq)
//	// returns: Right(6)
func Reduce[A, B any](f func(B, A) B, initial B) func(SeqResult[A]) Result[B] {
	return itereither.Reduce[error](f, initial)
}
