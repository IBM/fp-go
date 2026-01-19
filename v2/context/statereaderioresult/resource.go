// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioresult

import "github.com/IBM/fp-go/v2/statereaderioeither"

// WithResource constructs a function that creates a resource with state management, operates on it,
// and then releases the resource. This ensures proper resource cleanup even in the presence of errors,
// following the Resource Acquisition Is Initialization (RAII) pattern.
//
// The state is threaded through all operations: resource creation, usage, and release.
//
// The resource lifecycle with state management is:
//  1. onCreate: Acquires the resource (may modify state)
//  2. use: Operates on the resource with current state (provided as argument to the returned function)
//  3. onRelease: Releases the resource with current state (called regardless of success or failure)
//
// Type parameters:
//   - A: The type of the result produced by using the resource
//   - S: The state type that is threaded through all operations
//   - RES: The resource type
//   - ANY: The type returned by the release function (typically ignored)
//
// Parameters:
//   - onCreate: A stateful computation that acquires the resource
//   - onRelease: A stateful function that releases the resource, called with the resource and current state,
//     executed regardless of errors
//
// Returns:
//
//	A function that takes a resource-using function and returns a StateReaderIOResult that manages
//	the resource lifecycle with state
//
// Example:
//
//	type AppState struct {
//	    openFiles int
//	}
//
//	// Resource creation that updates state
//	openFile := func(filename string) StateReaderIOResult[AppState, *File] {
//	    return func(state AppState) ReaderIOResult[Pair[AppState, *File]] {
//	        return func(ctx context.Context) IOResult[Pair[AppState, *File]] {
//	            return func() Result[Pair[AppState, *File]] {
//	                file, err := os.Open(filename)
//	                if err != nil {
//	                    return result.Error[Pair[AppState, *File]](err)
//	                }
//	                newState := AppState{openFiles: state.openFiles + 1}
//	                return result.Of(pair.MakePair(newState, file))
//	            }
//	        }
//	    }
//	}
//
//	// Resource release that updates state
//	closeFile := func(f *File) StateReaderIOResult[AppState, int] {
//	    return func(state AppState) ReaderIOResult[Pair[AppState, int]] {
//	        return func(ctx context.Context) IOResult[Pair[AppState, int]] {
//	            return func() Result[Pair[AppState, int]] {
//	                f.Close()
//	                newState := AppState{openFiles: state.openFiles - 1}
//	                return result.Of(pair.MakePair(newState, 0))
//	            }
//	        }
//	    }
//	}
//
//	// Use the resource with automatic cleanup
//	withFile := WithResource(
//	    openFile("data.txt"),
//	    closeFile,
//	)
//
//	result := withFile(func(f *File) StateReaderIOResult[AppState, string] {
//	    return readContent(f) // File will be closed automatically
//	})
//
//	// Execute the computation
//	initialState := AppState{openFiles: 0}
//	ctx := t.Context()
//	outcome := result(initialState)(ctx)()
func WithResource[A, S, RES, ANY any](
	onCreate StateReaderIOResult[S, RES],
	onRelease Kleisli[S, RES, ANY],
) Kleisli[S, Kleisli[S, RES, A], A] {
	return statereaderioeither.WithResource[A](onCreate, onRelease)
}
