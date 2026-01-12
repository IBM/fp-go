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

package stateio

import (
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
)

// uncurryState transforms a curried function into an uncurried function that operates on pairs.
// This is an internal helper function used by WithResource to adapt StateIO computations
// to work with the IO resource management functions.
//
// It converts: func(A) io.Kleisli[S, B] -> io.Kleisli[Pair[S, A], B]
func uncurryState[S, A, B any](f func(A) io.Kleisli[S, B]) io.Kleisli[Pair[S, A], B] {
	return func(r Pair[S, A]) IO[B] {
		return f(pair.Tail(r))(pair.Head(r))
	}
}

// WithResource provides safe resource management for StateIO computations.
// It ensures that resources are properly acquired and released, even if errors occur.
//
// The function takes:
//   - onCreate: A StateIO computation that creates/acquires the resource
//   - onRelease: A Kleisli arrow that releases the resource (receives the resource, returns any value)
//
// It returns a Kleisli arrow that takes a resource-using computation and ensures proper cleanup.
//
// The pattern follows the bracket pattern (acquire-use-release):
//  1. Acquire the resource using onCreate
//  2. Use the resource with the provided computation
//  3. Release the resource using onRelease (guaranteed to run)
//
// Example:
//
//	// Create a file resource
//	openFile := func(s AppState) IO[Pair[AppState, *os.File]] {
//	    return io.Of(pair.MakePair(s, file))
//	}
//
//	// Release the file resource
//	closeFile := func(f *os.File) StateIO[AppState, error] {
//	    return FromIO[AppState](io.Of(f.Close()))
//	}
//
//	// Use the resource safely
//	withFile := WithResource[string, AppState, *os.File, error](
//	    openFile,
//	    closeFile,
//	)
//
//	// Apply to a computation that uses the file
//	result := withFile(func(f *os.File) StateIO[AppState, string] {
//	    // Use file f here
//	    return Of[AppState]("data")
//	})
func WithResource[A, S, RES, ANY any](
	onCreate StateIO[S, RES],
	onRelease Kleisli[S, RES, ANY],
) Kleisli[S, Kleisli[S, RES, A], A] {
	release := uncurryState(onRelease)
	return func(f Kleisli[S, RES, A]) StateIO[S, A] {
		use := uncurryState(f)
		return func(s S) IO[Pair[S, A]] {
			return io.WithResource[Pair[S, RES], Pair[S, A]](onCreate(s), release)(use)
		}
	}
}
