// Copyright (c) 2024 IBM Corp.
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

package builder

import (
	"errors"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// Test types for demonstration

type Person struct {
	Name string
	Age  int
}

type PersonBuilder struct {
	name string
	age  int
}

func (b PersonBuilder) WithName(name string) PersonBuilder {
	b.name = name
	return b
}

func (b PersonBuilder) WithAge(age int) PersonBuilder {
	b.age = age
	return b
}

func (b PersonBuilder) Build() Result[Person] {
	if b.name == "" {
		return result.Left[Person](errors.New("name is required"))
	}
	if b.age < 0 {
		return result.Left[Person](errors.New("age must be non-negative"))
	}
	if b.age > 150 {
		return result.Left[Person](errors.New("age must be realistic"))
	}
	return result.Of(Person{Name: b.name, Age: b.age})
}

func NewPersonBuilder(p Person) PersonBuilder {
	return PersonBuilder{name: p.Name, age: p.Age}
}

// Config example for additional test coverage

type Config struct {
	Host string
	Port int
}

type ConfigBuilder struct {
	host string
	port int
}

func (b ConfigBuilder) WithHost(host string) ConfigBuilder {
	b.host = host
	return b
}

func (b ConfigBuilder) WithPort(port int) ConfigBuilder {
	b.port = port
	return b
}

func (b ConfigBuilder) Build() Result[Config] {
	if b.host == "" {
		return result.Left[Config](errors.New("host is required"))
	}
	if b.port <= 0 || b.port > 65535 {
		return result.Left[Config](errors.New("port must be between 1 and 65535"))
	}
	return result.Of(Config{Host: b.host, Port: b.port})
}

func NewConfigBuilder(c Config) ConfigBuilder {
	return ConfigBuilder{host: c.Host, port: c.Port}
}

// Tests for Builder interface

func TestBuilder_SuccessfulBuild(t *testing.T) {
	builder := PersonBuilder{}.
		WithName("Alice").
		WithAge(30)

	res := builder.Build()

	assert.True(t, result.IsRight(res), "Build should succeed")
	person := result.ToOption(res)
	assert.True(t, O.IsSome(person), "Result should contain a person")

	p := O.GetOrElse(func() Person { return Person{} })(person)
	assert.Equal(t, "Alice", p.Name)
	assert.Equal(t, 30, p.Age)
}

func TestBuilder_ValidationFailure_MissingName(t *testing.T) {
	builder := PersonBuilder{}.WithAge(30)

	res := builder.Build()

	assert.True(t, result.IsLeft(res), "Build should fail when name is missing")
	err := result.Fold(
		func(e error) error { return e },
		func(Person) error { return errors.New("unexpected success") },
	)(res)
	assert.Equal(t, "name is required", err.Error())
}

func TestBuilder_ValidationFailure_NegativeAge(t *testing.T) {
	builder := PersonBuilder{}.
		WithName("Bob").
		WithAge(-5)

	res := builder.Build()

	assert.True(t, result.IsLeft(res), "Build should fail when age is negative")
	err := result.Fold(
		func(e error) error { return e },
		func(Person) error { return errors.New("unexpected success") },
	)(res)
	assert.Equal(t, "age must be non-negative", err.Error())
}

func TestBuilder_ValidationFailure_UnrealisticAge(t *testing.T) {
	builder := PersonBuilder{}.
		WithName("Charlie").
		WithAge(200)

	res := builder.Build()

	assert.True(t, result.IsLeft(res), "Build should fail when age is unrealistic")
	err := result.Fold(
		func(e error) error { return e },
		func(Person) error { return errors.New("unexpected success") },
	)(res)
	assert.Equal(t, "age must be realistic", err.Error())
}

func TestBuilder_ConfigSuccessfulBuild(t *testing.T) {
	builder := ConfigBuilder{}.
		WithHost("localhost").
		WithPort(8080)

	res := builder.Build()

	assert.True(t, result.IsRight(res), "Build should succeed")
	config := result.ToOption(res)
	assert.True(t, O.IsSome(config), "Result should contain a config")

	c := O.GetOrElse(func() Config { return Config{} })(config)
	assert.Equal(t, "localhost", c.Host)
	assert.Equal(t, 8080, c.Port)
}

func TestBuilder_ConfigValidationFailure_MissingHost(t *testing.T) {
	builder := ConfigBuilder{}.WithPort(8080)

	res := builder.Build()

	assert.True(t, result.IsLeft(res), "Build should fail when host is missing")
	err := result.Fold(
		func(e error) error { return e },
		func(Config) error { return errors.New("unexpected success") },
	)(res)
	assert.Equal(t, "host is required", err.Error())
}

func TestBuilder_ConfigValidationFailure_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero port", 0},
		{"negative port", -1},
		{"port too large", 70000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ConfigBuilder{}.
				WithHost("localhost").
				WithPort(tt.port)

			res := builder.Build()

			assert.True(t, result.IsLeft(res), "Build should fail for invalid port")
			err := result.Fold(
				func(e error) error { return e },
				func(Config) error { return errors.New("unexpected success") },
			)(res)
			assert.Equal(t, "port must be between 1 and 65535", err.Error())
		})
	}
}

// Tests for BuilderPrism

func TestBuilderPrism_GetOption_ValidBuilder(t *testing.T) {
	prism := BuilderPrism(NewPersonBuilder)

	builder := PersonBuilder{}.
		WithName("Alice").
		WithAge(30)

	personOpt := prism.GetOption(builder)

	assert.True(t, O.IsSome(personOpt), "GetOption should return Some for valid builder")
	person := O.GetOrElse(func() Person { return Person{} })(personOpt)
	assert.Equal(t, "Alice", person.Name)
	assert.Equal(t, 30, person.Age)
}

func TestBuilderPrism_GetOption_InvalidBuilder(t *testing.T) {
	prism := BuilderPrism(NewPersonBuilder)

	builder := PersonBuilder{}.WithAge(30) // Missing name

	personOpt := prism.GetOption(builder)

	assert.True(t, O.IsNone(personOpt), "GetOption should return None for invalid builder")
}

func TestBuilderPrism_ReverseGet(t *testing.T) {
	prism := BuilderPrism(NewPersonBuilder)

	person := Person{Name: "Bob", Age: 25}

	builder := prism.ReverseGet(person)

	assert.Equal(t, "Bob", builder.name)
	assert.Equal(t, 25, builder.age)

	// Verify the builder can build the same person
	res := builder.Build()
	assert.True(t, result.IsRight(res), "Builder from ReverseGet should be valid")

	rebuilt := O.GetOrElse(func() Person { return Person{} })(result.ToOption(res))
	assert.Equal(t, person, rebuilt)
}

func TestBuilderPrism_RoundTrip_ValidBuilder(t *testing.T) {
	prism := BuilderPrism(NewPersonBuilder)

	originalBuilder := PersonBuilder{}.
		WithName("Charlie").
		WithAge(35)

	// Extract person from builder
	personOpt := prism.GetOption(originalBuilder)
	assert.True(t, O.IsSome(personOpt), "Should extract person from valid builder")

	person := O.GetOrElse(func() Person { return Person{} })(personOpt)

	// Reconstruct builder from person
	rebuiltBuilder := prism.ReverseGet(person)

	// Verify the rebuilt builder produces the same person
	rebuiltRes := rebuiltBuilder.Build()
	assert.True(t, result.IsRight(rebuiltRes), "Rebuilt builder should be valid")

	rebuiltPerson := O.GetOrElse(func() Person { return Person{} })(result.ToOption(rebuiltRes))
	assert.Equal(t, person, rebuiltPerson)
}

func TestBuilderPrism_ConfigPrism(t *testing.T) {
	prism := BuilderPrism(NewConfigBuilder)

	builder := ConfigBuilder{}.
		WithHost("example.com").
		WithPort(443)

	configOpt := prism.GetOption(builder)

	assert.True(t, O.IsSome(configOpt), "GetOption should return Some for valid config builder")
	config := O.GetOrElse(func() Config { return Config{} })(configOpt)
	assert.Equal(t, "example.com", config.Host)
	assert.Equal(t, 443, config.Port)
}

func TestBuilderPrism_ConfigPrism_InvalidBuilder(t *testing.T) {
	prism := BuilderPrism(NewConfigBuilder)

	builder := ConfigBuilder{}.WithPort(8080) // Missing host

	configOpt := prism.GetOption(builder)

	assert.True(t, O.IsNone(configOpt), "GetOption should return None for invalid config builder")
}

func TestBuilderPrism_ConfigPrism_ReverseGet(t *testing.T) {
	prism := BuilderPrism(NewConfigBuilder)

	config := Config{Host: "api.example.com", Port: 9000}

	builder := prism.ReverseGet(config)

	assert.Equal(t, "api.example.com", builder.host)
	assert.Equal(t, 9000, builder.port)

	// Verify the builder can build the same config
	res := builder.Build()
	assert.True(t, result.IsRight(res), "Builder from ReverseGet should be valid")

	rebuilt := O.GetOrElse(func() Config { return Config{} })(result.ToOption(res))
	assert.Equal(t, config, rebuilt)
}

// Benchmark tests

func BenchmarkBuilder_SuccessfulBuild(b *testing.B) {
	builder := PersonBuilder{}.
		WithName("Alice").
		WithAge(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Build()
	}
}

func BenchmarkBuilder_FailedBuild(b *testing.B) {
	builder := PersonBuilder{}.WithAge(30) // Missing name

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Build()
	}
}

func BenchmarkBuilderPrism_GetOption(b *testing.B) {
	prism := BuilderPrism(NewPersonBuilder)
	builder := PersonBuilder{}.
		WithName("Alice").
		WithAge(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prism.GetOption(builder)
	}
}

func BenchmarkBuilderPrism_ReverseGet(b *testing.B) {
	prism := BuilderPrism(NewPersonBuilder)
	person := Person{Name: "Bob", Age: 25}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prism.ReverseGet(person)
	}
}
