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
	"github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/io"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[E, A any](eq EQ.Eq[Either[E, A]]) EQ.Eq[IOEither[E, A]] {
	return io.Eq(eq)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[E, A comparable]() EQ.Eq[IOEither[E, A]] {
	return Eq(either.FromStrictEquals[E, A]())
}
