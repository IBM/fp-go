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

package readerreaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	RO "github.com/IBM/fp-go/v2/readeroption"
	"github.com/stretchr/testify/assert"
)

type OuterConfig struct {
	database string
	logLevel string
}

type InnerConfig struct {
	apiKey  string
	timeout int
}

func TestOf(t *testing.T) {
	result := Of[OuterConfig, InnerConfig, error](42)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestRight(t *testing.T) {
	result := Right[OuterConfig, InnerConfig, error](42)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestLeft(t *testing.T) {
	err := errors.New("test error")
	result := Left[OuterConfig, InnerConfig, int](err)
	assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
}

func TestMap(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		Map[OuterConfig, InnerConfig, error](utils.Double),
	)
	assert.Equal(t, E.Right[error](2), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadMap(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadMap(fa, utils.Double)
	assert.Equal(t, E.Right[error](2), result(OuterConfig{})(InnerConfig{})())
}

func TestMapTo(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		MapTo[OuterConfig, InnerConfig, error, int]("mapped"),
	)
	assert.Equal(t, E.Right[error]("mapped"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadMapTo(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadMapTo(fa, "mapped")
	assert.Equal(t, E.Right[error]("mapped"), result(OuterConfig{})(InnerConfig{})())
}

func TestChain(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		Chain(func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			return Of[OuterConfig, InnerConfig, error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChain(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChain(fa, func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
		return Of[OuterConfig, InnerConfig, error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirst(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirst(func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			return Of[OuterConfig, InnerConfig, error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirst(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirst(fa, func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
		return Of[OuterConfig, InnerConfig, error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTap(t *testing.T) {
	sideEffect := 0
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		Tap(func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			sideEffect = v * 2
			return Of[OuterConfig, InnerConfig, error]("ignored")
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, 2, sideEffect)
}

func TestMonadTap(t *testing.T) {
	sideEffect := 0
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTap(fa, func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
		sideEffect = v * 2
		return Of[OuterConfig, InnerConfig, error]("ignored")
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, 2, sideEffect)
}

func TestFlatten(t *testing.T) {
	nested := Of[OuterConfig, InnerConfig, error](Of[OuterConfig, InnerConfig, error](42))
	result := Flatten(nested)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[OuterConfig, InnerConfig, error](utils.Double),
		Ap[int](Right[OuterConfig, InnerConfig, error](1)),
	)
	assert.Equal(t, E.Right[error](2), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadAp(t *testing.T) {
	fab := Right[OuterConfig, InnerConfig, error](utils.Double)
	fa := Right[OuterConfig, InnerConfig, error](1)
	result := MonadAp(fab, fa)
	assert.Equal(t, E.Right[error](2), result(OuterConfig{})(InnerConfig{})())
}

func TestMonadApSeq(t *testing.T) {
	fab := Right[OuterConfig, InnerConfig, error](utils.Double)
	fa := Right[OuterConfig, InnerConfig, error](1)
	result := MonadApSeq(fab, fa)
	assert.Equal(t, E.Right[error](2), result(OuterConfig{})(InnerConfig{})())
}

func TestMonadApPar(t *testing.T) {
	fab := Right[OuterConfig, InnerConfig, error](utils.Double)
	fa := Right[OuterConfig, InnerConfig, error](1)
	result := MonadApPar(fab, fa)
	assert.Equal(t, E.Right[error](2), result(OuterConfig{})(InnerConfig{})())
}

func TestFromEither(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		result := FromEither[OuterConfig, InnerConfig](E.Right[error](42))
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Left", func(t *testing.T) {
		err := errors.New("test error")
		result := FromEither[OuterConfig, InnerConfig](E.Left[int](err))
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromReader(t *testing.T) {
	reader := R.Of[OuterConfig](42)
	result := FromReader[InnerConfig, error](reader)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestRightReader(t *testing.T) {
	reader := R.Of[OuterConfig](42)
	result := RightReader[InnerConfig, error](reader)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestLeftReader(t *testing.T) {
	err := errors.New("test error")
	reader := R.Of[OuterConfig](err)
	result := LeftReader[InnerConfig, int](reader)
	assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
}

func TestFromIO(t *testing.T) {
	ioVal := io.Of(42)
	result := FromIO[OuterConfig, InnerConfig, error](ioVal)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestRightIO(t *testing.T) {
	ioVal := io.Of(42)
	result := RightIO[OuterConfig, InnerConfig, error](ioVal)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestLeftIO(t *testing.T) {
	err := errors.New("test error")
	ioVal := io.Of(err)
	result := LeftIO[OuterConfig, InnerConfig, int](ioVal)
	assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
}

func TestFromIOEither(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		ioe := IOE.Right[error](42)
		result := FromIOEither[OuterConfig, InnerConfig](ioe)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Left", func(t *testing.T) {
		err := errors.New("test error")
		ioe := IOE.Left[int](err)
		result := FromIOEither[OuterConfig, InnerConfig](ioe)
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromReaderIO(t *testing.T) {
	rio := readerio.Of[OuterConfig](42)
	result := FromReaderIO[InnerConfig, error](rio)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestRightReaderIO(t *testing.T) {
	rio := readerio.Of[OuterConfig](42)
	result := RightReaderIO[InnerConfig, error](rio)
	assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
}

func TestLeftReaderIO(t *testing.T) {
	err := errors.New("test error")
	rio := readerio.Of[OuterConfig](err)
	result := LeftReaderIO[InnerConfig, int](rio)
	assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
}

func TestFromReaderEither(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		re := RE.Right[OuterConfig, error](42)
		result := FromReaderEither[OuterConfig, InnerConfig](re)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Left", func(t *testing.T) {
		err := errors.New("test error")
		re := RE.Left[OuterConfig, int](err)
		result := FromReaderEither[OuterConfig, InnerConfig](re)
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromReaderIOEither(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		rioe := RIOE.Right[OuterConfig, error](42)
		result := FromReaderIOEither[InnerConfig](rioe)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Left", func(t *testing.T) {
		err := errors.New("test error")
		rioe := RIOE.Left[OuterConfig, int](err)
		result := FromReaderIOEither[InnerConfig](rioe)
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromReaderOption(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }

	t.Run("Some", func(t *testing.T) {
		ro := RO.Of[OuterConfig](42)
		result := FromReaderOption[OuterConfig, InnerConfig, int](onNone)(ro)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("None", func(t *testing.T) {
		ro := RO.None[OuterConfig, int]()
		result := FromReaderOption[OuterConfig, InnerConfig, int](onNone)(ro)
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromOption(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }

	t.Run("Some", func(t *testing.T) {
		opt := O.Some(42)
		result := FromOption[OuterConfig, InnerConfig, int](onNone)(opt)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("None", func(t *testing.T) {
		opt := O.None[int]()
		result := FromOption[OuterConfig, InnerConfig, int](onNone)(opt)
		assert.Equal(t, E.Left[int](err), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestFromPredicate(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("not positive: %d", n) }

	t.Run("Predicate true", func(t *testing.T) {
		result := FromPredicate[OuterConfig, InnerConfig](isPositive, onFalse)(5)
		assert.Equal(t, E.Right[error](5), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Predicate false", func(t *testing.T) {
		result := FromPredicate[OuterConfig, InnerConfig](isPositive, onFalse)(-5)
		expected := E.Left[int](fmt.Errorf("not positive: -5"))
		assert.Equal(t, expected, result(OuterConfig{})(InnerConfig{})())
	})
}

func TestAsk(t *testing.T) {
	outer := OuterConfig{database: "postgres", logLevel: "info"}
	result := Ask[OuterConfig, InnerConfig, error]()
	assert.Equal(t, E.Right[error](outer), result(outer)(InnerConfig{})())
}

func TestAsks(t *testing.T) {
	outer := OuterConfig{database: "postgres", logLevel: "info"}
	reader := R.Asks(func(cfg OuterConfig) string { return cfg.database })
	result := Asks[InnerConfig, error](reader)
	assert.Equal(t, E.Right[error]("postgres"), result(outer)(InnerConfig{})())
}

func TestLocal(t *testing.T) {
	outer1 := OuterConfig{database: "postgres", logLevel: "info"}
	outer2 := OuterConfig{database: "mysql", logLevel: "debug"}

	computation := Asks[InnerConfig, error](R.Asks(func(cfg OuterConfig) string {
		return cfg.database
	}))

	modified := Local[InnerConfig, error, string](func(cfg OuterConfig) OuterConfig {
		return outer2
	})(computation)

	assert.Equal(t, E.Right[error]("mysql"), modified(outer1)(InnerConfig{})())
}

func TestRead(t *testing.T) {
	outer := OuterConfig{database: "postgres", logLevel: "info"}
	computation := Asks[InnerConfig, error](R.Asks(func(cfg OuterConfig) string {
		return cfg.database
	}))

	result := Read[InnerConfig, error, string](outer)(computation)
	assert.Equal(t, E.Right[error]("postgres"), result(InnerConfig{})())
}

func TestChainEitherK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainEitherK[OuterConfig, InnerConfig](func(v int) E.Either[error, string] {
			return E.Right[error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainEitherK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainEitherK(fa, func(v int) E.Either[error, string] {
		return E.Right[error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirstEitherK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirstEitherK[OuterConfig, InnerConfig](func(v int) E.Either[error, string] {
			return E.Right[error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirstEitherK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirstEitherK(fa, func(v int) E.Either[error, string] {
		return E.Right[error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTapEitherK(t *testing.T) {
	sideEffect := ""
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		TapEitherK[OuterConfig, InnerConfig](func(v int) E.Either[error, string] {
			sideEffect = fmt.Sprintf("%d", v)
			return E.Right[error](sideEffect)
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, "1", sideEffect)
}

func TestMonadTapEitherK(t *testing.T) {
	sideEffect := ""
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTapEitherK(fa, func(v int) E.Either[error, string] {
		sideEffect = fmt.Sprintf("%d", v)
		return E.Right[error](sideEffect)
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, "1", sideEffect)
}

func TestChainReaderK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainReaderK[InnerConfig, error](func(v int) R.Reader[OuterConfig, string] {
			return R.Of[OuterConfig](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainReaderK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainReaderK(fa, func(v int) R.Reader[OuterConfig, string] {
		return R.Of[OuterConfig](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirstReaderK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirstReaderK[InnerConfig, error](func(v int) R.Reader[OuterConfig, string] {
			return R.Of[OuterConfig](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirstReaderK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirstReaderK(fa, func(v int) R.Reader[OuterConfig, string] {
		return R.Of[OuterConfig](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTapReaderK(t *testing.T) {
	sideEffect := ""
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		TapReaderK[InnerConfig, error](func(v int) R.Reader[OuterConfig, string] {
			sideEffect = fmt.Sprintf("%d", v)
			return R.Of[OuterConfig](sideEffect)
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, "1", sideEffect)
}

func TestMonadTapReaderK(t *testing.T) {
	sideEffect := ""
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTapReaderK(fa, func(v int) R.Reader[OuterConfig, string] {
		sideEffect = fmt.Sprintf("%d", v)
		return R.Of[OuterConfig](sideEffect)
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, "1", sideEffect)
}

func TestChainReaderIOK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainReaderIOK[InnerConfig, error](func(v int) readerio.ReaderIO[OuterConfig, string] {
			return readerio.Of[OuterConfig](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainReaderIOK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainReaderIOK(fa, func(v int) readerio.ReaderIO[OuterConfig, string] {
		return readerio.Of[OuterConfig](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirstReaderIOK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirstReaderIOK[InnerConfig, error](func(v int) readerio.ReaderIO[OuterConfig, string] {
			return readerio.Of[OuterConfig](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirstReaderIOK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirstReaderIOK(fa, func(v int) readerio.ReaderIO[OuterConfig, string] {
		return readerio.Of[OuterConfig](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTapReaderIOK(t *testing.T) {
	sideEffect := ""
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		TapReaderIOK[InnerConfig, error](func(v int) readerio.ReaderIO[OuterConfig, string] {
			sideEffect = fmt.Sprintf("%d", v)
			return readerio.Of[OuterConfig](sideEffect)
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, "1", sideEffect)
}

func TestMonadTapReaderIOK(t *testing.T) {
	sideEffect := ""
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTapReaderIOK(fa, func(v int) readerio.ReaderIO[OuterConfig, string] {
		sideEffect = fmt.Sprintf("%d", v)
		return readerio.Of[OuterConfig](sideEffect)
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, "1", sideEffect)
}

func TestChainReaderEitherK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainReaderEitherK[InnerConfig](func(v int) RE.ReaderEither[OuterConfig, error, string] {
			return RE.Right[OuterConfig, error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainReaderEitherK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainReaderEitherK(fa, func(v int) RE.ReaderEither[OuterConfig, error, string] {
		return RE.Right[OuterConfig, error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirstReaderEitherK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirstReaderEitherK[InnerConfig](func(v int) RE.ReaderEither[OuterConfig, error, string] {
			return RE.Right[OuterConfig, error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirstReaderEitherK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirstReaderEitherK(fa, func(v int) RE.ReaderEither[OuterConfig, error, string] {
		return RE.Right[OuterConfig, error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTapReaderEitherK(t *testing.T) {
	sideEffect := ""
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		TapReaderEitherK[InnerConfig](func(v int) RE.ReaderEither[OuterConfig, error, string] {
			sideEffect = fmt.Sprintf("%d", v)
			return RE.Right[OuterConfig, error](sideEffect)
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, "1", sideEffect)
}

func TestMonadTapReaderEitherK(t *testing.T) {
	sideEffect := ""
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTapReaderEitherK(fa, func(v int) RE.ReaderEither[OuterConfig, error, string] {
		sideEffect = fmt.Sprintf("%d", v)
		return RE.Right[OuterConfig, error](sideEffect)
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, "1", sideEffect)
}

func TestChainReaderOptionK(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }

	t.Run("Some", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainReaderOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) RO.ReaderOption[OuterConfig, string] {
				return RO.Some[OuterConfig](fmt.Sprintf("%d", v))
			}),
		)
		assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
	})

	t.Run("None", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainReaderOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) RO.ReaderOption[OuterConfig, string] {
				return RO.None[OuterConfig, string]()
			}),
		)
		assert.Equal(t, E.Left[string](err), g(OuterConfig{})(InnerConfig{})())
	})
}

func TestChainFirstReaderOptionK(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }

	t.Run("Some", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainFirstReaderOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) RO.ReaderOption[OuterConfig, string] {
				return RO.Some[OuterConfig](fmt.Sprintf("%d", v))
			}),
		)
		assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
	})

	t.Run("None", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainFirstReaderOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) RO.ReaderOption[OuterConfig, string] {
				return RO.None[OuterConfig, string]()
			}),
		)
		assert.Equal(t, E.Left[int](err), g(OuterConfig{})(InnerConfig{})())
	})
}

func TestTapReaderOptionK(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }
	sideEffect := ""

	t.Run("Some", func(t *testing.T) {
		sideEffect = ""
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			TapReaderOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) RO.ReaderOption[OuterConfig, string] {
				sideEffect = fmt.Sprintf("%d", v)
				return RO.Some[OuterConfig](sideEffect)
			}),
		)
		result := g(OuterConfig{})(InnerConfig{})()
		assert.Equal(t, E.Right[error](1), result)
		assert.Equal(t, "1", sideEffect)
	})
}

func TestChainIOEitherK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainIOEitherK[OuterConfig, InnerConfig](func(v int) IOE.IOEither[error, string] {
			return IOE.Right[error](fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainIOEitherK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainIOEitherK(fa, func(v int) IOE.IOEither[error, string] {
		return IOE.Right[error](fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainIOK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainIOK[OuterConfig, InnerConfig, error](func(v int) io.IO[string] {
			return io.Of(fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainIOK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainIOK(fa, func(v int) io.IO[string] {
		return io.Of(fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error]("1"), result(OuterConfig{})(InnerConfig{})())
}

func TestChainFirstIOK(t *testing.T) {
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		ChainFirstIOK[OuterConfig, InnerConfig, error](func(v int) io.IO[string] {
			return io.Of(fmt.Sprintf("%d", v))
		}),
	)
	assert.Equal(t, E.Right[error](1), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadChainFirstIOK(t *testing.T) {
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadChainFirstIOK(fa, func(v int) io.IO[string] {
		return io.Of(fmt.Sprintf("%d", v))
	})
	assert.Equal(t, E.Right[error](1), result(OuterConfig{})(InnerConfig{})())
}

func TestTapIOK(t *testing.T) {
	sideEffect := ""
	g := F.Pipe1(
		Of[OuterConfig, InnerConfig, error](1),
		TapIOK[OuterConfig, InnerConfig, error](func(v int) io.IO[string] {
			sideEffect = fmt.Sprintf("%d", v)
			return io.Of(sideEffect)
		}),
	)
	result := g(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), result)
	assert.Equal(t, "1", sideEffect)
}

func TestMonadTapIOK(t *testing.T) {
	sideEffect := ""
	fa := Of[OuterConfig, InnerConfig, error](1)
	result := MonadTapIOK(fa, func(v int) io.IO[string] {
		sideEffect = fmt.Sprintf("%d", v)
		return io.Of(sideEffect)
	})
	outcome := result(OuterConfig{})(InnerConfig{})()
	assert.Equal(t, E.Right[error](1), outcome)
	assert.Equal(t, "1", sideEffect)
}

func TestChainOptionK(t *testing.T) {
	err := errors.New("none")
	onNone := func() error { return err }

	t.Run("Some", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) O.Option[string] {
				return O.Some(fmt.Sprintf("%d", v))
			}),
		)
		assert.Equal(t, E.Right[error]("1"), g(OuterConfig{})(InnerConfig{})())
	})

	t.Run("None", func(t *testing.T) {
		g := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](1),
			ChainOptionK[OuterConfig, InnerConfig, int, string](onNone)(func(v int) O.Option[string] {
				return O.None[string]()
			}),
		)
		assert.Equal(t, E.Left[string](err), g(OuterConfig{})(InnerConfig{})())
	})
}

func TestMonadAlt(t *testing.T) {
	t.Run("First succeeds", func(t *testing.T) {
		first := Right[OuterConfig, InnerConfig, error](42)
		second := func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Right[OuterConfig, InnerConfig, error](99)
		}
		result := MonadAlt(first, second)
		assert.Equal(t, E.Right[error](42), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("First fails, second succeeds", func(t *testing.T) {
		err := errors.New("first error")
		first := Left[OuterConfig, InnerConfig, int](err)
		second := func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Right[OuterConfig, InnerConfig, error](99)
		}
		result := MonadAlt(first, second)
		assert.Equal(t, E.Right[error](99), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Both fail", func(t *testing.T) {
		err1 := errors.New("first error")
		err2 := errors.New("second error")
		first := Left[OuterConfig, InnerConfig, int](err1)
		second := func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Left[OuterConfig, InnerConfig, int](err2)
		}
		result := MonadAlt(first, second)
		assert.Equal(t, E.Left[int](err2), result(OuterConfig{})(InnerConfig{})())
	})
}

func TestAlt(t *testing.T) {
	t.Run("First succeeds", func(t *testing.T) {
		second := func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Right[OuterConfig, InnerConfig, error](99)
		}
		g := F.Pipe1(
			Right[OuterConfig, InnerConfig, error](42),
			Alt(second),
		)
		assert.Equal(t, E.Right[error](42), g(OuterConfig{})(InnerConfig{})())
	})

	t.Run("First fails, second succeeds", func(t *testing.T) {
		err := errors.New("first error")
		second := func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Right[OuterConfig, InnerConfig, error](99)
		}
		g := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](err),
			Alt(second),
		)
		assert.Equal(t, E.Right[error](99), g(OuterConfig{})(InnerConfig{})())
	})
}

func TestFlap(t *testing.T) {
	fab := Right[OuterConfig, InnerConfig, error](utils.Double)
	g := F.Pipe1(
		fab,
		Flap[OuterConfig, InnerConfig, error, int](1),
	)
	assert.Equal(t, E.Right[error](2), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadFlap(t *testing.T) {
	fab := Right[OuterConfig, InnerConfig, error](utils.Double)
	result := MonadFlap(fab, 1)
	assert.Equal(t, E.Right[error](2), result(OuterConfig{})(InnerConfig{})())
}

func TestMapLeft(t *testing.T) {
	err := errors.New("original error")
	g := F.Pipe1(
		Left[OuterConfig, InnerConfig, int](err),
		MapLeft[OuterConfig, InnerConfig, int](func(e error) string {
			return e.Error() + " transformed"
		}),
	)
	assert.Equal(t, E.Left[int]("original error transformed"), g(OuterConfig{})(InnerConfig{})())
}

func TestMonadMapLeft(t *testing.T) {
	err := errors.New("original error")
	fa := Left[OuterConfig, InnerConfig, int](err)
	result := MonadMapLeft(fa, func(e error) string {
		return e.Error() + " transformed"
	})
	assert.Equal(t, E.Left[int]("original error transformed"), result(OuterConfig{})(InnerConfig{})())
}

func TestMapLeftDoesNotAffectRight(t *testing.T) {
	g := F.Pipe1(
		Right[OuterConfig, InnerConfig, error](42),
		MapLeft[OuterConfig, InnerConfig, int](func(e error) string {
			return "should not be called"
		}),
	)

	assert.Equal(t, E.Right[string](42), g(OuterConfig{})(InnerConfig{})())
}

func TestMultiLayerContext(t *testing.T) {
	outer := OuterConfig{database: "postgres", logLevel: "info"}
	inner := InnerConfig{apiKey: "secret", timeout: 30}

	// Create a computation that uses both contexts
	computation := func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
		return func(c InnerConfig) IOE.IOEither[error, string] {
			return IOE.Right[error](fmt.Sprintf("db=%s, key=%s", r.database, c.apiKey))
		}
	}

	result := computation(outer)(inner)()
	assert.Equal(t, E.Right[error]("db=postgres, key=secret"), result)
}

func TestCompositionWithBothContexts(t *testing.T) {
	outer := OuterConfig{database: "postgres", logLevel: "info"}
	inner := InnerConfig{apiKey: "secret", timeout: 30}

	// Build a pipeline that uses both contexts
	pipeline := F.Pipe2(
		Ask[OuterConfig, InnerConfig, error](),
		Map[OuterConfig, InnerConfig, error](func(cfg OuterConfig) string {
			return cfg.database
		}),
		Chain(func(db string) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			return func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
				return func(c InnerConfig) IOE.IOEither[error, string] {
					return IOE.Right[error](fmt.Sprintf("%s:%s", db, c.apiKey))
				}
			}
		}),
	)

	result := pipeline(outer)(inner)()
	assert.Equal(t, E.Right[error]("postgres:secret"), result)
}

func TestChainFirstLeft_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := F.Pipe1(
			Right[OuterConfig, InnerConfig, error](42),
			ChainFirstLeft[int](sideEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "side effect should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestChainFirstLeft_Failure(t *testing.T) {
	t.Run("executes on Left value and preserves original error", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		var capturedErr error
		var capturedOuter OuterConfig
		var capturedInner InnerConfig

		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			return func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, F.Void] {
				capturedErr = e
				capturedOuter = r
				return func(c InnerConfig) IOE.IOEither[error, F.Void] {
					capturedInner = c
					return IOE.Of[error](F.VOID)
				}
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			ChainFirstLeft[int](sideEffect),
		)

		outer := OuterConfig{database: "test", logLevel: "debug"}
		inner := InnerConfig{apiKey: "key", timeout: 10}

		// Act
		result := computation(outer)(inner)()

		// Assert
		assert.Equal(t, originalErr, capturedErr, "should capture original error")
		assert.Equal(t, outer, capturedOuter, "should pass outer context")
		assert.Equal(t, inner, capturedInner, "should pass inner context")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})

	t.Run("preserves original error even if side effect returns different error", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		sideEffectErr := errors.New("side effect error")

		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			return Left[OuterConfig, InnerConfig, F.Void](sideEffectErr)
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			ChainFirstLeft[int](sideEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error, not side effect error")
	})
}

func TestMonadChainFirstLeft_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := Right[OuterConfig, InnerConfig, error](42)

		// Act
		result := MonadChainFirstLeft(computation, sideEffect)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "side effect should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestMonadChainFirstLeft_Failure(t *testing.T) {
	t.Run("executes on Left value and preserves original error", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		var capturedErr error

		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			capturedErr = e
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := Left[OuterConfig, InnerConfig, int](originalErr)

		// Act
		result := MonadChainFirstLeft(computation, sideEffect)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, originalErr, capturedErr, "should capture original error")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestTapLeft_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		tap := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := F.Pipe1(
			Right[OuterConfig, InnerConfig, error](42),
			TapLeft[int](tap),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "tap should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestTapLeft_Failure(t *testing.T) {
	t.Run("executes tap on Left value for logging", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("validation failed")
		var loggedErr error
		var logLevel string

		logError := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			return func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, F.Void] {
				loggedErr = e
				logLevel = r.logLevel
				return RIOE.Of[InnerConfig, error](F.VOID)
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeft[int](logError),
		)

		// Act
		result := computation(OuterConfig{logLevel: "ERROR"})(InnerConfig{})()

		// Assert
		assert.Equal(t, originalErr, loggedErr, "should log the error")
		assert.Equal(t, "ERROR", logLevel, "should pass log level")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})

	t.Run("chains multiple taps", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("network error")
		var tap1Executed, tap2Executed bool

		tap1 := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			tap1Executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		tap2 := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			tap2Executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := F.Pipe2(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeft[int](tap1),
			TapLeft[int](tap2),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.True(t, tap1Executed, "first tap should execute")
		assert.True(t, tap2Executed, "second tap should execute")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestMonadTapLeft_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		tap := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			executed = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := Right[OuterConfig, InnerConfig, error](42)

		// Act
		result := MonadTapLeft(computation, tap)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "tap should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestMonadTapLeft_Failure(t *testing.T) {
	t.Run("executes tap on Left value for metrics", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("database error")
		errorCount := 0

		recordMetric := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			errorCount++
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := Left[OuterConfig, InnerConfig, int](originalErr)

		// Act
		result := MonadTapLeft(computation, recordMetric)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, 1, errorCount, "should increment error count")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestChainFirstLeft_EdgeCases(t *testing.T) {
	t.Run("handles nil error in side effect", func(t *testing.T) {
		// Arrange
		var capturedErr error

		sideEffect := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			capturedErr = e
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		var nilErr error
		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](nilErr),
			ChainFirstLeft[int](sideEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Nil(t, capturedErr, "should handle nil error")
		assert.Equal(t, E.Left[int](nilErr), result)
	})
}

func TestTapLeft_Integration(t *testing.T) {
	t.Run("integrates with other operations", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("parse error")
		var logged bool

		logError := func(e error) ReaderReaderIOEither[OuterConfig, InnerConfig, error, F.Void] {
			logged = true
			return Right[OuterConfig, InnerConfig, error](F.VOID)
		}

		computation := F.Pipe3(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeft[int](logError),
			MapLeft[OuterConfig, InnerConfig, int](func(e error) error {
				return errors.New("wrapped: " + e.Error())
			}),
			Map[OuterConfig, InnerConfig, error](func(x int) int { return x * 2 }),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.True(t, logged, "should execute tap before MapLeft")
		err := F.Pipe1(result, E.Fold(
			F.Identity[error],
			func(_ int) error { t.Fatal("expected Left"); return nil },
		))
		assert.Equal(t, "wrapped: parse error", err.Error(), "should apply MapLeft transformation")
	})
}

func TestChainFirstLeftIOK_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		ioEffect := func(e error) io.IO[F.Void] {
			return func() F.Void {
				executed = true
				return F.VOID
			}
		}

		computation := F.Pipe1(
			Right[OuterConfig, InnerConfig, error](42),
			ChainFirstLeftIOK[int, OuterConfig, InnerConfig](ioEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "IO effect should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestChainFirstLeftIOK_Failure(t *testing.T) {
	t.Run("executes IO on Left value and preserves original error", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		var capturedErr error

		ioEffect := func(e error) io.IO[F.Void] {
			return func() F.Void {
				capturedErr = e
				return F.VOID
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			ChainFirstLeftIOK[int, OuterConfig, InnerConfig](ioEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, originalErr, capturedErr, "should capture original error in IO")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})

	t.Run("preserves original error even if IO returns different value", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		ioCounter := 0

		ioEffect := func(e error) io.IO[int] {
			return func() int {
				ioCounter++
				return 999
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			ChainFirstLeftIOK[int, OuterConfig, InnerConfig](ioEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, 1, ioCounter, "IO should execute")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error, not IO result")
	})
}

func TestMonadChainFirstLeftIOK_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		ioEffect := func(e error) io.IO[F.Void] {
			return func() F.Void {
				executed = true
				return F.VOID
			}
		}

		computation := Right[OuterConfig, InnerConfig, error](42)

		// Act
		result := MonadChainFirstLeftIOK(computation, ioEffect)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "IO effect should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestMonadChainFirstLeftIOK_Failure(t *testing.T) {
	t.Run("executes IO on Left value and preserves original error", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("original error")
		var capturedErr error

		ioEffect := func(e error) io.IO[F.Void] {
			return func() F.Void {
				capturedErr = e
				return F.VOID
			}
		}

		computation := Left[OuterConfig, InnerConfig, int](originalErr)

		// Act
		result := MonadChainFirstLeftIOK(computation, ioEffect)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, originalErr, capturedErr, "should capture original error in IO")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestTapLeftIOK_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		ioTap := func(e error) io.IO[F.Void] {
			return func() F.Void {
				executed = true
				return F.VOID
			}
		}

		computation := F.Pipe1(
			Right[OuterConfig, InnerConfig, error](42),
			TapLeftIOK[int, OuterConfig, InnerConfig](ioTap),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "IO tap should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestTapLeftIOK_Failure(t *testing.T) {
	t.Run("executes IO tap on Left value for logging", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("validation failed")
		var loggedErr error

		logToIO := func(e error) io.IO[F.Void] {
			return func() F.Void {
				loggedErr = e
				return F.VOID
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeftIOK[int, OuterConfig, InnerConfig](logToIO),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, originalErr, loggedErr, "should log the error via IO")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})

	t.Run("chains multiple IO taps", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("network error")
		var tap1Executed, tap2Executed bool

		ioTap1 := func(e error) io.IO[F.Void] {
			return func() F.Void {
				tap1Executed = true
				return F.VOID
			}
		}

		ioTap2 := func(e error) io.IO[F.Void] {
			return func() F.Void {
				tap2Executed = true
				return F.VOID
			}
		}

		computation := F.Pipe2(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeftIOK[int, OuterConfig, InnerConfig](ioTap1),
			TapLeftIOK[int, OuterConfig, InnerConfig](ioTap2),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.True(t, tap1Executed, "first IO tap should execute")
		assert.True(t, tap2Executed, "second IO tap should execute")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})

	t.Run("IO tap for metrics collection", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("database error")
		errorCount := 0

		recordMetric := func(e error) io.IO[F.Void] {
			return func() F.Void {
				errorCount++
				return F.VOID
			}
		}

		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeftIOK[int, OuterConfig, InnerConfig](recordMetric),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, 1, errorCount, "should increment error count via IO")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestMonadTapLeftIOK_Success(t *testing.T) {
	t.Run("does not execute on Right value", func(t *testing.T) {
		// Arrange
		executed := false
		ioTap := func(e error) io.IO[F.Void] {
			return func() F.Void {
				executed = true
				return F.VOID
			}
		}

		computation := Right[OuterConfig, InnerConfig, error](42)

		// Act
		result := MonadTapLeftIOK(computation, ioTap)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.False(t, executed, "IO tap should not execute on Right")
		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestMonadTapLeftIOK_Failure(t *testing.T) {
	t.Run("executes IO tap on Left value for metrics", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("database error")
		errorCount := 0

		recordMetric := func(e error) io.IO[F.Void] {
			return func() F.Void {
				errorCount++
				return F.VOID
			}
		}

		computation := Left[OuterConfig, InnerConfig, int](originalErr)

		// Act
		result := MonadTapLeftIOK(computation, recordMetric)(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Equal(t, 1, errorCount, "should increment error count via IO")
		assert.Equal(t, E.Left[int](originalErr), result, "should preserve original error")
	})
}

func TestTapLeftIOK_Integration(t *testing.T) {
	t.Run("integrates with other operations", func(t *testing.T) {
		// Arrange
		originalErr := errors.New("parse error")
		var logged bool

		logError := func(e error) io.IO[F.Void] {
			return func() F.Void {
				logged = true
				return F.VOID
			}
		}

		computation := F.Pipe3(
			Left[OuterConfig, InnerConfig, int](originalErr),
			TapLeftIOK[int, OuterConfig, InnerConfig](logError),
			MapLeft[OuterConfig, InnerConfig, int](func(e error) error {
				return errors.New("wrapped: " + e.Error())
			}),
			Map[OuterConfig, InnerConfig, error](func(x int) int { return x * 2 }),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.True(t, logged, "should execute IO tap before MapLeft")
		err := F.Pipe1(result, E.Fold(
			F.Identity[error],
			func(_ int) error { t.Fatal("expected Left"); return nil },
		))
		assert.Equal(t, "wrapped: parse error", err.Error(), "should apply MapLeft transformation")
	})
}

func TestChainFirstLeftIOK_EdgeCases(t *testing.T) {
	t.Run("handles nil error in IO effect", func(t *testing.T) {
		// Arrange
		var capturedErr error

		ioEffect := func(e error) io.IO[F.Void] {
			return func() F.Void {
				capturedErr = e
				return F.VOID
			}
		}

		var nilErr error
		computation := F.Pipe1(
			Left[OuterConfig, InnerConfig, int](nilErr),
			ChainFirstLeftIOK[int, OuterConfig, InnerConfig](ioEffect),
		)

		// Act
		result := computation(OuterConfig{})(InnerConfig{})()

		// Assert
		assert.Nil(t, capturedErr, "should handle nil error in IO")
		assert.Equal(t, E.Left[int](nilErr), result)
	})
}

func TestDefer_Success(t *testing.T) {
	t.Run("evaluates generator on each execution", func(t *testing.T) {
		// Arrange: count how many times the generator is invoked
		callCount := 0
		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			callCount++
			return Of[OuterConfig, InnerConfig, error](42)
		})

		// Generator must NOT be called until the computation is executed
		assert.Equal(t, 0, callCount, "generator must not be called before execution")

		// First execution
		assert.Equal(t, E.Right[error](42), computation(OuterConfig{})(InnerConfig{})())
		assert.Equal(t, 1, callCount)

		// Second execution re-invokes the generator
		assert.Equal(t, E.Right[error](42), computation(OuterConfig{})(InnerConfig{})())
		assert.Equal(t, 2, callCount)
	})

	t.Run("generator receives outer and inner environments", func(t *testing.T) {
		outer := OuterConfig{database: "postgres", logLevel: "info"}
		inner := InnerConfig{apiKey: "secret", timeout: 30}

		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			return func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
				return func(c InnerConfig) IOE.IOEither[error, string] {
					return IOE.Right[error](fmt.Sprintf("db=%s,key=%s", r.database, c.apiKey))
				}
			}
		})

		assert.Equal(t, E.Right[error]("db=postgres,key=secret"), computation(outer)(inner)())
	})
}

func TestDefer_Failure(t *testing.T) {
	t.Run("propagates Left produced by the deferred computation", func(t *testing.T) {
		expectedErr := errors.New("deferred error")
		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Left[OuterConfig, InnerConfig, int](expectedErr)
		})

		assert.Equal(t, E.Left[int](expectedErr), computation(OuterConfig{})(InnerConfig{})())
	})

	t.Run("generator can return different computations on each call", func(t *testing.T) {
		callCount := 0
		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			callCount++
			if callCount == 1 {
				return Left[OuterConfig, InnerConfig, int](errors.New("first call"))
			}
			return Of[OuterConfig, InnerConfig, error](callCount)
		})

		// First call: Left
		r1 := computation(OuterConfig{})(InnerConfig{})()
		assert.True(t, E.IsLeft(r1))
		assert.Equal(t, 1, callCount)

		// Second call: Right(2)
		assert.Equal(t, E.Right[error](2), computation(OuterConfig{})(InnerConfig{})())
		assert.Equal(t, 2, callCount)
	})
}

func TestDefer_EdgeCases(t *testing.T) {
	t.Run("deferred computation is lazy — generator not called at construction time", func(t *testing.T) {
		generated := false
		_ = Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			generated = true
			return Of[OuterConfig, InnerConfig, error](0)
		})
		assert.False(t, generated, "generator must not fire until the computation is executed")
	})

	t.Run("composes with Chain after deferral", func(t *testing.T) {
		computation := F.Pipe1(
			Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
				return Of[OuterConfig, InnerConfig, error](21)
			}),
			Chain(func(n int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
				return Of[OuterConfig, InnerConfig, error](n * 2)
			}),
		)

		assert.Equal(t, E.Right[error](42), computation(OuterConfig{})(InnerConfig{})())
	})

	t.Run("different outer environments produce different results", func(t *testing.T) {
		outer1 := OuterConfig{database: "primary"}
		outer2 := OuterConfig{database: "replica"}

		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
			return func(r OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
				return func(_ InnerConfig) IOE.IOEither[error, string] {
					return IOE.Right[error](r.database)
				}
			}
		})

		assert.Equal(t, E.Right[error]("primary"), computation(outer1)(InnerConfig{})())
		assert.Equal(t, E.Right[error]("replica"), computation(outer2)(InnerConfig{})())
	})
}

func TestDefer_Integration(t *testing.T) {
	t.Run("side effects captured at execution time, not construction time", func(t *testing.T) {
		var log []string

		computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			log = append(log, "generator called")
			return Of[OuterConfig, InnerConfig, error](1)
		})

		assert.Empty(t, log, "no side effects before execution")

		computation(OuterConfig{})(InnerConfig{})()
		assert.Equal(t, []string{"generator called"}, log)

		computation(OuterConfig{})(InnerConfig{})()
		assert.Equal(t, []string{"generator called", "generator called"}, log)
	})

	t.Run("deferred computation feeds into Map", func(t *testing.T) {
		computation := F.Pipe1(
			Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
				return Of[OuterConfig, InnerConfig, error](10)
			}),
			Map[OuterConfig, InnerConfig, error](func(n int) string {
				return fmt.Sprintf("value=%d", n)
			}),
		)

		assert.Equal(t, E.Right[error]("value=10"), computation(OuterConfig{})(InnerConfig{})())
	})
}

func ExampleDefer() {
	// Defer creates a ReaderReaderIOEither lazily: the generator function is
	// called every time the computation is executed, not when it is constructed.

	callCount := 0
	computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
		callCount++
		return Of[OuterConfig, InnerConfig, error](callCount * 10)
	})

	// Generator has not been called yet.
	fmt.Println("before execution:", callCount)

	r1 := computation(OuterConfig{})(InnerConfig{})()
	fmt.Println("after first execution:", E.IsRight(r1), callCount)

	r2 := computation(OuterConfig{})(InnerConfig{})()
	fmt.Println("after second execution:", E.IsRight(r2), callCount)

	// Output:
	// before execution: 0
	// after first execution: true 1
	// after second execution: true 2
}

func ExampleDefer_errorCase() {
	// Defer propagates the Left value produced by the deferred computation.

	expectedErr := errors.New("something failed")
	computation := Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
		return Left[OuterConfig, InnerConfig, int](expectedErr)
	})

	result := computation(OuterConfig{})(InnerConfig{})()
	fmt.Println(E.IsLeft(result))

	// Output:
	// true
}

func ExampleDefer_composition() {
	// Defer composes naturally with other operators in a pipeline.

	computation := F.Pipe1(
		Defer(func() ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Of[OuterConfig, InnerConfig, error](21)
		}),
		Chain(func(n int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, int] {
			return Of[OuterConfig, InnerConfig, error](n * 2)
		}),
	)

	result := computation(OuterConfig{})(InnerConfig{})()
	fmt.Println(E.IsRight(result))

	// Output:
	// true
}

// ---------------------------------------------------------------------------
// ReadIOEither
// ---------------------------------------------------------------------------

func TestReadIOEither_Success(t *testing.T) {
	t.Run("Right IOEither provides outer environment to computation", func(t *testing.T) {
		// Arrange: IOEither that successfully resolves the outer config
		configIO := IOE.Right[error](OuterConfig{database: "postgres", logLevel: "info"})

		// Computation reads the outer env and also the inner env
		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(inner InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](fmt.Sprintf("%s/%s", outer.database, inner.apiKey))
			}
		}

		// Act: collapse the two-layer stack into a single ReaderIOEither
		rioe := ReadIOEither[string, OuterConfig, InnerConfig, error](configIO)(computation)

		// Assert
		inner := InnerConfig{apiKey: "secret", timeout: 30}
		assert.Equal(t, E.Right[error]("postgres/secret"), rioe(inner)())
	})

	t.Run("inner environment is still threaded through after ReadIOEither", func(t *testing.T) {
		inner1 := InnerConfig{apiKey: "key-one", timeout: 10}
		inner2 := InnerConfig{apiKey: "key-two", timeout: 20}

		configIO := IOE.Right[error](OuterConfig{database: "db"})
		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(inner InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](fmt.Sprintf("%s:%s", outer.database, inner.apiKey))
			}
		}

		rioe := ReadIOEither[string, OuterConfig, InnerConfig, error](configIO)(computation)

		assert.Equal(t, E.Right[error]("db:key-one"), rioe(inner1)())
		assert.Equal(t, E.Right[error]("db:key-two"), rioe(inner2)())
	})

	t.Run("composes with Map in a pipeline before ReadIOEither", func(t *testing.T) {
		configIO := IOE.Right[error](OuterConfig{database: "pg"})

		result := F.Pipe2(
			Of[OuterConfig, InnerConfig, error](21),
			Map[OuterConfig, InnerConfig, error](N.Mul(2)),
			ReadIOEither[int, OuterConfig, InnerConfig, error](configIO),
		)(InnerConfig{})()

		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestReadIOEither_Failure(t *testing.T) {
	t.Run("Left IOEither short-circuits without executing the computation", func(t *testing.T) {
		// Arrange
		loadErr := errors.New("failed to load outer config")
		configIO := IOE.Left[OuterConfig](loadErr)

		executed := false
		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
			executed = true
			return RIOE.Of[InnerConfig, error](99)
		}

		// Act
		rioe := ReadIOEither[int, OuterConfig, InnerConfig, error](configIO)(computation)
		result := rioe(InnerConfig{})()

		// Assert
		assert.Equal(t, E.Left[int](loadErr), result)
		assert.False(t, executed, "computation must not run when IOEither is Left")
	})

	t.Run("Left IOEither error is preserved exactly", func(t *testing.T) {
		sentinel := errors.New("sentinel")
		rioe := ReadIOEither[int, OuterConfig, InnerConfig, error](IOE.Left[OuterConfig](sentinel))(
			func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
				return RIOE.Of[InnerConfig, error](0)
			},
		)
		got := rioe(InnerConfig{})()
		assert.Equal(t, E.Left[int](sentinel), got)
	})

	t.Run("Right IOEither but inner computation returns Left", func(t *testing.T) {
		computeErr := errors.New("compute error")
		configIO := IOE.Right[error](OuterConfig{database: "pg"})

		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
			return func(_ InnerConfig) IOE.IOEither[error, int] {
				return IOE.Left[int](computeErr)
			}
		}

		result := ReadIOEither[int, OuterConfig, InnerConfig, error](configIO)(computation)(InnerConfig{})()
		assert.Equal(t, E.Left[int](computeErr), result)
	})
}

func TestReadIOEither_EdgeCases(t *testing.T) {
	t.Run("IOEither is evaluated lazily on each invocation", func(t *testing.T) {
		callCount := 0
		configIO := func() E.Either[error, OuterConfig] {
			callCount++
			return E.Right[error](OuterConfig{database: fmt.Sprintf("db%d", callCount)})
		}

		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(_ InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](outer.database)
			}
		}

		rioe := ReadIOEither[string, OuterConfig, InnerConfig, error](configIO)(computation)

		assert.Equal(t, E.Right[error]("db1"), rioe(InnerConfig{})())
		assert.Equal(t, E.Right[error]("db2"), rioe(InnerConfig{})())
		assert.Equal(t, 2, callCount)
	})

	t.Run("result is itself a valid ReaderIOEither that can be re-run", func(t *testing.T) {
		configIO := IOE.Right[error](OuterConfig{database: "stable"})
		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return RIOE.Of[InnerConfig, error]("ok")
		}

		rioe := ReadIOEither[string, OuterConfig, InnerConfig, error](configIO)(computation)
		inner := InnerConfig{}

		assert.Equal(t, rioe(inner)(), rioe(inner)())
	})
}

func TestReadIOEither_Integration(t *testing.T) {
	t.Run("outer env from IOEither feeds a two-context computation", func(t *testing.T) {
		outer := OuterConfig{database: "prod", logLevel: "warn"}
		inner := InnerConfig{apiKey: "tok", timeout: 5}

		computation := Ask[OuterConfig, InnerConfig, error]()

		result := ReadIOEither[OuterConfig, OuterConfig, InnerConfig, error](IOE.Right[error](outer))(computation)(inner)()
		assert.Equal(t, E.Right[error](outer), result)
	})

	t.Run("ReadIOEither result composes with further RIOE operators", func(t *testing.T) {
		configIO := IOE.Right[error](OuterConfig{database: "pg"})
		base := Of[OuterConfig, InnerConfig, error](7)

		rioe := F.Pipe1(
			base,
			Map[OuterConfig, InnerConfig, error](N.Mul(6)),
		)

		result := ReadIOEither[int, OuterConfig, InnerConfig, error](configIO)(rioe)(InnerConfig{})()
		assert.Equal(t, E.Right[error](42), result)
	})
}

// ---------------------------------------------------------------------------
// ReadIO
// ---------------------------------------------------------------------------

func TestReadIO_Success(t *testing.T) {
	t.Run("IO provides outer environment to computation", func(t *testing.T) {
		configIO := func() OuterConfig {
			return OuterConfig{database: "postgres", logLevel: "info"}
		}

		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(inner InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](fmt.Sprintf("%s/%s", outer.database, inner.apiKey))
			}
		}

		rioe := ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)
		inner := InnerConfig{apiKey: "secret", timeout: 30}

		assert.Equal(t, E.Right[error]("postgres/secret"), rioe(inner)())
	})

	t.Run("inner environment is still threaded through after ReadIO", func(t *testing.T) {
		inner1 := InnerConfig{apiKey: "alpha", timeout: 1}
		inner2 := InnerConfig{apiKey: "beta", timeout: 2}

		configIO := func() OuterConfig { return OuterConfig{database: "db"} }
		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(inner InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](fmt.Sprintf("%s:%s", outer.database, inner.apiKey))
			}
		}

		rioe := ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)

		assert.Equal(t, E.Right[error]("db:alpha"), rioe(inner1)())
		assert.Equal(t, E.Right[error]("db:beta"), rioe(inner2)())
	})

	t.Run("IO is always successful — computation always receives the environment", func(t *testing.T) {
		// Unlike ReadIOEither, the IO cannot fail, so the computation is always reached
		executed := false
		configIO := func() OuterConfig { return OuterConfig{database: "always"} }
		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
			executed = true
			return RIOE.Of[InnerConfig, error](1)
		}

		ReadIO[InnerConfig, error, int, OuterConfig](configIO)(computation)(InnerConfig{})()
		assert.True(t, executed, "computation must always execute when using ReadIO")
	})

	t.Run("composes with Map in a pipeline before ReadIO", func(t *testing.T) {
		configIO := func() OuterConfig { return OuterConfig{database: "pg"} }

		result := F.Pipe2(
			Of[OuterConfig, InnerConfig, error](21),
			Map[OuterConfig, InnerConfig, error](N.Mul(2)),
			ReadIO[InnerConfig, error, int, OuterConfig](configIO),
		)(InnerConfig{})()

		assert.Equal(t, E.Right[error](42), result)
	})
}

func TestReadIO_Failure(t *testing.T) {
	t.Run("inner computation returning Left is propagated", func(t *testing.T) {
		computeErr := errors.New("compute error")
		configIO := func() OuterConfig { return OuterConfig{database: "pg"} }

		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
			return func(_ InnerConfig) IOE.IOEither[error, int] {
				return IOE.Left[int](computeErr)
			}
		}

		result := ReadIO[InnerConfig, error, int, OuterConfig](configIO)(computation)(InnerConfig{})()
		assert.Equal(t, E.Left[int](computeErr), result)
	})

	t.Run("error originates from inner env inspection", func(t *testing.T) {
		authErr := errors.New("unauthorized")
		configIO := func() OuterConfig { return OuterConfig{database: "pg"} }

		computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(inner InnerConfig) IOE.IOEither[error, string] {
				if inner.apiKey == "" {
					return IOE.Left[string](authErr)
				}
				return IOE.Right[error](inner.apiKey)
			}
		}

		// no apiKey → Left
		assert.Equal(t, E.Left[string](authErr),
			ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)(InnerConfig{})())

		// with apiKey → Right
		assert.Equal(t, E.Right[error]("tok"),
			ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)(InnerConfig{apiKey: "tok"})())
	})
}

func TestReadIO_EdgeCases(t *testing.T) {
	t.Run("IO is evaluated lazily on each invocation", func(t *testing.T) {
		callCount := 0
		configIO := func() OuterConfig {
			callCount++
			return OuterConfig{database: fmt.Sprintf("db%d", callCount)}
		}

		computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
			return func(_ InnerConfig) IOE.IOEither[error, string] {
				return IOE.Right[error](outer.database)
			}
		}

		rioe := ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)

		assert.Equal(t, E.Right[error]("db1"), rioe(InnerConfig{})())
		assert.Equal(t, E.Right[error]("db2"), rioe(InnerConfig{})())
		assert.Equal(t, 2, callCount)
	})

	t.Run("ReadIO vs ReadIOEither: ReadIO always reaches computation", func(t *testing.T) {
		// ReadIO  — IO cannot fail, so computation is always invoked
		// ReadIOEither — IOEither can be Left, bypassing computation entirely
		reachedViaIO := false
		reachedViaIOE := false

		ioFn := func() OuterConfig { return OuterConfig{} }
		ioEFn := IOE.Right[error](OuterConfig{})
		badIOE := IOE.Left[OuterConfig](errors.New("bad"))

		computation := func(reach *bool) func(OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
			return func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
				*reach = true
				return RIOE.Of[InnerConfig, error](1)
			}
		}

		ReadIO[InnerConfig, error, int, OuterConfig](ioFn)(computation(&reachedViaIO))(InnerConfig{})()
		ReadIOEither[int, OuterConfig, InnerConfig, error](ioEFn)(computation(&reachedViaIOE))(InnerConfig{})()

		assert.True(t, reachedViaIO)
		assert.True(t, reachedViaIOE)

		// With a Left IOEither the computation is bypassed
		reachedViaIOE = false
		ReadIOEither[int, OuterConfig, InnerConfig, error](badIOE)(computation(&reachedViaIOE))(InnerConfig{})()
		assert.False(t, reachedViaIOE)
	})
}

func TestReadIO_Integration(t *testing.T) {
	t.Run("outer env from IO feeds a two-context computation", func(t *testing.T) {
		outer := OuterConfig{database: "prod", logLevel: "warn"}
		inner := InnerConfig{apiKey: "tok", timeout: 5}

		computation := Ask[OuterConfig, InnerConfig, error]()

		result := ReadIO[InnerConfig, error, OuterConfig, OuterConfig](func() OuterConfig { return outer })(computation)(inner)()
		assert.Equal(t, E.Right[error](outer), result)
	})

	t.Run("ReadIO result composes with further RIOE operators", func(t *testing.T) {
		configIO := func() OuterConfig { return OuterConfig{database: "pg"} }
		base := F.Pipe1(
			Of[OuterConfig, InnerConfig, error](7),
			Map[OuterConfig, InnerConfig, error](N.Mul(6)),
		)

		result := ReadIO[InnerConfig, error, int, OuterConfig](configIO)(base)(InnerConfig{})()
		assert.Equal(t, E.Right[error](42), result)
	})
}

// ---------------------------------------------------------------------------
// Examples
// ---------------------------------------------------------------------------

func ExampleReadIOEither() {
	// ReadIOEither collapses a ReaderReaderIOEither into a ReaderIOEither by
	// supplying the outer environment from an IOEither.  When the IOEither is
	// Right the contained value is used as the outer env; when it is Left the
	// error propagates immediately and the inner computation is never called.

	configIO := IOE.Right[error](OuterConfig{database: "postgres"})

	computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
		return func(inner InnerConfig) IOE.IOEither[error, string] {
			return IOE.Right[error](fmt.Sprintf("%s/%s", outer.database, inner.apiKey))
		}
	}

	// Supply the outer env; still need to provide the inner env.
	rioe := ReadIOEither[string, OuterConfig, InnerConfig, error](configIO)(computation)
	result := rioe(InnerConfig{apiKey: "tok"})()
	fmt.Println(E.IsRight(result))

	// Output:
	// true
}

func ExampleReadIOEither_leftPropagation() {
	// When the IOEither is Left the error is forwarded and the computation
	// is never executed.

	loadErr := errors.New("config unavailable")
	configIO := IOE.Left[OuterConfig](loadErr)

	executed := false
	computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
		executed = true
		return RIOE.Of[InnerConfig, error](0)
	}

	rioe := ReadIOEither[int, OuterConfig, InnerConfig, error](configIO)(computation)
	result := rioe(InnerConfig{})()
	fmt.Println(E.IsLeft(result))
	fmt.Println(executed)

	// Output:
	// true
	// false
}

func ExampleReadIO() {
	// ReadIO collapses a ReaderReaderIOEither into a ReaderIOEither by
	// supplying the outer environment from a plain IO (which cannot fail).
	// The resulting ReaderIOEither still needs the inner environment C.

	configIO := func() OuterConfig {
		return OuterConfig{database: "postgres"}
	}

	computation := func(outer OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
		return func(inner InnerConfig) IOE.IOEither[error, string] {
			return IOE.Right[error](fmt.Sprintf("%s/%s", outer.database, inner.apiKey))
		}
	}

	rioe := ReadIO[InnerConfig, error, string, OuterConfig](configIO)(computation)
	result := rioe(InnerConfig{apiKey: "tok"})()
	fmt.Println(E.IsRight(result))

	// Output:
	// true
}

func ExampleReadIO_alwaysReachesComputation() {
	// Unlike ReadIOEither, ReadIO uses an infallible IO to load the outer env,
	// so the computation is always executed.

	executed := false
	configIO := func() OuterConfig { return OuterConfig{} }

	computation := func(_ OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, int] {
		executed = true
		return RIOE.Of[InnerConfig, error](42)
	}

	ReadIO[InnerConfig, error, int, OuterConfig](configIO)(computation)(InnerConfig{})()
	fmt.Println(executed)

	// Output:
	// true
}
