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
	C "github.com/IBM/fp-go/v2/internal/common"
	N "github.com/IBM/fp-go/v2/number"
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
func makeTimeoutOptional() C.Optional[Config, O.Option[int]] {
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
func makeNameOptional() C.Optional[Person, O.Option[string]] {
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
			O.Map(N.Mul(2)),
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
}

// TestCompose_BasicFunctionality tests basic composition behavior
func TestCompose_BasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when both optional and prism match", func(t *testing.T) {
		// Create an optional that focuses on an Option field
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		// Create a prism for Some values
		somePrism := PrismSome[int]()

		// Compose them
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}
		value := composed.GetOption(config)

		assert.True(t, O.IsSome(value))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(0))(value))
	})

	t.Run("GetOption returns None when optional doesn't match", func(t *testing.T) {
		// Create an optional that may not match
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				if c.Retries == O.None[int]() {
					return O.None[O.Option[int]]()
				}
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.None[int]()}
		value := composed.GetOption(config)

		assert.True(t, O.IsNone(value))
	})

	t.Run("GetOption returns None when prism doesn't match", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.None[int]()}
		value := composed.GetOption(config)

		assert.True(t, O.IsNone(value))
	})

	t.Run("Set updates value when both optional and prism match", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}
		updated := composed.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})

	t.Run("Set is no-op when optional doesn't match", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				if c.Retries == O.None[int]() {
					return O.None[O.Option[int]]()
				}
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.None[int]()}
		updated := composed.Set(60)(config)

		assert.Equal(t, config, updated)
	})

	t.Run("Set is no-op when prism doesn't match", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.None[int]()}
		updated := composed.Set(60)(config)

		assert.Equal(t, config, updated)
	})
}

// TestCompose_OptionalLaw1_GetSet tests Law 1 (No-op on None)
// Law: GetOption(s) = None => Set(b)(s) = s
func TestCompose_OptionalLaw1_GetSet(t *testing.T) {
	t.Run("Law 1: Set is no-op when GetOption returns None (optional doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				if c.Retries == O.None[int]() {
					return O.None[O.Option[int]]()
				}
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(composed.GetOption(config)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := composed.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
	})

	t.Run("Law 1: Set is no-op when GetOption returns None (prism doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(composed.GetOption(config)))

		// Law 1: Set must be a no-op when GetOption returns None
		updated := composed.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
	})

	t.Run("Law 1: Set updates when GetOption returns Some", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(composed.GetOption(config)))

		// Set should update when GetOption returns Some
		updated := composed.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
	})
}

// TestCompose_OptionalLaw2_SetGet tests Law 2 (Get what you Set)
// Law: GetOption(s) = Some(_) => GetOption(Set(b)(s)) = Some(b)
func TestCompose_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("Law 2: Can retrieve what was set when Some exists", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}

		// Verify GetOption returns Some
		assert.True(t, O.IsSome(composed.GetOption(config)))

		// Set a value and get it back
		newValue := 60
		updated := composed.Set(newValue)(config)
		retrieved := composed.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, newValue, O.GetOrElse(F.Constant(0))(retrieved))
	})

	t.Run("Law 2: Set is no-op when starting from None (optional doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				if c.Retries == O.None[int]() {
					return O.None[O.Option[int]]()
				}
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(composed.GetOption(config)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := composed.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.True(t, O.IsNone(composed.GetOption(updated)))
	})

	t.Run("Law 2: Set is no-op when starting from None (prism doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.None[int]()}

		// Verify GetOption returns None
		assert.True(t, O.IsNone(composed.GetOption(config)))

		// Set is a no-op when GetOption returns None (Law 1)
		updated := composed.Set(60)(config)

		// Verify the config is unchanged
		assert.Equal(t, config, updated)
		assert.True(t, O.IsNone(composed.GetOption(updated)))
	})

	t.Run("Law 2: Works with zero values", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}

		// Set to zero value
		updated := composed.Set(0)(config)
		retrieved := composed.GetOption(updated)

		assert.True(t, O.IsSome(retrieved))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(retrieved))
	})
}

// TestCompose_OptionalLaw3_SetSet tests Law 3 (Last Set Wins)
// Law: Set(c)(Set(b)(s)) = Set(c)(s)
func TestCompose_OptionalLaw3_SetSet(t *testing.T) {
	t.Run("Law 3: Last set wins when starting with Some", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30)}

		// Set twice
		setTwice := composed.Set(60)(composed.Set(45)(config))
		setOnce := composed.Set(60)(config)

		assert.Equal(t, setOnce.Timeout, setTwice.Timeout)
		assert.Equal(t, O.Some(60), setTwice.Timeout)
	})

	t.Run("Law 3: Both are no-op when starting with None (optional doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				if c.Retries == O.None[int]() {
					return O.None[O.Option[int]]()
				}
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.None[int]()}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := composed.Set(60)(composed.Set(45)(config))
		setOnce := composed.Set(60)(config)

		// Both should equal the original config (no-op)
		assert.Equal(t, config, setTwice)
		assert.Equal(t, config, setOnce)
	})

	t.Run("Law 3: Both are no-op when starting with None (prism doesn't match)", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.None[int]()}

		// Both Set operations are no-ops when GetOption returns None
		setTwice := composed.Set(60)(composed.Set(45)(config))
		setOnce := composed.Set(60)(config)

		// Both should equal the original config (no-op)
		assert.Equal(t, config, setTwice)
		assert.Equal(t, config, setOnce)
		assert.Equal(t, O.None[int](), setTwice.Timeout)
	})

	t.Run("Law 3: Works with multiple different values", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(10)}

		// Set multiple times
		setMultiple := composed.Set(40)(composed.Set(30)(composed.Set(20)(config)))
		setOnce := composed.Set(40)(config)

		assert.Equal(t, setOnce.Timeout, setMultiple.Timeout)
		assert.Equal(t, O.Some(40), setMultiple.Timeout)
	})
}

// TestCompose_EdgeCases tests edge cases
func TestCompose_EdgeCases(t *testing.T) {
	t.Run("Works with zero values", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(0)}

		result := composed.GetOption(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))

		updated := composed.Set(0)(config)
		assert.Equal(t, O.Some(0), updated.Timeout)
	})

	t.Run("Preserves other fields", func(t *testing.T) {
		configOptional := OPT.MakeOptional(
			func(c Config) O.Option[O.Option[int]] {
				return O.Some(c.Timeout)
			},
			func(c Config, opt O.Option[int]) Config {
				c.Timeout = opt
				return c
			},
		)

		somePrism := PrismSome[int]()
		composed := Compose[Config](somePrism)(configOptional)

		config := Config{Timeout: O.Some(30), Retries: O.Some(3)}
		updated := composed.Set(60)(config)

		assert.Equal(t, O.Some(60), updated.Timeout)
		assert.Equal(t, O.Some(3), updated.Retries)
	})
}

// TestPrismSome_TypeVariants tests PrismSome with different types
func TestPrismSome_TypeVariants(t *testing.T) {
	t.Run("works with string type", func(t *testing.T) {
		prism := PrismSome[string]()

		opt := O.Some("hello")
		result := prism.GetOption(opt)

		assert.Equal(t, opt, result)

		wrapped := prism.ReverseGet("world")
		assert.Equal(t, O.Some("world"), wrapped)
	})
}

// ---- AsOptional: equivalence and pipeline ----

// TestAsOptional_EquivalentToCommon verifies that AsOptional is a thin wrapper
// around the internal common.PrismAsOptional and produces identical behaviour.
func TestAsOptional_EquivalentToCommon(t *testing.T) {
	prismSA := makeSuccessPrism()

	method := AsOptional(prismSA)
	free := C.PrismAsOptional(prismSA)

	cases := []Result{Success{Value: 42}, Failure{Error: "oops"}}
	for _, r := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(r), method.GetOption(r))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(r), method.Set(99)(r))
		})
	}
}

// TestAsOptional_PipelineUsage verifies that AsOptional composes naturally
// with F.Pipe1 to chain into a longer pipeline.
func TestAsOptional_PipelineUsage(t *testing.T) {
	opt := AsOptional(makeSuccessPrism())

	// A GetOption → map pipeline that extracts the doubled value or -1.
	extract := func(r Result) int {
		return O.GetOrElse(F.Constant(-1))(
			O.Map(N.Mul(2))(opt.GetOption(r)),
		)
	}

	assert.Equal(t, 84, extract(Success{Value: 42}))
	assert.Equal(t, -1, extract(Failure{Error: "boom"}))
}

// ---- PrismSome: prism laws ----

// TestPrismSome_PrismLaw1_GetOptionReverseGet verifies the prism roundtrip law:
// GetOption(ReverseGet(a)) == Some(a).
func TestPrismSome_PrismLaw1_GetOptionReverseGet(t *testing.T) {
	p := PrismSome[int]()

	for _, v := range []int{0, 1, -99, 42} {
		t.Run("GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
			assert.Equal(t, O.Some(v), p.GetOption(p.ReverseGet(v)))
		})
	}
}

// TestPrismSome_PrismLaw2_GetSet verifies:
// if GetOption(s) == Some(a) then ReverseGet(a) == s (for Option types this means
// GetOption correctly reconstructs from itself).
func TestPrismSome_PrismLaw2_IdentityOnSome(t *testing.T) {
	p := PrismSome[string]()

	s := O.Some("hello")
	got := p.GetOption(s)
	assert.Equal(t, s, got, "GetOption of a Some should return that same Some")
}

// TestPrismSome_PipelineUsage verifies that PrismSome can be converted to an
// Optional and used in a F.Pipe1 pipeline.
func TestPrismSome_PipelineUsage(t *testing.T) {
	// Wrap PrismSome into an optional via AsOptional, then use Pipe1.
	opt := AsOptional(PrismSome[int]())

	result := F.Pipe1(O.Some(7), opt.Set(42))
	assert.Equal(t, O.Some(42), result)

	noOp := F.Pipe1(O.None[int](), opt.Set(42))
	assert.Equal(t, O.None[int](), noOp)
}

// ---- Some: equivalence and pipeline ----

// TestSome_EquivalentToCommon verifies that Some is a thin wrapper around
// common.OptionalSome and both produce identical results.
func TestSome_EquivalentToCommon(t *testing.T) {
	soa := makeTimeoutOptional()

	method := Some(soa)
	free := C.OptionalSome(soa)

	configs := []Config{
		{Timeout: O.Some(30), Retries: O.Some(3)},
		{Timeout: O.None[int](), Retries: O.Some(1)},
	}
	for _, c := range configs {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(c), method.GetOption(c))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(c), method.Set(99)(c))
		})
	}
}

// TestSome_PipelineUsage verifies that Some integrates with F.Pipe1.
func TestSome_PipelineUsage(t *testing.T) {
	opt := Some(makeTimeoutOptional())

	config := Config{Timeout: O.Some(10)}
	result := F.Pipe1(config, opt.Set(99))
	assert.Equal(t, O.Some(99), result.Timeout)
}

// TestSome_PreservesOtherFields verifies that Set through Some does not disturb
// fields that the optional does not touch.
func TestSome_PreservesOtherFields(t *testing.T) {
	opt := Some(makeTimeoutOptional())

	config := Config{Timeout: O.Some(5), Retries: O.Some(3)}
	updated := opt.Set(50)(config)

	assert.Equal(t, O.Some(50), updated.Timeout)
	assert.Equal(t, O.Some(3), updated.Retries)
}

// TestSome_MultipleTypes verifies that Some works with non-int type parameters.
func TestSome_MultipleTypes(t *testing.T) {
	nameOpt := Some(makeNameOptional())

	t.Run("string Some – GetOption match", func(t *testing.T) {
		p := Person{Name: O.Some("Alice"), Age: O.Some(30)}
		got := nameOpt.GetOption(p)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(got))
	})

	t.Run("string None – GetOption no match", func(t *testing.T) {
		p := Person{Name: O.None[string]()}
		assert.True(t, O.IsNone(nameOpt.GetOption(p)))
	})

	t.Run("string Set updates name", func(t *testing.T) {
		p := Person{Name: O.Some("Bob"), Age: O.Some(20)}
		updated := nameOpt.Set("Carol")(p)
		assert.Equal(t, O.Some("Carol"), updated.Name)
		assert.Equal(t, O.Some(20), updated.Age)
	})
}

// ---- Compose: equivalence, pipeline, sum-type prism ----

// makeConnectionOptional returns an Optional[Config, O.Option[int]] that always
// focuses on the Timeout field (reuses makeTimeoutOptional).
// For the Compose sum-type tests we also need an Optional that focuses on a
// non-Option intermediate type, so we define additional helpers below.

// TaggedValue is a simple sum type used to exercise Compose with a real
// discriminated union — as opposed to the Option-based examples in
// TestCompose_BasicFunctionality.
type TaggedValue interface{ isTaggedValue() }

type IntTagged struct{ V int }
type StrTagged struct{ S string }

func (IntTagged) isTaggedValue() {}
func (StrTagged) isTaggedValue() {}

type Payload struct {
	Tag TaggedValue
}

func makePayloadOptional() C.Optional[Payload, TaggedValue] {
	return OPT.MakeOptional(
		func(p Payload) O.Option[TaggedValue] {
			if p.Tag != nil {
				return O.Some(p.Tag)
			}
			return O.None[TaggedValue]()
		},
		func(p Payload, tv TaggedValue) Payload { p.Tag = tv; return p },
	)
}

func makeIntTaggedPrism() P.Prism[TaggedValue, int] {
	return P.MakePrism(
		func(tv TaggedValue) O.Option[int] {
			if it, ok := tv.(IntTagged); ok {
				return O.Some(it.V)
			}
			return O.None[int]()
		},
		func(v int) TaggedValue { return IntTagged{V: v} },
	)
}

// TestCompose_EquivalentToCommon verifies that Compose is a thin wrapper around
// common.OptionalComposePrism and produces identical behaviour.
func TestCompose_EquivalentToCommon(t *testing.T) {
	outer := makePayloadOptional()
	prismAB := makeIntTaggedPrism()

	method := F.Pipe1(outer, Compose[Payload](prismAB))
	free := C.OptionalComposePrism[Payload](prismAB)(outer)

	payloads := []Payload{
		{Tag: IntTagged{V: 7}},
		{Tag: StrTagged{S: "hi"}},
		{Tag: nil},
	}
	for _, p := range payloads {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(p), method.GetOption(p))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(p), method.Set(99)(p))
		})
	}
}

// TestCompose_SumTypePrism_GetOption tests Compose with a real discriminated
// union, verifying GetOption for each case.
func TestCompose_SumTypePrism_GetOption(t *testing.T) {
	composed := F.Pipe1(makePayloadOptional(), Compose[Payload](makeIntTaggedPrism()))

	t.Run("outer matches and prism matches", func(t *testing.T) {
		got := composed.GetOption(Payload{Tag: IntTagged{V: 42}})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(0))(got))
	})

	t.Run("outer matches but prism misses (StrTagged)", func(t *testing.T) {
		got := composed.GetOption(Payload{Tag: StrTagged{S: "hello"}})
		assert.True(t, O.IsNone(got))
	})

	t.Run("outer misses (nil tag)", func(t *testing.T) {
		got := composed.GetOption(Payload{Tag: nil})
		assert.True(t, O.IsNone(got))
	})
}

// TestCompose_SumTypePrism_Set tests Compose Set with a real discriminated union.
func TestCompose_SumTypePrism_Set(t *testing.T) {
	composed := F.Pipe1(makePayloadOptional(), Compose[Payload](makeIntTaggedPrism()))

	t.Run("updates to new IntTagged when both match", func(t *testing.T) {
		updated := composed.Set(99)(Payload{Tag: IntTagged{V: 1}})
		got := composed.GetOption(updated)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 99, O.GetOrElse(F.Constant(0))(got))
	})

	t.Run("no-op when prism misses (StrTagged)", func(t *testing.T) {
		original := Payload{Tag: StrTagged{S: "keep"}}
		result := composed.Set(99)(original)
		assert.Equal(t, original, result)
	})

	t.Run("no-op when outer optional misses (nil)", func(t *testing.T) {
		original := Payload{Tag: nil}
		result := composed.Set(99)(original)
		assert.Equal(t, original, result)
	})
}

// TestCompose_SumTypePrism_OptionalLaws verifies all three optional laws on the
// Compose result when a real sum-type prism is used.
func TestCompose_SumTypePrism_OptionalLaws(t *testing.T) {
	composed := F.Pipe1(makePayloadOptional(), Compose[Payload](makeIntTaggedPrism()))

	matching := Payload{Tag: IntTagged{V: 5}}
	nonMatchingOuter := Payload{Tag: nil}
	nonMatchingPrism := Payload{Tag: StrTagged{S: "str"}}

	t.Run("law GetSet: outer None => Set is no-op", func(t *testing.T) {
		assert.Equal(t, nonMatchingOuter, composed.Set(99)(nonMatchingOuter))
	})

	t.Run("law GetSet: prism miss => Set is no-op", func(t *testing.T) {
		assert.Equal(t, nonMatchingPrism, composed.Set(99)(nonMatchingPrism))
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := composed.Set(77)(matching)
		got := composed.GetOption(updated)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 77, O.GetOrElse(F.Constant(0))(got))
	})

	t.Run("law SetSet: Set(b)(Set(a)(s)) == Set(b)(s)", func(t *testing.T) {
		result1 := composed.Set(20)(composed.Set(10)(matching))
		result2 := composed.Set(20)(matching)
		assert.Equal(t, composed.GetOption(result2), composed.GetOption(result1))
	})
}

// TestCompose_PipelineUsage verifies that Compose integrates naturally with
// F.Pipe1 and produces the same result as direct application.
func TestCompose_PipelineUsage(t *testing.T) {
	outer := makePayloadOptional()
	prismAB := makeIntTaggedPrism()

	pipeResult := F.Pipe1(outer, Compose[Payload](prismAB))
	direct := Compose[Payload](prismAB)(outer)

	p := Payload{Tag: IntTagged{V: 3}}
	assert.Equal(t, direct.GetOption(p), pipeResult.GetOption(p))
	assert.Equal(t, direct.Set(99)(p), pipeResult.Set(99)(p))
}

// TestCompose_Chained verifies that two Compose calls can be chained to
// navigate three levels: Payload → TaggedValue → int → (via a second prism).
//
// The second prism only matches when the int is positive.
func TestCompose_Chained(t *testing.T) {
	positiveIntPrism := P.MakePrism(
		func(n int) O.Option[int] {
			if n > 0 {
				return O.Some(n)
			}
			return O.None[int]()
		},
		func(n int) int { return n },
	)

	// chain: Payload → TaggedValue (Optional) → int (Prism) → int (Prism)
	composed := F.Pipe1(
		F.Pipe1(makePayloadOptional(), Compose[Payload](makeIntTaggedPrism())),
		Compose[Payload](positiveIntPrism),
	)

	t.Run("all levels match (IntTagged positive)", func(t *testing.T) {
		got := composed.GetOption(Payload{Tag: IntTagged{V: 5}})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 5, O.GetOrElse(F.Constant(0))(got))
	})

	t.Run("second prism misses (IntTagged zero)", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(Payload{Tag: IntTagged{V: 0}})))
	})

	t.Run("first prism misses (StrTagged)", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(Payload{Tag: StrTagged{S: "hi"}})))
	})

	t.Run("outer misses (nil)", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(Payload{Tag: nil})))
	})

	t.Run("Set updates when all match", func(t *testing.T) {
		updated := composed.Set(99)(Payload{Tag: IntTagged{V: 5}})
		got := composed.GetOption(updated)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 99, O.GetOrElse(F.Constant(0))(got))
	})

	t.Run("Set is no-op when second prism misses", func(t *testing.T) {
		original := Payload{Tag: IntTagged{V: 0}}
		result := composed.Set(99)(original)
		assert.Equal(t, original, result)
	})
}
