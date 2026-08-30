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

package pair_test

import (
	"fmt"
	"strings"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/common"
	LP "github.com/IBM/fp-go/v2/optics/lens/pair"
	P "github.com/IBM/fp-go/v2/pair"
)

// ExampleHead_get demonstrates reading the head element via Head.
func ExampleHead_get() {
	p := P.MakePair("hello", 42)
	headLens := LP.Head[string, int]()

	fmt.Println(headLens.Get(p))
	// Output:
	// hello
}

// ExampleHead_set demonstrates replacing the head element via Head.
func ExampleHead_set() {
	p := P.MakePair("hello", 42)
	headLens := LP.Head[string, int]()

	updated := headLens.Set("world")(p)
	fmt.Println(updated)
	// Output:
	// Pair[string, int](world, 42)
}

// ExampleHead_modify demonstrates modifying the head element via Head.
func ExampleHead_modify() {
	p := P.MakePair("hello", 42)
	headLens := LP.Head[string, int]()

	toUpper := F.Pipe1(headLens, L.LensModify[P.Pair[string, int]](strings.ToUpper))
	fmt.Println(toUpper(p))
	// Output:
	// Pair[string, int](HELLO, 42)
}

// ExampleTail_get demonstrates reading the tail element via Tail.
func ExampleTail_get() {
	p := P.MakePair("hello", 42)
	tailLens := LP.Tail[string, int]()

	fmt.Println(tailLens.Get(p))
	// Output:
	// 42
}

// ExampleTail_set demonstrates replacing the tail element via Tail.
func ExampleTail_set() {
	p := P.MakePair("hello", 42)
	tailLens := LP.Tail[string, int]()

	updated := tailLens.Set(100)(p)
	fmt.Println(updated)
	// Output:
	// Pair[string, int](hello, 100)
}

// ExampleTail_modify demonstrates modifying the tail element via Tail.
func ExampleTail_modify() {
	p := P.MakePair("hello", 42)
	tailLens := LP.Tail[string, int]()

	double := F.Pipe1(tailLens, L.LensModify[P.Pair[string, int]](N.Mul(2)))
	fmt.Println(double(p))
	// Output:
	// Pair[string, int](hello, 84)
}
