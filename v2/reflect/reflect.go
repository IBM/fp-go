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

package option

import (
	R "reflect"

	F "github.com/IBM/fp-go/v2/function"
	G "github.com/IBM/fp-go/v2/reflect/generic"
)

func ReduceWithIndex[A any](f func(int, A, R.Value) A, initial A) func(R.Value) A {
	return func(val R.Value) A {
		count := val.Len()
		current := initial
		for i := 0; i < count; i++ {
			current = f(i, current, val.Index(i))
		}
		return current
	}
}

func Reduce[A any](f func(A, R.Value) A, initial A) func(R.Value) A {
	return ReduceWithIndex(F.Ignore1of3[int](f), initial)
}

func Map[A any](f func(R.Value) A) func(R.Value) []A {
	return G.Map[[]A](f)
}
