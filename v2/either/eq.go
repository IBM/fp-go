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

package either

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
)

// Constructs an equal predicate for an `Either`
func Eq[E, A any](e EQ.Eq[E], a EQ.Eq[A]) EQ.Eq[Either[E, A]] {
	// some convenient shortcuts
	eqa := F.Curry2(a.Equals)
	eqe := F.Curry2(e.Equals)

	fca := F.Bind2nd(Fold[E, A, bool], F.Constant1[A](false))
	fce := F.Bind1st(Fold[E, A, bool], F.Constant1[E](false))

	fld := Fold(F.Flow2(eqe, fca), F.Flow2(eqa, fce))

	return EQ.FromEquals(F.Uncurry2(fld))
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[E, A comparable]() EQ.Eq[Either[E, A]] {
	return Eq(EQ.FromStrictEquals[E](), EQ.FromStrictEquals[A]())
}
