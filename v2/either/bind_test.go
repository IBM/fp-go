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

package either

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) Either[error, string] {
	return Of[error]("Doe")
}

func getGivenName(s utils.WithLastName) Either[error, string] {
	return Of[error]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res, Of[error]("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		ApS(utils.SetLastName, Of[error]("Doe")),
		ApS(utils.SetGivenName, Of[error]("John")),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res, Of[error]("John Doe"))
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
			Right[error](Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Right[error](30)),
		)

		expected := Right[error](Person{Name: "Alice", Age: 30})
		assert.Equal(t, expected, result)
	})

	t.Run("ApSL with Left in context", func(t *testing.T) {
		result := F.Pipe1(
			Left[Person](assert.AnError),
			ApSL(ageLens, Right[error](30)),
		)

		expected := Left[Person](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("ApSL with Left in value", func(t *testing.T) {
		result := F.Pipe1(
			Right[error](Person{Name: "Alice", Age: 25}),
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
		increment := func(v int) Either[error, int] {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right[error](v + 1)
		}

		result := F.Pipe1(
			Right[error](Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		expected := Right[error](Counter{Value: 43})
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with failing transformation", func(t *testing.T) {
		increment := func(v int) Either[error, int] {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right[error](v + 1)
		}

		result := F.Pipe1(
			Right[error](Counter{Value: 100}),
			BindL(valueLens, increment),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with Left input", func(t *testing.T) {
		increment := func(v int) Either[error, int] {
			return Right[error](v + 1)
		}

		result := F.Pipe1(
			Left[Counter](assert.AnError),
			BindL(valueLens, increment),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("BindL with multiple operations", func(t *testing.T) {
		double := func(v int) Either[error, int] {
			return Right[error](v * 2)
		}

		addTen := func(v int) Either[error, int] {
			return Right[error](v + 10)
		}

		result := F.Pipe2(
			Right[error](Counter{Value: 5}),
			BindL(valueLens, double),
			BindL(valueLens, addTen),
		)

		expected := Right[error](Counter{Value: 20})
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
		double := N.Mul(2)

		result := F.Pipe1(
			Right[error](Counter{Value: 21}),
			LetL[error](valueLens, double),
		)

		expected := Right[error](Counter{Value: 42})
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with Left input", func(t *testing.T) {
		double := N.Mul(2)

		result := F.Pipe1(
			Left[Counter](assert.AnError),
			LetL[error](valueLens, double),
		)

		expected := Left[Counter](assert.AnError)
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with multiple transformations", func(t *testing.T) {
		double := N.Mul(2)
		addTen := N.Add(10)

		result := F.Pipe2(
			Right[error](Counter{Value: 5}),
			LetL[error](valueLens, double),
			LetL[error](valueLens, addTen),
		)

		expected := Right[error](Counter{Value: 20})
		assert.Equal(t, expected, result)
	})

	t.Run("LetL with identity transformation", func(t *testing.T) {
		identity := F.Identity[int]

		result := F.Pipe1(
			Right[error](Counter{Value: 42}),
			LetL[error](valueLens, identity),
		)

		expected := Right[error](Counter{Value: 42})
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
			Right[error](Config{Debug: true, Timeout: 30}),
			LetToL[error](debugLens, false),
		)

		expected := Right[error](Config{Debug: false, Timeout: 30})
		assert.Equal(t, expected, result)
	})

	t.Run("LetToL with Left input", func(t *testing.T) {
		result := F.Pipe1(
			Left[Config](assert.AnError),
			LetToL[error](debugLens, false),
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
			Right[error](Config{Debug: true, Timeout: 30}),
			LetToL[error](debugLens, false),
			LetToL[error](timeoutLens, 60),
		)

		expected := Right[error](Config{Debug: false, Timeout: 60})
		assert.Equal(t, expected, result)
	})

	t.Run("LetToL setting same value", func(t *testing.T) {
		result := F.Pipe1(
			Right[error](Config{Debug: false, Timeout: 30}),
			LetToL[error](debugLens, false),
		)

		expected := Right[error](Config{Debug: false, Timeout: 30})
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
		double := N.Mul(2)

		result := F.Pipe2(
			Right[error](Counter{Value: 100}),
			LetToL[error](valueLens, 10),
			LetL[error](valueLens, double),
		)

		expected := Right[error](Counter{Value: 20})
		assert.Equal(t, expected, result)
	})

	t.Run("Combine LetL and BindL", func(t *testing.T) {
		double := N.Mul(2)
		validate := func(v int) Either[error, int] {
			if v > 100 {
				return Left[int](assert.AnError)
			}
			return Right[error](v)
		}

		result := F.Pipe2(
			Right[error](Counter{Value: 25}),
			LetL[error](valueLens, double),
			BindL(valueLens, validate),
		)

		expected := Right[error](Counter{Value: 50})
		assert.Equal(t, expected, result)
	})

	t.Run("Combine ApSL and LetL", func(t *testing.T) {
		addFive := func(v int) int { return v + 5 }

		result := F.Pipe2(
			Right[error](Counter{Value: 10}),
			ApSL(valueLens, Right[error](20)),
			LetL[error](valueLens, addFive),
		)

		expected := Right[error](Counter{Value: 25})
		assert.Equal(t, expected, result)
	})
}
