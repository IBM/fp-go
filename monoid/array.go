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

package monoid

import (
	S "github.com/IBM/fp-go/semigroup"
)

func GenericConcatAll[GA ~[]A, A any](m Monoid[A]) func(GA) A {
	return S.GenericConcatAll[GA](S.MakeSemigroup(m.Concat))(m.Empty())
}

// ConcatAll concatenates all values using the monoid and the default empty value
func ConcatAll[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}

// Fold concatenates all values using the monoid and the default empty value
func Fold[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}
