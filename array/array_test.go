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

package array

import (
	"fmt"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestMap1(t *testing.T) {

	src := []string{"a", "b", "c"}

	up := Map(strings.ToUpper)(src)

	var up1 = []string{}
	for _, s := range src {
		up1 = append(up1, strings.ToUpper(s))
	}

	var up2 = []string{}
	for i := range src {
		up2 = append(up2, strings.ToUpper(src[i]))
	}

	assert.Equal(t, up, up1)
	assert.Equal(t, up, up2)
}

func TestMap(t *testing.T) {

	mapper := Map(utils.Upper)

	src := []string{"a", "b", "c"}

	dst := mapper(src)

	assert.Equal(t, dst, []string{"A", "B", "C"})
}

func TestReduce(t *testing.T) {

	values := MakeBy(101, F.Identity[int])

	sum := func(val int, current int) int {
		return val + current
	}
	reducer := Reduce(sum, 0)

	result := reducer(values)
	assert.Equal(t, result, 5050)

}

func TestEmpty(t *testing.T) {
	assert.True(t, IsNonEmpty(MakeBy(101, F.Identity[int])))
	assert.True(t, IsEmpty([]int{}))
}

func TestAp(t *testing.T) {
	assert.Equal(t,
		[]int{2, 4, 6, 3, 6, 9},
		F.Pipe1(
			[]func(int) int{
				utils.Double,
				utils.Triple,
			},
			Ap[int, int]([]int{1, 2, 3}),
		),
	)
}

func TestIntercalate(t *testing.T) {
	is := Intercalate(S.Monoid)("-")

	assert.Equal(t, "", is(Empty[string]()))
	assert.Equal(t, "a", is([]string{"a"}))
	assert.Equal(t, "a-b-c", is([]string{"a", "b", "c"}))
	assert.Equal(t, "a--c", is([]string{"a", "", "c"}))
	assert.Equal(t, "a-b", is([]string{"a", "b"}))
	assert.Equal(t, "a-b-c-d", is([]string{"a", "b", "c", "d"}))
}

func TestPrependAll(t *testing.T) {
	empty := Empty[int]()
	prep := PrependAll(0)
	assert.Equal(t, empty, prep(empty))
	assert.Equal(t, []int{0, 1, 0, 2, 0, 3}, prep([]int{1, 2, 3}))
	assert.Equal(t, []int{0, 1}, prep([]int{1}))
	assert.Equal(t, []int{0, 1, 0, 2, 0, 3, 0, 4}, prep([]int{1, 2, 3, 4}))
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Flatten([][]int{{1}, {2}, {3}}))
}

func TestLookup(t *testing.T) {
	data := []int{0, 1, 2}
	none := O.None[int]()

	assert.Equal(t, none, Lookup[int](-1)(data))
	assert.Equal(t, none, Lookup[int](10)(data))
	assert.Equal(t, O.Some(1), Lookup[int](1)(data))
}

func TestSlice(t *testing.T) {
	data := []int{0, 1, 2, 3}

	assert.Equal(t, []int{1, 2}, Slice[int](1, 3)(data))
}

func TestFrom(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, From(1, 2, 3))
}

func TestPartition(t *testing.T) {

	pred := func(n int) bool {
		return n > 2
	}

	assert.Equal(t, T.MakeTuple2(Empty[int](), Empty[int]()), Partition(pred)(Empty[int]()))
	assert.Equal(t, T.MakeTuple2(From(1), From(3)), Partition(pred)(From(1, 3)))
}

func TestFilterChain(t *testing.T) {
	src := From(1, 2, 3)

	f := func(i int) O.Option[[]string] {
		if i%2 != 0 {
			return O.Of(From(fmt.Sprintf("a%d", i), fmt.Sprintf("b%d", i)))
		}
		return O.None[[]string]()
	}

	res := FilterChain(f)(src)

	assert.Equal(t, From("a1", "b1", "a3", "b3"), res)
}

func TestFilterMap(t *testing.T) {
	src := From(1, 2, 3)

	f := func(i int) O.Option[string] {
		if i%2 != 0 {
			return O.Of(fmt.Sprintf("a%d", i))
		}
		return O.None[string]()
	}

	res := FilterMap(f)(src)

	assert.Equal(t, From("a1", "a3"), res)
}
