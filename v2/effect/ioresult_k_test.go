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
	"context"
	"errors"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// ioResultCounter returns an IOResult Kleisli that counts its executions and
// either doubles its input or fails with err, if err is not nil.
func ioResultCounter(calls *int, err error) ioresult.Kleisli[int, int] {
	return func(n int) IOResult[int] {
		return func() Result[int] {
			*calls++
			if err != nil {
				return result.Left[int](err)
			}
			return result.Of(n * 2)
		}
	}
}

// runIOResultK runs an effect against testConfig and returns its Result.
func runIOResultK[A any](eff Effect[TestConfig, A]) Result[A] {
	return eff(testConfig)(context.Background())()
}

func TestFromIOResult_Success(t *testing.T) {
	t.Run("lifts a successful IOResult", func(t *testing.T) {
		assert.Equal(t, result.Of(42), runIOResultK(FromIOResult[TestConfig](ioresult.Of(42))))
	})

	t.Run("defers execution until the effect is run", func(t *testing.T) {
		calls := 0
		eff := FromIOResult[TestConfig](ioResultCounter(&calls, nil)(21))
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(42), runIOResultK(eff))
		assert.Equal(t, result.Of(42), runIOResultK(eff))
		assert.Equal(t, 2, calls)
	})
}

func TestFromIOResult_Failure(t *testing.T) {
	t.Run("lifts a failed IOResult", func(t *testing.T) {
		err := errors.New("io failed")
		assert.Equal(t, result.Left[int](err), runIOResultK(FromIOResult[TestConfig](ioresult.Left[int](err))))
	})
}

func TestChainIOResultK_Success(t *testing.T) {
	t.Run("chains a successful IOResult", func(t *testing.T) {
		calls := 0
		eff := F.Pipe1(Of[TestConfig](21), ChainIOResultK[TestConfig](ioResultCounter(&calls, nil)))
		assert.Equal(t, result.Of(42), runIOResultK(eff))
		assert.Equal(t, 1, calls)
	})

	t.Run("changes the value type", func(t *testing.T) {
		eff := F.Pipe1(Of[TestConfig]("42"), ChainIOResultK[TestConfig](ioresult.Eitherize1(strconv.Atoi)))
		assert.Equal(t, result.Of(42), runIOResultK(eff))
	})

	t.Run("chains multiple IOResults in order", func(t *testing.T) {
		calls := 0
		f := ioResultCounter(&calls, nil)
		eff := F.Pipe2(Of[TestConfig](5), ChainIOResultK[TestConfig](f), ChainIOResultK[TestConfig](f))
		assert.Equal(t, result.Of(20), runIOResultK(eff))
		assert.Equal(t, 2, calls)
	})

	t.Run("monadic version", func(t *testing.T) {
		calls := 0
		eff := MonadChainIOResultK(Of[TestConfig](21), ioResultCounter(&calls, nil))
		assert.Equal(t, result.Of(42), runIOResultK(eff))
		assert.Equal(t, 1, calls)
	})
}

func TestChainIOResultK_Failure(t *testing.T) {
	t.Run("propagates the error of the IOResult", func(t *testing.T) {
		err := errors.New("io failed")
		eff := F.Pipe1(Of[TestConfig](21), ChainIOResultK[TestConfig](ioResultCounter(new(int), err)))
		assert.Equal(t, result.Left[int](err), runIOResultK(eff))
	})

	t.Run("propagates parse errors", func(t *testing.T) {
		eff := F.Pipe1(Of[TestConfig]("abc"), ChainIOResultK[TestConfig](ioresult.Eitherize1(strconv.Atoi)))
		assert.True(t, result.IsLeft(runIOResultK(eff)))
	})

	t.Run("short-circuits subsequent steps", func(t *testing.T) {
		err := errors.New("io failed")
		calls := 0
		eff := F.Pipe2(
			Of[TestConfig](21),
			ChainIOResultK[TestConfig](ioResultCounter(new(int), err)),
			ChainIOResultK[TestConfig](ioResultCounter(&calls, nil)),
		)
		assert.Equal(t, result.Left[int](err), runIOResultK(eff))
		assert.Equal(t, 0, calls)
	})

	t.Run("skips the IOResult on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		eff := MonadChainIOResultK(Fail[TestConfig, int](err), ioResultCounter(&calls, nil))
		assert.Equal(t, result.Left[int](err), runIOResultK(eff))
		assert.Equal(t, 0, calls)
	})
}

func TestChainFirstIOResultK_Success(t *testing.T) {
	t.Run("all variants keep the original value", func(t *testing.T) {
		calls := 0
		f := ioResultCounter(&calls, nil)
		for _, eff := range []Effect[TestConfig, int]{
			F.Pipe1(Of[TestConfig](21), ChainFirstIOResultK[TestConfig](f)),
			F.Pipe1(Of[TestConfig](21), TapIOResultK[TestConfig](f)),
			MonadChainFirstIOResultK(Of[TestConfig](21), f),
			MonadTapIOResultK(Of[TestConfig](21), f),
		} {
			assert.Equal(t, result.Of(21), runIOResultK(eff))
		}
		assert.Equal(t, 4, calls)
	})

	t.Run("side effect may have a different type", func(t *testing.T) {
		var seen string
		record := func(n int) IOResult[string] {
			return func() Result[string] {
				seen = strconv.Itoa(n)
				return result.Of(seen)
			}
		}
		eff := F.Pipe1(Of[TestConfig](7), TapIOResultK[TestConfig](record))
		assert.Equal(t, result.Of(7), runIOResultK(eff))
		assert.Equal(t, "7", seen)
	})
}

func TestChainFirstIOResultK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		f := ioResultCounter(new(int), err)
		for _, eff := range []Effect[TestConfig, int]{
			F.Pipe1(Of[TestConfig](21), ChainFirstIOResultK[TestConfig](f)),
			F.Pipe1(Of[TestConfig](21), TapIOResultK[TestConfig](f)),
			MonadChainFirstIOResultK(Of[TestConfig](21), f),
			MonadTapIOResultK(Of[TestConfig](21), f),
		} {
			assert.Equal(t, result.Left[int](err), runIOResultK(eff))
		}
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		f := ioResultCounter(&calls, nil)
		for _, eff := range []Effect[TestConfig, int]{
			F.Pipe1(Fail[TestConfig, int](err), ChainFirstIOResultK[TestConfig](f)),
			F.Pipe1(Fail[TestConfig, int](err), TapIOResultK[TestConfig](f)),
			MonadChainFirstIOResultK(Fail[TestConfig, int](err), f),
			MonadTapIOResultK(Fail[TestConfig, int](err), f),
		} {
			assert.Equal(t, result.Left[int](err), runIOResultK(eff))
		}
		assert.Equal(t, 0, calls)
	})
}

func TestIOResultK_EdgeCases(t *testing.T) {
	t.Run("ChainIOResultK agrees with Chain over FromIOResult", func(t *testing.T) {
		err := errors.New("io failed")
		for _, f := range []ioresult.Kleisli[int, int]{
			ioResultCounter(new(int), nil),
			ioResultCounter(new(int), err),
		} {
			for _, start := range []Effect[TestConfig, int]{Of[TestConfig](0), Fail[TestConfig, int](err)} {
				assert.Equal(t,
					runIOResultK(F.Pipe1(start, Chain(F.Flow2(f, FromIOResult[TestConfig, int])))),
					runIOResultK(F.Pipe1(start, ChainIOResultK[TestConfig](f))),
				)
			}
		}
	})

	t.Run("works with the zero value", func(t *testing.T) {
		eff := F.Pipe1(Of[TestConfig](0), ChainIOResultK[TestConfig](ioResultCounter(new(int), nil)))
		assert.Equal(t, result.Of(0), runIOResultK(eff))
	})
}

func TestIOResultK_Integration(t *testing.T) {
	type state struct {
		Raw    string
		Parsed int
		Double int
	}
	setParsed := func(v int) func(state) state { return func(s state) state { s.Parsed = v; return s } }
	setDouble := func(v int) func(state) state { return func(s state) state { s.Double = v; return s } }
	parse := func(s state) IOResult[int] { return ioresult.Eitherize1(strconv.Atoi)(s.Raw) }
	double := func(s state) IOResult[int] { return ioresult.Of(s.Parsed * 2) }
	audited := 0
	audit := func(s state) IOResult[state] {
		return func() Result[state] { audited++; return result.Of(s) }
	}

	pipeline := func(raw string) Effect[TestConfig, int] {
		return F.Pipe5(
			Do[TestConfig](state{Raw: raw}),
			BindIOResultK[TestConfig](setParsed, parse),
			BindIOResultK[TestConfig](setDouble, double),
			TapIOResultK[TestConfig](audit),
			Map[TestConfig](func(s state) int { return s.Double }),
			ChainIOResultK[TestConfig](ioresult.Of[int]),
		)
	}

	t.Run("composes with do-notation", func(t *testing.T) {
		assert.Equal(t, result.Of(42), runIOResultK(pipeline("21")))
		assert.Equal(t, 1, audited)
	})

	t.Run("fails fast inside do-notation", func(t *testing.T) {
		audited = 0
		assert.True(t, result.IsLeft(runIOResultK(pipeline("x"))))
		assert.Equal(t, 0, audited)
	})

	t.Run("runs through RunSync and Provide", func(t *testing.T) {
		value, err := runEffect(F.Pipe1(Of[TestContext](21), ChainIOResultK[TestContext](ioResultCounter(new(int), nil))), TestContext{})
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})
}

func BenchmarkChainIOResultK(b *testing.B) {
	eff := F.Pipe1(Of[TestConfig](21), ChainIOResultK[TestConfig](F.Flow2(N.Mul(2), ioresult.Of[int])))
	ctx := context.Background()
	for b.Loop() {
		eff(testConfig)(ctx)()
	}
}
