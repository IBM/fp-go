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

package either

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
)

var (
	errBench    = errors.New("benchmark error")
	benchResult Either[error, int]
	benchBool   bool
	benchInt    int
	benchString string
)

// Benchmark core constructors
func BenchmarkLeft(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = Left[int](errBench)
	}
}

func BenchmarkRight(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = Right[error](42)
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = Of[error](42)
	}
}

// Benchmark predicates
func BenchmarkIsLeft(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchBool = IsLeft(left)
	}
}

func BenchmarkIsRight(b *testing.B) {
	right := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchBool = IsRight(right)
	}
}

// Benchmark fold operations
func BenchmarkMonadFold_Right(b *testing.B) {
	right := Right[error](42)
	onLeft := func(e error) int { return 0 }
	onRight := func(a int) int { return a * 2 }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = MonadFold(right, onLeft, onRight)
	}
}

func BenchmarkMonadFold_Left(b *testing.B) {
	left := Left[int](errBench)
	onLeft := func(e error) int { return 0 }
	onRight := func(a int) int { return a * 2 }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = MonadFold(left, onLeft, onRight)
	}
}

func BenchmarkFold_Right(b *testing.B) {
	right := Right[error](42)
	folder := Fold(
		func(e error) int { return 0 },
		func(a int) int { return a * 2 },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = folder(right)
	}
}

func BenchmarkFold_Left(b *testing.B) {
	left := Left[int](errBench)
	folder := Fold(
		func(e error) int { return 0 },
		func(a int) int { return a * 2 },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = folder(left)
	}
}

// Benchmark unwrap operations
func BenchmarkUnwrap_Right(b *testing.B) {
	right := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt, _ = Unwrap(right)
	}
}

func BenchmarkUnwrap_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt, _ = Unwrap(left)
	}
}

func BenchmarkUnwrapError_Right(b *testing.B) {
	right := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt, _ = UnwrapError(right)
	}
}

func BenchmarkUnwrapError_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt, _ = UnwrapError(left)
	}
}

// Benchmark functor operations
func BenchmarkMonadMap_Right(b *testing.B) {
	right := Right[error](42)
	mapper := func(a int) int { return a * 2 }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadMap(right, mapper)
	}
}

func BenchmarkMonadMap_Left(b *testing.B) {
	left := Left[int](errBench)
	mapper := func(a int) int { return a * 2 }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadMap(left, mapper)
	}
}

func BenchmarkMap_Right(b *testing.B) {
	right := Right[error](42)
	mapper := Map[error](func(a int) int { return a * 2 })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = mapper(right)
	}
}

func BenchmarkMap_Left(b *testing.B) {
	left := Left[int](errBench)
	mapper := Map[error](func(a int) int { return a * 2 })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = mapper(left)
	}
}

func BenchmarkMapLeft_Right(b *testing.B) {
	right := Right[error](42)
	mapper := MapLeft[int](func(e error) string { return e.Error() })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mapper(right)
	}
}

func BenchmarkMapLeft_Left(b *testing.B) {
	left := Left[int](errBench)
	mapper := MapLeft[int](func(e error) string { return e.Error() })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mapper(left)
	}
}

func BenchmarkBiMap_Right(b *testing.B) {
	right := Right[error](42)
	mapper := BiMap(
		func(e error) string { return e.Error() },
		func(a int) string { return "value" },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mapper(right)
	}
}

func BenchmarkBiMap_Left(b *testing.B) {
	left := Left[int](errBench)
	mapper := BiMap(
		func(e error) string { return e.Error() },
		func(a int) string { return "value" },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mapper(left)
	}
}

// Benchmark monad operations
func BenchmarkMonadChain_Right(b *testing.B) {
	right := Right[error](42)
	chainer := func(a int) Either[error, int] { return Right[error](a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadChain(right, chainer)
	}
}

func BenchmarkMonadChain_Left(b *testing.B) {
	left := Left[int](errBench)
	chainer := func(a int) Either[error, int] { return Right[error](a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadChain(left, chainer)
	}
}

func BenchmarkChain_Right(b *testing.B) {
	right := Right[error](42)
	chainer := Chain[error](func(a int) Either[error, int] { return Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = chainer(right)
	}
}

func BenchmarkChain_Left(b *testing.B) {
	left := Left[int](errBench)
	chainer := Chain[error](func(a int) Either[error, int] { return Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = chainer(left)
	}
}

func BenchmarkChainFirst_Right(b *testing.B) {
	right := Right[error](42)
	chainer := ChainFirst(func(a int) Either[error, string] { return Right[error]("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = chainer(right)
	}
}

func BenchmarkChainFirst_Left(b *testing.B) {
	left := Left[int](errBench)
	chainer := ChainFirst[error](func(a int) Either[error, string] { return Right[error]("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = chainer(left)
	}
}

func BenchmarkFlatten_Right(b *testing.B) {
	nested := Right[error](Right[error](42))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = Flatten(nested)
	}
}

func BenchmarkFlatten_Left(b *testing.B) {
	nested := Left[Either[error, int]](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = Flatten(nested)
	}
}

// Benchmark applicative operations
func BenchmarkMonadAp_RightRight(b *testing.B) {
	fab := Right[error](func(a int) int { return a * 2 })
	fa := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadAp(fab, fa)
	}
}

func BenchmarkMonadAp_RightLeft(b *testing.B) {
	fab := Right[error](func(a int) int { return a * 2 })
	fa := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadAp(fab, fa)
	}
}

func BenchmarkMonadAp_LeftRight(b *testing.B) {
	fab := Left[func(int) int](errBench)
	fa := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadAp(fab, fa)
	}
}

func BenchmarkAp_RightRight(b *testing.B) {
	fab := Right[error](func(a int) int { return a * 2 })
	fa := Right[error](42)
	ap := Ap[int, error, int](fa)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = ap(fab)
	}
}

// Benchmark alternative operations
func BenchmarkAlt_RightRight(b *testing.B) {
	right := Right[error](42)
	alternative := Alt(func() Either[error, int] { return Right[error](99) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = alternative(right)
	}
}

func BenchmarkAlt_LeftRight(b *testing.B) {
	left := Left[int](errBench)
	alternative := Alt[error](func() Either[error, int] { return Right[error](99) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = alternative(left)
	}
}

func BenchmarkOrElse_Right(b *testing.B) {
	right := Right[error](42)
	recover := OrElse[error](func(e error) Either[error, int] { return Right[error](0) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = recover(right)
	}
}

func BenchmarkOrElse_Left(b *testing.B) {
	left := Left[int](errBench)
	recover := OrElse(func(e error) Either[error, int] { return Right[error](0) })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = recover(left)
	}
}

// Benchmark conversion operations
func BenchmarkTryCatch_Success(b *testing.B) {
	onThrow := func(err error) error { return err }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = TryCatch(42, nil, onThrow)
	}
}

func BenchmarkTryCatch_Error(b *testing.B) {
	onThrow := func(err error) error { return err }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = TryCatch(0, errBench, onThrow)
	}
}

func BenchmarkTryCatchError_Success(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = TryCatchError(42, nil)
	}
}

func BenchmarkTryCatchError_Error(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = TryCatchError(0, errBench)
	}
}

func BenchmarkSwap_Right(b *testing.B) {
	right := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Swap(right)
	}
}

func BenchmarkSwap_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Swap(left)
	}
}

func BenchmarkGetOrElse_Right(b *testing.B) {
	right := Right[error](42)
	getter := GetOrElse[error](func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = getter(right)
	}
}

func BenchmarkGetOrElse_Left(b *testing.B) {
	left := Left[int](errBench)
	getter := GetOrElse[error](func(e error) int { return 0 })
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchInt = getter(left)
	}
}

// Benchmark pipeline operations
func BenchmarkPipeline_Map_Right(b *testing.B) {
	right := Right[error](21)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe1(
			right,
			Map[error](func(x int) int { return x * 2 }),
		)
	}
}

func BenchmarkPipeline_Map_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe1(
			left,
			Map[error](func(x int) int { return x * 2 }),
		)
	}
}

func BenchmarkPipeline_Chain_Right(b *testing.B) {
	right := Right[error](21)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe1(
			right,
			Chain[error](func(x int) Either[error, int] { return Right[error](x * 2) }),
		)
	}
}

func BenchmarkPipeline_Chain_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe1(
			left,
			Chain[error](func(x int) Either[error, int] { return Right[error](x * 2) }),
		)
	}
}

func BenchmarkPipeline_Complex_Right(b *testing.B) {
	right := Right[error](10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe3(
			right,
			Map[error](func(x int) int { return x * 2 }),
			Chain[error](func(x int) Either[error, int] { return Right[error](x + 1) }),
			Map[error](func(x int) int { return x * 2 }),
		)
	}
}

func BenchmarkPipeline_Complex_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = F.Pipe3(
			left,
			Map[error](func(x int) int { return x * 2 }),
			Chain[error](func(x int) Either[error, int] { return Right[error](x + 1) }),
			Map[error](func(x int) int { return x * 2 }),
		)
	}
}

// Benchmark sequence operations
func BenchmarkMonadSequence2_RightRight(b *testing.B) {
	e1 := Right[error](10)
	e2 := Right[error](20)
	f := func(a, b int) Either[error, int] { return Right[error](a + b) }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadSequence2(e1, e2, f)
	}
}

func BenchmarkMonadSequence2_LeftRight(b *testing.B) {
	e1 := Left[int](errBench)
	e2 := Right[error](20)
	f := func(a, b int) Either[error, int] { return Right[error](a + b) }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadSequence2(e1, e2, f)
	}
}

func BenchmarkMonadSequence3_RightRightRight(b *testing.B) {
	e1 := Right[error](10)
	e2 := Right[error](20)
	e3 := Right[error](30)
	f := func(a, b, c int) Either[error, int] { return Right[error](a + b + c) }
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult = MonadSequence3(e1, e2, e3, f)
	}
}

// Benchmark do-notation operations
func BenchmarkDo(b *testing.B) {
	type State struct{ value int }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Do[error](State{})
	}
}

func BenchmarkBind_Right(b *testing.B) {
	type State struct{ value int }
	initial := Do[error](State{})
	binder := Bind[error, State, State](
		func(v int) func(State) State {
			return func(s State) State { return State{value: v} }
		},
		func(s State) Either[error, int] {
			return Right[error](42)
		},
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = binder(initial)
	}
}

func BenchmarkLet_Right(b *testing.B) {
	type State struct{ value int }
	initial := Right[error](State{value: 10})
	letter := Let[error, State, State](
		func(v int) func(State) State {
			return func(s State) State { return State{value: s.value + v} }
		},
		func(s State) int { return 32 },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = letter(initial)
	}
}

// Benchmark string formatting
func BenchmarkString_Right(b *testing.B) {
	right := Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchString = right.String()
	}
}

func BenchmarkString_Left(b *testing.B) {
	left := Left[int](errBench)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchString = left.String()
	}
}

// Made with Bob
