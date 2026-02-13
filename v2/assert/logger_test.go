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

package assert

import (
	"fmt"
	"testing"
)

// TestLogf_BasicInteger tests basic integer logging
func TestLogf_BasicInteger(t *testing.T) {
	logInt := Logf[int]("Processing value: %d")

	// This should not panic and should log the value
	logInt(42)(t)()

	// Test passes if no panic occurs
}

// TestLogf_BasicString tests basic string logging
func TestLogf_BasicString(t *testing.T) {
	logString := Logf[string]("String value: %s")

	logString("hello world")(t)()

	// Test passes if no panic occurs
}

// TestLogf_BasicFloat tests basic float logging
func TestLogf_BasicFloat(t *testing.T) {
	logFloat := Logf[float64]("Float value: %.2f")

	logFloat(3.14159)(t)()

	// Test passes if no panic occurs
}

// TestLogf_BasicBoolean tests basic boolean logging
func TestLogf_BasicBoolean(t *testing.T) {
	logBool := Logf[bool]("Boolean value: %t")

	logBool(true)(t)()
	logBool(false)(t)()

	// Test passes if no panic occurs
}

// TestLogf_ComplexStruct tests logging of complex structures
func TestLogf_ComplexStruct(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	logUser := Logf[User]("User: %+v")

	user := User{Name: "Alice", Age: 30}
	logUser(user)(t)()

	// Test passes if no panic occurs
}

// TestLogf_Slice tests logging of slices
func TestLogf_Slice(t *testing.T) {
	logSlice := Logf[[]int]("Slice: %v")

	numbers := []int{1, 2, 3, 4, 5}
	logSlice(numbers)(t)()

	// Test passes if no panic occurs
}

// TestLogf_Map tests logging of maps
func TestLogf_Map(t *testing.T) {
	logMap := Logf[map[string]int]("Map: %v")

	data := map[string]int{"a": 1, "b": 2, "c": 3}
	logMap(data)(t)()

	// Test passes if no panic occurs
}

// TestLogf_Pointer tests logging of pointers
func TestLogf_Pointer(t *testing.T) {
	logPtr := Logf[*int]("Pointer: %p")

	value := 42
	logPtr(&value)(t)()

	// Test passes if no panic occurs
}

// TestLogf_NilPointer tests logging of nil pointers
func TestLogf_NilPointer(t *testing.T) {
	logPtr := Logf[*int]("Pointer: %v")

	var nilPtr *int
	logPtr(nilPtr)(t)()

	// Test passes if no panic occurs
}

// TestLogf_EmptyString tests logging of empty strings
func TestLogf_EmptyString(t *testing.T) {
	logString := Logf[string]("String: '%s'")

	logString("")(t)()

	// Test passes if no panic occurs
}

// TestLogf_EmptySlice tests logging of empty slices
func TestLogf_EmptySlice(t *testing.T) {
	logSlice := Logf[[]int]("Slice: %v")

	logSlice([]int{})(t)()

	// Test passes if no panic occurs
}

// TestLogf_EmptyMap tests logging of empty maps
func TestLogf_EmptyMap(t *testing.T) {
	logMap := Logf[map[string]int]("Map: %v")

	logMap(map[string]int{})(t)()

	// Test passes if no panic occurs
}

// TestLogf_MultipleTypes tests using multiple loggers for different types
func TestLogf_MultipleTypes(t *testing.T) {
	logString := Logf[string]("String: %s")
	logInt := Logf[int]("Integer: %d")
	logFloat := Logf[float64]("Float: %.2f")

	logString("test")(t)()
	logInt(42)(t)()
	logFloat(3.14)(t)()

	// Test passes if no panic occurs
}

// TestLogf_WithinTestPipeline tests logging within a test pipeline
func TestLogf_WithinTestPipeline(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}

	config := Config{Host: "localhost", Port: 8080}

	logConfig := Logf[Config]("Testing config: %+v")
	logConfig(config)(t)()

	// Continue with assertions
	StringNotEmpty(config.Host)(t)
	That(func(port int) bool { return port > 0 })(config.Port)(t)

	// Test passes if no panic occurs and assertions pass
}

// TestLogf_NestedStructures tests logging of nested structures
func TestLogf_NestedStructures(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address Address
	}

	logPerson := Logf[Person]("Person: %+v")

	person := Person{
		Name: "Bob",
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	logPerson(person)(t)()

	// Test passes if no panic occurs
}

// TestLogf_Interface tests logging of interface values
func TestLogf_Interface(t *testing.T) {
	logAny := Logf[any]("Value: %v")

	logAny(42)(t)()
	logAny("string")(t)()
	logAny([]int{1, 2, 3})(t)()

	// Test passes if no panic occurs
}

// TestLogf_GoSyntaxFormat tests logging with Go-syntax format
func TestLogf_GoSyntaxFormat(t *testing.T) {
	type Point struct {
		X int
		Y int
	}

	logPoint := Logf[Point]("Point: %#v")

	point := Point{X: 10, Y: 20}
	logPoint(point)(t)()

	// Test passes if no panic occurs
}

// TestLogf_TypeFormat tests logging with type format
func TestLogf_TypeFormat(t *testing.T) {
	logType := Logf[any]("Type: %T, Value: %v")

	logType(42)(t)()
	logType("string")(t)()
	logType(3.14)(t)()

	// Test passes if no panic occurs
}

// TestLogf_LargeNumbers tests logging of large numbers
func TestLogf_LargeNumbers(t *testing.T) {
	logInt := Logf[int64]("Large number: %d")

	logInt(9223372036854775807)(t)() // Max int64

	// Test passes if no panic occurs
}

// TestLogf_NegativeNumbers tests logging of negative numbers
func TestLogf_NegativeNumbers(t *testing.T) {
	logInt := Logf[int]("Number: %d")

	logInt(-42)(t)()
	logInt(-100)(t)()

	// Test passes if no panic occurs
}

// TestLogf_SpecialFloats tests logging of special float values
func TestLogf_SpecialFloats(t *testing.T) {
	logFloat := Logf[float64]("Float: %v")

	logFloat(0.0)(t)()
	logFloat(-0.0)(t)()

	// Test passes if no panic occurs
}

// TestLogf_UnicodeStrings tests logging of unicode strings
func TestLogf_UnicodeStrings(t *testing.T) {
	logString := Logf[string]("Unicode: %s")

	logString("Hello, 世界")(t)()
	logString("Emoji: 🎉🎊")(t)()

	// Test passes if no panic occurs
}

// TestLogf_MultilineStrings tests logging of multiline strings
func TestLogf_MultilineStrings(t *testing.T) {
	logString := Logf[string]("Multiline:\n%s")

	multiline := `Line 1
Line 2
Line 3`

	logString(multiline)(t)()

	// Test passes if no panic occurs
}

// TestLogf_ReuseLogger tests reusing the same logger multiple times
func TestLogf_ReuseLogger(t *testing.T) {
	logInt := Logf[int]("Value: %d")

	for i := 0; i < 5; i++ {
		logInt(i)(t)()
	}

	// Test passes if no panic occurs
}

// TestLogf_ConditionalLogging tests conditional logging based on values
func TestLogf_ConditionalLogging(t *testing.T) {
	logDebug := Logf[string]("DEBUG: %s")

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		if v%2 == 0 {
			logDebug(fmt.Sprintf("Found even number: %d", v))(t)()
		}
	}

	// Test passes if no panic occurs
}

// TestLogf_WithAssertions tests combining logging with assertions
func TestLogf_WithAssertions(t *testing.T) {
	logInt := Logf[int]("Testing value: %d")

	value := 42
	logInt(value)(t)()

	// Perform assertion after logging
	Equal(42)(value)(t)

	// Test passes if assertion passes
}

// TestLogf_DebuggingFailures demonstrates using logging to debug test failures
func TestLogf_DebuggingFailures(t *testing.T) {
	logSlice := Logf[[]int]("Input slice: %v")
	logInt := Logf[int]("Computed sum: %d")

	numbers := []int{1, 2, 3, 4, 5}
	logSlice(numbers)(t)()

	sum := 0
	for _, n := range numbers {
		sum += n
	}
	logInt(sum)(t)()

	Equal(15)(sum)(t)

	// Test passes if assertion passes
}

// TestLogf_ComplexDataStructures tests logging of complex nested data
func TestLogf_ComplexDataStructures(t *testing.T) {
	type Metadata struct {
		Version string
		Tags    []string
	}

	type Document struct {
		ID       int
		Title    string
		Metadata Metadata
	}

	logDoc := Logf[Document]("Document: %+v")

	doc := Document{
		ID:    1,
		Title: "Test Document",
		Metadata: Metadata{
			Version: "1.0",
			Tags:    []string{"test", "example"},
		},
	}

	logDoc(doc)(t)()

	// Test passes if no panic occurs
}

// TestLogf_ArrayTypes tests logging of array types
func TestLogf_ArrayTypes(t *testing.T) {
	logArray := Logf[[5]int]("Array: %v")

	arr := [5]int{1, 2, 3, 4, 5}
	logArray(arr)(t)()

	// Test passes if no panic occurs
}

// TestLogf_ChannelTypes tests logging of channel types
func TestLogf_ChannelTypes(t *testing.T) {
	logChan := Logf[chan int]("Channel: %v")

	ch := make(chan int, 1)
	logChan(ch)(t)()
	close(ch)

	// Test passes if no panic occurs
}

// TestLogf_FunctionTypes tests logging of function types
func TestLogf_FunctionTypes(t *testing.T) {
	logFunc := Logf[func() int]("Function: %v")

	fn := func() int { return 42 }
	logFunc(fn)(t)()

	// Test passes if no panic occurs
}
