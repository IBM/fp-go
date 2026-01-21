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

package readerreaderioresult

import (
	"errors"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	computation := Of[AppConfig](42)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestRight(t *testing.T) {
	computation := Right[AppConfig](42)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestLeft(t *testing.T) {
	err := errors.New("test error")
	computation := Left[AppConfig, int](err)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Left[int](err), outcome)
}

func TestMonadMap(t *testing.T) {
	computation := MonadMap(
		Of[AppConfig](21),
		N.Mul(2),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestMap(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		Map[AppConfig](N.Mul(2)),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestMonadMapTo(t *testing.T) {
	computation := MonadMapTo(Of[AppConfig](21), 99)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(99), outcome)
}

func TestMapTo(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		MapTo[AppConfig, int](99),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(99), outcome)
}

func TestMonadChain(t *testing.T) {
	computation := MonadChain(
		Of[AppConfig](21),
		func(n int) ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](n * 2)
		},
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestChain(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		Chain(func(n int) ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](n * 2)
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestMonadChainFirst(t *testing.T) {
	sideEffect := 0
	computation := MonadChainFirst(
		Of[AppConfig](42),
		func(n int) ReaderReaderIOResult[AppConfig, string] {
			sideEffect = n
			return Of[AppConfig]("ignored")
		},
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirst(t *testing.T) {
	sideEffect := 0
	computation := F.Pipe1(
		Of[AppConfig](42),
		ChainFirst(func(n int) ReaderReaderIOResult[AppConfig, string] {
			sideEffect = n
			return Of[AppConfig]("ignored")
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 42, sideEffect)
}

func TestTap(t *testing.T) {
	sideEffect := 0
	computation := F.Pipe1(
		Of[AppConfig](42),
		Tap(func(n int) ReaderReaderIOResult[AppConfig, string] {
			sideEffect = n
			return Of[AppConfig]("ignored")
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
	assert.Equal(t, 42, sideEffect)
}

func TestFlatten(t *testing.T) {
	nested := Of[AppConfig](Of[AppConfig](42))
	computation := Flatten(nested)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestFromEither(t *testing.T) {
	t.Run("right", func(t *testing.T) {
		computation := FromEither[AppConfig](either.Right[error](42))
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("left", func(t *testing.T) {
		err := errors.New("test error")
		computation := FromEither[AppConfig](either.Left[int](err))
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		computation := FromResult[AppConfig](result.Of(42))
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("test error")
		computation := FromResult[AppConfig](result.Left[int](err))
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromReader(t *testing.T) {
	computation := FromReader(func(cfg AppConfig) int {
		return len(cfg.DatabaseURL)
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(20), outcome) // len("postgres://localhost")
}

func TestRightReader(t *testing.T) {
	computation := RightReader(func(cfg AppConfig) int {
		return len(cfg.LogLevel)
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(4), outcome) // len("info")
}

func TestLeftReader(t *testing.T) {
	err := errors.New("test error")
	computation := LeftReader[int](func(cfg AppConfig) error {
		return err
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestFromIO(t *testing.T) {
	computation := FromIO[AppConfig](func() int { return 42 })
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestRightIO(t *testing.T) {
	computation := RightIO[AppConfig](func() int { return 42 })
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestLeftIO(t *testing.T) {
	err := errors.New("test error")
	computation := LeftIO[AppConfig, int](func() error { return err })
	outcome := computation(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestFromIOEither(t *testing.T) {
	t.Run("right", func(t *testing.T) {
		computation := FromIOEither[AppConfig](ioeither.Of[error](42))
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("left", func(t *testing.T) {
		err := errors.New("test error")
		computation := FromIOEither[AppConfig](ioeither.Left[int](err))
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromIOResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		computation := FromIOResult[AppConfig](func() result.Result[int] {
			return result.Of(42)
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("test error")
		computation := FromIOResult[AppConfig](func() result.Result[int] {
			return result.Left[int](err)
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromReaderIO(t *testing.T) {
	computation := FromReaderIO(func(cfg AppConfig) io.IO[int] {
		return func() int { return len(cfg.DatabaseURL) }
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(20), outcome)
}

func TestRightReaderIO(t *testing.T) {
	computation := RightReaderIO(func(cfg AppConfig) io.IO[int] {
		return func() int { return len(cfg.LogLevel) }
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(4), outcome)
}

func TestLeftReaderIO(t *testing.T) {
	err := errors.New("test error")
	computation := LeftReaderIO[int](func(cfg AppConfig) io.IO[error] {
		return func() error { return err }
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestFromReaderEither(t *testing.T) {
	t.Run("right", func(t *testing.T) {
		computation := FromReaderEither(func(cfg AppConfig) either.Either[error, int] {
			return either.Right[error](len(cfg.DatabaseURL))
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(20), outcome)
	})

	t.Run("left", func(t *testing.T) {
		err := errors.New("test error")
		computation := FromReaderEither(func(cfg AppConfig) either.Either[error, int] {
			return either.Left[int](err)
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestAsk(t *testing.T) {
	computation := Ask[AppConfig]()
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(defaultConfig), outcome)
}

func TestAsks(t *testing.T) {
	computation := Asks(func(cfg AppConfig) string {
		return cfg.DatabaseURL
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of("postgres://localhost"), outcome)
}

func TestFromOption(t *testing.T) {
	err := errors.New("none error")

	t.Run("some", func(t *testing.T) {
		computation := FromOption[AppConfig, int](func() error { return err })(option.Some(42))
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("none", func(t *testing.T) {
		computation := FromOption[AppConfig, int](func() error { return err })(option.None[int]())
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromPredicate(t *testing.T) {
	isPositive := func(n int) bool { return n > 0 }
	onFalse := func(n int) error { return errors.New("not positive") }

	t.Run("predicate true", func(t *testing.T) {
		computation := FromPredicate[AppConfig](isPositive, onFalse)(42)
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("predicate false", func(t *testing.T) {
		computation := FromPredicate[AppConfig](isPositive, onFalse)(-5)
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestMonadAlt(t *testing.T) {
	err := errors.New("first error")

	t.Run("first succeeds", func(t *testing.T) {
		first := Of[AppConfig](42)
		second := func() ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](99)
		}
		computation := MonadAlt(first, second)
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("first fails, second succeeds", func(t *testing.T) {
		first := Left[AppConfig, int](err)
		second := func() ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](99)
		}
		computation := MonadAlt(first, second)
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(99), outcome)
	})

	t.Run("both fail", func(t *testing.T) {
		first := Left[AppConfig, int](err)
		second := func() ReaderReaderIOResult[AppConfig, int] {
			return Left[AppConfig, int](errors.New("second error"))
		}
		computation := MonadAlt(first, second)
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestAlt(t *testing.T) {
	err := errors.New("first error")

	computation := F.Pipe1(
		Left[AppConfig, int](err),
		Alt(func() ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](99)
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(99), outcome)
}

func TestMonadFlap(t *testing.T) {
	fab := Of[AppConfig](N.Mul(2))
	computation := MonadFlap(fab, 21)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestFlap(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](N.Mul(2)),
		Flap[AppConfig, int](21),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestMonadMapLeft(t *testing.T) {
	err := errors.New("original error")
	computation := MonadMapLeft(
		Left[AppConfig, int](err),
		func(e error) error { return errors.New("mapped: " + e.Error()) },
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
	result.Fold(
		func(e error) any {
			assert.Contains(t, e.Error(), "mapped:")
			return nil
		},
		func(v int) any {
			t.Fatal("should be left")
			return nil
		},
	)(outcome)
}

func TestMapLeft(t *testing.T) {
	err := errors.New("original error")
	computation := F.Pipe1(
		Left[AppConfig, int](err),
		MapLeft[AppConfig, int](func(e error) error {
			return errors.New("mapped: " + e.Error())
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.True(t, result.IsLeft(outcome))
}

func TestLocal(t *testing.T) {
	type OtherConfig struct {
		URL string
	}

	computation := F.Pipe1(
		Asks(func(cfg AppConfig) string {
			return cfg.DatabaseURL
		}),
		Local[string](func(other OtherConfig) AppConfig {
			return AppConfig{DatabaseURL: other.URL, LogLevel: "debug"}
		}),
	)

	outcome := computation(OtherConfig{URL: "test-url"})(t.Context())()
	assert.Equal(t, result.Of("test-url"), outcome)
}

func TestRead(t *testing.T) {
	computation := Asks(func(cfg AppConfig) string {
		return cfg.DatabaseURL
	})

	reader := Read[string](defaultConfig)
	outcome := reader(computation)(t.Context())()
	assert.Equal(t, result.Of("postgres://localhost"), outcome)
}

func TestReadIOEither(t *testing.T) {
	computation := Asks(func(cfg AppConfig) string {
		return cfg.DatabaseURL
	})

	rio := ioeither.Of[error](defaultConfig)
	reader := ReadIOEither[string](rio)
	outcome := reader(computation)(t.Context())()
	assert.Equal(t, result.Of("postgres://localhost"), outcome)
}

func TestReadIO(t *testing.T) {
	computation := Asks(func(cfg AppConfig) string {
		return cfg.DatabaseURL
	})

	rio := func() AppConfig { return defaultConfig }
	reader := ReadIO[string](rio)
	outcome := reader(computation)(t.Context())()
	assert.Equal(t, result.Of("postgres://localhost"), outcome)
}

func TestMonadChainLeft(t *testing.T) {
	err := errors.New("original error")
	computation := MonadChainLeft(
		Left[AppConfig, int](err),
		func(e error) ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](99)
		},
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(99), outcome)
}

func TestChainLeft(t *testing.T) {
	err := errors.New("original error")
	computation := F.Pipe1(
		Left[AppConfig, int](err),
		ChainLeft(func(e error) ReaderReaderIOResult[AppConfig, int] {
			return Of[AppConfig](99)
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(99), outcome)
}

func TestDelay(t *testing.T) {
	start := time.Now()
	computation := F.Pipe1(
		Of[AppConfig](42),
		Delay[AppConfig, int](50*time.Millisecond),
	)
	outcome := computation(defaultConfig)(t.Context())()
	elapsed := time.Since(start)

	assert.Equal(t, result.Of(42), outcome)
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
}

func TestChainEitherK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		ChainEitherK[AppConfig](func(n int) either.Either[error, int] {
			return either.Right[error](n * 2)
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestChainReaderK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](10),
		ChainReaderK(func(n int) reader.Reader[AppConfig, int] {
			return func(cfg AppConfig) int {
				return n + len(cfg.LogLevel)
			}
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(14), outcome) // 10 + len("info")
}

func TestChainReaderIOK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](10),
		ChainReaderIOK(func(n int) readerio.ReaderIO[AppConfig, int] {
			return func(cfg AppConfig) io.IO[int] {
				return func() int {
					return n + len(cfg.DatabaseURL)
				}
			}
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(30), outcome) // 10 + 20
}

func TestChainReaderEitherK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](10),
		ChainReaderEitherK(func(n int) RE.ReaderEither[AppConfig, error, int] {
			return func(cfg AppConfig) either.Either[error, int] {
				return either.Right[error](n + len(cfg.LogLevel))
			}
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(14), outcome)
}

func TestChainReaderOptionK(t *testing.T) {
	onNone := func() error { return errors.New("none") }

	t.Run("some", func(t *testing.T) {
		computation := F.Pipe1(
			Of[AppConfig](10),
			ChainReaderOptionK[AppConfig, int, int](onNone)(func(n int) readeroption.ReaderOption[AppConfig, int] {
				return func(cfg AppConfig) option.Option[int] {
					return option.Some(n + len(cfg.LogLevel))
				}
			}),
		)
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(14), outcome)
	})

	t.Run("none", func(t *testing.T) {
		computation := F.Pipe1(
			Of[AppConfig](10),
			ChainReaderOptionK[AppConfig, int, int](onNone)(func(n int) readeroption.ReaderOption[AppConfig, int] {
				return func(cfg AppConfig) option.Option[int] {
					return option.None[int]()
				}
			}),
		)
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestChainIOEitherK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		ChainIOEitherK[AppConfig](func(n int) ioeither.IOEither[error, int] {
			return ioeither.Of[error](n * 2)
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestChainIOK(t *testing.T) {
	computation := F.Pipe1(
		Of[AppConfig](21),
		ChainIOK[AppConfig](func(n int) io.IO[int] {
			return func() int { return n * 2 }
		}),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestChainOptionK(t *testing.T) {
	onNone := func() error { return errors.New("none") }

	t.Run("some", func(t *testing.T) {
		computation := F.Pipe1(
			Of[AppConfig](21),
			ChainOptionK[AppConfig, int, int](onNone)(func(n int) option.Option[int] {
				return option.Some(n * 2)
			}),
		)
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("none", func(t *testing.T) {
		computation := F.Pipe1(
			Of[AppConfig](21),
			ChainOptionK[AppConfig, int, int](onNone)(func(n int) option.Option[int] {
				return option.None[int]()
			}),
		)
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestFromReaderIOResult(t *testing.T) {
	computation := FromReaderIOResult(func(cfg AppConfig) ioresult.IOResult[int] {
		return func() result.Result[int] {
			return result.Of(len(cfg.DatabaseURL))
		}
	})
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(20), outcome)
}

func TestFromReaderOption(t *testing.T) {
	onNone := func() error { return errors.New("none") }

	t.Run("some", func(t *testing.T) {
		computation := FromReaderOption[AppConfig, int](onNone)(func(cfg AppConfig) option.Option[int] {
			return option.Some(len(cfg.DatabaseURL))
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(20), outcome)
	})

	t.Run("none", func(t *testing.T) {
		computation := FromReaderOption[AppConfig, int](onNone)(func(cfg AppConfig) option.Option[int] {
			return option.None[int]()
		})
		outcome := computation(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestMonadAp(t *testing.T) {
	fab := Of[AppConfig](N.Mul(2))
	fa := Of[AppConfig](21)
	computation := MonadAp(fab, fa)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}

func TestAp(t *testing.T) {
	fa := Of[AppConfig](21)
	computation := F.Pipe1(
		Of[AppConfig](N.Mul(2)),
		Ap[int](fa),
	)
	outcome := computation(defaultConfig)(t.Context())()
	assert.Equal(t, result.Of(42), outcome)
}
