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
	"strconv"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestNilSlice_IsEmpty verifies that IsEmpty handles nil slices correctly
func TestNilSlice_IsEmpty(t *testing.T) {
	var nilSlice []int
	assert.True(t, IsEmpty(nilSlice), "nil slice should be empty")
}

// TestNilSlice_IsNonEmpty verifies that IsNonEmpty handles nil slices correctly
func TestNilSlice_IsNonEmpty(t *testing.T) {
	var nilSlice []int
	assert.False(t, IsNonEmpty(nilSlice), "nil slice should not be non-empty")
}

// TestNilSlice_MonadMap verifies that MonadMap handles nil slices correctly
func TestNilSlice_MonadMap(t *testing.T) {
	var nilSlice []int
	result := MonadMap(nilSlice, strconv.Itoa)
	assert.Empty(t, result, "MonadMap should return empty slice for nil input")
}

// TestNilSlice_MonadMapRef verifies that MonadMapRef handles nil slices correctly
func TestNilSlice_MonadMapRef(t *testing.T) {
	var nilSlice []int
	result := MonadMapRef(nilSlice, func(v *int) string {
		return fmt.Sprintf("%d", *v)
	})
	assert.NotNil(t, result, "MonadMapRef should return non-nil slice")
	assert.Equal(t, 0, len(result), "MonadMapRef should return empty slice for nil input")
}

// TestNilSlice_Map verifies that Map handles nil slices correctly
func TestNilSlice_Map(t *testing.T) {
	var nilSlice []int
	mapper := Map(strconv.Itoa)
	result := mapper(nilSlice)
	assert.Empty(t, result, "Map should return empty slice for nil input")
}

// TestNilSlice_MapRef verifies that MapRef handles nil slices correctly
func TestNilSlice_MapRef(t *testing.T) {
	var nilSlice []int
	mapper := MapRef(func(v *int) string {
		return fmt.Sprintf("%d", *v)
	})
	result := mapper(nilSlice)
	assert.NotNil(t, result, "MapRef should return non-nil slice")
	assert.Equal(t, 0, len(result), "MapRef should return empty slice for nil input")
}

// TestNilSlice_MapWithIndex verifies that MapWithIndex handles nil slices correctly
func TestNilSlice_MapWithIndex(t *testing.T) {
	var nilSlice []int
	mapper := MapWithIndex(func(i int, v int) string {
		return fmt.Sprintf("%d:%d", i, v)
	})
	result := mapper(nilSlice)
	assert.Empty(t, result, "MapWithIndex should return empty slice for nil input")
}

// TestNilSlice_Filter verifies that Filter handles nil slices correctly
func TestNilSlice_Filter(t *testing.T) {
	var nilSlice []int
	filter := Filter(func(v int) bool {
		return v > 0
	})
	result := filter(nilSlice)
	assert.True(t, IsEmpty(result), "Filter should return empty slice for nil input")
}

// TestNilSlice_FilterWithIndex verifies that FilterWithIndex handles nil slices correctly
func TestNilSlice_FilterWithIndex(t *testing.T) {
	var nilSlice []int
	filter := FilterWithIndex(func(i int, v int) bool {
		return v > 0
	})
	result := filter(nilSlice)
	assert.True(t, IsEmpty(result), "FilterWithIndex should return empty slice for nil input")
}

// TestNilSlice_FilterRef verifies that FilterRef handles nil slices correctly
func TestNilSlice_FilterRef(t *testing.T) {
	var nilSlice []int
	filter := FilterRef(func(v *int) bool {
		return *v > 0
	})
	result := filter(nilSlice)
	assert.True(t, IsEmpty(result), "FilterRef should return empty slice for nil input")
}

// TestNilSlice_MonadFilterMap verifies that MonadFilterMap handles nil slices correctly
func TestNilSlice_MonadFilterMap(t *testing.T) {
	var nilSlice []int
	result := MonadFilterMap(nilSlice, func(v int) O.Option[string] {
		return O.Some(fmt.Sprintf("%d", v))
	})
	assert.True(t, IsEmpty(result), "MonadFilterMap should return empty slice for nil input")
}

// TestNilSlice_MonadFilterMapWithIndex verifies that MonadFilterMapWithIndex handles nil slices correctly
func TestNilSlice_MonadFilterMapWithIndex(t *testing.T) {
	var nilSlice []int
	result := MonadFilterMapWithIndex(nilSlice, func(i int, v int) O.Option[string] {
		return O.Some(fmt.Sprintf("%d:%d", i, v))
	})
	assert.True(t, IsEmpty(result), "MonadFilterMapWithIndex should return empty slice for nil input")
}

// TestNilSlice_FilterMap verifies that FilterMap handles nil slices correctly
func TestNilSlice_FilterMap(t *testing.T) {
	var nilSlice []int
	filter := FilterMap(func(v int) O.Option[string] {
		return O.Some(fmt.Sprintf("%d", v))
	})
	result := filter(nilSlice)
	assert.True(t, IsEmpty(result), "FilterMap should return empty slice for nil input")
}

// TestNilSlice_FilterMapWithIndex verifies that FilterMapWithIndex handles nil slices correctly
func TestNilSlice_FilterMapWithIndex(t *testing.T) {
	var nilSlice []int
	filter := FilterMapWithIndex(func(i int, v int) O.Option[string] {
		return O.Some(fmt.Sprintf("%d:%d", i, v))
	})
	result := filter(nilSlice)
	assert.True(t, IsEmpty(result), "FilterMapWithIndex should return empty slice for nil input")
}

// TestNilSlice_MonadReduce verifies that MonadReduce handles nil slices correctly
func TestNilSlice_MonadReduce(t *testing.T) {
	var nilSlice []int
	result := MonadReduce(nilSlice, func(acc int, v int) int {
		return acc + v
	}, 10)
	assert.Equal(t, 10, result, "MonadReduce should return initial value for nil slice")
}

// TestNilSlice_MonadReduceWithIndex verifies that MonadReduceWithIndex handles nil slices correctly
func TestNilSlice_MonadReduceWithIndex(t *testing.T) {
	var nilSlice []int
	result := MonadReduceWithIndex(nilSlice, func(i int, acc int, v int) int {
		return acc + v
	}, 10)
	assert.Equal(t, 10, result, "MonadReduceWithIndex should return initial value for nil slice")
}

// TestNilSlice_Reduce verifies that Reduce handles nil slices correctly
func TestNilSlice_Reduce(t *testing.T) {
	var nilSlice []int
	reducer := Reduce(func(acc int, v int) int {
		return acc + v
	}, 10)
	result := reducer(nilSlice)
	assert.Equal(t, 10, result, "Reduce should return initial value for nil slice")
}

// TestNilSlice_ReduceWithIndex verifies that ReduceWithIndex handles nil slices correctly
func TestNilSlice_ReduceWithIndex(t *testing.T) {
	var nilSlice []int
	reducer := ReduceWithIndex(func(i int, acc int, v int) int {
		return acc + v
	}, 10)
	result := reducer(nilSlice)
	assert.Equal(t, 10, result, "ReduceWithIndex should return initial value for nil slice")
}

// TestNilSlice_ReduceRight verifies that ReduceRight handles nil slices correctly
func TestNilSlice_ReduceRight(t *testing.T) {
	var nilSlice []int
	reducer := ReduceRight(func(v int, acc int) int {
		return acc + v
	}, 10)
	result := reducer(nilSlice)
	assert.Equal(t, 10, result, "ReduceRight should return initial value for nil slice")
}

// TestNilSlice_ReduceRightWithIndex verifies that ReduceRightWithIndex handles nil slices correctly
func TestNilSlice_ReduceRightWithIndex(t *testing.T) {
	var nilSlice []int
	reducer := ReduceRightWithIndex(func(i int, v int, acc int) int {
		return acc + v
	}, 10)
	result := reducer(nilSlice)
	assert.Equal(t, 10, result, "ReduceRightWithIndex should return initial value for nil slice")
}

// TestNilSlice_ReduceRef verifies that ReduceRef handles nil slices correctly
func TestNilSlice_ReduceRef(t *testing.T) {
	var nilSlice []int
	reducer := ReduceRef(func(acc int, v *int) int {
		return acc + *v
	}, 10)
	result := reducer(nilSlice)
	assert.Equal(t, 10, result, "ReduceRef should return initial value for nil slice")
}

// TestNilSlice_Append verifies that Append handles nil slices correctly
func TestNilSlice_Append(t *testing.T) {
	var nilSlice []int
	result := Append(nilSlice, 42)
	assert.NotNil(t, result, "Append should return non-nil slice")
	assert.Equal(t, 1, len(result), "Append should create slice with one element")
	assert.Equal(t, 42, result[0], "Append should add element correctly")
}

// TestNilSlice_MonadChain verifies that MonadChain handles nil slices correctly
func TestNilSlice_MonadChain(t *testing.T) {
	var nilSlice []int
	result := MonadChain(nilSlice, func(v int) []string {
		return []string{fmt.Sprintf("%d", v)}
	})
	// nil is a valid representation of the empty array.
	assert.Empty(t, result, "MonadChain should return an empty slice for nil input")
}

// TestNilSlice_Chain verifies that Chain handles nil slices correctly
func TestNilSlice_Chain(t *testing.T) {
	var nilSlice []int
	chain := Chain(func(v int) []string {
		return []string{fmt.Sprintf("%d", v)}
	})
	result := chain(nilSlice)
	// nil is a valid representation of the empty array.
	assert.Empty(t, result, "Chain should return an empty slice for nil input")
}

// TestNilSlice_MonadAp verifies that MonadAp handles nil slices correctly
func TestNilSlice_MonadAp(t *testing.T) {
	var nilFuncs []func(int) string
	var nilValues []int

	// nil is a valid representation of the empty array, so these cases only
	// assert that an empty result is produced.

	// nil functions, nil values
	result1 := MonadAp(nilFuncs, nilValues)
	assert.Empty(t, result1, "MonadAp should return an empty slice for nil inputs")

	// nil functions, non-nil values
	nonNilValues := []int{1, 2, 3}
	result2 := MonadAp(nilFuncs, nonNilValues)
	assert.Empty(t, result2, "MonadAp should return an empty slice when functions are nil")

	// non-nil functions, nil values
	nonNilFuncs := []func(int) string{func(v int) string { return fmt.Sprintf("%d", v) }}
	result3 := MonadAp(nonNilFuncs, nilValues)
	assert.Empty(t, result3, "MonadAp should return an empty slice when values are nil")
}

// TestNilSlice_Ap verifies that Ap handles nil slices correctly
func TestNilSlice_Ap(t *testing.T) {
	var nilValues []int
	ap := Ap[string](nilValues)

	var nilFuncs []func(int) string
	result := ap(nilFuncs)
	// nil is a valid representation of the empty array.
	assert.Empty(t, result, "Ap should return an empty slice for nil inputs")
}

// TestNilSlice_Head verifies that Head handles nil slices correctly
func TestNilSlice_Head(t *testing.T) {
	var nilSlice []int
	result := Head(nilSlice)
	assert.True(t, O.IsNone(result), "Head should return None for nil slice")
}

// TestNilSlice_First verifies that First handles nil slices correctly
func TestNilSlice_First(t *testing.T) {
	var nilSlice []int
	result := First(nilSlice)
	assert.True(t, O.IsNone(result), "First should return None for nil slice")
}

// TestNilSlice_Last verifies that Last handles nil slices correctly
func TestNilSlice_Last(t *testing.T) {
	var nilSlice []int
	result := Last(nilSlice)
	assert.True(t, O.IsNone(result), "Last should return None for nil slice")
}

// TestNilSlice_Tail verifies that Tail handles nil slices correctly
func TestNilSlice_Tail(t *testing.T) {
	var nilSlice []int
	result := Tail(nilSlice)
	assert.True(t, O.IsNone(result), "Tail should return None for nil slice")
}

// TestNilSlice_Flatten verifies that Flatten handles nil slices correctly
func TestNilSlice_Flatten(t *testing.T) {
	var nilSlice [][]int
	result := Flatten(nilSlice)
	// nil is a valid representation of the empty array, so only the emptiness
	// of the result is asserted.
	assert.Empty(t, result, "Flatten should return an empty slice for nil input")
}

// TestNilSlice_Lookup verifies that Lookup handles nil slices correctly
func TestNilSlice_Lookup(t *testing.T) {
	var nilSlice []int
	lookup := Lookup[int](0)
	result := lookup(nilSlice)
	assert.True(t, O.IsNone(result), "Lookup should return None for nil slice")
}

// TestNilSlice_Size verifies that Size handles nil slices correctly
func TestNilSlice_Size(t *testing.T) {
	var nilSlice []int
	result := Size(nilSlice)
	assert.Equal(t, 0, result, "Size should return 0 for nil slice")
}

// TestNilSlice_MonadPartition verifies that MonadPartition handles nil slices correctly
func TestNilSlice_MonadPartition(t *testing.T) {
	var nilSlice []int
	result := MonadPartition(nilSlice, func(v int) bool {
		return v > 0
	})
	// nil is a valid representation of the empty array.
	assert.Empty(t, P.Head(result), "MonadPartition left should be empty for nil input")
	assert.Empty(t, P.Tail(result), "MonadPartition right should be empty for nil input")
}

// TestNilSlice_Partition verifies that Partition handles nil slices correctly
func TestNilSlice_Partition(t *testing.T) {
	var nilSlice []int
	partition := Partition(func(v int) bool {
		return v > 0
	})
	result := partition(nilSlice)
	// nil is a valid representation of the empty array.
	assert.Empty(t, P.Head(result), "Partition left should be empty for nil input")
	assert.Empty(t, P.Tail(result), "Partition right should be empty for nil input")
}

// TestNilSlice_IsNil verifies that IsNil handles nil slices correctly
func TestNilSlice_IsNil(t *testing.T) {
	var nilSlice []int
	assert.True(t, IsNil(nilSlice), "IsNil should return true for nil slice")

	nonNilSlice := []int{}
	assert.False(t, IsNil(nonNilSlice), "IsNil should return false for non-nil empty slice")
}

// TestNilSlice_IsNonNil verifies that IsNonNil handles nil slices correctly
func TestNilSlice_IsNonNil(t *testing.T) {
	var nilSlice []int
	assert.False(t, IsNonNil(nilSlice), "IsNonNil should return false for nil slice")

	nonNilSlice := []int{}
	assert.True(t, IsNonNil(nonNilSlice), "IsNonNil should return true for non-nil empty slice")
}

// TestNilSlice_Copy verifies that Copy handles nil slices correctly
func TestNilSlice_Copy(t *testing.T) {
	var nilSlice []int
	result := Copy(nilSlice)
	assert.NotNil(t, result, "Copy should return non-nil slice")
	assert.Equal(t, 0, len(result), "Copy should return empty slice for nil input")
}

// TestNilSlice_FoldMap verifies that FoldMap handles nil slices correctly
func TestNilSlice_FoldMap(t *testing.T) {
	var nilSlice []int
	monoid := S.Monoid
	foldMap := FoldMap[int](monoid)(strconv.Itoa)
	result := foldMap(nilSlice)
	assert.Equal(t, "", result, "FoldMap should return empty value for nil slice")
}

// TestNilSlice_FoldMapWithIndex verifies that FoldMapWithIndex handles nil slices correctly
func TestNilSlice_FoldMapWithIndex(t *testing.T) {
	var nilSlice []int
	monoid := S.Monoid
	foldMap := FoldMapWithIndex[int](monoid)(func(i int, v int) string {
		return fmt.Sprintf("%d:%d", i, v)
	})
	result := foldMap(nilSlice)
	assert.Equal(t, "", result, "FoldMapWithIndex should return empty value for nil slice")
}

// TestNilSlice_Fold verifies that Fold handles nil slices correctly
func TestNilSlice_Fold(t *testing.T) {
	var nilSlice []string
	monoid := S.Monoid
	fold := Fold(monoid)
	result := fold(nilSlice)
	assert.Equal(t, "", result, "Fold should return empty value for nil slice")
}

// TestNilSlice_Concat verifies that Concat handles nil slices correctly
func TestNilSlice_Concat(t *testing.T) {
	var nilSlice []int
	nonNilSlice := []int{1, 2, 3}

	// nil concat non-nil
	concat1 := Concat(nonNilSlice)
	result1 := concat1(nilSlice)
	assert.Equal(t, nonNilSlice, result1, "nil concat non-nil should return non-nil slice")

	// non-nil concat nil
	concat2 := Concat(nilSlice)
	result2 := concat2(nonNilSlice)
	assert.Equal(t, nonNilSlice, result2, "non-nil concat nil should return non-nil slice")

	// nil concat nil
	concat3 := Concat(nilSlice)
	result3 := concat3(nilSlice)
	assert.Nil(t, result3, "nil concat nil should return nil")
}

// TestNilSlice_MonadFlap verifies that MonadFlap handles nil slices correctly
func TestNilSlice_MonadFlap(t *testing.T) {
	var nilSlice []func(int) string
	result := MonadFlap(nilSlice, 42)
	assert.Empty(t, result, "MonadFlap should return empty slice for nil input")
}

// TestNilSlice_Flap verifies that Flap handles nil slices correctly
func TestNilSlice_Flap(t *testing.T) {
	var nilSlice []func(int) string
	flap := Flap[string](42)
	result := flap(nilSlice)
	assert.Empty(t, result, "Flap should return empty slice for nil input")
}

// TestNilSlice_Reverse verifies that Reverse handles nil slices correctly
func TestNilSlice_Reverse(t *testing.T) {
	var nilSlice []int
	result := Reverse(nilSlice)
	assert.Nil(t, result, "Reverse should return nil for nil slice")
}

// TestNilSlice_Extend verifies that Extend handles nil slices correctly
func TestNilSlice_Extend(t *testing.T) {
	var nilSlice []int
	extend := Extend(func(as []int) string {
		return fmt.Sprintf("%v", as)
	})
	result := extend(nilSlice)
	assert.True(t, IsEmpty(result), "Extend should return empty slice for nil input")
}

// TestNilSlice_Empty verifies that Empty returns a nil (empty) slice
func TestNilSlice_Empty(t *testing.T) {
	result := Empty[int]()
	assert.True(t, IsEmpty(result), "Empty should return an empty slice")
	assert.True(t, IsNil(result), "Empty should return nil")
}

// TestNilSlice_Zero verifies that Zero returns a nil (empty) slice
func TestNilSlice_Zero(t *testing.T) {
	result := Zero[int]()
	assert.True(t, IsEmpty(result), "Zero should return an empty slice")
	assert.True(t, IsNil(result), "Zero should return nil")
}

// TestNilSlice_ConstNil verifies that ConstNil returns a nil slice
func TestNilSlice_ConstNil(t *testing.T) {
	result := ConstNil[int]()
	assert.Nil(t, result, "ConstNil should return nil slice")
	assert.True(t, IsNil(result), "ConstNil should return nil slice")
}

// TestNilSlice_Of verifies that Of creates a proper singleton slice
func TestNilSlice_Of(t *testing.T) {
	result := Of(42)
	assert.NotNil(t, result, "Of should return non-nil slice")
	assert.Equal(t, 1, len(result), "Of should create slice with one element")
	assert.Equal(t, 42, result[0], "Of should set value correctly")
}

// TestNilSlice_EmptyResults verifies that operations producing an empty array
// return an empty array
func TestNilSlice_EmptyResults(t *testing.T) {
	data := []int{1, 2, 3}
	isNeg := func(x int) bool { return x < 0 }

	assert.True(t, IsEmpty(Slice[int](1, 1)(data)), "Slice with empty range")
	assert.True(t, IsEmpty(Slice[int](5, 10)(data)), "Slice beyond the end")
	assert.True(t, IsEmpty(SliceRight[int](3)(data)), "SliceRight at the end")
	assert.True(t, IsEmpty(SliceRight[int](5)(data)), "SliceRight beyond the end")
	assert.True(t, IsEmpty(MakeBy(0, func(i int) int { return i })), "MakeBy(0)")
	assert.True(t, IsEmpty(Replicate(0, 1)), "Replicate(0)")
	assert.True(t, IsEmpty(FromOption(O.None[int]())), "FromOption(None)")
	assert.True(t, IsEmpty(Filter(isNeg)(data)), "Filter without matches")
	assert.True(t, IsEmpty(FilterWithIndex(func(_ int, x int) bool { return isNeg(x) })(data)), "FilterWithIndex without matches")
	assert.True(t, IsEmpty(FilterRef(func(x *int) bool { return isNeg(*x) })(data)), "FilterRef without matches")
	assert.True(t, IsEmpty(FilterMap(func(int) Option[int] { return O.None[int]() })(data)), "FilterMap without matches")
	assert.True(t, IsEmpty(FilterMapWithIndex(func(int, int) Option[int] { return O.None[int]() })(data)), "FilterMapWithIndex without matches")
	assert.True(t, IsEmpty(FilterMapRef(func(x *int) bool { return isNeg(*x) }, func(x *int) int { return *x })(data)), "FilterMapRef without matches")
	assert.True(t, IsEmpty(Extend(Size[int])(Empty[int]())), "Extend of empty array")
	assert.True(t, IsEmpty(Monoid[int]().Empty()), "Monoid empty")
}

// TestNilSlice_TraverseEmpty verifies that traversing an empty array yields an empty array
func TestNilSlice_TraverseEmpty(t *testing.T) {
	result := MonadTraverse(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		Empty[int](),
		func(n int) O.Option[string] { return O.Some(strconv.Itoa(n)) },
	)
	assert.True(t, O.IsSome(result))
	assert.True(t, IsEmpty(O.GetOrElse(func() []string { return []string{"fallback"} })(result)))
}
