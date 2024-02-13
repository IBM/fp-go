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

package state

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/state/generic"
)

// Constructs an equal predicate for a [State]
func Eq[S, A any](w EQ.Eq[S], a EQ.Eq[A]) func(S) EQ.Eq[State[S, A]] {
	return G.Eq[State[S, A]](w, a)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[S, A comparable]() func(S) EQ.Eq[State[S, A]] {
	return G.FromStrictEquals[State[S, A]]()
}
