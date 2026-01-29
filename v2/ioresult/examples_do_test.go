// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"fmt"
	"log"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	T "github.com/IBM/fp-go/v2/tuple"
)

func ExampleIOEither_do() {
	foo := Of("foo")
	bar := Of(1)

	// quux consumes the state of three bindings and returns an [IO] instead of an [IOEither]
	quux := func(t T.Tuple3[string, int, string]) IO[Void] {
		return io.FromImpure(func() {
			log.Printf("t1: %s, t2: %d, t3: %s", t.F1, t.F2, t.F3)
		})
	}

	transform := func(t T.Tuple3[string, int, string]) int {
		return len(t.F1) + t.F2 + len(t.F3)
	}

	b := F.Pipe5(
		foo,
		BindTo(T.Of[string]),
		ApS(T.Push1[string, int], bar),
		Bind(T.Push2[string, int, string], func(t T.Tuple2[string, int]) IOResult[string] {
			return Of(fmt.Sprintf("%s%d", t.F1, t.F2))
		}),
		ChainFirstIOK(quux),
		Map(transform),
	)

	fmt.Println(b())

	// Output:
	// Right[int](8)
}
