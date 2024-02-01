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
	S "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
)

func doubleAndLog(data int) Writer[[]string, int] {
	return func() T.Tuple3[int, []string, S.Semigroup[[]string]] {
		result := data * 2
		return T.MakeTuple3(result, A.Of(fmt.Sprintf("Doubled %d -> %d", data, result)), sg)
	}
}

func ExampleWriter_logging() {

	res := F.Pipe4(
		10,
		Of[int](monoid),
		Chain(doubleAndLog),
		Chain(doubleAndLog),
		Execute[[]string, int],
	)

	fmt.Println(res)

	// Output: [Doubled 10 -> 20 Doubled 20 -> 40]
}
