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

package itereither

import (
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
)

// ExampleMonadTapIOK runs an IO side effect for every Right element.
func ExampleMonadTapIOK() {
	seq := iter.From(E.Right[string](7), E.Left[int]("err"))
	for e := range MonadTapIOK(seq, func(n int) func() int {
		return func() int { fmt.Println("log:", n); return n }
	}) {
		fmt.Println(e)
	}
	// Output:
	// log: 7
	// Right[int](7)
	// Left[string](err)
}

// ExampleMonadTapIOEitherK runs a fallible side effect for every Right element;
// a failing side effect replaces the element.
func ExampleMonadTapIOEitherK() {
	check := func(n int) func() E.Either[string, int] {
		return func() E.Either[string, int] {
			if n > 0 {
				return E.Right[string](n)
			}
			return E.Left[int]("not positive")
		}
	}
	for e := range MonadTapIOEitherK(iter.From(E.Right[string](7), E.Right[string](-7)), check) {
		fmt.Println(e)
	}
	// Output:
	// Right[int](7)
	// Left[string](not positive)
}
