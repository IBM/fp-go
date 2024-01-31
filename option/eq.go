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

package option

import (
	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
)

// Constructs an equal predicate for an `Option`
func Eq[A any](a EQ.Eq[A]) EQ.Eq[Option[A]] {
	// some convenient shortcuts
	fld := Fold(
		F.Constant(Fold(F.ConstTrue, F.Constant1[A](false))),
		F.Flow2(F.Curry2(a.Equals), F.Bind1st(Fold[A, bool], F.ConstFalse)),
	)
	// convert to an equals predicate
	return EQ.FromEquals(F.Uncurry2(fld))
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[A comparable]() EQ.Eq[Option[A]] {
	return Eq(EQ.FromStrictEquals[A]())
}
