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

package readerioeither

import "github.com/IBM/fp-go/v2/ioeither"

// WithResource constructs a function that creates a resource, operates on it, and then releases the resource.
// This ensures proper resource cleanup even in the presence of errors, following the Resource Acquisition Is Initialization (RAII) pattern.
//
// The resource lifecycle is:
//  1. onCreate: Acquires the resource
//  2. use: Operates on the resource (provided as argument to the returned function)
//  3. onRelease: Releases the resource (called regardless of success or failure)
//
// Type parameters:
//   - A: The type of the result produced by using the resource
//   - L: The context type
//   - E: The error type
//   - R: The resource type
//   - ANY: The type returned by the release function (typically ignored)
//
// Parameters:
//   - onCreate: A computation that acquires the resource
//   - onRelease: A function that releases the resource, called with the resource and executed regardless of errors
//
// Returns:
//
//	A function that takes a resource-using function and returns a ReaderIOEither that manages the resource lifecycle
//
// Example:
//
//	withFile := WithResource(
//	    openFile("data.txt"),
//	    func(f *File) ReaderIOEither[Config, error, int] {
//	        return closeFile(f)
//	    },
//	)
//	result := withFile(func(f *File) ReaderIOEither[Config, error, string] {
//	    return readContent(f)
//	})
func WithResource[A, L, E, R, ANY any](onCreate ReaderIOEither[L, E, R], onRelease Kleisli[L, E, R, ANY]) Kleisli[L, E, Kleisli[L, E, R, A], A] {
	return func(f func(R) ReaderIOEither[L, E, A]) ReaderIOEither[L, E, A] {
		return func(l L) ioeither.IOEither[E, A] {
			// dispatch to the generic implementation
			return ioeither.WithResource[A](
				onCreate(l),
				func(r R) ioeither.IOEither[E, ANY] {
					return onRelease(r)(l)
				},
			)(func(r R) ioeither.IOEither[E, A] {
				return f(r)(l)
			})
		}
	}
}
