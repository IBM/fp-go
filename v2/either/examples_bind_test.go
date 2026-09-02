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
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/internal/common"
)

// Counter is a simple structure used in the lens-based do-notation examples.
type Counter struct {
	Value int
}

var counterLens = L.MakeLens(
	func(c Counter) int { return c.Value },
	func(c Counter, v int) Counter { c.Value = v; return c },
)

// ExampleApSL demonstrates using ApSL to set a field via a lens in applicative style.
// This is equivalent to ApS but uses a Lens rather than a manually-written curried setter.
func ExampleApSL() {
	// Use ApSL to set the Value field of Counter using a lens.
	result := F.Pipe1(
		Of[error](Counter{}),
		ApSL(counterLens, Of[error](42)),
	)
	fmt.Println(result)

	// Output:
	// Right[either.Counter]({42})
}

// ExampleApSL_leftContext demonstrates that ApSL short-circuits when the context is Left.
func ExampleApSL_leftContext() {
	result := F.Pipe1(
		Left[Counter](fmt.Errorf("context error")),
		ApSL(counterLens, Of[error](42)),
	)
	fmt.Println(IsLeft(result))

	// Output:
	// true
}

// ExampleBindL demonstrates using BindL to transform a focused field monadically.
// The lens reads the current field value, passes it to f, and writes the new value back.
func ExampleBindL() {
	increment := func(v int) Either[error, int] {
		return Of[error](v + 1)
	}

	result := F.Pipe1(
		Of[error](Counter{Value: 41}),
		BindL(counterLens, increment),
	)
	fmt.Println(result)

	// Output:
	// Right[either.Counter]({42})
}

// ExampleBindL_error demonstrates that BindL propagates an error returned by f.
func ExampleBindL_error() {
	alwaysFail := func(_ int) Either[error, int] {
		return Left[int](fmt.Errorf("overflow"))
	}

	result := F.Pipe1(
		Of[error](Counter{Value: 100}),
		BindL(counterLens, alwaysFail),
	)
	fmt.Println(IsLeft(result))

	// Output:
	// true
}

// ExampleLetL demonstrates using LetL to apply a pure transformation to a focused field.
func ExampleLetL() {
	double := func(v int) int { return v * 2 }

	result := F.Pipe1(
		Of[error](Counter{Value: 21}),
		LetL[error](counterLens, double),
	)
	fmt.Println(result)

	// Output:
	// Right[either.Counter]({42})
}

// ExampleLetToL demonstrates using LetToL to replace a focused field with a constant value.
func ExampleLetToL() {
	result := F.Pipe1(
		Of[error](Counter{Value: 0}),
		LetToL[error](counterLens, 99),
	)
	fmt.Println(result)

	// Output:
	// Right[either.Counter]({99})
}
