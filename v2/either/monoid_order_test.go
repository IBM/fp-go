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

package either

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

func TestApplySemigroupOrder(t *testing.T) {
	sg := ApplySemigroup[error, string](orderMonoid)

	t.Run("combines values in argument order", func(t *testing.T) {
		assert.Equal(t, Right[error]("ab"), sg.Concat(Right[error]("a"), Right[error]("b")))
	})

	t.Run("reports the error of the first argument when both fail", func(t *testing.T) {
		assert.Equal(t, Left[string](errOrderA), sg.Concat(Left[string](errOrderA), Left[string](errOrderB)))
	})
}

func TestApplicativeMonoidOrder(t *testing.T) {
	m := ApplicativeMonoid[error](orderMonoid)

	t.Run("combines values in argument order", func(t *testing.T) {
		assert.Equal(t, Right[error]("ab"), m.Concat(Right[error]("a"), Right[error]("b")))
	})

	t.Run("empty lifts the empty value", func(t *testing.T) {
		assert.Equal(t, Right[error](""), m.Empty())
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := Right[error]("a"), Right[error]("b"), Right[error]("c")
		assert.Equal(t, Right[error]("abc"), m.Concat(m.Concat(a, b), c))
		assert.Equal(t, Right[error]("abc"), m.Concat(a, m.Concat(b, c)))
	})

	t.Run("propagates the first failure", func(t *testing.T) {
		assert.Equal(t, Left[string](errOrderA), m.Concat(Left[string](errOrderA), Right[error]("b")))
		assert.Equal(t, Left[string](errOrderB), m.Concat(Right[error]("a"), Left[string](errOrderB)))
		assert.Equal(t, Left[string](errOrderA), m.Concat(Left[string](errOrderA), Left[string](errOrderB)))
	})
}

func TestAlternativeMonoidOrder(t *testing.T) {
	m := AlternativeMonoid[error](orderMonoid)
	a, b := Right[error]("a"), Right[error]("b")
	failA, failB := Left[string](errOrderA), Left[string](errOrderB)

	t.Run("combines two successes in argument order", func(t *testing.T) {
		assert.Equal(t, Right[error]("ab"), m.Concat(a, b))
	})

	t.Run("falls back to the success", func(t *testing.T) {
		assert.Equal(t, b, m.Concat(failA, b))
		assert.Equal(t, a, m.Concat(a, failB))
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		assert.Equal(t, failB, m.Concat(failA, failB))
	})

	t.Run("empty lifts the empty value", func(t *testing.T) {
		assert.Equal(t, Right[error](""), m.Empty())
	})
}

func TestAltMonoidOrder(t *testing.T) {
	m := AltMonoid[error, string](F.Constant(Left[string](errOrderEmpty)))
	a, b := Right[error]("a"), Right[error]("b")
	failA, failB := Left[string](errOrderA), Left[string](errOrderB)

	t.Run("prefers the first success", func(t *testing.T) {
		assert.Equal(t, a, m.Concat(a, b))
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		assert.Equal(t, b, m.Concat(failA, b))
		assert.Equal(t, failB, m.Concat(failA, failB))
	})

	t.Run("empty is the zero value and an identity", func(t *testing.T) {
		assert.Equal(t, Left[string](errOrderEmpty), m.Empty())
		assert.Equal(t, a, m.Concat(m.Empty(), a))
		assert.Equal(t, a, m.Concat(a, m.Empty()))
	})
}

func TestAltSemigroupOrder(t *testing.T) {
	sg := AltSemigroup[error, string]()

	assert.Equal(t, Right[error]("a"), sg.Concat(Right[error]("a"), Right[error]("b")))
	assert.Equal(t, Right[error]("b"), sg.Concat(Left[string](errOrderA), Right[error]("b")))
	assert.Equal(t, Left[string](errOrderB), sg.Concat(Left[string](errOrderA), Left[string](errOrderB)))
}
