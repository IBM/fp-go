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

package effect

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type combConfig struct {
	Factor int
}

func combFactor(c combConfig) int { return c.Factor }

var combCfg = combConfig{Factor: 3}

func combScale(n int) Effect[combConfig, int] {
	return F.Pipe1(Asks(combFactor), Map[combConfig](N.Mul(n)))
}

// combRun runs an effect with combCfg and the given context.
func combRun[A any](ctx context.Context, eff Effect[combConfig, A]) Result[A] {
	return eff(combCfg)(ctx)()
}

// combProbe records the peak number of effects running at the same time.
type combProbe struct {
	running atomic.Int32
	peak    atomic.Int32
}

func (p *combProbe) track(d time.Duration) Kleisli[combConfig, int, int] {
	return func(n int) Effect[combConfig, int] {
		return func(c combConfig) Thunk[int] {
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
						return result.Of(n * c.Factor)
					case <-ctx.Done():
						return result.Left[int](context.Cause(ctx))
					}
				}
			}
		}
	}
}

// recorder collects the names of executed steps in order.
type recorder struct {
	mu    sync.Mutex
	steps []string
}

func (r *recorder) step(name string) Effect[combConfig, Void] {
	return FromIO[combConfig](func() Void {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.steps = append(r.steps, name)
		return F.VOID
	})
}

func TestBracket_Success(t *testing.T) {
	t.Run("acquires, uses and releases with the outcome of use", func(t *testing.T) {
		rec := &recorder{}
		var released Result[int]
		eff := Bracket(
			rec.step("acquire"),
			func(Void) Effect[combConfig, int] { return combScale(2) },
			func(_ Void, res Result[int]) Effect[combConfig, Void] {
				released = res
				return rec.step("release")
			},
		)

		assert.Equal(t, result.Of(6), combRun(t.Context(), eff))
		assert.Equal(t, []string{"acquire", "release"}, rec.steps)
		assert.Equal(t, result.Of(6), released)
	})
}

func TestBracket_Failure(t *testing.T) {
	boom := errors.New("boom")

	t.Run("releases when use fails", func(t *testing.T) {
		rec := &recorder{}
		eff := Bracket(
			rec.step("acquire"),
			func(Void) Effect[combConfig, int] { return Fail[combConfig, int](boom) },
			func(Void, Result[int]) Effect[combConfig, Void] { return rec.step("release") },
		)

		assert.Equal(t, result.Left[int](boom), combRun(t.Context(), eff))
		assert.Equal(t, []string{"acquire", "release"}, rec.steps)
	})

	t.Run("does not use or release when acquire fails", func(t *testing.T) {
		rec := &recorder{}
		eff := Bracket(
			Fail[combConfig, Void](boom),
			func(Void) Effect[combConfig, int] {
				return F.Pipe1(rec.step("use"), Map[combConfig](F.Constant1[Void](1)))
			},
			func(Void, Result[int]) Effect[combConfig, Void] { return rec.step("release") },
		)

		assert.Equal(t, result.Left[int](boom), combRun(t.Context(), eff))
		assert.Empty(t, rec.steps)
	})
}

func TestWithResource_Success(t *testing.T) {
	rec := &recorder{}
	withRes := WithResource[int](rec.step("create"), func(Void) Effect[combConfig, Void] { return rec.step("release") })

	res := combRun(t.Context(), withRes(func(Void) Effect[combConfig, int] { return combScale(5) }))

	assert.Equal(t, result.Of(15), res)
	assert.Equal(t, []string{"create", "release"}, rec.steps)
}

func TestWithResource_Failure(t *testing.T) {
	boom := errors.New("boom")
	rec := &recorder{}
	withRes := WithResource[int](rec.step("create"), func(Void) Effect[combConfig, Void] { return rec.step("release") })

	res := combRun(t.Context(), withRes(func(Void) Effect[combConfig, int] { return Fail[combConfig, int](boom) }))

	assert.Equal(t, result.Left[int](boom), res)
	assert.Equal(t, []string{"create", "release"}, rec.steps)
}

func TestWithLock_Success(t *testing.T) {
	t.Run("enforces mutual exclusion", func(t *testing.T) {
		var mu sync.Mutex
		lock := FromIO[combConfig](func() context.CancelFunc {
			mu.Lock()
			return mu.Unlock
		})
		probe := &combProbe{}

		res := combRun(t.Context(), TraverseArrayPar(F.Flow2(
			probe.track(2*time.Millisecond),
			WithLock[combConfig, int](lock),
		))([]int{1, 2, 3, 4}))

		assert.Equal(t, result.Of([]int{3, 6, 9, 12}), res)
		assert.Equal(t, int32(1), probe.peak.Load())
	})
}

func TestWithLock_Failure(t *testing.T) {
	boom := errors.New("boom")
	var mu sync.Mutex
	lock := FromIO[combConfig](func() context.CancelFunc {
		mu.Lock()
		return mu.Unlock
	})

	res := combRun(t.Context(), WithLock[combConfig, int](lock)(Fail[combConfig, int](boom)))

	assert.Equal(t, result.Left[int](boom), res)
	assert.True(t, mu.TryLock(), "lock must be released after a failure")
}

func TestWithSemaphore_Success(t *testing.T) {
	sem := make(chan Void, 2)
	probe := &combProbe{}

	res := combRun(t.Context(), TraverseArrayPar(F.Flow2(
		probe.track(5*time.Millisecond),
		WithSemaphore[combConfig, int](sem),
	))(make([]int, 8)))

	assert.True(t, result.IsRight(res))
	assert.LessOrEqual(t, probe.peak.Load(), int32(2))
	assert.Empty(t, sem)
}

func TestWithSemaphore_Failure(t *testing.T) {
	sem := make(chan Void, 1)
	sem <- F.VOID
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()

	res := combRun(ctx, WithSemaphore[combConfig, int](sem)(combScale(1)))

	assert.Equal(t, result.Left[int](context.DeadlineExceeded), res)
}

func TestTraverseArrayVariants_Success(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := result.Of([]int{3, 6, 9, 12})

	assert.Equal(t, expected, combRun(t.Context(), TraverseArraySeq(combScale)(input)))
	assert.Equal(t, expected, combRun(t.Context(), TraverseArrayPar(combScale)(input)))
	assert.Equal(t, expected, combRun(t.Context(), TraverseArrayParN(2, combScale)(input)))
}

func TestTraverseArrayVariants_Failure(t *testing.T) {
	boom := errors.New("boom")
	failOn2 := func(n int) Effect[combConfig, int] {
		if n == 2 {
			return Fail[combConfig, int](boom)
		}
		return combScale(n)
	}
	input := []int{1, 2, 3}

	assert.Equal(t, result.Left[[]int](boom), combRun(t.Context(), TraverseArraySeq(failOn2)(input)))
	assert.Equal(t, result.Left[[]int](boom), combRun(t.Context(), TraverseArrayPar(failOn2)(input)))
	assert.Equal(t, result.Left[[]int](boom), combRun(t.Context(), TraverseArrayParN(2, failOn2)(input)))
}

func TestTraverseArrayVariants_EdgeCases(t *testing.T) {
	t.Run("Seq runs one at a time", func(t *testing.T) {
		probe := &combProbe{}
		res := combRun(t.Context(), TraverseArraySeq(probe.track(time.Millisecond))(make([]int, 5)))
		assert.True(t, result.IsRight(res))
		assert.Equal(t, int32(1), probe.peak.Load())
	})

	t.Run("ParN bounds concurrency", func(t *testing.T) {
		probe := &combProbe{}
		res := combRun(t.Context(), TraverseArrayParN(3, probe.track(5*time.Millisecond))(make([]int, 12)))
		assert.True(t, result.IsRight(res))
		assert.LessOrEqual(t, probe.peak.Load(), int32(3))
	})

	t.Run("dependencies are visible to every element", func(t *testing.T) {
		res := TraverseArrayParN(2, combScale)([]int{1, 2})(combConfig{Factor: 7})(t.Context())()
		assert.Equal(t, result.Of([]int{7, 14}), res)
	})

	t.Run("empty input", func(t *testing.T) {
		res := combRun(t.Context(), TraverseArrayParN(2, combScale)(nil))
		assert.True(t, result.IsRight(res))
		assert.Empty(t, F.Pipe1(res, result.GetOrElse(F.Constant1[error, []int](nil))))
	})
}

func TestApSeqApPar_Success(t *testing.T) {
	add := Of[combConfig](N.Add(1))

	assert.Equal(t, result.Of(7), combRun(t.Context(), ApSeq[int](combScale(2))(add)))
	assert.Equal(t, result.Of(7), combRun(t.Context(), ApPar[int](combScale(2))(add)))
}

func TestApSeqApPar_Failure(t *testing.T) {
	boom := errors.New("boom")
	add := Of[combConfig](N.Add(1))

	assert.Equal(t, result.Left[int](boom), combRun(t.Context(), ApSeq[int](Fail[combConfig, int](boom))(add)))
	assert.Equal(t, result.Left[int](boom), combRun(t.Context(), ApPar[int](Fail[combConfig, int](boom))(add)))
}

func TestApSeqApPar_EdgeCases(t *testing.T) {
	t.Run("ApPar runs both effects concurrently", func(t *testing.T) {
		slowFn := F.Pipe1(Of[combConfig](N.Add(1)), Delay[combConfig, func(int) int](30*time.Millisecond))
		slowArg := F.Pipe1(combScale(1), Delay[combConfig, int](30*time.Millisecond))

		start := time.Now()
		res := combRun(t.Context(), ApPar[int](slowArg)(slowFn))

		assert.Equal(t, result.Of(4), res)
		assert.Less(t, time.Since(start), 55*time.Millisecond)
	})
}

func TestDelay_Success(t *testing.T) {
	start := time.Now()
	res := combRun(t.Context(), Delay[combConfig, int](10*time.Millisecond)(combScale(1)))

	assert.Equal(t, result.Of(3), res)
	assert.GreaterOrEqual(t, time.Since(start), 10*time.Millisecond)
}

func TestDelay_Failure(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	res := combRun(ctx, Delay[combConfig, int](time.Second)(combScale(1)))

	assert.Equal(t, result.Left[int](context.Canceled), res)
}

func TestTimer_Success(t *testing.T) {
	start := time.Now()
	res := combRun(t.Context(), Timer[combConfig](5*time.Millisecond))

	ts, err := result.Unwrap(res)
	assert.NoError(t, err)
	assert.False(t, ts.Before(start.Add(5*time.Millisecond)))
}

func TestWithTimeout_Success(t *testing.T) {
	res := combRun(t.Context(), WithTimeout[combConfig, int](time.Second)(combScale(2)))
	assert.Equal(t, result.Of(6), res)
}

func TestWithTimeout_Failure(t *testing.T) {
	start := time.Now()
	slow := Delay[combConfig, int](time.Second)(combScale(2))

	res := combRun(t.Context(), WithTimeout[combConfig, int](10*time.Millisecond)(slow))

	assert.Equal(t, result.Left[int](context.DeadlineExceeded), res)
	assert.Less(t, time.Since(start), 500*time.Millisecond)
}

// combLogContext returns a context whose logger writes text records into buf.
func combLogContext(ctx context.Context, buf *bytes.Buffer) context.Context {
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return pair.Second(logging.WithLogger(slog.New(handler))(ctx))
}

func TestLogging_Success(t *testing.T) {
	t.Run("TapSLog logs the value", func(t *testing.T) {
		var buf bytes.Buffer
		res := combRun(combLogContext(t.Context(), &buf), TapSLog[combConfig, int]("scaled")(combScale(2)))

		assert.Equal(t, result.Of(6), res)
		assert.Contains(t, buf.String(), "level=INFO msg=scaled value=6")
	})

	t.Run("TapSLogDebug logs at debug level", func(t *testing.T) {
		var buf bytes.Buffer
		combRun(combLogContext(t.Context(), &buf), TapSLogDebug[combConfig, int]("scaled")(combScale(2)))

		assert.Contains(t, buf.String(), "level=DEBUG msg=scaled value=6")
	})

	t.Run("LogEntryExit logs entry and exit", func(t *testing.T) {
		var buf bytes.Buffer
		res := combRun(combLogContext(t.Context(), &buf), LogEntryExit[combConfig, int]("op")(combScale(2)))

		assert.Equal(t, result.Of(6), res)
		assert.Contains(t, buf.String(), "[entering]")
		assert.Contains(t, buf.String(), "[exiting ]")
	})

	t.Run("AnnotateLogs adds attributes inside the effect only", func(t *testing.T) {
		var buf bytes.Buffer
		eff := F.Pipe3(
			combScale(2),
			TapSLog[combConfig, int]("inner"),
			AnnotateLogs[combConfig, int]("requestID", "r-1"),
			TapSLog[combConfig, int]("outer"),
		)

		res := combRun(combLogContext(t.Context(), &buf), eff)

		assert.Equal(t, result.Of(6), res)
		assert.Contains(t, buf.String(), "msg=inner requestID=r-1")
		assert.Contains(t, buf.String(), "msg=outer value=6")
		assert.NotContains(t, buf.String(), "msg=outer requestID")
	})
}

func TestLogging_Failure(t *testing.T) {
	boom := errors.New("boom")
	var buf bytes.Buffer

	res := combRun(combLogContext(t.Context(), &buf), TapSLog[combConfig, int]("failed")(Fail[combConfig, int](boom)))

	assert.Equal(t, result.Left[int](boom), res)
	assert.Contains(t, buf.String(), "error=boom")
}
