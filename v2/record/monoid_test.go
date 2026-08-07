//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package record

import (
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestUnionMonoid(t *testing.T) {
	m := UnionMonoid[string](S.Semigroup)

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b1b2",
		"c": "c1c2",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestUnionFirstMonoid(t *testing.T) {
	m := UnionFirstMonoid[string, string]()

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestUnionLastMonoid(t *testing.T) {
	m := UnionLastMonoid[string, string]()

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}

func ExampleUnionMonoid() {
	// UnionMonoid keeps all keys from both maps and combines values for shared
	// keys using the provided semigroup (here: string concatenation).
	m := UnionMonoid[string](S.Semigroup)

	x := map[string]string{"a": "hello", "b": "world", "c": "foo"}
	y := map[string]string{"b": "!", "c": "bar", "d": "extra"}

	result := m.Concat(x, y)

	// Print keys in deterministic order
	fmt.Println(result["a"])
	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(result["d"])
	fmt.Println(len(result))
	// Output:
	// hello
	// world!
	// foobar
	// extra
	// 4
}

func ExampleUnionLastMonoid() {
	// UnionLastMonoid keeps all keys from both maps; when a key is shared
	// the value from the right (second) map wins.
	m := UnionLastMonoid[string, string]()

	x := map[string]string{"a": "left-only", "b": "left-b", "c": "left-c"}
	y := map[string]string{"b": "right-b", "c": "right-c", "d": "right-only"}

	result := m.Concat(x, y)

	fmt.Println(result["a"])
	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(result["d"])
	fmt.Println(len(result))
	// Output:
	// left-only
	// right-b
	// right-c
	// right-only
	// 4
}

func ExampleUnionFirstMonoid() {
	// UnionFirstMonoid keeps all keys from both maps; when a key is shared
	// the value from the left (first) map wins.
	m := UnionFirstMonoid[string, string]()

	x := map[string]string{"a": "left-only", "b": "left-b", "c": "left-c"}
	y := map[string]string{"b": "right-b", "c": "right-c", "d": "right-only"}

	result := m.Concat(x, y)

	fmt.Println(result["a"])
	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(result["d"])
	fmt.Println(len(result))
	// Output:
	// left-only
	// left-b
	// left-c
	// right-only
	// 4
}

func ExampleMergeMonoid() {
	// MergeMonoid is an alias for UnionLastMonoid: all keys are kept and the
	// right (second) map wins for shared keys.
	m := MergeMonoid[string, string]()

	x := map[string]string{"a": "left-only", "b": "left-b", "c": "left-c"}
	y := map[string]string{"b": "right-b", "c": "right-c", "d": "right-only"}

	result := m.Concat(x, y)

	fmt.Println(result["a"])
	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(result["d"])
	fmt.Println(len(result))
	// Output:
	// left-only
	// right-b
	// right-c
	// right-only
	// 4
}

func TestIntersectionMonoid(t *testing.T) {
	m := IntersectionMonoid[string](S.Semigroup)

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	// only shared keys survive; values are concatenated by the semigroup
	res := map[string]string{
		"b": "b1b2",
		"c": "c1c2",
	}

	// intersection with empty is always empty (absorbing element)
	assert.Empty(t, m.Concat(x, m.Empty()))
	assert.Empty(t, m.Concat(m.Empty(), x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestIntersectionLastMonoid(t *testing.T) {
	m := IntersectionLastMonoid[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	// only shared keys survive; last (right) value wins
	res := map[string]string{
		"b": "b2",
		"c": "c2",
	}

	// intersection with empty is always empty (absorbing element)
	assert.Empty(t, m.Concat(x, m.Empty()))
	assert.Empty(t, m.Concat(m.Empty(), x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestIntersectionFirstMonoid(t *testing.T) {
	m := IntersectionFirstMonoid[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	// only shared keys survive; first (left) value wins
	res := map[string]string{
		"b": "b1",
		"c": "c1",
	}

	// intersection with empty is always empty (absorbing element)
	assert.Empty(t, m.Concat(x, m.Empty()))
	assert.Empty(t, m.Concat(m.Empty(), x))

	assert.Equal(t, res, m.Concat(x, y))
}

func TestIntersectionMonoidNoOverlap(t *testing.T) {
	m := IntersectionMonoid[string](N.SemigroupSum[int]())

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"c": 3, "d": 4}

	// no shared keys → empty result
	assert.Empty(t, m.Concat(map1, map2))
}

func ExampleIntersectionMonoid() {
	// IntersectionMonoid keeps only keys present in both maps and combines
	// values using the provided semigroup (here: string concatenation).
	m := IntersectionMonoid[string](S.Semigroup)

	x := map[string]string{"a": "hello", "b": "world", "c": "foo"}
	y := map[string]string{"b": "!", "c": "bar", "d": "extra"}

	result := m.Concat(x, y)

	// Print keys in deterministic order
	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(len(result))
	// Output:
	// world!
	// foobar
	// 2
}

func ExampleIntersectionLastMonoid() {
	// IntersectionLastMonoid keeps only keys present in both maps;
	// when a key is shared the value from the right (second) map wins.
	m := IntersectionLastMonoid[string, string]()

	x := map[string]string{"a": "left-only", "b": "left-b", "c": "left-c"}
	y := map[string]string{"b": "right-b", "c": "right-c", "d": "right-only"}

	result := m.Concat(x, y)

	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(len(result))
	// Output:
	// right-b
	// right-c
	// 2
}

func ExampleIntersectionFirstMonoid() {
	// IntersectionFirstMonoid keeps only keys present in both maps;
	// when a key is shared the value from the left (first) map wins.
	m := IntersectionFirstMonoid[string, string]()

	x := map[string]string{"a": "left-only", "b": "left-b", "c": "left-c"}
	y := map[string]string{"b": "right-b", "c": "right-c", "d": "right-only"}

	result := m.Concat(x, y)

	fmt.Println(result["b"])
	fmt.Println(result["c"])
	fmt.Println(len(result))
	// Output:
	// left-b
	// left-c
	// 2
}
