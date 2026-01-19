// Copyright (c) 2024 - 2025 IBM Corp.
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

package testing

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	EQ "github.com/IBM/fp-go/v2/eq"
	"github.com/stretchr/testify/assert"
)

func TestMonadLaws(t *testing.T) {
	// some comparison
	eqs := A.Eq(EQ.FromStrictEquals[string]())
	eqa := EQ.FromStrictEquals[bool]()
	eqb := EQ.FromStrictEquals[int]()
	eqc := EQ.FromStrictEquals[string]()

	ab := func(a bool) int {
		if a {
			return 1
		}
		return 0
	}

	bc := func(b int) string {
		return fmt.Sprintf("value %d", b)
	}

	laws := AssertLaws(t, eqs, eqa, eqb, eqc, ab, bc, A.Empty[string](), t.Context())

	assert.True(t, laws(true))
	assert.True(t, laws(false))
}
