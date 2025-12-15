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
	opt := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(opt, func(x int) Option[int] {
			return Some(x + 1)
		})
	}
}

// Benchmark moderate chain (3 steps)
func BenchmarkChain_3Steps(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(
			MonadChain(
				MonadChain(
					opt,
					func(x int) Option[int] { return Some(x + 1) },
				),
				func(x int) Option[int] { return Some(x * 2) },
			),
			func(x int) Option[int] { return Some(x - 5) },
		)
	}
}

// Benchmark deep chain (5 steps)
func BenchmarkChain_5Steps(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(
			MonadChain(
				MonadChain(
					MonadChain(
						MonadChain(
							opt,
							func(x int) Option[int] { return Some(x + 1) },
						),
						func(x int) Option[int] { return Some(x * 2) },
					),
					func(x int) Option[int] { return Some(x - 5) },
				),
				func(x int) Option[int] { return Some(x * 10) },
			),
			func(x int) Option[int] { return Some(x + 100) },
		)
	}
}

// Benchmark very deep chain (10 steps)
func BenchmarkChain_10Steps(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonadChain(
			MonadChain(
				MonadChain(
					MonadChain(
						MonadChain(
							MonadChain(
								MonadChain(
									MonadChain(
										MonadChain(
											MonadChain(
												opt,
												func(x int) Option[int] { return Some(x + 1) },
											),
											func(x int) Option[int] { return Some(x * 2) },
										),
										func(x int) Option[int] { return Some(x - 5) },
									),
									func(x int) Option[int] { return Some(x * 10) },
								),
								func(x int) Option[int] { return Some(x + 100) },
							),
							func(x int) Option[int] { return Some(x - 50) },
						),
						func(x int) Option[int] { return Some(x * 3) },
					),
					func(x int) Option[int] { return Some(x + 20) },
				),
				func(x int) Option[int] { return Some(x / 2) },
			),
			func(x int) Option[int] { return Some(x - 10) },
		)
	}
}

// Benchmark Map-based chain (should be faster due to inlining)
func BenchmarkMap_5Steps(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(N.Sub(10))(
			Map(N.Div(2))(
				Map(N.Add(20))(
					Map(N.Mul(3))(
						Map(N.Add(1))(opt),
					),
				),
			),
		)
	}
}

// Real-world example: parsing and validating user input
func BenchmarkChain_RealWorld_Validation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := Some("42")
		_ = MonadChain(
			MonadChain(
				MonadChain(
					input,
					// Step 1: Validate not empty
					func(s string) Option[string] {
						if S.IsNonEmpty(s) {
							return Some(s)
						}
						return None[string]()
					},
				),
				// Step 2: Parse to int (simulated)
				func(s string) Option[int] {
					// Simplified: just check if numeric
					if s == "42" {
						return Some(42)
					}
					return None[int]()
				},
			),
			// Step 3: Validate range
			func(n int) Option[int] {
				if n > 0 && n < 100 {
					return Some(n)
				}
				return None[int]()
			},
		)
	}
}
