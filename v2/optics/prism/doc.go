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

/*
Package prism provides prisms - optics for focusing on variants within sum types.

# Overview

A Prism is an optic used to select a specific variant from a sum type (also known as
tagged unions, discriminated unions, or algebraic data types). Unlike lenses which
focus on fields that always exist, prisms focus on values that may or may not be
present depending on which variant is active.

Prisms are essential for working with:
  - Either types (Left/Right)
  - Option types (Some/None)
  - Result types (Success/Failure)
  - Custom sum types
  - Error handling patterns

# Mathematical Foundation

A Prism[S, A] consists of two operations:
  - GetOption: S → Option[A] (try to extract A from S, may fail)
  - ReverseGet: A → S (construct S from A, always succeeds)

Prisms must satisfy the prism laws:
 1. GetOptionReverseGet: prism.GetOption(prism.ReverseGet(a)) == Some(a)
 2. ReverseGetGetOption: if GetOption(s) == Some(a), then ReverseGet(a) produces equivalent s

These laws ensure that:
  - Constructing and then extracting always succeeds
  - Extracting and then constructing preserves the value

# Basic Usage

Creating a prism for an Either type:

	type Result interface{ isResult() }
	type Success struct{ Value int }
	type Failure struct{ Error string }

	func (Success) isResult() {}
	func (Failure) isResult() {}

	successPrism := prism.MakePrism(
		func(r Result) option.Option[int] {
			if s, ok := r.(Success); ok {
				return option.Some(s.Value)
			}
			return option.None[int]()
		},
		func(v int) Result {
			return Success{Value: v}
		},
	)

	// Try to extract value from Success
	result := Success{Value: 42}
	value := successPrism.GetOption(result) // Some(42)

	// Try to extract value from Failure
	failure := Failure{Error: "oops"}
	value = successPrism.GetOption(failure) // None[int]

	// Construct a Success from a value
	newResult := successPrism.ReverseGet(100) // Success{Value: 100}

# Identity Prism

The identity prism focuses on the entire value:

	idPrism := prism.Id[int]()

	value := idPrism.GetOption(42)    // Some(42)
	result := idPrism.ReverseGet(42)  // 42

# From Predicate

Create a prism that matches values satisfying a predicate:

	positivePrism := prism.FromPredicate(func(n int) bool {
		return n > 0
	})

	// Matches positive numbers
	value := positivePrism.GetOption(42)  // Some(42)
	value = positivePrism.GetOption(-5)   // None[int]

	// ReverseGet always succeeds (doesn't check predicate)
	result := positivePrism.ReverseGet(42) // 42

# Composing Prisms

Prisms can be composed to focus on nested sum types:

	type Outer interface{ isOuter() }
	type OuterA struct{ Inner Inner }
	type OuterB struct{ Value string }

	type Inner interface{ isInner() }
	type InnerX struct{ Data int }
	type InnerY struct{ Info string }

	outerAPrism := prism.MakePrism(
		func(o Outer) option.Option[Inner] {
			if a, ok := o.(OuterA); ok {
				return option.Some(a.Inner)
			}
			return option.None[Inner]()
		},
		func(i Inner) Outer { return OuterA{Inner: i} },
	)

	innerXPrism := prism.MakePrism(
		func(i Inner) option.Option[int] {
			if x, ok := i.(InnerX); ok {
				return option.Some(x.Data)
			}
			return option.None[int]()
		},
		func(d int) Inner { return InnerX{Data: d} },
	)

	// Compose to access InnerX data from Outer
	composedPrism := F.Pipe1(
		outerAPrism,
		prism.Compose[Outer](innerXPrism),
	)

	outer := OuterA{Inner: InnerX{Data: 42}}
	data := composedPrism.GetOption(outer) // Some(42)

# Modifying Through Prisms

Apply transformations when the variant matches:

	type Status interface{ isStatus() }
	type Active struct{ Count int }
	type Inactive struct{}

	activePrism := prism.MakePrism(
		func(s Status) option.Option[int] {
			if a, ok := s.(Active); ok {
				return option.Some(a.Count)
			}
			return option.None[int]()
		},
		func(count int) Status { return Active{Count: count} },
	)

	status := Active{Count: 5}

	// Increment count if active
	updated := prism.Set(10)(activePrism)(status)
	// Result: Active{Count: 10}

	// Try to modify inactive status (no change)
	inactive := Inactive{}
	unchanged := prism.Set(10)(activePrism)(inactive)
	// Result: Inactive{} (unchanged)

# Working with Option Types

The Some function creates a prism focused on the Some variant of an Option:

	type Config struct {
		Timeout option.Option[int]
	}

	timeoutPrism := prism.MakePrism(
		func(c Config) option.Option[option.Option[int]] {
			return option.Some(c.Timeout)
		},
		func(t option.Option[int]) Config {
			return Config{Timeout: t}
		},
	)

	// Focus on the Some value
	somePrism := prism.Some(timeoutPrism)

	config := Config{Timeout: option.Some(30)}
	timeout := somePrism.GetOption(config) // Some(30)

	configNone := Config{Timeout: option.None[int]()}
	timeout = somePrism.GetOption(configNone) // None[int]

# Bidirectional Mapping

Transform the focus type of a prism:

	type Message interface{ isMessage() }
	type TextMessage struct{ Content string }
	type ImageMessage struct{ URL string }

	textPrism := prism.MakePrism(
		func(m Message) option.Option[string] {
			if t, ok := m.(TextMessage); ok {
				return option.Some(t.Content)
			}
			return option.None[string]()
		},
		func(content string) Message {
			return TextMessage{Content: content}
		},
	)

	// Map to uppercase
	upperPrism := F.Pipe1(
		textPrism,
		prism.IMap(
			strings.ToUpper,
			strings.ToLower,
		),
	)

	msg := TextMessage{Content: "hello"}
	upper := upperPrism.GetOption(msg) // Some("HELLO")

# Real-World Example: Error Handling

	type Result[T any] interface{ isResult() }

	type Success[T any] struct {
		Value T
	}

	type Failure[T any] struct {
		Error error
	}

	func (Success[T]) isResult() {}
	func (Failure[T]) isResult() {}

	func SuccessPrism[T any]() prism.Prism[Result[T], T] {
		return prism.MakePrism(
			func(r Result[T]) option.Option[T] {
				if s, ok := r.(Success[T]); ok {
					return option.Some(s.Value)
				}
				return option.None[T]()
			},
			func(v T) Result[T] {
				return Success[T]{Value: v}
			},
		)
	}

	func FailurePrism[T any]() prism.Prism[Result[T], error] {
		return prism.MakePrism(
			func(r Result[T]) option.Option[error] {
				if f, ok := r.(Failure[T]); ok {
					return option.Some(f.Error)
				}
				return option.None[error]()
			},
			func(e error) Result[T] {
				return Failure[T]{Error: e}
			},
		)
	}

	// Use the prisms
	result := Success[int]{Value: 42}

	successPrism := SuccessPrism[int]()
	value := successPrism.GetOption(result) // Some(42)

	failurePrism := FailurePrism[int]()
	err := failurePrism.GetOption(result) // None[error]

	// Transform success values
	doubled := prism.Set(84)(successPrism)(result)
	// Result: Success{Value: 84}

# Real-World Example: JSON Parsing

	type JSONValue interface{ isJSON() }
	type JSONString struct{ Value string }
	type JSONNumber struct{ Value float64 }
	type JSONBool struct{ Value bool }
	type JSONNull struct{}

	stringPrism := prism.MakePrism(
		func(j JSONValue) option.Option[string] {
			if s, ok := j.(JSONString); ok {
				return option.Some(s.Value)
			}
			return option.None[string]()
		},
		func(s string) JSONValue {
			return JSONString{Value: s}
		},
	)

	numberPrism := prism.MakePrism(
		func(j JSONValue) option.Option[float64] {
			if n, ok := j.(JSONNumber); ok {
				return option.Some(n.Value)
			}
			return option.None[float64]()
		},
		func(n float64) JSONValue {
			return JSONNumber{Value: n}
		},
	)

	// Parse and extract values
	jsonValue := JSONString{Value: "hello"}
	str := stringPrism.GetOption(jsonValue) // Some("hello")
	num := numberPrism.GetOption(jsonValue) // None[float64]

# Real-World Example: State Machine

	type State interface{ isState() }
	type Idle struct{}
	type Running struct{ Progress int }
	type Completed struct{ Result string }
	type Failed struct{ Error error }

	runningPrism := prism.MakePrism(
		func(s State) option.Option[int] {
			if r, ok := s.(Running); ok {
				return option.Some(r.Progress)
			}
			return option.None[int]()
		},
		func(progress int) State {
			return Running{Progress: progress}
		},
	)

	// Update progress if running
	state := Running{Progress: 50}
	updated := prism.Set(75)(runningPrism)(state)
	// Result: Running{Progress: 75}

	// Try to update when not running (no change)
	idle := Idle{}
	unchanged := prism.Set(75)(runningPrism)(idle)
	// Result: Idle{} (unchanged)

# Prisms vs Lenses

While both are optics, they serve different purposes:

**Prisms:**
  - Focus on variants of sum types
  - GetOption may fail (returns Option)
  - ReverseGet always succeeds
  - Used for pattern matching
  - Example: Either[Error, Value] → Value

**Lenses:**
  - Focus on fields of product types
  - Get always succeeds
  - Set always succeeds
  - Used for field access
  - Example: Person → Name

# Prism Laws

Prisms must satisfy two laws:

**Law 1: GetOptionReverseGet**

	prism.GetOption(prism.ReverseGet(a)) == Some(a)

Constructing a value and then extracting it always succeeds.

**Law 2: ReverseGetGetOption**

	if prism.GetOption(s) == Some(a)
	then prism.ReverseGet(a) should produce equivalent s

If extraction succeeds, reconstructing should produce an equivalent value.

# Performance Considerations

Prisms are efficient:
  - No reflection - uses type assertions
  - Minimal allocations
  - Composition creates function closures
  - GetOption short-circuits on mismatch

For best performance:
  - Cache composed prisms
  - Use type switches for multiple prisms
  - Consider batch operations when possible

# Type Safety

Prisms are fully type-safe:
  - Compile-time type checking
  - Type assertions are explicit
  - Generic type parameters ensure correctness
  - Composition maintains type relationships

# Built-in Prisms

The package provides many useful prisms for common transformations:

**Type Conversion & Parsing:**
  - FromEncoding(enc): Base64 encoding/decoding - Prism[string, []byte]
  - ParseURL(): URL parsing/formatting - Prism[string, *url.URL]
  - ParseDate(layout): Date parsing with custom layouts - Prism[string, time.Time]
  - ParseInt(): Integer string parsing - Prism[string, int]
  - ParseInt64(): 64-bit integer parsing - Prism[string, int64]
  - ParseBool(): Boolean string parsing - Prism[string, bool]
  - ParseFloat32(): 32-bit float parsing - Prism[string, float32]
  - ParseFloat64(): 64-bit float parsing - Prism[string, float64]

**Type Assertion & Extraction:**
  - InstanceOf[T](): Safe type assertion from any - Prism[any, T]
  - Deref[T](): Safe pointer dereferencing (filters nil) - Prism[*T, *T]

**Container/Wrapper Prisms:**
  - FromEither[E, T](): Extract Right values - Prism[Either[E, T], T]
    ReverseGet wraps into Right (acts as success constructor)
  - FromResult[T](): Extract success from Result - Prism[Result[T], T]
    ReverseGet wraps into success Result
  - FromOption[T](): Extract Some values - Prism[Option[T], T]
    ReverseGet wraps into Some (acts as Some constructor)

**Validation Prisms:**
  - FromZero[T](): Match only zero/default values - Prism[T, T]
  - FromNonZero[T](): Match only non-zero values - Prism[T, T]

**Pattern Matching:**
  - RegexMatcher(re): Extract regex matches with groups - Prism[string, Match]
  - RegexNamedMatcher(re): Extract named regex groups - Prism[string, NamedMatch]

Example using built-in prisms:

	// Parse and validate an integer from a string
	intPrism := prism.ParseInt()
	value := intPrism.GetOption("42")  // Some(42)
	invalid := intPrism.GetOption("abc")  // None[int]()

	// Extract success values from Either
	resultPrism := prism.FromEither[error, int]()
	success := either.Right[error](100)
	value = resultPrism.GetOption(success)  // Some(100)

	// ReverseGet acts as a constructor
	wrapped := resultPrism.ReverseGet(42)  // Right(42)

	// Compose prisms for complex transformations
	// Parse string to int, then wrap in Option
	composed := F.Pipe1(
		prism.ParseInt(),
		prism.Compose[string](prism.FromOption[int]()),
	)

# Function Reference

Core Functions:
  - MakePrism: Create a prism from GetOption and ReverseGet functions
  - Id: Create an identity prism
  - FromPredicate: Create a prism from a predicate function
  - Compose: Compose two prisms

Transformation:
  - Set: Set a value through a prism (no-op if variant doesn't match)
  - IMap: Bidirectionally map a prism

Specialized:
  - Some: Focus on the Some variant of an Option

# Related Packages

  - github.com/IBM/fp-go/v2/optics/lens: Lenses for product types
  - github.com/IBM/fp-go/v2/optics/iso: Isomorphisms
  - github.com/IBM/fp-go/v2/optics/optional: Optional optics
  - github.com/IBM/fp-go/v2/option: Optional values
  - github.com/IBM/fp-go/v2/either: Sum types
  - github.com/IBM/fp-go/v2/function: Function composition
*/
package prism
