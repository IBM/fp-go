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

package readerresult

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid to pin the argument
// order of the curried Map, Ap and Alt operations used by the monoids.

var (
	orderMonoid   = M.MakeMonoid(func(a, b string) string { return a + b }, "")
	errOrderA     = errors.New("a")
	errOrderB     = errors.New("b")
	errOrderEmpty = errors.New("empty")
)

// orderTagged returns a reader result that appends tag to its environment
func orderTagged(tag string) ReaderResult[string, string] {
	return func(env string) (string, error) {
		return env + tag, nil
	}
}

func TestApplicativeMonoidOrder(t *testing.T) {
	m := ApplicativeMonoid[string](orderMonoid)

	t.Run("combines the results in argument order", func(t *testing.T) {
		v, err := m.Concat(orderTagged(":a"), orderTagged(":b"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:ax:b", v)
	})

	t.Run("empty yields the empty value", func(t *testing.T) {
		v, err := m.Empty()("x")
		assert.NoError(t, err)
		assert.Equal(t, "", v)
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := orderTagged(":a"), orderTagged(":b"), orderTagged(":c")
		left, _ := m.Concat(m.Concat(a, b), c)("x")
		right, _ := m.Concat(a, m.Concat(b, c))("x")
		assert.Equal(t, "x:ax:bx:c", left)
		assert.Equal(t, left, right)
	})

	t.Run("reports the error of the first argument when both fail", func(t *testing.T) {
		_, err := m.Concat(Left[string, string](errOrderA), Left[string, string](errOrderB))("x")
		assert.ErrorIs(t, err, errOrderA)
	})
}

func TestAlternativeMonoidOrder(t *testing.T) {
	m := AlternativeMonoid[string](orderMonoid)

	t.Run("combines two successes in argument order", func(t *testing.T) {
		v, err := m.Concat(orderTagged(":a"), orderTagged(":b"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:ax:b", v)
	})

	t.Run("falls back to the success", func(t *testing.T) {
		v, err := m.Concat(Left[string, string](errOrderA), orderTagged(":b"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:b", v)

		v, err = m.Concat(orderTagged(":a"), Left[string, string](errOrderB))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:a", v)
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		_, err := m.Concat(Left[string, string](errOrderA), Left[string, string](errOrderB))("x")
		assert.ErrorIs(t, err, errOrderB)
	})
}

func TestAltMonoidOrder(t *testing.T) {
	m := AltMonoid[string, string](F.Constant(Left[string, string](errOrderEmpty)))

	t.Run("prefers the first success", func(t *testing.T) {
		v, err := m.Concat(orderTagged(":a"), orderTagged(":b"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:a", v)
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		v, err := m.Concat(Left[string, string](errOrderA), orderTagged(":b"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:b", v)
	})

	t.Run("empty is the zero value and an identity", func(t *testing.T) {
		_, err := m.Empty()("x")
		assert.ErrorIs(t, err, errOrderEmpty)

		v, err := m.Concat(m.Empty(), orderTagged(":a"))("x")
		assert.NoError(t, err)
		assert.Equal(t, "x:a", v)
	})
}
