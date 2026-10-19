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
	"sync/atomic"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// concurrencyProbe records the maximum number of computations that run at the same time.
type concurrencyProbe struct {
	running atomic.Int32
	peak    atomic.Int32
}

func (p *concurrencyProbe) track(d time.Duration) Kleisli[int, int] {
	return func(n int) ReaderIOResult[int] {
		return func(ctx context.Context) IOResult[int] {
			return func() Result[int] {
				cur := p.running.Add(1)
				defer p.running.Add(-1)
				for {
					old := p.peak.Load()
					if cur <= old || p.peak.CompareAndSwap(old, cur) {
						break
					}
				}
				select {
				case <-time.After(d):
					return result.Of(n * 2)
				case <-ctx.Done():
					return result.Left[int](context.Cause(ctx))
				}
			}
		}
	}
}

func TestTraverseArrayParN_Success(t *testing.T) {
	t.Run("preserves the order of the input", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		res := TraverseArrayParN(3, F.Flow2(N.Mul(2), Of[int]))(input)(t.Context())()
		assert.Equal(t, result.Of([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}), res)
	})

	t.Run("bounds the number of concurrent computations", func(t *testing.T) {
		probe := &concurrencyProbe{}
		input := make([]int, 20)
		res := TraverseArrayParN(3, probe.track(5*time.Millisecond))(input)(t.Context())()
		assert.True(t, result.IsRight(res))
		assert.LessOrEqual(t, probe.peak.Load(), int32(3))
		assert.Greater(t, probe.peak.Load(), int32(1))
	})

	t.Run("matches TraverseArrayPar", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5}
		f := F.Flow2(N.Add(1), Of[int])
		assert.Equal(t,
			TraverseArrayPar(f)(input)(t.Context())(),
			TraverseArrayParN(2, f)(input)(t.Context())(),
		)
	})
}

func TestTraverseArrayParN_Failure(t *testing.T) {
	t.Run("returns the error and cancels the other computations", func(t *testing.T) {
		boom := errors.New("boom")
		var started atomic.Int32
		f := func(n int) ReaderIOResult[int] {
			return func(ctx context.Context) IOResult[int] {
				return func() Result[int] {
					started.Add(1)
					if n == 0 {
						return result.Left[int](boom)
					}
					select {
					case <-time.After(time.Second):
						return result.Of(n)
					case <-ctx.Done():
						return result.Left[int](context.Cause(ctx))
					}
				}
			}
		}
		input := []int{1, 0, 2, 3, 4, 5, 6, 7, 8, 9}

		start := time.Now()
		res := TraverseArrayParN(2, f)(input)(t.Context())()

		assert.Equal(t, result.Left[[]int](boom), res)
		assert.Less(t, time.Since(start), 500*time.Millisecond)
		assert.Less(t, started.Load(), int32(len(input)))
	})

	t.Run("fails when the outer context is cancelled", func(t *testing.T) {
		probe := &concurrencyProbe{}
		ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
		defer cancel()

		res := TraverseArrayParN(2, probe.track(time.Second))(make([]int, 10))(ctx)()

		assert.Equal(t, result.Left[[]int](context.DeadlineExceeded), res)
	})
}

func TestTraverseArrayParN_EdgeCases(t *testing.T) {
	double := F.Flow2(N.Mul(2), Of[int])

	t.Run("empty input", func(t *testing.T) {
		res := TraverseArrayParN(3, double)(nil)(t.Context())()
		assert.Equal(t, result.Of([]int{}), res)
	})

	t.Run("non-positive concurrency means unbounded", func(t *testing.T) {
		probe := &concurrencyProbe{}
		input := make([]int, 8)
		res := TraverseArrayParN(0, probe.track(20*time.Millisecond))(input)(t.Context())()
		assert.True(t, result.IsRight(res))
		assert.Equal(t, int32(8), probe.peak.Load())
	})

	t.Run("concurrency of one runs sequentially", func(t *testing.T) {
		probe := &concurrencyProbe{}
		res := TraverseArrayParN(1, probe.track(time.Millisecond))(make([]int, 5))(t.Context())()
		assert.True(t, result.IsRight(res))
		assert.Equal(t, int32(1), probe.peak.Load())
	})

	t.Run("already cancelled context does not run anything", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		var calls atomic.Int32
		f := func(n int) ReaderIOResult[int] {
			calls.Add(1)
			return Of(n)
		}
		res := TraverseArrayParN(2, f)([]int{1, 2, 3})(ctx)()
		assert.Equal(t, result.Left[[]int](context.Canceled), res)
		assert.Equal(t, int32(0), calls.Load())
	})
}

func BenchmarkTraverseArrayParN(b *testing.B) {
	input := make([]int, 1000)
	f := TraverseArrayParN(8, F.Flow2(N.Mul(2), Of[int]))
	ctx := b.Context()
	b.ReportAllocs()
	for b.Loop() {
		f(input)(ctx)()
	}
}
