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
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOR "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	N "github.com/IBM/fp-go/v2/number"
)

type benchConfig struct {
	value int
}

var (
	benchErr    = errors.New("benchmark error")
	benchCfg    = benchConfig{value: 100}
	benchResult Result[int]
	benchRIOE   ReaderIOResult[benchConfig, int]
	benchInt    int
)

// Benchmark core constructors
func BenchmarkLeft(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Left[benchConfig, int](benchErr)
	}
}

func BenchmarkRight(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Right[benchConfig](42)
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Of[benchConfig](42)
	}
}

func BenchmarkFromEither_Right(b *testing.B) {
	either := E.Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromEither[benchConfig](either)
	}
}

func BenchmarkFromEither_Left(b *testing.B) {
	either := E.Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromEither[benchConfig](either)
	}
}

func BenchmarkFromIO(b *testing.B) {
	io := func() int { return 42 }
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIO[benchConfig, error](io)
	}
}

func BenchmarkFromIOResult_Right(b *testing.B) {
	ioe := IOR.Of(42)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIOResult[benchConfig](ioe)
	}
}

func BenchmarkFromIOResult_Left(b *testing.B) {
	ioe := IOR.Left[int](benchErr)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIOResult[benchConfig](ioe)
	}
}

// Benchmark execution
func BenchmarkExecute_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := rioe(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}

func BenchmarkExecute_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := rioe(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}

// Benchmark functor operations
func BenchmarkMonadMap_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	mapper := N.Mul(2)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadMap(rioe, mapper)
	}
}

func BenchmarkMonadMap_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	mapper := N.Mul(2)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadMap(rioe, mapper)
	}
}

func BenchmarkMap_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	mapper := Map[benchConfig](N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

func BenchmarkMap_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	mapper := Map[benchConfig](N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

func BenchmarkMapTo_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	mapper := MapTo[benchConfig, int](99)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

// Benchmark monad operations
func BenchmarkMonadChain_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	chainer := func(a int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadChain(rioe, chainer)
	}
}

func BenchmarkMonadChain_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	chainer := func(a int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadChain(rioe, chainer)
	}
}

func BenchmarkChain_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	chainer := Chain(func(a int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChain_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	chainer := Chain(func(a int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainFirst_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	chainer := ChainFirst(func(a int) ReaderIOResult[benchConfig, string] { return Right[benchConfig]("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainFirst_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	chainer := ChainFirst(func(a int) ReaderIOResult[benchConfig, string] { return Right[benchConfig]("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkFlatten_Right(b *testing.B) {
	nested := Right[benchConfig](Right[benchConfig](42))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Flatten(nested)
	}
}

func BenchmarkFlatten_Left(b *testing.B) {
	nested := Left[benchConfig, ReaderIOResult[benchConfig, int]](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Flatten(nested)
	}
}

// Benchmark applicative operations
func BenchmarkMonadApSeq_RightRight(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Right[benchConfig](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApSeq_RightLeft(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApSeq_LeftRight(b *testing.B) {
	fab := Left[benchConfig, func(int) int](benchErr)
	fa := Right[benchConfig](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApPar_RightRight(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Right[benchConfig](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

func BenchmarkMonadApPar_RightLeft(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

func BenchmarkMonadApPar_LeftRight(b *testing.B) {
	fab := Left[benchConfig, func(int) int](benchErr)
	fa := Right[benchConfig](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

// Benchmark execution of applicative operations
func BenchmarkExecuteApSeq_RightRight(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Right[benchConfig](42)
	rioe := MonadApSeq(fab, fa)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := rioe(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}

func BenchmarkExecuteApPar_RightRight(b *testing.B) {
	fab := Right[benchConfig](N.Mul(2))
	fa := Right[benchConfig](42)
	rioe := MonadApPar(fab, fa)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := rioe(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}

// Benchmark chain operations with different types
func BenchmarkChainEitherK_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	chainer := ChainEitherK[benchConfig](func(a int) E.Either[error, int] { return E.Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainEitherK_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	chainer := ChainEitherK[benchConfig](func(a int) E.Either[error, int] { return E.Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOK_Right(b *testing.B) {
	rioe := Right[benchConfig](42)
	chainer := ChainIOK[benchConfig](func(a int) func() int { return func() int { return a * 2 } })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOK_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	chainer := ChainIOK[benchConfig](func(a int) func() int { return func() int { return a * 2 } })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

// Benchmark context operations
func BenchmarkAsk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Ask[benchConfig]()
	}
}

func BenchmarkAsks(b *testing.B) {
	reader := func(cfg benchConfig) int { return cfg.value }
	b.ReportAllocs()
	for b.Loop() {
		_ = Asks(reader)
	}
}

// Benchmark pipeline operations
func BenchmarkPipeline_Map_Right(b *testing.B) {
	rioe := Right[benchConfig](21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Map[benchConfig](N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Map_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Map[benchConfig](N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Chain_Right(b *testing.B) {
	rioe := Right[benchConfig](21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Chain(func(x int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](x * 2) }),
		)
	}
}

func BenchmarkPipeline_Chain_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Chain(func(x int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](x * 2) }),
		)
	}
}

func BenchmarkPipeline_Complex_Right(b *testing.B) {
	rioe := Right[benchConfig](10)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe3(
			rioe,
			Map[benchConfig](N.Mul(2)),
			Chain(func(x int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](x + 1) }),
			Map[benchConfig](N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Complex_Left(b *testing.B) {
	rioe := Left[benchConfig, int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe3(
			rioe,
			Map[benchConfig](N.Mul(2)),
			Chain(func(x int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](x + 1) }),
			Map[benchConfig](N.Mul(2)),
		)
	}
}

func BenchmarkExecutePipeline_Complex_Right(b *testing.B) {
	rioe := F.Pipe3(
		Right[benchConfig](10),
		Map[benchConfig](N.Mul(2)),
		Chain(func(x int) ReaderIOResult[benchConfig, int] { return Right[benchConfig](x + 1) }),
		Map[benchConfig](N.Mul(2)),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := rioe(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}

// Benchmark Local operation
func BenchmarkLocal(b *testing.B) {
	rioe := Asks(func(cfg benchConfig) int { return cfg.value })
	localOp := Local[int](func(cfg benchConfig) benchConfig { return benchConfig{value: cfg.value * 2} })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = localOp(rioe)
	}
}

func BenchmarkExecuteLocal(b *testing.B) {
	rioe := Asks(func(cfg benchConfig) int { return cfg.value })
	localOp := Local[int](func(cfg benchConfig) benchConfig { return benchConfig{value: cfg.value * 2} })
	modified := localOp(rioe)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		val, err := modified(benchCfg)()
		if err != nil {
			benchResult = E.Left[int](err)
		} else {
			benchResult = E.Right[error](val)
		}
	}
}
