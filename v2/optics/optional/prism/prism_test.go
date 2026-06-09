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

package prism

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// Test types for AsOptional
type Result interface {
	isResult()
}

type Success struct {
	Value int
}

func (Success) isResult() {}

type Failure struct {
	Error string
}

func (Failure) isResult() {}

type Config struct {
	Timeout O.Option[int]
	Retries O.Option[int]
}

type Person struct {
	Name O.Option[string]
	Age  O.Option[int]
}

// Helper to create a prism for Success variant
func makeSuccessPrism() P.Prism[Result, int] {
	return P.MakePrism(
		func(r Result) O.Option[int] {
			if s, ok := r.(Success); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) Result { return Success{Value: v} },
	)
}

// Helper to create an optional for Config.Timeout
func makeTimeoutOptional() OPT.Optional[Config, O.Option[int]] {
	return OPT.MakeOptional(
		func(c Config) O.Option[O.Option[int]] {
			return O.Some(c.Timeout)
		},
		func(c Config, opt O.Option[int]) Config {
			c.Timeout = opt
			return c
		},
	)
}

// Helper to create an optional for Person.Name
func makeNameOptional() OPT.Optional[Person, O.Option[string]] {
	return OPT.MakeOptional(
		func(p Person) O.Option[O.Option[string]] {
			return O.Some(p.Name)
		},
		func(p Person, opt O.Option[string]) Person {
			p.Name = opt
			return p
		},
	)
}

// TestAsOptional_BasicConversion tests basic conversion functionality
func TestAsOptional_BasicConversion(t *testing.T) {
	t.Run("converts prism to optional", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns Some for matching variant", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 100}
		value := optional.GetOption(result)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns None for non-matching variant", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}
		value := optional.GetOption(result)

		assert.True(t, O.IsNone(value))
	})
}

// TestAsOptional_Set tests the Set operation
func TestAsOptional_Set(t *testing.T) {
	t.Run("sets value when variant matches", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}
		updated := optional.Set(100)(result)

		assert.Equal(t, Success{Value: 100}, updated)
	})

	t.Run("Set is no-op when variant doesn't match", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}
		updated := optional.Set(100)(result)

		assert.Equal(t, result, updated)
	})

	t.Run("Set with zero value", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}
		updated := optional.Set(0)(result)

		assert.Equal(t, Success{Value: 0}, updated)
	})
}

// TestAsOptional_OptionalLaw1_GetSet tests Law 1 (No-op on None)
// Law: GetOption(s) = None => Set(a)(s) = s
func TestAsOptional_OptionalLaw1_GetSet(t *testing.T) {
	t.Run("Law 1: Set is no-op when GetOption returns None", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(optional.GetOption(result)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := optional.Set(100)(result)

		// Verify the result is unchanged
		assert.Equal(t, result, updated)
	})

	t.Run("Law 1: Set updates when GetOption returns Some", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(optional.GetOption(result)))

		// Set should update when GetOption returns Some
		updated := optional.Set(100)(result)

		assert.Equal(t, Success{Value: 100}, updated)
	})

	t.Run("Law 1: Multiple Set operations on None are all no-ops", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}

		// Multiple Set operations should all be no-ops
		updated := optional.Set(100)(optional.Set(50)(result))

		assert.Equal(t, result, updated)
	})
}

// TestAsOptional_OptionalLaw2_SetGet tests Law 2 (Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
func TestAsOptional_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("Law 2: Can retrieve what was set when Some exists", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}

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
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(optional.GetOption(result)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := optional.Set(100)(result)

		// Verify the result is unchanged
		assert.Equal(t, result, updated)
		assert.True(t, O.IsNone(optional.GetOption(updated)))
	})

	t.Run("Law 2: Works with zero values", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}

		// Set to zero value
		updated := optional.Set(0)(result)
		retrieved := optional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(retrieved))
	})

	t.Run("Law 2: Works with negative values", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}

		// Set to negative value
		updated := optional.Set(-10)(result)
		retrieved := optional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, -10, O.GetOrElse(F.Constant(0))(retrieved))
	})
}

// TestAsOptional_OptionalLaw3_SetSet tests Law 3 (Last Set Wins)
// Law: Set(b)(Set(a)(s)) = Set(b)(s)
func TestAsOptional_OptionalLaw3_SetSet(t *testing.T) {
	t.Run("Law 3: Last set wins when starting with Some", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 42}

		// Set twice
		setTwice := optional.Set(100)(optional.Set(50)(result))
		setOnce := optional.Set(100)(result)

		assert.Equal(t, setOnce, setTwice)
		assert.Equal(t, Success{Value: 100}, setTwice)
	})

	t.Run("Law 3: Both are no-op when starting with None", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Failure{Error: "failed"}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := optional.Set(100)(optional.Set(50)(result))
		setOnce := optional.Set(100)(result)

		// Both should equal the original result (no-op)
		assert.Equal(t, result, setTwice)
		assert.Equal(t, result, setOnce)
	})

	t.Run("Law 3: Works with multiple different values", func(t *testing.T) {
		prism := makeSuccessPrism()
		optional := AsOptional(prism)

		result := Success{Value: 1}

		// Set multiple times
		setMultiple := optional.Set(4)(optional.Set(3)(optional.Set(2)(result)))
		setOnce := optional.Set(4)(result)

		assert.Equal(t, setOnce, setMultiple)
		assert.Equal(t, Success{Value: 4}, setMultiple)
	})
}

// TestSome_BasicConversion tests basic conversion functionality
func TestSome_BasicConversion(t *testing.T) {
	t.Run("converts optional to focus on Some", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}
		value := valueOptional.GetOption(config)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns Some when inner Option is Some", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(100)}
		value := valueOptional.GetOption(config)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns None when inner Option is None", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}
		value := valueOptional.GetOption(config)

		assert.True(t, O.IsNone(value))
	})

	t.Run("works with string values", func(t *testing.T) {
		nameOptional := makeNameOptional()
		valueOptional := Some(nameOptional)

		person := Person{Name: O.Some("Alice")}
		value := valueOptional.GetOption(person)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(value))
	})
}

// TestSome_Set tests the Set operation
func TestSome_Set(t *testing.T) {
	t.Run("sets value when inner Option is Some", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}
		updated := valueOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})

	t.Run("Set is no-op when inner Option is None", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}
		updated := valueOptional.Set(60)(config)

		assert.Equal(t, config, updated)
		assert.Equal(t, O.None[int](), updated.Timeout)
	})

	t.Run("preserves other fields", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30), Retries: O.Some(3)}
		updated := valueOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
		assert.Equal(t, O.Some(3), updated.Retries)
	})

	t.Run("works with string values", func(t *testing.T) {
		nameOptional := makeNameOptional()
		valueOptional := Some(nameOptional)

		person := Person{Name: O.Some("Alice")}
		updated := valueOptional.Set("Bob")(person)

		assert.Equal(t, O.Some("Bob"), updated.Name)
	})
}

// TestSome_OptionalLaw1_GetSet tests Law 1 (No-op on None)
// Law: GetOption(s) = None => Set(a)(s) = s
func TestSome_OptionalLaw1_GetSet(t *testing.T) {
	t.Run("Law 1: Set is no-op when GetOption returns None", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(valueOptional.GetOption(config)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := valueOptional.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.Equal(t, O.None[int](), updated.Timeout)
	})

	t.Run("Law 1: Set updates when GetOption returns Some", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(valueOptional.GetOption(config)))

		// Set should update when GetOption returns Some
		updated := valueOptional.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})

	t.Run("Law 1: Multiple Set operations on None are all no-ops", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}

		// Multiple Set operations should all be no-ops
		updated := valueOptional.Set(60)(valueOptional.Set(30)(config))

		assert.Equal(t, config, updated)
		assert.Equal(t, O.None[int](), updated.Timeout)
	})
}

// TestSome_OptionalLaw2_SetGet tests Law 2 (Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
func TestSome_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("Law 2: Can retrieve what was set when Some exists", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(valueOptional.GetOption(config)))

		// Set a value and get it back
		newValue := 60
		updated := valueOptional.Set(newValue)(config)
		retrieved := valueOptional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, newValue, O.GetOrElse(F.Constant(0))(retrieved))
	})

	t.Run("Law 2: Set is no-op when starting from None", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(valueOptional.GetOption(config)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := valueOptional.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.True(t, O.IsNone(valueOptional.GetOption(updated)))
	})

	t.Run("Law 2: Works with zero values", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Set to zero value
		updated := valueOptional.Set(0)(config)
		retrieved := valueOptional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(retrieved))
	})

	t.Run("Law 2: Works with string values", func(t *testing.T) {
		nameOptional := makeNameOptional()
		valueOptional := Some(nameOptional)

		person := Person{Name: O.Some("Alice")}

		// Set a new value
		updated := valueOptional.Set("Bob")(person)
		retrieved := valueOptional.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, "Bob", O.GetOrElse(F.Constant(""))(retrieved))
	})
}

// TestSome_OptionalLaw3_SetSet tests Law 3 (Last Set Wins)
// Law: Set(b)(Set(a)(s)) = Set(b)(s)
func TestSome_OptionalLaw3_SetSet(t *testing.T) {
	t.Run("Law 3: Last set wins when starting with Some", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Set twice
		setTwice := valueOptional.Set(60)(valueOptional.Set(45)(config))
		setOnce := valueOptional.Set(60)(config)

		assert.Equal(t, setOnce.Timeout, setTwice.Timeout)
		assert.Equal(t, O.Some(60), setTwice.Timeout)
	})

	t.Run("Law 3: Both are no-op when starting with None", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.None[int]()}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := valueOptional.Set(60)(valueOptional.Set(45)(config))
		setOnce := valueOptional.Set(60)(config)

		// Both should equal the original config (no-op)
		assert.Equal(t, config, setTwice)
		assert.Equal(t, config, setOnce)
		assert.Equal(t, O.None[int](), setTwice.Timeout)
	})

	t.Run("Law 3: Works with multiple different values", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(10)}

		// Set multiple times
		setMultiple := valueOptional.Set(40)(valueOptional.Set(30)(valueOptional.Set(20)(config)))
		setOnce := valueOptional.Set(40)(config)

		assert.Equal(t, setOnce.Timeout, setMultiple.Timeout)
		assert.Equal(t, O.Some(40), setMultiple.Timeout)
	})

	t.Run("Law 3: Works with string values", func(t *testing.T) {
		nameOptional := makeNameOptional()
		valueOptional := Some(nameOptional)

		person := Person{Name: O.Some("Alice")}

		// Set twice
		setTwice := valueOptional.Set("Charlie")(valueOptional.Set("Bob")(person))
		setOnce := valueOptional.Set("Charlie")(person)

		assert.Equal(t, setOnce.Name, setTwice.Name)
		assert.Equal(t, O.Some("Charlie"), setTwice.Name)
	})
}

// TestSome_EdgeCases tests edge cases
func TestSome_EdgeCases(t *testing.T) {
	t.Run("works with zero values", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(0)}

		result := valueOptional.GetOption(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))

		updated := valueOptional.Set(0)(config)
		assert.Equal(t, O.Some(0), updated.Timeout)
	})

	t.Run("works with empty strings", func(t *testing.T) {
		nameOptional := makeNameOptional()
		valueOptional := Some(nameOptional)

		person := Person{Name: O.Some("")}

		result := valueOptional.GetOption(person)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "", O.GetOrElse(F.Constant("default"))(result))

		updated := valueOptional.Set("")(person)
		assert.Equal(t, O.Some(""), updated.Name)
	})

	t.Run("works with negative numbers", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(-10)}

		result := valueOptional.GetOption(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, -10, O.GetOrElse(F.Constant(0))(result))

		updated := valueOptional.Set(-20)(config)
		assert.Equal(t, O.Some(-20), updated.Timeout)
	})
}

// TestSome_Integration tests integration scenarios
func TestSome_Integration(t *testing.T) {
	t.Run("can be used with option operations", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Use with option map
		result := F.Pipe1(
			valueOptional.GetOption(config),
			O.Map(func(x int) int { return x * 2 }),
		)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("can be used with option chain", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		valueOptional := Some(timeoutOptional)

		config := Config{Timeout: O.Some(30)}

		// Use with option chain
		result := F.Pipe1(
			valueOptional.GetOption(config),
			O.Chain(func(x int) O.Option[int] {
				if x > 20 {
					return O.Some(x * 2)
				}
				return O.None[int]()
			}),
		)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 60, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("can chain updates on multiple fields", func(t *testing.T) {
		timeoutOptional := makeTimeoutOptional()
		timeoutValue := Some(timeoutOptional)

		config := Config{
			Timeout: O.Some(30),
			Retries: O.Some(3),
		}

		// Chain updates
		updated := timeoutValue.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
		assert.Equal(t, O.Some(3), updated.Retries)
	})
}

// TestPrismSome tests the PrismSome helper function
func TestPrismSome(t *testing.T) {
	t.Run("GetOption returns the Option itself", func(t *testing.T) {
		prism := PrismSome[int]()

		opt := O.Some(42)
		result := prism.GetOption(opt)

		assert.Equal(t, opt, result)
	})

	t.Run("GetOption works with None", func(t *testing.T) {
		prism := PrismSome[int]()

		opt := O.None[int]()
		result := prism.GetOption(opt)

		assert.Equal(t, opt, result)
	})

	t.Run("ReverseGet wraps in Some", func(t *testing.T) {
		prism := PrismSome[int]()

		wrapped := prism.ReverseGet(42)

		assert.Equal(t, O.Some(42), wrapped)
	})

	t.Run("works with string type", func(t *testing.T) {
		prism := PrismSome[string]()

		opt := O.Some("hello")
		result := prism.GetOption(opt)

		assert.Equal(t, opt, result)

		wrapped := prism.ReverseGet("world")
		assert.Equal(t, O.Some("world"), wrapped)
	})
}

// Made with Bob
