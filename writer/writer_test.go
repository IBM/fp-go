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
	P "github.com/IBM/fp-go/pair"
)

func doubleAndLog(data int) Writer[[]string, int] {
	return func() P.Pair[int, []string] {
		result := data * 2
		return P.MakePair(result, A.Of(fmt.Sprintf("Doubled %d -> %d", data, result)))
	}
}

func ExampleWriter_logging() {

	res := F.Pipe3(
		Of[int](monoid, 10),
		Chain(sg, doubleAndLog),
		Chain(sg, doubleAndLog),
		Execute[[]string, int],
	)

	fmt.Println(res)

	// Output: [Doubled 10 -> 20 Doubled 20 -> 40]
}
