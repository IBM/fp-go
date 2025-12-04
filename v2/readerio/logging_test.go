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

package readerio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestConfig struct {
		AppName string
		Debug   bool
	}

	TestUser struct {
		Name string
		Age  int
	}

	TestData struct {
		ID    int
		Value string
	}
)

func TestLogf(t *testing.T) {
	l := Logf[TestConfig, int]("[%v] Value: %d")
	rio := l(42)

	config := TestConfig{AppName: "TestApp", Debug: true}
	result := rio(config)()

	assert.Equal(t, 42, result)
}

func TestLogfReturnsOriginalValue(t *testing.T) {
	l := Logf[TestConfig, TestUser]("[%v] User: %+v")
	rio := l(TestUser{Name: "Alice", Age: 30})

	config := TestConfig{AppName: "TestApp"}
	result := rio(config)()

	assert.Equal(t, TestUser{Name: "Alice", Age: 30}, result)
}

func TestLogfWithDifferentTypes(t *testing.T) {
	// Test with string value
	l1 := Logf[string, string]("[%s] Message: %s")
	rio1 := l1("hello")
	result1 := rio1("context")()
	assert.Equal(t, "hello", result1)

	// Test with struct value
	l2 := Logf[int, TestData]("[%d] Data: %+v")
	rio2 := l2(TestData{ID: 123, Value: "test"})
	result2 := rio2(999)()
	assert.Equal(t, TestData{ID: 123, Value: "test"}, result2)
}

func TestPrintf(t *testing.T) {
	l := Printf[TestConfig, int]("[%v] Value: %d\n")
	rio := l(42)

	config := TestConfig{AppName: "TestApp", Debug: false}
	result := rio(config)()

	assert.Equal(t, 42, result)
}

func TestPrintfReturnsOriginalValue(t *testing.T) {
	l := Printf[TestConfig, TestUser]("[%v] User: %+v\n")
	rio := l(TestUser{Name: "Bob", Age: 25})

	config := TestConfig{AppName: "TestApp"}
	result := rio(config)()

	assert.Equal(t, TestUser{Name: "Bob", Age: 25}, result)
}

func TestPrintfWithDifferentTypes(t *testing.T) {
	// Test with float value
	l1 := Printf[string, float64]("[%s] Number: %.2f\n")
	rio1 := l1(3.14159)
	result1 := rio1("PI")()
	assert.Equal(t, 3.14159, result1)

	// Test with bool value
	l2 := Printf[int, bool]("[%d] Flag: %v\n")
	rio2 := l2(true)
	result2 := rio2(1)()
	assert.True(t, result2)
}

func TestLogGo(t *testing.T) {
	l := LogGo[TestConfig, TestUser]("[{{.R.AppName}}] User: {{.A.Name}}, Age: {{.A.Age}}")
	rio := l(TestUser{Name: "Charlie", Age: 35})

	config := TestConfig{AppName: "MyApp", Debug: true}
	result := rio(config)()

	assert.Equal(t, TestUser{Name: "Charlie", Age: 35}, result)
}

func TestLogGoReturnsOriginalValue(t *testing.T) {
	l := LogGo[TestConfig, TestData]("{{.R.AppName}}: Data {{.A.ID}} - {{.A.Value}}")
	rio := l(TestData{ID: 456, Value: "test data"})

	config := TestConfig{AppName: "TestApp"}
	result := rio(config)()

	assert.Equal(t, TestData{ID: 456, Value: "test data"}, result)
}

func TestLogGoWithInvalidTemplate(t *testing.T) {
	// Invalid template syntax - should not panic
	l := LogGo[TestConfig, int]("Value: {{.A.MissingField")
	rio := l(42)

	config := TestConfig{AppName: "TestApp"}
	assert.NotPanics(t, func() {
		result := rio(config)()
		assert.Equal(t, 42, result)
	})
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

	l := LogGo[TestConfig, Person]("[{{.R.AppName}}] Person: {{.A.Name}} from {{.A.Address.City}}")
	rio := l(Person{
		Name:    "David",
		Address: Address{Street: "Main St", City: "NYC"},
	})

	config := TestConfig{AppName: "TestApp"}
	result := rio(config)()

	assert.Equal(t, "David", result.Name)
	assert.Equal(t, "NYC", result.Address.City)
}

func TestLogGoWithConditionalTemplate(t *testing.T) {
	l := LogGo[TestConfig, int]("{{if .R.Debug}}[DEBUG] {{end}}Value: {{.A}}")
	rio := l(100)

	config := TestConfig{AppName: "TestApp", Debug: true}
	result := rio(config)()

	assert.Equal(t, 100, result)
}

func TestPrintGo(t *testing.T) {
	l := PrintGo[TestConfig, TestUser]("[{{.R.AppName}}] User: {{.A.Name}}, Age: {{.A.Age}}")
	rio := l(TestUser{Name: "Eve", Age: 28})

	config := TestConfig{AppName: "MyApp", Debug: false}
	result := rio(config)()

	assert.Equal(t, TestUser{Name: "Eve", Age: 28}, result)
}

func TestPrintGoReturnsOriginalValue(t *testing.T) {
	l := PrintGo[TestConfig, TestData]("{{.R.AppName}}: {{.A.ID}} - {{.A.Value}}")
	rio := l(TestData{ID: 789, Value: "sample"})

	config := TestConfig{AppName: "TestApp"}
	result := rio(config)()

	assert.Equal(t, TestData{ID: 789, Value: "sample"}, result)
}

func TestPrintGoWithInvalidTemplate(t *testing.T) {
	// Invalid template syntax - should not panic
	l := PrintGo[TestConfig, string]("Value: {{.")
	rio := l("test")

	config := TestConfig{AppName: "TestApp"}
	assert.NotPanics(t, func() {
		result := rio(config)()
		assert.Equal(t, "test", result)
	})
}

func TestPrintGoWithComplexTemplate(t *testing.T) {
	type Score struct {
		Player string
		Points int
	}

	l := PrintGo[TestConfig, Score]("{{if .R.Debug}}[DEBUG] {{end}}{{.A.Player}}: {{.A.Points}} points")
	rio := l(Score{Player: "Alice", Points: 100})

	config := TestConfig{AppName: "GameApp", Debug: true}
	result := rio(config)()

	assert.Equal(t, "Alice", result.Player)
	assert.Equal(t, 100, result.Points)
}

func TestLogGoInPipeline(t *testing.T) {
	config := TestConfig{AppName: "PipelineApp", Debug: true}

	// Create a pipeline using Chain and logging
	pipeline := MonadChain(
		LogGo[TestConfig, TestData]("[{{.R.AppName}}] Processing: {{.A.ID}}")(TestData{ID: 10, Value: "initial"}),
		func(d TestData) ReaderIO[TestConfig, TestData] {
			return Of[TestConfig](TestData{ID: d.ID * 2, Value: d.Value + "_processed"})
		},
	)

	result := pipeline(config)()

	assert.Equal(t, 20, result.ID)
	assert.Equal(t, "initial_processed", result.Value)
}

func TestPrintGoInPipeline(t *testing.T) {
	config := TestConfig{AppName: "PrintApp", Debug: false}

	pipeline := MonadChain(
		PrintGo[TestConfig, string]("[{{.R.AppName}}] Input: {{.A}}")("hello"),
		func(s string) ReaderIO[TestConfig, string] {
			return Of[TestConfig](s + " world")
		},
	)

	result := pipeline(config)()

	assert.Equal(t, "hello world", result)
}

func TestLogfInPipeline(t *testing.T) {
	config := TestConfig{AppName: "LogfApp"}

	pipeline := MonadChain(
		Logf[TestConfig, int]("[%v] Value: %d")(5),
		func(n int) ReaderIO[TestConfig, int] {
			return Of[TestConfig](n * 3)
		},
	)

	result := pipeline(config)()

	assert.Equal(t, 15, result)
}

func TestPrintfInPipeline(t *testing.T) {
	config := TestConfig{AppName: "PrintfApp"}

	pipeline := MonadChain(
		Printf[TestConfig, float64]("[%v] Number: %.1f\n")(2.5),
		func(n float64) ReaderIO[TestConfig, float64] {
			return Of[TestConfig](n * 2)
		},
	)

	result := pipeline(config)()

	assert.Equal(t, 5.0, result)
}

func TestMultipleLoggersInPipeline(t *testing.T) {
	config := TestConfig{AppName: "MultiApp", Debug: true}

	pipeline := MonadChain(
		Logf[TestConfig, int]("[%v] Initial: %d")(10),
		func(n int) ReaderIO[TestConfig, int] {
			return MonadChain(
				LogGo[TestConfig, int]("[{{.R.AppName}}] After add: {{.A}}")(n+5),
				func(n int) ReaderIO[TestConfig, int] {
					return Of[TestConfig](n * 2)
				},
			)
		},
	)

	result := pipeline(config)()

	assert.Equal(t, 30, result)
}

func TestLogGoWithNestedStructs(t *testing.T) {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Name  string
		Inner Inner
	}

	l := LogGo[TestConfig, Outer]("[{{.R.AppName}}] {{.A.Name}}: {{.A.Inner.Value}}")
	rio := l(Outer{Name: "Test", Inner: Inner{Value: 42}})

	config := TestConfig{AppName: "NestedApp"}
	result := rio(config)()

	assert.Equal(t, "Test", result.Name)
	assert.Equal(t, 42, result.Inner.Value)
}

func TestPrintGoWithNestedStructs(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}
	type Request struct {
		Method string
		Config Config
	}

	l := PrintGo[TestConfig, Request]("{{.A.Method}} -> {{.A.Config.Host}}:{{.A.Config.Port}}")
	rio := l(Request{
		Method: "GET",
		Config: Config{Host: "localhost", Port: 8080},
	})

	config := TestConfig{AppName: "HTTPApp"}
	result := rio(config)()

	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, "localhost", result.Config.Host)
	assert.Equal(t, 8080, result.Config.Port)
}

func TestLogGoWithEmptyTemplate(t *testing.T) {
	l := LogGo[TestConfig, int]("")
	rio := l(42)

	config := TestConfig{AppName: "EmptyApp"}
	result := rio(config)()

	assert.Equal(t, 42, result)
}

func TestPrintGoWithEmptyTemplate(t *testing.T) {
	l := PrintGo[TestConfig, string]("")
	rio := l("test")

	config := TestConfig{AppName: "EmptyApp"}
	result := rio(config)()

	assert.Equal(t, "test", result)
}
