// Copyright (c) 2023 - 2025 IBM Corp.
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
	"sync/atomic"
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid to pin the argument
// order of the curried Map and Ap operations used by the monoids. The default Ap may
// run effects in parallel, so only the values are checked, not the effect order.

var orderMonoid = M.MakeMonoid(func(a, b string) string { return a + b }, "")

func TestApplySemigroupOrder(t *testing.T) {
	sg := ApplySemigroup[string](orderMonoid)

	assert.Equal(t, "ab", sg.Concat(Of("a"), Of("b"))())
}

func TestApplicativeMonoidOrder(t *testing.T) {
	m := ApplicativeMonoid(orderMonoid)

	t.Run("combines values in argument order", func(t *testing.T) {
		assert.Equal(t, "ab", m.Concat(Of("a"), Of("b"))())
	})

	t.Run("empty lifts the empty value", func(t *testing.T) {
		assert.Equal(t, "", m.Empty()())
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := Of("a"), Of("b"), Of("c")
		assert.Equal(t, "abc", m.Concat(m.Concat(a, b), c)())
		assert.Equal(t, "abc", m.Concat(a, m.Concat(b, c))())
	})

	t.Run("runs every effect exactly once", func(t *testing.T) {
		var calls atomic.Int32
		count := func(s string) IO[string] {
			return func() string {
				calls.Add(1)
				return s
			}
		}
		assert.Equal(t, "abc", m.Concat(count("a"), m.Concat(count("b"), count("c")))())
		assert.Equal(t, int32(3), calls.Load())
	})
}
