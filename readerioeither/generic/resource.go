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

package generic

import (
	ET "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GEA ~func(L) TEA,
	GER ~func(L) TER,
	GEANY ~func(L) TEANY,

	TEA ~func() ET.Either[E, A],
	TER ~func() ET.Either[E, R],
	TEANY ~func() ET.Either[E, ANY],

	L, E, R, A, ANY any](onCreate GER, onRelease func(R) GEANY) func(func(R) GEA) GEA {

	return func(f func(R) GEA) GEA {
		return func(l L) TEA {
			// dispatch to the generic implementation
			return IOE.WithResource[TEA](
				onCreate(l),
				func(r R) TEANY {
					return onRelease(r)(l)
				},
			)(func(r R) TEA {
				return f(r)(l)
			})
		}
	}
}
