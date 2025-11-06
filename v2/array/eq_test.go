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

package array

import (
	"testing"

	E "github.com/IBM/fp-go/v2/eq"
	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {
	intEq := Eq(E.FromStrictEquals[int]())

	// Test equal arrays
	assert.True(t, intEq.Equals([]int{1, 2, 3}, []int{1, 2, 3}))

	// Test different lengths
	assert.False(t, intEq.Equals([]int{1, 2, 3}, []int{1, 2}))

	// Test different values
	assert.False(t, intEq.Equals([]int{1, 2, 3}, []int{1, 2, 4}))

	// Test empty arrays
	assert.True(t, intEq.Equals([]int{}, []int{}))

	// Test string arrays
	stringEq := Eq(E.FromStrictEquals[string]())
	assert.True(t, stringEq.Equals([]string{"a", "b"}, []string{"a", "b"}))
	assert.False(t, stringEq.Equals([]string{"a", "b"}, []string{"a", "c"}))
}
