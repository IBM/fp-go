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

package readerioeither

import (
	"context"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		Map[context.Context, error](utils.Double),
	)

	assert.Equal(t, E.Of[error](2), g(context.Background())())
}

func TestOrLeft(t *testing.T) {
	f := OrLeft[int](func(s string) readerio.ReaderIO[context.Context, string] {
		return readerio.Of[context.Context](s + "!")
	})

	g1 := F.Pipe1(
		Right[context.Context, string](1),
		f,
	)

	g2 := F.Pipe1(
		Left[context.Context, int]("a"),
		f,
	)

	assert.Equal(t, E.Of[string](1), g1(context.Background())())
	assert.Equal(t, E.Left[int]("a!"), g2(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context, error](utils.Double),
		Ap[int](Right[context.Context, error](1)),
	)

	assert.Equal(t, E.Right[error](2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		ChainReaderK[error](func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, E.Right[error]("1"), g(context.Background())())
}

func TestOrElseWFunc(t *testing.T) {
	type Config struct {
		retryEnabled bool
	}

	// Test with Right - should pass through unchanged
	t.Run("Right passes through", func(t *testing.T) {
		rioe := Right[Config, string](42)
		handler := OrElse(func(err string) ReaderIOEither[Config, int, int] {
			return Left[Config, int](999)
		})
		result := handler(rioe)(Config{retryEnabled: true})()
		assert.Equal(t, E.Right[int](42), result)
	})

	// Test with Left - error type widening
	t.Run("Left with error type widening", func(t *testing.T) {
		rioe := Left[Config, int]("network error")
		handler := OrElse(func(err string) ReaderIOEither[Config, int, int] {
			return func(cfg Config) IOEither[int, int] {
				if cfg.retryEnabled {
					return IOE.Right[int](100)
				}
				return IOE.Left[int](404)
			}
		})
		result := handler(rioe)(Config{retryEnabled: true})()
		assert.Equal(t, E.Right[int](100), result)
	})
}

func TestChainLeftFunc(t *testing.T) {
	type Config struct {
		errorCode int
	}

	// Test with Right - should pass through unchanged
	t.Run("Right passes through", func(t *testing.T) {
		g := F.Pipe1(
			Right[Config, string](42),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				return Left[Config, int](999)
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Right[int](42), result)
	})

	// Test with Left - error transformation with config
	t.Run("Left transforms error with config", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("error"),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				return func(cfg Config) IOEither[int, int] {
					return IOE.Left[int](cfg.errorCode)
				}
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Left[int](500), result)
	})

	// Test with Left - successful recovery
	t.Run("Left recovers successfully", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("recoverable"),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				if err == "recoverable" {
					return Right[Config, int](999)
				}
				return Left[Config, int](0)
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Right[int](999), result)
	})
}

func TestChainFirstLeftFunc(t *testing.T) {
	type Config struct {
		logEnabled bool
	}

	logged := false

	// Test with Right - should not call function
	t.Run("Right does not call function", func(t *testing.T) {
		logged = false
		g := F.Pipe1(
			Right[Config, string](42),
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				logged = true
				return Right[Config, int]("logged")
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Right[string](42), result)
		assert.False(t, logged)
	})

	// Test with Left - calls function but preserves original error
	t.Run("Left calls function but preserves error", func(t *testing.T) {
		logged = false
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				return func(cfg Config) IOEither[int, string] {
					if cfg.logEnabled {
						logged = true
					}
					return IOE.Right[int]("side effect done")
				}
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Left[int]("original error"), result)
		assert.True(t, logged)
	})

	// Test with Left - preserves original error even if side effect fails
	t.Run("Left preserves error even if side effect fails", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				return Left[Config, string](999) // Side effect fails
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Left[int]("original error"), result)
	})
}

func TestTapLeft(t *testing.T) {
	// TapLeft is an alias for ChainFirstLeft, so just a basic sanity test
	type Config struct{}

	sideEffectRan := false

	g := F.Pipe1(
		Left[Config, int]("error"),
		TapLeft[int](func(err string) ReaderIOEither[Config, string, int] {
			sideEffectRan = true
			return Right[Config, string](0)
		}),
	)

	result := g(Config{})()
	assert.Equal(t, E.Left[int]("error"), result)
	assert.True(t, sideEffectRan)
}
