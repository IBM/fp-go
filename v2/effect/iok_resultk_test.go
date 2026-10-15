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
	"errors"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// requirePositive fails for numbers that are not positive.
func requirePositive(n int) Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

func TestMonadChainIOK_Success(t *testing.T) {
	eff := MonadChainIOK(Of[TestConfig](21), F.Flow2(N.Mul(2), io.Of[int]))
	assert.Equal(t, result.Of(42), runIOResultK(eff))
}

func TestMonadChainIOK_Failure(t *testing.T) {
	err := errors.New("upstream")
	calls := 0
	eff := MonadChainIOK(Fail[TestConfig, int](err), func(n int) IO[int] {
		return func() int { calls++; return n }
	})
	assert.Equal(t, result.Left[int](err), runIOResultK(eff))
	assert.Equal(t, 0, calls)
}

func TestMonadChainFirstIOK(t *testing.T) {
	t.Run("keeps the value and runs the side effect", func(t *testing.T) {
		var seen []int
		record := func(n int) IO[string] {
			return func() string { seen = append(seen, n); return strconv.Itoa(n) }
		}
		assert.Equal(t, result.Of(7), runIOResultK(MonadChainFirstIOK(Of[TestConfig](7), record)))
		assert.Equal(t, result.Of(8), runIOResultK(MonadTapIOK(Of[TestConfig](8), record)))
		assert.Equal(t, []int{7, 8}, seen)
	})

	t.Run("skips the side effect on failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		record := func(n int) IO[int] { return func() int { calls++; return n } }
		assert.Equal(t, result.Left[int](err), runIOResultK(MonadTapIOK(Fail[TestConfig, int](err), record)))
		assert.Equal(t, 0, calls)
	})

	t.Run("agrees with the curried form", func(t *testing.T) {
		f := F.Flow2(N.Mul(2), io.Of[int])
		assert.Equal(t,
			runIOResultK(F.Pipe1(Of[TestConfig](3), ChainFirstIOK[TestConfig](f))),
			runIOResultK(MonadChainFirstIOK(Of[TestConfig](3), f)),
		)
	})
}

func TestMonadChainResultK(t *testing.T) {
	assert.Equal(t, result.Of(5), runIOResultK(MonadChainResultK(Of[TestConfig](5), requirePositive)))
	assert.True(t, result.IsLeft(runIOResultK(MonadChainResultK(Of[TestConfig](-5), requirePositive))))
	assert.Equal(t,
		runIOResultK(F.Pipe1(Of[TestConfig](-1), ChainResultK[TestConfig](requirePositive))),
		runIOResultK(MonadChainResultK(Of[TestConfig](-1), requirePositive)),
	)
}

func TestChainFirstResultK_Success(t *testing.T) {
	check := F.Flow2(requirePositive, result.Map(strconv.Itoa))
	for _, eff := range []Effect[TestConfig, int]{
		F.Pipe1(Of[TestConfig](7), ChainFirstResultK[TestConfig](check)),
		F.Pipe1(Of[TestConfig](7), TapResultK[TestConfig](check)),
		MonadChainFirstResultK(Of[TestConfig](7), check),
		MonadTapResultK(Of[TestConfig](7), check),
	} {
		assert.Equal(t, result.Of(7), runIOResultK(eff))
	}
}

func TestChainFirstResultK_Failure(t *testing.T) {
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, eff := range []Effect[TestConfig, int]{
			F.Pipe1(Of[TestConfig](-7), ChainFirstResultK[TestConfig](requirePositive)),
			F.Pipe1(Of[TestConfig](-7), TapResultK[TestConfig](requirePositive)),
			MonadChainFirstResultK(Of[TestConfig](-7), requirePositive),
			MonadTapResultK(Of[TestConfig](-7), requirePositive),
		} {
			assert.Equal(t, result.Left[int](errors.New("not positive")), runIOResultK(eff))
		}
	})

	t.Run("upstream error is preserved", func(t *testing.T) {
		err := errors.New("upstream")
		assert.Equal(t, result.Left[int](err), runIOResultK(MonadTapResultK(Fail[TestConfig, int](err), requirePositive)))
	})
}
