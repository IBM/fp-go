// Copyright (c) 2024 - 2025 IBM Corp.
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

package pair

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestZeroWithIntegers tests Zero function with integer types
func TestZeroWithIntegers(t *testing.T) {
	p := Zero[int, int]()

	assert.Equal(t, 0, Head(p), "Head should be zero value for int")
	assert.Equal(t, 0, Tail(p), "Tail should be zero value for int")
}

// TestZeroWithStrings tests Zero function with string types
func TestZeroWithStrings(t *testing.T) {
	p := Zero[string, string]()

	assert.Equal(t, "", Head(p), "Head should be zero value for string")
	assert.Equal(t, "", Tail(p), "Tail should be zero value for string")
}

// TestZeroWithMixedTypes tests Zero function with different types
func TestZeroWithMixedTypes(t *testing.T) {
	p := Zero[string, int]()

	assert.Equal(t, "", Head(p), "Head should be zero value for string")
	assert.Equal(t, 0, Tail(p), "Tail should be zero value for int")
}

// TestZeroWithBooleans tests Zero function with boolean types
func TestZeroWithBooleans(t *testing.T) {
	p := Zero[bool, bool]()

	assert.Equal(t, false, Head(p), "Head should be zero value for bool")
	assert.Equal(t, false, Tail(p), "Tail should be zero value for bool")
}

// TestZeroWithFloats tests Zero function with float types
func TestZeroWithFloats(t *testing.T) {
	p := Zero[float64, float32]()

	assert.Equal(t, 0.0, Head(p), "Head should be zero value for float64")
	assert.Equal(t, float32(0.0), Tail(p), "Tail should be zero value for float32")
}

// TestZeroWithPointers tests Zero function with pointer types
func TestZeroWithPointers(t *testing.T) {
	p := Zero[*int, *string]()

	assert.Nil(t, Head(p), "Head should be nil for pointer type")
	assert.Nil(t, Tail(p), "Tail should be nil for pointer type")
}

// TestZeroWithSlices tests Zero function with slice types
func TestZeroWithSlices(t *testing.T) {
	p := Zero[[]int, []string]()

	assert.Nil(t, Head(p), "Head should be nil for slice type")
	assert.Nil(t, Tail(p), "Tail should be nil for slice type")
}

// TestZeroWithMaps tests Zero function with map types
func TestZeroWithMaps(t *testing.T) {
	p := Zero[map[string]int, map[int]string]()

	assert.Nil(t, Head(p), "Head should be nil for map type")
	assert.Nil(t, Tail(p), "Tail should be nil for map type")
}

// TestZeroWithStructs tests Zero function with struct types
func TestZeroWithStructs(t *testing.T) {
	type TestStruct struct {
		Field1 int
		Field2 string
	}

	p := Zero[TestStruct, TestStruct]()

	expected := TestStruct{Field1: 0, Field2: ""}
	assert.Equal(t, expected, Head(p), "Head should be zero value for struct")
	assert.Equal(t, expected, Tail(p), "Tail should be zero value for struct")
}

// TestZeroWithInterfaces tests Zero function with interface types
func TestZeroWithInterfaces(t *testing.T) {
	p := Zero[interface{}, interface{}]()

	assert.Nil(t, Head(p), "Head should be nil for interface type")
	assert.Nil(t, Tail(p), "Tail should be nil for interface type")
}

// TestZeroWithChannels tests Zero function with channel types
func TestZeroWithChannels(t *testing.T) {
	p := Zero[chan int, chan string]()

	assert.Nil(t, Head(p), "Head should be nil for channel type")
	assert.Nil(t, Tail(p), "Tail should be nil for channel type")
}

// TestZeroWithFunctions tests Zero function with function types
func TestZeroWithFunctions(t *testing.T) {
	p := Zero[func() int, func(string) bool]()

	assert.Nil(t, Head(p), "Head should be nil for function type")
	assert.Nil(t, Tail(p), "Tail should be nil for function type")
}

// TestZeroCanBeUsedWithOtherFunctions tests that Zero pairs work with other pair functions
func TestZeroCanBeUsedWithOtherFunctions(t *testing.T) {
	p := Zero[int, string]()

	// Test with Head and Tail
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "", Tail(p))

	// Test with First and Second
	assert.Equal(t, 0, First(p))
	assert.Equal(t, "", Second(p))

	// Test with Swap
	swapped := Swap(p)
	assert.Equal(t, "", Head(swapped))
	assert.Equal(t, 0, Tail(swapped))

	// Test with Map
	mapped := MonadMapTail(p, func(s string) int { return len(s) })
	assert.Equal(t, 0, Head(mapped))
	assert.Equal(t, 0, Tail(mapped))
}

// TestZeroEquality tests that multiple Zero calls produce equal pairs
func TestZeroEquality(t *testing.T) {
	p1 := Zero[int, string]()
	p2 := Zero[int, string]()

	assert.Equal(t, Head(p1), Head(p2), "Heads should be equal")
	assert.Equal(t, Tail(p1), Tail(p2), "Tails should be equal")
}

// TestZeroWithComplexTypes tests Zero with more complex nested types
func TestZeroWithComplexTypes(t *testing.T) {
	type ComplexType struct {
		Nested map[string][]int
		Ptr    *string
	}

	p := Zero[ComplexType, []map[string]int]()

	expectedHead := ComplexType{Nested: nil, Ptr: nil}
	assert.Equal(t, expectedHead, Head(p), "Head should be zero value for complex struct")
	assert.Nil(t, Tail(p), "Tail should be nil for slice of maps")
}

// TestUnpackWithIntegers tests Unpack function with integer types
func TestUnpackWithIntegers(t *testing.T) {
	p := MakePair(42, 100)
	head, tail := Unpack(p)

	assert.Equal(t, 42, head, "Head should be 42")
	assert.Equal(t, 100, tail, "Tail should be 100")
}

// TestUnpackWithStrings tests Unpack function with string types
func TestUnpackWithStrings(t *testing.T) {
	p := MakePair("hello", "world")
	head, tail := Unpack(p)

	assert.Equal(t, "hello", head, "Head should be 'hello'")
	assert.Equal(t, "world", tail, "Tail should be 'world'")
}

// TestUnpackWithMixedTypes tests Unpack function with different types
func TestUnpackWithMixedTypes(t *testing.T) {
	p := MakePair("Alice", 30)
	name, age := Unpack(p)

	assert.Equal(t, "Alice", name, "Name should be 'Alice'")
	assert.Equal(t, 30, age, "Age should be 30")
}

// TestUnpackWithBooleans tests Unpack function with boolean types
func TestUnpackWithBooleans(t *testing.T) {
	p := MakePair(true, false)
	head, tail := Unpack(p)

	assert.Equal(t, true, head, "Head should be true")
	assert.Equal(t, false, tail, "Tail should be false")
}

// TestUnpackWithFloats tests Unpack function with float types
func TestUnpackWithFloats(t *testing.T) {
	p := MakePair(3.14, float32(2.71))
	head, tail := Unpack(p)

	assert.Equal(t, 3.14, head, "Head should be 3.14")
	assert.Equal(t, float32(2.71), tail, "Tail should be 2.71")
}

// TestUnpackWithPointers tests Unpack function with pointer types
func TestUnpackWithPointers(t *testing.T) {
	x := 42
	y := "test"
	p := MakePair(&x, &y)
	head, tail := Unpack(p)

	assert.Equal(t, &x, head, "Head should point to x")
	assert.Equal(t, &y, tail, "Tail should point to y")
	assert.Equal(t, 42, *head, "Dereferenced head should be 42")
	assert.Equal(t, "test", *tail, "Dereferenced tail should be 'test'")
}

// TestUnpackWithSlices tests Unpack function with slice types
func TestUnpackWithSlices(t *testing.T) {
	p := MakePair([]int{1, 2, 3}, []string{"a", "b", "c"})
	head, tail := Unpack(p)

	assert.Equal(t, []int{1, 2, 3}, head, "Head should be [1, 2, 3]")
	assert.Equal(t, []string{"a", "b", "c"}, tail, "Tail should be ['a', 'b', 'c']")
}

// TestUnpackWithMaps tests Unpack function with map types
func TestUnpackWithMaps(t *testing.T) {
	m1 := map[string]int{"one": 1, "two": 2}
	m2 := map[int]string{1: "one", 2: "two"}
	p := MakePair(m1, m2)
	head, tail := Unpack(p)

	assert.Equal(t, m1, head, "Head should be the first map")
	assert.Equal(t, m2, tail, "Tail should be the second map")
}

// TestUnpackWithStructs tests Unpack function with struct types
func TestUnpackWithStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	p := MakePair(p1, p2)
	head, tail := Unpack(p)

	assert.Equal(t, p1, head, "Head should be Alice")
	assert.Equal(t, p2, tail, "Tail should be Bob")
}

// TestUnpackWithFunctions tests Unpack function with function types
func TestUnpackWithFunctions(t *testing.T) {
	f1 := func(x int) int { return x * 2 }
	f2 := func(x int) int { return x + 10 }
	p := MakePair(f1, f2)
	head, tail := Unpack(p)

	assert.Equal(t, 20, head(10), "Head function should double the input")
	assert.Equal(t, 20, tail(10), "Tail function should add 10 to the input")
}

// TestUnpackWithZeroPair tests Unpack function with zero-valued pair
func TestUnpackWithZeroPair(t *testing.T) {
	p := Zero[int, string]()
	head, tail := Unpack(p)

	assert.Equal(t, 0, head, "Head should be zero value for int")
	assert.Equal(t, "", tail, "Tail should be zero value for string")
}

// TestUnpackWithNilValues tests Unpack function with nil values
func TestUnpackWithNilValues(t *testing.T) {
	p := MakePair[*int, *string](nil, nil)
	head, tail := Unpack(p)

	assert.Nil(t, head, "Head should be nil")
	assert.Nil(t, tail, "Tail should be nil")
}

// TestUnpackInverseMakePair tests that Unpack is the inverse of MakePair
func TestUnpackInverseMakePair(t *testing.T) {
	original := MakePair("test", 123)
	head, tail := Unpack(original)
	reconstructed := MakePair(head, tail)

	assert.Equal(t, Head(original), Head(reconstructed), "Heads should be equal")
	assert.Equal(t, Tail(original), Tail(reconstructed), "Tails should be equal")
}

// TestUnpackWithOf tests Unpack with a pair created by Of
func TestUnpackWithOf(t *testing.T) {
	p := Of(42)
	head, tail := Unpack(p)

	assert.Equal(t, 42, head, "Head should be 42")
	assert.Equal(t, 42, tail, "Tail should be 42")
}

// TestUnpackWithSwap tests Unpack after swapping a pair
func TestUnpackWithSwap(t *testing.T) {
	original := MakePair("hello", 42)
	swapped := Swap(original)
	head, tail := Unpack(swapped)

	assert.Equal(t, 42, head, "Head should be 42 after swap")
	assert.Equal(t, "hello", tail, "Tail should be 'hello' after swap")
}

// TestUnpackWithMappedPair tests Unpack with a mapped pair
func TestUnpackWithMappedPair(t *testing.T) {
	original := MakePair(5, "hello")
	mapped := MonadMapTail(original, func(s string) int { return len(s) })
	head, tail := Unpack(mapped)

	assert.Equal(t, 5, head, "Head should remain 5")
	assert.Equal(t, 5, tail, "Tail should be length of 'hello'")
}

// TestUnpackWithComplexTypes tests Unpack with complex nested types
func TestUnpackWithComplexTypes(t *testing.T) {
	type ComplexType struct {
		Data   map[string][]int
		Nested *ComplexType
	}

	c1 := ComplexType{
		Data:   map[string][]int{"key": {1, 2, 3}},
		Nested: nil,
	}
	c2 := ComplexType{
		Data:   map[string][]int{"other": {4, 5, 6}},
		Nested: &c1,
	}

	p := MakePair(c1, c2)
	head, tail := Unpack(p)

	assert.Equal(t, c1, head, "Head should be c1")
	assert.Equal(t, c2, tail, "Tail should be c2")
	assert.NotNil(t, tail.Nested, "Tail's nested field should not be nil")
}

// TestUnpackMultipleAssignments tests that Unpack can be used in multiple assignments
func TestUnpackMultipleAssignments(t *testing.T) {
	p1 := MakePair(1, "one")
	p2 := MakePair(2, "two")

	h1, t1 := Unpack(p1)
	h2, t2 := Unpack(p2)

	assert.Equal(t, 1, h1)
	assert.Equal(t, "one", t1)
	assert.Equal(t, 2, h2)
	assert.Equal(t, "two", t2)
}

// TestUnpackWithChannels tests Unpack function with channel types
func TestUnpackWithChannels(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan string, 1)
	ch1 <- 42
	ch2 <- "test"

	p := MakePair(ch1, ch2)
	head, tail := Unpack(p)

	assert.Equal(t, 42, <-head, "Should receive 42 from head channel")
	assert.Equal(t, "test", <-tail, "Should receive 'test' from tail channel")
}

// TestUnpackWithInterfaces tests Unpack function with interface types
func TestUnpackWithInterfaces(t *testing.T) {
	var i1 interface{} = 42
	var i2 interface{} = "test"

	p := MakePair(i1, i2)
	head, tail := Unpack(p)

	assert.Equal(t, 42, head, "Head should be 42")
	assert.Equal(t, "test", tail, "Tail should be 'test'")
}
