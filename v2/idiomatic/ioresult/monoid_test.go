//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package ioresult

import (
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid(t *testing.T) {
	m := ApplicativeMonoid(S.Monoid)

	// good cases
	result1, err1 := m.Concat(Of("a"), Of("b"))()
	assert.NoError(t, err1)
	assert.Equal(t, "ab", result1)

	result2, err2 := m.Concat(Of("a"), m.Empty())()
	assert.NoError(t, err2)
	assert.Equal(t, "a", result2)

	result3, err3 := m.Concat(m.Empty(), Of("b"))()
	assert.NoError(t, err3)
	assert.Equal(t, "b", result3)

	// bad cases
	e1 := fmt.Errorf("e1")
	e2 := fmt.Errorf("e2")

	_, err4 := m.Concat(Left[string](e1), Of("b"))()
	assert.Error(t, err4)
	assert.Equal(t, e1, err4)

	_, err5 := m.Concat(Left[string](e1), Left[string](e2))()
	assert.Error(t, err5)
	assert.Equal(t, e1, err5)

	_, err6 := m.Concat(Of("a"), Left[string](e2))()
	assert.Error(t, err6)
	assert.Equal(t, e2, err6)
}

func TestApplicativeMonoidSeq(t *testing.T) {
	m := ApplicativeMonoidSeq(S.Monoid)

	t.Run("Sequential concatenation of successful values", func(t *testing.T) {
		result, err := m.Concat(Of("hello"), Of(" world"))()
		assert.NoError(t, err)
		assert.Equal(t, "hello world", result)
	})

	t.Run("Empty element is identity (left)", func(t *testing.T) {
		result, err := m.Concat(m.Empty(), Of("test"))()
		assert.NoError(t, err)
		assert.Equal(t, "test", result)
	})

	t.Run("Empty element is identity (right)", func(t *testing.T) {
		result, err := m.Concat(Of("test"), m.Empty())()
		assert.NoError(t, err)
		assert.Equal(t, "test", result)
	})

	t.Run("First error short-circuits", func(t *testing.T) {
		e1 := fmt.Errorf("error1")
		_, err := m.Concat(Left[string](e1), Of("world"))()
		assert.Error(t, err)
		assert.Equal(t, e1, err)
	})

	t.Run("Second error after first succeeds", func(t *testing.T) {
		e2 := fmt.Errorf("error2")
		_, err := m.Concat(Of("hello"), Left[string](e2))()
		assert.Error(t, err)
		assert.Equal(t, e2, err)
	})

	t.Run("Multiple concatenations", func(t *testing.T) {
		m := ApplicativeMonoidSeq(S.Monoid)
		result, err := m.Concat(
			m.Concat(Of("a"), Of("b")),
			m.Concat(Of("c"), Of("d")),
		)()
		assert.NoError(t, err)
		assert.Equal(t, "abcd", result)
	})
}

func TestApplicativeMonoidPar(t *testing.T) {
	m := ApplicativeMonoidPar(N.MonoidSum[int]())

	t.Run("Parallel concatenation of successful values", func(t *testing.T) {
		result, err := m.Concat(Of(10), Of(20))()
		assert.NoError(t, err)
		assert.Equal(t, 30, result)
	})

	t.Run("Empty element is identity (left)", func(t *testing.T) {
		result, err := m.Concat(m.Empty(), Of(42))()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("Empty element is identity (right)", func(t *testing.T) {
		result, err := m.Concat(Of(42), m.Empty())()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("Both empty returns empty", func(t *testing.T) {
		result, err := m.Empty()()
		assert.NoError(t, err)
		assert.Equal(t, 0, result) // 0 is the identity for sum
	})

	t.Run("First error in parallel", func(t *testing.T) {
		e1 := fmt.Errorf("error1")
		_, err := m.Concat(Left[int](e1), Of(20))()
		assert.Error(t, err)
		assert.Equal(t, e1, err)
	})

	t.Run("Second error in parallel", func(t *testing.T) {
		e2 := fmt.Errorf("error2")
		_, err := m.Concat(Of(10), Left[int](e2))()
		assert.Error(t, err)
		assert.Equal(t, e2, err)
	})

	t.Run("Associativity property", func(t *testing.T) {
		// (a <> b) <> c == a <> (b <> c)
		a := Of(1)
		b := Of(2)
		c := Of(3)

		left, leftErr := m.Concat(m.Concat(a, b), c)()
		right, rightErr := m.Concat(a, m.Concat(b, c))()

		assert.NoError(t, leftErr)
		assert.NoError(t, rightErr)
		assert.Equal(t, 6, left)
		assert.Equal(t, 6, right)
	})

	t.Run("Identity property (left)", func(t *testing.T) {
		// empty <> a == a
		a := Of(42)
		result, err := m.Concat(m.Empty(), a)()
		expected, expectedErr := a()

		assert.NoError(t, err)
		assert.NoError(t, expectedErr)
		assert.Equal(t, expected, result)
	})

	t.Run("Identity property (right)", func(t *testing.T) {
		// a <> empty == a
		a := Of(42)
		result, err := m.Concat(a, m.Empty())()
		expected, expectedErr := a()

		assert.NoError(t, err)
		assert.NoError(t, expectedErr)
		assert.Equal(t, expected, result)
	})
}

func TestMonoidWithProduct(t *testing.T) {
	m := ApplicativeMonoid(N.MonoidProduct[int]())

	t.Run("Multiply successful values", func(t *testing.T) {
		result, err := m.Concat(Of(3), Of(4))()
		assert.NoError(t, err)
		assert.Equal(t, 12, result)
	})

	t.Run("Identity is 1 for product", func(t *testing.T) {
		result, err := m.Empty()()
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("Multiply with identity", func(t *testing.T) {
		result, err := m.Concat(Of(5), m.Empty())()
		assert.NoError(t, err)
		assert.Equal(t, 5, result)
	})
}

func TestMonoidLaws(t *testing.T) {
	// Test that all three monoid variants satisfy monoid laws
	testMonoidLaws := func(name string, m Monoid[string]) {
		t.Run(name, func(t *testing.T) {
			a := Of("a")
			b := Of("b")
			c := Of("c")

			t.Run("Associativity: (a <> b) <> c == a <> (b <> c)", func(t *testing.T) {
				left, _ := m.Concat(m.Concat(a, b), c)()
				right, _ := m.Concat(a, m.Concat(b, c))()
				assert.Equal(t, left, right)
			})

			t.Run("Left identity: empty <> a == a", func(t *testing.T) {
				result, _ := m.Concat(m.Empty(), a)()
				expected, _ := a()
				assert.Equal(t, expected, result)
			})

			t.Run("Right identity: a <> empty == a", func(t *testing.T) {
				result, _ := m.Concat(a, m.Empty())()
				expected, _ := a()
				assert.Equal(t, expected, result)
			})
		})
	}

	testMonoidLaws("ApplicativeMonoid", ApplicativeMonoid(S.Monoid))
	testMonoidLaws("ApplicativeMonoidSeq", ApplicativeMonoidSeq(S.Monoid))
	testMonoidLaws("ApplicativeMonoidPar", ApplicativeMonoidPar(S.Monoid))
}
