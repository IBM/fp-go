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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type TestType struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestToType(t *testing.T) {

	generic := map[string]any{"a": "value", "b": 1}

	assert.True(t, E.IsRight(ToTypeE[TestType](generic)))
	assert.True(t, E.IsRight(ToTypeE[TestType](&generic)))

	assert.Equal(t, E.Right[error](&TestType{A: "value", B: 1}), ToTypeE[*TestType](&generic))
	assert.Equal(t, E.Right[error](TestType{A: "value", B: 1}), F.Pipe1(ToTypeE[*TestType](&generic), E.Map[error](F.Deref[TestType])))
}

func TestToTypeESuccess(t *testing.T) {
	type Source struct {
		Name  string
		Value int
	}

	type Target struct {
		Name  string `json:"Name"`
		Value int    `json:"Value"`
	}

	src := Source{Name: "test", Value: 42}
	result := ToTypeE[Target](src)

	assert.True(t, E.IsRight(result))

	target := E.GetOrElse(func(error) Target { return Target{} })(result)
	assert.Equal(t, "test", target.Name)
	assert.Equal(t, 42, target.Value)
}

func TestToTypeEFromMap(t *testing.T) {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	data := map[string]any{"host": "localhost", "port": 8080}
	result := ToTypeE[Config](data)

	assert.True(t, E.IsRight(result))

	config := E.GetOrElse(func(error) Config { return Config{} })(result)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 8080, config.Port)
}

func TestToTypeEError(t *testing.T) {
	type Target struct {
		Value int `json:"value"`
	}

	// Invalid source that can't be marshaled
	src := make(chan int)
	result := ToTypeE[Target](src)

	assert.True(t, E.IsLeft(result))
}

func TestToTypeETypeMismatch(t *testing.T) {
	type Target struct {
		Value int `json:"value"`
	}

	// Source has string where int is expected
	data := map[string]any{"value": "not a number"}
	result := ToTypeE[Target](data)

	assert.True(t, E.IsLeft(result))
}

func TestToTypeOSuccess(t *testing.T) {
	type Source struct {
		Name  string
		Value int
	}

	type Target struct {
		Name  string `json:"Name"`
		Value int    `json:"Value"`
	}

	src := Source{Name: "test", Value: 42}
	result := ToTypeO[Target](src)

	assert.True(t, O.IsSome(result))

	target := O.GetOrElse(func() Target { return Target{} })(result)
	assert.Equal(t, "test", target.Name)
	assert.Equal(t, 42, target.Value)
}

func TestToTypeOFromMap(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data := map[string]any{"name": "Alice", "age": 30}
	result := ToTypeO[Person](data)

	assert.True(t, O.IsSome(result))

	person := O.GetOrElse(func() Person { return Person{} })(result)
	assert.Equal(t, "Alice", person.Name)
	assert.Equal(t, 30, person.Age)
}

func TestToTypeOError(t *testing.T) {
	type Target struct {
		Value int `json:"value"`
	}

	// Invalid source that can't be marshaled
	src := make(chan int)
	result := ToTypeO[Target](src)

	assert.True(t, O.IsNone(result))
}

func TestToTypeOTypeMismatch(t *testing.T) {
	type Target struct {
		Value int `json:"value"`
	}

	// Source has string where int is expected
	data := map[string]any{"value": "not a number"}
	result := ToTypeO[Target](data)

	assert.True(t, O.IsNone(result))
}

func TestToTypeEWithNestedStructs(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	data := map[string]any{
		"name": "Bob",
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "Springfield",
		},
	}

	result := ToTypeE[Person](data)
	assert.True(t, E.IsRight(result))

	person := E.GetOrElse(func(error) Person { return Person{} })(result)
	assert.Equal(t, "Bob", person.Name)
	assert.Equal(t, "123 Main St", person.Address.Street)
	assert.Equal(t, "Springfield", person.Address.City)
}

func TestToTypeOWithNestedStructs(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	data := map[string]any{
		"name": "Charlie",
		"address": map[string]any{
			"street": "456 Oak Ave",
			"city":   "Shelbyville",
		},
	}

	result := ToTypeO[Person](data)
	assert.True(t, O.IsSome(result))

	person := O.GetOrElse(func() Person { return Person{} })(result)
	assert.Equal(t, "Charlie", person.Name)
	assert.Equal(t, "456 Oak Ave", person.Address.Street)
	assert.Equal(t, "Shelbyville", person.Address.City)
}

func TestToTypeEWithSlice(t *testing.T) {
	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := []map[string]any{
		{"id": 1, "name": "Item1"},
		{"id": 2, "name": "Item2"},
	}

	result := ToTypeE[[]Item](data)
	assert.True(t, E.IsRight(result))

	items := E.GetOrElse(func(error) []Item { return []Item{} })(result)
	assert.Len(t, items, 2)
	assert.Equal(t, 1, items[0].ID)
	assert.Equal(t, "Item1", items[0].Name)
	assert.Equal(t, 2, items[1].ID)
	assert.Equal(t, "Item2", items[1].Name)
}

func TestToTypeOWithSlice(t *testing.T) {
	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := []map[string]any{
		{"id": 3, "name": "Item3"},
		{"id": 4, "name": "Item4"},
	}

	result := ToTypeO[[]Item](data)
	assert.True(t, O.IsSome(result))

	items := O.GetOrElse(func() []Item { return []Item{} })(result)
	assert.Len(t, items, 2)
	assert.Equal(t, 3, items[0].ID)
	assert.Equal(t, "Item3", items[0].Name)
	assert.Equal(t, 4, items[1].ID)
	assert.Equal(t, "Item4", items[1].Name)
}

func TestToTypeEWithPointer(t *testing.T) {
	type Data struct {
		Value *int `json:"value"`
	}

	val := 100
	src := map[string]any{"value": val}
	result := ToTypeE[Data](src)

	assert.True(t, E.IsRight(result))

	data := E.GetOrElse(func(error) Data { return Data{} })(result)
	assert.NotNil(t, data.Value)
	assert.Equal(t, 100, *data.Value)
}

func TestToTypeOWithNullValue(t *testing.T) {
	type Data struct {
		Value *int `json:"value"`
	}

	src := map[string]any{"value": nil}
	result := ToTypeO[Data](src)

	assert.True(t, O.IsSome(result))

	data := O.GetOrElse(func() Data { return Data{} })(result)
	assert.Nil(t, data.Value)
}
