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

package io

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/io/generic"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[A any](e EQ.Eq[A]) EQ.Eq[IO[A]] {
	return G.Eq[IO[A]](e)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[A comparable]() EQ.Eq[IO[A]] {
	return G.FromStrictEquals[IO[A]]()
}
