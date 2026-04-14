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

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
)

const (
	defaultBufferSize = 8
)

// MergeBuf merges multiple sequences concurrently into a single sequence.
// It spawns a goroutine for each input sequence and merges their elements through
// a buffered channel, allowing concurrent production from all sources. The output
// order is non-deterministic and depends on the timing of concurrent producers.
//
// This function is useful for combining results from multiple concurrent operations,
// processing data from multiple sources in parallel, or implementing fan-in patterns
// where multiple producers feed into a single consumer.
//
// Type Parameters:
//   - T: The type of elements in the sequences
//
// Parameters:
//   - iterables: A slice of sequences to merge. If empty, returns an empty sequence.
//   - bufSize: The buffer size for the internal channel. Negative values are treated as 0 (unbuffered).
//     A larger buffer allows more elements to be produced ahead of consumption,
//     reducing contention between producers but using more memory.
//     A buffer of 0 creates an unbuffered channel requiring synchronization.
//
// Returns:
//   - Seq[T]: A new sequence that yields elements from all input sequences in non-deterministic order
//
// Behavior:
//   - Spawns one goroutine per input sequence to produce elements concurrently
//   - Elements from different sequences are interleaved non-deterministically
//   - Properly handles early termination: if the consumer stops iterating (yield returns false),
//     all producer goroutines are signaled to stop and cleaned up
//   - The output channel is closed when all input sequences are exhausted
//   - No goroutines leak even with early termination
//   - Thread-safe: multiple producers can safely send to the shared channel
//
// Example Usage:
//
//	// MergeBuf three sequences concurrently
//	seq1 := From(1, 2, 3)
//	seq2 := From(4, 5, 6)
//	seq3 := From(7, 8, 9)
//	merged := MergeBuf([]Seq[int]{seq1, seq2, seq3}, 10)
//
//	// Elements appear in non-deterministic order
//	for v := range merged {
//	    fmt.Println(v) // May print: 1, 4, 7, 2, 5, 8, 3, 6, 9 (order varies)
//	}
//
// Example with Early Termination:
//
//	seq1 := From(1, 2, 3, 4, 5)
//	seq2 := From(6, 7, 8, 9, 10)
//	merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
//
//	// Stop after 3 elements - all producer goroutines will be properly cleaned up
//	count := 0
//	for v := range merged {
//	    fmt.Println(v)
//	    count++
//	    if count >= 3 {
//	        break
//	    }
//	}
//
// Example with Unbuffered Channel:
//
//	// bufSize of 0 creates an unbuffered channel
//	seq1 := From(1, 2, 3)
//	seq2 := From(4, 5, 6)
//	merged := MergeBuf([]Seq[int]{seq1, seq2}, 0)
//
//	// Producers and consumer are synchronized
//	for v := range merged {
//	    fmt.Println(v)
//	}
//
// See Also:
//   - Async: Converts a single sequence to asynchronous
//   - From: Creates a sequence from values
//   - MonadChain: Sequentially chains sequences (deterministic order)
func MergeBuf[T any](iterables []Seq[T], bufSize int) Seq[T] {
	return F.Pipe2(
		iterables,
		slices.Values,
		MergeAll[T](bufSize),
	)
}

// Merge merges multiple sequences concurrently into a single sequence using a default buffer size.
// This is a convenience wrapper around MergeBuf that uses a default buffer size of 8.
//
// Type Parameters:
//   - T: The type of elements in the sequences
//
// Parameters:
//   - iterables: A slice of sequences to merge. If empty, returns an empty sequence.
//
// Returns:
//   - Seq[T]: A new sequence that yields elements from all input sequences in non-deterministic order
//
// Behavior:
//   - Uses a default buffer size of 8 for the internal channel
//   - Spawns one goroutine per input sequence to produce elements concurrently
//   - Elements from different sequences are interleaved non-deterministically
//   - Properly handles early termination with goroutine cleanup
//   - Thread-safe: multiple producers can safely send to the shared channel
//
// Example:
//
//	seq1 := From(1, 2, 3)
//	seq2 := From(4, 5, 6)
//	seq3 := From(7, 8, 9)
//	merged := Merge([]Seq[int]{seq1, seq2, seq3})
//
//	// Elements appear in non-deterministic order
//	for v := range merged {
//	    fmt.Println(v) // May print: 1, 4, 7, 2, 5, 8, 3, 6, 9 (order varies)
//	}
//
// See Also:
//   - MergeBuf: Merge with custom buffer size
//   - MergeAll: Merges a sequence of sequences
//   - Async: Converts a single sequence to asynchronous
func Merge[T any](iterables []Seq[T]) Seq[T] {
	return MergeBuf(iterables, defaultBufferSize)
}

// MergeMonoid creates a Monoid for merging sequences concurrently.
// The monoid combines two sequences by merging them concurrently with the specified
// buffer size, and uses an empty sequence as the identity element.
//
// A Monoid is an algebraic structure with an associative binary operation (concat)
// and an identity element (empty). For sequences, the concat operation merges two
// sequences concurrently, and the identity is an empty sequence.
//
// This is useful for functional composition patterns where you need to combine
// multiple sequences using monoid operations like Reduce, FoldMap, or when working
// with monadic operations that require a monoid instance.
//
// Marble Diagram (Concurrent Merging):
//
//	Seq1:   --1--2--3--|
//	Seq2:   --4--5--6--|
//	Merge:  --1-4-2-5-3-6--|
//	        (non-deterministic order)
//
// Marble Diagram (vs ConcatMonoid):
//
//	MergeMonoid (concurrent):
//	  Seq1:   --1--2--3--|
//	  Seq2:   --4--5--6--|
//	  Result: --1-4-2-5-3-6--|
//	          (elements interleaved)
//
//	ConcatMonoid (sequential):
//	  Seq1:   --1--2--3--|
//	  Seq2:              --4--5--6--|
//	  Result: --1--2--3--4--5--6--|
//	          (deterministic order)
//
// Type Parameters:
//   - T: The type of elements in the sequences
//
// Parameters:
//   - bufSize: The buffer size for the internal channel used during merging.
//     This buffer size will be used for all merge operations performed by the monoid.
//     Negative values are treated as 0 (unbuffered).
//
// Returns:
//   - Monoid[Seq[T]]: A monoid instance with:
//   - Concat: Merges two sequences concurrently using Merge
//   - Empty: Returns an empty sequence
//
// Properties:
//   - Identity: concat(empty, x) = concat(x, empty) = x
//   - Associativity: concat(concat(a, b), c) = concat(a, concat(b, c))
//     Note: Due to concurrent execution, element order may vary between equivalent expressions
//
// Example Usage:
//
//	// Create a monoid for merging integer sequences
//	monoid := MergeMonoid[int](10)
//
//	// Use with Reduce to merge multiple sequences
//	sequences := []Seq[int]{
//	    From(1, 2, 3),
//	    From(4, 5, 6),
//	    From(7, 8, 9),
//	}
//	merged := MonadReduce(From(sequences...), monoid.Concat, monoid.Empty)
//	// merged contains all elements from all sequences (order non-deterministic)
//
// Example with Empty Identity:
//
//	monoid := MergeMonoid[int](5)
//	seq := From(1, 2, 3)
//
//	// Merging with empty is identity
//	result1 := monoid.Concat(monoid.Empty, seq)  // same as seq
//	result2 := monoid.Concat(seq, monoid.Empty)  // same as seq
//
// Example with FoldMap:
//
//	// Convert each number to a sequence and merge all results
//	monoid := MergeMonoid[int](10)
//	numbers := From(1, 2, 3)
//	result := MonadFoldMap(numbers, func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	}, monoid)
//	// result contains: 1, 10, 100, 2, 20, 200, 3, 30, 300 (order varies)
//
// See Also:
//   - Merge: The underlying merge function
//   - MergeAll: Merges multiple sequences at once
//   - Empty: Creates an empty sequence
func MergeMonoid[T any](bufSize int) M.Monoid[Seq[T]] {
	return M.MakeMonoid(
		func(l, r Seq[T]) Seq[T] {
			return MergeBuf(A.From(l, r), bufSize)
		},
		Empty[T](),
	)
}

// MergeAll creates an operator that flattens and merges a sequence of sequences concurrently.
// It takes a sequence of sequences (Seq[Seq[T]]) and produces a single flat sequence (Seq[T])
// by spawning a goroutine for each inner sequence as it arrives, merging all their elements
// through a buffered channel. This enables dynamic concurrent processing where inner sequences
// can be produced and consumed concurrently.
//
// Unlike Merge which takes a pre-defined slice of sequences, MergeAll processes sequences
// dynamically as they are produced by the outer sequence. This makes it ideal for scenarios
// where the number of sequences isn't known upfront or where sequences are generated on-the-fly.
//
// Type Parameters:
//   - T: The type of elements in the inner sequences
//
// Parameters:
//   - bufSize: The buffer size for the internal channel. Negative values are treated as 0 (unbuffered).
//     A larger buffer allows more elements to be produced ahead of consumption,
//     reducing contention between producers but using more memory.
//
// Returns:
//   - Operator[Seq[T], T]: A function that takes a sequence of sequences and returns a flat sequence
//
// Behavior:
//   - Spawns one goroutine for the outer sequence to iterate and spawn inner producers
//   - Spawns one goroutine per inner sequence as it arrives from the outer sequence
//   - Elements from different inner sequences are interleaved non-deterministically
//   - Properly handles early termination: if the consumer stops iterating, all goroutines are cleaned up
//   - The output channel is closed when both the outer sequence and all inner sequences are exhausted
//   - No goroutines leak even with early termination
//   - Thread-safe: multiple producers can safely send to the shared channel
//
// Example Usage:
//
//	// Create a sequence of sequences dynamically
//	outer := From(
//	    From(1, 2, 3),
//	    From(4, 5, 6),
//	    From(7, 8, 9),
//	)
//	mergeAll := MergeAll[int](10)
//	merged := mergeAll(outer)
//
//	// Elements appear in non-deterministic order
//	for v := range merged {
//	    fmt.Println(v) // May print: 1, 4, 7, 2, 5, 8, 3, 6, 9 (order varies)
//	}
//
// Example with Dynamic Generation:
//
//	// Generate sequences on-the-fly
//	outer := Map(func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	})(From(1, 2, 3))
//	mergeAll := MergeAll[int](10)
//	merged := mergeAll(outer)
//
//	// Yields: 1, 10, 100, 2, 20, 200, 3, 30, 300 (order varies)
//	for v := range merged {
//	    fmt.Println(v)
//	}
//
// Example with Early Termination:
//
//	outer := From(
//	    From(1, 2, 3, 4, 5),
//	    From(6, 7, 8, 9, 10),
//	    From(11, 12, 13, 14, 15),
//	)
//	mergeAll := MergeAll[int](5)
//	merged := mergeAll(outer)
//
//	// Stop after 5 elements - all goroutines will be properly cleaned up
//	count := 0
//	for v := range merged {
//	    fmt.Println(v)
//	    count++
//	    if count >= 5 {
//	        break
//	    }
//	}
//
// Example with Chain:
//
//	// Use with Chain to flatten nested sequences
//	numbers := From(1, 2, 3)
//	result := Chain(func(n int) Seq[int] {
//	    return From(n, n*10)
//	})(numbers)
//	// This is equivalent to: MergeAll[int](0)(Map(...)(numbers))
//
// See Also:
//   - Merge: Merges a pre-defined slice of sequences
//   - Chain: Sequentially flattens sequences (deterministic order)
//   - Flatten: Flattens nested sequences sequentially
//   - Async: Converts a single sequence to asynchronous
func MergeAll[T any](bufSize int) Operator[Seq[T], T] {
	buf := N.Max(bufSize, 0)

	return func(s Seq[Seq[T]]) Seq[T] {

		return func(yield func(T) bool) {

			ch := make(chan T, buf)
			done := make(chan Void)
			var wg sync.WaitGroup

			// Outer producer: iterates the outer Seq and spawns an inner
			// goroutine for each inner Seq it emits.
			wg.Add(1)
			go func() {
				defer wg.Done()
				s(func(inner Seq[T]) bool {
					select {
					case <-done:
						return false
					default:
					}

					wg.Add(1)
					go func(seq Seq[T]) {
						defer wg.Done()
						seq(func(v T) bool {
							select {
							case ch <- v:
								return true
							case <-done:
								return false
							}
						})
					}(inner)

					return true
				})
			}()

			// Close ch once the outer producer and all inner producers finish.
			go func() {
				wg.Wait()
				close(ch)
			}()

			// On exit, signal cancellation and drain so no producer blocks
			// forever on `ch <- v`.
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

// MergeMapBuf applies a function that returns a sequence to each element and merges the results concurrently.
// This is the concurrent version of Chain (flatMap), where each mapped sequence is processed in parallel
// rather than sequentially. It combines Map and MergeAll into a single operation.
//
// Unlike Chain which processes sequences sequentially (deterministic order), MergeMapBuf spawns a goroutine
// for each mapped sequence and merges their elements concurrently through a buffered channel. This makes
// it ideal for I/O-bound operations, parallel data processing, or when the order of results doesn't matter.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of elements in the output sequences
//
// Parameters:
//   - f: A function that transforms each input element into a sequence of output elements
//   - bufSize: The buffer size for the internal channel. Negative values are treated as 0 (unbuffered).
//     A larger buffer allows more elements to be produced ahead of consumption,
//     reducing contention between producers but using more memory.
//
// Returns:
//   - Operator[A, B]: A function that takes a sequence of A and returns a flat sequence of B
//
// Behavior:
//   - Applies f to each element in the input sequence to produce inner sequences
//   - Spawns one goroutine per inner sequence to produce elements concurrently
//   - Elements from different inner sequences are interleaved non-deterministically
//   - Properly handles early termination: if the consumer stops iterating, all goroutines are cleaned up
//   - No goroutines leak even with early termination
//   - Thread-safe: multiple producers can safely send to the shared channel
//
// Comparison with Chain:
//   - Chain: Sequential processing, deterministic order, no concurrency overhead
//   - MergeMapBuf: Concurrent processing, non-deterministic order, better for I/O-bound tasks
//
// Example Usage:
//
//	// Expand each number into a sequence concurrently
//	expand := MergeMapBuf(func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	}, 10)
//	seq := From(1, 2, 3)
//	result := expand(seq)
//
//	// Yields: 1, 10, 100, 2, 20, 200, 3, 30, 300 (order varies)
//	for v := range result {
//	    fmt.Println(v)
//	}
//
// Example with I/O Operations:
//
//	// Fetch data concurrently for each ID
//	fetchData := MergeMapBuf(func(id int) Seq[string] {
//	    // Simulate I/O operation
//	    data := fetchFromAPI(id)
//	    return From(data...)
//	}, 20)
//	ids := From(1, 2, 3, 4, 5)
//	results := fetchData(ids)
//
//	// All fetches happen concurrently
//	for data := range results {
//	    fmt.Println(data)
//	}
//
// Example with Early Termination:
//
//	expand := MergeMapBuf(func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	}, 5)
//	seq := From(1, 2, 3, 4, 5)
//	result := expand(seq)
//
//	// Stop after 5 elements - all goroutines will be properly cleaned up
//	count := 0
//	for v := range result {
//	    fmt.Println(v)
//	    count++
//	    if count >= 5 {
//	        break
//	    }
//	}
//
// Example with Unbuffered Channel:
//
//	// bufSize of 0 creates an unbuffered channel
//	expand := MergeMapBuf(func(n int) Seq[int] {
//	    return From(n, n*2)
//	}, 0)
//	seq := From(1, 2, 3)
//	result := expand(seq)
//
//	// Producers and consumer are synchronized
//	for v := range result {
//	    fmt.Println(v)
//	}
//
// See Also:
//   - Chain: Sequential version (deterministic order)
//   - MergeAll: Merges pre-existing sequences concurrently
//   - Map: Transforms elements without flattening
//   - Async: Converts a single sequence to asynchronous
func MergeMapBuf[A, B any](f func(A) Seq[B], bufSize int) Operator[A, B] {
	return F.Flow2(
		Map(f),
		MergeAll[B](bufSize),
	)
}

// MergeMap applies a function that returns a sequence to each element and merges the results concurrently using a default buffer size.
// This is a convenience wrapper around MergeMapBuf that uses a default buffer size of 8.
// It's the concurrent version of Chain (flatMap), where each mapped sequence is processed in parallel.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of elements in the output sequences
//
// Parameters:
//   - f: A function that transforms each input element into a sequence of output elements
//
// Returns:
//   - Operator[A, B]: A function that takes a sequence of A and returns a flat sequence of B
//
// Behavior:
//   - Uses a default buffer size of 8 for the internal channel
//   - Applies f to each element in the input sequence to produce inner sequences
//   - Spawns one goroutine per inner sequence to produce elements concurrently
//   - Elements from different inner sequences are interleaved non-deterministically
//   - Properly handles early termination with goroutine cleanup
//   - Thread-safe: multiple producers can safely send to the shared channel
//
// Comparison with Chain:
//   - Chain: Sequential processing, deterministic order, no concurrency overhead
//   - MergeMap: Concurrent processing, non-deterministic order, better for I/O-bound tasks
//
// Example:
//
//	// Expand each number into a sequence concurrently
//	expand := MergeMap(func(n int) Seq[int] {
//	    return From(n, n*10, n*100)
//	})
//	seq := From(1, 2, 3)
//	result := expand(seq)
//
//	// Yields: 1, 10, 100, 2, 20, 200, 3, 30, 300 (order varies)
//	for v := range result {
//	    fmt.Println(v)
//	}
//
// See Also:
//   - MergeMapBuf: MergeMap with custom buffer size
//   - Chain: Sequential version (deterministic order)
//   - MergeAll: Merges pre-existing sequences concurrently
//   - Map: Transforms elements without flattening
func MergeMap[A, B any](f func(A) Seq[B]) Operator[A, B] {
	return MergeMapBuf(f, defaultBufferSize)
}
