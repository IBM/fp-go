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

package reader

import (
	"context"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	n := 10

	readers := A.MakeBy(n, Of[context.Context, int])
	exp := A.MakeBy(n, F.Identity[int])

	g := F.Pipe1(
		readers,
		SequenceArray[context.Context, int],
	)

	assert.Equal(t, exp, g(context.Background()))
}

func TestTraverseArray(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 10}

	multiply := func(n int) Reader[Config, int] {
		return Asks(func(c Config) int { return n * c.Multiplier })
	}

	transform := TraverseArray(multiply)
	r := transform([]int{1, 2, 3})
	result := r(config)

	assert.Equal(t, []int{10, 20, 30}, result)
}

func TestTraverseArrayWithIndex(t *testing.T) {
	type Config struct{ Prefix string }
	config := Config{Prefix: "item"}

	addIndexPrefix := func(i int, s string) Reader[Config, string] {
		return Asks(func(c Config) string {
			// Simple string formatting
			idx := string(rune('0' + i))
			return c.Prefix + "[" + idx + "]:" + s
		})
	}

	transform := TraverseArrayWithIndex(addIndexPrefix)
	r := transform([]string{"a", "b", "c"})
	result := r(config)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "item[0]:a", result[0])
	assert.Equal(t, "item[1]:b", result[1])
	assert.Equal(t, "item[2]:c", result[2])
}

func TestMonadTraverseArray(t *testing.T) {
	type Config struct{ Prefix string }
	config := Config{Prefix: "num"}

	numbers := []int{1, 2, 3}
	addPrefix := func(n int) Reader[Config, string] {
		return Asks(func(c Config) string {
			return c.Prefix + string(rune('0'+n))
		})
	}

	r := MonadTraverseArray(numbers, addPrefix)
	result := r(config)

	assert.Equal(t, 3, len(result))
	assert.Contains(t, result[0], "num")
}

func TestMonadReduceArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	readers := []Reader[Config, int]{
		Asks(func(c Config) int { return c.Base + 1 }),
		Asks(func(c Config) int { return c.Base + 2 }),
		Asks(func(c Config) int { return c.Base + 3 }),
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadReduceArray(readers, sum, 0)
	result := r(config)

	assert.Equal(t, 36, result) // 11 + 12 + 13
}

func TestReduceArray(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 5}

	product := func(acc, val int) int { return acc * val }
	reducer := ReduceArray[Config](product, 1)

	readers := []Reader[Config, int]{
		Asks(func(c Config) int { return c.Multiplier * 2 }),
		Asks(func(c Config) int { return c.Multiplier * 3 }),
	}

	r := reducer(readers)
	result := r(config)

	assert.Equal(t, 150, result) // 10 * 15
}

func TestMonadReduceArrayM(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	readers := []Reader[Config, int]{
		Asks(func(c Config) int { return c.Factor }),
		Asks(func(c Config) int { return c.Factor * 2 }),
		Asks(func(c Config) int { return c.Factor * 3 }),
	}

	intAddMonoid := N.MonoidSum[int]()

	r := MonadReduceArrayM(readers, intAddMonoid)
	result := r(config)

	assert.Equal(t, 30, result) // 5 + 10 + 15
}

func TestReduceArrayM(t *testing.T) {
	type Config struct{ Scale int }
	config := Config{Scale: 3}

	intMultMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)

	reducer := ReduceArrayM[Config](intMultMonoid)

	readers := []Reader[Config, int]{
		Asks(func(c Config) int { return c.Scale }),
		Asks(func(c Config) int { return c.Scale * 2 }),
	}

	r := reducer(readers)
	result := r(config)

	assert.Equal(t, 18, result) // 3 * 6
}

func TestMonadTraverseReduceArray(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 10}

	numbers := []int{1, 2, 3, 4}
	multiply := func(n int) Reader[Config, int] {
		return Asks(func(c Config) int { return n * c.Multiplier })
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadTraverseReduceArray(numbers, multiply, sum, 0)
	result := r(config)

	assert.Equal(t, 100, result) // 10 + 20 + 30 + 40
}

func TestTraverseReduceArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	addBase := func(n int) Reader[Config, int] {
		return Asks(func(c Config) int { return n + c.Base })
	}

	product := func(acc, val int) int { return acc * val }
	transformer := TraverseReduceArray(addBase, product, 1)

	r := transformer([]int{2, 3, 4})
	result := r(config)

	assert.Equal(t, 2184, result) // 12 * 13 * 14
}

func TestMonadTraverseReduceArrayM(t *testing.T) {
	type Config struct{ Offset int }
	config := Config{Offset: 100}

	numbers := []int{1, 2, 3}
	addOffset := func(n int) Reader[Config, int] {
		return Asks(func(c Config) int { return n + c.Offset })
	}

	intSumMonoid := N.MonoidSum[int]()

	r := MonadTraverseReduceArrayM(numbers, addOffset, intSumMonoid)
	result := r(config)

	assert.Equal(t, 306, result) // 101 + 102 + 103
}

func TestTraverseReduceArrayM(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	scale := func(n int) Reader[Config, int] {
		return Asks(func(c Config) int { return n * c.Factor })
	}

	intProdMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)

	transformer := TraverseReduceArrayM(scale, intProdMonoid)
	r := transformer([]int{2, 3, 4})
	result := r(config)

	assert.Equal(t, 3000, result) // 10 * 15 * 20
}
