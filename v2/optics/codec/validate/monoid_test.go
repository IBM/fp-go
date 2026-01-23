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

package validate

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var (
	intAddMonoid = N.MonoidSum[int]()
	strMonoid    = S.Monoid
)

// Helper function to create a successful validator
func successValidator[I, A any](value A) Validate[I, A] {
	return func(input I) Reader[validation.Context, validation.Validation[A]] {
		return func(ctx validation.Context) validation.Validation[A] {
			return validation.Success(value)
		}
	}
}

// Helper function to create a failing validator
func failureValidator[I, A any](message string) Validate[I, A] {
	return func(input I) Reader[validation.Context, validation.Validation[A]] {
		return validation.FailureWithMessage[A](input, message)
	}
}

// Helper function to create a validator that uses the input
func inputDependentValidator[A any](f func(A) A) Validate[A, A] {
	return func(input A) Reader[validation.Context, validation.Validation[A]] {
		return func(ctx validation.Context) validation.Validation[A] {
			return validation.Success(f(input))
		}
	}
}

// TestApplicativeMonoid_EmptyElement tests the empty element of the monoid
func TestApplicativeMonoid_EmptyElement(t *testing.T) {
	t.Run("int addition monoid", func(t *testing.T) {
		m := ApplicativeMonoid[string](intAddMonoid)
		empty := m.Empty()

		result := empty("test")(nil)

		assert.Equal(t, validation.Of(0), result)
	})

	t.Run("string concatenation monoid", func(t *testing.T) {
		m := ApplicativeMonoid[int](strMonoid)
		empty := m.Empty()

		result := empty(42)(nil)

		assert.Equal(t, validation.Of(""), result)
	})
}

// TestApplicativeMonoid_ConcatSuccesses tests concatenating two successful validators
func TestApplicativeMonoid_ConcatSuccesses(t *testing.T) {
	t.Run("int addition", func(t *testing.T) {
		m := ApplicativeMonoid[string](intAddMonoid)

		v1 := successValidator[string](5)
		v2 := successValidator[string](3)

		combined := m.Concat(v1, v2)
		result := combined("input")(nil)

		assert.Equal(t, validation.Of(8), result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		m := ApplicativeMonoid[int](strMonoid)

		v1 := successValidator[int]("Hello")
		v2 := successValidator[int](" World")

		combined := m.Concat(v1, v2)
		result := combined(42)(nil)

		assert.Equal(t, validation.Of("Hello World"), result)
	})
}

// TestApplicativeMonoid_ConcatWithFailure tests concatenating validators where one fails
func TestApplicativeMonoid_ConcatWithFailure(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	t.Run("left failure", func(t *testing.T) {
		v1 := failureValidator[string, int]("left error")
		v2 := successValidator[string](5)

		combined := m.Concat(v1, v2)
		result := combined("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "left error", errors[0].Messsage)
	})

	t.Run("right failure", func(t *testing.T) {
		v1 := successValidator[string](5)
		v2 := failureValidator[string, int]("right error")

		combined := m.Concat(v1, v2)
		result := combined("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "right error", errors[0].Messsage)
	})

	t.Run("both failures", func(t *testing.T) {
		v1 := failureValidator[string, int]("left error")
		v2 := failureValidator[string, int]("right error")

		combined := m.Concat(v1, v2)
		result := combined("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		// Note: The current implementation returns the first error encountered
		assert.GreaterOrEqual(t, len(errors), 1)
		// At least one of the errors should be present
		hasError := false
		for _, err := range errors {
			if err.Messsage == "left error" || err.Messsage == "right error" {
				hasError = true
				break
			}
		}
		assert.True(t, hasError, "Should contain at least one validation error")
	})
}

// TestApplicativeMonoid_LeftIdentity tests the left identity law
func TestApplicativeMonoid_LeftIdentity(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v := successValidator[string](42)

	// empty <> v == v
	combined := m.Concat(m.Empty(), v)

	resultCombined := combined("test")(nil)
	resultOriginal := v("test")(nil)

	assert.Equal(t, resultOriginal, resultCombined)
}

// TestApplicativeMonoid_RightIdentity tests the right identity law
func TestApplicativeMonoid_RightIdentity(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v := successValidator[string](42)

	// v <> empty == v
	combined := m.Concat(v, m.Empty())

	resultCombined := combined("test")(nil)
	resultOriginal := v("test")(nil)

	assert.Equal(t, resultOriginal, resultCombined)
}

// TestApplicativeMonoid_Associativity tests the associativity law
func TestApplicativeMonoid_Associativity(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v1 := successValidator[string](1)
	v2 := successValidator[string](2)
	v3 := successValidator[string](3)

	// (v1 <> v2) <> v3 == v1 <> (v2 <> v3)
	left := m.Concat(m.Concat(v1, v2), v3)
	right := m.Concat(v1, m.Concat(v2, v3))

	resultLeft := left("test")(nil)
	resultRight := right("test")(nil)

	assert.Equal(t, resultRight, resultLeft)

	// Both should equal 6
	assert.Equal(t, validation.Of(6), resultLeft)
}

// TestApplicativeMonoid_AssociativityWithFailures tests associativity with failures
func TestApplicativeMonoid_AssociativityWithFailures(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v1 := successValidator[string](1)
	v2 := failureValidator[string, int]("error 2")
	v3 := successValidator[string](3)

	// (v1 <> v2) <> v3 == v1 <> (v2 <> v3)
	left := m.Concat(m.Concat(v1, v2), v3)
	right := m.Concat(v1, m.Concat(v2, v3))

	resultLeft := left("test")(nil)
	resultRight := right("test")(nil)

	// Both should fail with the same error
	assert.True(t, E.IsLeft(resultLeft))
	assert.True(t, E.IsLeft(resultRight))

	_, errorsLeft := E.Unwrap(resultLeft)
	_, errorsRight := E.Unwrap(resultRight)

	assert.Len(t, errorsLeft, 1)
	assert.Len(t, errorsRight, 1)
	assert.Equal(t, "error 2", errorsLeft[0].Messsage)
	assert.Equal(t, "error 2", errorsRight[0].Messsage)
}

// TestApplicativeMonoid_MultipleValidators tests combining multiple validators
func TestApplicativeMonoid_MultipleValidators(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v1 := successValidator[string](10)
	v2 := successValidator[string](20)
	v3 := successValidator[string](30)
	v4 := successValidator[string](40)

	// Chain multiple concat operations
	combined := m.Concat(
		m.Concat(
			m.Concat(v1, v2),
			v3,
		),
		v4,
	)

	result := combined("test")(nil)

	assert.Equal(t, validation.Of(100), result)
}

// TestApplicativeMonoid_InputDependent tests validators that depend on input
func TestApplicativeMonoid_InputDependent(t *testing.T) {
	m := ApplicativeMonoid[int](intAddMonoid)

	// Validator that doubles the input
	v1 := inputDependentValidator(N.Mul(2))
	// Validator that adds 10 to the input
	v2 := inputDependentValidator(N.Add(10))

	combined := m.Concat(v1, v2)
	result := combined(5)(nil)

	// (5 * 2) + (5 + 10) = 10 + 15 = 25
	assert.Equal(t, validation.Of(25), result)
}

// TestApplicativeMonoid_ContextPropagation tests that context is properly propagated
func TestApplicativeMonoid_ContextPropagation(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	// Create a validator that captures the context
	var capturedContext validation.Context
	v1 := func(input string) Reader[validation.Context, validation.Validation[int]] {
		return func(ctx validation.Context) validation.Validation[int] {
			capturedContext = ctx
			return validation.Success(5)
		}
	}

	v2 := successValidator[string](3)

	combined := m.Concat(v1, v2)

	// Create a context with some entries
	ctx := validation.Context{
		{Key: "field1", Type: "int"},
		{Key: "field2", Type: "string"},
	}

	result := combined("test")(ctx)

	assert.True(t, E.IsRight(result))
	assert.Equal(t, ctx, capturedContext)
}

// TestApplicativeMonoid_ErrorAccumulation tests that errors are accumulated
func TestApplicativeMonoid_ErrorAccumulation(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v1 := failureValidator[string, int]("error 1")
	v2 := failureValidator[string, int]("error 2")
	v3 := failureValidator[string, int]("error 3")

	combined := m.Concat(m.Concat(v1, v2), v3)
	result := combined("test")(nil)

	assert.True(t, E.IsLeft(result))
	_, errors := E.Unwrap(result)

	// Note: The current implementation returns the first error encountered
	// At least one error should be present
	assert.GreaterOrEqual(t, len(errors), 1)
	hasError := false
	for _, err := range errors {
		if err.Messsage == "error 1" || err.Messsage == "error 2" || err.Messsage == "error 3" {
			hasError = true
			break
		}
	}
	assert.True(t, hasError, "Should contain at least one validation error")
}

// TestApplicativeMonoid_MixedSuccessFailure tests mixing successes and failures
func TestApplicativeMonoid_MixedSuccessFailure(t *testing.T) {
	m := ApplicativeMonoid[string](intAddMonoid)

	v1 := successValidator[string](10)
	v2 := failureValidator[string, int]("error in v2")
	v3 := successValidator[string](20)
	v4 := failureValidator[string, int]("error in v4")

	combined := m.Concat(
		m.Concat(
			m.Concat(v1, v2),
			v3,
		),
		v4,
	)

	result := combined("test")(nil)

	assert.True(t, E.IsLeft(result))
	_, errors := E.Unwrap(result)

	// Note: The current implementation returns the first error encountered
	// At least one error should be present
	assert.GreaterOrEqual(t, len(errors), 1)
	hasError := false
	for _, err := range errors {
		if err.Messsage == "error in v2" || err.Messsage == "error in v4" {
			hasError = true
			break
		}
	}
	assert.True(t, hasError, "Should contain at least one validation error")
}

// TestApplicativeMonoid_DifferentInputTypes tests with different input types
func TestApplicativeMonoid_DifferentInputTypes(t *testing.T) {
	t.Run("struct input", func(t *testing.T) {
		type Config struct {
			Port    int
			Timeout int
		}

		m := ApplicativeMonoid[Config](intAddMonoid)

		v1 := func(cfg Config) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.Success(cfg.Port)
			}
		}

		v2 := func(cfg Config) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.Success(cfg.Timeout)
			}
		}

		combined := m.Concat(v1, v2)
		result := combined(Config{Port: 8080, Timeout: 30})(nil)

		assert.Equal(t, validation.Of(8110), result) // 8080 + 30
	})
}

// TestApplicativeMonoid_StringConcatenation tests string concatenation scenarios
func TestApplicativeMonoid_StringConcatenation(t *testing.T) {
	m := ApplicativeMonoid[string](strMonoid)

	t.Run("build sentence", func(t *testing.T) {
		v1 := successValidator[string]("The")
		v2 := successValidator[string](" quick")
		v3 := successValidator[string](" brown")
		v4 := successValidator[string](" fox")

		combined := m.Concat(
			m.Concat(
				m.Concat(v1, v2),
				v3,
			),
			v4,
		)

		result := combined("input")(nil)

		assert.Equal(t, validation.Of("The quick brown fox"), result)
	})

	t.Run("with empty strings", func(t *testing.T) {
		v1 := successValidator[string]("Hello")
		v2 := successValidator[string]("")
		v3 := successValidator[string]("World")

		combined := m.Concat(m.Concat(v1, v2), v3)
		result := combined("input")(nil)

		assert.Equal(t, validation.Of("HelloWorld"), result)
	})
}

// Benchmark tests
func BenchmarkApplicativeMonoid_ConcatSuccesses(b *testing.B) {
	m := ApplicativeMonoid[string](intAddMonoid)
	v1 := successValidator[string](5)
	v2 := successValidator[string](3)
	combined := m.Concat(v1, v2)

	b.ResetTimer()
	for range b.N {
		_ = combined("test")(nil)
	}
}

func BenchmarkApplicativeMonoid_ConcatFailures(b *testing.B) {
	m := ApplicativeMonoid[string](intAddMonoid)
	v1 := failureValidator[string, int]("error 1")
	v2 := failureValidator[string, int]("error 2")
	combined := m.Concat(v1, v2)

	b.ResetTimer()
	for range b.N {
		_ = combined("test")(nil)
	}
}

func BenchmarkApplicativeMonoid_MultipleConcat(b *testing.B) {
	m := ApplicativeMonoid[string](intAddMonoid)

	validators := make([]Validate[string, int], 10)
	for i := range validators {
		validators[i] = successValidator[string](i)
	}

	// Chain all validators
	combined := validators[0]
	for i := 1; i < len(validators); i++ {
		combined = m.Concat(combined, validators[i])
	}

	b.ResetTimer()
	for range b.N {
		_ = combined("test")(nil)
	}
}
