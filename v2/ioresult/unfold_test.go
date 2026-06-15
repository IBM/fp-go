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

package ioresult

import (
	"errors"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// collectSeqRaw drains a SeqResult into a slice of Result values.
func collectSeqRaw[A any](seq SeqResult[A]) []Result[A] {
	var rs []Result[A]
	for r := range seq {
		rs = append(rs, r)
	}
	return rs
}

// collectSeq drains a SeqResult into a slice of values, returning the first error encountered.
func collectSeq[A any](seq SeqResult[A]) ([]A, error) {
	var values []A
	for r := range seq {
		v, err := R.Unwrap(r)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}
	return values, nil
}

// countDownStep is a reusable step function emitting n, n-1, …, 1 and stopping at 0.
func countDownStep(n int) IOResult[O.Option[P.Pair[int, int]]] {
	if n == 0 {
		return Of(O.None[P.Pair[int, int]]())
	}
	// Head = next seed, Tail = emitted value
	return Of(O.Some(P.MakePair(n-1, n)))
}

func TestUnfold_EmptySequenceWhenSeedTerminatesImmediately(t *testing.T) {
	seq := Unfold(countDownStep)(0)
	values, err := collectSeq(seq)
	require.NoError(t, err)
	assert.Empty(t, values)
}

func TestUnfold_FiniteSequenceCountdown(t *testing.T) {
	seq := Unfold(countDownStep)(3)
	values, err := collectSeq(seq)
	require.NoError(t, err)
	assert.Equal(t, []int{3, 2, 1}, values)
}

func TestUnfold_SingleElement(t *testing.T) {
	seq := Unfold(countDownStep)(1)
	values, err := collectSeq(seq)
	require.NoError(t, err)
	assert.Equal(t, []int{1}, values)
}

func TestUnfold_ErrorOnFirstStep(t *testing.T) {
	sentinel := errors.New("immediate failure")
	step := func(_ int) IOResult[O.Option[P.Pair[int, int]]] {
		return Left[O.Option[P.Pair[int, int]]](sentinel)
	}
	rs := collectSeqRaw(Unfold(step)(0))
	require.Len(t, rs, 1)
	assert.True(t, R.IsLeft(rs[0]))
	_, err := R.Unwrap(rs[0])
	assert.ErrorIs(t, err, sentinel)
}

func TestUnfold_ErrorMidSequence(t *testing.T) {
	sentinel := errors.New("mid-sequence failure")
	// Emits Right(3), Right(2), Right(1), then Left(sentinel) when seed reaches 0
	step := func(n int) IOResult[O.Option[P.Pair[int, int]]] {
		if n == 0 {
			return Left[O.Option[P.Pair[int, int]]](sentinel)
		}
		return Of(O.Some(P.MakePair(n-1, n)))
	}
	rs := collectSeqRaw(Unfold(step)(3))
	require.Len(t, rs, 4)
	assert.True(t, R.IsRight(rs[0]))
	assert.True(t, R.IsRight(rs[1]))
	assert.True(t, R.IsRight(rs[2]))
	assert.True(t, R.IsLeft(rs[3]))
	_, err := R.Unwrap(rs[3])
	assert.ErrorIs(t, err, sentinel)
}

func TestUnfold_EarlyTerminationByConsumer(t *testing.T) {
	calls := 0
	step := func(n int) IOResult[O.Option[P.Pair[int, int]]] {
		calls++
		return Of(O.Some(P.MakePair(n+1, n)))
	}
	collected := 0
	for r := range Unfold(step)(0) {
		assert.True(t, R.IsRight(r))
		collected++
		if collected == 3 {
			break
		}
	}
	assert.Equal(t, 3, collected)
	assert.LessOrEqual(t, calls, 3)
}

func TestUnfold_StringUnfolding(t *testing.T) {
	// Unfold bytes from a string seed
	step := func(s string) IOResult[O.Option[P.Pair[string, byte]]] {
		if len(s) == 0 {
			return Of(O.None[P.Pair[string, byte]]())
		}
		return Of(O.Some(P.MakePair(s[1:], s[0])))
	}
	rs := collectSeqRaw(Unfold(step)("abc"))
	require.Len(t, rs, 3)
	v0, _ := R.Unwrap(rs[0])
	v1, _ := R.Unwrap(rs[1])
	v2, _ := R.Unwrap(rs[2])
	assert.Equal(t, byte('a'), v0)
	assert.Equal(t, byte('b'), v1)
	assert.Equal(t, byte('c'), v2)
}
