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

package lens

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// Test types for composition

// Result is a sum type representing computation results
type Result interface {
	isResult()
}

type Success struct {
	Value   int
	Message string
}

func (Success) isResult() {}

type Failure struct {
	Error string
	Code  int
}

func (Failure) isResult() {}

type Pending struct {
	Progress int
}

func (Pending) isResult() {}

// Helper to create a prism for Success variant
func successPrism() P.Prism[Result, Success] {
	return P.MakePrism(
		func(r Result) O.Option[Success] {
			if s, ok := r.(Success); ok {
				return O.Some(s)
			}
			return O.None[Success]()
		},
		func(s Success) Result { return s },
	)
}

// Helper to create a prism for Failure variant
func failurePrism() P.Prism[Result, Failure] {
	return P.MakePrism(
		func(r Result) O.Option[Failure] {
			if f, ok := r.(Failure); ok {
				return O.Some(f)
			}
			return O.None[Failure]()
		},
		func(f Failure) Result { return f },
	)
}

// Helper to create a lens for Success.Value field
func successValueLens() L.Lens[Success, int] {
	return L.MakeLens(
		func(s Success) int { return s.Value },
		func(s Success, v int) Success { s.Value = v; return s },
	)
}

// Helper to create a lens for Success.Message field
func successMessageLens() L.Lens[Success, string] {
	return L.MakeLens(
		func(s Success) string { return s.Message },
		func(s Success, m string) Success { s.Message = m; return s },
	)
}

// Helper to create a lens for Failure.Code field
func failureCodeLens() L.Lens[Failure, int] {
	return L.MakeLens(
		func(f Failure) int { return f.Code },
		func(f Failure, c int) Failure { f.Code = c; return f },
	)
}

// TestCompose_BasicFunctionality tests basic composition behavior
func TestCompose_BasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when prism matches", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns None when prism doesn't match", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}
		value := optional.GetOption(result)

		assert.True(t, O.IsNone(value))
	})

	t.Run("Set updates value when prism matches", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}
		updated := optional.Set(100)(result)

		// Verify the update
		value := optional.GetOption(updated)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))

		// Verify other fields are preserved
		if s, ok := updated.(Success); ok {
			assert.Equal(t, "ok", s.Message)
		} else {
			t.Fatal("Expected Success type")
		}
	})

	t.Run("Set is no-op when prism doesn't match", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}
		updated := optional.Set(100)(result)

		// Verify nothing changed (no-op)
		assert.Equal(t, result, updated)
	})
}

// TestCompose_OptionalLaw1_GetSet tests Law 1 (No-op on None)
// Law: GetOption(s) = None => Set(b)(s) = s
func TestCompose_OptionalLaw1_GetSet(t *testing.T) {
	t.Run("Law 1: Set is no-op when GetOption returns None", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(optional.GetOption(result)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := optional.Set(100)(result)

		// Verify the result is unchanged
		assert.Equal(t, result, updated)
	})

	t.Run("Law 1: Set updates when GetOption returns Some", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(optional.GetOption(result)))

		// Set should update when GetOption returns Some
		updated := optional.Set(100)(result)

		value := optional.GetOption(updated)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Law 1: Multiple Set operations on None are all no-ops", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}

		// Multiple Set operations should all be no-ops
		updated := optional.Set(100)(optional.Set(50)(result))

		assert.Equal(t, result, updated)
	})
}

// TestCompose_OptionalLaw2_SetGet tests Law 2 (Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(b)(s)) = Some(b)
func TestCompose_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("Law 2: Can retrieve what was set when Some exists", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(optional.GetOption(result)))

		// Set a value and get it back
		newValue := 100
		updated := optional.Set(newValue)(result)
		retrieved := optional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, newValue, O.GetOrElse(F.Constant(0))(retrieved))
	})

	t.Run("Law 2: Set is no-op when starting from None", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(optional.GetOption(result)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := optional.Set(100)(result)

		// Verify the result is unchanged
		assert.Equal(t, result, updated)
		assert.True(t, O.IsNone(optional.GetOption(updated)))
	})

	t.Run("Law 2: Works with zero values", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}

		// Set to zero value
		updated := optional.Set(0)(result)
		retrieved := optional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(retrieved))
	})

	t.Run("Law 2: Works with string fields", func(t *testing.T) {
		prism := successPrism()
		lens := successMessageLens()
		optional := Compose[Result, Success, string](lens)(prism)

		result := Success{Value: 42, Message: "ok"}

		// Set a new message
		updated := optional.Set("updated")(result)
		retrieved := optional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, "updated", O.GetOrElse(F.Constant(""))(retrieved))
	})
}

// TestCompose_OptionalLaw3_SetSet tests Law 3 (Last Set Wins)
// Law: Set(c)(Set(b)(s)) = Set(c)(s)
func TestCompose_OptionalLaw3_SetSet(t *testing.T) {
	t.Run("Law 3: Last set wins when starting with Some", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}

		// Set twice
		setTwice := optional.Set(100)(optional.Set(50)(result))
		setOnce := optional.Set(100)(result)

		// They should be equal
		assert.Equal(t, setOnce, setTwice)

		// Verify final value
		value := optional.GetOption(setTwice)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Law 3: Both are no-op when starting with None", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := optional.Set(100)(optional.Set(50)(result))
		setOnce := optional.Set(100)(result)

		// Both should equal the original result (no-op)
		assert.Equal(t, result, setTwice)
		assert.Equal(t, result, setOnce)
	})

	t.Run("Law 3: Works with multiple different values", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 1, Message: "ok"}

		// Set multiple times
		setMultiple := optional.Set(4)(optional.Set(3)(optional.Set(2)(result)))
		setOnce := optional.Set(4)(result)

		assert.Equal(t, setOnce, setMultiple)

		// Verify final value
		value := optional.GetOption(setMultiple)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 4, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Law 3: Works with string fields", func(t *testing.T) {
		prism := successPrism()
		lens := successMessageLens()
		optional := Compose[Result, Success, string](lens)(prism)

		result := Success{Value: 42, Message: "initial"}

		// Set twice
		setTwice := optional.Set("final")(optional.Set("intermediate")(result))
		setOnce := optional.Set("final")(result)

		assert.Equal(t, setOnce, setTwice)

		// Verify final value
		value := optional.GetOption(setTwice)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, "final", O.GetOrElse(F.Constant(""))(value))
	})
}

// TestCompose_MultipleVariants tests composition with different prism variants
func TestCompose_MultipleVariants(t *testing.T) {
	t.Run("Success variant with Value field", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "ok"}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Success variant with Message field", func(t *testing.T) {
		prism := successPrism()
		lens := successMessageLens()
		optional := Compose[Result, Success, string](lens)(prism)

		result := Success{Value: 42, Message: "hello"}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(value))
	})

	t.Run("Failure variant with Code field", func(t *testing.T) {
		prism := failurePrism()
		lens := failureCodeLens()
		optional := Compose[Result, Failure, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 500, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Cross-variant no-op", func(t *testing.T) {
		// Try to use Success optional on Failure result
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Failure{Error: "failed", Code: 500}

		// GetOption should return None
		value := optional.GetOption(result)
		assert.True(t, O.IsNone(value))

		// Set should be no-op
		updated := optional.Set(100)(result)
		assert.Equal(t, result, updated)
	})
}

// TestCompose_EdgeCases tests edge cases and boundary conditions
func TestCompose_EdgeCases(t *testing.T) {
	t.Run("Works with zero values", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 0, Message: ""}

		value := optional.GetOption(result)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(value))

		updated := optional.Set(0)(result)
		retrievedValue := optional.GetOption(updated)
		assert.True(t, O.IsSome(retrievedValue))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(retrievedValue))
	})

	t.Run("Works with empty strings", func(t *testing.T) {
		prism := successPrism()
		lens := successMessageLens()
		optional := Compose[Result, Success, string](lens)(prism)

		result := Success{Value: 42, Message: ""}

		value := optional.GetOption(result)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, "", O.GetOrElse(F.Constant("default"))(value))

		updated := optional.Set("")(result)
		retrievedValue := optional.GetOption(updated)
		assert.True(t, O.IsSome(retrievedValue))
		assert.Equal(t, "", O.GetOrElse(F.Constant("default"))(retrievedValue))
	})

	t.Run("Works with negative numbers", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: -10, Message: "ok"}

		value := optional.GetOption(result)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, -10, O.GetOrElse(F.Constant(0))(value))

		updated := optional.Set(-20)(result)
		retrievedValue := optional.GetOption(updated)
		assert.True(t, O.IsSome(retrievedValue))
		assert.Equal(t, -20, O.GetOrElse(F.Constant(0))(retrievedValue))
	})

	t.Run("Preserves other fields when updating", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 42, Message: "important"}

		updated := optional.Set(100)(result)

		// Verify Value was updated
		value := optional.GetOption(updated)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))

		// Verify Message was preserved
		if s, ok := updated.(Success); ok {
			assert.Equal(t, "important", s.Message)
		} else {
			t.Fatal("Expected Success type")
		}
	})
}

// TestCompose_Integration tests integration scenarios
func TestCompose_Integration(t *testing.T) {
	t.Run("Can be used with option operations", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 30, Message: "ok"}

		// Use with option map
		value := F.Pipe1(
			optional.GetOption(result),
			O.Map(func(x int) int { return x * 2 }),
		)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Can be used with option chain", func(t *testing.T) {
		prism := successPrism()
		lens := successValueLens()
		optional := Compose[Result, Success, int](lens)(prism)

		result := Success{Value: 30, Message: "ok"}

		// Use with option chain
		value := F.Pipe1(
			optional.GetOption(result),
			O.Chain(func(x int) O.Option[int] {
				if x > 20 {
					return O.Some(x * 2)
				}
				return O.None[int]()
			}),
		)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("Can chain multiple updates", func(t *testing.T) {
		prism := successPrism()
		valueLens := successValueLens()
		messageLens := successMessageLens()

		valueOptional := Compose[Result, Success, int](valueLens)(prism)
		messageOptional := Compose[Result, Success, string](messageLens)(prism)

		result := Success{Value: 42, Message: "initial"}

		// Chain updates to different fields
		updated := messageOptional.Set("updated")(valueOptional.Set(100)(result))

		// Verify both updates
		value := valueOptional.GetOption(updated)
		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))

		message := messageOptional.GetOption(updated)
		assert.True(t, O.IsSome(message))
		assert.Equal(t, "updated", O.GetOrElse(F.Constant(""))(message))
	})
}

// TestCompose_DocumentationExample tests the example from the documentation
func TestCompose_DocumentationExample(t *testing.T) {
	// Prism to focus on Success variant
	successPrism := P.MakePrism(
		func(r Result) O.Option[Success] {
			if s, ok := r.(Success); ok {
				return O.Some(s)
			}
			return O.None[Success]()
		},
		func(s Success) Result { return s },
	)

	// Lens to focus on Value field within Success
	valueLens := L.MakeLens(
		func(s Success) int { return s.Value },
		func(s Success, v int) Success { s.Value = v; return s },
	)

	// Compose to create Optional[Result, int]
	resultValueOptional := Compose[Result, Success, int](valueLens)(successPrism)

	// Use the optional
	result := Success{Value: 42}
	value := resultValueOptional.GetOption(result) // Some(42)
	assert.True(t, O.IsSome(value))
	assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(value))

	updated := resultValueOptional.Set(100)(result) // Success{Value: 100}
	updatedValue := resultValueOptional.GetOption(updated)
	assert.True(t, O.IsSome(updatedValue))
	assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(updatedValue))

	// Set is no-op when prism doesn't match (Law 1)
	failure := Failure{Error: "failed"}
	unchanged := resultValueOptional.Set(100)(failure) // failure (unchanged)
	assert.Equal(t, failure, unchanged)
}

// Made with Bob
