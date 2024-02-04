// Copyright (c) 2023 IBM Corp.
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

// Package reader implements the example in [https://dev.to/gcanti/getting-started-with-fp-ts-reader-1ie5]
package example1

import (
	"fmt"

	F "github.com/IBM/fp-go/function"
	N "github.com/IBM/fp-go/number"
	I "github.com/IBM/fp-go/number/integer"
	"github.com/IBM/fp-go/ord"
	R "github.com/IBM/fp-go/reader"
	S "github.com/IBM/fp-go/string"
)

type (
	I18n struct {
		True  string
		False string
	}

	Dependencies struct {
		I18n I18n
	}
)

var (
	// g: func(int) R.Reader[*Dependencies, string], note how the implementation does not depend on the dependencies
	g = F.Flow2(
		ord.Gt(I.Ord)(2),
		f,
	)

	// h: func(string) R.Reader[*Dependencies, string], note how the implementation does not depend on the dependencies
	h = F.Flow3(
		S.Size,
		N.Add(1),
		g,
	)
)

func f(b bool) R.Reader[*Dependencies, string] {
	return func(deps *Dependencies) string {
		if b {
			return deps.I18n.True
		}
		return deps.I18n.False
	}
}

func ExampleReader() {

	deps := Dependencies{I18n: I18n{True: "vero", False: "falso"}}

	fmt.Println(h("foo")(&deps))
	fmt.Println(h("a")(&deps))

	// Output:
	// vero
	// falso
}
