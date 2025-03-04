// Copyright (c) 2023 IBM Corp.
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

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[A, L, E, R, ANY any](onCreate ReaderIOEither[L, E, R], onRelease func(R) ReaderIOEither[L, E, ANY]) func(func(R) ReaderIOEither[L, E, A]) ReaderIOEither[L, E, A] {
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
