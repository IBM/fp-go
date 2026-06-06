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

package iter

import (
	"slices"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// ConcatBuf concatenates multiple sequences into a single sequence with
// deterministic output order. All inner producers run concurrently (each in its
// own goroutine), but a drainer goroutine reads them strictly left-to-right, so
// the consumer always observes elements from iterables[0] before iterables[1],
// iterables[1] before iterables[2], and so on.
//
// bufSize controls the capacity of each per-sequence channel and of the inners
// coordination channel. A larger value lets producers run further ahead of the
// drainer, reducing synchronisation overhead at the cost of memory. Negative
// values are clamped to 0 (unbuffered).
//
// This is the order-preserving, concurrent-producers counterpart of MergeBuf.
// Use it when you need deterministic output order but still want producers to
// run in parallel (e.g. multiple I/O-bound sequences).
//
// Marble Diagram:
//
//	Seq1 (goroutine): --1--2--3--|
//	Seq2 (goroutine): --4--5--6--|   (runs concurrently with Seq1)
//	ConcatBuf output: --1--2--3--4--5--6--|   (drained in order)
//
// Comparison with MergeBuf:
//
//	ConcatBuf: concurrent producers, sequential output
//	  Result: 1, 2, 3, 4, 5, 6  (always)
//
//	MergeBuf: concurrent producers, non-deterministic output
//	  Result: 1, 4, 2, 5, 3, 6  (order varies)
//
// Example:
//
//	result := ConcatBuf([]Seq[int]{From(1, 2, 3), From(4, 5, 6)}, 8)
//	for v := range result {
//	    fmt.Println(v) // prints: 1, 2, 3, 4, 5, 6
//	}
//
// See Also:
//   - Concat: Convenience wrapper using defaultBufferSize
//   - ConcatAll: Operator form — takes Seq[Seq[T]] instead of a slice
//   - MergeBuf: Concurrent version with non-deterministic order
func ConcatBuf[T any](iterables []Seq[T], bufSize int) Seq[T] {
	return F.Pipe2(
		iterables,
		slices.Values,
		ConcatAll[T](bufSize),
	)
}

// Concat is a convenience wrapper around ConcatBuf that uses defaultBufferSize.
// It runs all inner producers concurrently but drains them in order, yielding a
// deterministic output sequence.
//
// See ConcatBuf for full semantics and the bufSize trade-offs.
//
// See Also:
//   - ConcatBuf: Customisable buffer size
//   - Merge: Concurrent, non-deterministic order
//
//go:inline
func Concat[T any](iterables []Seq[T]) Seq[T] {
	return ConcatBuf(iterables, defaultBufferSize)
}

// ConcatSeq concatenates multiple sequences into a single sequence using
// purely sequential nested iteration — no goroutines, no channels. Each input
// sequence is fully drained before the next one is started.
//
// This is the goroutine-free counterpart of Concat. Use it when the inner
// sequences are cheap and synchronous and goroutine overhead is undesirable.
//
// See Also:
//   - Concat: Concurrent producers, same output order, uses defaultBufferSize
//   - ConcatPar: Always concurrent (bypasses the sequential optimisation)
//   - ConcatBuf: Custom buffer size
//
//go:inline
func ConcatSeq[T any](iterables []Seq[T]) Seq[T] {
	return F.Pipe2(
		iterables,
		slices.Values,
		ConcatAllSeq[T](),
	)
}

// ConcatPar concatenates multiple sequences into a single sequence using
// concurrent inner producers drained in order, always via ConcatAllPar with
// defaultBufferSize. Unlike Concat, it never selects the goroutine-free
// sequential path.
//
// Use this when you need concurrent producers regardless of buffer size, for
// example to ensure forward progress in I/O-bound pipelines.
//
// See Also:
//   - Concat: Dispatches to sequential when bufSize == 1
//   - ConcatSeq: Always sequential, no goroutines
//   - ConcatBuf: Custom buffer size
//
//go:inline
func ConcatPar[T any](iterables []Seq[T]) Seq[T] {
	return F.Pipe2(
		iterables,
		slices.Values,
		ConcatAllPar[T](defaultBufferSize),
	)
}

// ConcatMapBuf maps each element to a sequence via f and flattens the results
// with deterministic output order. Each mapped sequence runs in its own goroutine
// (concurrent producers), but a drainer goroutine forwards values to the consumer
// strictly in input order, so all elements of f(a₁) appear before any element of
// f(a₂).
//
// bufSize controls the capacity of each per-mapped-sequence channel. A larger
// value lets producers run further ahead of the drainer. Negative values are
// clamped to 0 (unbuffered, fully synchronised).
//
// Because producers run concurrently, side effects inside f or inside the returned
// sequences may execute in any order. Only the *output* order is guaranteed.
// Use MergeMapBuf when output order does not matter and you want maximum throughput.
//
// Marble Diagram:
//
//	Input:  --1------2------3----->
//	f(x) => [x, x*10]            (each f(x) runs in its own goroutine)
//	Output: --1-10---2-20---3-30-> (drained in input order)
//
// Example:
//
//	expand := ConcatMapBuf(func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	}, 8)
//	// Always yields: 1, 10, 100, 2, 20, 200, 3, 30, 300
//	for v := range expand(From(1, 2, 3)) {
//	    fmt.Println(v)
//	}
//
// See Also:
//   - ConcatMap: Convenience wrapper using defaultBufferSize
//   - MergeMapBuf: Concurrent, non-deterministic output order
//
//go:inline
func ConcatMapBuf[A, B any](f Kleisli[A, B], bufSize int) Operator[A, B] {
	return F.Flow2(
		Map(f),
		ConcatAll[B](bufSize),
	)
}

// ConcatMapBufPar maps each element to a sequence via f and flattens the
// results using ConcatAllPar directly, bypassing the bufSize == 1 sequential
// optimisation of ConcatMapBuf. All mapped sequences always run in their own
// goroutines and are drained in input order.
//
// Use this variant when you need the concurrent-producers model regardless of
// bufSize, for example to ensure forward progress in I/O-bound pipelines even
// when bufSize == 1.
//
// See ConcatMapBuf for full semantics. The only difference is that this
// function never selects the goroutine-free sequential path.
//
// See Also:
//   - ConcatMapBuf: Dispatches to sequential for bufSize == 1
//   - ConcatMapPar: Convenience wrapper using defaultBufferSize
//   - MergeMapBuf: Concurrent, non-deterministic output order
//
//go:inline
func ConcatMapBufPar[A, B any](f Kleisli[A, B], bufSize int) Operator[A, B] {
	return F.Flow2(
		Map(f),
		ConcatAllPar[B](bufSize),
	)
}

// MonadConcatMap is the uncurried form of ConcatMap. It maps f over fa and
// flattens the results with deterministic output order using concurrent inner
// producers (see ConcatMapBuf for the full semantics).
//
// See Also:
//   - ConcatMap: The curried/operator form
//   - MonadChain: Identical operation under a different name
//   - MonadMergeMap: Non-deterministic output order
func MonadConcatMap[A, B any](fa Seq[A], f Kleisli[A, B]) Seq[B] {
	return ConcatMap(f)(fa)
}

// MonadConcatMapSeq is the uncurried form of ConcatMapSeq. It maps f over fa
// and flattens the results using purely sequential iteration — no goroutines,
// no channels. Each inner sequence is fully drained before f is called with
// the next element.
//
// See Also:
//   - ConcatMapSeq: The curried/operator form
//   - MonadConcatMap: Uses the bufSize-dispatched path (concurrent by default)
//   - MonadConcatMapPar: Always concurrent
func MonadConcatMapSeq[A, B any](fa Seq[A], f Kleisli[A, B]) Seq[B] {
	return ConcatMapSeq(f)(fa)
}

// MonadConcatMapPar is the uncurried form of ConcatMapPar. It maps f over fa
// and flattens the results using ConcatAllPar with defaultBufferSize — all
// mapped sequences run concurrently in their own goroutines and are drained in
// input order.
//
// See Also:
//   - ConcatMapPar: The curried/operator form
//   - MonadConcatMap: Dispatches to sequential for small bufSizes
//   - MonadConcatMapSeq: Always sequential
func MonadConcatMapPar[A, B any](fa Seq[A], f Kleisli[A, B]) Seq[B] {
	return ConcatMapPar(f)(fa)
}

// ConcatAllPar flattens a sequence of sequences into a single sequence with
// deterministic output order using three goroutine roles per invocation.
//
// # Concurrency model
//
//  1. Outer (O): iterates the outer sequence s.  For every inner Seq it spawns
//     an Iₙ goroutine and passes its output channel out_n to the Drainer via
//     the inners coordination channel.  O closes inners when s is exhausted or
//     done is closed.
//
//  2. Inner producers (Iₙ, one per inner Seq): drain one inner Seq into a
//     dedicated buffered channel out_n, then close it.  Not tracked by a
//     WaitGroup — they exit as soon as their Seq is exhausted or done is closed.
//
//  3. Drainer (D): for-ranges over inners in arrival order.  For each out_n it
//     drains all values into the output channel ch, then moves to the next
//     channel.  D closes ch when inners is closed (i.e. after O exits).
//
// # Defer ordering
//
// O defers close(inners): fires when O exits, signalling D that no more
// inner channels will arrive.
//
// Each Iₙ defers close(out_n): fires when Iₙ exits, signalling D that out_n
// is fully drained and it can move to the next channel.
//
// D defers close(ch): fires after the for-range over inners exits (i.e. after
// O has closed inners and all out_n channels have been drained), signalling the
// consumer that no more values will arrive.
//
// # Cancellation via done
//
// The consumer defers close(done).  All goroutines observe cancellation via select:
//   - O: checks done when handing out_n to inners.
//   - Each Iₙ: checks done when sending a value to out_n.
//   - D: checks done when forwarding a value to ch.
//
// When done is closed each goroutine reaches its done-check on the next send or
// receive and returns.  Their defers then close their owned channels in
// dependency order (Iₙ closes out_n → D drains out_n and closes ch).
// Goroutines may run briefly after the consumer returns; the done channel
// ensures they terminate promptly.
//
// # Consumer cleanup
//
//	defer close(done) // signal all goroutines to stop
//	for v := range ch {
//	    if !yield(v) { return }
//	}
//
// # Note on RxJS
//
// The RxJS concatAll operator subscribes to inner observables one at a time
// (strictly sequential).  This implementation starts all Iₙ goroutines
// concurrently for lower latency while preserving output order via the drainer.
//
// # Parameters
//
//   - bufSize: capacity of each out_n and of the inners coordination channel.
//     Larger values let producers run further ahead of the drainer.
//     Negative values are clamped to 0 (unbuffered).
//
// Marble Diagram:
//
//	Seq1 (goroutine I₁): --1--2--3--|
//	Seq2 (goroutine I₂): --4--5--6--|   (concurrent with I₁)
//	Drainer output:       --1--2--3--4--5--6--|
func ConcatAllPar[T any](bufSize int) Operator[Seq[T], T] {
	buf := N.Max(bufSize, 0)
	return func(s Seq[Seq[T]]) Seq[T] {
		return func(yield func(T) bool) {
			ch := make(chan T, buf)
			done := make(chan Void)
			inners := make(chan chan T, buf)

			go func() {
				defer close(inners)
				for inner := range s {
					out := make(chan T, buf)
					go func() {
						defer close(out)
						for v := range inner {
							select {
							case out <- v:
							case <-done:
								return
							}
						}
					}()
					select {
					case inners <- out:
					case <-done:
						return
					}
				}
			}()

			go func() {
				defer close(ch)
				for out := range inners {
					for v := range out {
						select {
						case ch <- v:
						case <-done:
							return
						}
					}
				}
			}()

			defer close(done)
			for v := range ch {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ConcatAllSeq flattens a Seq[Seq[T]] into a Seq[T] using purely sequential
// nested iteration — no goroutines, no channels. Each inner sequence is fully
// drained before the next outer element is consumed.
//
// This is the implementation selected by ConcatAll when bufSize == 1. It
// eliminates all synchronisation overhead, making it the fastest choice for
// pipelines that are already single-threaded or where the inner sequences are
// cheap to produce synchronously.
//
// Compare with ConcatAllPar, which spawns one goroutine per inner sequence and
// forwards results to the consumer in order via a drainer goroutine. Even
// ConcatAllPar(0) (unbuffered channels) still incurs goroutine-spawn and
// channel overhead; ConcatAllSeq avoids all of it.
//
// Marble Diagram (sequential — seq[n+1] starts only after seq[n] is drained):
//
//	Seq1: --1--2--3--|
//	Seq2:             --4--5--6--|
//	Output: --1--2--3--4--5--6--|
//
// Example:
//
//	outer := From(From(1, 2, 3), From(4, 5, 6))
//	for v := range ConcatAllSeq[int]()(outer) {
//	    fmt.Println(v) // prints: 1, 2, 3, 4, 5, 6
//	}
//
// See Also:
//   - ConcatAll: Dispatcher — selects ConcatAllSeq for bufSize == 1
//   - ConcatAllPar: Concurrent inner producers, sequential consumer output
func ConcatAllSeq[T any]() Operator[Seq[T], T] {
	return func(s Seq[Seq[T]]) Seq[T] {
		return func(yield func(T) bool) {
			for outer := range s {
				for inner := range outer {
					if !yield(inner) {
						return
					}
				}
			}
		}
	}
}

// ConcatAll flattens a Seq[Seq[T]] into a Seq[T] with deterministic output
// order, dispatching to the most efficient implementation for the requested
// buffer size:
//
//   - bufSize == 1: delegates to ConcatAllSeq — purely sequential iteration,
//     zero goroutines, zero channels. Each inner sequence is fully consumed
//     before the next one is started. Best when concurrency adds no value or
//     when goroutine overhead must be avoided.
//
//   - all other values (including 0 and negative): delegates to
//     ConcatAllPar(bufSize) — all inner producers run concurrently (one
//     goroutine each), a drainer goroutine forwards values in arrival order,
//     channels have capacity max(bufSize, 0). Negative values → unbuffered.
//
// The dispatch is non-monotonic with respect to concurrency: bufSize == 0
// (unbuffered) still uses goroutines and channels, whereas bufSize == 1 is
// fully goroutine-free and sequential. This is intentional: the
// goroutine-free path is selected only at bufSize == 1 as an explicit
// performance mode, not as a generalisation of "small buffer".
//
// See Also:
//   - ConcatAllSeq: Sequential implementation (no goroutines)
//   - ConcatAllPar: Concurrent implementation with configurable buffer
//   - ConcatBuf: Slice-input convenience wrapper
func ConcatAll[T any](bufSize int) Operator[Seq[T], T] {
	switch bufSize {
	case 1:
		return ConcatAllSeq[T]()
	default:
		return ConcatAllPar[T](bufSize)
	}
}
