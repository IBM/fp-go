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
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) (string, error) {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) (string, error) {
	return Of("John")
}

func TestBind(t *testing.T) {

	res, err := Pipe4(
		utils.Empty,
		Do,
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	AssertEq(Of("John Doe"))(res, err)(t)
}

func TestApS(t *testing.T) {

	res, err := Pipe4(
		utils.Empty,
		Do,
		ApS(utils.SetLastName)(Of("Doe")),
		ApS(utils.SetGivenName)(Of("John")),
		Map(utils.GetFullName),
	)

	AssertEq(Of("John Doe"))(res, err)(t)
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
		result, err := Pipe2(
			Person{Name: "Alice", Age: 25},
			Right,
			ApSL(ageLens)(Right(30)),
		)

		AssertEq(Right(Person{Name: "Alice", Age: 30}))(result, err)(t)
	})

	t.Run("ApSL with Left in context", func(t *testing.T) {
		result, err := Pipe2(
			assert.AnError,
			Left[Person],
			ApSL(ageLens)(Right(30)),
		)

		AssertEq(Left[Person](assert.AnError))(result, err)(t)
	})

	t.Run("ApSL with Left in value", func(t *testing.T) {
		result, err := Pipe2(
			Person{Name: "Alice", Age: 25},
			Right,
			ApSL(ageLens)(Left[int](assert.AnError)),
		)

		AssertEq(Left[Person](assert.AnError))(result, err)(t)
	})

	t.Run("ApSL with both Left", func(t *testing.T) {
		result, err := Pipe2(
			assert.AnError,
			Left[Person],
			ApSL(ageLens)(Left[int](assert.AnError)),
		)

		AssertEq(Left[Person](assert.AnError))(result, err)(t)
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
		increment := func(v int) (int, error) {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right(v + 1)
		}

		result, err := Pipe2(
			Counter{Value: 42},
			Right,
			BindL(valueLens, increment),
		)

		AssertEq(Right(Counter{Value: 43}))(result, err)(t)
	})

	t.Run("BindL with failing transformation", func(t *testing.T) {
		increment := func(v int) (int, error) {
			if v >= 100 {
				return Left[int](assert.AnError)
			}
			return Right(v + 1)
		}

		result, err := Pipe2(
			Counter{Value: 100},
			Right,
			BindL(valueLens, increment),
		)

		AssertEq(Left[Counter](assert.AnError))(result, err)(t)
	})

	t.Run("BindL with Left input", func(t *testing.T) {
		increment := func(v int) (int, error) {
			return Right(v + 1)
		}

		result, err := Pipe2(
			assert.AnError,
			Left[Counter],
			BindL(valueLens, increment),
		)

		AssertEq(Left[Counter](assert.AnError))(result, err)(t)
	})

	t.Run("BindL with multiple operations", func(t *testing.T) {
		double := func(v int) (int, error) {
			return Right(v * 2)
		}

		addTen := func(v int) (int, error) {
			return Right(v + 10)
		}

		result, err := Pipe3(
			Counter{Value: 5},
			Right,
			BindL(valueLens, double),
			BindL(valueLens, addTen),
		)

		AssertEq(Right(Counter{Value: 20}))(result, err)(t)
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

		result, err := Pipe2(
			Counter{Value: 21},
			Right,
			LetL(valueLens, double),
		)
		AssertEq(Right(Counter{Value: 42}))(result, err)(t)
	})

	t.Run("LetL with Left input", func(t *testing.T) {
		double := N.Mul(2)

		result, err := Pipe2(
			assert.AnError,
			Left[Counter],
			LetL(valueLens, double),
		)

		AssertEq(Left[Counter](assert.AnError))(result, err)(t)
	})

	t.Run("LetL with multiple transformations", func(t *testing.T) {
		double := N.Mul(2)
		addTen := N.Add(10)

		result, err := Pipe3(
			Counter{Value: 5},
			Right,
			LetL(valueLens, double),
			LetL(valueLens, addTen),
		)

		AssertEq(Right(Counter{Value: 20}))(result, err)(t)
	})

	t.Run("LetL with identity transformation", func(t *testing.T) {
		identity := F.Identity[int]

		result, err := Pipe2(
			Counter{Value: 42},
			Right,
			LetL(valueLens, identity),
		)

		AssertEq(Right(Counter{Value: 42}))(result, err)(t)
	})
}

func TestLetToL(t *testing.T) {
	// Create a lens for the Debug field
	debugLens := L.MakeLens(
		func(c Config) bool { return c.Debug },
		func(c Config, d bool) Config { c.Debug = d; return c },
	)

	t.Run("LetToL with constant value", func(t *testing.T) {
		result, err := Pipe2(
			Config{Debug: true, Timeout: 30},
			Right,
			LetToL(debugLens, false),
		)

		AssertEq(Right(Config{Debug: false, Timeout: 30}))(result, err)(t)
	})

	t.Run("LetToL with Left input", func(t *testing.T) {
		result, err := Pipe2(
			assert.AnError,
			Left[Config],
			LetToL(debugLens, false),
		)

		AssertEq(Left[Config](assert.AnError))(result, err)(t)
	})

	t.Run("LetToL with multiple fields", func(t *testing.T) {
		timeoutLens := L.MakeLens(
			func(c Config) int { return c.Timeout },
			func(c Config, t int) Config { c.Timeout = t; return c },
		)

		result, err := Pipe3(
			Config{Debug: true, Timeout: 30},
			Right,
			LetToL(debugLens, false),
			LetToL(timeoutLens, 60),
		)

		AssertEq(Right(Config{Debug: false, Timeout: 60}))(result, err)(t)
	})

	t.Run("LetToL setting same value", func(t *testing.T) {
		result, err := Pipe2(
			Config{Debug: false, Timeout: 30},
			Right,
			LetToL(debugLens, false),
		)

		AssertEq(Right(Config{Debug: false, Timeout: 30}))(result, err)(t)
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

		result, err := Pipe3(
			Counter{Value: 100},
			Right,
			LetToL(valueLens, 10),
			LetL(valueLens, double),
		)

		AssertEq(Right(Counter{Value: 20}))(result, err)(t)
	})

	t.Run("Combine LetL and BindL", func(t *testing.T) {
		double := N.Mul(2)
		validate := func(v int) (int, error) {
			if v > 100 {
				return Left[int](assert.AnError)
			}
			return Right(v)
		}

		result, err := Pipe3(
			Counter{Value: 25},
			Right,
			LetL(valueLens, double),
			BindL(valueLens, validate),
		)

		AssertEq(Right(Counter{Value: 50}))(result, err)(t)
	})

	t.Run("Combine ApSL and LetL", func(t *testing.T) {
		addFive := func(v int) int { return v + 5 }

		result, err := Pipe3(
			Counter{Value: 10},
			Right,
			ApSL(valueLens)(Right(20)),
			LetL(valueLens, addFive),
		)

		AssertEq(Right(Counter{Value: 25}))(result, err)(t)
	})
}
