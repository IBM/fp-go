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

package result_test

import (
	"errors"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	R "github.com/IBM/fp-go/v2/result"
)

func ExampleDo() {
	type State struct {
		x, y int
	}
	result := R.Do(State{})
	fmt.Println(R.IsRight(result))
	// Output:
	// true
}

func ExampleBind() {
	type State struct {
		value int
	}
	result := F.Pipe1(
		R.Do(State{}),
		R.Bind(
			func(v int) func(State) State {
				return func(s State) State { return State{value: v} }
			},
			func(s State) R.Result[int] {
				return R.Of(42)
			},
		),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](State{}))(result))
	// Output:
	// {42}
}

func ExampleLet() {
	type State struct {
		value int
	}
	result := F.Pipe1(
		R.Of(State{value: 10}),
		R.Let(
			func(v int) func(State) State {
				return func(s State) State { return State{value: s.value + v} }
			},
			func(s State) int { return 32 },
		),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](State{}))(result))
	// Output:
	// {42}
}

func ExampleLetTo() {
	type State struct {
		name string
	}
	result := F.Pipe1(
		R.Of(State{}),
		R.LetTo(
			func(n string) func(State) State {
				return func(s State) State { return State{name: n} }
			},
			"Alice",
		),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](State{}))(result))
	// Output:
	// {Alice}
}

func ExampleBindTo() {
	type State struct {
		value int
	}
	result := F.Pipe1(
		R.Of(42),
		R.BindTo(func(v int) State { return State{value: v} }),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](State{}))(result))
	// Output:
	// {42}
}

func ExampleApS() {
	type State struct {
		x, y int
	}
	result := F.Pipe1(
		R.Of(State{x: 10}),
		R.ApS(
			func(y int) func(State) State {
				return func(s State) State { return State{x: s.x, y: y} }
			},
			R.Of(32),
		),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](State{}))(result))
	// Output:
	// {10 32}
}

func ExampleApSL() {
	type Person struct {
		Name string
		Age  int
	}

	ageLens := lens.MakeLens(
		func(p Person) int { return p.Age },
		func(p Person, a int) Person { p.Age = a; return p },
	)

	result := F.Pipe1(
		R.Of(Person{Name: "Alice", Age: 25}),
		R.ApSL(ageLens, R.Of(30)),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](Person{}))(result))
	// Output:
	// {Alice 30}
}

func ExampleBindL() {
	type Counter struct {
		Value int
	}

	valueLens := lens.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	// Increment the counter, but fail if it would exceed 100
	increment := func(v int) R.Result[int] {
		if v >= 100 {
			return R.Left[int](errors.New("counter overflow"))
		}
		return R.Of(v + 1)
	}

	result := F.Pipe1(
		R.Of(Counter{Value: 42}),
		R.BindL(valueLens, increment),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](Counter{}))(result))
	// Output:
	// {43}
}

func ExampleLetL() {
	type Counter struct {
		Value int
	}

	valueLens := lens.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter { c.Value = v; return c },
	)

	// Double the counter value
	double := func(v int) int { return v * 2 }

	result := F.Pipe1(
		R.Of(Counter{Value: 21}),
		R.LetL(valueLens, double),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](Counter{}))(result))
	// Output:
	// {42}
}

func ExampleLetToL() {
	type Config struct {
		Debug   bool
		Timeout int
	}

	debugLens := lens.MakeLens(
		func(c Config) bool { return c.Debug },
		func(c Config, d bool) Config { c.Debug = d; return c },
	)

	result := F.Pipe1(
		R.Of(Config{Debug: true, Timeout: 30}),
		R.LetToL(debugLens, false),
	)
	fmt.Println(R.GetOrElse(F.Constant1[error](Config{}))(result))
	// Output:
	// {false 30}
}
