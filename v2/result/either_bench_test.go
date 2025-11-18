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

package result

import (
	"errors"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
)

var (
	errBench       = errors.New("benchmark error")
	benchResultInt Result[int]
	benchBool      bool
	benchInt       int
)

// Benchmark core constructors
func BenchmarkLeft(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = Left[int](errBench)
	}
}

func BenchmarkRight(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = Right(42)
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = Of(42)
	}
}

// Benchmark predicates
func BenchmarkIsLeft(b *testing.B) {
	val := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchBool = IsLeft(val)
	}
}

func BenchmarkIsRight(b *testing.B) {
	val := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchBool = IsRight(val)
	}
}

// Benchmark fold operations
func BenchmarkFold_Right(b *testing.B) {
	val := Right(42)
	folder := Fold(
		func(e error) int { return 0 },
		N.Mul(2),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = folder(val)
	}
}

func BenchmarkFold_Left(b *testing.B) {
	val := Left[int](errBench)
	folder := Fold(
		func(e error) int { return 0 },
		N.Mul(2),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = folder(val)
	}
}

// Benchmark functor operations
func BenchmarkMap_Right(b *testing.B) {
	val := Right(42)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = mapper(val)
	}
}

func BenchmarkMap_Left(b *testing.B) {
	val := Left[int](errBench)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = mapper(val)
	}
}

// Benchmark monad operations
func BenchmarkChain_Right(b *testing.B) {
	val := Right(42)
	chainer := Chain(func(a int) Result[int] { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = chainer(val)
	}
}

func BenchmarkChain_Left(b *testing.B) {
	val := Left[int](errBench)
	chainer := Chain(func(a int) Result[int] { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = chainer(val)
	}
}

func BenchmarkChainFirst_Right(b *testing.B) {
	val := Right(42)
	chainer := ChainFirst(func(a int) Result[string] { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = chainer(val)
	}
}

func BenchmarkChainFirst_Left(b *testing.B) {
	val := Left[int](errBench)
	chainer := ChainFirst(func(a int) Result[string] { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = chainer(val)
	}
}

// Benchmark alternative operations
func BenchmarkAlt_RightRight(b *testing.B) {
	val := Right(42)
	alternative := Alt(func() Result[int] { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = alternative(val)
	}
}

func BenchmarkAlt_LeftRight(b *testing.B) {
	val := Left[int](errBench)
	alternative := Alt(func() Result[int] { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = alternative(val)
	}
}

func BenchmarkOrElse_Right(b *testing.B) {
	val := Right(42)
	recover := OrElse(func(e error) Result[int] { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = recover(val)
	}
}

func BenchmarkOrElse_Left(b *testing.B) {
	val := Left[int](errBench)
	recover := OrElse(func(e error) Result[int] { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt = recover(val)
	}
}

// Benchmark GetOrElse
func BenchmarkGetOrElse_Right(b *testing.B) {
	val := Right(42)
	getter := GetOrElse(func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = getter(val)
	}
}

func BenchmarkGetOrElse_Left(b *testing.B) {
	val := Left[int](errBench)
	getter := GetOrElse(func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = getter(val)
	}
}
