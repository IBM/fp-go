// Copyright (c) 2025 IBM Corp.
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

package iter

import (
	"fmt"
	"slices"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// Test types
type User struct {
	Name string
	Age  int
}

type State struct {
	Value  int
	Double int
	Status string
}

// TestDo tests the Do function
func TestDo(t *testing.T) {
	t.Run("creates sequence with single element", func(t *testing.T) {
		result := Do(42)
		values := slices.Collect(result)
		assert.Equal(t, A.Of(42), values)
	})

	t.Run("creates sequence with struct", func(t *testing.T) {
		user := User{Name: "Alice", Age: 30}
		result := Do(user)
		values := slices.Collect(result)
		assert.Equal(t, A.Of(user), values)
	})

	t.Run("creates sequence with zero value", func(t *testing.T) {
		result := Do(State{})
		values := slices.Collect(result)
		assert.Equal(t, []State{{Value: 0, Double: 0, Status: ""}}, values)
	})
}

// TestBind tests the Bind function
func TestBind(t *testing.T) {
	t.Run("binds simple value", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		getValues := func(s State) Seq[int] {
			return From(1, 2, 3)
		}

		bindOp := Bind(setValue, getValues)
		result := bindOp(Do(State{}))

		values := slices.Collect(result)
		expected := []State{
			{Value: 1},
			{Value: 2},
			{Value: 3},
		}
		assert.Equal(t, expected, values)
	})

	t.Run("chains multiple binds", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		getValues := func(s State) Seq[int] {
			return From(5)
		}

		computeDouble := func(s State) Seq[int] {
			return From(s.Value * 2)
		}

		result := F.Flow2(
			Bind(setValue, getValues),
			Bind(setDouble, computeDouble),
		)(Do(State{}))

		values := slices.Collect(result)
		expected := []State{{Value: 5, Double: 10}}
		assert.Equal(t, expected, values)
	})

	t.Run("binds with multiple results", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		multiplyValues := func(s State) Seq[int] {
			return From(s.Value, s.Value*2, s.Value*3)
		}

		bindOp := Bind(setValue, multiplyValues)
		result := bindOp(Do(State{Value: 2}))

		values := slices.Collect(result)
		expected := []State{
			{Value: 2},
			{Value: 4},
			{Value: 6},
		}
		assert.Equal(t, expected, values)
	})

	t.Run("binds with empty sequence", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		emptySeq := func(s State) Seq[int] {
			return Empty[int]()
		}

		bindOp := Bind(setValue, emptySeq)
		result := bindOp(Do(State{Value: 5}))

		values := slices.Collect(result)
		assert.Empty(t, values)
	})
}

// TestLet tests the Let function
func TestLet(t *testing.T) {
	t.Run("computes value from state", func(t *testing.T) {
		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		computeDouble := func(s State) int {
			return s.Value * 2
		}

		letOp := Let(setDouble, computeDouble)
		result := letOp(Do(State{Value: 5}))

		values := slices.Collect(result)
		expected := []State{{Value: 5, Double: 10}}
		assert.Equal(t, expected, values)
	})

	t.Run("chains multiple lets", func(t *testing.T) {
		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		computeDouble := func(s State) int {
			return s.Value * 2
		}

		computeValue := func(s State) int {
			return s.Double / 2
		}

		result := F.Flow2(
			Let(setDouble, computeDouble),
			Let(setValue, computeValue),
		)(Do(State{Value: 7}))

		values := slices.Collect(result)
		expected := []State{{Value: 7, Double: 14}}
		assert.Equal(t, expected, values)
	})

	t.Run("computes complex transformation", func(t *testing.T) {
		setStatus := func(s string) func(State) State {
			return func(st State) State {
				st.Status = s
				return st
			}
		}

		computeStatus := func(s State) string {
			if s.Value > 10 {
				return "high"
			}
			return "low"
		}

		letOp := Let(setStatus, computeStatus)
		result := letOp(Do(State{Value: 15}))

		values := slices.Collect(result)
		expected := []State{{Value: 15, Status: "high"}}
		assert.Equal(t, expected, values)
	})
}

// TestLetTo tests the LetTo function
func TestLetTo(t *testing.T) {
	t.Run("sets constant value", func(t *testing.T) {
		setStatus := func(s string) func(State) State {
			return func(st State) State {
				st.Status = s
				return st
			}
		}

		letToOp := LetTo(setStatus, "active")
		result := letToOp(Do(State{Value: 5}))

		values := slices.Collect(result)
		expected := []State{{Value: 5, Status: "active"}}
		assert.Equal(t, expected, values)
	})

	t.Run("chains multiple LetTo calls", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setStatus := func(st string) func(State) State {
			return func(s State) State {
				s.Status = st
				return s
			}
		}

		result := F.Flow2(
			LetTo(setValue, 42),
			LetTo(setStatus, "ready"),
		)(Do(State{}))

		values := slices.Collect(result)
		expected := []State{{Value: 42, Status: "ready"}}
		assert.Equal(t, expected, values)
	})

	t.Run("sets zero value", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		letToOp := LetTo(setValue, 0)
		result := letToOp(Do(State{Value: 100}))

		values := slices.Collect(result)
		expected := []State{{Value: 0}}
		assert.Equal(t, expected, values)
	})
}

// TestBindTo tests the BindTo function
func TestBindTo(t *testing.T) {
	t.Run("wraps values into structure", func(t *testing.T) {
		createState := func(v int) State {
			return State{Value: v}
		}

		bindToOp := BindTo(createState)
		result := bindToOp(From(1, 2, 3))

		values := slices.Collect(result)
		expected := []State{
			{Value: 1},
			{Value: 2},
			{Value: 3},
		}
		assert.Equal(t, expected, values)
	})

	t.Run("wraps into complex structure", func(t *testing.T) {
		createUser := func(name string) User {
			return User{Name: name, Age: 0}
		}

		bindToOp := BindTo(createUser)
		result := bindToOp(From("Alice", "Bob", "Charlie"))

		values := slices.Collect(result)
		expected := []User{
			{Name: "Alice", Age: 0},
			{Name: "Bob", Age: 0},
			{Name: "Charlie", Age: 0},
		}
		assert.Equal(t, expected, values)
	})

	t.Run("wraps empty sequence", func(t *testing.T) {
		createState := func(v int) State {
			return State{Value: v}
		}

		bindToOp := BindTo(createState)
		result := bindToOp(Empty[int]())

		values := slices.Collect(result)
		assert.Empty(t, values)
	})
}

// TestApS tests the ApS function
func TestApS(t *testing.T) {
	t.Run("applies sequence of values", func(t *testing.T) {
		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		doubles := From(10, 20, 30)

		apOp := ApS(setDouble, doubles)
		result := apOp(Do(State{Value: 5}))

		values := slices.Collect(result)
		expected := []State{
			{Value: 5, Double: 10},
			{Value: 5, Double: 20},
			{Value: 5, Double: 30},
		}
		assert.Equal(t, expected, values)
	})

	t.Run("applies with empty sequence", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		apOp := ApS(setValue, Empty[int]())
		result := apOp(Do(State{Value: 5}))

		values := slices.Collect(result)
		assert.Empty(t, values)
	})

	t.Run("chains multiple ApS calls", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		values := From(1, 2)
		doubles := From(10, 20)

		result := F.Flow2(
			ApS(setValue, values),
			ApS(setDouble, doubles),
		)(Do(State{}))

		results := slices.Collect(result)
		// Cartesian product: 2 values × 2 doubles = 4 results
		assert.Len(t, results, 4)
	})
}

// TestDoNotationChain tests a complete do-notation chain
func TestDoNotationChain(t *testing.T) {
	t.Run("complex do-notation chain", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		setStatus := func(st string) func(State) State {
			return func(s State) State {
				s.Status = st
				return s
			}
		}

		getValues := func(s State) Seq[int] {
			return From(5, 10)
		}

		computeDouble := func(s State) int {
			return s.Value * 2
		}

		result := F.Flow3(
			Bind(setValue, getValues),
			Let(setDouble, computeDouble),
			LetTo(setStatus, "computed"),
		)(Do(State{}))

		results := slices.Collect(result)
		expected := []State{
			{Value: 5, Double: 10, Status: "computed"},
			{Value: 10, Double: 20, Status: "computed"},
		}
		assert.Equal(t, expected, results)
	})

	t.Run("mixed bind and let operations", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		getInitial := func(s State) Seq[int] {
			return From(3)
		}

		multiplyValue := func(s State) Seq[int] {
			return From(s.Value*2, s.Value*3)
		}

		result := F.Flow2(
			Bind(setValue, getInitial),
			Bind(setDouble, multiplyValue),
		)(Do(State{}))

		results := slices.Collect(result)
		expected := []State{
			{Value: 3, Double: 6},
			{Value: 3, Double: 9},
		}
		assert.Equal(t, expected, results)
	})
}

// TestEdgeCases tests edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("bind with single element", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		getSingle := func(s State) Seq[int] {
			return From(42)
		}

		bindOp := Bind(setValue, getSingle)
		result := bindOp(Do(State{}))

		results := slices.Collect(result)
		expected := []State{{Value: 42}}
		assert.Equal(t, expected, results)
	})

	t.Run("multiple binds with cartesian product", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		setDouble := func(d int) func(State) State {
			return func(s State) State {
				s.Double = d
				return s
			}
		}

		getValues := func(s State) Seq[int] {
			return From(1, 2)
		}

		getDoubles := func(s State) Seq[int] {
			return From(10, 20)
		}

		result := F.Flow2(
			Bind(setValue, getValues),
			Bind(setDouble, getDoubles),
		)(Do(State{}))

		results := slices.Collect(result)
		// Should produce cartesian product: 2 × 2 = 4 results
		assert.Len(t, results, 4)
	})

	t.Run("let with identity function", func(t *testing.T) {
		setValue := func(v int) func(State) State {
			return func(s State) State {
				s.Value = v
				return s
			}
		}

		identity := func(s State) int {
			return s.Value
		}

		letOp := Let(setValue, identity)
		result := letOp(Do(State{Value: 99}))

		results := slices.Collect(result)
		expected := []State{{Value: 99}}
		assert.Equal(t, expected, results)
	})
}

// Benchmark tests
func BenchmarkBind(b *testing.B) {
	setValue := func(v int) func(State) State {
		return func(s State) State {
			s.Value = v
			return s
		}
	}

	getValues := func(s State) Seq[int] {
		return From(1, 2, 3, 4, 5)
	}

	bindOp := Bind(setValue, getValues)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := bindOp(Do(State{}))
		// Consume the sequence
		for range result {
		}
	}
}

func BenchmarkLet(b *testing.B) {
	setDouble := func(d int) func(State) State {
		return func(s State) State {
			s.Double = d
			return s
		}
	}

	computeDouble := func(s State) int {
		return s.Value * 2
	}

	letOp := Let(setDouble, computeDouble)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := letOp(Do(State{Value: 5}))
		// Consume the sequence
		for range result {
		}
	}
}

func BenchmarkDoNotationChain(b *testing.B) {
	setValue := func(v int) func(State) State {
		return func(s State) State {
			s.Value = v
			return s
		}
	}

	setDouble := func(d int) func(State) State {
		return func(s State) State {
			s.Double = d
			return s
		}
	}

	setStatus := func(st string) func(State) State {
		return func(s State) State {
			s.Status = st
			return s
		}
	}

	getValues := func(s State) Seq[int] {
		return From(5, 10, 15)
	}

	computeDouble := func(s State) int {
		return s.Value * 2
	}

	chain := F.Flow3(
		Bind(setValue, getValues),
		Let(setDouble, computeDouble),
		LetTo(setStatus, "computed"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := chain(Do(State{}))
		// Consume the sequence
		for range result {
		}
	}
}

// Example tests for documentation
func ExampleDo() {
	result := Do(42)
	for v := range result {
		fmt.Println(v)
	}
	// Output: 42
}

func ExampleBind() {
	setValue := func(v int) func(State) State {
		return func(s State) State {
			s.Value = v
			return s
		}
	}

	getValues := func(s State) Seq[int] {
		return From(1, 2, 3)
	}

	bindOp := Bind(setValue, getValues)
	result := bindOp(Do(State{}))

	for s := range result {
		fmt.Printf("Value: %d\n", s.Value)
	}
	// Output:
	// Value: 1
	// Value: 2
	// Value: 3
}

func ExampleLet() {
	setDouble := func(d int) func(State) State {
		return func(s State) State {
			s.Double = d
			return s
		}
	}

	computeDouble := func(s State) int {
		return s.Value * 2
	}

	letOp := Let(setDouble, computeDouble)
	result := letOp(Do(State{Value: 5}))

	for s := range result {
		fmt.Printf("Value: %d, Double: %d\n", s.Value, s.Double)
	}
	// Output: Value: 5, Double: 10
}

func ExampleLetTo() {
	setStatus := func(s string) func(State) State {
		return func(st State) State {
			st.Status = s
			return st
		}
	}

	letToOp := LetTo(setStatus, "active")
	result := letToOp(Do(State{Value: 5}))

	for s := range result {
		fmt.Printf("Value: %d, Status: %s\n", s.Value, s.Status)
	}
	// Output: Value: 5, Status: active
}
