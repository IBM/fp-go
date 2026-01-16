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

package option

import (
	"encoding/json"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	SampleData struct {
		Value    string
		OptValue Option[string]
	}
)

func TestJson(t *testing.T) {

	sample := SampleData{
		Value:    "value",
		OptValue: Of("optValue"),
	}

	data, err := json.Marshal(&sample)
	require.NoError(t, err)

	var deser SampleData
	err = json.Unmarshal(data, &deser)
	require.NoError(t, err)

	assert.Equal(t, sample, deser)

	sample = SampleData{
		Value:    "value",
		OptValue: None[string](),
	}

	data, err = json.Marshal(&sample)
	require.NoError(t, err)

	err = json.Unmarshal(data, &deser)
	require.NoError(t, err)

	assert.Equal(t, sample, deser)
}

func TestDefault(t *testing.T) {
	var e Option[string]

	assert.Equal(t, None[string](), e)
}

func TestReduce(t *testing.T) {

	assert.Equal(t, 2, F.Pipe1(None[int](), Reduce(utils.Sum, 2)))
	assert.Equal(t, 5, F.Pipe1(Some(3), Reduce(utils.Sum, 2)))
}

func TestIsNone(t *testing.T) {
	assert.True(t, IsNone(None[int]()))
	assert.False(t, IsNone(Of(1)))
}

func TestIsSome(t *testing.T) {
	assert.True(t, IsSome(Of(1)))
	assert.False(t, IsSome(None[int]()))
}

func TestMapOption(t *testing.T) {

	assert.Equal(t, F.Pipe1(Some(2), Map(utils.Double)), Some(4))

	assert.Equal(t, F.Pipe1(None[int](), Map(utils.Double)), None[int]())
}

func TestTryCachOption(t *testing.T) {

	res := TryCatch(utils.Error)

	assert.Equal(t, None[int](), res)
}

func TestAp(t *testing.T) {
	assert.Equal(t, Some(4), F.Pipe1(
		Some(utils.Double),
		Ap[int](Some(2)),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		Some(utils.Double),
		Ap[int](None[int]()),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[func(int) int](),
		Ap[int](Some(2)),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[func(int) int](),
		Ap[int](None[int]()),
	))
}

func TestChain(t *testing.T) {
	f := func(n int) Option[int] { return Some(n * 2) }
	g := func(_ int) Option[int] { return None[int]() }

	assert.Equal(t, Some(2), F.Pipe1(
		Some(1),
		Chain(f),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[int](),
		Chain(f),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		Some(1),
		Chain(g),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[int](),
		Chain(g),
	))
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, Of(1), F.Pipe1(Of(Of(1)), Flatten[int]))
}

func TestFold(t *testing.T) {
	f := F.Constant("none")
	g := func(s string) string { return fmt.Sprintf("some%d", len(s)) }

	fold := Fold(f, g)

	assert.Equal(t, "none", fold(None[string]()))
	assert.Equal(t, "some3", fold(Some("abc")))
}

func TestFromPredicate(t *testing.T) {
	p := func(n int) bool { return n > 2 }
	f := FromPredicate(p)

	assert.Equal(t, None[int](), f(1))
	assert.Equal(t, Some(3), f(3))
}

func TestAlt(t *testing.T) {
	assert.Equal(t, Some(1), F.Pipe1(Some(1), Alt(F.Constant(Some(2)))))
	assert.Equal(t, Some(2), F.Pipe1(Some(2), Alt(F.Constant(None[int]()))))
	assert.Equal(t, Some(1), F.Pipe1(None[int](), Alt(F.Constant(Some(1)))))
	assert.Equal(t, None[int](), F.Pipe1(None[int](), Alt(F.Constant(None[int]()))))
}

// TestZeroWithIntegers tests Zero function with integer types
func TestZeroWithIntegers(t *testing.T) {
	o := Zero[int]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithStrings tests Zero function with string types
func TestZeroWithStrings(t *testing.T) {
	o := Zero[string]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithBooleans tests Zero function with boolean types
func TestZeroWithBooleans(t *testing.T) {
	o := Zero[bool]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithFloats tests Zero function with float types
func TestZeroWithFloats(t *testing.T) {
	o := Zero[float64]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithPointers tests Zero function with pointer types
func TestZeroWithPointers(t *testing.T) {
	o := Zero[*int]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithSlices tests Zero function with slice types
func TestZeroWithSlices(t *testing.T) {
	o := Zero[[]int]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithMaps tests Zero function with map types
func TestZeroWithMaps(t *testing.T) {
	o := Zero[map[string]int]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithStructs tests Zero function with struct types
func TestZeroWithStructs(t *testing.T) {
	type TestStruct struct {
		Field1 int
		Field2 string
	}

	o := Zero[TestStruct]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithInterfaces tests Zero function with interface types
func TestZeroWithInterfaces(t *testing.T) {
	o := Zero[interface{}]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroIsNotSomeWithZeroValue tests that Zero is different from Some(zero value)
func TestZeroIsNotSomeWithZeroValue(t *testing.T) {
	// Zero returns None
	zero := Zero[int]()
	assert.True(t, IsNone(zero), "Zero should be None")

	// Some with zero value is different
	someZero := Some(0)
	assert.True(t, IsSome(someZero), "Some(0) should be Some")

	// They are not equal
	assert.NotEqual(t, zero, someZero, "Zero (None) should not equal Some(0)")
}

// TestZeroCanBeUsedWithOtherFunctions tests that Zero Options work with other option functions
func TestZeroCanBeUsedWithOtherFunctions(t *testing.T) {
	o := Zero[int]()

	// Test with Map - should remain None
	mapped := MonadMap(o, func(n int) string {
		return fmt.Sprintf("%d", n)
	})
	assert.True(t, IsNone(mapped), "Mapped Zero should still be None")

	// Test with Chain - should remain None
	chained := MonadChain(o, func(n int) Option[string] {
		return Some(fmt.Sprintf("value: %d", n))
	})
	assert.True(t, IsNone(chained), "Chained Zero should still be None")

	// Test with Fold - should use onNone branch
	folded := MonadFold(o,
		func() string { return "none" },
		func(n int) string { return fmt.Sprintf("some: %d", n) },
	)
	assert.Equal(t, "none", folded, "Folded Zero should use onNone branch")

	// Test with GetOrElse
	value := GetOrElse(func() int { return 42 })(o)
	assert.Equal(t, 42, value, "GetOrElse on Zero should return default value")
}

// TestZeroEquality tests that multiple Zero calls produce equal Options
func TestZeroEquality(t *testing.T) {
	o1 := Zero[int]()
	o2 := Zero[int]()

	assert.Equal(t, IsNone(o1), IsNone(o2), "Both should be None")
	assert.Equal(t, IsSome(o1), IsSome(o2), "Both should not be Some")
	assert.Equal(t, o1, o2, "Zero values should be equal")
}

// TestZeroWithComplexTypes tests Zero with more complex nested types
func TestZeroWithComplexTypes(t *testing.T) {
	type ComplexType struct {
		Nested map[string][]int
		Ptr    *string
	}

	o := Zero[ComplexType]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroWithNestedOption tests Zero with nested Option type
func TestZeroWithNestedOption(t *testing.T) {
	o := Zero[Option[int]]()

	assert.True(t, IsNone(o), "Zero should create a None value")
	assert.False(t, IsSome(o), "Zero should not create a Some value")
}

// TestZeroIsAlwaysNone tests that Zero never creates a Some value
func TestZeroIsAlwaysNone(t *testing.T) {
	// Test with various types
	o1 := Zero[int]()
	o2 := Zero[string]()
	o3 := Zero[bool]()
	o4 := Zero[*int]()
	o5 := Zero[[]string]()

	assert.True(t, IsNone(o1), "Zero should always be None")
	assert.True(t, IsNone(o2), "Zero should always be None")
	assert.True(t, IsNone(o3), "Zero should always be None")
	assert.True(t, IsNone(o4), "Zero should always be None")
	assert.True(t, IsNone(o5), "Zero should always be None")

	assert.False(t, IsSome(o1), "Zero should never be Some")
	assert.False(t, IsSome(o2), "Zero should never be Some")
	assert.False(t, IsSome(o3), "Zero should never be Some")
	assert.False(t, IsSome(o4), "Zero should never be Some")
	assert.False(t, IsSome(o5), "Zero should never be Some")
}

// TestZeroEqualsNone tests that Zero is equivalent to None
func TestZeroEqualsNone(t *testing.T) {
	zero := Zero[int]()
	none := None[int]()

	assert.Equal(t, zero, none, "Zero should be equal to None")
	assert.Equal(t, IsNone(zero), IsNone(none), "Both should be None")
	assert.Equal(t, IsSome(zero), IsSome(none), "Both should not be Some")
}

// TestZeroEqualsDefaultInitialization tests that Zero returns the same value as default initialization
func TestZeroEqualsDefaultInitialization(t *testing.T) {
	// Default initialization of Option
	var defaultInit Option[int]

	// Zero function
	zero := Zero[int]()

	// They should be equal
	assert.Equal(t, defaultInit, zero, "Zero should equal default initialization")
	assert.Equal(t, IsNone(defaultInit), IsNone(zero), "Both should be None")
	assert.Equal(t, IsSome(defaultInit), IsSome(zero), "Both should not be Some")
}
