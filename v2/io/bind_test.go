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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) IO[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) IO[string] {
	return Of("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(), "John Doe")
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(), "John Doe")
}

// Test types for lens-based operations
type Counter struct {
	Value int
}

type Person struct {
	Name string
	Age  int
}

func TestBindL(t *testing.T) {
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("BindL with successful transformation", func(t *testing.T) {
		increment := func(v int) IO[int] {
			return Of(v + 1)
		}

		result := F.Pipe1(
			Of(Counter{Value: 42}),
			BindL(valueLens, increment),
		)

		assert.Equal(t, Counter{Value: 43}, result())
	})

	t.Run("BindL with multiple operations", func(t *testing.T) {
		double := func(v int) IO[int] {
			return Of(v * 2)
		}

		addTen := func(v int) IO[int] {
			return Of(v + 10)
		}

		result := F.Pipe2(
			Of(Counter{Value: 5}),
			BindL(valueLens, double),
			BindL(valueLens, addTen),
		)

		assert.Equal(t, Counter{Value: 20}, result())
	})
}

func TestLetL(t *testing.T) {
	valueLens := L.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	t.Run("LetL with pure transformation", func(t *testing.T) {
		double := func(v int) int { return v * 2 }

		result := F.Pipe1(
			Of(Counter{Value: 21}),
			LetL(valueLens, double),
		)

		assert.Equal(t, Counter{Value: 42}, result())
	})

	t.Run("LetL with multiple transformations", func(t *testing.T) {
		double := func(v int) int { return v * 2 }
		addTen := func(v int) int { return v + 10 }

		result := F.Pipe2(
			Of(Counter{Value: 5}),
			LetL(valueLens, double),
			LetL(valueLens, addTen),
		)

		assert.Equal(t, Counter{Value: 20}, result())
	})
}

func TestLetToL(t *testing.T) {
	ageLens := L.MakeLens(
		func(p Person) int { return p.Age },
		func(p Person, a int) Person { p.Age = a; return p },
	)

	t.Run("LetToL with constant value", func(t *testing.T) {
		result := F.Pipe1(
			Of(Person{Name: "Alice", Age: 25}),
			LetToL(ageLens, 30),
		)

		assert.Equal(t, Person{Name: "Alice", Age: 30}, result())
	})

	t.Run("LetToL with multiple fields", func(t *testing.T) {
		nameLens := L.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, n string) Person { p.Name = n; return p },
		)

		result := F.Pipe2(
			Of(Person{Name: "Alice", Age: 25}),
			LetToL(ageLens, 30),
			LetToL(nameLens, "Bob"),
		)

		assert.Equal(t, Person{Name: "Bob", Age: 30}, result())
	})
}

func TestApSL(t *testing.T) {
	ageLens := L.MakeLens(
		func(p Person) int { return p.Age },
		func(p Person, a int) Person { p.Age = a; return p },
	)

	t.Run("ApSL with value", func(t *testing.T) {
		result := F.Pipe1(
			Of(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Of(30)),
		)

		assert.Equal(t, Person{Name: "Alice", Age: 30}, result())
	})

	t.Run("ApSL with chaining", func(t *testing.T) {
		nameLens := L.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, n string) Person { p.Name = n; return p },
		)

		result := F.Pipe2(
			Of(Person{Name: "Alice", Age: 25}),
			ApSL(ageLens, Of(30)),
			ApSL(nameLens, Of("Bob")),
		)

		assert.Equal(t, Person{Name: "Bob", Age: 30}, result())
	})
}
