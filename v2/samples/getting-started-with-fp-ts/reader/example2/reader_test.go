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

// Package example2 implements the second example in [https://dev.to/gcanti/getting-started-with-fp-ts-reader-1ie5]
package example2

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	ID "github.com/IBM/fp-go/v2/identity"
	N "github.com/IBM/fp-go/v2/number"
	I "github.com/IBM/fp-go/v2/number/integer"
	"github.com/IBM/fp-go/v2/ord"
	R "github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/string"
)

type (
	I18n struct {
		True  string
		False string
	}

	Dependencies struct {
		I18n       I18n
		LowerBound int
	}
)

var (
	// h: func(string) R.Reader[*Dependencies, string], note how the implementation does not depend on the dependencies
	h = F.Flow3(
		S.Size,
		N.Add(1),
		g,
	)
)

func getLowerBound(deps *Dependencies) int {
	return deps.LowerBound
}

func g(n int) R.Reader[*Dependencies, string] {
	return F.Pipe1(
		R.Ask[*Dependencies](),
		R.Chain(F.Flow4(
			getLowerBound,
			ord.Gt(I.Ord),
			ID.Flap[bool](n),
			f,
		)),
	)
}

func f(b bool) R.Reader[*Dependencies, string] {
	return func(deps *Dependencies) string {
		if b {
			return deps.I18n.True
		}
		return deps.I18n.False
	}
}

func ExampleReader() {

	deps1 := Dependencies{I18n: I18n{True: "vero", False: "falso"}, LowerBound: 2}

	hFoo := h("foo")
	hA := h("a")

	fmt.Println(hFoo(&deps1))
	fmt.Println(hA(&deps1))

	deps2 := Dependencies{I18n: I18n{True: "vero", False: "falso"}, LowerBound: 4}

	fmt.Println(hFoo(&deps2))
	fmt.Println(hA(&deps2))

	// Output:
	// vero
	// falso
	// falso
	// falso
}
