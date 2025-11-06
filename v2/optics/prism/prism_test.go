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
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	somePrism := MakePrism(F.Identity[O.Option[int]], O.Some[int])

	assert.Equal(t, O.Some(1), somePrism.GetOption(O.Some(1)))
}

func TestId(t *testing.T) {
	idPrism := Id[int]()

	// GetOption always returns Some for identity
	assert.Equal(t, O.Some(42), idPrism.GetOption(42))

	// ReverseGet is identity
	assert.Equal(t, 42, idPrism.ReverseGet(42))
}

func TestFromPredicate(t *testing.T) {
	// Prism for positive numbers
	positivePrism := FromPredicate(func(n int) bool {
		return n > 0
	})

	// Matches positive numbers
	assert.Equal(t, O.Some(42), positivePrism.GetOption(42))
	assert.Equal(t, O.Some(1), positivePrism.GetOption(1))

	// Doesn't match non-positive numbers
	assert.Equal(t, O.None[int](), positivePrism.GetOption(0))
	assert.Equal(t, O.None[int](), positivePrism.GetOption(-5))

	// ReverseGet always succeeds (doesn't check predicate)
	assert.Equal(t, 42, positivePrism.ReverseGet(42))
	assert.Equal(t, -5, positivePrism.ReverseGet(-5))
}

func TestCompose(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[O.Option[int]],
		O.Some[int],
	)

	// Prism for positive numbers
	positivePrism := FromPredicate(func(n int) bool {
		return n > 0
	})

	// Compose: Option[int] -> int (if Some and positive)
	composedPrism := F.Pipe1(
		somePrism,
		Compose[O.Option[int]](positivePrism),
	)

	// Test with Some positive
	assert.Equal(t, O.Some(42), composedPrism.GetOption(O.Some(42)))

	// Test with Some non-positive
	assert.Equal(t, O.None[int](), composedPrism.GetOption(O.Some(-5)))

	// Test with None
	assert.Equal(t, O.None[int](), composedPrism.GetOption(O.None[int]()))

	// ReverseGet constructs Some
	assert.Equal(t, O.Some(42), composedPrism.ReverseGet(42))
}

func TestSet(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[O.Option[int]],
		O.Some[int],
	)

	// Set value when it matches
	result := Set[O.Option[int], int](100)(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(100), result)

	// No change when it doesn't match
	result = Set[O.Option[int], int](100)(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

func TestSomeFunction(t *testing.T) {
	// Prism that focuses on an Option field
	type Config struct {
		Timeout O.Option[int]
	}

	configPrism := MakePrism(
		func(c Config) O.Option[O.Option[int]] {
			return O.Some(c.Timeout)
		},
		func(t O.Option[int]) Config {
			return Config{Timeout: t}
		},
	)

	// Focus on the Some value
	somePrism := Some(configPrism)

	// Extract from Some
	config := Config{Timeout: O.Some(30)}
	assert.Equal(t, O.Some(30), somePrism.GetOption(config))

	// Extract from None
	configNone := Config{Timeout: O.None[int]()}
	assert.Equal(t, O.None[int](), somePrism.GetOption(configNone))

	// ReverseGet constructs Config with Some
	result := somePrism.ReverseGet(60)
	assert.Equal(t, Config{Timeout: O.Some(60)}, result)
}

func TestIMap(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[O.Option[int]],
		O.Some[int],
	)

	// Map to string and back
	stringPrism := F.Pipe1(
		somePrism,
		IMap[O.Option[int]](
			func(n int) string {
				if n == 42 {
					return "42"
				}
				return "100"
			},
			func(s string) int {
				if s == "42" {
					return 42
				}
				return 100
			},
		),
	)

	// GetOption maps the value
	result := stringPrism.GetOption(O.Some(42))
	assert.Equal(t, O.Some("42"), result)

	// GetOption on None
	result = stringPrism.GetOption(O.None[int]())
	assert.Equal(t, O.None[string](), result)

	// ReverseGet maps back
	opt := stringPrism.ReverseGet("100")
	assert.Equal(t, O.Some(100), opt)
}

func TestPrismLaws(t *testing.T) {
	// Test prism laws with a simple prism
	somePrism := MakePrism(
		F.Identity[O.Option[int]],
		O.Some[int],
	)

	// Law 1: GetOptionReverseGet
	// prism.GetOption(prism.ReverseGet(a)) == Some(a)
	a := 42
	result := somePrism.GetOption(somePrism.ReverseGet(a))
	assert.Equal(t, O.Some(a), result)

	// Law 2: ReverseGetGetOption
	// if GetOption(s) == Some(a), then ReverseGet(a) should produce equivalent s
	s := O.Some(42)
	extracted := somePrism.GetOption(s)
	if O.IsSome(extracted) {
		reconstructed := somePrism.ReverseGet(O.GetOrElse(F.Constant(0))(extracted))
		assert.Equal(t, s, reconstructed)
	}
}

func TestPrismModifyOption(t *testing.T) {
	// Test the internal prismModifyOption function through Set
	somePrism := MakePrism(
		F.Identity[O.Option[int]],
		O.Some[int],
	)

	// Modify when match
	setFn := Set[O.Option[int], int](100)
	result := setFn(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(100), result)

	// No modification when no match
	result = setFn(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

// Custom sum type for testing
type testResult interface{ isResult() }
type testSuccess struct{ Value int }
type testFailure struct{ Error string }

func (testSuccess) isResult() {}
func (testFailure) isResult() {}

func TestPrismWithCustomType(t *testing.T) {
	// Create prism for Success variant
	successPrism := MakePrism(
		func(r testResult) O.Option[int] {
			if s, ok := r.(testSuccess); ok {
				return O.Some(s.Value)
			}
			return O.None[int]()
		},
		func(v int) testResult {
			return testSuccess{Value: v}
		},
	)

	// Test GetOption with Success
	success := testSuccess{Value: 42}
	assert.Equal(t, O.Some(42), successPrism.GetOption(success))

	// Test GetOption with Failure
	failure := testFailure{Error: "oops"}
	assert.Equal(t, O.None[int](), successPrism.GetOption(failure))

	// Test ReverseGet
	result := successPrism.ReverseGet(100)
	assert.Equal(t, testSuccess{Value: 100}, result)

	// Test Set with Success
	setFn := Set[testResult, int](200)
	updated := setFn(successPrism)(success)
	assert.Equal(t, testSuccess{Value: 200}, updated)

	// Test Set with Failure (no change)
	unchanged := setFn(successPrism)(failure)
	assert.Equal(t, failure, unchanged)
}
