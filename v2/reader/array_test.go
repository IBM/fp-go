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

package reader

import (
	"context"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
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
