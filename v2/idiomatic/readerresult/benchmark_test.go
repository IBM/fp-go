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

package readerresult

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

type BenchContext struct {
	Value int
}

// Benchmark basic operations

func BenchmarkOf(b *testing.B) {
	ctx := BenchContext{Value: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := Of[BenchContext](i)
		_, _ = rr(ctx)
	}
}

func BenchmarkLeft(b *testing.B) {
	ctx := BenchContext{Value: 42}
	err := testError
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := Left[BenchContext, int](err)
		_, _ = rr(ctx)
	}
}

func BenchmarkMap(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr := Of[BenchContext](10)
	double := N.Mul(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := F.Pipe1(rr, Map[BenchContext](double))
		_, _ = mapped(ctx)
	}
}

func BenchmarkMapChain(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr := Of[BenchContext](1)
	double := N.Mul(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := F.Pipe3(
			rr,
			Map[BenchContext](double),
			Map[BenchContext](double),
			Map[BenchContext](double),
		)
		_, _ = result(ctx)
	}
}

func BenchmarkChain(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr := Of[BenchContext](10)
	addOne := func(x int) ReaderResult[BenchContext, int] {
		return Of[BenchContext](x + 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chained := F.Pipe1(rr, Chain(addOne))
		_, _ = chained(ctx)
	}
}

func BenchmarkChainDeep(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr := Of[BenchContext](1)
	addOne := func(x int) ReaderResult[BenchContext, int] {
		return Of[BenchContext](x + 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := F.Pipe5(
			rr,
			Chain(addOne),
			Chain(addOne),
			Chain(addOne),
			Chain(addOne),
			Chain(addOne),
		)
		_, _ = result(ctx)
	}
}

func BenchmarkAp(b *testing.B) {
	ctx := BenchContext{Value: 42}
	fab := Of[BenchContext](N.Mul(2))
	fa := Of[BenchContext](21)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := MonadAp(fab, fa)
		_, _ = result(ctx)
	}
}

func BenchmarkSequenceT2(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr1 := Of[BenchContext](10)
	rr2 := Of[BenchContext](20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := SequenceT2(rr1, rr2)
		_, _ = result(ctx)
	}
}

func BenchmarkSequenceT4(b *testing.B) {
	ctx := BenchContext{Value: 42}
	rr1 := Of[BenchContext](10)
	rr2 := Of[BenchContext](20)
	rr3 := Of[BenchContext](30)
	rr4 := Of[BenchContext](40)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := SequenceT4(rr1, rr2, rr3, rr4)
		_, _ = result(ctx)
	}
}

func BenchmarkDoNotation(b *testing.B) {
	ctx := context.Background()

	type State struct {
		A int
		B int
		C int
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := F.Pipe3(
			Do[context.Context](State{}),
			Bind(
				func(a int) func(State) State {
					return func(s State) State { s.A = a; return s }
				},
				func(s State) ReaderResult[context.Context, int] {
					return Of[context.Context](10)
				},
			),
			Bind(
				func(b int) func(State) State {
					return func(s State) State { s.B = b; return s }
				},
				func(s State) ReaderResult[context.Context, int] {
					return Of[context.Context](s.A * 2)
				},
			),
			Bind(
				func(c int) func(State) State {
					return func(s State) State { s.C = c; return s }
				},
				func(s State) ReaderResult[context.Context, int] {
					return Of[context.Context](s.A + s.B)
				},
			),
		)
		_, _ = result(ctx)
	}
}

func BenchmarkErrorPropagation(b *testing.B) {
	ctx := BenchContext{Value: 42}
	err := testError
	rr := Left[BenchContext, int](err)
	double := N.Mul(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := F.Pipe5(
			rr,
			Map[BenchContext](double),
			Map[BenchContext](double),
			Map[BenchContext](double),
			Map[BenchContext](double),
			Map[BenchContext](double),
		)
		_, _ = result(ctx)
	}
}

func BenchmarkTraverseArray(b *testing.B) {
	ctx := BenchContext{Value: 42}
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	kleisli := func(x int) ReaderResult[BenchContext, int] {
		return Of[BenchContext](x * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		traversed := TraverseArray(kleisli)
		result := traversed(arr)
		_, _ = result(ctx)
	}
}

func BenchmarkSequenceArray(b *testing.B) {
	ctx := BenchContext{Value: 42}
	arr := []ReaderResult[BenchContext, int]{
		Of[BenchContext](1),
		Of[BenchContext](2),
		Of[BenchContext](3),
		Of[BenchContext](4),
		Of[BenchContext](5),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := SequenceArray(arr)
		_, _ = result(ctx)
	}
}

// Real-world scenario benchmarks

func BenchmarkRealWorldPipeline(b *testing.B) {
	type Config struct {
		Multiplier int
		Offset     int
	}

	ctx := Config{Multiplier: 5, Offset: 10}

	type State struct {
		Input  int
		Result int
	}

	getMultiplier := func(cfg Config) int { return cfg.Multiplier }
	getOffset := func(cfg Config) int { return cfg.Offset }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		step1 := Bind(
			func(m int) func(State) State {
				return func(s State) State { s.Result = s.Input * m; return s }
			},
			func(s State) ReaderResult[Config, int] {
				return Asks(getMultiplier)
			},
		)
		step2 := Bind(
			func(off int) func(State) State {
				return func(s State) State { s.Result += off; return s }
			},
			func(s State) ReaderResult[Config, int] {
				return Asks(getOffset)
			},
		)
		result := F.Pipe3(
			Do[Config](State{Input: 10}),
			step1,
			step2,
			Map[Config](func(s State) int { return s.Result }),
		)
		_, _ = result(ctx)
	}
}
