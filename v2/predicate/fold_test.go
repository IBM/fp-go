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

package predicate

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestFold_TrueBranch(t *testing.T) {
	t.Run("calls onTrue when predicate is satisfied", func(t *testing.T) {
		result := F.Pipe1(isPositive, Fold(
			func(n int) string { return "negative or zero" },
			func(n int) string { return "positive" },
		))(5)
		assert.Equal(t, "positive", result)
	})

	t.Run("onTrue receives the original input value", func(t *testing.T) {
		result := F.Pipe1(isPositive, Fold(
			func(n int) int { return -n },
			func(n int) int { return n * 10 },
		))(3)
		assert.Equal(t, 30, result)
	})
}

func TestFold_FalseBranch(t *testing.T) {
	t.Run("calls onFalse when predicate is not satisfied", func(t *testing.T) {
		result := F.Pipe1(isPositive, Fold(
			func(n int) string { return "negative or zero" },
			func(n int) string { return "positive" },
		))(-1)
		assert.Equal(t, "negative or zero", result)
	})

	t.Run("onFalse receives the original input value", func(t *testing.T) {
		result := F.Pipe1(isPositive, Fold(
			func(n int) int { return n - 1 },
			func(n int) int { return n * 10 },
		))(0)
		assert.Equal(t, -1, result)
	})
}

func TestFold_EdgeCases(t *testing.T) {
	t.Run("zero is routed to onFalse by isPositive", func(t *testing.T) {
		result := F.Pipe1(isPositive, Fold(
			func(n int) string { return "not positive" },
			func(n int) string { return "positive" },
		))(0)
		assert.Equal(t, "not positive", result)
	})

	t.Run("works with Always predicate — always takes onTrue branch", func(t *testing.T) {
		for i := range 5 {
			result := F.Pipe1(Always[int](), Fold(
				func(n int) string { return "false" },
				func(n int) string { return "true" },
			))(i)
			assert.Equal(t, "true", result)
		}
	})

	t.Run("works with Never predicate — always takes onFalse branch", func(t *testing.T) {
		for i := range 5 {
			result := F.Pipe1(Never[int](), Fold(
				func(n int) string { return "false" },
				func(n int) string { return "true" },
			))(i)
			assert.Equal(t, "false", result)
		}
	})
}

func TestFold_TypeVariety(t *testing.T) {
	t.Run("maps int predicate to bool output", func(t *testing.T) {
		classify := F.Pipe1(isEven, Fold(
			func(_ int) bool { return false },
			func(_ int) bool { return true },
		))
		assert.True(t, classify(4))
		assert.False(t, classify(3))
	})

	t.Run("maps string predicate to int output", func(t *testing.T) {
		isEmpty := func(s string) bool { return s == "" }
		length := F.Pipe1(isEmpty, Fold(
			func(s string) int { return len(s) },
			func(_ string) int { return 0 },
		))
		assert.Equal(t, 0, length(""))
		assert.Equal(t, 5, length("hello"))
	})
}

func TestFold_Integration(t *testing.T) {
	t.Run("composes with And to classify numbers", func(t *testing.T) {
		import_and := F.Pipe1(isPositive, And(isEven))
		classify := F.Pipe1(import_and, Fold(
			func(n int) string { return "other" },
			func(n int) string { return "positive even" },
		))
		assert.Equal(t, "positive even", classify(4))
		assert.Equal(t, "other", classify(3))
		assert.Equal(t, "other", classify(-2))
	})

	t.Run("composes with Not to invert branch selection", func(t *testing.T) {
		result := F.Pipe1(Not(isPositive), Fold(
			func(n int) string { return "positive" },
			func(n int) string { return "not positive" },
		))(5)
		assert.Equal(t, "positive", result)
	})
}
