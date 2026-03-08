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

package record

import (
	"fmt"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	SG "github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestNilMap_IsEmpty verifies that IsEmpty handles nil maps correctly
func TestNilMap_IsEmpty(t *testing.T) {
	var nilMap Record[string, int]
	assert.True(t, IsEmpty(nilMap), "nil map should be empty")
}

// TestNilMap_IsNonEmpty verifies that IsNonEmpty handles nil maps correctly
func TestNilMap_IsNonEmpty(t *testing.T) {
	var nilMap Record[string, int]
	assert.False(t, IsNonEmpty(nilMap), "nil map should not be non-empty")
}

// TestNilMap_Keys verifies that Keys handles nil maps correctly
func TestNilMap_Keys(t *testing.T) {
	var nilMap Record[string, int]
	keys := Keys(nilMap)
	assert.NotNil(t, keys, "Keys should return non-nil slice")
	assert.Equal(t, 0, len(keys), "Keys should return empty slice for nil map")
}

// TestNilMap_Values verifies that Values handles nil maps correctly
func TestNilMap_Values(t *testing.T) {
	var nilMap Record[string, int]
	values := Values(nilMap)
	assert.NotNil(t, values, "Values should return non-nil slice")
	assert.Equal(t, 0, len(values), "Values should return empty slice for nil map")
}

// TestNilMap_Collect verifies that Collect handles nil maps correctly
func TestNilMap_Collect(t *testing.T) {
	var nilMap Record[string, int]
	collector := Collect(func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := collector(nilMap)
	assert.NotNil(t, result, "Collect should return non-nil slice")
	assert.Equal(t, 0, len(result), "Collect should return empty slice for nil map")
}

// TestNilMap_Reduce verifies that Reduce handles nil maps correctly
func TestNilMap_Reduce(t *testing.T) {
	var nilMap Record[string, int]
	reducer := Reduce[string](func(acc int, v int) int {
		return acc + v
	}, 10)
	result := reducer(nilMap)
	assert.Equal(t, 10, result, "Reduce should return initial value for nil map")
}

// TestNilMap_ReduceWithIndex verifies that ReduceWithIndex handles nil maps correctly
func TestNilMap_ReduceWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	reducer := ReduceWithIndex(func(k string, acc int, v int) int {
		return acc + v
	}, 10)
	result := reducer(nilMap)
	assert.Equal(t, 10, result, "ReduceWithIndex should return initial value for nil map")
}

// TestNilMap_ReduceRef verifies that ReduceRef handles nil maps correctly
func TestNilMap_ReduceRef(t *testing.T) {
	var nilMap Record[string, int]
	reducer := ReduceRef[string](func(acc int, v *int) int {
		return acc + *v
	}, 10)
	result := reducer(nilMap)
	assert.Equal(t, 10, result, "ReduceRef should return initial value for nil map")
}

// TestNilMap_ReduceRefWithIndex verifies that ReduceRefWithIndex handles nil maps correctly
func TestNilMap_ReduceRefWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	reducer := ReduceRefWithIndex(func(k string, acc int, v *int) int {
		return acc + *v
	}, 10)
	result := reducer(nilMap)
	assert.Equal(t, 10, result, "ReduceRefWithIndex should return initial value for nil map")
}

// TestNilMap_MonadMap verifies that MonadMap handles nil maps correctly
func TestNilMap_MonadMap(t *testing.T) {
	var nilMap Record[string, int]
	result := MonadMap(nilMap, func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	assert.NotNil(t, result, "MonadMap should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadMap should return empty map for nil input")
}

// TestNilMap_MonadMapWithIndex verifies that MonadMapWithIndex handles nil maps correctly
func TestNilMap_MonadMapWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	result := MonadMapWithIndex(nilMap, func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	assert.NotNil(t, result, "MonadMapWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadMapWithIndex should return empty map for nil input")
}

// TestNilMap_MonadMapRefWithIndex verifies that MonadMapRefWithIndex handles nil maps correctly
func TestNilMap_MonadMapRefWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	result := MonadMapRefWithIndex(nilMap, func(k string, v *int) string {
		return fmt.Sprintf("%s=%d", k, *v)
	})
	assert.NotNil(t, result, "MonadMapRefWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadMapRefWithIndex should return empty map for nil input")
}

// TestNilMap_MonadMapRef verifies that MonadMapRef handles nil maps correctly
func TestNilMap_MonadMapRef(t *testing.T) {
	var nilMap Record[string, int]
	result := MonadMapRef(nilMap, func(v *int) string {
		return fmt.Sprintf("%d", *v)
	})
	assert.NotNil(t, result, "MonadMapRef should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadMapRef should return empty map for nil input")
}

// TestNilMap_Map verifies that Map handles nil maps correctly
func TestNilMap_Map(t *testing.T) {
	var nilMap Record[string, int]
	mapper := Map[string](func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	result := mapper(nilMap)
	assert.NotNil(t, result, "Map should return non-nil map")
	assert.Equal(t, 0, len(result), "Map should return empty map for nil input")
}

// TestNilMap_MapRef verifies that MapRef handles nil maps correctly
func TestNilMap_MapRef(t *testing.T) {
	var nilMap Record[string, int]
	mapper := MapRef[string](func(v *int) string {
		return fmt.Sprintf("%d", *v)
	})
	result := mapper(nilMap)
	assert.NotNil(t, result, "MapRef should return non-nil map")
	assert.Equal(t, 0, len(result), "MapRef should return empty map for nil input")
}

// TestNilMap_MapWithIndex verifies that MapWithIndex handles nil maps correctly
func TestNilMap_MapWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	mapper := MapWithIndex[string](func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := mapper(nilMap)
	assert.NotNil(t, result, "MapWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "MapWithIndex should return empty map for nil input")
}

// TestNilMap_MapRefWithIndex verifies that MapRefWithIndex handles nil maps correctly
func TestNilMap_MapRefWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	mapper := MapRefWithIndex[string](func(k string, v *int) string {
		return fmt.Sprintf("%s=%d", k, *v)
	})
	result := mapper(nilMap)
	assert.NotNil(t, result, "MapRefWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "MapRefWithIndex should return empty map for nil input")
}

// TestNilMap_Lookup verifies that Lookup handles nil maps correctly
func TestNilMap_Lookup(t *testing.T) {
	var nilMap Record[string, int]
	lookup := Lookup[int]("key")
	result := lookup(nilMap)
	assert.True(t, O.IsNone(result), "Lookup should return None for nil map")
}

// TestNilMap_MonadLookup verifies that MonadLookup handles nil maps correctly
func TestNilMap_MonadLookup(t *testing.T) {
	var nilMap Record[string, int]
	result := MonadLookup(nilMap, "key")
	assert.True(t, O.IsNone(result), "MonadLookup should return None for nil map")
}

// TestNilMap_Has verifies that Has handles nil maps correctly
func TestNilMap_Has(t *testing.T) {
	var nilMap Record[string, int]
	result := Has("key", nilMap)
	assert.False(t, result, "Has should return false for nil map")
}

// TestNilMap_Union verifies that Union handles nil maps correctly
func TestNilMap_Union(t *testing.T) {
	var nilMap Record[string, int]
	nonNilMap := Record[string, int]{"a": 1, "b": 2}

	semigroup := SG.Last[int]()
	union := Union[string](semigroup)

	// nil union non-nil
	result1 := union(nonNilMap)(nilMap)
	assert.Equal(t, nonNilMap, result1, "nil union non-nil should return non-nil map")

	// non-nil union nil
	result2 := union(nilMap)(nonNilMap)
	assert.Equal(t, nonNilMap, result2, "non-nil union nil should return non-nil map")

	// nil union nil - returns nil when both inputs are nil (optimization)
	result3 := union(nilMap)(nilMap)
	assert.Nil(t, result3, "nil union nil returns nil")
}

// TestNilMap_Merge verifies that Merge handles nil maps correctly
func TestNilMap_Merge(t *testing.T) {
	var nilMap Record[string, int]
	nonNilMap := Record[string, int]{"a": 1, "b": 2}

	// nil merge non-nil
	result1 := Merge(nonNilMap)(nilMap)
	assert.Equal(t, nonNilMap, result1, "nil merge non-nil should return non-nil map")

	// non-nil merge nil
	result2 := Merge(nilMap)(nonNilMap)
	assert.Equal(t, nonNilMap, result2, "non-nil merge nil should return non-nil map")

	// nil merge nil - returns nil when both inputs are nil (optimization)
	result3 := Merge(nilMap)(nilMap)
	assert.Nil(t, result3, "nil merge nil returns nil")
}

// TestNilMap_Size verifies that Size handles nil maps correctly
func TestNilMap_Size(t *testing.T) {
	var nilMap Record[string, int]
	result := Size(nilMap)
	assert.Equal(t, 0, result, "Size should return 0 for nil map")
}

// TestNilMap_ToArray verifies that ToArray handles nil maps correctly
func TestNilMap_ToArray(t *testing.T) {
	var nilMap Record[string, int]
	result := ToArray(nilMap)
	assert.NotNil(t, result, "ToArray should return non-nil slice")
	assert.Equal(t, 0, len(result), "ToArray should return empty slice for nil map")
}

// TestNilMap_ToEntries verifies that ToEntries handles nil maps correctly
func TestNilMap_ToEntries(t *testing.T) {
	var nilMap Record[string, int]
	result := ToEntries(nilMap)
	assert.NotNil(t, result, "ToEntries should return non-nil slice")
	assert.Equal(t, 0, len(result), "ToEntries should return empty slice for nil map")
}

// TestNilMap_UpsertAt verifies that UpsertAt handles nil maps correctly
func TestNilMap_UpsertAt(t *testing.T) {
	var nilMap Record[string, int]
	upsert := UpsertAt("key", 42)
	result := upsert(nilMap)
	assert.NotNil(t, result, "UpsertAt should return non-nil map")
	assert.Equal(t, 1, len(result), "UpsertAt should create map with one entry")
	assert.Equal(t, 42, result["key"], "UpsertAt should insert value correctly")
}

// TestNilMap_DeleteAt verifies that DeleteAt handles nil maps correctly
func TestNilMap_DeleteAt(t *testing.T) {
	var nilMap Record[string, int]
	deleteFunc := DeleteAt[string, int]("key")
	result := deleteFunc(nilMap)
	assert.NotNil(t, result, "DeleteAt should return non-nil map")
	assert.Equal(t, 0, len(result), "DeleteAt should return empty map for nil input")
}

// TestNilMap_Filter verifies that Filter handles nil maps correctly
func TestNilMap_Filter(t *testing.T) {
	var nilMap Record[string, int]
	filter := Filter[string, int](func(k string) bool {
		return true
	})
	result := filter(nilMap)
	assert.NotNil(t, result, "Filter should return non-nil map")
	assert.Equal(t, 0, len(result), "Filter should return empty map for nil input")
}

// TestNilMap_FilterWithIndex verifies that FilterWithIndex handles nil maps correctly
func TestNilMap_FilterWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	filter := FilterWithIndex[string, int](func(k string, v int) bool {
		return true
	})
	result := filter(nilMap)
	assert.NotNil(t, result, "FilterWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "FilterWithIndex should return empty map for nil input")
}

// TestNilMap_IsNil verifies that IsNil handles nil maps correctly
func TestNilMap_IsNil(t *testing.T) {
	var nilMap Record[string, int]
	assert.True(t, IsNil(nilMap), "IsNil should return true for nil map")

	nonNilMap := Record[string, int]{}
	assert.False(t, IsNil(nonNilMap), "IsNil should return false for non-nil empty map")
}

// TestNilMap_IsNonNil verifies that IsNonNil handles nil maps correctly
func TestNilMap_IsNonNil(t *testing.T) {
	var nilMap Record[string, int]
	assert.False(t, IsNonNil(nilMap), "IsNonNil should return false for nil map")

	nonNilMap := Record[string, int]{}
	assert.True(t, IsNonNil(nonNilMap), "IsNonNil should return true for non-nil empty map")
}

// TestNilMap_MonadChainWithIndex verifies that MonadChainWithIndex handles nil maps correctly
func TestNilMap_MonadChainWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	monoid := MergeMonoid[string, string]()
	result := MonadChainWithIndex(monoid, nilMap, func(k string, v int) Record[string, string] {
		return Record[string, string]{k: fmt.Sprintf("%d", v)}
	})
	assert.NotNil(t, result, "MonadChainWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadChainWithIndex should return empty map for nil input")
}

// TestNilMap_MonadChain verifies that MonadChain handles nil maps correctly
func TestNilMap_MonadChain(t *testing.T) {
	var nilMap Record[string, int]
	monoid := MergeMonoid[string, string]()
	result := MonadChain(monoid, nilMap, func(v int) Record[string, string] {
		return Record[string, string]{"key": fmt.Sprintf("%d", v)}
	})
	assert.NotNil(t, result, "MonadChain should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadChain should return empty map for nil input")
}

// TestNilMap_ChainWithIndex verifies that ChainWithIndex handles nil maps correctly
func TestNilMap_ChainWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	monoid := MergeMonoid[string, string]()
	chain := ChainWithIndex[int, string](monoid)(func(k string, v int) Record[string, string] {
		return Record[string, string]{k: fmt.Sprintf("%d", v)}
	})
	result := chain(nilMap)
	assert.NotNil(t, result, "ChainWithIndex should return non-nil map")
	assert.Equal(t, 0, len(result), "ChainWithIndex should return empty map for nil input")
}

// TestNilMap_Chain verifies that Chain handles nil maps correctly
func TestNilMap_Chain(t *testing.T) {
	var nilMap Record[string, int]
	monoid := MergeMonoid[string, string]()
	chain := Chain[int, string](monoid)(func(v int) Record[string, string] {
		return Record[string, string]{"key": fmt.Sprintf("%d", v)}
	})
	result := chain(nilMap)
	assert.NotNil(t, result, "Chain should return non-nil map")
	assert.Equal(t, 0, len(result), "Chain should return empty map for nil input")
}

// TestNilMap_Flatten verifies that Flatten handles nil maps correctly
func TestNilMap_Flatten(t *testing.T) {
	var nilMap Record[string, Record[string, int]]
	monoid := MergeMonoid[string, int]()
	flatten := Flatten[string, int](monoid)
	result := flatten(nilMap)
	assert.NotNil(t, result, "Flatten should return non-nil map")
	assert.Equal(t, 0, len(result), "Flatten should return empty map for nil input")
}

// TestNilMap_Copy verifies that Copy handles nil maps correctly
func TestNilMap_Copy(t *testing.T) {
	var nilMap Record[string, int]
	result := Copy(nilMap)
	assert.NotNil(t, result, "Copy should return non-nil map")
	assert.Equal(t, 0, len(result), "Copy should return empty map for nil input")
}

// TestNilMap_Clone verifies that Clone handles nil maps correctly
func TestNilMap_Clone(t *testing.T) {
	var nilMap Record[string, int]
	clone := Clone[string, int](func(v int) int { return v * 2 })
	result := clone(nilMap)
	assert.NotNil(t, result, "Clone should return non-nil map")
	assert.Equal(t, 0, len(result), "Clone should return empty map for nil input")
}

// TestNilMap_FromArray verifies that FromArray handles nil/empty arrays correctly
func TestNilMap_FromArray(t *testing.T) {
	semigroup := SG.Last[int]()
	fromArray := FromArray[string, int](semigroup)

	// Test with nil slice
	var nilSlice Entries[string, int]
	result1 := fromArray(nilSlice)
	assert.NotNil(t, result1, "FromArray should return non-nil map for nil slice")
	assert.Equal(t, 0, len(result1), "FromArray should return empty map for nil slice")

	// Test with empty slice
	emptySlice := Entries[string, int]{}
	result2 := fromArray(emptySlice)
	assert.NotNil(t, result2, "FromArray should return non-nil map for empty slice")
	assert.Equal(t, 0, len(result2), "FromArray should return empty map for empty slice")
}

// TestNilMap_MonadAp verifies that MonadAp handles nil maps correctly
func TestNilMap_MonadAp(t *testing.T) {
	var nilFab Record[string, func(int) string]
	var nilFa Record[string, int]
	monoid := MergeMonoid[string, string]()

	// nil functions, nil values
	result1 := MonadAp(monoid, nilFab, nilFa)
	assert.NotNil(t, result1, "MonadAp should return non-nil map")
	assert.Equal(t, 0, len(result1), "MonadAp should return empty map for nil inputs")

	// nil functions, non-nil values
	nonNilFa := Record[string, int]{"a": 1}
	result2 := MonadAp(monoid, nilFab, nonNilFa)
	assert.NotNil(t, result2, "MonadAp should return non-nil map")
	assert.Equal(t, 0, len(result2), "MonadAp should return empty map when functions are nil")

	// non-nil functions, nil values
	nonNilFab := Record[string, func(int) string]{"a": func(v int) string { return fmt.Sprintf("%d", v) }}
	result3 := MonadAp(monoid, nonNilFab, nilFa)
	assert.NotNil(t, result3, "MonadAp should return non-nil map")
	assert.Equal(t, 0, len(result3), "MonadAp should return empty map when values are nil")
}

// TestNilMap_Of verifies that Of creates a proper singleton map
func TestNilMap_Of(t *testing.T) {
	result := Of("key", 42)
	assert.NotNil(t, result, "Of should return non-nil map")
	assert.Equal(t, 1, len(result), "Of should create map with one entry")
	assert.Equal(t, 42, result["key"], "Of should set value correctly")
}

// TestNilMap_FromEntries verifies that FromEntries handles nil/empty slices correctly
func TestNilMap_FromEntries(t *testing.T) {
	// Test with nil slice
	var nilSlice Entries[string, int]
	result1 := FromEntries(nilSlice)
	assert.NotNil(t, result1, "FromEntries should return non-nil map for nil slice")
	assert.Equal(t, 0, len(result1), "FromEntries should return empty map for nil slice")

	// Test with empty slice
	emptySlice := Entries[string, int]{}
	result2 := FromEntries(emptySlice)
	assert.NotNil(t, result2, "FromEntries should return non-nil map for empty slice")
	assert.Equal(t, 0, len(result2), "FromEntries should return empty map for empty slice")

	// Test with actual entries
	entries := Entries[string, int]{
		P.MakePair("a", 1),
		P.MakePair("b", 2),
	}
	result3 := FromEntries(entries)
	assert.NotNil(t, result3, "FromEntries should return non-nil map")
	assert.Equal(t, 2, len(result3), "FromEntries should create map with correct size")
	assert.Equal(t, 1, result3["a"], "FromEntries should set values correctly")
	assert.Equal(t, 2, result3["b"], "FromEntries should set values correctly")
}

// TestNilMap_Singleton verifies that Singleton creates a proper singleton map
func TestNilMap_Singleton(t *testing.T) {
	result := Singleton("key", 42)
	assert.NotNil(t, result, "Singleton should return non-nil map")
	assert.Equal(t, 1, len(result), "Singleton should create map with one entry")
	assert.Equal(t, 42, result["key"], "Singleton should set value correctly")
}

// TestNilMap_Empty verifies that Empty creates an empty non-nil map
func TestNilMap_Empty(t *testing.T) {
	result := Empty[string, int]()
	assert.NotNil(t, result, "Empty should return non-nil map")
	assert.Equal(t, 0, len(result), "Empty should return empty map")
	assert.False(t, IsNil(result), "Empty should not return nil map")
}

// TestNilMap_ConstNil verifies that ConstNil returns a nil map
func TestNilMap_ConstNil(t *testing.T) {
	result := ConstNil[string, int]()
	assert.Nil(t, result, "ConstNil should return nil map")
	assert.True(t, IsNil(result), "ConstNil should return nil map")
}

// TestNilMap_FoldMap verifies that FoldMap handles nil maps correctly
func TestNilMap_FoldMap(t *testing.T) {
	var nilMap Record[string, int]
	monoid := S.Monoid
	foldMap := FoldMap[string, int, string](monoid)(func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	result := foldMap(nilMap)
	assert.Equal(t, "", result, "FoldMap should return empty value for nil map")
}

// TestNilMap_FoldMapWithIndex verifies that FoldMapWithIndex handles nil maps correctly
func TestNilMap_FoldMapWithIndex(t *testing.T) {
	var nilMap Record[string, int]
	monoid := S.Monoid
	foldMap := FoldMapWithIndex[string, int, string](monoid)(func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := foldMap(nilMap)
	assert.Equal(t, "", result, "FoldMapWithIndex should return empty value for nil map")
}

// TestNilMap_Fold verifies that Fold handles nil maps correctly
func TestNilMap_Fold(t *testing.T) {
	var nilMap Record[string, string]
	monoid := S.Monoid
	fold := Fold[string](monoid)
	result := fold(nilMap)
	assert.Equal(t, "", result, "Fold should return empty value for nil map")
}

// TestNilMap_MonadFlap verifies that MonadFlap handles nil maps correctly
func TestNilMap_MonadFlap(t *testing.T) {
	var nilMap Record[string, func(int) string]
	result := MonadFlap(nilMap, 42)
	assert.NotNil(t, result, "MonadFlap should return non-nil map")
	assert.Equal(t, 0, len(result), "MonadFlap should return empty map for nil input")
}
