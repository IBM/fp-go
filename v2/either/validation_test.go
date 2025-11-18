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

package either

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	SG "github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestMonadApV_BothRight tests MonadApV when both function and value are Right
func TestMonadApV_BothRight(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := MonadApV[int, int](sg)

	// Both are Right - should apply function
	fab := Right[string](N.Mul(2))
	fa := Right[string](21)

	result := applyV(fab, fa)

	assert.True(t, IsRight(result))
	assert.Equal(t, Right[string](42), result)
}

// TestMonadApV_BothLeft tests MonadApV when both function and value are Left
func TestMonadApV_BothLeft(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := MonadApV[int, int](sg)

	// Both are Left - should combine errors
	fab := Left[func(int) int]("error1")
	fa := Left[int]("error2")

	result := applyV(fab, fa)

	assert.True(t, IsLeft(result))
	// When both are Left, errors are combined as: fa error + fab error
	assert.Equal(t, Left[int]("error1; error2"), result)
}

// TestMonadApV_LeftFunction tests MonadApV when function is Left and value is Right
func TestMonadApV_LeftFunction(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := MonadApV[int, int](sg)

	// Function is Left, value is Right - should return function's error
	fab := Left[func(int) int]("function error")
	fa := Right[string](21)

	result := applyV(fab, fa)

	assert.True(t, IsLeft(result))
	assert.Equal(t, Left[int]("function error"), result)
}

// TestMonadApV_LeftValue tests MonadApV when function is Right and value is Left
func TestMonadApV_LeftValue(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := MonadApV[int, int](sg)

	// Function is Right, value is Left - should return value's error
	fab := Right[string](N.Mul(2))
	fa := Left[int]("value error")

	result := applyV(fab, fa)

	assert.True(t, IsLeft(result))
	assert.Equal(t, Left[int]("value error"), result)
}

// TestMonadApV_WithSliceSemigroup tests MonadApV with a slice-based semigroup
func TestMonadApV_WithSliceSemigroup(t *testing.T) {
	// Create a semigroup that concatenates slices
	sg := SG.MakeSemigroup(func(a, b []string) []string {
		return append(a, b...)
	})

	// Create the validation applicative
	applyV := MonadApV[string, string](sg)

	// Both are Left with slice errors
	fab := Left[func(string) string]([]string{"error1", "error2"})
	fa := Left[string]([]string{"error3", "error4"})

	result := applyV(fab, fa)

	assert.True(t, IsLeft(result))
	// When both are Left, errors are combined as: fa errors + fab errors
	expected := Left[string]([]string{"error1", "error2", "error3", "error4"})
	assert.Equal(t, expected, result)
}

// TestMonadApV_ComplexFunction tests MonadApV with a more complex function
func TestMonadApV_ComplexFunction(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := SG.MakeSemigroup(func(a, b string) string {
		return a + " | " + b
	})

	// Create the validation applicative
	applyV := MonadApV[string, int](sg)

	// Test with a function that transforms the value
	fab := Right[string](func(x int) string {
		if x > 0 {
			return "positive"
		}
		return "non-positive"
	})
	fa := Right[string](42)

	result := applyV(fab, fa)

	assert.True(t, IsRight(result))
	assert.Equal(t, Right[string]("positive"), result)
}

// TestApV_BothRight tests ApV when both function and value are Right
func TestApV_BothRight(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Both are Right - should apply function
	fa := Right[string](21)
	fab := Right[string](N.Mul(2))

	result := applyV(fa)(fab)

	assert.True(t, IsRight(result))
	assert.Equal(t, Right[string](42), result)
}

// TestApV_BothLeft tests ApV when both function and value are Left
func TestApV_BothLeft(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Both are Left - should combine errors
	fa := Left[int]("error2")
	fab := Left[func(int) int]("error1")

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	// When both are Left, errors are combined as: fa error + fab error
	assert.Equal(t, Left[int]("error1; error2"), result)
}

// TestApV_LeftFunction tests ApV when function is Left and value is Right
func TestApV_LeftFunction(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Function is Left, value is Right - should return function's error
	fa := Right[string](21)
	fab := Left[func(int) int]("function error")

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	assert.Equal(t, Left[int]("function error"), result)
}

// TestApV_LeftValue tests ApV when function is Right and value is Left
func TestApV_LeftValue(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup("; ")

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Function is Right, value is Left - should return value's error
	fa := Left[int]("value error")
	fab := Right[string](N.Mul(2))

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	assert.Equal(t, Left[int]("value error"), result)
}

// TestApV_Composition tests ApV with function composition
func TestApV_Composition(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := SG.MakeSemigroup(func(a, b string) string {
		return a + " & " + b
	})

	// Create the validation applicative
	applyV := ApV[string, int](sg)

	// Test composition with pipe
	fa := Right[string](10)
	fab := Right[string](func(x int) string {
		return F.Pipe1(x, func(n int) string {
			if n >= 10 {
				return "large"
			}
			return "small"
		})
	})

	result := F.Pipe1(fa, applyV)(fab)

	assert.True(t, IsRight(result))
	assert.Equal(t, Right[string]("large"), result)
}

// TestApV_WithStructSemigroup tests ApV with a custom struct semigroup
func TestApV_WithStructSemigroup(t *testing.T) {
	type ValidationErrors struct {
		Errors []string
	}

	// Create a semigroup that combines validation errors
	sg := SG.MakeSemigroup(func(a, b ValidationErrors) ValidationErrors {
		return ValidationErrors{
			Errors: append(append([]string{}, a.Errors...), b.Errors...),
		}
	})

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Both are Left with validation errors
	fa := Left[int](ValidationErrors{Errors: []string{"field2: invalid"}})
	fab := Left[func(int) int](ValidationErrors{Errors: []string{"field1: required"}})

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	// When both are Left, errors are combined as: fa errors + fab errors
	expected := Left[int](ValidationErrors{
		Errors: []string{"field1: required", "field2: invalid"},
	})
	assert.Equal(t, expected, result)
}

// TestApV_MultipleValidations tests ApV with multiple validation steps
func TestApV_MultipleValidations(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := SG.MakeSemigroup(func(a, b string) string {
		return a + ", " + b
	})

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Simulate multiple validation failures
	validation1 := Left[int]("age must be positive")
	validation2 := Left[func(int) int]("name is required")

	result := applyV(validation1)(validation2)

	assert.True(t, IsLeft(result))
	// When both are Left, errors are combined as: validation1 error + validation2 error
	assert.Equal(t, Left[int]("name is required, age must be positive"), result)
}

// TestMonadApV_DifferentTypes tests MonadApV with different input and output types
func TestMonadApV_DifferentTypes(t *testing.T) {
	// Create a semigroup for string concatenation
	sg := S.IntersperseSemigroup(" + ")

	// Create the validation applicative
	applyV := MonadApV[string, int](sg)

	// Function converts int to string
	fab := Right[string](func(x int) string {
		return F.Pipe1(x, func(n int) string {
			if n == 0 {
				return "zero"
			} else if n > 0 {
				return "positive"
			}
			return "negative"
		})
	})
	fa := Right[string](-5)

	result := applyV(fab, fa)

	assert.True(t, IsRight(result))
	assert.Equal(t, Right[string]("negative"), result)
}

// TestApV_FirstSemigroup tests ApV with First semigroup (always returns first error)
func TestApV_FirstSemigroup(t *testing.T) {
	// Use First semigroup which always returns the first value
	sg := SG.First[string]()

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Both are Left - should return first error
	fa := Left[int]("error2")
	fab := Left[func(int) int]("error1")

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	// First semigroup returns the first value, which is fab's error
	assert.Equal(t, Left[int]("error1"), result)
}

// TestApV_LastSemigroup tests ApV with Last semigroup (always returns last error)
func TestApV_LastSemigroup(t *testing.T) {
	// Use Last semigroup which always returns the last value
	sg := SG.Last[string]()

	// Create the validation applicative
	applyV := ApV[int, int](sg)

	// Both are Left - should return last error
	fa := Left[int]("error2")
	fab := Left[func(int) int]("error1")

	result := applyV(fa)(fab)

	assert.True(t, IsLeft(result))
	// Last semigroup returns the last value, which is fa's error
	assert.Equal(t, Left[int]("error2"), result)
}
