// All rights reserved.
// Copyright (c) 2023 - 2025 IBM Corp.
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
	"sort"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/internal/utils"
	Mg "github.com/IBM/fp-go/v2/magma"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Keys(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestValues(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Values(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"A", "B", "C"}, keys)
}

func TestMap(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	assert.Equal(t, expected, Map[string](utils.Upper)(data))
}

func TestLookup(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	assert.Equal(t, O.Some("a"), Lookup[string]("a")(data))
	assert.Equal(t, O.None[string](), Lookup[string]("a1")(data))
}

func TestFilterChain(t *testing.T) {
	src := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	f := func(k string, value int) O.Option[map[string]string] {
		if value%2 != 0 {
			return O.Of(map[string]string{
				k: fmt.Sprintf("%s%d", k, value),
			})
		}
		return O.None[map[string]string]()
	}

	// monoid
	monoid := MergeMonoid[string, string]()

	res := FilterChainWithIndex[int](monoid)(f)(src)

	assert.Equal(t, map[string]string{
		"a": "a1",
		"c": "c3",
	}, res)
}

func ExampleFoldMap() {
	src := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	fold := FoldMapOrd[string, string](S.Ord)(S.Monoid)(strings.ToUpper)

	fmt.Println(fold(src))

	// Output: ABC

}

func ExampleValuesOrd() {
	src := map[string]string{
		"c": "a",
		"b": "b",
		"a": "c",
	}

	getValues := ValuesOrd[string](S.Ord)

	fmt.Println(getValues(src))

	// Output: [c b a]

}

func TestCopyVsClone(t *testing.T) {
	slc := []string{"b", "c"}
	src := map[string][]string{
		"a": slc,
	}
	// make a shallow copy
	cpy := Copy(src)
	// make a deep copy
	cln := Clone[string](A.Copy[string])(src)

	assert.Equal(t, cpy, cln)
	// make a modification to the original slice
	slc[0] = "d"
	assert.NotEqual(t, cpy, cln)
	assert.Equal(t, src, cpy)
}

func TestFromArrayMap(t *testing.T) {
	src1 := A.From("a", "b", "c", "a")
	frm := FromArrayMap[string, string](Mg.Second[string]())

	f := frm(P.Of[string])

	res1 := f(src1)

	assert.Equal(t, map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}, res1)

	src2 := A.From("A", "B", "C", "A")

	res2 := f(src2)

	assert.Equal(t, map[string]string{
		"A": "A",
		"B": "B",
		"C": "C",
	}, res2)
}

func TestEmpty(t *testing.T) {
	nonEmpty := map[string]string{
		"a": "A",
		"b": "B",
	}
	empty := Empty[string, string]()

	assert.True(t, IsEmpty(empty))
	assert.False(t, IsEmpty(nonEmpty))
	assert.False(t, IsNonEmpty(empty))
	assert.True(t, IsNonEmpty(nonEmpty))
}

func TestHas(t *testing.T) {
	nonEmpty := map[string]string{
		"a": "A",
		"b": "B",
	}
	assert.True(t, Has("a", nonEmpty))
	assert.False(t, Has("c", nonEmpty))
}

func TestCollect(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	collector := Collect(func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := collector(data)
	sort.Strings(result)
	assert.Equal(t, []string{"a=1", "b=2", "c=3"}, result)
}

func TestCollectOrd(t *testing.T) {
	data := map[string]int{
		"c": 3,
		"a": 1,
		"b": 2,
	}
	collector := CollectOrd[int, string](S.Ord)(func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := collector(data)
	assert.Equal(t, []string{"a=1", "b=2", "c=3"}, result)
}

func TestReduce(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	sum := Reduce[string](func(acc, v int) int {
		return acc + v
	}, 0)
	result := sum(data)
	assert.Equal(t, 6, result)
}

func TestReduceWithIndex(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	concat := ReduceWithIndex(func(k string, acc string, v int) string {
		if acc == "" {
			return fmt.Sprintf("%s:%d", k, v)
		}
		return fmt.Sprintf("%s,%s:%d", acc, k, v)
	}, "")
	result := concat(data)
	// Result order is non-deterministic, so check it contains all parts
	assert.Contains(t, result, "a:1")
	assert.Contains(t, result, "b:2")
	assert.Contains(t, result, "c:3")
}

func TestMonadMap(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	result := MonadMap(data, func(v int) int { return v * 2 })
	assert.Equal(t, map[string]int{"a": 2, "b": 4, "c": 6}, result)
}

func TestMonadMapWithIndex(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	result := MonadMapWithIndex(data, func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	assert.Equal(t, map[string]string{"a": "a=1", "b": "b=2"}, result)
}

func TestMapWithIndex(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	mapper := MapWithIndex(func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	result := mapper(data)
	assert.Equal(t, map[string]string{"a": "a=1", "b": "b=2"}, result)
}

func TestMonadLookup(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	assert.Equal(t, O.Some(1), MonadLookup(data, "a"))
	assert.Equal(t, O.None[int](), MonadLookup(data, "c"))
}

func TestMerge(t *testing.T) {
	left := map[string]int{"a": 1, "b": 2}
	right := map[string]int{"b": 3, "c": 4}
	result := Merge(right)(left)
	assert.Equal(t, map[string]int{"a": 1, "b": 3, "c": 4}, result)
}

func TestSize(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.Equal(t, 3, Size(data))
	assert.Equal(t, 0, Size(Empty[string, int]()))
}

func TestToArray(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := ToArray(data)
	assert.Len(t, result, 2)
	// Check both entries exist (order is non-deterministic)
	found := make(map[string]int)
	for _, entry := range result {
		found[P.Head(entry)] = P.Tail(entry)
	}
	assert.Equal(t, data, found)
}

func TestToEntries(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := ToEntries(data)
	assert.Len(t, result, 2)
}

func TestFromEntries(t *testing.T) {
	entries := Entries[string, int]{
		P.MakePair("a", 1),
		P.MakePair("b", 2),
	}
	result := FromEntries(entries)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)
}

func TestUpsertAt(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := UpsertAt("c", 3)(data)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, result)
	// Original should be unchanged
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, data)

	// Update existing
	result2 := UpsertAt("a", 10)(data)
	assert.Equal(t, map[string]int{"a": 10, "b": 2}, result2)
}

func TestDeleteAt(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	result := DeleteAt[string, int]("b")(data)
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, result)
	// Original should be unchanged
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, data)
}

func TestSingleton(t *testing.T) {
	result := Singleton("key", 42)
	assert.Equal(t, map[string]int{"key": 42}, result)
}

func TestFilterMapWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	filter := FilterMapWithIndex(func(k string, v int) O.Option[int] {
		if v%2 == 0 {
			return O.Some(v * 10)
		}
		return O.None[int]()
	})
	result := filter(data)
	assert.Equal(t, map[string]int{"b": 20}, result)
}

func TestFilterMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	filter := FilterMap[string](func(v int) O.Option[int] {
		if v%2 == 0 {
			return O.Some(v * 10)
		}
		return O.None[int]()
	})
	result := filter(data)
	assert.Equal(t, map[string]int{"b": 20}, result)
}

func TestFilter(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	filter := Filter[string, int](func(k string) bool {
		return k != "b"
	})
	result := filter(data)
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, result)
}

func TestFilterWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	filter := FilterWithIndex(func(k string, v int) bool {
		return v%2 == 0
	})
	result := filter(data)
	assert.Equal(t, map[string]int{"b": 2}, result)
}

func TestIsNil(t *testing.T) {
	var nilMap map[string]int
	nonNilMap := map[string]int{}

	assert.True(t, IsNil(nilMap))
	assert.False(t, IsNil(nonNilMap))
}

func TestIsNonNil(t *testing.T) {
	var nilMap map[string]int
	nonNilMap := map[string]int{}

	assert.False(t, IsNonNil(nilMap))
	assert.True(t, IsNonNil(nonNilMap))
}

func TestConstNil(t *testing.T) {
	result := ConstNil[string, int]()
	assert.Nil(t, result)
	assert.True(t, IsNil(result))
}

func TestMonadChain(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	monoid := MergeMonoid[string, int]()
	result := MonadChain(monoid, data, func(v int) map[string]int {
		return map[string]int{
			fmt.Sprintf("x%d", v): v * 10,
		}
	})
	assert.Equal(t, map[string]int{"x1": 10, "x2": 20}, result)
}

func TestChain(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	monoid := MergeMonoid[string, int]()
	chain := Chain[int](monoid)(func(v int) map[string]int {
		return map[string]int{
			fmt.Sprintf("x%d", v): v * 10,
		}
	})
	result := chain(data)
	assert.Equal(t, map[string]int{"x1": 10, "x2": 20}, result)
}

func TestFlatten(t *testing.T) {
	nested := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"z": 3},
	}
	monoid := MergeMonoid[string, int]()
	flatten := Flatten(monoid)
	result := flatten(nested)
	assert.Equal(t, map[string]int{"x": 1, "y": 2, "z": 3}, result)
}

func TestFoldMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	// Use string monoid for simplicity
	fold := FoldMap[string, int](S.Monoid)(func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	result := fold(data)
	// Result contains all digits but order is non-deterministic
	assert.Contains(t, result, "1")
	assert.Contains(t, result, "2")
	assert.Contains(t, result, "3")
}

func TestFold(t *testing.T) {
	data := map[string]string{"a": "A", "b": "B", "c": "C"}
	fold := Fold[string](S.Monoid)
	result := fold(data)
	// Result contains all letters but order is non-deterministic
	assert.Contains(t, result, "A")
	assert.Contains(t, result, "B")
	assert.Contains(t, result, "C")
}

func TestKeysOrd(t *testing.T) {
	data := map[string]int{"c": 3, "a": 1, "b": 2}
	keys := KeysOrd[int](S.Ord)(data)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestMonadFlap(t *testing.T) {
	fns := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
		"triple": func(x int) int { return x * 3 },
	}
	result := MonadFlap(fns, 5)
	assert.Equal(t, map[string]int{"double": 10, "triple": 15}, result)
}

func TestFlap(t *testing.T) {
	fns := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
		"triple": func(x int) int { return x * 3 },
	}
	flap := Flap[int, string](5)
	result := flap(fns)
	assert.Equal(t, map[string]int{"double": 10, "triple": 15}, result)
}

func TestFromArray(t *testing.T) {
	entries := Entries[string, int]{
		P.MakePair("a", 1),
		P.MakePair("b", 2),
		P.MakePair("a", 3), // Duplicate key
	}
	// Use Second magma to keep last value
	from := FromArray[string](Mg.Second[int]())
	result := from(entries)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, result)
}

func TestMonadAp(t *testing.T) {
	fns := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
	}
	vals := map[string]int{
		"double": 5,
	}
	monoid := MergeMonoid[string, int]()
	result := MonadAp(monoid, fns, vals)
	assert.Equal(t, map[string]int{"double": 10}, result)
}

func TestAp(t *testing.T) {
	fns := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
	}
	vals := map[string]int{
		"double": 5,
	}
	monoid := MergeMonoid[string, int]()
	ap := Ap[int](monoid)(vals)
	result := ap(fns)
	assert.Equal(t, map[string]int{"double": 10}, result)
}

func TestOf(t *testing.T) {
	result := Of("key", 42)
	assert.Equal(t, map[string]int{"key": 42}, result)
}

func TestReduceRef(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := ReduceRef[string](func(acc int, v *int) int {
		return acc + *v
	}, 0)
	result := sum(data)
	assert.Equal(t, 6, result)
}

func TestReduceRefWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	concat := ReduceRefWithIndex(func(k string, acc string, v *int) string {
		if acc == "" {
			return fmt.Sprintf("%s:%d", k, *v)
		}
		return fmt.Sprintf("%s,%s:%d", acc, k, *v)
	}, "")
	result := concat(data)
	assert.Contains(t, result, "a:1")
	assert.Contains(t, result, "b:2")
}

func TestMonadMapRef(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := MonadMapRef(data, func(v *int) int { return *v * 2 })
	assert.Equal(t, map[string]int{"a": 2, "b": 4}, result)
}

func TestMapRef(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	mapper := MapRef[string](func(v *int) int { return *v * 2 })
	result := mapper(data)
	assert.Equal(t, map[string]int{"a": 2, "b": 4}, result)
}

func TestMonadMapRefWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := MonadMapRefWithIndex(data, func(k string, v *int) string {
		return fmt.Sprintf("%s=%d", k, *v)
	})
	assert.Equal(t, map[string]string{"a": "a=1", "b": "b=2"}, result)
}

func TestMapRefWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	mapper := MapRefWithIndex(func(k string, v *int) string {
		return fmt.Sprintf("%s=%d", k, *v)
	})
	result := mapper(data)
	assert.Equal(t, map[string]string{"a": "a=1", "b": "b=2"}, result)
}

func TestUnion(t *testing.T) {
	left := map[string]int{"a": 1, "b": 2}
	right := map[string]int{"b": 3, "c": 4}
	// Union combines maps, with the magma resolving conflicts
	// The order is union(left)(right), which means right is merged into left
	// First magma keeps the first value (from right in this case)
	union := Union[string](Mg.First[int]())
	result := union(left)(right)
	assert.Equal(t, map[string]int{"a": 1, "b": 3, "c": 4}, result)

	// Second magma keeps the second value (from left in this case)
	union2 := Union[string](Mg.Second[int]())
	result2 := union2(left)(right)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 4}, result2)
}

func TestMonadChainWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	monoid := MergeMonoid[string, string]()
	result := MonadChainWithIndex(monoid, data, func(k string, v int) map[string]string {
		return map[string]string{
			fmt.Sprintf("%s%d", k, v): fmt.Sprintf("val%d", v),
		}
	})
	assert.Equal(t, map[string]string{"a1": "val1", "b2": "val2"}, result)
}

func TestChainWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	monoid := MergeMonoid[string, string]()
	chain := ChainWithIndex[int](monoid)(func(k string, v int) map[string]string {
		return map[string]string{
			fmt.Sprintf("%s%d", k, v): fmt.Sprintf("val%d", v),
		}
	})
	result := chain(data)
	assert.Equal(t, map[string]string{"a1": "val1", "b2": "val2"}, result)
}

func TestFilterChainWithIndex(t *testing.T) {
	src := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	f := func(k string, value int) O.Option[map[string]string] {
		if value%2 != 0 {
			return O.Of(map[string]string{
				k: fmt.Sprintf("%s%d", k, value),
			})
		}
		return O.None[map[string]string]()
	}

	monoid := MergeMonoid[string, string]()
	res := FilterChainWithIndex[int](monoid)(f)(src)

	assert.Equal(t, map[string]string{
		"a": "a1",
		"c": "c3",
	}, res)
}

func TestFoldMapWithIndex(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	fold := FoldMapWithIndex[string, int](S.Monoid)(func(k string, v int) string {
		return fmt.Sprintf("%s:%d", k, v)
	})
	result := fold(data)
	// Result contains all pairs but order is non-deterministic
	assert.Contains(t, result, "a:1")
	assert.Contains(t, result, "b:2")
	assert.Contains(t, result, "c:3")
}

func TestReduceOrd(t *testing.T) {
	data := map[string]int{"c": 3, "a": 1, "b": 2}
	sum := ReduceOrd[int, int](S.Ord)(func(acc, v int) int {
		return acc + v
	}, 0)
	result := sum(data)
	assert.Equal(t, 6, result)
}

func TestReduceOrdWithIndex(t *testing.T) {
	data := map[string]int{"c": 3, "a": 1, "b": 2}
	concat := ReduceOrdWithIndex[int, string](S.Ord)(func(k string, acc string, v int) string {
		if acc == "" {
			return fmt.Sprintf("%s:%d", k, v)
		}
		return fmt.Sprintf("%s,%s:%d", acc, k, v)
	}, "")
	result := concat(data)
	// With Ord, keys should be in order
	assert.Equal(t, "a:1,b:2,c:3", result)
}

func TestFoldMapOrdWithIndex(t *testing.T) {
	data := map[string]int{"c": 3, "a": 1, "b": 2}
	fold := FoldMapOrdWithIndex[string, int, string](S.Ord)(S.Monoid)(func(k string, v int) string {
		return fmt.Sprintf("%s:%d,", k, v)
	})
	result := fold(data)
	assert.Equal(t, "a:1,b:2,c:3,", result)
}

func TestFoldOrd(t *testing.T) {
	data := map[string]string{"c": "C", "a": "A", "b": "B"}
	fold := FoldOrd[string](S.Ord)(S.Monoid)
	result := fold(data)
	assert.Equal(t, "ABC", result)
}

func TestFromFoldableMap(t *testing.T) {
	src := A.From("a", "b", "c", "a")
	// Create a reducer function
	reducer := A.Reduce[string, map[string]string]
	from := FromFoldableMap(
		Mg.Second[string](),
		reducer,
	)
	f := from(P.Of[string])
	result := f(src)
	assert.Equal(t, map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}, result)
}

func TestFromFoldable(t *testing.T) {
	entries := Entries[string, int]{
		P.MakePair("a", 1),
		P.MakePair("b", 2),
		P.MakePair("a", 3), // Duplicate key
	}
	reducer := A.Reduce[Entry[string, int], map[string]int]
	from := FromFoldable(
		Mg.Second[int](),
		reducer,
	)
	result := from(entries)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, result)
}
