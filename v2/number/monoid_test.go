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

package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test monoid laws for MonoidSum
func TestMonoidSumLaws(t *testing.T) {
	sumMonoid := MonoidSum[int]()
	testData := []int{0, 1, 1000, -1, -100, 42}

	// Test identity laws manually
	for _, x := range testData {
		assert.Equal(t, x, sumMonoid.Concat(sumMonoid.Empty(), x),
			"Left identity failed for %d", x)
		assert.Equal(t, x, sumMonoid.Concat(x, sumMonoid.Empty()),
			"Right identity failed for %d", x)
	}

	// Test associativity
	assert.Equal(t,
		sumMonoid.Concat(sumMonoid.Concat(1, 2), 3),
		sumMonoid.Concat(1, sumMonoid.Concat(2, 3)),
	)
}

// Test monoid laws for MonoidProduct
func TestMonoidProductLaws(t *testing.T) {
	prodMonoid := MonoidProduct[int]()
	testData := []int{1, 2, 10, -1, -5, 42}

	// Test identity laws
	for _, x := range testData {
		assert.Equal(t, x, prodMonoid.Concat(prodMonoid.Empty(), x),
			"Left identity failed for %d", x)
		assert.Equal(t, x, prodMonoid.Concat(x, prodMonoid.Empty()),
			"Right identity failed for %d", x)
	}

	// Test associativity
	assert.Equal(t,
		prodMonoid.Concat(prodMonoid.Concat(2, 3), 4),
		prodMonoid.Concat(2, prodMonoid.Concat(3, 4)),
	)
}
