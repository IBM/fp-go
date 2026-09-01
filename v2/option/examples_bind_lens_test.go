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

package option

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/internal/common"
)

type counter struct {
	Value int
	Label string
}

var valueLens = L.MakeLens(
	func(c counter) int { return c.Value },
	func(c counter, v int) counter { c.Value = v; return c },
)

// ExampleApSL demonstrates applying an Option value to a struct field via a lens.
func ExampleApSL() {
	result := F.Pipe1(
		Do(counter{Label: "score"}),
		ApSL(valueLens, Some(42)),
	)
	fmt.Println(result)
	// Output:
	// Some[option.counter]({42 score})
}

// ExampleApSL_none shows that None propagates when the provided Option is None.
func ExampleApSL_none() {
	result := F.Pipe1(
		Do(counter{Label: "score"}),
		ApSL(valueLens, None[int]()),
	)
	fmt.Println(result)
	// Output:
	// None[option.counter]
}

// ExampleBindL demonstrates applying a monadic computation to the field focused by a lens.
func ExampleBindL() {
	increment := func(v int) Option[int] {
		if v < 100 {
			return Some(v + 1)
		}
		return None[int]()
	}

	result := F.Pipe1(
		Some(counter{Value: 41, Label: "score"}),
		BindL(valueLens, increment),
	)
	fmt.Println(result)
	// Output:
	// Some[option.counter]({42 score})
}

// ExampleBindL_none shows that None is returned when the computation fails.
func ExampleBindL_none() {
	increment := func(v int) Option[int] {
		if v < 100 {
			return Some(v + 1)
		}
		return None[int]()
	}

	result := F.Pipe1(
		Some(counter{Value: 100, Label: "score"}),
		BindL(valueLens, increment),
	)
	fmt.Println(result)
	// Output:
	// None[option.counter]
}

// ExampleLetL demonstrates applying a pure transformation to the field focused by a lens.
func ExampleLetL() {
	result := F.Pipe1(
		Some(counter{Value: 21, Label: "score"}),
		LetL(valueLens, func(v int) int { return v * 2 }),
	)
	fmt.Println(result)
	// Output:
	// Some[option.counter]({42 score})
}

// ExampleLetToL demonstrates replacing the field focused by a lens with a constant value.
func ExampleLetToL() {
	result := F.Pipe1(
		Some(counter{Value: 0, Label: "score"}),
		LetToL(valueLens, 100),
	)
	fmt.Println(result)
	// Output:
	// Some[option.counter]({100 score})
}
