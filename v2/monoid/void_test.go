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

package monoid

import (
	"testing"

	"github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// TestVoidMonoid_Basic tests basic VoidMonoid functionality
func TestVoidMonoid_Basic(t *testing.T) {
	m := VoidMonoid()

	// Test Empty returns VOID
	empty := m.Empty()
	assert.Equal(t, function.VOID, empty)

	// Test Concat returns VOID (since all Void values are identical)
	result := m.Concat(function.VOID, function.VOID)
	assert.Equal(t, function.VOID, result)
}

// TestVoidMonoid_Laws verifies VoidMonoid satisfies monoid laws
func TestVoidMonoid_Laws(t *testing.T) {
	m := VoidMonoid()

	// Since Void has only one value, we test with that value
	v := function.VOID

	// Left Identity: Concat(Empty(), x) = x
	t.Run("left identity", func(t *testing.T) {
		result := m.Concat(m.Empty(), v)
		assert.Equal(t, v, result, "Left identity law failed")
	})

	// Right Identity: Concat(x, Empty()) = x
	t.Run("right identity", func(t *testing.T) {
		result := m.Concat(v, m.Empty())
		assert.Equal(t, v, result, "Right identity law failed")
	})

	// Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
	t.Run("associativity", func(t *testing.T) {
		left := m.Concat(m.Concat(v, v), v)
		right := m.Concat(v, m.Concat(v, v))
		assert.Equal(t, left, right, "Associativity law failed")
	})

	// All results should be VOID
	t.Run("all operations return VOID", func(t *testing.T) {
		assert.Equal(t, function.VOID, m.Concat(v, v))
		assert.Equal(t, function.VOID, m.Empty())
		assert.Equal(t, function.VOID, m.Concat(m.Empty(), v))
		assert.Equal(t, function.VOID, m.Concat(v, m.Empty()))
	})
}

// TestVoidMonoid_ConcatAll tests combining multiple Void values
func TestVoidMonoid_ConcatAll(t *testing.T) {
	m := VoidMonoid()
	concatAll := ConcatAll(m)

	tests := []struct {
		name     string
		input    []Void
		expected Void
	}{
		{
			name:     "empty slice",
			input:    []Void{},
			expected: function.VOID,
		},
		{
			name:     "single element",
			input:    []Void{function.VOID},
			expected: function.VOID,
		},
		{
			name:     "multiple elements",
			input:    []Void{function.VOID, function.VOID, function.VOID},
			expected: function.VOID,
		},
		{
			name:     "many elements",
			input:    make([]Void, 100),
			expected: function.VOID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize slice with VOID values
			for i := range tt.input {
				tt.input[i] = function.VOID
			}
			result := concatAll(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestVoidMonoid_Fold tests the Fold function with VoidMonoid
func TestVoidMonoid_Fold(t *testing.T) {
	m := VoidMonoid()
	fold := Fold(m)

	// Fold should behave identically to ConcatAll
	voids := []Void{function.VOID, function.VOID, function.VOID}
	result := fold(voids)
	assert.Equal(t, function.VOID, result)

	// Empty fold
	emptyResult := fold([]Void{})
	assert.Equal(t, function.VOID, emptyResult)
}

// TestVoidMonoid_Reverse tests that Reverse doesn't affect VoidMonoid
func TestVoidMonoid_Reverse(t *testing.T) {
	m := VoidMonoid()
	reversed := Reverse(m)

	// Since all Void values are identical, reverse should have no effect
	v := function.VOID

	assert.Equal(t, m.Concat(v, v), reversed.Concat(v, v))
	assert.Equal(t, m.Empty(), reversed.Empty())

	// Test identity laws still hold
	assert.Equal(t, v, reversed.Concat(reversed.Empty(), v))
	assert.Equal(t, v, reversed.Concat(v, reversed.Empty()))
}

// TestVoidMonoid_ToSemigroup tests conversion to Semigroup
func TestVoidMonoid_ToSemigroup(t *testing.T) {
	m := VoidMonoid()
	sg := ToSemigroup(m)

	// Should work as a semigroup
	result := sg.Concat(function.VOID, function.VOID)
	assert.Equal(t, function.VOID, result)

	// Verify it's the same underlying operation
	assert.Equal(t, m.Concat(function.VOID, function.VOID), sg.Concat(function.VOID, function.VOID))
}

// TestVoidMonoid_FunctionMonoid tests VoidMonoid with FunctionMonoid
func TestVoidMonoid_FunctionMonoid(t *testing.T) {
	m := VoidMonoid()
	funcMonoid := FunctionMonoid[string](m)

	// Create functions that return Void
	f1 := func(s string) Void { return function.VOID }
	f2 := func(s string) Void { return function.VOID }

	// Combine functions
	combined := funcMonoid.Concat(f1, f2)

	// Test combined function
	result := combined("test")
	assert.Equal(t, function.VOID, result)

	// Test empty function
	emptyFunc := funcMonoid.Empty()
	assert.Equal(t, function.VOID, emptyFunc("anything"))
}

// TestVoidMonoid_PracticalUsage demonstrates practical usage patterns
func TestVoidMonoid_PracticalUsage(t *testing.T) {
	m := VoidMonoid()

	// Simulate tracking that operations occurred without caring about results
	type Action func() Void

	actions := []Action{
		func() Void { return function.VOID }, // Action 1
		func() Void { return function.VOID }, // Action 2
		func() Void { return function.VOID }, // Action 3
	}

	// Execute all actions and collect results
	results := make([]Void, len(actions))
	for i, action := range actions {
		results[i] = action()
	}

	// Combine all results (all are VOID)
	finalResult := ConcatAll(m)(results)
	assert.Equal(t, function.VOID, finalResult)
}

// TestVoidMonoid_EdgeCases tests edge cases
func TestVoidMonoid_EdgeCases(t *testing.T) {
	m := VoidMonoid()

	t.Run("multiple concatenations", func(t *testing.T) {
		// Chain multiple Concat operations
		result := m.Concat(
			m.Concat(
				m.Concat(function.VOID, function.VOID),
				function.VOID,
			),
			function.VOID,
		)
		assert.Equal(t, function.VOID, result)
	})

	t.Run("concat with empty", func(t *testing.T) {
		// Various combinations with Empty()
		assert.Equal(t, function.VOID, m.Concat(m.Empty(), m.Empty()))
		assert.Equal(t, function.VOID, m.Concat(m.Concat(m.Empty(), function.VOID), m.Empty()))
	})

	t.Run("large slice", func(t *testing.T) {
		// Test with a large number of elements
		largeSlice := make([]Void, 10000)
		for i := range largeSlice {
			largeSlice[i] = function.VOID
		}
		result := ConcatAll(m)(largeSlice)
		assert.Equal(t, function.VOID, result)
	})
}

// TestVoidMonoid_TypeSafety verifies type safety
func TestVoidMonoid_TypeSafety(t *testing.T) {
	m := VoidMonoid()

	// Verify it implements Monoid interface
	var _ Monoid[Void] = m

	// Verify Empty returns correct type
	empty := m.Empty()
	var _ Void = empty

	// Verify Concat returns correct type
	result := m.Concat(function.VOID, function.VOID)
	var _ Void = result
}

// BenchmarkVoidMonoid_Concat benchmarks the Concat operation
func BenchmarkVoidMonoid_Concat(b *testing.B) {
	m := VoidMonoid()
	v := function.VOID

	b.ResetTimer()
	for b.Loop() {
		_ = m.Concat(v, v)
	}
}

// BenchmarkVoidMonoid_ConcatAll benchmarks combining multiple Void values
func BenchmarkVoidMonoid_ConcatAll(b *testing.B) {
	m := VoidMonoid()
	concatAll := ConcatAll(m)

	voids := make([]Void, 1000)
	for i := range voids {
		voids[i] = function.VOID
	}

	b.ResetTimer()
	for b.Loop() {
		_ = concatAll(voids)
	}
}

// BenchmarkVoidMonoid_Empty benchmarks the Empty operation
func BenchmarkVoidMonoid_Empty(b *testing.B) {
	m := VoidMonoid()

	b.ResetTimer()
	for b.Loop() {
		_ = m.Empty()
	}
}

// Made with Bob
