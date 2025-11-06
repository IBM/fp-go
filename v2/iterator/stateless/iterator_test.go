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

package stateless

import (
	"fmt"
	"math"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {

	result := F.Pipe2(
		A.From(1, 2, 3),
		FromArray[int],
		Reduce(utils.Sum, 0),
	)

	assert.Equal(t, 6, result)
}

func TestChain(t *testing.T) {

	outer := From(1, 2, 3)

	inner := func(data int) Iterator[string] {
		return F.Pipe2(
			A.From(0, 1),
			FromArray[int],
			Map(func(idx int) string {
				return fmt.Sprintf("item[%d][%d]", data, idx)
			}),
		)
	}

	total := F.Pipe2(
		outer,
		Chain(inner),
		ToArray[string],
	)

	assert.Equal(t, A.From("item[1][0]", "item[1][1]", "item[2][0]", "item[2][1]", "item[3][0]", "item[3][1]"), total)
}

func isPrimeNumber(num int) bool {
	if num <= 2 {
		return true
	}
	sqRoot := int(math.Sqrt(float64(num)))
	for i := 2; i <= sqRoot; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func TestFilterMap(t *testing.T) {

	it := F.Pipe3(
		MakeBy(utils.Inc),
		Take[int](100),
		FilterMap(O.FromPredicate(isPrimeNumber)),
		ToArray[int],
	)

	assert.Equal(t, A.From(1, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97), it)
}

func TestAp(t *testing.T) {

	f := F.Curry3(func(s1 string, n int, s2 string) string {
		return fmt.Sprintf("%s-%d-%s", s1, n, s2)
	})

	it := F.Pipe4(
		Of(f),
		Ap[func(int) func(string) string](From("a", "b")),
		Ap[func(string) string](From(1, 2)),
		Ap[string](From("c", "d")),
		ToArray[string],
	)

	assert.Equal(t, A.From("a-1-c", "a-1-d", "a-2-c", "a-2-d", "b-1-c", "b-1-d", "b-2-c", "b-2-d"), it)
}

func ExampleFoldMap() {
	src := From("a", "b", "c")

	fold := FoldMap[string](S.Monoid)(strings.ToUpper)

	fmt.Println(fold(src))

	// Output: ABC

}
