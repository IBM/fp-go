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

package fibheap

import (
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	"github.com/stretchr/testify/assert"
)

var intOrd = ord.FromStrictCompare[int]()

// ---- helpers ---------------------------------------------------------------

func mustSome[A any](t *testing.T, opt Option[A]) A {
	t.Helper()
	v, ok := O.Unwrap(opt)
	assert.True(t, ok, "expected Some, got None")
	return v
}

func extractAll(t *testing.T, h Heap[int]) []int {
	t.Helper()
	out := make([]int, 0, Size(h))
	for !IsEmpty(h) {
		var v Option[int]
		h, v = ExtractMin[int](intOrd)(h)
		out = append(out, mustSome(t, v))
	}
	return out
}

// ---- Empty / IsEmpty / Size ------------------------------------------------

func TestEmpty(t *testing.T) {
	h := Empty[int]()
	assert.True(t, IsEmpty(h))
	assert.Equal(t, 0, Size(h))
	assert.Equal(t, O.None[int](), FindMin(h))
}

// ---- Insert / FindMin ------------------------------------------------------

func TestInsertSingle(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(42)(h)
	assert.False(t, IsEmpty(h))
	assert.Equal(t, 1, Size(h))
	assert.Equal(t, O.Some(42), FindMin(h))
}

func TestInsertMaintainsMin(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(5)(h)
	h, _ = Insert(intOrd)(3)(h)
	h, _ = Insert(intOrd)(8)(h)
	h, _ = Insert(intOrd)(1)(h)
	h, _ = Insert(intOrd)(4)(h)
	assert.Equal(t, O.Some(1), FindMin(h))
	assert.Equal(t, 5, Size(h))
}

// ---- ExtractMin ------------------------------------------------------------

func TestExtractMinEmpty(t *testing.T) {
	h := Empty[int]()
	h2, val := ExtractMin[int](intOrd)(h)
	assert.Equal(t, O.None[int](), val)
	assert.True(t, IsEmpty(h2))
}

func TestExtractMinSingle(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(7)(h)
	h, val := ExtractMin[int](intOrd)(h)
	assert.Equal(t, O.Some(7), val)
	assert.True(t, IsEmpty(h))
}

func TestExtractMinOrder(t *testing.T) {
	keys := []int{5, 3, 8, 1, 4, 7, 2, 6}
	h := Empty[int]()
	for _, k := range keys {
		h, _ = Insert(intOrd)(k)(h)
	}
	got := extractAll(t, h)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, got)
}

func TestExtractMinDecreasing(t *testing.T) {
	h := Empty[int]()
	for i := 10; i >= 1; i-- {
		h, _ = Insert(intOrd)(i)(h)
	}
	got := extractAll(t, h)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, got)
}

// ---- Merge -----------------------------------------------------------------

func TestMergeEmpty(t *testing.T) {
	h1 := Empty[int]()
	h2 := Empty[int]()
	m := Merge(intOrd)(h1)(h2)
	assert.True(t, IsEmpty(m))
}

func TestMergeLeftEmpty(t *testing.T) {
	h1 := Empty[int]()
	h2 := Empty[int]()
	h2, _ = Insert(intOrd)(5)(h2)
	m := Merge(intOrd)(h1)(h2)
	assert.Equal(t, O.Some(5), FindMin(m))
	assert.Equal(t, 1, Size(m))
}

func TestMergeRightEmpty(t *testing.T) {
	h1 := Empty[int]()
	h1, _ = Insert(intOrd)(3)(h1)
	h2 := Empty[int]()
	m := Merge(intOrd)(h1)(h2)
	assert.Equal(t, O.Some(3), FindMin(m))
}

func TestMergeBothNonEmpty(t *testing.T) {
	h1 := Empty[int]()
	for _, k := range []int{4, 2, 6} {
		h1, _ = Insert(intOrd)(k)(h1)
	}
	h2 := Empty[int]()
	for _, k := range []int{5, 1, 3} {
		h2, _ = Insert(intOrd)(k)(h2)
	}
	m := Merge(intOrd)(h1)(h2)
	assert.Equal(t, 6, Size(m))
	got := extractAll(t, m)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

// ---- DecreaseKey -----------------------------------------------------------

func TestDecreaseKeyUpdatesMin(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(10)(h)
	var handle Handle[int]
	h, handle = Insert(intOrd)(20)(h)
	assert.Equal(t, O.Some(10), FindMin(h))

	h = DecreaseKey(intOrd)(5)(handle)(h)
	assert.Equal(t, O.Some(5), FindMin(h))
}

func TestDecreaseKeyThenExtract(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(10)(h)
	var h1, h2 Handle[int]
	h, h1 = Insert(intOrd)(8)(h)
	h, h2 = Insert(intOrd)(6)(h)
	_ = h2

	h = DecreaseKey(intOrd)(1)(h1)(h)

	got := extractAll(t, h)
	assert.Equal(t, []int{1, 6, 10}, got)
}

func TestDecreaseKeyPanicsOnIncrease(t *testing.T) {
	h := Empty[int]()
	var handle Handle[int]
	h, handle = Insert(intOrd)(5)(h)
	assert.Panics(t, func() {
		DecreaseKey(intOrd)(10)(handle)(h)
	})
}

// ---- Delete ----------------------------------------------------------------

func TestDeleteMin(t *testing.T) {
	h := Empty[int]()
	var minHandle Handle[int]
	h, minHandle = Insert(intOrd)(1)(h)
	h, _ = Insert(intOrd)(2)(h)
	h, _ = Insert(intOrd)(3)(h)

	h = Delete(intOrd)(minHandle)(h)
	assert.Equal(t, 2, Size(h))
	assert.Equal(t, O.Some(2), FindMin(h))
}

func TestDeleteMiddle(t *testing.T) {
	h := Empty[int]()
	h, _ = Insert(intOrd)(1)(h)
	var mid Handle[int]
	h, mid = Insert(intOrd)(5)(h)
	h, _ = Insert(intOrd)(3)(h)

	h = Delete(intOrd)(mid)(h)
	got := extractAll(t, h)
	assert.Equal(t, []int{1, 3}, got)
}

func TestDeleteAll(t *testing.T) {
	keys := []int{9, 3, 7, 1}
	handles := make([]Handle[int], len(keys))
	h := Empty[int]()
	for i, k := range keys {
		h, handles[i] = Insert(intOrd)(k)(h)
	}
	for _, handle := range handles {
		h = Delete(intOrd)(handle)(h)
	}
	assert.True(t, IsEmpty(h))
}

// ---- Custom Ord (struct keys) ----------------------------------------------

type Priority struct {
	Value int
	Name  string
}

func TestCustomOrd(t *testing.T) {
	byValue := ord.MakeOrd(
		func(a, b Priority) int {
			if a.Value < b.Value {
				return -1
			}
			if a.Value > b.Value {
				return 1
			}
			return 0
		},
		func(a, b Priority) bool { return a.Value == b.Value },
	)

	h := Empty[Priority]()
	h, _ = Insert(byValue)(Priority{3, "c"})(h)
	h, _ = Insert(byValue)(Priority{1, "a"})(h)
	h, _ = Insert(byValue)(Priority{2, "b"})(h)

	h, v := ExtractMin[Priority](byValue)(h)
	assert.Equal(t, O.Some(Priority{1, "a"}), v)

	_, v = ExtractMin[Priority](byValue)(h)
	assert.Equal(t, O.Some(Priority{2, "b"}), v)
}

// ---- Size invariants -------------------------------------------------------

func TestSizeAfterOperations(t *testing.T) {
	h := Empty[int]()
	assert.Equal(t, 0, Size(h))

	h, _ = Insert(intOrd)(1)(h)
	assert.Equal(t, 1, Size(h))

	h, _ = Insert(intOrd)(2)(h)
	assert.Equal(t, 2, Size(h))

	h, _ = ExtractMin[int](intOrd)(h)
	assert.Equal(t, 1, Size(h))

	h, _ = ExtractMin[int](intOrd)(h)
	assert.Equal(t, 0, Size(h))
	assert.True(t, IsEmpty(h))
}
