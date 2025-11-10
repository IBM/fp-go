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

package result

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) Result[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) Result[string] {
	return Of("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res, Of("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res, Of("John Doe"))
}

// Test types for lens-based operations
type Counter struct {
	Value int
}

type Person struct {
	Name string
	Age  int
}

type Config struct {
	Debug   bool
	Timeout int
}

func TestApSL(t *testing.T) {
	// Create a lens for the Age field
	ageLens := L.MakeLens(
		func(p Person) int { return p.Age },
		func(p Person, a int) Person { p.Age = a; return p },
	)

	t.Run("ApSL with Right value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Right(30)),
		)

		expected := Right(Person{Name: "Alice", Age: 30})
		assert.Equal(t, expected, result)
	})

	t.Run("ApSL with Left in context", func(t *testing.T) {
		result := F.Pipe1(
			Left[Person](assert.AnError),
			ApSL(ageLens, Right(30)),
		)

		expected := Left[Person](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("ApSL with Left in value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Left[int](assert.AnError)),
		)

		expected := Left[Person](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("ApSL with both Left", func(t *testing.T) {
		result := F.Pipe1(
			Left[Person](assert.AnError),
			ApSL(ageLens, Left[int](assert.AnError)),
		)

		expected := Left[Person](assert.AnError)
		assert.Equal(t, expected, result)
	})
}

func TestBindL(t *testing.T) {
	// Create a lens for the Value field
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("BindL with successful transformation", func(t *testing.T) {
		// Increment the counter, but fail if it would exceed 100
		increment := func(v int) Result[int] {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right(v + 1)
		}

		result := F.Pipe1(
			Right(Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		expected := Right(Counter{Value: 43})
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with failing transformation", func(t *testing.T) {
		increment := func(v int) Result[int] {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right(v + 1)
		}

		result := F.Pipe1(
			Right(Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with Left input", func(t *testing.T) {
		increment := func(v int) Result[int] {
			return Right(v + 1)
		}

		result := F.Pipe1(
			Left[Counter](assert.AnError),
			BindL(valueLens, increment),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with multiple operations", func(t *testing.T) {
		double := func(v int) Result[int] {
			return Right(v * 2)
		}

		addTen := func(v int) Result[int] {
			return Right(v + 10)
		}

		result := F.Pipe2(
			Right(Counter{Value: 5}),
			BindL(valueLens, double),
			BindL(valueLens, addTen),
		)

		expected := Right(Counter{Value: 20})
		assert.Equal(t, expected, result)
	})
}

func TestLetL(t *testing.T) {
	// Create a lens for the Value field
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("LetL with pure transformation", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Right(Counter{Value: 21}),
			LetL(valueLens, double),
		)

		expected := Right(Counter{Value: 42})
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with Left input", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Left[Counter](assert.AnError),
			LetL(valueLens, double),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with multiple transformations", func(t *testing.T) {
		double := func(v int) int { return v * 2 }
		addTen := func(v int) int { return v + 10 }

		result := F.Pipe2(
			Right(Counter{Value: 5}),
			LetL(valueLens, double),
			LetL(valueLens, addTen),
		)

		expected := Right(Counter{Value: 20})
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with identity transformation", func(t *testing.T) {
		identity := func(v int) int { return v }

		result := F.Pipe1(
			Right(Counter{Value: 42}),
			LetL(valueLens, identity),
		)

		expected := Right(Counter{Value: 42})
		assert.Equal(t, expected, result)
	})
}

func TestLetToL(t *testing.T) {
	// Create a lens for the Debug field
	debugLens := L.MakeLens(
		func(c Config) bool { return c.Debug },
		func(c Config, d bool) Config { c.Debug = d; return c },
	)

	t.Run("LetToL with constant value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Config{Debug: true, Timeout: 30}),
			LetToL(debugLens, false),
		)

		expected := Right(Config{Debug: false, Timeout: 30})
		assert.Equal(t, expected, result)
	})

	t.Run("LetToL with Left input", func(t *testing.T) {
		result := F.Pipe1(
			Left[Config](assert.AnError),
			LetToL(debugLens, false),
		)

		expected := Left[Config](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("LetToL with multiple fields", func(t *testing.T) {
		timeoutLens := L.MakeLens(
			func(c Config) int { return c.Timeout },
			func(c Config, t int) Config { c.Timeout = t; return c },
		)

		result := F.Pipe2(
			Right(Config{Debug: true, Timeout: 30}),
			LetToL(debugLens, false),
			LetToL(timeoutLens, 60),
		)

		expected := Right(Config{Debug: false, Timeout: 60})
		assert.Equal(t, expected, result)
	})

	t.Run("LetToL setting same value", func(t *testing.T) {
		result := F.Pipe1(
			Right(Config{Debug: false, Timeout: 30}),
			LetToL(debugLens, false),
		)

		expected := Right(Config{Debug: false, Timeout: 30})
		assert.Equal(t, expected, result)
	})
}

func TestLensOperationsCombined(t *testing.T) {
	// Test combining different lens operations
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("Combine LetToL and LetL", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe2(
			Right(Counter{Value: 100}),
			LetToL(valueLens, 10),
			LetL(valueLens, double),
		)

		expected := Right(Counter{Value: 20})
		assert.Equal(t, expected, result)
	})

	t.Run("Combine LetL and BindL", func(t *testing.T) {
		double := func(v int) int { return v * 2 }
		validate := func(v int) Result[int] {
			if v > 100 {
				return Left[int](assert.AnError)
			}
			return Right(v)
		}

		result := F.Pipe2(
			Right(Counter{Value: 25}),
			LetL(valueLens, double),
			BindL(valueLens, validate),
		)

		expected := Right(Counter{Value: 50})
		assert.Equal(t, expected, result)
	})

	t.Run("Combine ApSL and LetL", func(t *testing.T) {
		addFive := func(v int) int { return v + 5 }

		result := F.Pipe2(
			Right(Counter{Value: 10}),
			ApSL(valueLens, Right(20)),
			LetL(valueLens, addFive),
		)

		expected := Right(Counter{Value: 25})
		assert.Equal(t, expected, result)
	})
}
