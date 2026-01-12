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

package iter

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Helper function to collect sequence into a slice
func toSlice[T any](seq Seq[T]) []T {
	return slices.Collect(seq)
}

// Helper function to collect Seq2 into a map
func toMap[K comparable, V any](seq Seq2[K, V]) map[K]V {
	return maps.Collect(seq)
}

func TestOf(t *testing.T) {
	seq := Of(42)
	result := toSlice(seq)
	assert.Equal(t, A.Of(42), result)
}

func TestOf2(t *testing.T) {
	seq := Of2("key", 100)
	result := toMap(seq)
	assert.Equal(t, R.Of("key", 100), result)
}

func TestFrom(t *testing.T) {
	seq := From(1, 2, 3, 4, 5)
	result := toSlice(seq)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestEmpty(t *testing.T) {
	seq := Empty[int]()
	result := toSlice(seq)
	assert.Empty(t, result)
}

func TestMonadMap(t *testing.T) {
	seq := From(1, 2, 3)
	doubled := MonadMap(seq, N.Mul(2))
	result := toSlice(doubled)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMap(t *testing.T) {
	seq := From(1, 2, 3)
	double := Map(N.Mul(2))
	result := toSlice(double(seq))
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMonadMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	indexed := MonadMapWithIndex(seq, func(i int, s string) string {
		return fmt.Sprintf("%d:%s", i, s)
	})
	result := toSlice(indexed)
	assert.Equal(t, []string{"0:a", "1:b", "2:c"}, result)
}

func TestMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	indexer := MapWithIndex(func(i int, s string) string {
		return fmt.Sprintf("%d:%s", i, s)
	})
	result := toSlice(indexer(seq))
	assert.Equal(t, []string{"0:a", "1:b", "2:c"}, result)
}

func TestMonadMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	doubled := MonadMapWithKey(seq, func(k string, v int) int { return v * 2 })
	result := toMap(doubled)
	assert.Equal(t, map[string]int{"x": 20}, result)
}

func TestMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	doubler := MapWithKey(func(k string, v int) int { return v * 2 })
	result := toMap(doubler(seq))
	assert.Equal(t, map[string]int{"x": 20}, result)
}

func TestMonadFilter(t *testing.T) {
	seq := From(1, 2, 3, 4, 5)
	evens := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
	result := toSlice(evens)
	assert.Equal(t, []int{2, 4}, result)
}

func TestFilter(t *testing.T) {
	seq := From(1, 2, 3, 4, 5)
	isEven := Filter(func(x int) bool { return x%2 == 0 })
	result := toSlice(isEven(seq))
	assert.Equal(t, []int{2, 4}, result)
}

func TestMonadFilterWithIndex(t *testing.T) {
	seq := From("a", "b", "c", "d")
	oddIndices := MonadFilterWithIndex(seq, func(i int, _ string) bool { return i%2 == 1 })
	result := toSlice(oddIndices)
	assert.Equal(t, []string{"b", "d"}, result)
}

func TestFilterWithIndex(t *testing.T) {
	seq := From("a", "b", "c", "d")
	oddIndexFilter := FilterWithIndex(func(i int, _ string) bool { return i%2 == 1 })
	result := toSlice(oddIndexFilter(seq))
	assert.Equal(t, []string{"b", "d"}, result)
}

func TestMonadFilterWithKey(t *testing.T) {
	seq := Of2("x", 10)
	filtered := MonadFilterWithKey(seq, func(k string, v int) bool { return v > 5 })
	result := toMap(filtered)
	assert.Equal(t, map[string]int{"x": 10}, result)

	seq2 := Of2("y", 3)
	filtered2 := MonadFilterWithKey(seq2, func(k string, v int) bool { return v > 5 })
	result2 := toMap(filtered2)
	assert.Equal(t, map[string]int{}, result2)
}

func TestFilterWithKey(t *testing.T) {
	seq := Of2("x", 10)
	filter := FilterWithKey(func(k string, v int) bool { return v > 5 })
	result := toMap(filter(seq))
	assert.Equal(t, map[string]int{"x": 10}, result)
}

func TestMonadFilterMap(t *testing.T) {
	seq := From(1, 2, 3, 4)
	result := MonadFilterMap(seq, func(x int) Option[int] {
		if x%2 == 0 {
			return O.Some(x * 10)
		}
		return O.None[int]()
	})
	assert.Equal(t, []int{20, 40}, toSlice(result))
}

func TestFilterMap(t *testing.T) {
	seq := From(1, 2, 3, 4)
	filterMapper := FilterMap(func(x int) Option[int] {
		if x%2 == 0 {
			return O.Some(x * 10)
		}
		return O.None[int]()
	})
	result := toSlice(filterMapper(seq))
	assert.Equal(t, []int{20, 40}, result)
}

func TestMonadFilterMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	result := MonadFilterMapWithIndex(seq, func(i int, s string) Option[string] {
		if i%2 == 0 {
			return O.Some(strings.ToUpper(s))
		}
		return O.None[string]()
	})
	assert.Equal(t, []string{"A", "C"}, toSlice(result))
}

func TestFilterMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	filterMapper := FilterMapWithIndex(func(i int, s string) Option[string] {
		if i%2 == 0 {
			return O.Some(strings.ToUpper(s))
		}
		return O.None[string]()
	})
	result := toSlice(filterMapper(seq))
	assert.Equal(t, []string{"A", "C"}, result)
}

func TestMonadFilterMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	result := MonadFilterMapWithKey(seq, func(k string, v int) Option[int] {
		if v > 5 {
			return O.Some(v * 2)
		}
		return O.None[int]()
	})
	assert.Equal(t, map[string]int{"x": 20}, toMap(result))
}

func TestFilterMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	filterMapper := FilterMapWithKey(func(k string, v int) Option[int] {
		if v > 5 {
			return O.Some(v * 2)
		}
		return O.None[int]()
	})
	result := toMap(filterMapper(seq))
	assert.Equal(t, map[string]int{"x": 20}, result)
}

func TestMonadChain(t *testing.T) {
	seq := From(1, 2)
	result := MonadChain(seq, func(x int) Seq[int] {
		return From(x, x*10)
	})
	assert.Equal(t, []int{1, 10, 2, 20}, toSlice(result))
}

func TestChain(t *testing.T) {
	seq := From(1, 2)
	chainer := Chain(func(x int) Seq[int] {
		return From(x, x*10)
	})
	result := toSlice(chainer(seq))
	assert.Equal(t, []int{1, 10, 2, 20}, result)
}

func TestFlatten(t *testing.T) {
	seq := From(From(1, 2), From(3, 4))
	result := Flatten(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, toSlice(result))
}

func TestMonadAp(t *testing.T) {
	fns := From(
		N.Mul(2),
		N.Add(10),
	)
	vals := From(1, 2)
	result := MonadAp(fns, vals)
	assert.Equal(t, []int{2, 4, 11, 12}, toSlice(result))
}

func TestAp(t *testing.T) {
	fns := From(
		N.Mul(2),
		N.Add(10),
	)
	vals := From(1, 2)
	applier := Ap[int](vals)
	result := toSlice(applier(fns))
	assert.Equal(t, []int{2, 4, 11, 12}, result)
}

func TestApCurried(t *testing.T) {
	f := F.Curry3(func(s1 string, n int, s2 string) string {
		return fmt.Sprintf("%s-%d-%s", s1, n, s2)
	})

	result := F.Pipe4(
		Of(f),
		Ap[func(int) func(string) string](From("a", "b")),
		Ap[func(string) string](From(1, 2)),
		Ap[string](From("c", "d")),
		toSlice[string],
	)

	expected := []string{"a-1-c", "a-1-d", "a-2-c", "a-2-d", "b-1-c", "b-1-d", "b-2-c", "b-2-d"}
	assert.Equal(t, expected, result)
}

func TestMakeBy(t *testing.T) {
	seq := MakeBy(5, func(i int) int { return i * i })
	result := toSlice(seq)
	assert.Equal(t, []int{0, 1, 4, 9, 16}, result)
}

func TestMakeByZero(t *testing.T) {
	seq := MakeBy(0, F.Identity)
	result := toSlice(seq)
	assert.Empty(t, result)
}

func TestMakeByNegative(t *testing.T) {
	seq := MakeBy(-5, F.Identity)
	result := toSlice(seq)
	assert.Empty(t, result)
}

func TestReplicate(t *testing.T) {
	seq := Replicate(3, "hello")
	result := toSlice(seq)
	assert.Equal(t, []string{"hello", "hello", "hello"}, result)
}

func TestMonadReduce(t *testing.T) {
	seq := From(1, 2, 3, 4)
	sum := MonadReduce(seq, func(acc, x int) int { return acc + x }, 0)
	assert.Equal(t, 10, sum)
}

func TestReduce(t *testing.T) {
	seq := From(1, 2, 3, 4)
	sum := Reduce(func(acc, x int) int { return acc + x }, 0)
	result := sum(seq)
	assert.Equal(t, 10, result)
}

func TestMonadReduceWithIndex(t *testing.T) {
	seq := From(10, 20, 30)
	result := MonadReduceWithIndex(seq, func(i, acc, x int) int {
		return acc + (i * x)
	}, 0)
	// 0*10 + 1*20 + 2*30 = 0 + 20 + 60 = 80
	assert.Equal(t, 80, result)
}

func TestReduceWithIndex(t *testing.T) {
	seq := From(10, 20, 30)
	reducer := ReduceWithIndex(func(i, acc, x int) int {
		return acc + (i * x)
	}, 0)
	result := reducer(seq)
	assert.Equal(t, 80, result)
}

func TestMonadReduceWithKey(t *testing.T) {
	seq := Of2("x", 10)
	result := MonadReduceWithKey(seq, func(k string, acc, v int) int {
		return acc + v
	}, 0)
	assert.Equal(t, 10, result)
}

func TestReduceWithKey(t *testing.T) {
	seq := Of2("x", 10)
	reducer := ReduceWithKey(func(k string, acc, v int) int {
		return acc + v
	}, 0)
	result := reducer(seq)
	assert.Equal(t, 10, result)
}

func TestMonadFold(t *testing.T) {
	seq := From("Hello", " ", "World")
	result := MonadFold(seq, S.Monoid)
	assert.Equal(t, "Hello World", result)
}

func TestFold(t *testing.T) {
	seq := From("Hello", " ", "World")
	folder := Fold(S.Monoid)
	result := folder(seq)
	assert.Equal(t, "Hello World", result)
}

func TestMonadFoldMap(t *testing.T) {
	seq := From(1, 2, 3)
	result := MonadFoldMap(seq, strconv.Itoa, S.Monoid)
	assert.Equal(t, "123", result)
}

func TestFoldMap(t *testing.T) {
	seq := From(1, 2, 3)
	folder := FoldMap[int](S.Monoid)(strconv.Itoa)
	result := folder(seq)
	assert.Equal(t, "123", result)
}

func TestMonadFoldMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	result := MonadFoldMapWithIndex(seq, func(i int, s string) string {
		return fmt.Sprintf("%d:%s ", i, s)
	}, S.Monoid)
	assert.Equal(t, "0:a 1:b 2:c ", result)
}

func TestFoldMapWithIndex(t *testing.T) {
	seq := From("a", "b", "c")
	folder := FoldMapWithIndex[string](S.Monoid)(func(i int, s string) string {
		return fmt.Sprintf("%d:%s ", i, s)
	})
	result := folder(seq)
	assert.Equal(t, "0:a 1:b 2:c ", result)
}

func TestMonadFoldMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	result := MonadFoldMapWithKey(seq, func(k string, v int) string {
		return fmt.Sprintf("%s:%d ", k, v)
	}, S.Monoid)
	assert.Equal(t, "x:10 ", result)
}

func TestFoldMapWithKey(t *testing.T) {
	seq := Of2("x", 10)
	folder := FoldMapWithKey[string, int](S.Monoid)(func(k string, v int) string {
		return fmt.Sprintf("%s:%d ", k, v)
	})
	result := folder(seq)
	assert.Equal(t, "x:10 ", result)
}

func TestMonadFlap(t *testing.T) {
	fns := From(
		N.Mul(2),
		N.Add(10),
	)
	result := MonadFlap(fns, 5)
	assert.Equal(t, []int{10, 15}, toSlice(result))
}

func TestFlap(t *testing.T) {
	fns := From(
		N.Mul(2),
		N.Add(10),
	)
	flapper := Flap[int](5)
	result := toSlice(flapper(fns))
	assert.Equal(t, []int{10, 15}, result)
}

func TestPrepend(t *testing.T) {
	seq := From(2, 3, 4)
	result := Prepend(1)(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, toSlice(result))
}

func TestAppend(t *testing.T) {
	seq := From(1, 2, 3)
	result := Append(4)(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, toSlice(result))
}

func TestMonadZip(t *testing.T) {
	seqA := From(1, 2, 3)
	seqB := From("a", "b")
	result := MonadZip(seqB, seqA)

	var pairs []string
	for a, b := range result {
		pairs = append(pairs, fmt.Sprintf("%d:%s", b, a))
	}
	assert.Equal(t, []string{"1:a", "2:b"}, pairs)
}

func TestZip(t *testing.T) {
	seqA := From(1, 2, 3)
	seqB := From("a", "b", "c")
	zipWithA := Zip[int](seqB)
	result := zipWithA(seqA)

	var pairs []string
	for a, b := range result {
		pairs = append(pairs, fmt.Sprintf("%d:%s", a, b))
	}
	assert.Equal(t, []string{"1:a", "2:b", "3:c"}, pairs)
}

func TestMonoid(t *testing.T) {
	m := Monoid[int]()
	seq1 := From(1, 2)
	seq2 := From(3, 4)
	result := m.Concat(seq1, seq2)
	assert.Equal(t, []int{1, 2, 3, 4}, toSlice(result))
}

func TestMonoidEmpty(t *testing.T) {
	m := Monoid[int]()
	empty := m.Empty()
	assert.Empty(t, toSlice(empty))
}

func TestMonoidAssociativity(t *testing.T) {
	m := Monoid[int]()
	seq1 := From(1, 2)
	seq2 := From(3, 4)
	seq3 := From(5, 6)

	// (seq1 + seq2) + seq3
	left := m.Concat(m.Concat(seq1, seq2), seq3)
	// seq1 + (seq2 + seq3)
	right := m.Concat(seq1, m.Concat(seq2, seq3))

	assert.Equal(t, toSlice(left), toSlice(right))
}

func TestMonoidIdentity(t *testing.T) {
	m := Monoid[int]()
	seq := From(1, 2, 3)
	empty := m.Empty()

	// seq + empty = seq
	leftIdentity := m.Concat(seq, empty)
	assert.Equal(t, []int{1, 2, 3}, toSlice(leftIdentity))

	// empty + seq = seq
	rightIdentity := m.Concat(empty, seq)
	assert.Equal(t, []int{1, 2, 3}, toSlice(rightIdentity))
}

func TestPipelineComposition(t *testing.T) {
	// Test a complex pipeline
	result := F.Pipe4(
		From(1, 2, 3, 4, 5, 6),
		Filter(func(x int) bool { return x%2 == 0 }),
		Map(N.Mul(10)),
		Prepend(0),
		toSlice[int],
	)
	assert.Equal(t, []int{0, 20, 40, 60}, result)
}

func TestLazyEvaluation(t *testing.T) {
	// Test that operations are lazy
	callCount := 0
	seq := From(1, 2, 3, 4, 5)
	mapped := MonadMap(seq, func(x int) int {
		callCount++
		return x * 2
	})

	// No calls yet since we haven't iterated
	assert.Equal(t, 0, callCount)

	// Iterate only first 2 elements
	count := 0
	for range mapped {
		count++
		if count == 2 {
			break
		}
	}

	// Should have called the function only twice
	assert.Equal(t, 2, callCount)
}

func ExampleFoldMap() {
	seq := From("a", "b", "c")
	fold := FoldMap[string](S.Monoid)(strings.ToUpper)
	result := fold(seq)
	fmt.Println(result)
	// Output: ABC
}

func ExampleChain() {
	seq := From(1, 2)
	result := F.Pipe2(
		seq,
		Chain(func(x int) Seq[int] {
			return From(x, x*10)
		}),
		toSlice[int],
	)
	fmt.Println(result)
	// Output: [1 10 2 20]
}

func ExampleMonoid() {
	m := Monoid[int]()
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)
	combined := m.Concat(seq1, seq2)
	result := toSlice(combined)
	fmt.Println(result)
	// Output: [1 2 3 4 5 6]
}

func TestMonadMapToArray(t *testing.T) {
	seq := From(1, 2, 3)
	result := MonadMapToArray(seq, N.Mul(2))
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMonadMapToArrayEmpty(t *testing.T) {
	seq := Empty[int]()
	result := MonadMapToArray(seq, N.Mul(2))
	assert.Empty(t, result)
}

func TestMapToArray(t *testing.T) {
	seq := From(1, 2, 3)
	mapper := MapToArray(N.Mul(2))
	result := mapper(seq)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestMapToArrayIdentity(t *testing.T) {
	seq := From("a", "b", "c")
	mapper := MapToArray(F.Identity[string])
	result := mapper(seq)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}
