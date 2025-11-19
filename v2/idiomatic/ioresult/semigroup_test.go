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

package ioresult

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAltSemigroup(t *testing.T) {
	sg := AltSemigroup[int]()

	t.Run("First succeeds, second not evaluated", func(t *testing.T) {
		first := Of(42)
		second := Of(100)

		result, err := sg.Concat(first, second)()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("First fails, second succeeds", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Of(100)

		result, err := sg.Concat(first, second)()
		assert.NoError(t, err)
		assert.Equal(t, 100, result)
	})

	t.Run("Both fail, returns second error", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Left[int](errors.New("second error"))

		_, err := sg.Concat(first, second)()
		assert.Error(t, err)
		assert.Equal(t, "second error", err.Error())
	})

	t.Run("Both succeed, returns first", func(t *testing.T) {
		first := Of(42)
		second := Of(100)

		result, err := sg.Concat(first, second)()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("Associativity property", func(t *testing.T) {
		// (a <> b) <> c == a <> (b <> c)
		a := Left[int](errors.New("a"))
		b := Left[int](errors.New("b"))
		c := Of(42)

		// Left associative: (a <> b) <> c
		left, leftErr := sg.Concat(sg.Concat(a, b), c)()

		// Right associative: a <> (b <> c)
		right, rightErr := sg.Concat(a, sg.Concat(b, c))()

		assert.NoError(t, leftErr)
		assert.NoError(t, rightErr)
		assert.Equal(t, left, right)
		assert.Equal(t, 42, left)
	})

	t.Run("Multiple alternatives", func(t *testing.T) {
		sg := AltSemigroup[string]()

		first := Left[string](errors.New("error1"))
		second := Left[string](errors.New("error2"))
		third := Left[string](errors.New("error3"))
		fourth := Of("success")

		result, err := sg.Concat(
			sg.Concat(sg.Concat(first, second), third),
			fourth,
		)()

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})
}

func TestAltSemigroupWithStrings(t *testing.T) {
	sg := AltSemigroup[string]()

	t.Run("Concatenate successful string results", func(t *testing.T) {
		first := Of("hello")
		second := Of("world")

		result, err := sg.Concat(first, second)()
		assert.NoError(t, err)
		assert.Equal(t, "hello", result) // First one wins
	})

	t.Run("Fallback to second on first failure", func(t *testing.T) {
		first := Left[string](errors.New("network error"))
		second := Of("fallback value")

		result, err := sg.Concat(first, second)()
		assert.NoError(t, err)
		assert.Equal(t, "fallback value", result)
	})
}
