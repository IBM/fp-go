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

package ioresult

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

func BenchmarkOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Of(42)
	}
}

func BenchmarkMap(b *testing.B) {
	io := Of(42)
	f := N.Mul(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(f)(io)
	}
}

func BenchmarkChain(b *testing.B) {
	io := Of(42)
	f := func(x int) IOResult[int] { return Of(x * 2) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Chain(f)(io)
	}
}

func BenchmarkBind(b *testing.B) {
	type Data struct {
		Value int
	}

	io := Of(Data{Value: 0})
	f := func(d Data) IOResult[int] { return Of(d.Value * 2) }
	setter := func(v int) func(Data) Data {
		return func(d Data) Data {
			d.Value = v
			return d
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Bind(setter, f)(io)
	}
}

func BenchmarkPipeline(b *testing.B) {
	f1 := N.Add(1)
	f2 := N.Mul(2)
	f3 := N.Sub(3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = F.Pipe3(
			Of(42),
			Map(f1),
			Map(f2),
			Map(f3),
		)
	}
}

func BenchmarkExecute(b *testing.B) {
	io := Of(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io()
	}
}

func BenchmarkExecutePipeline(b *testing.B) {
	f1 := N.Add(1)
	f2 := N.Mul(2)
	f3 := N.Sub(3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		io := F.Pipe3(
			Of(42),
			Map(f1),
			Map(f2),
			Map(f3),
		)
		_, _ = io()
	}
}

func BenchmarkChainSequence(b *testing.B) {
	f1 := func(x int) IOResult[int] { return Of(x + 1) }
	f2 := func(x int) IOResult[int] { return Of(x * 2) }
	f3 := func(x int) IOResult[int] { return Of(x - 3) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = F.Pipe3(
			Of(42),
			Chain(f1),
			Chain(f2),
			Chain(f3),
		)
	}
}

func BenchmarkLeft(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Left[int](F.Constant[error](nil)())
	}
}

func BenchmarkMapWithError(b *testing.B) {
	io := Left[int](F.Constant[error](nil)())
	f := N.Mul(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(f)(io)
	}
}

func BenchmarkChainWithError(b *testing.B) {
	io := Left[int](F.Constant[error](nil)())
	f := func(x int) IOResult[int] { return Of(x * 2) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Chain(f)(io)
	}
}

func BenchmarkApFirst(b *testing.B) {
	first := Of(42)
	second := Of("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ApFirst[int](second)(first)
	}
}

func BenchmarkApSecond(b *testing.B) {
	first := Of(42)
	second := Of("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ApSecond[int](second)(first)
	}
}

func BenchmarkMonadApFirst(b *testing.B) {
	first := Of(42)
	second := Of("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadApFirst(first, second)
	}
}

func BenchmarkMonadApSecond(b *testing.B) {
	first := Of(42)
	second := Of("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadApSecond(first, second)
	}
}

func BenchmarkFunctor(b *testing.B) {
	functor := Functor[int, int]()
	io := Of(42)
	f := N.Mul(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = functor.Map(f)(io)
	}
}

func BenchmarkMonad(b *testing.B) {
	monad := Monad[int, int]()
	f := func(x int) IOResult[int] { return Of(x * 2) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = monad.Chain(f)(monad.Of(42))
	}
}

func BenchmarkPointed(b *testing.B) {
	pointed := Pointed[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pointed.Of(42)
	}
}

func BenchmarkTraverseArray(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	f := func(x int) IOResult[int] { return Of(x * 2) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TraverseArray(f)(data)
	}
}

func BenchmarkSequenceArray(b *testing.B) {
	data := []IOResult[int]{Of(1), Of(2), Of(3), Of(4), Of(5)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SequenceArray(data)
	}
}

func BenchmarkAlt(b *testing.B) {
	first := Left[int](F.Constant[error](nil)())
	second := func() IOResult[int] { return Of(42) }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Alt(second)(first)
	}
}

func BenchmarkGetOrElse(b *testing.B) {
	io := Of(42)
	def := func(error) func() int { return func() int { return 0 } }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetOrElse(def)(io)()
	}
}

func BenchmarkFold(b *testing.B) {
	io := Of(42)
	onLeft := func(error) func() int { return func() int { return 0 } }
	onRight := func(x int) func() int { return func() int { return x } }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Fold(onLeft, onRight)(io)()
	}
}

func BenchmarkFromIO(b *testing.B) {
	ioVal := func() int { return 42 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromIO(ioVal)
	}
}

func BenchmarkChainIOK(b *testing.B) {
	io := Of(42)
	f := func(x int) func() int { return func() int { return x * 2 } }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ChainIOK(f)(io)
	}
}

func BenchmarkChainFirst(b *testing.B) {
	io := Of(42)
	f := func(x int) IOResult[string] { return Of("test") }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ChainFirst(f)(io)
	}
}
