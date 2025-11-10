// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package testing

import (
	"testing"

	"github.com/IBM/fp-go/v2/eq"
	IOET "github.com/IBM/fp-go/v2/ioeither/testing"
)

// AssertLaws asserts the apply monad laws for the `IOEither` monad
func AssertLaws[A, B, C any](t *testing.T,
	eqa eq.Eq[A],
	eqb eq.Eq[B],
	eqc eq.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return IOET.AssertLaws(t, eq.FromStrictEquals[error](), eqa, eqb, eqc, ab, bc)
}
