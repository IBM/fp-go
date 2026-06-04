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

package io

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// counter returns a step function that emits 0, 1, …, limit-1 from an integer seed.
func counter(limit int) Kleisli[int, Option[Pair[int, int]]] {
	return func(n int) IO[Option[Pair[int, int]]] {
		if n >= limit {
			return Of(option.None[Pair[int, int]]())
		}
		return Of(option.Some(pair.MakePair(n+1, n)))
	}
}

// collect drains a Seq[A] into a slice.
func collect[A any](seq Seq[A]) []A {
	var out []A
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func TestUnfold_EmitsExpectedValues(t *testing.T) {
	values := collect(Unfold(counter(5))(0))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, values)
}

func TestUnfold_EmptySequenceWhenSeedTerminatesImmediately(t *testing.T) {
	values := collect(Unfold(counter(0))(0))
	assert.Empty(t, values)
}

func TestUnfold_SingleElement(t *testing.T) {
	values := collect(Unfold(counter(1))(0))
	assert.Equal(t, []int{0}, values)
}

func TestUnfold_StringSeed(t *testing.T) {
	step := func(s string) IO[Option[Pair[string, byte]]] {
		if len(s) == 0 {
			return Of(option.None[Pair[string, byte]]())
		}
		return Of(option.Some(pair.MakePair(s[1:], s[0])))
	}

	values := collect(Unfold(step)("abc"))
	require.Len(t, values, 3)
	assert.Equal(t, byte('a'), values[0])
	assert.Equal(t, byte('b'), values[1])
	assert.Equal(t, byte('c'), values[2])
}

func TestUnfold_EarlyTerminationByConsumer(t *testing.T) {
	calls := 0
	step := func(n int) IO[Option[Pair[int, int]]] {
		calls++
		return Of(option.Some(pair.MakePair(n+1, n)))
	}

	collected := 0
	for range Unfold(step)(0) {
		collected++
		if collected == 3 {
			break
		}
	}

	assert.Equal(t, 3, collected)
	assert.LessOrEqual(t, calls, 3)
}

func TestUnfold_LargeSequence(t *testing.T) {
	const n = 1000
	values := collect(Unfold(counter(n))(0))
	require.Len(t, values, n)
	assert.Equal(t, 0, values[0])
	assert.Equal(t, n-1, values[n-1])
}
