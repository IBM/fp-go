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

package readerioresult

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// collectSeq drains a SeqResult into a slice, stopping on the first error.
func collectSeq[A any](seq SeqResult[A]) ([]A, error) {
	var values []A
	for r := range seq {
		v, err := result.Unwrap(r)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}
	return values, nil
}

// collectSeqRaw drains a SeqResult into a slice of Result without any processing.
func collectSeqRaw[A any](seq SeqResult[A]) []Result[A] {
	var rs []Result[A]
	for r := range seq {
		rs = append(rs, r)
	}
	return rs
}

// counter returns a step function that emits 0, 1, …, n-1 from an integer seed.
func counter(limit int) func(int) ReaderIOResult[Option[Pair[int, int]]] {
	return func(n int) ReaderIOResult[Option[Pair[int, int]]] {
		if n >= limit {
			return Of(option.None[Pair[int, int]]())
		}
		return Of(option.Some(pair.MakePair(n+1, n)))
	}
}

func TestUnfold_EmitsExpectedValues(t *testing.T) {
	seq := Unfold(counter(5), 0)
	values, err := collectSeq(seq(t.Context()))
	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, values)
}

func TestUnfold_EmptySequenceWhenSeedTerminatesImmediately(t *testing.T) {
	seq := Unfold(counter(0), 0)
	values, err := collectSeq(seq(t.Context()))
	require.NoError(t, err)
	assert.Empty(t, values)
}

func TestUnfold_SingleElement(t *testing.T) {
	seq := Unfold(counter(1), 0)
	values, err := collectSeq(seq(t.Context()))
	require.NoError(t, err)
	assert.Equal(t, []int{0}, values)
}

func TestUnfold_PropagatesStepError(t *testing.T) {
	sentinel := errors.New("step failed")
	step := func(n int) ReaderIOResult[Option[Pair[int, int]]] {
		if n == 3 {
			return Left[Option[Pair[int, int]]](sentinel)
		}
		return Of(option.Some(pair.MakePair(n+1, n)))
	}
	seq := Unfold(step, 0)
	rs := collectSeqRaw(seq(t.Context()))
	// Expect Right(0), Right(1), Right(2), Left(sentinel)
	require.Len(t, rs, 4)
	assert.True(t, result.IsRight(rs[0]))
	assert.True(t, result.IsRight(rs[1]))
	assert.True(t, result.IsRight(rs[2]))
	assert.True(t, result.IsLeft(rs[3]))
	_, err := result.Unwrap(rs[3])
	assert.ErrorIs(t, err, sentinel)
}

func TestUnfold_ErrorOnFirstStep(t *testing.T) {
	sentinel := errors.New("immediate error")
	step := func(_ int) ReaderIOResult[Option[Pair[int, int]]] {
		return Left[Option[Pair[int, int]]](sentinel)
	}
	seq := Unfold(step, 0)
	rs := collectSeqRaw(seq(t.Context()))
	require.Len(t, rs, 1)
	assert.True(t, result.IsLeft(rs[0]))
	_, err := result.Unwrap(rs[0])
	assert.ErrorIs(t, err, sentinel)
}

func TestUnfold_CancelledContextBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // cancel before any iteration

	seq := Unfold(counter(5), 0)
	rs := collectSeqRaw(seq(ctx))
	require.Len(t, rs, 1)
	assert.True(t, result.IsLeft(rs[0]))
	_, err := result.Unwrap(rs[0])
	assert.ErrorIs(t, err, context.Canceled)
}

func TestUnfold_CancelledContextMidSequence(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	step := func(n int) ReaderIOResult[Option[Pair[int, int]]] {
		if n == 3 {
			cancel() // cancel partway through
		}
		return Of(option.Some(pair.MakePair(n+1, n)))
	}

	seq := Unfold(step, 0)
	rs := collectSeqRaw(seq(ctx))

	// Values 0, 1, 2 succeed; then cancellation is detected before the next iteration
	require.GreaterOrEqual(t, len(rs), 3)
	for i := 0; i < 3; i++ {
		assert.True(t, result.IsRight(rs[i]), "result[%d] should be Right", i)
	}
	last := rs[len(rs)-1]
	assert.True(t, result.IsLeft(last))
	_, err := result.Unwrap(last)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestUnfold_EarlyTerminationByConsumer(t *testing.T) {
	calls := 0
	step := func(n int) ReaderIOResult[Option[Pair[int, int]]] {
		calls++
		return Of(option.Some(pair.MakePair(n+1, n)))
	}

	seq := Unfold(step, 0)
	collected := 0
	for r := range seq(t.Context()) {
		assert.True(t, result.IsRight(r))
		collected++
		if collected == 3 {
			break // consumer stops after 3 items
		}
	}

	assert.Equal(t, 3, collected)
	// Step should have been called at most 3 times
	assert.LessOrEqual(t, calls, 3)
}

func TestUnfold_StringSeed(t *testing.T) {
	// Unfold characters from a string seed
	step := func(s string) ReaderIOResult[Option[Pair[string, byte]]] {
		if len(s) == 0 {
			return Of(option.None[Pair[string, byte]]())
		}
		return Of(option.Some(pair.MakePair(s[1:], s[0])))
	}

	seq := Unfold(step, "abc")
	rs := collectSeqRaw(seq(t.Context()))
	require.Len(t, rs, 3)
	v0, _ := result.Unwrap(rs[0])
	v1, _ := result.Unwrap(rs[1])
	v2, _ := result.Unwrap(rs[2])
	assert.Equal(t, byte('a'), v0)
	assert.Equal(t, byte('b'), v1)
	assert.Equal(t, byte('c'), v2)
}

func TestUnfold_ContextValueVisible(t *testing.T) {
	type ctxKey string
	const key ctxKey = "id"

	step := func(n int) ReaderIOResult[Option[Pair[int, string]]] {
		return func(ctx context.Context) IOEither[Option[Pair[int, string]]] {
			return func() Either[Option[Pair[int, string]]] {
				id, _ := ctx.Value(key).(string)
				if n >= 2 {
					return result.Of(option.None[Pair[int, string]]())
				}
				return result.Of(option.Some(pair.MakePair(n+1, id)))
			}
		}
	}

	ctx := context.WithValue(t.Context(), key, "tenant-42")
	seq := Unfold(step, 0)
	values, err := collectSeq(seq(ctx))
	require.NoError(t, err)
	assert.Equal(t, []string{"tenant-42", "tenant-42"}, values)
}
