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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/readerioeither"
)

func uncurryState[S, R, E, A, B any](f func(A) readerioeither.Kleisli[R, E, S, B]) readerioeither.Kleisli[R, E, Pair[S, A], B] {
	return func(r Pair[S, A]) ReaderIOEither[R, E, B] {
		return f(pair.Tail(r))(pair.Head(r))
	}
}

// WithResource constructs a function that creates a resource with state management, operates on it, and then releases the resource.
// This ensures proper resource cleanup even in the presence of errors, following the Resource Acquisition Is Initialization (RAII) pattern.
// The state is threaded through all operations: resource creation, usage, and release.
//
// The resource lifecycle with state management is:
//  1. onCreate: Acquires the resource (may modify state)
//  2. use: Operates on the resource with current state (provided as argument to the returned function)
//  3. onRelease: Releases the resource with current state (called regardless of success or failure)
//
// Type parameters:
//   - S: The state type that is threaded through all operations
//   - R: The reader/context type
//   - E: The error type
//   - RES: The resource type
//   - A: The type of the result produced by using the resource
//   - ANY: The type returned by the release function (typically ignored)
//
// Parameters:
//   - onCreate: A stateful computation that acquires the resource
//   - onRelease: A stateful function that releases the resource, called with the resource and current state, executed regardless of errors
//
// Returns:
//
//	A function that takes a resource-using function and returns a StateReaderIOEither that manages the resource lifecycle with state
//
// Example:
//
//	type AppState struct {
//	    openFiles int
//	}
//
//	withFile := WithResource(
//	    openFile("data.txt"), // Increments openFiles in state
//	    func(f *File) StateReaderIOEither[AppState, Config, error, int] {
//	        return closeFile(f) // Decrements openFiles in state
//	    },
//	)
//	result := withFile(func(f *File) StateReaderIOEither[AppState, Config, error, string] {
//	    return readContent(f)
//	})
func WithResource[A, S, R, E, RES, ANY any](
	onCreate StateReaderIOEither[S, R, E, RES],
	onRelease Kleisli[S, R, E, RES, ANY],
) Kleisli[S, R, E, Kleisli[S, R, E, RES, A], A] {
	release := uncurryState(onRelease)
	return func(f Kleisli[S, R, E, RES, A]) StateReaderIOEither[S, R, E, A] {
		use := uncurryState(f)
		return func(s S) ReaderIOEither[R, E, Pair[S, A]] {
			return readerioeither.WithResource[Pair[S, A]](onCreate(s), release)(use)
		}
	}
}
