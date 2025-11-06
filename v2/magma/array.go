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

package magma

import (
	F "github.com/IBM/fp-go/v2/function"
	AR "github.com/IBM/fp-go/v2/internal/array"
)

func GenericMonadConcatAll[GA ~[]A, A any](m Magma[A]) func(GA, A) A {
	return func(as GA, first A) A {
		return AR.Reduce(as, m.Concat, first)
	}
}

// GenericConcatAll concats all items using the semigroup and a starting value
func GenericConcatAll[GA ~[]A, A any](m Magma[A]) func(A) func(GA) A {
	ca := GenericMonadConcatAll[GA](m)
	return func(a A) func(GA) A {
		return F.Bind2nd(ca, a)
	}
}

func MonadConcatAll[A any](m Magma[A]) func([]A, A) A {
	return GenericMonadConcatAll[[]A](m)
}

// ConcatAll concats all items using the semigroup and a starting value
func ConcatAll[A any](m Magma[A]) func(A) func([]A) A {
	return GenericConcatAll[[]A](m)
}
