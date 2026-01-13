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

package tuple_test

import (
	"fmt"

	S "github.com/IBM/fp-go/v2/string"
	"github.com/IBM/fp-go/v2/tuple"
)

func ExampleOf() {
	t := tuple.Of(42)
	fmt.Println(t)
	// Output:
	// Tuple1[int](42)
}

func ExampleFirst() {
	t := tuple.MakeTuple2("hello", 42)
	s := tuple.First(t)
	fmt.Println(s)
	// Output:
	// hello
}

func ExampleSecond() {
	t := tuple.MakeTuple2("hello", 42)
	n := tuple.Second(t)
	fmt.Println(n)
	// Output:
	// 42
}

func ExampleSwap() {
	t := tuple.MakeTuple2("hello", 42)
	swapped := tuple.Swap(t)
	fmt.Println(swapped)
	// Output:
	// Tuple2[int, string](42, hello)
}

func ExampleOf2() {
	pairWith42 := tuple.Of2[string](42)
	t := pairWith42("hello")
	fmt.Println(t)
	// Output:
	// Tuple2[string, int](hello, 42)
}

func ExampleBiMap() {
	t := tuple.MakeTuple2(5, "hello")
	mapper := tuple.BiMap(
		S.Size,
		func(n int) string { return fmt.Sprintf("%d", n*2) },
	)
	result := mapper(t)
	fmt.Println(result)
	// Output:
	// Tuple2[string, int](10, 5)
}
