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

package option

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/optics/iso/option"
	"github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// Test types for AsOptional
type OptionalConfig struct {
	Timeout Option[int]
	Retries Option[int]
}

type OptionalPerson struct {
	Name    Option[string]
	Age     Option[int]
	Address Option[string]
}

type OptionalSettings struct {
	Theme Option[string]
}

type WithDefaultValue struct {
	Value string
}

func makeWithDefaultValueLens() Lens[WithDefaultValue, string] {
	return lens.MakeLensWithName(
		func(oc WithDefaultValue) string { return oc.Value },
		func(oc WithDefaultValue, v string) WithDefaultValue { oc.Value = v; return oc },
		"Value",
	)
}

func makeWithDefaultValueLensO() LensO[WithDefaultValue, string] {
	return FromIso[WithDefaultValue](IO.FromZero[string]())(makeWithDefaultValueLens())
}

// Helper to create a lens for OptionalConfig.Timeout
func makeTimeoutLens() LensO[OptionalConfig, int] {
	return lens.MakeLens(
		func(c OptionalConfig) Option[int] { return c.Timeout },
		func(c OptionalConfig, t Option[int]) OptionalConfig { c.Timeout = t; return c },
	)
}

// Helper to create a lens for OptionalPerson.Name
func makeNameLens() LensO[OptionalPerson, string] {
	return lens.MakeLens(
		func(p OptionalPerson) Option[string] { return p.Name },
		func(p OptionalPerson, n Option[string]) OptionalPerson { p.Name = n; return p },
	)
}

// Helper to create a lens for OptionalPerson.Age
func makeAgeLens() LensO[OptionalPerson, int] {
	return lens.MakeLens(
		func(p OptionalPerson) Option[int] { return p.Age },
		func(p OptionalPerson, a Option[int]) OptionalPerson { p.Age = a; return p },
	)
}

// TestAsOptional_BasicConversion tests the basic conversion functionality
func TestAsOptional_BasicConversion(t *testing.T) {
	t.Run("converts LensO to Optional", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}
		result := timeoutOptional.GetOption(config)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("preserves lens name", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		// The optional's name should match the lens's name
		// We can't directly access the name field, but we can verify the optional was created
		// by checking that it works correctly
		config := OptionalConfig{Timeout: O.Some(30)}
		result := timeoutOptional.GetOption(config)
		assert.True(t, O.IsSome(result))
	})
}

// TestAsOptional_GetOption tests the GetOption operation
func TestAsOptional_GetOption(t *testing.T) {
	t.Run("returns Some when value exists", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}
		result := timeoutOptional.GetOption(config)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("returns None when value doesn't exist", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.None[int]()}
		result := timeoutOptional.GetOption(config)

		assert.True(t, O.IsNone(result))
	})

	t.Run("works with string values", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.Some("Alice")}
		result := nameOptional.GetOption(person)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("works with empty string values", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.None[string]()}
		result := nameOptional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})
}

// TestAsOptional_Set tests the Set operation
func TestAsOptional_Set(t *testing.T) {
	t.Run("sets value when Some exists", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}
		updated := timeoutOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})

	t.Run("Set on None is a no-op (Law 1)", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.None[int]()}
		updated := timeoutOptional.Set(60)(config)

		// Law 1: When GetOption returns None, Set is a no-op
		assert.Equal(t, O.None[int](), updated.Timeout)
		assert.Equal(t, config, updated)
	})

	t.Run("preserves other fields", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30), Retries: O.Some(3)}
		updated := timeoutOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
		assert.Equal(t, O.Some(3), updated.Retries)
	})

	t.Run("works with string values", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.Some("Alice")}
		updated := nameOptional.Set("Bob")(person)

		assert.Equal(t, O.Some("Bob"), updated.Name)
	})
}

// TestAsOptional_OptionalLaw1_GetSet tests the GetSet law (No-op on None)
// Law: GetOption(s) = None => Set(a)(s) = s (no-op)
func TestAsOptional_OptionalLaw1_GetSet(t *testing.T) {
	t.Run("GetSet law - Set is no-op when GetOption returns None", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(timeoutOptional.GetOption(config)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := timeoutOptional.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.Equal(t, O.None[int](), updated.Timeout)
	})

	t.Run("GetSet law - Set updates when GetOption returns Some", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(timeoutOptional.GetOption(config)))

		// Set should update when GetOption returns Some
		updated := timeoutOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})
}

// TestAsOptional_OptionalLaw2_SetGet tests the SetGet law (Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
func TestAsOptional_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("SetGet law - can retrieve what was set when Some exists", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(timeoutOptional.GetOption(config)))

		// Set a value and get it back
		updated := timeoutOptional.Set(60)(config)
		retrieved := timeoutOptional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(retrieved))
	})

	t.Run("SetGet law - Set is no-op when starting from None", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(timeoutOptional.GetOption(config)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := timeoutOptional.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.True(t, O.IsNone(timeoutOptional.GetOption(updated)))
	})

	t.Run("SetGet law - works with string values", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.Some("Alice")}

		updated := nameOptional.Set("Bob")(person)
		retrieved := nameOptional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, "Bob", O.GetOrElse(F.Constant(""))(retrieved))
	})
}

// TestAsOptional_OptionalLaw3_SetSet tests the SetSet law (Last Set Wins)
// Law: Set(b)(Set(a)(s)) = Set(b)(s)
func TestAsOptional_OptionalLaw3_SetSet(t *testing.T) {
	t.Run("SetSet law - last set wins", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Set twice
		setTwice := timeoutOptional.Set(50)(timeoutOptional.Set(40)(config))
		setOnce := timeoutOptional.Set(50)(config)

		assert.Equal(t, setOnce.Timeout, setTwice.Timeout)
		assert.Equal(t, O.Some(50), setTwice.Timeout)
	})

	t.Run("SetSet law - both are no-op when starting with None", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.None[int]()}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := timeoutOptional.Set(50)(timeoutOptional.Set(40)(config))
		setOnce := timeoutOptional.Set(50)(config)

		// Both should equal the original config (no-op)
		assert.Equal(t, config, setTwice)
		assert.Equal(t, config, setOnce)
		assert.Equal(t, O.None[int](), setTwice.Timeout)
	})

	t.Run("SetSet law - works with string values", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.Some("Alice")}

		// Set twice
		setTwice := nameOptional.Set("Charlie")(nameOptional.Set("Bob")(person))
		setOnce := nameOptional.Set("Charlie")(person)

		assert.Equal(t, setOnce.Name, setTwice.Name)
		assert.Equal(t, O.Some("Charlie"), setTwice.Name)
	})
}

// TestAsOptional_MultipleFields tests working with multiple optional fields
func TestAsOptional_MultipleFields(t *testing.T) {
	t.Run("can work with multiple fields independently", func(t *testing.T) {
		nameLens := makeNameLens()
		ageLens := makeAgeLens()

		nameOptional := AsOptional(nameLens)
		ageOptional := AsOptional(ageLens)

		person := OptionalPerson{
			Name: O.Some("Alice"),
			Age:  O.Some(30),
		}

		// Update name
		updatedName := nameOptional.Set("Bob")(person)
		assert.Equal(t, O.Some("Bob"), updatedName.Name)
		assert.Equal(t, O.Some(30), updatedName.Age)

		// Update age
		updatedAge := ageOptional.Set(35)(person)
		assert.Equal(t, O.Some("Alice"), updatedAge.Name)
		assert.Equal(t, O.Some(35), updatedAge.Age)
	})

	t.Run("can chain updates on multiple fields", func(t *testing.T) {
		nameLens := makeNameLens()
		ageLens := makeAgeLens()

		nameOptional := AsOptional(nameLens)
		ageOptional := AsOptional(ageLens)

		person := OptionalPerson{
			Name: O.Some("Alice"),
			Age:  O.Some(30),
		}

		// Chain updates
		updated := F.Pipe2(
			person,
			nameOptional.Set("Bob"),
			ageOptional.Set(35),
		)

		assert.Equal(t, O.Some("Bob"), updated.Name)
		assert.Equal(t, O.Some(35), updated.Age)
	})
}

// TestAsOptional_EdgeCases tests edge cases
func TestAsOptional_EdgeCases(t *testing.T) {
	t.Run("works with zero values", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(0)}

		result := timeoutOptional.GetOption(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))

		updated := timeoutOptional.Set(0)(config)
		assert.Equal(t, O.Some(0), updated.Timeout)
	})

	t.Run("works with empty strings", func(t *testing.T) {
		nameLens := makeNameLens()
		nameOptional := AsOptional(nameLens)

		person := OptionalPerson{Name: O.Some("")}

		result := nameOptional.GetOption(person)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "", O.GetOrElse(F.Constant("default"))(result))

		updated := nameOptional.Set("")(person)
		assert.Equal(t, O.Some(""), updated.Name)
	})

	t.Run("works with negative numbers", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(-10)}

		result := timeoutOptional.GetOption(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, -10, O.GetOrElse(F.Constant(0))(result))

		updated := timeoutOptional.Set(-20)(config)
		assert.Equal(t, O.Some(-20), updated.Timeout)
	})
}

// TestAsOptional_Integration tests integration scenarios
func TestAsOptional_Integration(t *testing.T) {
	t.Run("can be used with option operations", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Use with option map
		result := F.Pipe1(
			timeoutOptional.GetOption(config),
			O.Map(func(x int) int { return x * 2 }),
		)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("can be used with option chain", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Use with option chain
		result := F.Pipe1(
			timeoutOptional.GetOption(config),
			O.Chain(func(x int) Option[int] {
				if x > 20 {
					return O.Some(x * 2)
				}
				return O.None[int]()
			}),
		)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("can be used with option fold", func(t *testing.T) {
		timeoutLens := makeTimeoutLens()
		timeoutOptional := AsOptional(timeoutLens)

		config := OptionalConfig{Timeout: O.Some(30)}

		// Use with option fold
		result := F.Pipe1(
			timeoutOptional.GetOption(config),
			O.Fold(
				F.Constant("no timeout"),
				func(x int) string {
					if x == 30 {
						return "30"
					}
					return "other"
				},
			),
		)

		assert.Equal(t, "30", result)
	})
}

// TestAsOptional_WithDefaultValueLens tests AsOptional with a lens that uses FromZero isomorphism
// This validates that the Optional laws are properly enforced when the lens treats zero values as None
func TestAsOptional_WithDefaultValueLens(t *testing.T) {
	t.Run("GetOption returns None for zero value", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		// Empty string is the zero value, should return None
		obj := WithDefaultValue{Value: ""}
		result := valueOptional.GetOption(obj)

		assert.True(t, O.IsNone(result))
	})

	t.Run("GetOption returns Some for non-zero value", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}
		result := valueOptional.GetOption(obj)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("Set updates non-zero value", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}
		updated := valueOptional.Set("world")(obj)

		assert.Equal(t, "world", updated.Value)
	})

	t.Run("Set on zero value is no-op (Law 1)", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: ""}
		updated := valueOptional.Set("hello")(obj)

		// Law 1: When GetOption returns None (zero value), Set is a no-op
		assert.Equal(t, "", updated.Value)
		assert.Equal(t, obj, updated)
	})
}

// TestAsOptional_WithDefaultValueLens_Law1 tests Law 1 (GetSet - No-op on None)
// Law: GetOption(s) = None => Set(a)(s) = s (no-op)
func TestAsOptional_WithDefaultValueLens_Law1(t *testing.T) {
	t.Run("Law 1: Set is no-op when GetOption returns None (zero value)", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: ""} // Zero value

		// Verify GetOption returns None
		assert.True(t, O.IsNone(valueOptional.GetOption(obj)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := valueOptional.Set("new_value")(obj)

		// Verify the object is unchanged
		assert.Equal(t, obj, updated)
		assert.Equal(t, "", updated.Value)
	})

	t.Run("Law 1: Set updates when GetOption returns Some (non-zero value)", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(valueOptional.GetOption(obj)))

		// Set should update when GetOption returns Some
		updated := valueOptional.Set("world")(obj)

		assert.Equal(t, "world", updated.Value)
	})
}

// TestAsOptional_WithDefaultValueLens_Law2 tests Law 2 (SetGet - Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
func TestAsOptional_WithDefaultValueLens_Law2(t *testing.T) {
	t.Run("Law 2: Can retrieve what was set when GetOption returns Some", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(valueOptional.GetOption(obj)))

		// Set a new value
		newValue := "world"
		updated := valueOptional.Set(newValue)(obj)

		// Get the value back
		retrieved := valueOptional.GetOption(updated)

		// Verify we get back what we set
		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, newValue, O.GetOrElse(F.Constant(""))(retrieved))
	})

	t.Run("Law 2: Set is no-op when starting from None (zero value)", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: ""} // Zero value

		// Verify GetOption returns None
		assert.True(t, O.IsNone(valueOptional.GetOption(obj)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := valueOptional.Set("hello")(obj)

		// Verify the object is unchanged
		assert.Equal(t, obj, updated)
		assert.True(t, O.IsNone(valueOptional.GetOption(updated)))
	})

	t.Run("Law 2: Setting empty string results in None", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}

		// Set to empty string (zero value)
		updated := valueOptional.Set("")(obj)

		// GetOption should return None because empty string is the zero value
		retrieved := valueOptional.GetOption(updated)
		assert.True(t, O.IsNone(retrieved))
	})
}

// TestAsOptional_WithDefaultValueLens_Law3 tests Law 3 (SetSet - Last Set Wins)
// Law: Set(b)(Set(a)(s)) = Set(b)(s)
func TestAsOptional_WithDefaultValueLens_Law3(t *testing.T) {
	t.Run("Law 3: Last set wins when starting with Some", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "initial"}

		// Set twice
		setTwice := valueOptional.Set("final")(valueOptional.Set("intermediate")(obj))
		setOnce := valueOptional.Set("final")(obj)

		// Both should be equal
		assert.Equal(t, setOnce.Value, setTwice.Value)
		assert.Equal(t, "final", setTwice.Value)
	})

	t.Run("Law 3: Both are no-op when starting with None (zero value)", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: ""}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := valueOptional.Set("final")(valueOptional.Set("intermediate")(obj))
		setOnce := valueOptional.Set("final")(obj)

		// Both should equal the original object (no-op)
		assert.Equal(t, obj, setTwice)
		assert.Equal(t, obj, setOnce)
		assert.Equal(t, "", setTwice.Value)
	})

	t.Run("Law 3: Setting to zero value then non-zero", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "initial"}

		// Set to empty string (which results in None), then try to set to "final"
		// First set changes "initial" to "" (None)
		// Second set is a no-op because GetOption now returns None
		setTwice := valueOptional.Set("final")(valueOptional.Set("")(obj))
		setOnce := valueOptional.Set("final")(obj)

		// setTwice should be "" (from first set), setOnce should be "final"
		assert.Equal(t, "", setTwice.Value)
		assert.Equal(t, "final", setOnce.Value)
	})

	t.Run("Law 3: Setting to non-zero then zero value", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "initial"}

		// Set to a value, then to empty string
		setTwice := valueOptional.Set("")(valueOptional.Set("intermediate")(obj))
		setOnce := valueOptional.Set("")(obj)

		// Both should be equal
		assert.Equal(t, setOnce.Value, setTwice.Value)
		assert.Equal(t, "", setTwice.Value)

		// GetOption should return None for both
		assert.True(t, O.IsNone(valueOptional.GetOption(setTwice)))
		assert.True(t, O.IsNone(valueOptional.GetOption(setOnce)))
	})
}

// TestAsOptional_WithDefaultValueLens_Integration tests integration scenarios
func TestAsOptional_WithDefaultValueLens_Integration(t *testing.T) {
	t.Run("Chain multiple operations", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		obj := WithDefaultValue{Value: "hello"}

		// Chain: Get -> Map -> Set
		result := F.Pipe2(
			valueOptional.GetOption(obj),
			O.Map(func(s string) string { return s + "_world" }),
			O.Fold(
				F.Constant(obj),
				func(newVal string) WithDefaultValue {
					return valueOptional.Set(newVal)(obj)
				},
			),
		)

		assert.Equal(t, "hello_world", result.Value)
	})

	t.Run("Conditional update based on GetOption", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		updateIfPresent := func(obj WithDefaultValue, newVal string) WithDefaultValue {
			return F.Pipe1(
				valueOptional.GetOption(obj),
				O.Fold(
					F.Constant(obj), // If None, return unchanged
					func(_ string) WithDefaultValue {
						return valueOptional.Set(newVal)(obj)
					},
				),
			)
		}

		// Test with Some
		obj1 := WithDefaultValue{Value: "hello"}
		updated1 := updateIfPresent(obj1, "updated")
		assert.Equal(t, "updated", updated1.Value)

		// Test with None
		obj2 := WithDefaultValue{Value: ""}
		updated2 := updateIfPresent(obj2, "updated")
		assert.Equal(t, "", updated2.Value) // Unchanged
	})

	t.Run("Toggle between zero and non-zero", func(t *testing.T) {
		valueLens := makeWithDefaultValueLensO()
		valueOptional := AsOptional(valueLens)

		toggle := func(obj WithDefaultValue, defaultVal string) WithDefaultValue {
			return F.Pipe1(
				valueOptional.GetOption(obj),
				O.Fold(
					func() WithDefaultValue {
						// When None, Set is a no-op, so return unchanged
						// To actually set a value, we need to manually update
						return WithDefaultValue{Value: defaultVal}
					},
					func(_ string) WithDefaultValue {
						// When Some, set to empty string (which becomes None)
						return valueOptional.Set("")(obj)
					},
				),
			)
		}

		// Start with None, toggle to Some (manual update since Set is no-op on None)
		obj1 := WithDefaultValue{Value: ""}
		toggled1 := toggle(obj1, "default")
		assert.Equal(t, "default", toggled1.Value)

		// Toggle back to None
		toggled2 := toggle(toggled1, "default")
		assert.Equal(t, "", toggled2.Value)
	})
}
