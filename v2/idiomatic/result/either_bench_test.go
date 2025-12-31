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
	benchResultInt int
	benchResultErr error
	benchBool      bool
	benchInt       int
	benchString    string
)

// Benchmark core constructors
func BenchmarkLeft(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Left[int](errBench)
	}
}

func BenchmarkRight(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Right(42)
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Of(42)
	}
}

// Benchmark predicates
func BenchmarkIsLeft(b *testing.B) {
	val, err := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchBool = IsLeft(val, err)
	}
}

func BenchmarkIsRight(b *testing.B) {
	val, err := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchBool = IsRight(val, err)
	}
}

// Benchmark fold operations
func BenchmarkFold_Right(b *testing.B) {
	val, err := Right(42)
	folder := Fold(
		func(e error) int { return 0 },
		N.Mul(2),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = folder(val, err)
	}
}

func BenchmarkFold_Left(b *testing.B) {
	val, err := Left[int](errBench)
	folder := Fold(
		func(e error) int { return 0 },
		N.Mul(2),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = folder(val, err)
	}
}

// Benchmark functor operations
func BenchmarkMap_Right(b *testing.B) {
	val, err := Right(42)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

func BenchmarkMap_Left(b *testing.B) {
	val, err := Left[int](errBench)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

// Benchmark monad operations
func BenchmarkChain_Right(b *testing.B) {
	val, err := Right(42)
	chainer := Chain(func(a int) (int, error) { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

func BenchmarkChain_Left(b *testing.B) {
	val, err := Left[int](errBench)
	chainer := Chain(func(a int) (int, error) { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

func BenchmarkChainFirst_Right(b *testing.B) {
	val, err := Right(42)
	chainer := ChainFirst(func(a int) (string, error) { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

func BenchmarkChainFirst_Left(b *testing.B) {
	val, err := Left[int](errBench)
	chainer := ChainFirst(func(a int) (string, error) { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

// Benchmark alternative operations
func BenchmarkAlt_RightRight(b *testing.B) {
	val, err := Right(42)
	alternative := Alt(func() (int, error) { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = alternative(val, err)
	}
}

func BenchmarkAlt_LeftRight(b *testing.B) {
	val, err := Left[int](errBench)
	alternative := Alt(func() (int, error) { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = alternative(val, err)
	}
}

func BenchmarkOrElse_Right(b *testing.B) {
	val, err := Right(42)
	recover := OrElse(func(e error) (int, error) { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = recover(val, err)
	}
}

func BenchmarkOrElse_Left(b *testing.B) {
	val, err := Left[int](errBench)
	recover := OrElse(func(e error) (int, error) { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = recover(val, err)
	}
}

// Benchmark GetOrElse
func BenchmarkGetOrElse_Right(b *testing.B) {
	val, err := Right(42)
	getter := GetOrElse(func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = getter(val, err)
	}
}

func BenchmarkGetOrElse_Left(b *testing.B) {
	val, err := Left[int](errBench)
	getter := GetOrElse(func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = getter(val, err)
	}
}

// Benchmark pipeline operations
func BenchmarkPipeline_Map_Right(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe2(
			21,
			Right[int],
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Map_Left(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe2(
			0,
			func(int) (int, error) { return Left[int](errBench) },
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Chain_Right(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe2(
			21,
			Right[int],
			Chain(func(x int) (int, error) { return Right(x * 2) }),
		)
	}
}

func BenchmarkPipeline_Chain_Left(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe2(
			0,
			func(int) (int, error) { return Left[int](errBench) },
			Chain(func(x int) (int, error) { return Right(x * 2) }),
		)
	}
}

func BenchmarkPipeline_Complex_Right(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe4(
			10,
			Right[int],
			Map(N.Mul(2)),
			Chain(func(x int) (int, error) { return Right(x + 1) }),
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Complex_Left(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = Pipe4(
			0,
			func(int) (int, error) { return Left[int](errBench) },
			Map(N.Mul(2)),
			Chain(func(x int) (int, error) { return Right(x + 1) }),
			Map(N.Mul(2)),
		)
	}
}

// Benchmark string formatting
func BenchmarkToString_Right(b *testing.B) {
	val, err := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchString = ToString(val, err)
	}
}

func BenchmarkToString_Left(b *testing.B) {
	val, err := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchString = ToString(val, err)
	}
}

// Benchmark BiMap
func BenchmarkBiMap_Right(b *testing.B) {
	val, err := Right(42)
	wrapErr := func(e error) error { return e }
	mapper := BiMap(wrapErr, N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

func BenchmarkBiMap_Left(b *testing.B) {
	val, err := Left[int](errBench)
	wrapErr := func(e error) error { return e }
	mapper := BiMap(wrapErr, N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

// Benchmark MapTo
func BenchmarkMapTo_Right(b *testing.B) {
	val, err := Right(42)
	mapper := MapTo[int](99)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

func BenchmarkMapTo_Left(b *testing.B) {
	val, err := Left[int](errBench)
	mapper := MapTo[int](99)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

// Benchmark MapLeft
func BenchmarkMapLeft_Right(b *testing.B) {
	val, err := Right(42)
	mapper := MapLeft[int](func(e error) error { return e })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

func BenchmarkMapLeft_Left(b *testing.B) {
	val, err := Left[int](errBench)
	mapper := MapLeft[int](func(e error) error { return e })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = mapper(val, err)
	}
}

// Benchmark ChainTo
func BenchmarkChainTo_Right(b *testing.B) {
	val, err := Right(42)
	chainer := ChainTo[int](99, nil)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

func BenchmarkChainTo_Left(b *testing.B) {
	val, err := Left[int](errBench)
	chainer := ChainTo[int](99, nil)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = chainer(val, err)
	}
}

// Benchmark Reduce
func BenchmarkReduce_Right(b *testing.B) {
	val, err := Right(42)
	reducer := Reduce(func(acc, v int) int { return acc + v }, 10)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = reducer(val, err)
	}
}

func BenchmarkReduce_Left(b *testing.B) {
	val, err := Left[int](errBench)
	reducer := Reduce(func(acc, v int) int { return acc + v }, 10)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchInt = reducer(val, err)
	}
}

// Benchmark FromPredicate
func BenchmarkFromPredicate_Pass(b *testing.B) {
	pred := FromPredicate(
		N.MoreThan(0),
		func(x int) error { return errBench },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = pred(42)
	}
}

func BenchmarkFromPredicate_Fail(b *testing.B) {
	pred := FromPredicate(
		N.MoreThan(0),
		func(x int) error { return errBench },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = pred(-1)
	}
}

// Benchmark Flap
func BenchmarkFlap_Right(b *testing.B) {
	fn, ferr := Right(N.Mul(2))
	flapper := Flap[int](21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = flapper(fn, ferr)
	}
}

func BenchmarkFlap_Left(b *testing.B) {
	fn, ferr := Left[func(int) int](errBench)
	flapper := Flap[int](21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = flapper(fn, ferr)
	}
}

// Benchmark ToOption
func BenchmarkToOption_Right(b *testing.B) {
	val, err := Right(42)
	var resVal int
	var resOk bool
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		resVal, resOk = ToOption(val, err)
		benchInt = resVal
		benchBool = resOk
	}
}

func BenchmarkToOption_Left(b *testing.B) {
	val, err := Left[int](errBench)
	var resVal int
	var resOk bool
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		resVal, resOk = ToOption(val, err)
		benchInt = resVal
		benchBool = resOk
	}
}

// Benchmark FromOption
func BenchmarkFromOption_Some(b *testing.B) {
	converter := FromOption[int](func() error { return errBench })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = converter(42, true)
	}
}

func BenchmarkFromOption_None(b *testing.B) {
	converter := FromOption[int](func() error { return errBench })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultInt, benchResultErr = converter(0, false)
	}
}

// Benchmark ToError
func BenchmarkToError_Right(b *testing.B) {
	val, err := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultErr = ToError(val, err)
	}
}

func BenchmarkToError_Left(b *testing.B) {
	val, err := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResultErr = ToError(val, err)
	}
}

// Benchmark TraverseArray
func BenchmarkTraverseArray_Success(b *testing.B) {
	input := []int{1, 2, 3, 4, 5}
	traverse := TraverseArray(func(x int) (int, error) {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, benchResultErr = traverse(input)
	}
}

func BenchmarkTraverseArray_Error(b *testing.B) {
	input := []int{1, 2, 3, 4, 5}
	traverse := TraverseArray(func(x int) (int, error) {
		if x == 3 {
			return Left[int](errBench)
		}
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, benchResultErr = traverse(input)
	}
}
