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
		_, _ = Some(42)
	}
}

func BenchmarkNone(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = None[int]()
	}
}

// Benchmark basic operations
func BenchmarkIsSome(b *testing.B) {
	v, ok := Some(42)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsSome(v, ok)
	}
}

func BenchmarkMap(b *testing.B) {
	v, ok := Some(21)
	mapper := Map(N.Mul(2))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mapper(v, ok)
	}
}

func BenchmarkChain(b *testing.B) {
	v, ok := Some(21)
	chainer := Chain(func(x int) (int, bool) {
		if x > 0 {
			return x * 2, true
		}
		return 0, false
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chainer(v, ok)
	}
}

func BenchmarkFilter(b *testing.B) {
	v, ok := Some(42)
	filter := Filter(N.MoreThan(0))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = filter(v, ok)
	}
}

func BenchmarkGetOrElse(b *testing.B) {
	v, ok := Some(42)
	getter := GetOrElse(func() int { return 0 })
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getter(v, ok)
	}
}

// Benchmark collection operations
func BenchmarkTraverseArray_Small(b *testing.B) {
	data := []int{1, 2, 3, 4, 5}
	traverser := TraverseArray(func(x int) (int, bool) {
		return x * 2, true
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = traverser(data)
	}
}

func BenchmarkTraverseArray_Large(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	traverser := TraverseArray(func(x int) (int, bool) {
		return x * 2, true
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = traverser(data)
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
		s1, ok1 := Do(State{})
		s2, ok2 := Bind(
			func(x int) func(State) State {
				return func(s State) State {
					s.x = x
					return s
				}
			},
			func(s State) (int, bool) { return 10, true },
		)(s1, ok1)
		_, _ = Bind(
			func(y int) func(State) State {
				return func(s State) State {
					s.y = y
					return s
				}
			},
			func(s State) (int, bool) { return 20, true },
		)(s2, ok2)
	}
}

// Benchmark conversions
func BenchmarkFromPredicate(b *testing.B) {
	pred := FromPredicate(N.MoreThan(0))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pred(42)
	}
}

func BenchmarkFromNillable(b *testing.B) {
	val := 42
	ptr := &val
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromNillable(ptr)
	}
}

// Benchmark complex chains
func BenchmarkComplexChain(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1, ok1 := Some(1)
		v2, ok2 := Chain(func(x int) (int, bool) { return x + 1, true })(v1, ok1)
		v3, ok3 := Chain(func(x int) (int, bool) { return x * 2, true })(v2, ok2)
		_, _ = Chain(func(x int) (int, bool) { return x - 5, true })(v3, ok3)
	}
}
