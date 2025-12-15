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
	S "github.com/IBM/fp-go/v2/string"
)

// Benchmark shallow chain (1 step)
func BenchmarkChain_1Step(b *testing.B) {
	v, ok := Some(1)
	chainer := Chain(func(x int) (int, bool) { return x + 1, true })
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chainer(v, ok)
	}
}

// Benchmark moderate chain (3 steps)
func BenchmarkChain_3Steps(b *testing.B) {
	v, ok := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1, ok1 := Chain(func(x int) (int, bool) { return x + 1, true })(v, ok)
		v2, ok2 := Chain(func(x int) (int, bool) { return x * 2, true })(v1, ok1)
		_, _ = Chain(func(x int) (int, bool) { return x - 5, true })(v2, ok2)
	}
}

// Benchmark deep chain (5 steps)
func BenchmarkChain_5Steps(b *testing.B) {
	v, ok := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1, ok1 := Chain(func(x int) (int, bool) { return x + 1, true })(v, ok)
		v2, ok2 := Chain(func(x int) (int, bool) { return x * 2, true })(v1, ok1)
		v3, ok3 := Chain(func(x int) (int, bool) { return x - 5, true })(v2, ok2)
		v4, ok4 := Chain(func(x int) (int, bool) { return x * 10, true })(v3, ok3)
		_, _ = Chain(func(x int) (int, bool) { return x + 100, true })(v4, ok4)
	}
}

// Benchmark very deep chain (10 steps)
func BenchmarkChain_10Steps(b *testing.B) {
	v, ok := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1, ok1 := Chain(func(x int) (int, bool) { return x + 1, true })(v, ok)
		v2, ok2 := Chain(func(x int) (int, bool) { return x * 2, true })(v1, ok1)
		v3, ok3 := Chain(func(x int) (int, bool) { return x - 5, true })(v2, ok2)
		v4, ok4 := Chain(func(x int) (int, bool) { return x * 10, true })(v3, ok3)
		v5, ok5 := Chain(func(x int) (int, bool) { return x + 100, true })(v4, ok4)
		v6, ok6 := Chain(func(x int) (int, bool) { return x - 50, true })(v5, ok5)
		v7, ok7 := Chain(func(x int) (int, bool) { return x * 3, true })(v6, ok6)
		v8, ok8 := Chain(func(x int) (int, bool) { return x + 20, true })(v7, ok7)
		v9, ok9 := Chain(func(x int) (int, bool) { return x / 2, true })(v8, ok8)
		_, _ = Chain(func(x int) (int, bool) { return x - 10, true })(v9, ok9)
	}
}

// Benchmark Map-based chain (should be faster due to inlining)
func BenchmarkMap_5Steps(b *testing.B) {
	v, ok := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1, ok1 := Map(N.Add(1))(v, ok)
		v2, ok2 := Map(N.Mul(3))(v1, ok1)
		v3, ok3 := Map(N.Add(20))(v2, ok2)
		v4, ok4 := Map(N.Div(2))(v3, ok3)
		_, _ = Map(N.Sub(10))(v4, ok4)
	}
}

// Real-world example: parsing and validating user input
func BenchmarkChain_RealWorld_Validation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, sok := Some("42")

		// Step 1: Validate not empty
		v1, ok1 := Chain(func(s string) (string, bool) {
			if S.IsNonEmpty(s) {
				return s, true
			}
			return "", false
		})(s, sok)

		// Step 2: Parse to int (simulated)
		v2, ok2 := Chain(func(s string) (int, bool) {
			if s == "42" {
				return 42, true
			}
			return 0, false
		})(v1, ok1)

		// Step 3: Validate range
		_, _ = Chain(func(n int) (int, bool) {
			if n > 0 && n < 100 {
				return n, true
			}
			return 0, false
		})(v2, ok2)
	}
}
