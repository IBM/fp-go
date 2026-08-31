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

package common

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestMakePrism(t *testing.T) {
	somePrism := MakePrism(F.Identity[Option[int]], O.Some[int])

	assert.Equal(t, O.Some(1), somePrism.GetOption(O.Some(1)))
}

func TestPrismId(t *testing.T) {
	idPrism := PrismId[int]()

	// GetOption always returns Some for identity
	assert.Equal(t, O.Some(42), idPrism.GetOption(42))

	// ReverseGet is identity
	assert.Equal(t, 42, idPrism.ReverseGet(42))
}

func TestPrismFromPredicate(t *testing.T) {
	// Prism for positive numbers
	positivePrism := PrismFromPredicate(func(n int) bool {
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

func TestPrismComposePrism(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Prism for positive numbers
	positivePrism := PrismFromPredicate(func(n int) bool {
		return n > 0
	})

	// Compose: Option[int] -> int (if Some and positive)
	composedPrism := F.Pipe1(
		somePrism,
		PrismComposePrism[Option[int]](positivePrism),
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

func TestPrismSet(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Set value when it matches
	result := PrismSet[Option[int]](100)(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(100), result)

	// No change when it doesn't match
	result = PrismSet[Option[int]](100)(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

func TestPrismSome(t *testing.T) {
	// Prism that focuses on an Option field
	type Config struct {
		Timeout Option[int]
	}

	configPrism := MakePrism(
		func(c Config) Option[Option[int]] {
			return O.Some(c.Timeout)
		},
		func(t Option[int]) Config {
			return Config{Timeout: t}
		},
	)

	// Focus on the Some value
	somePrism := PrismSome(configPrism)

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

func TestPrismIMap(t *testing.T) {
	// Prism for Some values
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Map to string and back
	stringPrism := F.Pipe1(
		somePrism,
		PrismIMap[Option[int]](
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
		F.Identity[Option[int]],
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

func TestPrismSet_ViaModifyOption(t *testing.T) {
	// Test Set, which delegates to prismModifyOption
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Modify when match
	setFn := PrismSet[Option[int]](100)
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
		func(r testResult) Option[int] {
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
	setFn := PrismSet[testResult](200)
	updated := setFn(successPrism)(success)
	assert.Equal(t, testSuccess{Value: 200}, updated)

	// Test Set with Failure (no change)
	unchanged := setFn(successPrism)(failure)
	assert.Equal(t, failure, unchanged)
}

// TestPrismSet_Modify tests Set, which delegates to the internal prismModify function
func TestPrismSet_Modify(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Test modify with matching value
	result := PrismSet[Option[int]](84)(somePrism)(O.Some(42))
	assert.Equal(t, O.Some(84), result)

	// Test that original is returned when no match
	result = PrismSet[Option[int]](100)(somePrism)(O.None[int]())
	assert.Equal(t, O.None[int](), result)
}

// TestPrismSet_WithPredicate tests Set applied to a PrismFromPredicate prism
func TestPrismSet_WithPredicate(t *testing.T) {
	// Create a prism for positive numbers
	positivePrism := PrismFromPredicate(N.MoreThan(0))

	// Modify positive number
	setter := PrismSet[int](100)
	result := setter(positivePrism)(42)
	assert.Equal(t, 100, result)

	// Try to modify negative number (no change)
	result = setter(positivePrism)(-5)
	assert.Equal(t, -5, result)
}

// TestPrismAsTraversal tests converting a prism to a traversal
func TestPrismAsTraversal(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Simple identity functor for testing
	type Identity[A any] struct{ Value A }

	fof := func(s Option[int]) Identity[Option[int]] {
		return Identity[Option[int]]{Value: s}
	}

	fmap := func(f func(int) Option[int]) func(Identity[int]) Identity[Option[int]] {
		return func(ia Identity[int]) Identity[Option[int]] {
			return Identity[Option[int]]{Value: f(ia.Value)}
		}
	}

	type TraversalFunc func(func(int) Identity[int]) func(Option[int]) Identity[Option[int]]
	traversal := PrismAsTraversal[TraversalFunc](fof, fmap)(somePrism)

	// Test with Some value
	f := func(n int) Identity[int] {
		return Identity[int]{Value: n * 2}
	}
	result := traversal(f)(O.Some(21))
	assert.Equal(t, O.Some(42), result.Value)

	// Test with None value
	result = traversal(f)(O.None[int]())
	assert.Equal(t, O.None[int](), result.Value)
}

// Test types for composition chain test
type testOuter struct{ Middle Option[testInner] }
type testInner struct{ Value Option[int] }

// TestPrismCompositionChain tests composing multiple prisms
func TestPrismCompositionChain(t *testing.T) {
	// Three-level composition

	outerPrism := MakePrism(
		func(o testOuter) Option[Option[testInner]] {
			return O.Some(o.Middle)
		},
		func(m Option[testInner]) testOuter {
			return testOuter{Middle: m}
		},
	)

	middlePrism := MakePrism(
		F.Identity[Option[testInner]],
		O.Some[testInner],
	)

	innerPrism := MakePrism(
		func(i testInner) Option[Option[int]] {
			return O.Some(i.Value)
		},
		func(v Option[int]) testInner {
			return testInner{Value: v}
		},
	)

	// Compose all three
	composed := F.Pipe2(
		outerPrism,
		PrismComposePrism[testOuter](middlePrism),
		PrismComposePrism[testOuter](innerPrism),
	)

	// Further compose to get the int value
	finalPrism := PrismSome(composed)

	// Test extraction through all layers
	outer := testOuter{Middle: O.Some(testInner{Value: O.Some(42)})}
	value := finalPrism.GetOption(outer)
	assert.Equal(t, O.Some(42), value)

	// Test with None at middle layer
	outerNone := testOuter{Middle: O.None[testInner]()}
	value = finalPrism.GetOption(outerNone)
	assert.Equal(t, O.None[int](), value)

	// Test with None at inner layer
	outerInnerNone := testOuter{Middle: O.Some(testInner{Value: O.None[int]()})}
	value = finalPrism.GetOption(outerInnerNone)
	assert.Equal(t, O.None[int](), value)
}

// TestPrismSetMultipleTimes tests setting values multiple times
func TestPrismSetMultipleTimes(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Chain multiple sets
	result := F.Pipe3(
		O.Some(10),
		PrismSet[Option[int]](20)(somePrism),
		PrismSet[Option[int]](30)(somePrism),
		PrismSet[Option[int]](40)(somePrism),
	)

	assert.Equal(t, O.Some(40), result)
}

// TestPrismIMap_Bidirectional tests that PrismIMap maintains bidirectionality
func TestPrismIMap_Bidirectional(t *testing.T) {
	somePrism := MakePrism(
		F.Identity[Option[int]],
		O.Some[int],
	)

	// Map int to string and back
	stringPrism := F.Pipe1(
		somePrism,
		PrismIMap[Option[int]](
			func(n int) string {
				if n == 42 {
					return "forty-two"
				}
				return "other"
			},
			func(s string) int {
				if s == "forty-two" {
					return 42
				}
				return 0
			},
		),
	)

	// Test GetOption with mapping
	result := stringPrism.GetOption(O.Some(42))
	assert.Equal(t, O.Some("forty-two"), result)

	// Test ReverseGet with reverse mapping
	opt := stringPrism.ReverseGet("forty-two")
	assert.Equal(t, O.Some(42), opt)

	// Verify round-trip
	value := stringPrism.GetOption(stringPrism.ReverseGet("forty-two"))
	assert.Equal(t, O.Some("forty-two"), value)
}

// Test types for complex sum type test
type Shape interface{ isShape() }
type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ Base, Height float64 }

func (Circle) isShape()    {}
func (Rectangle) isShape() {}
func (Triangle) isShape()  {}

// TestPrismWithComplexSumType tests prism with a more complex sum type
func TestPrismWithComplexSumType(t *testing.T) {

	// Prism for Circle
	circlePrism := MakePrism(
		func(s Shape) Option[float64] {
			if c, ok := s.(Circle); ok {
				return O.Some(c.Radius)
			}
			return O.None[float64]()
		},
		func(r float64) Shape {
			return Circle{Radius: r}
		},
	)

	// Prism for Rectangle
	rectanglePrism := MakePrism(
		func(s Shape) Option[struct{ Width, Height float64 }] {
			if r, ok := s.(Rectangle); ok {
				return O.Some(struct{ Width, Height float64 }{r.Width, r.Height})
			}
			return O.None[struct{ Width, Height float64 }]()
		},
		func(dims struct{ Width, Height float64 }) Shape {
			return Rectangle{Width: dims.Width, Height: dims.Height}
		},
	)

	// Test Circle prism
	circle := Circle{Radius: 5.0}
	radius := circlePrism.GetOption(circle)
	assert.Equal(t, O.Some(5.0), radius)

	// Circle prism doesn't match Rectangle
	rect := Rectangle{Width: 10, Height: 20}
	radius = circlePrism.GetOption(rect)
	assert.Equal(t, O.None[float64](), radius)

	// Rectangle prism matches Rectangle
	dims := rectanglePrism.GetOption(rect)
	assert.True(t, O.IsSome(dims))

	// Test ReverseGet
	newCircle := circlePrism.ReverseGet(7.5)
	assert.Equal(t, Circle{Radius: 7.5}, newCircle)
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("prism with zero value", func(t *testing.T) {
		somePrism := MakePrism(
			F.Identity[Option[int]],
			O.Some[int],
		)

		// Zero value should work fine
		result := somePrism.GetOption(O.Some(0))
		assert.Equal(t, O.Some(0), result)

		opt := somePrism.ReverseGet(0)
		assert.Equal(t, O.Some(0), opt)
	})

	t.Run("prism with empty string", func(t *testing.T) {
		somePrism := MakePrism(
			F.Identity[Option[string]],
			O.Some[string],
		)

		result := somePrism.GetOption(O.Some(""))
		assert.Equal(t, O.Some(""), result)
	})

	t.Run("identity prism with nil pointer", func(t *testing.T) {
		type MyStruct struct{ Value int }
		idPrism := PrismId[*MyStruct]()

		var nilPtr *MyStruct
		result := idPrism.GetOption(nilPtr)
		assert.Equal(t, O.Some(nilPtr), result)
	})
}

// ---------------------------------------------------------------------------
// MakePrismWithName
// ---------------------------------------------------------------------------

func TestMakePrismWithName(t *testing.T) {
	p := MakePrismWithName(
		F.Identity[Option[int]],
		O.Some[int],
		"MyPrism",
	)

	t.Run("stores the assigned name", func(t *testing.T) {
		assert.Equal(t, "MyPrism", p.String())
	})

	t.Run("GetOption returns Some on a matching value", func(t *testing.T) {
		assert.Equal(t, O.Some(7), p.GetOption(O.Some(7)))
	})

	t.Run("GetOption returns None on a non-matching value", func(t *testing.T) {
		assert.Equal(t, O.None[int](), p.GetOption(O.None[int]()))
	})

	t.Run("ReverseGet constructs the source value", func(t *testing.T) {
		assert.Equal(t, O.Some(99), p.ReverseGet(99))
	})
}

// ---------------------------------------------------------------------------
// prismName formatting methods (String, Format, LogValue)
// ---------------------------------------------------------------------------

func TestPrismName_String(t *testing.T) {
	p := MakePrismWithName(F.Identity[Option[int]], O.Some[int], "Test.Prism")
	assert.Equal(t, "Test.Prism", p.String())
}

func TestPrismName_Format(t *testing.T) {
	p := MakePrismWithName(F.Identity[Option[int]], O.Some[int], "Test.Prism")

	t.Run("%s verb uses the prism name", func(t *testing.T) {
		assert.Equal(t, "Test.Prism", fmt.Sprintf("%s", p))
	})

	t.Run("%v verb uses the prism name", func(t *testing.T) {
		assert.Equal(t, "Test.Prism", fmt.Sprintf("%v", p))
	})

	t.Run("%q verb produces a quoted prism name", func(t *testing.T) {
		assert.Equal(t, `"Test.Prism"`, fmt.Sprintf("%q", p))
	})
}

func TestPrismName_LogValue(t *testing.T) {
	p := MakePrismWithName(F.Identity[Option[int]], O.Some[int], "Log.Prism")
	lv := p.LogValue()
	assert.Equal(t, "Log.Prism", lv.String())
}

// ---------------------------------------------------------------------------
// prismModifyOption — directly assert Option return values
// ---------------------------------------------------------------------------

func TestPrismModifyOption_Some(t *testing.T) {
	somePrism := MakePrism(F.Identity[Option[int]], O.Some[int])

	t.Run("returns Some(modified S) when prism matches", func(t *testing.T) {
		result := prismModifyOption(func(n int) int { return n * 2 }, somePrism, O.Some(21))
		assert.Equal(t, O.Some(O.Some(42)), result)
	})

	t.Run("returns None when prism does not match", func(t *testing.T) {
		result := prismModifyOption(func(n int) int { return n * 2 }, somePrism, O.None[int]())
		assert.Equal(t, O.None[Option[int]](), result)
	})
}

// ---------------------------------------------------------------------------
// prismSet (deprecated) — same semantics as Set
// ---------------------------------------------------------------------------

func TestPrismSet_Deprecated(t *testing.T) {
	somePrism := MakePrism(F.Identity[Option[int]], O.Some[int])

	t.Run("sets the value when prism matches", func(t *testing.T) {
		setter := prismSet[Option[int]](100)
		assert.Equal(t, O.Some(100), setter(somePrism)(O.Some(42)))
	})

	t.Run("returns original when prism does not match", func(t *testing.T) {
		setter := prismSet[Option[int]](100)
		assert.Equal(t, O.None[int](), setter(somePrism)(O.None[int]()))
	})

	t.Run("produces identical results to Set", func(t *testing.T) {
		via_set := PrismSet[Option[int]](55)(somePrism)(O.Some(1))
		via_prismSet := prismSet[Option[int]](55)(somePrism)(O.Some(1))
		assert.Equal(t, via_set, via_prismSet)
	})
}

// ---------------------------------------------------------------------------
// Id + Set interaction
// ---------------------------------------------------------------------------

func TestPrismId_Set(t *testing.T) {
	idPrism := PrismId[int]()

	t.Run("Set replaces the value because Id always matches", func(t *testing.T) {
		result := PrismSet[int](99)(idPrism)(42)
		assert.Equal(t, 99, result)
	})

	t.Run("Set with zero value", func(t *testing.T) {
		result := PrismSet[int](0)(idPrism)(42)
		assert.Equal(t, 0, result)
	})
}
