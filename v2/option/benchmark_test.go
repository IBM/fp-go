// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"testing"

	N "github.com/IBM/fp-go/v2/number"
)

// Benchmark basic construction
func BenchmarkSome(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Some(42)
	}
}

func BenchmarkNone(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = None[int]()
	}
}

// Benchmark basic operations
func BenchmarkIsSome(b *testing.B) {
	opt := Some(42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsSome(opt)
	}
}

func BenchmarkMap(b *testing.B) {
	opt := Some(21)
	mapper := Map(N.Mul(2))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper(opt)
	}
}

func BenchmarkChain(b *testing.B) {
	opt := Some(21)
	chainer := Chain(func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chainer(opt)
	}
}

func BenchmarkFilter(b *testing.B) {
	opt := Some(42)
	filter := Filter(N.MoreThan(0))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(opt)
	}
}

func BenchmarkGetOrElse(b *testing.B) {
	opt := Some(42)
	getter := GetOrElse(func() int { return 0 })
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getter(opt)
	}
}

// Benchmark collection operations
func BenchmarkTraverseArray_Small(b *testing.B) {
	data := []int{1, 2, 3, 4, 5}
	traverser := TraverseArray(func(x int) Option[int] {
		return Some(x * 2)
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = traverser(data)
	}
}

func BenchmarkTraverseArray_Large(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	traverser := TraverseArray(func(x int) Option[int] {
		return Some(x * 2)
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = traverser(data)
	}
}

func BenchmarkSequenceArray_Small(b *testing.B) {
	data := []Option[int]{Some(1), Some(2), Some(3), Some(4), Some(5)}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SequenceArray(data)
	}
}

func BenchmarkCompactArray_Small(b *testing.B) {
	data := []Option[int]{Some(1), None[int](), Some(3), None[int](), Some(5)}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CompactArray(data)
	}
}

// Benchmark do-notation
func BenchmarkDoBind(b *testing.B) {
	type State struct {
		x int
		y int
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(
			MonadChain(
				Of(State{}),
				func(s State) Option[State] {
					s.x = 10
					return Some(s)
				},
			),
			func(s State) Option[State] {
				s.y = 20
				return Some(s)
			},
		)
	}
}

// Benchmark conversions
func BenchmarkFromPredicate(b *testing.B) {
	pred := FromPredicate(N.MoreThan(0))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pred(42)
	}
}

func BenchmarkFromNillable(b *testing.B) {
	val := 42
	ptr := &val
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromNillable(ptr)
	}
}

func BenchmarkTryCatch(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TryCatch(func() (int, error) {
			return 42, nil
		})
	}
}

// Benchmark complex chains
func BenchmarkComplexChain(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(
			MonadChain(
				MonadChain(
					Some(1),
					func(x int) Option[int] { return Some(x + 1) },
				),
				func(x int) Option[int] { return Some(x * 2) },
			),
			func(x int) Option[int] { return Some(x - 5) },
		)
	}
}
