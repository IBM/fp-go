// Copyright (c) 2025 IBM Corp.
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

package function_test

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
	"github.com/stretchr/testify/assert"
)

// TestTernary tests the behaviour previously covered by function.Ternary,
// now expressed via predicate.Fold (the non-deprecated equivalent).
func TestTernary(t *testing.T) {
	t.Run("applies onTrue when predicate is true", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		double := func(n int) int { return n * 2 }
		negate := func(n int) int { return -n }

		transform := P.Fold(negate, double)(isPositive)

		assert.Equal(t, 10, transform(5))
		assert.Equal(t, 20, transform(10))
	})

	t.Run("applies onFalse when predicate is false", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		double := func(n int) int { return n * 2 }
		negate := func(n int) int { return -n }

		transform := P.Fold(negate, double)(isPositive)

		assert.Equal(t, 3, transform(-3))
		assert.Equal(t, 5, transform(-5))
		assert.Equal(t, 0, transform(0))
	})

	t.Run("works with string classification", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		classify := P.Fold(
			F.Constant1[int]("non-positive"),
			F.Constant1[int]("positive"),
		)(isPositive)

		assert.Equal(t, "positive", classify(5))
		assert.Equal(t, "non-positive", classify(-3))
		assert.Equal(t, "non-positive", classify(0))
	})
}
