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

package io

import (
	"bytes"
	"log"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l := Logger[int]()
	lio := l("out")
	assert.NotPanics(t, func() { lio(10)() })
}

func TestLoggerWithCustomLogger(t *testing.T) {
	var buf bytes.Buffer
	customLogger := log.New(&buf, "", 0)

	l := Logger[int](customLogger)
	lio := l("test value")

	result := lio(42)()

	assert.Equal(t, 42, result)
	assert.Contains(t, buf.String(), "test value")
	assert.Contains(t, buf.String(), "42")
}

func TestLoggerReturnsOriginalValue(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value int
	}

	l := Logger[TestStruct]()
	lio := l("test")

	input := TestStruct{Name: "test", Value: 100}
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestLogf(t *testing.T) {
	l := Logf[int]
	lio := l("Value is %d")
	assert.NotPanics(t, func() { lio(10)() })
}

func TestLogfReturnsOriginalValue(t *testing.T) {
	l := Logf[string]
	lio := l("String: %s")

	input := "hello"
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestPrintfLogger(t *testing.T) {
	l := Printf[int]
	lio := l("Value: %d\n")
	assert.NotPanics(t, func() { lio(10)() })
}

func TestPrintfLoggerReturnsOriginalValue(t *testing.T) {
	l := Printf[float64]
	lio := l("Number: %.2f\n")

	input := 3.14159
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestLogGo(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	l := LogGo[User]
	lio := l("User: {{.Name}}, Age: {{.Age}}")

	input := User{Name: "Alice", Age: 30}
	assert.NotPanics(t, func() { lio(input)() })
}

func TestLogGoReturnsOriginalValue(t *testing.T) {
	type Product struct {
		ID    int
		Name  string
		Price float64
	}

	l := LogGo[Product]
	lio := l("Product: {{.Name}} ({{.ID}})")

	input := Product{ID: 123, Name: "Widget", Price: 19.99}
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestLogGoWithInvalidTemplate(t *testing.T) {
	l := LogGo[int]
	// Invalid template syntax
	lio := l("Value: {{.MissingField")

	// Should not panic even with invalid template
	assert.NotPanics(t, func() { lio(42)() })
}

func TestLogGoWithComplexTemplate(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}
	type Person struct {
		Name    string
		Address Address
	}

	l := LogGo[Person]
	lio := l("Person: {{.Name}} from {{.Address.City}}")

	input := Person{
		Name:    "Bob",
		Address: Address{Street: "Main St", City: "NYC"},
	}
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestPrintGo(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	l := PrintGo[User]
	lio := l("User: {{.Name}}, Age: {{.Age}}")

	input := User{Name: "Charlie", Age: 25}
	assert.NotPanics(t, func() { lio(input)() })
}

func TestPrintGoReturnsOriginalValue(t *testing.T) {
	type Score struct {
		Player string
		Points int
	}

	l := PrintGo[Score]
	lio := l("{{.Player}}: {{.Points}} points")

	input := Score{Player: "Alice", Points: 100}
	result := lio(input)()

	assert.Equal(t, input, result)
}

func TestPrintGoWithInvalidTemplate(t *testing.T) {
	l := PrintGo[string]
	// Invalid template syntax
	lio := l("Value: {{.}")

	// Should not panic even with invalid template
	assert.NotPanics(t, func() { lio("test")() })
}

func TestLogGoInPipeline(t *testing.T) {
	type Data struct {
		Value int
	}

	input := Data{Value: 10}

	result := F.Pipe2(
		Of(input),
		ChainFirst(LogGo[Data]("Processing: {{.Value}}")),
		Map(func(d Data) Data {
			return Data{Value: d.Value * 2}
		}),
	)()

	assert.Equal(t, 20, result.Value)
}

func TestPrintGoInPipeline(t *testing.T) {
	input := "hello"

	result := F.Pipe2(
		Of(input),
		ChainFirst(PrintGo[string]("Input: {{.}}")),
		Map(func(s string) string {
			return s + " world"
		}),
	)()

	assert.Equal(t, "hello world", result)
}
