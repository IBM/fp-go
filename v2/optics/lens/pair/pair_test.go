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

package pair

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/common"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// TestHead_Get verifies that Head correctly extracts the head element.
func TestHead_Get(t *testing.T) {
	t.Run("string head", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		assert.Equal(t, "hello", Head[string, int]().Get(p))
	})

	t.Run("int head", func(t *testing.T) {
		p := pair.MakePair(7, true)
		assert.Equal(t, 7, Head[int, bool]().Get(p))
	})

	t.Run("zero value head", func(t *testing.T) {
		p := pair.MakePair("", 0)
		assert.Equal(t, "", Head[string, int]().Get(p))
	})
}

// TestHead_Set verifies that Head correctly replaces the head while
// preserving the tail and leaving the original pair unchanged.
func TestHead_Set(t *testing.T) {
	t.Run("replaces head", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		updated := Head[string, int]().Set("world")(p)
		assert.Equal(t, "world", pair.Head(updated))
		assert.Equal(t, 42, pair.Tail(updated))
	})

	t.Run("preserves original", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		_ = Head[string, int]().Set("world")(p)
		assert.Equal(t, "hello", pair.Head(p))
	})

	t.Run("tail unchanged", func(t *testing.T) {
		p := pair.MakePair("hello", 99)
		updated := Head[string, int]().Set("changed")(p)
		assert.Equal(t, 99, pair.Tail(updated))
	})
}

// TestTail_Get verifies that Tail correctly extracts the tail element.
func TestTail_Get(t *testing.T) {
	t.Run("int tail", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		assert.Equal(t, 42, Tail[string, int]().Get(p))
	})

	t.Run("bool tail", func(t *testing.T) {
		p := pair.MakePair(7, true)
		assert.Equal(t, true, Tail[int, bool]().Get(p))
	})

	t.Run("zero value tail", func(t *testing.T) {
		p := pair.MakePair("key", 0)
		assert.Equal(t, 0, Tail[string, int]().Get(p))
	})
}

// TestTail_Set verifies that Tail correctly replaces the tail while
// preserving the head and leaving the original pair unchanged.
func TestTail_Set(t *testing.T) {
	t.Run("replaces tail", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		updated := Tail[string, int]().Set(100)(p)
		assert.Equal(t, "hello", pair.Head(updated))
		assert.Equal(t, 100, pair.Tail(updated))
	})

	t.Run("preserves original", func(t *testing.T) {
		p := pair.MakePair("hello", 42)
		_ = Tail[string, int]().Set(100)(p)
		assert.Equal(t, 42, pair.Tail(p))
	})

	t.Run("head unchanged", func(t *testing.T) {
		p := pair.MakePair("original", 1)
		updated := Tail[string, int]().Set(99)(p)
		assert.Equal(t, "original", pair.Head(updated))
	})
}

// TestHead_LensLaws verifies that Head satisfies all three lens laws:
//
//	get(set(a)(s)) = a
//	set(get(s))(s) = s
//	set(a)(set(a)(s)) = set(a)(s)
func TestHead_LensLaws(t *testing.T) {
	eqString := EQT.Eq[string]()
	eqPair := EQT.Eq[Pair[string, int]]()

	laws := LT.AssertLaws(t, eqString, eqPair)(Head[string, int]())

	assert.True(t, laws(pair.MakePair("hello", 42), "world"))
	assert.True(t, laws(pair.MakePair("", 0), "non-empty"))
	assert.True(t, laws(pair.MakePair("same", 1), "same"))
}

// TestTail_LensLaws verifies that Tail satisfies all three lens laws:
//
//	get(set(a)(s)) = a
//	set(get(s))(s) = s
//	set(a)(set(a)(s)) = set(a)(s)
func TestTail_LensLaws(t *testing.T) {
	eqInt := EQT.Eq[int]()
	eqPair := EQT.Eq[Pair[string, int]]()

	laws := LT.AssertLaws(t, eqInt, eqPair)(Tail[string, int]())

	assert.True(t, laws(pair.MakePair("hello", 42), 100))
	assert.True(t, laws(pair.MakePair("", 0), -1))
	assert.True(t, laws(pair.MakePair("key", 7), 7))
}

// TestHead_Modify verifies modifying the head through the lens.
func TestHead_Modify(t *testing.T) {
	p := pair.MakePair("hello", 42)
	toUpper := F.Pipe1(
		Head[string, int](),
		L.LensModify[Pair[string, int]](func(s string) string { return s + "!" }),
	)

	result := toUpper(p)
	assert.Equal(t, "hello!", pair.Head(result))
	assert.Equal(t, 42, pair.Tail(result))
	assert.Equal(t, "hello", pair.Head(p)) // original unchanged
}

// TestTail_Modify verifies modifying the tail through the lens.
func TestTail_Modify(t *testing.T) {
	p := pair.MakePair("hello", 42)
	double := F.Pipe1(
		Tail[string, int](),
		L.LensModify[Pair[string, int]](func(n int) int { return n * 2 }),
	)

	result := double(p)
	assert.Equal(t, "hello", pair.Head(result))
	assert.Equal(t, 84, pair.Tail(result))
	assert.Equal(t, 42, pair.Tail(p)) // original unchanged
}

// TestHead_Compose verifies that Head can be composed with another lens.
func TestHead_Compose(t *testing.T) {
	// Pair[Pair[string, bool], int] — compose Head on outer then Head on inner
	outerHead := Head[Pair[string, bool], int]()
	innerHead := Head[string, bool]()

	// Composed lens: outer -> head (Pair[string, bool]) -> head (string)
	composed := F.Pipe1(outerHead, L.LensComposeLens[Pair[Pair[string, bool], int]](innerHead))

	p := pair.MakePair(pair.MakePair("deep", true), 99)
	assert.Equal(t, "deep", composed.Get(p))

	updated := composed.Set("changed")(p)
	assert.Equal(t, "changed", pair.Head(pair.Head(updated)))
	assert.Equal(t, true, pair.Tail(pair.Head(updated)))
	assert.Equal(t, 99, pair.Tail(updated))
}

// TestTail_Compose verifies that Tail can be composed with another lens.
func TestTail_Compose(t *testing.T) {
	// Pair[string, Pair[int, bool]] — compose Tail on outer then Head on inner
	outerTail := Tail[string, Pair[int, bool]]()
	innerHead := Head[int, bool]()

	// Composed lens: outer -> tail (Pair[int, bool]) -> head (int)
	composed := F.Pipe1(outerTail, L.LensComposeLens[Pair[string, Pair[int, bool]]](innerHead))

	p := pair.MakePair("outer", pair.MakePair(7, false))
	assert.Equal(t, 7, composed.Get(p))

	updated := composed.Set(42)(p)
	assert.Equal(t, "outer", pair.Head(updated))
	assert.Equal(t, 42, pair.Head(pair.Tail(updated)))
	assert.Equal(t, false, pair.Tail(pair.Tail(updated)))
}

// TestHead_Name verifies that the lens has a meaningful name.
func TestHead_Name(t *testing.T) {
	l := Head[string, int]()
	assert.Equal(t, "pair.Head", l.String())
}

// TestTail_Name verifies that the lens has a meaningful name.
func TestTail_Name(t *testing.T) {
	l := Tail[string, int]()
	assert.Equal(t, "pair.Tail", l.String())
}
