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

package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
)

// Bind applies a function to an input state and merges the result into that state
func Bind[E, A, S1, S2 any](s func(A) func(S1) S2, f func(S1) IOEither[E, A]) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.Bind[IOEither[E, S1], IOEither[E, S2], IOEither[E, A], func(S1) IOEither[E, A]](s, f)
}

// BindTo initializes some state based on a value
func BindTo[
	E, A, S2 any](s func(A) S2) func(IOEither[E, A]) IOEither[E, S2] {
	return G.BindTo[IOEither[E, S2], IOEither[E, A]](s)
}

func ApS[
	E, A, S1, S2 any,
](s func(A) func(S1) S2, fa IOEither[E, A]) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.ApS[IOEither[E, S1], IOEither[E, S2], IOEither[E, A], IOEither[E, func(S1) S2]](s, fa)
}

func Let[E, A, S1, S2 any](
	s func(A) func(S1) S2,
	f func(S1) A,
) func(IOEither[E, S1]) IOEither[E, S2] {
	return G.Let[IOEither[E, S1], IOEither[E, S2]](s, f)
}
