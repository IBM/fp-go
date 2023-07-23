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

package writer

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	M "github.com/IBM/fp-go/monoid"
	T "github.com/IBM/fp-go/tuple"
)

func doubleAndLog(data int) Writer[[]string, int] {
	return func() T.Tuple2[int, []string] {
		result := data * 2
		return T.MakeTuple2(result, A.Of(fmt.Sprintf("Doubled %d -> %d", data, result)))
	}
}

func ExampleLoggingWriter() {

	m := A.Monoid[string]()
	s := M.ToSemigroup(m)

	res := F.Pipe3(
		10,
		Of[int](m),
		Chain[int, int](s)(doubleAndLog),
		Chain[int, int](s)(doubleAndLog),
	)

	fmt.Println(res())

	// Output: {40 [Doubled 10 -> 20 Doubled 20 -> 40]}

}
