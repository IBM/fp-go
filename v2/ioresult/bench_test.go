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
		_ = io()
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
		_ = io()
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
