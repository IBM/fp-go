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

package itereither

import (
	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/file"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromiter"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
)

type (
	// Seq is a single-value iterator sequence from Go 1.23+.
	// It represents a lazy sequence of values that can be iterated using range.
	//
	// Type Parameters:
	//   - A: The type of elements in the sequence
	Seq[A any] = iter.Seq[A]

	// Either represents a value that can be one of two types: Left (error) or Right (success).
	// It is used for computations that may fail, providing a type-safe alternative to
	// error handling with multiple return values.
	//
	// Type Parameters:
	//   - E: The error type (Left)
	//   - A: The success type (Right)
	Either[E, A any] = either.Either[E, A]

	// SeqEither represents a lazy sequence of Either values.
	// Each element in the sequence can be either a Left (error) or Right (success).
	// This combines the lazy evaluation of sequences with the error handling of Either.
	//
	// Type Parameters:
	//   - E: The error type
	//   - A: The success type
	//
	// Example:
	//
	//	seq := iter.From(
	//	    either.Right[string](1),
	//	    either.Right[string](2),
	//	    either.Left[int]("error"),
	//	)
	//	// seq is a SeqEither[string, int]
	SeqEither[E, A any] = Seq[Either[E, A]]

	// Kleisli represents a function that takes a value and returns a SeqEither.
	// This is the monadic bind operation for SeqEither, enabling composition of
	// computations that may produce sequences of results or errors.
	//
	// Type Parameters:
	//   - E: The error type
	//   - A: The input type
	//   - B: The element type of the output sequence
	//
	// Example:
	//
	//	validate := func(x int) SeqEither[string, int] {
	//	    if x > 0 {
	//	        return Right[string](x)
	//	    }
	//	    return Left[int]("must be positive")
	//	}
	//	// validate is a Kleisli[string, int, int]
	Kleisli[E, A, B any] = R.Reader[A, SeqEither[E, B]]

	// Operator represents a transformation from one SeqEither to another.
	// It takes a SeqEither[E, A] and returns a SeqEither[E, B], allowing for
	// composition of sequence transformations in a functional pipeline.
	//
	// Type Parameters:
	//   - E: The error type (preserved across the transformation)
	//   - A: The element type of the input sequence
	//   - B: The element type of the output sequence
	//
	// Example:
	//
	//	double := Map[string](N.Mul(2))
	//	// double is an Operator[string, int, int]
	Operator[E, A, B any] = Kleisli[E, SeqEither[E, A], B]
)

// FromIO converts an IO computation into a single-element SeqEither containing a Right value.
// The IO computation is executed when the sequence is consumed, and its result is wrapped
// in a Right, creating a successful SeqEither.
//
// This function bridges IO computations with SeqEither, allowing side effects that always
// succeed to be integrated into error-handling pipelines.
//
// Type Parameters:
//   - E: The error type (not used, but required for type consistency)
//   - A: The type of value produced by the IO computation
//
// Parameters:
//   - mr: An IO computation that produces a value of type A
//
// Returns:
//   - SeqEither[E, A]: A sequence containing a single Right value with the IO result
//
// Example:
//
//	getCurrentTime := func() time.Time { return time.Now() }
//	seq := FromIO[error](getCurrentTime)
//	// seq yields: Right(current time) when consumed
//
// See Also:
//   - FromLazy: Converts a Lazy computation to SeqEither
//   - FromIOEither: Converts an IOEither to SeqEither
//   - Right: Creates a SeqEither from a pure value
func FromIO[E, A any](mr IO[A]) SeqEither[E, A] {
	return iter.FromIO(eithert.RightF(io.MonadMap[A, Either[E, A]], mr))
}

// FromLazy converts a Lazy computation into a single-element SeqEither containing a Right value.
// The Lazy computation is executed when the sequence is consumed, and its result is wrapped
// in a Right, creating a successful SeqEither.
//
// This is a convenience function that provides semantic clarity when working with lazy
// computations. It delegates to FromIO since Lazy is an alias for IO.
//
// Type Parameters:
//   - E: The error type (not used, but required for type consistency)
//   - A: The type of value produced by the Lazy computation
//
// Parameters:
//   - mr: A Lazy computation that produces a value of type A
//
// Returns:
//   - SeqEither[E, A]: A sequence containing a single Right value with the Lazy result
//
// Example:
//
//	expensiveCalc := func() int {
//	    // Expensive computation
//	    return 42
//	}
//	seq := FromLazy[error](expensiveCalc)
//	// Computation only runs when sequence is consumed
//	// seq yields: Right(42)
//
// See Also:
//   - FromIO: Converts an IO computation to SeqEither
//   - FromIOEither: Converts an IOEither to SeqEither
//   - Right: Creates a SeqEither from a pure value
func FromLazy[E, A any](mr Lazy[A]) SeqEither[E, A] {
	return iter.FromIO(eithert.RightF(io.MonadMap[A, Either[E, A]], mr))
}

// FromIOEither converts an IOEither computation into a single-element SeqEither.
// The IOEither computation is executed when the sequence is consumed, and its result
// (either Left or Right) becomes the single element of the sequence.
//
// This function bridges IOEither computations with SeqEither, allowing deferred I/O operations
// that may fail to be integrated into sequence-based error-handling pipelines.
//
// Type Parameters:
//   - E: The error type
//   - A: The success type
//
// Parameters:
//   - mr: An IOEither computation that produces Either[E, A]
//
// Returns:
//   - SeqEither[E, A]: A sequence containing the single Either value from the deferred IOEither
//
// Example:
//
//	readConfig := func() either.Either[error, Config] {
//	    data, err := os.ReadFile("config.json")
//	    if err != nil {
//	        return either.Left[Config](err)
//	    }
//	    return either.Right[error](parseConfig(data))
//	}
//	seq := FromIOEither(readConfig)
//	// The IOEither is only executed when the sequence is iterated.
//
// See Also:
//   - FromIO: Converts an IO computation to SeqEither (always Right)
//   - FromLazy: Converts a Lazy computation to SeqEither (always Right)
//   - FromEither: Converts a pure Either to SeqEither
func FromIOEither[E, A any](mr IOEither[E, A]) SeqEither[E, A] {
	return iter.FromIO(mr)
}

// Left constructs an [SeqEither] that represents a failure with an error value of type E
func Left[A, E any](l E) SeqEither[E, A] {
	return eithert.Left(iter.Of[Either[E, A]], l)
}

// Right constructs an [SeqEither] that represents a successful computation with a value of type A
func Right[E, A any](r A) SeqEither[E, A] {
	return eithert.Right(iter.Of[Either[E, A]], r)
}

// Of constructs an [SeqEither] that represents a successful computation with a value of type A.
// This is an alias for [Right] and is the canonical way to lift a pure value into the SeqEither context.
func Of[E, A any](r A) SeqEither[E, A] {
	return Right[E](r)
}

// MonadOf is an alias for [Of], provided for consistency with monad naming conventions
func MonadOf[E, A any](r A) SeqEither[E, A] {
	return Of[E](r)
}

// LeftSeq constructs an [SeqEither] from an [Seq] that produces an error value
func LeftSeq[A, E any](ml Seq[E]) SeqEither[E, A] {
	return eithert.LeftF(iter.MonadMap[E, Either[E, A]], ml)
}

// RightSeq constructs an [SeqEither] from an [Seq] that produces a success value
func RightSeq[E, A any](mr Seq[A]) SeqEither[E, A] {
	return eithert.RightF(iter.MonadMap[A, Either[E, A]], mr)
}

// FromEither lifts an [Either] value into the [SeqEither] context
func FromEither[E, A any](e Either[E, A]) SeqEither[E, A] {
	return iter.Of(e)
}

// FromOption converts an Option[A] to a single-element SeqEither.
// Some(a) becomes Right(a). None becomes Left(onNone()), using the provided
// function to compute the error value.
//
// Type Parameters:
//   - A: The success type
//   - E: The error type
//
// Parameters:
//   - onNone: A function that produces the error value when the Option is None
//
// Returns:
//   - A Kleisli function that converts Option[A] to SeqEither[E, A]
//
// Example:
//
//	fromOpt := FromOption[int](func() string { return "missing value" })
//	seq := fromOpt(option.Some(42))   // yields: Right(42)
//	seq2 := fromOpt(option.None[int]()) // yields: Left("missing value")
func FromOption[A, E any](onNone func() E) Kleisli[E, O.Option[A], A] {
	return fromeither.FromOption(
		FromEither[E, A],
		onNone,
	)
}

// ChainOptionK chains an Option-returning function into the SeqEither context.
// For each Right(a) in the sequence it applies f(a); if f returns Some(b) the
// result is Right(b), if f returns None the result is Left(onNone()).
// Left values in the input sequence pass through unchanged.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//   - E: The error type
//
// Parameters:
//   - onNone: Produces the error value when f returns None
//
// Returns:
//   - A function that takes an Option-returning Kleisli and produces an Operator
//
// Example:
//
//	safeDiv := ChainOptionK[int, int](func() string { return "division by zero" })(
//	    func(n int) option.Option[int] {
//	        if n == 0 { return option.None[int]() }
//	        return option.Some(100 / n)
//	    },
//	)
//	result := safeDiv(Right[string](5))  // yields: Right(20)
//	result2 := safeDiv(Right[string](0)) // yields: Left("division by zero")
func ChainOptionK[A, B, E any](onNone func() E) func(O.Kleisli[A, B]) Operator[E, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[E, A, B],
		FromEither[E, B],
		onNone,
	)
}

// MonadChainSeqK is the uncurried form of ChainSeqK.
// For each Right(a) in ma it calls f(a) (which returns a plain Seq[B]), lifts
// every element of that sequence into Right, and concatenates them in order.
// Left values in the input pass through unchanged.
//
// This is useful when you have an infallible sequence-producing function that
// you want to use inside a SeqEither pipeline.
//
// Example:
//
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := MonadChainSeqK(seq, func(n int) iter.Seq[int] {
//	    return iter.From(n, n*10)
//	})
//	// yields: Right(1), Right(10), Right(2), Right(20)
func MonadChainSeqK[E, A, B any](ma SeqEither[E, A], f iter.Kleisli[A, B]) SeqEither[E, B] {
	return fromiter.MonadChainIOK(
		MonadChain[E, A, B],
		FromSeq[E, B],
		ma,
		f,
	)
}

// ChainSeqK is the curried version of MonadChainSeqK.
// It returns an Operator that chains an infallible Seq-returning function into
// the SeqEither context, lifting each produced value into Right.
//
// Example:
//
//	expand := ChainSeqK[string](func(n int) iter.Seq[int] {
//	    return iter.From(n, n*10)
//	})
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := expand(seq)
//	// yields: Right(1), Right(10), Right(2), Right(20)
func ChainSeqK[E, A, B any](f iter.Kleisli[A, B]) Operator[E, A, B] {
	return fromiter.ChainIOK(
		Chain[E, A, B],
		FromSeq[E, B],
		f,
	)
}

// MonadMergeMapSeqK is the uncurried form of MergeMapSeqK.
// It is identical to MonadChainSeqK except it uses MergeMap semantics:
// results from different inner sequences may be interleaved rather than
// strictly concatenated in input order.
//
// Use MonadChainSeqK (ordered) when result order matters; use
// MonadMergeMapSeqK when throughput is more important than ordering.
//
// Example:
//
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := MonadMergeMapSeqK(seq, func(n int) iter.Seq[int] {
//	    return iter.From(n * 10)
//	})
//	// yields: Right(10), Right(20) in some order
func MonadMergeMapSeqK[E, A, B any](ma SeqEither[E, A], f iter.Kleisli[A, B]) SeqEither[E, B] {
	return fromiter.MonadChainIOK(
		MonadMergeMap[E, A, B],
		FromSeq[E, B],
		ma,
		f,
	)
}

// MergeMapSeqK is the curried version of MonadMergeMapSeqK.
// It returns an Operator that uses MergeMap (interleaving) semantics when
// expanding Right values through an infallible Seq-producing function.
//
// Example:
//
//	expand := MergeMapSeqK[string](func(n int) iter.Seq[int] {
//	    return iter.From(n * 10)
//	})
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := expand(seq)
//	// yields: Right(10), Right(20) in some order
func MergeMapSeqK[E, A, B any](f iter.Kleisli[A, B]) Operator[E, A, B] {
	return fromiter.ChainIOK(
		MergeMap[E, A, B],
		FromSeq[E, B],
		f,
	)
}

// FromSeq creates an [SeqEither] from an [Seq] instance, invoking [Seq] for each invocation of [SeqEither]
func FromSeq[E, A any](mr Seq[A]) SeqEither[E, A] {
	return RightSeq[E](mr)
}

// MonadMap applies a function to the value inside a successful SeqEither, leaving errors unchanged.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---R(3)---|
//	f(x) = x * 2
//	Output: ---R(2)---R(4)---L(e)---R(6)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func MonadMap[E, A, B any](fa SeqEither[E, A], f func(A) B) SeqEither[E, B] {
	return eithert.MonadMap(iter.MonadMap[Either[E, A], Either[E, B]], fa, f)
}

// Map returns a function that applies a transformation to the value inside a successful SeqEither.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---R(3)---|
//	f(x) = x * 2
//	Output: ---R(2)---R(4)---L(e)---R(6)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return eithert.Map(iter.Map[Either[E, A], Either[E, B]], f)
}

// MonadMapTo replaces the value inside a successful [SeqEither] with a constant value
func MonadMapTo[E, A, B any](fa SeqEither[E, A], b B) SeqEither[E, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

// MapTo returns a function that replaces the value inside a successful [SeqEither] with a constant value
func MapTo[E, A, B any](b B) Operator[E, A, B] {
	return Map[E](function.Constant1[A](b))
}

// MonadChain sequences two SeqEither computations, where the second depends on the result of the first.
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
func MonadChain[E, A, B any](fa SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, B] {
	return eithert.MonadChain(iter.MonadChain[Either[E, A], Either[E, B]], iter.MonadOf[Either[E, B]], fa, f)
}

// MonadChainSeq sequences two SeqEither computations using purely sequential
// iteration — no goroutines, no channels. Each Right value is bound to f and
// the inner sequences are drained one at a time before the next input element
// is processed. Left values pass through unchanged.
//
// Marble diagram:
//
//	Input:  ---R(1)-------R(2)---L(e)---|
//	f(1) -> ---R(10)---R(11)---|
//	f(2) -> ---R(20)---R(21)---|
//	Output: ---R(10)---R(11)---R(20)---R(21)---L(e)---|
//
// See Also:
//   - MonadChain: Concurrent inner producers, sequential drain
//   - MonadChainPar: Concurrent inner producers and concurrent drain
func MonadChainSeq[E, A, B any](fa SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, B] {
	return eithert.MonadChain(iter.MonadChainSeq[Either[E, A], Either[E, B]], iter.MonadOf[Either[E, B]], fa, f)
}

// MonadChainPar sequences two SeqEither computations using concurrent inner
// producers drained in input order. Right values are bound to f, and all inner
// sequences run concurrently. Left values pass through unchanged.
//
// Marble diagram:
//
//	Input:  ---R(1)-------R(2)---|
//	f(1) -> ---R(10)------R(11)---|
//	f(2) -> ------R(20)------R(21)---|
//	Output: ---R(10)---R(20)---R(11)---R(21)---|  (order may vary)
//
// See Also:
//   - MonadChain: Concurrent producers, same output order
//   - MonadChainSeq: Always sequential, no goroutines
func MonadChainPar[E, A, B any](fa SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, B] {
	return eithert.MonadChain(iter.MonadChainPar[Either[E, A], Either[E, B]], iter.MonadOf[Either[E, B]], fa, f)
}

// MonadMergeMap sequences two SeqEither computations, where the second depends on the result of the first.
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
func MonadMergeMap[E, A, B any](fa SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, B] {
	return eithert.MonadChain(iter.MonadMergeMap[Either[E, A], Either[E, B]], iter.MonadOf[Either[E, B]], fa, f)
}

// Chain returns a function that sequences two [SeqEither] computations
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return eithert.Chain(iter.Chain[Either[E, A], Either[E, B]], iter.Of[Either[E, B]], f)
}

// ChainSeq is the curried version of MonadChainSeq. It returns an Operator that
// sequences two SeqEither computations using purely sequential iteration — no
// goroutines, no channels. Left values pass through unchanged.
//
// See Also:
//   - Chain: Concurrent inner producers, sequential drain
//   - ChainPar: Always concurrent
func ChainSeq[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return eithert.Chain(iter.ChainSeq[Either[E, A], Either[E, B]], iter.Of[Either[E, B]], f)
}

// ChainPar is the curried version of MonadChainPar. It returns an Operator that
// sequences two SeqEither computations with concurrent inner producers drained
// in input order. Left values pass through unchanged.
//
// See Also:
//   - Chain: Dispatches to sequential when buffer size is 1
//   - ChainSeq: Always sequential, no goroutines
func ChainPar[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return eithert.Chain(iter.ChainPar[Either[E, A], Either[E, B]], iter.Of[Either[E, B]], f)
}

// MergeMap returns a function that sequences two [SeqEither] computations
func MergeMap[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return eithert.Chain(iter.MergeMap[Either[E, A], Either[E, B]], iter.Of[Either[E, B]], f)
}

// MonadChainEitherK chains an [either.Kleisli] function into a SeqEither
// pipeline. For each Right(a) in ma it calls f(a), lifting the resulting Either
// into the sequence. Left values pass through unchanged.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---|
//	f(1) = Right(10), f(2) = Left("bad")
//	Output: ---R(10)---L("bad")---L(e)---|
func MonadChainEitherK[E, A, B any](ma SeqEither[E, A], f either.Kleisli[E, A, B]) SeqEither[E, B] {
	return fromeither.MonadChainEitherK(
		MonadChainSeq[E, A, B],
		FromEither[E, B],
		ma,
		f,
	)
}

// ChainEitherK is the curried version of MonadChainEitherK. It returns an
// Operator that chains an [either.Kleisli] function into a SeqEither pipeline.
// For each Right(a) in the input sequence it applies f(a) and lifts the
// resulting Either into the sequence. Left values pass through unchanged.
func ChainEitherK[E, A, B any](f either.Kleisli[E, A, B]) Operator[E, A, B] {
	return fromeither.ChainEitherK(
		ChainSeq[E, A, B],
		FromEither[E, B],
		f,
	)
}

// MonadAp applies a function wrapped in an [SeqEither] to a value wrapped in an [SeqEither]
func MonadAp[B, E, A any](mab SeqEither[E, func(A) B], ma SeqEither[E, A]) SeqEither[E, B] {
	return eithert.MonadAp(
		iter.MonadAp[Either[E, B], Either[E, A]],
		iter.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// MonadApSeq applies a function wrapped in a SeqEither to a value wrapped in a
// SeqEither using purely sequential iteration — no goroutines, no channels. Each
// function is applied to every value before the next function is processed.
// Left values in either argument propagate as errors.
//
// See Also:
//   - MonadAp: Concurrent inner producers, same output order
//   - MonadApPar: Always concurrent
func MonadApSeq[B, E, A any](mab SeqEither[E, func(A) B], ma SeqEither[E, A]) SeqEither[E, B] {
	return eithert.MonadAp(
		iter.MonadApSeq[Either[E, B], Either[E, A]],
		iter.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// MonadApPar applies a function wrapped in a SeqEither to a value wrapped in a
// SeqEither using concurrent inner producers drained in input order.
// Left values in either argument propagate as errors.
//
// See Also:
//   - MonadAp: Dispatches to sequential when buffer size is 1
//   - MonadApSeq: Always sequential, no goroutines
func MonadApPar[B, E, A any](mab SeqEither[E, func(A) B], ma SeqEither[E, A]) SeqEither[E, B] {
	return eithert.MonadAp(
		iter.MonadApPar[Either[E, B], Either[E, A]],
		iter.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// Ap applies a function wrapped in an [SeqEither] to a value wrapped in an [SeqEither].
// This is an alias of [ApPar] which applies the function and value in parallel.
func Ap[B, E, A any](ma SeqEither[E, A]) Operator[E, func(A) B, B] {
	return eithert.Ap(
		iter.Ap[Either[E, B], Either[E, A]],
		iter.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

// ApSeq is the curried version of MonadApSeq. It returns an Operator that
// applies a function wrapped in a SeqEither to the value sequence ma using
// purely sequential iteration — no goroutines, no channels.
//
// See Also:
//   - Ap: Concurrent producers, same output order
//   - ApPar: Always concurrent
func ApSeq[B, E, A any](ma SeqEither[E, A]) Operator[E, func(A) B, B] {
	return eithert.Ap(
		iter.ApSeq[Either[E, B], Either[E, A]],
		iter.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

// ApPar is the curried version of MonadApPar. It returns an Operator that
// applies a function wrapped in a SeqEither to the value sequence ma using
// concurrent inner producers drained in input order.
//
// See Also:
//   - Ap: Dispatches to sequential when buffer size is 1
//   - ApSeq: Always sequential, no goroutines
func ApPar[B, E, A any](ma SeqEither[E, A]) Operator[E, func(A) B, B] {
	return eithert.Ap(
		iter.ApPar[Either[E, B], Either[E, A]],
		iter.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

// Flatten removes one level of nesting from a nested [SeqEither]
func Flatten[E, A any](mma SeqEither[E, SeqEither[E, A]]) SeqEither[E, A] {
	return MonadChain(mma, function.Identity[SeqEither[E, A]])
}

// FlattenSeq removes one level of nesting from a nested SeqEither using purely
// sequential iteration — no goroutines, no channels. Equivalent to
// MonadChainSeq with the identity function.
//
// See Also:
//   - Flatten: Concurrent inner producers, sequential drain
//   - FlattenPar: Concurrent inner producers and concurrent drain
func FlattenSeq[E, A any](mma SeqEither[E, SeqEither[E, A]]) SeqEither[E, A] {
	return MonadChainSeq(mma, function.Identity[SeqEither[E, A]])
}

// FlattenPar removes one level of nesting from a nested SeqEither using
// concurrent inner producers drained in input order. Equivalent to
// MonadChainPar with the identity function.
//
// See Also:
//   - Flatten: Concurrent inner producers, sequential drain
//   - FlattenSeq: Always sequential, no goroutines
func FlattenPar[E, A any](mma SeqEither[E, SeqEither[E, A]]) SeqEither[E, A] {
	return MonadChainPar(mma, function.Identity[SeqEither[E, A]])
}

// MonadMapLeft applies a function to the error value of a failed SeqEither, leaving successful values unchanged.
//
// Marble diagram:
//
//	Input:  ---L(e1)---R(1)---L(e2)---R(2)---|
//	f(e) = "error: " + e
//	Output: ---L("error: e1")---R(1)---L("error: e2")---R(2)---|
//
// Where R(x) represents Right(x) and L(e) represents Left(e).
func MonadMapLeft[A, E1, E2 any](fa SeqEither[E1, A], f func(E1) E2) SeqEither[E2, A] {
	return eithert.MonadMapLeft(
		iter.MonadMap[Either[E1, A], Either[E2, A]],
		fa,
		f,
	)
}

// MapLeft returns a function that applies a transformation to the error value of a failed [SeqEither]
func MapLeft[A, E1, E2 any](f func(E1) E2) func(SeqEither[E1, A]) SeqEither[E2, A] {
	return eithert.MapLeft(
		iter.Map[Either[E1, A], Either[E2, A]],
		f,
	)
}

// MonadBiMap applies one function to the error value and another to the success value of a SeqEither.
//
// Marble diagram:
//
//	Input:  ---L(e1)---R(1)---L(e2)---R(2)---|
//	f(e) = len(e), g(x) = x * 2
//	Output: ---L(3)---R(2)---L(3)---R(4)---|
//
// Both Left and Right values are transformed according to their respective functions.
func MonadBiMap[E1, E2, A, B any](fa SeqEither[E1, A], f func(E1) E2, g func(A) B) SeqEither[E2, B] {
	return eithert.MonadBiMap(iter.MonadMap[Either[E1, A], Either[E2, B]], fa, f, g)
}

// BiMap returns a function that maps a pair of functions over the two type arguments of the bifunctor
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(SeqEither[E1, A]) SeqEither[E2, B] {
	return eithert.BiMap(iter.Map[Either[E1, A], Either[E2, B]], f, g)
}

// Fold converts an [SeqEither] into an [Seq] by providing handlers for both the error and success cases
func Fold[E, A, B any](onLeft iter.Kleisli[E, B], onRight iter.Kleisli[A, B]) func(SeqEither[E, A]) Seq[B] {
	return eithert.MatchE(iter.MonadChain[Either[E, A], B], onLeft, onRight)
}

// GetOrElse extracts the value from a successful [SeqEither] or computes a default value from the error
func GetOrElse[E, A any](onLeft iter.Kleisli[E, A]) func(SeqEither[E, A]) Seq[A] {
	return eithert.GetOrElse(iter.MonadChain[Either[E, A], A], iter.MonadOf[A], onLeft)
}

// GetOrElseOf extracts the value from a successful [SeqEither] or computes a default value from the error
func GetOrElseOf[E, A any](onLeft func(E) A) func(SeqEither[E, A]) Seq[A] {
	return eithert.GetOrElseOf(iter.MonadChain[Either[E, A], A], iter.MonadOf[A], onLeft)
}

// MonadChainTo sequences two [SeqEither] computations, discarding the result of the first
func MonadChainTo[A, E, B any](fa SeqEither[E, A], fb SeqEither[E, B]) SeqEither[E, B] {
	return MonadChain(fa, function.Constant1[A](fb))
}

// ChainTo returns a function that sequences two [SeqEither] computations, discarding the result of the first
func ChainTo[A, E, B any](fb SeqEither[E, B]) Operator[E, A, B] {
	return Chain(function.Constant1[A](fb))
}

// MonadChainToSeq sequences an [SeqEither] with an [Seq], discarding the result of the first
func MonadChainToSeq[E, A, B any](fa SeqEither[E, A], fb Seq[B]) SeqEither[E, B] {
	return MonadChainTo(fa, FromSeq[E](fb))
}

// ChainToSeq returns a function that sequences an [SeqEither] with an [Seq], discarding the result of the first
func ChainToSeq[E, A, B any](fb Seq[B]) Operator[E, A, B] {
	return ChainTo[A](FromSeq[E](fb))
}

// MonadChainFirst executes a side-effecting [SeqEither] computation but returns the original value
func MonadChainFirst[E, A, B any](ma SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, A] {
	return chain.MonadChainFirst(
		MonadChain[E, A, A],
		MonadMap[E, B, A],
		ma,
		f,
	)
}

// MonadTap is an alias for [MonadChainFirst], executing a side effect while preserving the original value
func MonadTap[E, A, B any](ma SeqEither[E, A], f Kleisli[E, A, B]) SeqEither[E, A] {
	return MonadChainFirst(ma, f)
}

// ChainFirst returns a function that executes a side-effecting [SeqEither] computation but returns the original value
func ChainFirst[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A] {
	return chain.ChainFirst(
		Chain[E, A, A],
		Map[E, B, A],
		f,
	)
}

// Tap is an alias for [ChainFirst], executing a side effect while preserving the original value
func Tap[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A] {
	return ChainFirst(f)
}

// MonadChainFirstEitherK executes a side-effecting [Either] computation but returns the original [SeqEither] value
func MonadChainFirstEitherK[A, E, B any](ma SeqEither[E, A], f either.Kleisli[E, A, B]) SeqEither[E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChainSeq[E, A, A],
		MonadMap[E, B, A],
		FromEither[E, B],
		ma,
		f,
	)
}

// ChainFirstEitherK returns a function that executes a side-effecting [Either] computation but returns the original value
func ChainFirstEitherK[A, E, B any](f either.Kleisli[E, A, B]) Operator[E, A, A] {
	return fromeither.ChainFirstEitherK(
		ChainSeq[E, A, A],
		Map[E, B, A],
		FromEither[E, B],
		f,
	)
}

// MonadMergeMapFirstEitherK is the MergeMap variant of MonadChainFirstEitherK.
// For each Right(a) in ma it calls f(a) as a side effect: if f returns Right,
// the original Right(a) is passed through; if f returns Left, that error
// propagates. Left values in the input pass through unchanged without calling f.
//
// Use this instead of MonadChainFirstEitherK when inner work can be processed
// concurrently and result order does not need to match input order.
func MonadMergeMapFirstEitherK[A, E, B any](ma SeqEither[E, A], f either.Kleisli[E, A, B]) SeqEither[E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadMergeMap[E, A, A],
		MonadMap[E, B, A],
		FromEither[E, B],
		ma,
		f,
	)
}

// MergeMapFirstEitherK is the curried version of MonadMergeMapFirstEitherK. It
// returns an Operator that executes a side-effecting [Either] computation for
// each Right value: if f returns Right, the original value is passed through;
// if f returns Left, that error propagates. Left values in the input pass
// through unchanged without calling f.
func MergeMapFirstEitherK[A, E, B any](f either.Kleisli[E, A, B]) Operator[E, A, A] {
	return fromeither.ChainFirstEitherK(
		MergeMap[E, A, A],
		Map[E, B, A],
		FromEither[E, B],
		f,
	)
}

// MonadTapEitherK is an alias for [MonadChainFirstEitherK], executing an [Either] side effect while preserving the original value
func MonadTapEitherK[A, E, B any](ma SeqEither[E, A], f either.Kleisli[E, A, B]) SeqEither[E, A] {
	return MonadChainFirstEitherK(ma, f)
}

// TapEitherK is an alias for [ChainFirstEitherK], executing an [Either] side effect while preserving the original value
func TapEitherK[A, E, B any](f either.Kleisli[E, A, B]) Operator[E, A, A] {
	return ChainFirstEitherK(f)
}

// MonadFold eliminates an [SeqEither] by providing handlers for both error and success cases, returning an [Seq]
func MonadFold[E, A, B any](ma SeqEither[E, A], onLeft iter.Kleisli[E, B], onRight iter.Kleisli[A, B]) Seq[B] {
	return eithert.FoldE(iter.MonadChain[Either[E, A], B], ma, onLeft, onRight)
}

// WithResource constructs a function that safely manages a resource with automatic cleanup.
// It creates a resource, operates on it, and ensures the resource is released even if an error occurs.
func WithResource[A, E, R, ANY any](onCreate SeqEither[E, R], onRelease Kleisli[E, R, ANY]) Kleisli[E, Kleisli[E, R, A], A] {
	return file.WithResource(
		MonadChain[E, R, A],
		MonadFold[E, A, Either[E, A]],
		MonadFold[E, ANY, Either[E, A]],
		MonadMap[E, ANY, A],
		Left[A, E],
	)(function.Constant(onCreate), onRelease)
}

// Swap exchanges the error and success type parameters of an [SeqEither]
func Swap[E, A any](val SeqEither[E, A]) SeqEither[A, E] {
	return MonadFold(val, Right[A, E], Left[E, A])
}

// MonadAlt provides an alternative SeqEither computation if the first one fails.
//
// Marble diagram:
//
//	First:  ---L(e1)---L(e2)---|
//	Second: ---R(1)---R(2)---|
//	Output: ---R(1)---R(2)---L(e2)---|
//
// When a Left is encountered, it's replaced with values from the second sequence.
// Right values from the first sequence pass through unchanged.
func MonadAlt[E, A any](first SeqEither[E, A], second lazy.Lazy[SeqEither[E, A]]) SeqEither[E, A] {
	return eithert.MonadAlt(
		iter.Of[Either[E, A]],
		iter.MonadChain[Either[E, A], Either[E, A]],

		first,
		second,
	)
}

// Alt returns a function that provides an alternative [SeqEither] computation if the first one fails
func Alt[E, A any](second lazy.Lazy[SeqEither[E, A]]) Operator[E, A, A] {
	return function.Bind2nd(MonadAlt[E, A], second)
}

// MonadFlap applies a value to a function wrapped in an [SeqEither]
func MonadFlap[E, B, A any](fab SeqEither[E, func(A) B], a A) SeqEither[E, B] {
	return functor.MonadFlap(MonadMap[E, func(A) B, B], fab, a)
}

// Flap returns a function that applies a value to a function wrapped in an [SeqEither]
func Flap[E, B, A any](a A) Operator[E, func(A) B, B] {
	return functor.Flap(Map[E, func(A) B, B], a)
}

// MonadChainLeft chains a computation on the left (error) side of a SeqEither.
// If the input is a Left value, it applies the function f to transform the error and potentially
// change the error type from EA to EB. If the input is a Right value, it passes through unchanged.
//
// Note: MonadChainLeft is identical to OrElse - both provide the same functionality for error recovery.
//
// This is useful for error recovery or error transformation scenarios where you want to handle
// errors by performing another computation that may also fail.
//
// Parameters:
//   - fa: The input SeqEither that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a SeqEither with error type EB
//
// Returns:
//   - A SeqEither with the potentially transformed error type EB
//
// Example:
//
//	// Recover from a specific error by trying an alternative computation
//	result := MonadChainLeft(
//	    Left[int]("network error"),
//	    func(err string) SeqEither[string, int] {
//	        if err == "network error" {
//	            return Right[string](42) // recover with default value
//	        }
//	        return Left[int]("unrecoverable: " + err)
//	    },
//	)
func MonadChainLeft[EA, EB, A any](fa SeqEither[EA, A], f Kleisli[EB, EA, A]) SeqEither[EB, A] {
	return eithert.MonadChainLeft(
		iter.MonadChain[Either[EA, A], Either[EB, A]],
		iter.MonadOf[Either[EB, A]],
		fa,
		f,
	)
}

// ChainLeft is the curried version of MonadChainLeft.
// It returns a function that chains a computation on the left (error) side of a SeqEither.
//
// Note: ChainLeft is identical to OrElse - both provide the same functionality for error recovery.
//
// This is particularly useful in functional composition pipelines where you want to handle
// errors by performing another computation that may also fail.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a SeqEither with error type EB
//
// Returns:
//   - A function that transforms a SeqEither with error type EA to one with error type EB
//
// Example:
//
//	// Create a reusable error handler
//	recoverFromNetworkError := ChainLeft(func(err string) SeqEither[string, int] {
//	    if strings.Contains(err, "network") {
//	        return Right[string](0) // return default value
//	    }
//	    return Left[int](err) // propagate other errors
//	})
//
//	result := F.Pipe1(
//	    Left[int]("network timeout"),
//	    recoverFromNetworkError,
//	)
func ChainLeft[EA, EB, A any](f Kleisli[EB, EA, A]) func(SeqEither[EA, A]) SeqEither[EB, A] {
	return eithert.ChainLeft(
		iter.Chain[Either[EA, A], Either[EB, A]],
		iter.Of[Either[EB, A]],
		f,
	)
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
//   - ma: The input SeqEither that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a SeqEither (typically for side effects)
//
// Returns:
//   - A SeqEither with the original error preserved if input was Left, or the original Right value
//
// Example:
//
//	// Log errors but always preserve the original error
//	result := MonadChainFirstLeft(
//	    Left[int]("database error"),
//	    func(err string) SeqEither[string, int] {
//	        return FromIO[string](func() int {
//	            log.Printf("Error occurred: %s", err)
//	            return 0
//	        })
//	    },
//	)
//	// result will always be Left("database error"), even though f returns Right
func MonadChainFirstLeft[A, EA, EB, B any](ma SeqEither[EA, A], f Kleisli[EB, EA, B]) SeqEither[EA, A] {
	return eithert.MonadChainFirstLeft(
		iter.MonadChain[Either[EA, A], Either[EA, A]],
		iter.MonadMap[Either[EB, B], Either[EA, A]],
		iter.MonadOf[Either[EA, A]],
		ma,
		f,
	)
}

// MonadTapLeft is an alias for [MonadChainFirstLeft]. It executes a
// side-effecting computation on each Left (error) value and always returns the
// original Left error unchanged. Right values pass through without calling f.
//
//go:inline
func MonadTapLeft[A, EA, EB, B any](ma SeqEither[EA, A], f Kleisli[EB, EA, B]) SeqEither[EA, A] {
	return MonadChainFirstLeft(ma, f)
}

// ChainFirstLeft is the curried version of MonadChainFirstLeft.
// It returns a function that chains a computation on the left (error) side while always preserving the original error.
//
// This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// ensuring the error path is preserved.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a SeqEither (typically for side effects)
//
// Returns:
//   - An Operator that performs the side effect but always returns the original error if input was Left
//
// Example:
//
//	// Create a reusable error logger
//	logError := ChainFirstLeft(func(err string) SeqEither[any, int] {
//	    return FromIO[any](func() int {
//	        log.Printf("Error: %s", err)
//	        return 0
//	    })
//	})
//
//	result := F.Pipe1(
//	    Left[int]("validation failed"),
//	    logError, // logs the error
//	)
//	// result is always Left("validation failed"), even though f returns Right
func ChainFirstLeft[A, EA, EB, B any](f Kleisli[EB, EA, B]) Operator[EA, A, A] {
	return eithert.ChainFirstLeft(
		iter.Chain[Either[EA, A], Either[EA, A]],
		iter.Map[Either[EB, B], Either[EA, A]],
		iter.Of[Either[EA, A]],
		f,
	)
}

// TapLeft is an alias for [ChainFirstLeft]. It returns an Operator that
// executes a side-effecting computation on each Left (error) value and always
// returns the original Left error unchanged. Right values pass through
// unchanged without calling f.
//
//go:inline
func TapLeft[A, EA, EB, B any](f Kleisli[EB, EA, B]) Operator[EA, A, A] {
	return ChainFirstLeft[A](f)
}

// OrElse recovers from a Left (error) by providing an alternative computation.
// If the SeqEither is Right, it returns the value unchanged.
// If the SeqEither is Left, it applies the provided function to the error value,
// which returns a new SeqEither that replaces the original.
//
// Note: OrElse is identical to ChainLeft - both provide the same functionality for error recovery.
//
// This is useful for error recovery, fallback logic, or chaining alternative IO computations.
// The error type can be widened from E1 to E2, allowing transformation of error types.
//
// Example:
//
//	// Recover from specific errors with fallback IO operations
//	recover := ioeither.OrElse(func(err error) ioeither.SeqEither[error, int] {
//	    if err.Error() == "not found" {
//	        return ioeither.Right[error](0) // default value
//	    }
//	    return ioeither.Left[int](err) // propagate other errors
//	})
//	result := recover(ioeither.Left[int](errors.New("not found"))) // Right(0)
//	result := recover(ioeither.Right[error](42)) // Right(42) - unchanged
//
//go:inline
func OrElse[E1, E2, A any](onLeft Kleisli[E2, E1, A]) Kleisli[E2, SeqEither[E1, A], A] {
	return Fold(onLeft, Of[E2, A])
}

// MonadReduce reduces a SeqEither to a single Either value by applying a function to each
// Right element and an accumulator, starting with an initial value. If any Left is encountered,
// reduction stops immediately and returns that Left value.
//
// Type Parameters:
//   - E: The error type
//   - A: The element type in the sequence
//   - B: The accumulator and result type
//
// Parameters:
//   - fa: The SeqEither to reduce
//   - f: The reducer function that combines the accumulator with each element
//   - initial: The initial accumulator value
//
// Returns:
//   - Either[E, B]: Left with the first error encountered, or Right with the final accumulated value
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
//	seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
//	sum := MonadReduce(seq, func(acc, x int) int { return acc + x }, 0)
//	// returns: E.Right[string](6)
//
//	seqWithError := iter.From(E.Right[string](1), E.Left[int]("error"), E.Right[string](3))
//	result := MonadReduce(seqWithError, func(acc, x int) int { return acc + x }, 0)
//	// returns: E.Left[int]("error")
//
//go:inline
func MonadReduce[E, A, B any](fa SeqEither[E, A], f func(B, A) B, initial B) IOEither[E, B] {
	return func() Either[E, B] {
		current := initial
		for ea := range fa {
			a, e := either.Unwrap(ea)
			if either.IsLeft(ea) {
				return either.Left[B](e)
			}
			current = f(current, a)
		}
		return either.Of[E](current)
	}
}

// Reduce returns a function that reduces a SeqEither to a single Either value.
// This is the curried version of MonadReduce.
//
// The return type is IOEither because a SeqEither represents a dynamic lazy sequence of
// Either values. Reducing it requires consuming the iterator, so the final Either must
// also remain deferred as a lazy effect. IOEither is the single-value counterpart to
// SeqEither: both are lazy, and both represent work that is only performed when evaluated.
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
//   - A function that takes a SeqEither and returns IOEither[E, B]
//
// Example:
//
//	sum := Reduce(func(acc, x int) int { return acc + x }, 0)
//	seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
//	resultIO := sum(seq)
//	result := resultIO()
//	// returns: E.Right[string](6)
func Reduce[E, A, B any](f func(B, A) B, initial B) func(SeqEither[E, A]) IOEither[E, B] {
	return func(fa SeqEither[E, A]) IOEither[E, B] {
		return MonadReduce(fa, f, initial)
	}
}

// Collect materializes a SeqEither into an Either containing a slice.
// It consumes all elements from the sequence, collecting Right values into a slice.
// If any Left (error) is encountered during iteration, collection stops immediately
// and returns that Left value.
//
// The return type is IOEither because collecting a SeqEither is a reduction over a dynamic
// lazy source. The final Either is not available until iteration has happened, so it must
// stay deferred. IOEither is the single-value counterpart to SeqEither: both are lazy, and
// both represent work that is only performed when evaluated.
//
// This function eagerly evaluates the entire sequence and accumulates all successful
// values into a slice in memory. It short-circuits on the first error.
//
// Type Parameters:
//   - E: The error type
//   - T: The type of elements in the sequence
//
// Parameters:
//   - fa: The SeqEither to collect into a slice
//
// Returns:
//   - IOEither[E, []T]: A deferred computation that yields the first error encountered or a slice of all values
//
// Marble Diagram:
//
//	Input:  --R(1)--R(2)--R(3)--R(4)--R(5)--|
//	Output: IO(Right([1, 2, 3, 4, 5]))
//
//	Input:  --R(1)--R(2)--L(e)--R(4)--R(5)--|
//	Output: IO(Left(e))
//
// Example - Success case:
//
//	seq := iter.From(
//	    either.Right[string](1),
//	    either.Right[string](2),
//	    either.Right[string](3),
//	)
//	resultIO := Collect(seq)
//	result := resultIO()
//	// result = either.Right[string]([]int{1, 2, 3})
//
// Example - Error case:
//
//	seq := iter.From(
//	    either.Right[string](1),
//	    either.Left[int]("error"),
//	    either.Right[string](3),
//	)
//	resultIO := Collect(seq)
//	result := resultIO()
//	// result = either.Left[[]int]("error")
//
// Example - Empty sequence:
//
//	seq := iter.Empty[either.Either[string, int]]()
//	resultIO := Collect(seq)
//	result := resultIO()
//	// result = either.Right[string]([]int{})
//
// See Also:
//   - MonadReduce: Reduces a SeqEither to a single value
//   - Fold: Converts SeqEither to Seq by handling both cases
//   - GetOrElse: Extracts value or provides default
func Collect[E, T any](fa SeqEither[E, T]) IOEither[E, []T] {
	return MonadReduce(fa, A.Append, nil)
}

// MonadChainFirstIOK is the uncurried form of ChainFirstIOK.
// For each Right(a) in ma it executes the IO computation f(a) for its side
// effect, then discards the result and passes the original Right(a) through.
// Left values pass through unchanged without calling f.
//
// This is the SeqEither equivalent of iter.MonadChainFirstIOK: use it to
// attach IO side effects (logging, metrics, tracing) to Right values in a
// SeqEither pipeline without altering the values.
//
// Example:
//
//	var logged []int
//	logIO := func(n int) func() string {
//	    return func() string { logged = append(logged, n); return "ok" }
//	}
//	seq := iter.From(E.Right[string](1), E.Left[int]("err"), E.Right[string](3))
//	result := MonadChainFirstIOK(seq, logIO)
//	// yields: Right(1), Left("err"), Right(3)
//	// logged == [1, 3] (Left is skipped)
func MonadChainFirstIOK[E, A, B any](ma SeqEither[E, A], f io.Kleisli[A, B]) SeqEither[E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChainSeq[E, A, A],
		MonadMap[E, B, A],
		FromIO,
		ma,
		f,
	)
}

// ChainFirstIOK is the curried version of MonadChainFirstIOK.
// It returns an Operator that executes an IO-returning function for side
// effects on each Right value while passing the original element through
// unchanged. Left values are not affected.
//
// Example:
//
//	var seen []int
//	logOp := ChainFirstIOK[string](func(n int) func() string {
//	    return func() string { seen = append(seen, n); return "logged" }
//	})
//	seq := iter.From(E.Right[string](10), E.Right[string](20))
//	result := logOp(seq)
//	// yields: Right(10), Right(20)
//	// seen == [10, 20]
func ChainFirstIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, A] {
	return fromio.ChainFirstIOK(
		ChainSeq[E, A, A],
		Map[E, B, A],
		FromIO,
		f,
	)
}

// TapIOK is an alias for ChainFirstIOK.
// It executes an IO-returning function for its side effect on each Right value
// and passes the original element through unchanged.
//
// Example:
//
//	logOp := TapIOK[string](func(n int) func() string {
//	    return func() string { fmt.Printf("element: %d\n", n); return "ok" }
//	})
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := logOp(seq)
//	// yields: Right(1), Right(2) (and prints each element on iteration)
func TapIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, A] {
	return ChainFirstIOK[E](f)
}

// MonadChainIOK is the uncurried form of ChainIOK.
// For each Right(a) in ma it calls f(a) to obtain an IO[B], executes it,
// and lifts the result into Right(b), replacing the original element.
// Left values pass through unchanged without calling f.
//
// Example:
//
//	toStringIO := func(n int) func() string {
//	    return func() string { return fmt.Sprintf("item-%d", n) }
//	}
//	seq := iter.From(E.Right[string](1), E.Left[int]("err"), E.Right[string](3))
//	result := MonadChainIOK(seq, toStringIO)
//	// yields: Right("item-1"), Left("err"), Right("item-3")
func MonadChainIOK[E, A, B any](ma SeqEither[E, A], f io.Kleisli[A, B]) SeqEither[E, B] {
	return fromio.MonadChainIOK(
		MonadChainSeq[E, A, B],
		FromIO,
		ma,
		f,
	)
}

// ChainIOK is the curried version of MonadChainIOK.
// It returns an Operator that maps each Right value through an IO-producing
// function, replacing each Right(a) with Right(f(a)()), and leaving Left
// values unchanged.
//
// Example:
//
//	stringify := ChainIOK[string](func(n int) func() string {
//	    return func() string { return fmt.Sprintf("item-%d", n) }
//	})
//	seq := iter.From(E.Right[string](1), E.Right[string](2))
//	result := stringify(seq)
//	// yields: Right("item-1"), Right("item-2")
func ChainIOK[E, A, B any](f io.Kleisli[A, B]) Operator[E, A, B] {
	return fromio.ChainIOK(
		ChainSeq[E, A, B],
		FromIO,
		f,
	)
}

// MonadChainIOEitherK chains an [ioeither.Kleisli] function into a SeqEither
// pipeline. For each Right(a) in ma it calls f(a) to obtain an IOEither[E, B],
// executes it, and lifts the result into the sequence. Left values pass through
// unchanged without calling f.
//
// Marble diagram:
//
//	Input:  ---R(1)---R(2)---L(e)---|
//	f(1) = IO(Right(10)), f(2) = IO(Left("bad"))
//	Output: ---R(10)---L("bad")---L(e)---|
func MonadChainIOEitherK[E, A, B any](ma SeqEither[E, A], f ioeither.Kleisli[E, A, B]) SeqEither[E, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChainSeq[E, A, B],
		FromIOEither,
		ma,
		f,
	)
}

// ChainIOEitherK is the curried version of MonadChainIOEitherK. It returns an
// Operator that chains an [ioeither.Kleisli] function into a SeqEither pipeline.
// For each Right(a) in the input sequence it calls f(a), executes the resulting
// IOEither, and lifts the result. Left values pass through unchanged.
func ChainIOEitherK[E, A, B any](f ioeither.Kleisli[E, A, B]) Operator[E, A, B] {
	return fromioeither.ChainIOEitherK(
		ChainSeq[E, A, B],
		FromIOEither,
		f,
	)
}

// MonadChainFirstIOEitherK executes a side-effecting [ioeither.Kleisli]
// computation for each Right value and returns the original value unchanged.
// For each Right(a) in ma it calls f(a), executes the resulting IOEither as a
// side effect, then passes the original Right(a) through. Left values pass
// through unchanged without calling f.
func MonadChainFirstIOEitherK[E, A, B any](ma SeqEither[E, A], f ioeither.Kleisli[E, A, B]) SeqEither[E, A] {
	return fromioeither.MonadChainFirstIOEitherK(
		MonadChainSeq[E, A, A],
		MonadMap[E, B, A],
		FromIOEither,
		ma,
		f,
	)
}

// ChainFirstIOEitherK is the curried version of MonadChainFirstIOEitherK. It
// returns an Operator that executes a side-effecting [ioeither.Kleisli]
// computation on each Right value and always returns the original value
// unchanged. Left values pass through without calling f.
func ChainFirstIOEitherK[E, A, B any](f ioeither.Kleisli[E, A, B]) Operator[E, A, A] {
	return fromioeither.ChainFirstIOEitherK(
		ChainSeq[E, A, A],
		Map[E, B, A],
		FromIOEither,
		f,
	)
}

// TapIOEitherK is an alias for [ChainFirstIOEitherK]. It executes a
// side-effecting [ioeither.Kleisli] computation on each Right value and always
// returns the original value unchanged. Left values pass through without
// calling f.
func TapIOEitherK[E, A, B any](f ioeither.Kleisli[E, A, B]) Operator[E, A, A] {
	return ChainFirstIOEitherK(f)
}
