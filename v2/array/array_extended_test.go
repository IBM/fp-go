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

package array

import (
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestReplicate(t *testing.T) {
	result := Replicate(3, "a")
	assert.Equal(t, []string{"a", "a", "a"}, result)

	empty := Replicate(0, 42)
	assert.Equal(t, []int{}, empty)
}

func TestMonadMap(t *testing.T) {
	src := []int{1, 2, 3}
	result := MonadMap(src, func(x int) int { return x * 2 })
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMonadMapRef(t *testing.T) {
	src := []int{1, 2, 3}
	result := MonadMapRef(src, func(x *int) int { return *x * 2 })
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMapWithIndex(t *testing.T) {
	src := []string{"a", "b", "c"}
	mapper := MapWithIndex(func(i int, s string) string {
		return fmt.Sprintf("%d:%s", i, s)
	})
	result := mapper(src)
	assert.Equal(t, []string{"0:a", "1:b", "2:c"}, result)
}

func TestMapRef(t *testing.T) {
	src := []int{1, 2, 3}
	mapper := MapRef(func(x *int) int { return *x * 2 })
	result := mapper(src)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestFilterWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	filter := FilterWithIndex(func(i, x int) bool {
		return i%2 == 0 && x > 2
	})
	result := filter(src)
	assert.Equal(t, []int{3, 5}, result)
}

func TestFilterRef(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	filter := FilterRef(func(x *int) bool { return *x > 2 })
	result := filter(src)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestMonadFilterMap(t *testing.T) {
	src := []int{1, 2, 3, 4}
	result := MonadFilterMap(src, func(x int) O.Option[string] {
		if x%2 == 0 {
			return O.Some(fmt.Sprintf("even:%d", x))
		}
		return O.None[string]()
	})
	assert.Equal(t, []string{"even:2", "even:4"}, result)
}

func TestMonadFilterMapWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4}
	result := MonadFilterMapWithIndex(src, func(i, x int) O.Option[string] {
		if i%2 == 0 {
			return O.Some(fmt.Sprintf("%d:%d", i, x))
		}
		return O.None[string]()
	})
	assert.Equal(t, []string{"0:1", "2:3"}, result)
}

func TestFilterMapWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4}
	filter := FilterMapWithIndex(func(i, x int) O.Option[string] {
		if i%2 == 0 {
			return O.Some(fmt.Sprintf("%d:%d", i, x))
		}
		return O.None[string]()
	})
	result := filter(src)
	assert.Equal(t, []string{"0:1", "2:3"}, result)
}

func TestFilterMapRef(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	filter := FilterMapRef(
		func(x *int) bool { return *x > 2 },
		func(x *int) string { return fmt.Sprintf("val:%d", *x) },
	)
	result := filter(src)
	assert.Equal(t, []string{"val:3", "val:4", "val:5"}, result)
}

func TestReduceWithIndex(t *testing.T) {
	src := []int{1, 2, 3}
	reducer := ReduceWithIndex(func(i, acc, x int) int {
		return acc + i + x
	}, 0)
	result := reducer(src)
	assert.Equal(t, 9, result) // 0 + (0+1) + (1+2) + (2+3) = 9
}

func TestReduceRightWithIndex(t *testing.T) {
	src := []string{"a", "b", "c"}
	reducer := ReduceRightWithIndex(func(i int, x, acc string) string {
		return fmt.Sprintf("%s%d:%s", acc, i, x)
	}, "")
	result := reducer(src)
	assert.Equal(t, "2:c1:b0:a", result)
}

func TestReduceRef(t *testing.T) {
	src := []int{1, 2, 3}
	reducer := ReduceRef(func(acc int, x *int) int {
		return acc + *x
	}, 0)
	result := reducer(src)
	assert.Equal(t, 6, result)
}

func TestZero(t *testing.T) {
	result := Zero[int]()
	assert.Equal(t, []int{}, result)
	assert.True(t, IsEmpty(result))
}

func TestMonadChain(t *testing.T) {
	src := []int{1, 2, 3}
	result := MonadChain(src, func(x int) []int {
		return []int{x, x * 10}
	})
	assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, result)
}

func TestChain(t *testing.T) {
	src := []int{1, 2, 3}
	chain := Chain(func(x int) []int {
		return []int{x, x * 10}
	})
	result := chain(src)
	assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, result)
}

func TestMonadAp(t *testing.T) {
	fns := []func(int) int{
		func(x int) int { return x * 2 },
		func(x int) int { return x + 10 },
	}
	values := []int{1, 2}
	result := MonadAp(fns, values)
	assert.Equal(t, []int{2, 4, 11, 12}, result)
}

func TestMatchLeft(t *testing.T) {
	matcher := MatchLeft(
		func() string { return "empty" },
		func(head int, tail []int) string {
			return fmt.Sprintf("head:%d,tail:%v", head, tail)
		},
	)

	assert.Equal(t, "empty", matcher([]int{}))
	assert.Equal(t, "head:1,tail:[2 3]", matcher([]int{1, 2, 3}))
}

func TestTail(t *testing.T) {
	assert.Equal(t, O.None[[]int](), Tail([]int{}))
	assert.Equal(t, O.Some([]int{2, 3}), Tail([]int{1, 2, 3}))
	assert.Equal(t, O.Some([]int{}), Tail([]int{1}))
}

func TestFirst(t *testing.T) {
	assert.Equal(t, O.None[int](), First([]int{}))
	assert.Equal(t, O.Some(1), First([]int{1, 2, 3}))
}

func TestLast(t *testing.T) {
	assert.Equal(t, O.None[int](), Last([]int{}))
	assert.Equal(t, O.Some(3), Last([]int{1, 2, 3}))
	assert.Equal(t, O.Some(1), Last([]int{1}))
}

func TestUpsertAt(t *testing.T) {
	src := []int{1, 2, 3}
	upsert := UpsertAt(99)

	result1 := upsert(src)
	assert.Equal(t, []int{1, 2, 3, 99}, result1)
}

func TestSize(t *testing.T) {
	assert.Equal(t, 0, Size([]int{}))
	assert.Equal(t, 3, Size([]int{1, 2, 3}))
}

func TestMonadPartition(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	result := MonadPartition(src, func(x int) bool { return x > 2 })
	assert.Equal(t, []int{1, 2}, result.F1)
	assert.Equal(t, []int{3, 4, 5}, result.F2)
}

func TestIsNil(t *testing.T) {
	var nilSlice []int
	assert.True(t, IsNil(nilSlice))
	assert.False(t, IsNil([]int{}))
	assert.False(t, IsNil([]int{1}))
}

func TestIsNonNil(t *testing.T) {
	var nilSlice []int
	assert.False(t, IsNonNil(nilSlice))
	assert.True(t, IsNonNil([]int{}))
	assert.True(t, IsNonNil([]int{1}))
}

func TestConstNil(t *testing.T) {
	result := ConstNil[int]()
	assert.True(t, IsNil(result))
}

func TestSliceRight(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	slicer := SliceRight[int](2)
	result := slicer(src)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestCopy(t *testing.T) {
	src := []int{1, 2, 3}
	copied := Copy(src)
	assert.Equal(t, src, copied)
	// Verify it's a different slice
	copied[0] = 99
	assert.Equal(t, 1, src[0])
	assert.Equal(t, 99, copied[0])
}

func TestClone(t *testing.T) {
	src := []int{1, 2, 3}
	cloner := Clone(func(x int) int { return x * 2 })
	result := cloner(src)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestFoldMapWithIndex(t *testing.T) {
	src := []string{"a", "b", "c"}
	folder := FoldMapWithIndex[string](S.Monoid)(func(i int, s string) string {
		return fmt.Sprintf("%d:%s", i, s)
	})
	result := folder(src)
	assert.Equal(t, "0:a1:b2:c", result)
}

func TestFold(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	folder := Fold(N.MonoidSum[int]())
	result := folder(src)
	assert.Equal(t, 15, result)
}

func TestPush(t *testing.T) {
	src := []int{1, 2, 3}
	pusher := Push(4)
	result := pusher(src)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestMonadFlap(t *testing.T) {
	fns := []func(int) string{
		func(x int) string { return fmt.Sprintf("a%d", x) },
		func(x int) string { return fmt.Sprintf("b%d", x) },
	}
	result := MonadFlap(fns, 5)
	assert.Equal(t, []string{"a5", "b5"}, result)
}

func TestFlap(t *testing.T) {
	fns := []func(int) string{
		func(x int) string { return fmt.Sprintf("a%d", x) },
		func(x int) string { return fmt.Sprintf("b%d", x) },
	}
	flapper := Flap[string](5)
	result := flapper(fns)
	assert.Equal(t, []string{"a5", "b5"}, result)
}

func TestPrepend(t *testing.T) {
	src := []int{2, 3, 4}
	prepender := Prepend(1)
	result := prepender(src)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}
