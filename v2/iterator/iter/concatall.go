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
	"sync"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// ConcatBuf concatenates multiple sequences into a single sequence with
// deterministic output order. All inner producers run concurrently (each in its
// own goroutine), but a drainer goroutine reads them strictly left-to-right, so
// the consumer always observes elements from iterables[0] before iterables[1],
// iterables[1] before iterables[2], and so on.
//
// bufSize controls the capacity of each per-sequence channel and of the
// innerChannels coordination channel. A larger value lets producers run further
// ahead of the drainer, reducing synchronisation overhead at the cost of memory.
// Negative values are clamped to 0 (unbuffered).
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

// ConcatAll flattens a sequence of sequences into a single sequence.
//
// # Concurrency model
//
// Four goroutine roles are involved for each invocation:
//
//  1. Outer (O): iterates the outer sequence s.  For every inner Seq it spawns
//     an inner producer goroutine and hands its channel to the Drainer via
//     innerChannels.
//
//  2. Inner producers (Iₙ, one per inner Seq): pull values from one inner Seq
//     and push them into a dedicated buffered channel innerCh_n.  Each Iₙ is
//     tracked by the WaitGroup.
//
//  3. Drainer (D): reads innerCh_n channels from innerChannels strictly in
//     arrival order and forwards their values to the single output channel ch.
//     This is the mechanism that guarantees output order even though all Iₙ run
//     in parallel.
//
//  4. Closer (C): calls wg.Wait() and then closes ch once every tracked
//     goroutine (O, all Iₙ, D) has exited.  C is not itself tracked by wg.
//
// # WaitGroup invariant
//
// O and D are registered (wg.Add) before C starts.  Therefore wg ≥ 2 when C
// begins waiting, so C cannot see a spurious zero caused by O finishing quickly
// before D is registered.  Each Iₙ is registered inside O's loop body while O
// still holds wg ≥ 1, so the counter never touches zero between registrations.
//
// # Defer ordering inside O and each Iₙ
//
// Go defers are LIFO, so inside O:
//   - close(innerChannels) runs first  → D's for-range exits
//   - wg.Done() runs second            → C unblocks only after D has seen the close
//
// Inside each Iₙ:
//   - close(innerCh_n) runs first      → D's inner for-range exits
//   - wg.Done() runs second            → C unblocks only after innerCh_n is drained
//
// # Cancellation via done
//
// All goroutines check the done channel via select:
//   - O checks at the top of each callback invocation and when handing innerCh
//     to innerChannels.
//   - Each Iₙ checks when sending a value to its innerCh.
//   - D checks between consecutive inner channels and when forwarding to ch.
//
// When done is closed, every goroutine reaches its done-check and returns.
// Their defers close their owned channels and call wg.Done().  After all have
// exited, C calls close(ch), which unblocks the consumer's drain loop.
//
// # Consumer cleanup
//
// The consumer's defer (see below) runs on both normal and early exit:
//
//	defer func() {
//	    close(done)      // signal all goroutines to stop
//	    for range ch {}  // drain until C closes ch (blocking until all goroutines exit)
//	}()
//
// This ensures the consumer never returns while any goroutine is still live.
// On normal completion done is closed after ch is already shut; the drain loop
// exits immediately and the close is harmless to the already-finished goroutines.
//
// # Note on RxJS
//
// The RxJS concatAll operator subscribes to inner observables one at a time
// (strictly sequential).  This implementation starts all inner goroutines
// concurrently for lower latency while preserving output order via the drainer.
//
// # Parameters
//
//   - bufSize: capacity of each innerCh_n and of the innerChannels coordination
//     channel.  Larger values let producers run further ahead of the drainer.
//     Negative values are clamped to 0 (unbuffered).
//
// Marble Diagram:
//
//	Seq1 (goroutine I₁): --1--2--3--|
//	Seq2 (goroutine I₂): --4--5--6--|   (concurrent with I₁)
//	Drainer output:       --1--2--3--4--5--6--|
func ConcatAll[T any](bufSize int) Operator[Seq[T], T] {
	buf := N.Max(bufSize, 0)

	return func(s Seq[Seq[T]]) Seq[T] {
		return func(yield func(T) bool) {

			// ch is the single output channel read by the consumer.
			// Closed by the Closer (C) after wg reaches zero.
			ch := make(chan T, buf)

			// done is closed by the consumer's defer to cancel all goroutines.
			done := make(chan Void)

			// innerChannels delivers per-sequence channels from O to D in order.
			innerChannels := make(chan chan T, buf)

			// wg tracks O + all Iₙ + D.  C calls wg.Wait() to know when to
			// close ch.  O and D are added before C starts (see invariant above).
			var wg sync.WaitGroup

			// O — Outer goroutine.
			// Iterates s, spawns one Iₙ per inner Seq, hands its channel to D.
			// close(innerChannels) runs before wg.Done() (LIFO defers).
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(innerChannels)
				s(func(inner Seq[T]) bool {
					select {
					case <-done:
						return false
					default:
					}
					innerCh := make(chan T, buf)
					// Iₙ — Inner producer goroutine.
					// Pulls from inner Seq, pushes to innerCh.
					// close(innerCh) runs before wg.Done() (LIFO defers).
					wg.Add(1)
					go func(seq Seq[T], out chan T) {
						defer wg.Done()
						defer close(out)
						seq(func(v T) bool {
							select {
							case out <- v:
								return true
							case <-done:
								return false
							}
						})
					}(inner, innerCh)
					select {
					case innerChannels <- innerCh:
						return true
					case <-done:
						return false
					}
				})
			}()

			// D — Drainer goroutine.
			// Reads innerCh channels in arrival order and forwards values to ch.
			// Registered before C starts (see WaitGroup invariant above).
			wg.Add(1)
			go func() {
				defer wg.Done()
				for innerCh := range innerChannels {
					// Check cancellation between consecutive inner channels.
					select {
					case <-done:
						return
					default:
					}
					for v := range innerCh {
						// innerCh is closed by Iₙ's defer when Iₙ exits, so
						// this for-range terminates even without checking done.
						select {
						case ch <- v:
						case <-done:
							return
						}
					}
				}
			}()

			// C — Closer goroutine (not tracked by wg).
			// Closes ch once every tracked goroutine has called wg.Done().
			// Must start after all fixed wg.Add calls (O and D above).
			go func() {
				wg.Wait()
				close(ch)
			}()

			// Consumer cleanup: signal cancellation, then drain ch until C
			// closes it.  This blocks until every goroutine has exited,
			// guaranteeing no goroutine outlives the consumer.
			defer func() {
				close(done)
				for range ch {
				}
			}()

			for v := range ch {
				if !yield(v) {
					return
				}
			}
		}
	}
}
