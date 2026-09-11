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

package readeroption

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid to pin the argument
// order of the curried Map, Ap and Alt operations used by the monoids.

var orderMonoid = M.MakeMonoid(func(a, b string) string { return a + b }, "")

// orderTagged returns a reader option that appends tag to its environment
func orderTagged(tag string) ReaderOption[string, string] {
	return func(env string) option.Option[string] {
		return option.Some(env + tag)
	}
}

func TestApplicativeMonoidOrder(t *testing.T) {
	m := ApplicativeMonoid[string](orderMonoid)

	t.Run("combines the results in argument order", func(t *testing.T) {
		assert.Equal(t, option.Some("x:ax:b"), m.Concat(orderTagged(":a"), orderTagged(":b"))("x"))
	})

	t.Run("is None when either argument is None", func(t *testing.T) {
		assert.Equal(t, option.None[string](), m.Concat(None[string, string](), orderTagged(":b"))("x"))
		assert.Equal(t, option.None[string](), m.Concat(orderTagged(":a"), None[string, string]())("x"))
	})

	t.Run("empty yields the empty value", func(t *testing.T) {
		assert.Equal(t, option.Some(""), m.Empty()("x"))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := orderTagged(":a"), orderTagged(":b"), orderTagged(":c")
		assert.Equal(t, option.Some("x:ax:bx:c"), m.Concat(m.Concat(a, b), c)("x"))
		assert.Equal(t, option.Some("x:ax:bx:c"), m.Concat(a, m.Concat(b, c))("x"))
	})
}

func TestAlternativeMonoidOrder(t *testing.T) {
	m := AlternativeMonoid[string](orderMonoid)

	t.Run("combines two values in argument order", func(t *testing.T) {
		assert.Equal(t, option.Some("x:ax:b"), m.Concat(orderTagged(":a"), orderTagged(":b"))("x"))
	})

	t.Run("falls back to the present value", func(t *testing.T) {
		assert.Equal(t, option.Some("x:b"), m.Concat(None[string, string](), orderTagged(":b"))("x"))
		assert.Equal(t, option.Some("x:a"), m.Concat(orderTagged(":a"), None[string, string]())("x"))
	})

	t.Run("is None when both are None", func(t *testing.T) {
		assert.Equal(t, option.None[string](), m.Concat(None[string, string](), None[string, string]())("x"))
	})
}

func TestAltMonoidOrder(t *testing.T) {
	m := AltMonoid[string, string]()

	t.Run("prefers the first value", func(t *testing.T) {
		assert.Equal(t, option.Some("x:a"), m.Concat(orderTagged(":a"), orderTagged(":b"))("x"))
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		assert.Equal(t, option.Some("x:b"), m.Concat(None[string, string](), orderTagged(":b"))("x"))
	})

	t.Run("empty is None and an identity", func(t *testing.T) {
		assert.Equal(t, option.None[string](), m.Empty()("x"))
		assert.Equal(t, option.Some("x:a"), m.Concat(m.Empty(), orderTagged(":a"))("x"))
		assert.Equal(t, option.Some("x:a"), m.Concat(orderTagged(":a"), m.Empty())("x"))
	})
}
