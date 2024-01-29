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

package iooption

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/iooption/generic"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[A any](e EQ.Eq[A]) EQ.Eq[IOOption[A]] {
	return G.Eq[IOOption[A]](e)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[A comparable]() EQ.Eq[IOOption[A]] {
	return G.FromStrictEquals[IOOption[A]]()
}
