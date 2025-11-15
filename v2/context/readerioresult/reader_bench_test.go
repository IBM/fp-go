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
	"context"
	"errors"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
)

var (
	benchErr    = errors.New("benchmark error")
	benchCtx    = context.Background()
	benchResult Either[int]
	benchRIOE   ReaderIOResult[int]
	benchInt    int
)

// Benchmark core constructors
func BenchmarkLeft(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Left[int](benchErr)
	}
}

func BenchmarkRight(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Right(42)
	}
}

func BenchmarkOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Of(42)
	}
}

func BenchmarkFromEither_Right(b *testing.B) {
	either := E.Right[error](42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromEither(either)
	}
}

func BenchmarkFromEither_Left(b *testing.B) {
	either := E.Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromEither(either)
	}
}

func BenchmarkFromIO(b *testing.B) {
	io := func() int { return 42 }
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIO(io)
	}
}

func BenchmarkFromIOEither_Right(b *testing.B) {
	ioe := IOE.Of[error](42)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIOEither(ioe)
	}
}

func BenchmarkFromIOEither_Left(b *testing.B) {
	ioe := IOE.Left[int](benchErr)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = FromIOEither(ioe)
	}
}

// Benchmark execution
func BenchmarkExecute_Right(b *testing.B) {
	rioe := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

func BenchmarkExecute_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

func BenchmarkExecute_WithContext(b *testing.B) {
	rioe := Right(42)
	ctx, cancel := context.WithCancel(benchCtx)
	defer cancel()
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(ctx)()
	}
}

// Benchmark functor operations
func BenchmarkMonadMap_Right(b *testing.B) {
	rioe := Right(42)
	mapper := N.Mul(2)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadMap(rioe, mapper)
	}
}

func BenchmarkMonadMap_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	mapper := N.Mul(2)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadMap(rioe, mapper)
	}
}

func BenchmarkMap_Right(b *testing.B) {
	rioe := Right(42)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

func BenchmarkMap_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	mapper := Map(N.Mul(2))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

func BenchmarkMapTo_Right(b *testing.B) {
	rioe := Right(42)
	mapper := MapTo[int](99)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = mapper(rioe)
	}
}

// Benchmark monad operations
func BenchmarkMonadChain_Right(b *testing.B) {
	rioe := Right(42)
	chainer := func(a int) ReaderIOResult[int] { return Right(a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadChain(rioe, chainer)
	}
}

func BenchmarkMonadChain_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := func(a int) ReaderIOResult[int] { return Right(a * 2) }
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadChain(rioe, chainer)
	}
}

func BenchmarkChain_Right(b *testing.B) {
	rioe := Right(42)
	chainer := Chain(func(a int) ReaderIOResult[int] { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChain_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := Chain(func(a int) ReaderIOResult[int] { return Right(a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainFirst_Right(b *testing.B) {
	rioe := Right(42)
	chainer := ChainFirst(func(a int) ReaderIOResult[string] { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainFirst_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := ChainFirst(func(a int) ReaderIOResult[string] { return Right("logged") })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkFlatten_Right(b *testing.B) {
	nested := Right(Right(42))
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Flatten(nested)
	}
}

func BenchmarkFlatten_Left(b *testing.B) {
	nested := Left[ReaderIOResult[int]](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Flatten(nested)
	}
}

// Benchmark applicative operations
func BenchmarkMonadApSeq_RightRight(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApSeq_RightLeft(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApSeq_LeftRight(b *testing.B) {
	fab := Left[func(int) int](benchErr)
	fa := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApSeq(fab, fa)
	}
}

func BenchmarkMonadApPar_RightRight(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

func BenchmarkMonadApPar_RightLeft(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

func BenchmarkMonadApPar_LeftRight(b *testing.B) {
	fab := Left[func(int) int](benchErr)
	fa := Right(42)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = MonadApPar(fab, fa)
	}
}

// Benchmark execution of applicative operations
func BenchmarkExecuteApSeq_RightRight(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Right(42)
	rioe := MonadApSeq(fab, fa)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

func BenchmarkExecuteApPar_RightRight(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Right(42)
	rioe := MonadApPar(fab, fa)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

// Benchmark alternative operations
func BenchmarkAlt_RightRight(b *testing.B) {
	rioe := Right(42)
	alternative := Alt(func() ReaderIOResult[int] { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = alternative(rioe)
	}
}

func BenchmarkAlt_LeftRight(b *testing.B) {
	rioe := Left[int](benchErr)
	alternative := Alt(func() ReaderIOResult[int] { return Right(99) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = alternative(rioe)
	}
}

func BenchmarkOrElse_Right(b *testing.B) {
	rioe := Right(42)
	recover := OrElse(func(e error) ReaderIOResult[int] { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = recover(rioe)
	}
}

func BenchmarkOrElse_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	recover := OrElse(func(e error) ReaderIOResult[int] { return Right(0) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = recover(rioe)
	}
}

// Benchmark chain operations with different types
func BenchmarkChainEitherK_Right(b *testing.B) {
	rioe := Right(42)
	chainer := ChainEitherK(func(a int) Either[int] { return E.Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainEitherK_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := ChainEitherK(func(a int) Either[int] { return E.Right[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOK_Right(b *testing.B) {
	rioe := Right(42)
	chainer := ChainIOK(func(a int) func() int { return func() int { return a * 2 } })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOK_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := ChainIOK(func(a int) func() int { return func() int { return a * 2 } })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOEitherK_Right(b *testing.B) {
	rioe := Right(42)
	chainer := ChainIOEitherK(func(a int) IOEither[int] { return IOE.Of[error](a * 2) })
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = chainer(rioe)
	}
}

func BenchmarkChainIOEitherK_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	chainer := ChainIOEitherK(func(a int) IOEither[int] { return IOE.Of[error](a * 2) })
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
		_ = Ask()
	}
}

func BenchmarkDefer(b *testing.B) {
	gen := func() ReaderIOResult[int] { return Right(42) }
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Defer(gen)
	}
}

func BenchmarkMemoize(b *testing.B) {
	rioe := Right(42)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Memoize(rioe)
	}
}

// Benchmark delay operations
func BenchmarkDelay_Construction(b *testing.B) {
	rioe := Right(42)
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = Delay[int](time.Millisecond)(rioe)
	}
}

func BenchmarkTimer_Construction(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Timer(time.Millisecond)
	}
}

// Benchmark TryCatch
func BenchmarkTryCatch_Success(b *testing.B) {
	f := func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return 42, nil }
	}
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = TryCatch(f)
	}
}

func BenchmarkTryCatch_Error(b *testing.B) {
	f := func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return 0, benchErr }
	}
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = TryCatch(f)
	}
}

func BenchmarkExecuteTryCatch_Success(b *testing.B) {
	f := func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return 42, nil }
	}
	rioe := TryCatch(f)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

func BenchmarkExecuteTryCatch_Error(b *testing.B) {
	f := func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return 0, benchErr }
	}
	rioe := TryCatch(f)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

// Benchmark pipeline operations
func BenchmarkPipeline_Map_Right(b *testing.B) {
	rioe := Right(21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Map_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Chain_Right(b *testing.B) {
	rioe := Right(21)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Chain(func(x int) ReaderIOResult[int] { return Right(x * 2) }),
		)
	}
}

func BenchmarkPipeline_Chain_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe1(
			rioe,
			Chain(func(x int) ReaderIOResult[int] { return Right(x * 2) }),
		)
	}
}

func BenchmarkPipeline_Complex_Right(b *testing.B) {
	rioe := Right(10)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe3(
			rioe,
			Map(N.Mul(2)),
			Chain(func(x int) ReaderIOResult[int] { return Right(x + 1) }),
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkPipeline_Complex_Left(b *testing.B) {
	rioe := Left[int](benchErr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchRIOE = F.Pipe3(
			rioe,
			Map(N.Mul(2)),
			Chain(func(x int) ReaderIOResult[int] { return Right(x + 1) }),
			Map(N.Mul(2)),
		)
	}
}

func BenchmarkExecutePipeline_Complex_Right(b *testing.B) {
	rioe := F.Pipe3(
		Right(10),
		Map(N.Mul(2)),
		Chain(func(x int) ReaderIOResult[int] { return Right(x + 1) }),
		Map(N.Mul(2)),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

// Benchmark do-notation operations
func BenchmarkDo(b *testing.B) {
	type State struct{ value int }
	b.ReportAllocs()
	for b.Loop() {
		_ = Do(State{})
	}
}

func BenchmarkBind_Right(b *testing.B) {
	type State struct{ value int }
	initial := Do(State{})
	binder := Bind(
		func(v int) func(State) State {
			return func(s State) State { return State{value: v} }
		},
		func(s State) ReaderIOResult[int] {
			return Right(42)
		},
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = binder(initial)
	}
}

func BenchmarkLet_Right(b *testing.B) {
	type State struct{ value int }
	initial := Right(State{value: 10})
	letter := Let(
		func(v int) func(State) State {
			return func(s State) State { return State{value: s.value + v} }
		},
		func(s State) int { return 32 },
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = letter(initial)
	}
}

func BenchmarkApS_Right(b *testing.B) {
	type State struct{ value int }
	initial := Right(State{value: 10})
	aps := ApS(
		func(v int) func(State) State {
			return func(s State) State { return State{value: v} }
		},
		Right(42),
	)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = aps(initial)
	}
}

// Benchmark traverse operations
func BenchmarkTraverseArray_Empty(b *testing.B) {
	arr := []int{}
	traverser := TraverseArray(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(arr)
	}
}

func BenchmarkTraverseArray_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	traverser := TraverseArray(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(arr)
	}
}

func BenchmarkTraverseArray_Medium(b *testing.B) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}
	traverser := TraverseArray(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(arr)
	}
}

func BenchmarkTraverseArraySeq_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	traverser := TraverseArraySeq(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(arr)
	}
}

func BenchmarkTraverseArrayPar_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	traverser := TraverseArrayPar(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(arr)
	}
}

func BenchmarkSequenceArray_Small(b *testing.B) {
	arr := []ReaderIOResult[int]{
		Right(1),
		Right(2),
		Right(3),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = SequenceArray(arr)
	}
}

func BenchmarkExecuteTraverseArray_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	rioe := TraverseArray(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})(arr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rioe(benchCtx)()
	}
}

func BenchmarkExecuteTraverseArraySeq_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	rioe := TraverseArraySeq(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})(arr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rioe(benchCtx)()
	}
}

func BenchmarkExecuteTraverseArrayPar_Small(b *testing.B) {
	arr := []int{1, 2, 3}
	rioe := TraverseArrayPar(func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})(arr)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rioe(benchCtx)()
	}
}

// Benchmark record operations
func BenchmarkTraverseRecord_Small(b *testing.B) {
	rec := map[string]int{"a": 1, "b": 2, "c": 3}
	traverser := TraverseRecord[string](func(x int) ReaderIOResult[int] {
		return Right(x * 2)
	})
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = traverser(rec)
	}
}

func BenchmarkSequenceRecord_Small(b *testing.B) {
	rec := map[string]ReaderIOResult[int]{
		"a": Right(1),
		"b": Right(2),
		"c": Right(3),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = SequenceRecord(rec)
	}
}

// Benchmark resource management
func BenchmarkWithResource_Success(b *testing.B) {
	acquire := Right(42)
	release := func(int) ReaderIOResult[int] { return Right(0) }
	body := func(x int) ReaderIOResult[int] { return Right(x * 2) }

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = WithResource[int](acquire, release)(body)
	}
}

func BenchmarkExecuteWithResource_Success(b *testing.B) {
	acquire := Right(42)
	release := func(int) ReaderIOResult[int] { return Right(0) }
	body := func(x int) ReaderIOResult[int] { return Right(x * 2) }
	rioe := WithResource[int](acquire, release)(body)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

func BenchmarkExecuteWithResource_ErrorInBody(b *testing.B) {
	acquire := Right(42)
	release := func(int) ReaderIOResult[int] { return Right(0) }
	body := func(x int) ReaderIOResult[int] { return Left[int](benchErr) }
	rioe := WithResource[int](acquire, release)(body)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(benchCtx)()
	}
}

// Benchmark context cancellation
func BenchmarkExecute_CanceledContext(b *testing.B) {
	rioe := Right(42)
	ctx, cancel := context.WithCancel(benchCtx)
	cancel() // Cancel immediately

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(ctx)()
	}
}

func BenchmarkExecuteApPar_CanceledContext(b *testing.B) {
	fab := Right(N.Mul(2))
	fa := Right(42)
	rioe := MonadApPar(fab, fa)
	ctx, cancel := context.WithCancel(benchCtx)
	cancel() // Cancel immediately

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		benchResult = rioe(ctx)()
	}
}
