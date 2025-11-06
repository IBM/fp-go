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

package semigroup

import (
	M "github.com/IBM/fp-go/v2/magma"
)

func GenericMonadConcatAll[GA ~[]A, A any](s Semigroup[A]) func(GA, A) A {
	return M.GenericMonadConcatAll[GA](M.MakeMagma(s.Concat))
}

func GenericConcatAll[GA ~[]A, A any](s Semigroup[A]) func(A) func(GA) A {
	return M.GenericConcatAll[GA](M.MakeMagma(s.Concat))
}

func MonadConcatAll[A any](s Semigroup[A]) func([]A, A) A {
	return GenericMonadConcatAll[[]A](s)
}

func ConcatAll[A any](s Semigroup[A]) func(A) func([]A) A {
	return GenericConcatAll[[]A](s)
}
