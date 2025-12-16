// Copyright (c) 2024 - 2025 IBM Corp.
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

package pair_test

import (
	"fmt"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/tuple"
)

func ExampleOf() {
	p := pair.Of(42)
	fmt.Println(p)
	// Output:
	// Pair[int, int](42, 42)
}

func ExampleFromTuple() {
	t := tuple.MakeTuple2("hello", 42)
	p := pair.FromTuple(t)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleFromHead() {
	makePair := pair.FromHead[int]("hello")
	p := makePair(42)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleFromTail() {
	makePair := pair.FromTail[string](42)
	p := makePair("hello")
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleToTuple() {
	p := pair.MakePair("hello", 42)
	t := pair.ToTuple(p)
	fmt.Println(t)
	// Output:
	// Tuple2[string, int](hello, 42)
}

func ExampleMakePair() {
	p := pair.MakePair("hello", 42)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleHead() {
	p := pair.MakePair("hello", 42)
	h := pair.Head(p)
	fmt.Println(h)
	// Output:
	// hello
}

func ExampleTail() {
	p := pair.MakePair("hello", 42)
	t := pair.Tail(p)
	fmt.Println(t)
	// Output:
	// 42
}

func ExampleFirst() {
	p := pair.MakePair("hello", 42)
	f := pair.First(p)
	fmt.Println(f)
	// Output:
	// hello
}

func ExampleSecond() {
	p := pair.MakePair("hello", 42)
	s := pair.Second(p)
	fmt.Println(s)
	// Output:
	// 42
}

func ExampleSwap() {
	p := pair.MakePair("hello", 42)
	swapped := pair.Swap(p)
	fmt.Println(swapped)
	// Output:
	// Pair[int, string](42, hello)
}

func ExampleMap() {
	p := pair.MakePair("hello", 42)
	doubled := pair.Map[string](N.Mul(2))(p)
	fmt.Println(doubled)
	// Output:
	// Pair[string, int](hello, 84)
}

func ExampleMapHead() {
	p := pair.MakePair("hello", 42)
	upper := pair.MapHead[int](func(s string) string { return s + "!" })(p)
	fmt.Println(upper)
	// Output:
	// Pair[string, int](hello!, 42)
}

func ExampleBiMap() {
	p := pair.MakePair("hello", 42)
	result := pair.BiMap(
		func(s string) string { return s + "!" },
		N.Mul(2),
	)(p)
	fmt.Println(result)
	// Output:
	// Pair[string, int](hello!, 84)
}
