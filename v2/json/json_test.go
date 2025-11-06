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

package json

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type Json map[string]any

func TestJsonMarshal(t *testing.T) {

	resRight := Unmarshal[Json]([]byte("{\"a\": \"b\"}"))
	assert.True(t, E.IsRight(resRight))

	resLeft := Unmarshal[Json]([]byte("{\"a\""))
	assert.True(t, E.IsLeft(resLeft))

	res1 := F.Pipe1(
		resRight,
		E.Chain(Marshal[Json]),
	)
	fmt.Println(res1)
}

func TestUnmarshalSuccess(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data := []byte(`{"name":"Alice","age":30}`)
	result := Unmarshal[Person](data)

	assert.True(t, E.IsRight(result))

	person := E.GetOrElse(func(error) Person { return Person{} })(result)
	assert.Equal(t, "Alice", person.Name)
	assert.Equal(t, 30, person.Age)
}

func TestUnmarshalError(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Invalid JSON
	data := []byte(`{"name":"Alice","age":`)
	result := Unmarshal[Person](data)

	assert.True(t, E.IsLeft(result))
}

func TestUnmarshalTypeMismatch(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Age is a string instead of int
	data := []byte(`{"name":"Alice","age":"thirty"}`)
	result := Unmarshal[Person](data)

	assert.True(t, E.IsLeft(result))
}

func TestMarshalSuccess(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := Person{Name: "Bob", Age: 25}
	result := Marshal(person)

	assert.True(t, E.IsRight(result))

	jsonBytes := E.GetOrElse(func(error) []byte { return []byte{} })(result)
	assert.Contains(t, string(jsonBytes), `"name":"Bob"`)
	assert.Contains(t, string(jsonBytes), `"age":25`)
}

func TestMarshalWithMap(t *testing.T) {
	data := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	result := Marshal(data)
	assert.True(t, E.IsRight(result))
}

func TestMarshalWithSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	result := Marshal(data)

	assert.True(t, E.IsRight(result))

	jsonBytes := E.GetOrElse(func(error) []byte { return []byte{} })(result)
	assert.Equal(t, "[1,2,3,4,5]", string(jsonBytes))
}

func TestMarshalWithNil(t *testing.T) {
	var data *string
	result := Marshal(data)

	assert.True(t, E.IsRight(result))

	jsonBytes := E.GetOrElse(func(error) []byte { return []byte{} })(result)
	assert.Equal(t, "null", string(jsonBytes))
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	original := Config{Host: "localhost", Port: 8080}

	result := F.Pipe2(
		original,
		Marshal[Config],
		E.Chain(Unmarshal[Config]),
	)

	assert.True(t, E.IsRight(result))

	recovered := E.GetOrElse(func(error) Config { return Config{} })(result)
	assert.Equal(t, original, recovered)
}

func TestMarshalWithNestedStructs(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	person := Person{
		Name: "Charlie",
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	result := Marshal(person)
	assert.True(t, E.IsRight(result))

	jsonBytes := E.GetOrElse(func(error) []byte { return []byte{} })(result)
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"name":"Charlie"`)
	assert.Contains(t, jsonStr, `"street":"123 Main St"`)
	assert.Contains(t, jsonStr, `"city":"Springfield"`)
}

func TestUnmarshalEmptyJSON(t *testing.T) {
	type Empty struct{}

	data := []byte(`{}`)
	result := Unmarshal[Empty](data)

	assert.True(t, E.IsRight(result))
}

func TestUnmarshalArray(t *testing.T) {
	data := []byte(`[1, 2, 3, 4, 5]`)
	result := Unmarshal[[]int](data)

	assert.True(t, E.IsRight(result))

	arr := E.GetOrElse(func(error) []int { return []int{} })(result)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr)
}

func TestUnmarshalNull(t *testing.T) {
	data := []byte(`null`)
	result := Unmarshal[*string](data)

	assert.True(t, E.IsRight(result))

	ptr := E.GetOrElse(func(error) *string { return nil })(result)
	assert.Nil(t, ptr)
}
