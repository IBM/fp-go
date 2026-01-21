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
		Chain[OuterConfig, InnerConfig, error](func(v int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
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
		result := FromEither[OuterConfig, InnerConfig, error, int](E.Left[int](err))
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
		result := FromIOEither[OuterConfig, InnerConfig, error](ioe)
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
	isPositive := func(n int) bool { return n > 0 }
	onFalse := func(n int) error { return fmt.Errorf("not positive: %d", n) }

	t.Run("Predicate true", func(t *testing.T) {
		result := FromPredicate[OuterConfig, InnerConfig](isPositive, onFalse)(5)
		assert.Equal(t, E.Right[error](5), result(OuterConfig{})(InnerConfig{})())
	})

	t.Run("Predicate false", func(t *testing.T) {
		result := FromPredicate[OuterConfig, InnerConfig, error](isPositive, onFalse)(-5)
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
		ChainEitherK[OuterConfig, InnerConfig, error](func(v int) E.Either[error, string] {
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
