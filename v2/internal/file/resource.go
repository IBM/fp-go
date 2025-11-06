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

package file

import (
	F "github.com/IBM/fp-go/v2/function"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GA,
	GR,
	GANY,
	E, R, A, ANY any](
	mchain func(GR, func(R) GA) GA,
	mfold1 func(GA, func(E) GA, func(A) GA) GA,
	mfold2 func(GANY, func(E) GA, func(ANY) GA) GA,
	mmap func(GANY, func(ANY) A) GA,
	left func(E) GA,
) func(onCreate func() GR, onRelease func(R) GANY) func(func(R) GA) GA {

	return func(onCreate func() GR, onRelease func(R) GANY) func(func(R) GA) GA {

		return func(f func(R) GA) GA {
			return mchain(
				onCreate(), func(r R) GA {
					// handle errors
					return mfold1(
						f(r),
						func(e E) GA {
							// the original error
							err := left(e)
							// if resource processing produced and error, still release the resource but return the first error
							return mfold2(
								onRelease(r),
								F.Constant1[E](err),
								F.Constant1[ANY](err),
							)
						},
						func(a A) GA {
							// if resource processing succeeded, release the resource. If this fails return failure, else the original error
							return F.Pipe1(
								onRelease(r),
								F.Bind2nd(mmap, F.Constant1[ANY](a)),
							)
						})
				},
			)
		}
	}
}
