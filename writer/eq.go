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

package writer

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/writer/generic"
)

// Constructs an equal predicate for a [Writer]
func Eq[W, A any](w EQ.Eq[W], a EQ.Eq[A]) EQ.Eq[Writer[W, A]] {
	return G.Eq[Writer[W, A]](w, a)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[W, A comparable]() EQ.Eq[Writer[W, A]] {
	return G.FromStrictEquals[Writer[W, A]]()
}
