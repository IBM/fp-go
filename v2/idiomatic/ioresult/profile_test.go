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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// Benchmark the closure allocations in Bind
func BenchmarkBindAllocations(b *testing.B) {
	type Data struct {
		Value int
	}

	setter := func(v int) func(Data) Data {
		return func(d Data) Data {
			d.Value = v
			return d
		}
	}

	b.Run("Bind", func(b *testing.B) {
		io := Of(Data{Value: 0})
		f := func(d Data) IOResult[int] { return Of(d.Value * 2) }

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := Bind(setter, f)(io)
			_, _ = result()
		}
	})

	b.Run("DirectChainMap", func(b *testing.B) {
		io := Of(Data{Value: 0})
		f := func(d Data) IOResult[int] { return Of(d.Value * 2) }

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Manual inlined version of Bind to see baseline
			result := Chain(func(s1 Data) IOResult[Data] {
				return Map(func(b int) Data {
					return setter(b)(s1)
				})(f(s1))
			})(io)
			_, _ = result()
		}
	})
}

// Benchmark Map with different patterns
func BenchmarkMapPatterns(b *testing.B) {
	b.Run("SimpleFunction", func(b *testing.B) {
		io := Of(42)
		f := N.Mul(2)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := Map(f)(io)
			_, _ = result()
		}
	})

	b.Run("InlinedLambda", func(b *testing.B) {
		io := Of(42)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := Map(N.Mul(2))(io)
			_, _ = result()
		}
	})

	b.Run("NestedMaps", func(b *testing.B) {
		io := Of(42)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := F.Pipe3(
				io,
				Map(N.Add(1)),
				Map(N.Mul(2)),
				Map(N.Sub(3)),
			)
			_, _ = result()
		}
	})
}

// Benchmark Of patterns
func BenchmarkOfPatterns(b *testing.B) {
	b.Run("IntValue", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			io := Of(42)
			_, _ = io()
		}
	})

	b.Run("StructValue", func(b *testing.B) {
		type Data struct {
			A int
			B string
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			io := Of(Data{A: 42, B: "test"})
			_, _ = io()
		}
	})

	b.Run("PointerValue", func(b *testing.B) {
		type Data struct {
			A int
			B string
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			io := Of(&Data{A: 42, B: "test"})
			_, _ = io()
		}
	})
}

// Benchmark the internal chain implementation
func BenchmarkChainPatterns(b *testing.B) {
	b.Run("SimpleChain", func(b *testing.B) {
		io := Of(42)
		f := func(x int) IOResult[int] { return Of(x * 2) }

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := Chain(f)(io)
			_, _ = result()
		}
	})

	b.Run("ChainSequence", func(b *testing.B) {
		f1 := func(x int) IOResult[int] { return Of(x + 1) }
		f2 := func(x int) IOResult[int] { return Of(x * 2) }
		f3 := func(x int) IOResult[int] { return Of(x - 3) }

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := F.Pipe3(
				Of(42),
				Chain(f1),
				Chain(f2),
				Chain(f3),
			)
			_, _ = result()
		}
	})
}

// Benchmark error handling paths
func BenchmarkErrorPaths(b *testing.B) {
	b.Run("SuccessPath", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := F.Pipe2(
				Of(42),
				Map(N.Mul(2)),
				Chain(func(x int) IOResult[int] { return Of(x + 1) }),
			)
			_, _ = result()
		}
	})

	b.Run("ErrorPath", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			result := F.Pipe2(
				Left[int](errors.New("error")),
				Map(N.Mul(2)),
				Chain(func(x int) IOResult[int] { return Of(x + 1) }),
			)
			_, _ = result()
		}
	})
}
