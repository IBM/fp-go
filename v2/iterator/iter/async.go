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
	N "github.com/IBM/fp-go/v2/number"
)

// Async converts a synchronous sequence into an asynchronous buffered sequence.
// It spawns a goroutine to consume the input sequence and sends values through
// a buffered channel, allowing concurrent production and consumption of elements.
//
// The function provides backpressure control through the buffer size and properly
// handles early termination when the consumer stops iterating. This is useful for
// decoupling producers and consumers, enabling pipeline parallelism, or when you
// need to process sequences concurrently.
//
// # Type Parameters
//
//   - T: The type of elements in the sequence
//
// # Parameters
//
//   - input: The source sequence to be consumed asynchronously
//   - bufSize: The buffer size for the channel. Negative values are treated as 0 (unbuffered).
//     A larger buffer allows more elements to be produced ahead of consumption,
//     but uses more memory. A buffer of 0 creates an unbuffered channel requiring
//     synchronization between producer and consumer.
//
// # Returns
//
//   - Seq[T]: A new sequence that yields elements from the input sequence asynchronously
//
// # Behavior
//
//   - Spawns a goroutine that consumes the input sequence
//   - Elements are sent through a buffered channel to the output sequence
//   - Properly handles early termination: if the consumer stops iterating (yield returns false),
//     the producer goroutine is signaled to stop via a done channel
//   - Both the producer goroutine and the done channel are properly cleaned up
//   - The channel is closed when the input sequence is exhausted or early termination occurs
//
// # Example Usage
//
//	// Create an async sequence with a buffer of 10
//	seq := From(1, 2, 3, 4, 5)
//	async := Async(seq, 10)
//
//	// Elements are produced concurrently
//	for v := range async {
//	    fmt.Println(v) // Prints: 1, 2, 3, 4, 5
//	}
//
// # Example with Early Termination
//
//	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	async := Async(seq, 5)
//
//	// Stop after 3 elements - producer goroutine will be properly cleaned up
//	count := 0
//	for v := range async {
//	    fmt.Println(v)
//	    count++
//	    if count >= 3 {
//	        break
//	    }
//	}
//
// # Example with Unbuffered Channel
//
//	// bufSize of 0 creates an unbuffered channel
//	seq := From(1, 2, 3)
//	async := Async(seq, 0)
//
//	// Producer and consumer are synchronized
//	for v := range async {
//	    fmt.Println(v)
//	}
//
// # See Also
//
//   - From: Creates a sequence from values
//   - Map: Transforms sequence elements
//   - Filter: Filters sequence elements
func Async[T any](input Seq[T], bufSize int) Seq[T] {
	return func(yield func(T) bool) {
		ch := make(chan T, N.Max(bufSize, 0))
		done := make(chan Void)

		go func() {
			defer close(ch)
			for v := range input {
				select {
				case ch <- v:
				case <-done:
					return
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

// Async2 converts a synchronous key-value sequence into an asynchronous buffered sequence.
// It spawns a goroutine to consume the input sequence and sends key-value pairs through
// a buffered channel, allowing concurrent production and consumption of elements.
//
// This function is the Seq2 variant of Async, providing the same asynchronous behavior
// for key-value sequences. It internally converts the Seq2 to a sequence of Pairs,
// applies Async, and converts back to Seq2.
//
// # Type Parameters
//
//   - K: The type of keys in the sequence
//   - V: The type of values in the sequence
//
// # Parameters
//
//   - input: The source key-value sequence to be consumed asynchronously
//   - bufSize: The buffer size for the channel. Negative values are treated as 0 (unbuffered).
//     A larger buffer allows more elements to be produced ahead of consumption,
//     but uses more memory. A buffer of 0 creates an unbuffered channel requiring
//     synchronization between producer and consumer.
//
// # Returns
//
//   - Seq2[K, V]: A new key-value sequence that yields elements from the input sequence asynchronously
//
// # Behavior
//
//   - Spawns a goroutine that consumes the input key-value sequence
//   - Key-value pairs are sent through a buffered channel to the output sequence
//   - Properly handles early termination: if the consumer stops iterating (yield returns false),
//     the producer goroutine is signaled to stop via a done channel
//   - Both the producer goroutine and the done channel are properly cleaned up
//   - The channel is closed when the input sequence is exhausted or early termination occurs
//
// # Example Usage
//
//	// Create an async key-value sequence with a buffer of 10
//	seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
//	async := Async2(seq, 10)
//
//	// Elements are produced concurrently
//	for k, v := range async {
//	    fmt.Printf("%d: %s\n", k, v)
//	}
//	// Output:
//	// 1: a
//	// 2: b
//	// 3: c
//
// # Example with Early Termination
//
//	seq := MonadZip(From(1, 2, 3, 4, 5), From("a", "b", "c", "d", "e"))
//	async := Async2(seq, 5)
//
//	// Stop after 2 pairs - producer goroutine will be properly cleaned up
//	count := 0
//	for k, v := range async {
//	    fmt.Printf("%d: %s\n", k, v)
//	    count++
//	    if count >= 2 {
//	        break
//	    }
//	}
//
// # See Also
//
//   - Async: Asynchronous sequence for single-value sequences
//   - ToSeqPair: Converts Seq2 to Seq of Pairs
//   - FromSeqPair: Converts Seq of Pairs to Seq2
//   - MonadZip: Creates key-value sequences from two sequences
func Async2[K, V any](input Seq2[K, V], bufSize int) Seq2[K, V] {
	return FromSeqPair(Async(ToSeqPair(input), bufSize))
}
