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

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type ctxKey struct{}

// testMultiplier reads the multiplier of the test configuration.
func testMultiplier(c TestConfig) int { return c.Multiplier }

func TestFromReaderIO_Success(t *testing.T) {
	t.Run("reads the context and runs the IO", func(t *testing.T) {
		eff := FromReaderIO(F.Flow2(testMultiplier, io.Of[int]))
		assert.Equal(t, result.Of(testConfig.Multiplier), runIOResultK(eff))
	})

	t.Run("defers execution until the effect is run", func(t *testing.T) {
		calls := 0
		eff := FromReaderIO(func(c TestConfig) IO[int] {
			return func() int { calls++; return c.Multiplier }
		})
		assert.Equal(t, 0, calls)
		runIOResultK(eff)
		runIOResultK(eff)
		assert.Equal(t, 2, calls)
	})
}

func TestFromReaderIOResult_Success(t *testing.T) {
	t.Run("lifts a successful ReaderIOResult", func(t *testing.T) {
		eff := FromReaderIOResult(F.Flow2(testMultiplier, ioresult.Of[int]))
		assert.Equal(t, result.Of(testConfig.Multiplier), runIOResultK(eff))
	})

	t.Run("agrees with the readerioresult value", func(t *testing.T) {
		rio := RIOR.Asks[TestConfig](testMultiplier)
		assert.Equal(t, rio(testConfig)(), runIOResultK(FromReaderIOResult(rio)))
	})
}

func TestFromReaderIOResult_Failure(t *testing.T) {
	t.Run("lifts a failed ReaderIOResult", func(t *testing.T) {
		err := errors.New("rio failed")
		eff := FromReaderIOResult(RIOR.Left[TestConfig, int](err))
		assert.Equal(t, result.Left[int](err), runIOResultK(eff))
	})
}

func TestFromReaderThunk_Success(t *testing.T) {
	t.Run("sees both the context and context.Context", func(t *testing.T) {
		eff := FromReaderThunk(func(c TestConfig) Thunk[string] {
			return func(ctx context.Context) IOResult[string] {
				return ioresult.Of(c.Prefix + ctx.Value(ctxKey{}).(string))
			}
		})
		ctx := context.WithValue(context.Background(), ctxKey{}, "-ctx")
		assert.Equal(t, result.Of(testConfig.Prefix+"-ctx"), eff(testConfig)(ctx)())
	})
}

func TestFromReaderThunk_Failure(t *testing.T) {
	t.Run("lifts a failed Thunk", func(t *testing.T) {
		err := errors.New("thunk failed")
		eff := FromReaderThunk(F.Constant1[TestConfig](thunk.Left[int](err)))
		assert.Equal(t, result.Left[int](err), runIOResultK(eff))
	})
}

func TestFlatten_Success(t *testing.T) {
	t.Run("runs the inner effect with the same context", func(t *testing.T) {
		eff := Flatten(Of[TestConfig](Asks(testMultiplier)))
		assert.Equal(t, result.Of(testConfig.Multiplier), runIOResultK(eff))
	})
}

func TestFlatten_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("outer failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(Flatten(Fail[TestConfig, Effect[TestConfig, int]](err))))
	})

	t.Run("inner failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(Flatten(Of[TestConfig](Fail[TestConfig, int](err)))))
	})
}

func TestFlattenIO_Success(t *testing.T) {
	t.Run("runs the produced IO", func(t *testing.T) {
		calls := 0
		eff := F.Pipe2(
			Of[TestConfig](21),
			Map[TestConfig](func(n int) IO[int] { return func() int { calls++; return n * 2 } }),
			FlattenIO[TestConfig, int],
		)
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(42), runIOResultK(eff))
		assert.Equal(t, 1, calls)
	})
}

func TestFlattenIO_Failure(t *testing.T) {
	t.Run("upstream failure skips the IO", func(t *testing.T) {
		err := errors.New("upstream")
		assert.Equal(t, result.Left[int](err), runIOResultK(FlattenIO(Fail[TestConfig, IO[int]](err))))
	})
}

func TestFlattenIOResult_Success(t *testing.T) {
	t.Run("runs the produced IOResult", func(t *testing.T) {
		eff := F.Pipe2(
			Of[TestConfig]("42"),
			Map[TestConfig](ioresult.Eitherize1(strconv.Atoi)),
			FlattenIOResult[TestConfig, int],
		)
		assert.Equal(t, result.Of(42), runIOResultK(eff))
	})
}

func TestFlattenIOResult_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("failed IOResult fails the effect", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(FlattenIOResult(Of[TestConfig](ioresult.Left[int](err)))))
	})

	t.Run("upstream failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(FlattenIOResult(Fail[TestConfig, IOResult[int]](err))))
	})
}

func TestFlattenThunk_Success(t *testing.T) {
	t.Run("runs the produced Thunk with context.Context", func(t *testing.T) {
		fromCtx := func(ctx context.Context) IOResult[string] {
			return ioresult.Of(ctx.Value(ctxKey{}).(string))
		}
		eff := FlattenThunk(Of[TestConfig](Thunk[string](fromCtx)))
		ctx := context.WithValue(context.Background(), ctxKey{}, "value")
		assert.Equal(t, result.Of("value"), eff(testConfig)(ctx)())
	})
}

func TestFlattenThunk_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("failed Thunk fails the effect", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(FlattenThunk(Of[TestConfig](thunk.Left[int](err)))))
	})

	t.Run("upstream failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), runIOResultK(FlattenThunk(Fail[TestConfig, Thunk[int]](err))))
	})
}

func TestFlatten_EdgeCases(t *testing.T) {
	t.Run("flatten variants agree with their Chain counterparts", func(t *testing.T) {
		start := Of[TestConfig](21)
		double := F.Flow2(N.Mul(2), ioresult.Of[int])
		assert.Equal(t,
			runIOResultK(F.Pipe1(start, ChainIOResultK[TestConfig](double))),
			runIOResultK(FlattenIOResult(F.Pipe1(start, Map[TestConfig](double)))),
		)
		assert.Equal(t,
			runIOResultK(F.Pipe1(start, ChainThunkK[TestConfig](thunk.Of[int]))),
			runIOResultK(FlattenThunk(F.Pipe1(start, Map[TestConfig](thunk.Of[int])))),
		)
	})
}
