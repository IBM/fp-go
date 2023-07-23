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

package either

import (
	F "github.com/IBM/fp-go/function"
)

// constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[E, R, A any](onCreate func() Either[E, R], onRelease func(R) Either[E, any]) func(func(R) Either[E, A]) Either[E, A] {

	return func(f func(R) Either[E, A]) Either[E, A] {
		return MonadChain(
			onCreate(), func(r R) Either[E, A] {
				// run the code and make sure to release as quickly as possible
				res := f(r)
				released := onRelease(r)
				// handle the errors
				return MonadFold(
					res,
					Left[A, E],
					func(a A) Either[E, A] {
						return F.Pipe1(
							released,
							MapTo[E, any](a),
						)
					})
			},
		)
	}
}
