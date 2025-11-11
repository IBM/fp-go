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

package benchmarks

import (
	"math/big"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
)

var (
	createStringSet  = createRandom(createRandomString(256))(256)
	createIntDataSet = createRandom(randInt(10000))(256)

	globalResult any
)

func BenchmarkMap(b *testing.B) {

	data := createStringSet()

	var benchResult []string

	b.Run("functional", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResult = F.Pipe1(
				data,
				A.Map(strings.ToUpper),
			)
		}
	})

	b.Run("idiomatic", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var result = make([]string, 0, len(data))
			for _, value := range data {
				result = append(result, strings.ToUpper(value))
			}
			benchResult = result
		}
	})

	globalResult = benchResult
}

func isEven(data int) bool {
	return data%2 == 0
}

func isPrime(data int) bool {
	return big.NewInt(int64(data)).ProbablyPrime(0)
}

func BenchmarkMapThenFilter(b *testing.B) {

	data := createIntDataSet()
	var benchResult []int

	b.Run("functional isPrime", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResult = F.Pipe2(
				data,
				A.Filter(isPrime),
				A.Map(N.Div(2)),
			)
		}
	})

	b.Run("idiomatic isPrime", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var result []int
			for _, value := range data {
				if isPrime(value) {
					result = append(result, value/2)
				}
			}
			benchResult = result
		}
	})
	b.Run("functional isEven", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResult = F.Pipe2(
				data,
				A.Filter(isEven),
				A.Map(N.Div(2)),
			)
		}
	})

	b.Run("idiomatic isEven", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var result []int
			for _, value := range data {
				if isEven(value) {
					result = append(result, value/2)
				}
			}
			benchResult = result
		}
	})

	globalResult = benchResult
}

func BenchmarkFilterMap(b *testing.B) {

	data := createIntDataSet()
	var benchResult []int

	b.Run("functional isPrime", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResult = F.Pipe1(
				data,
				A.FilterMap(F.Flow2(
					O.FromPredicate(isPrime),
					O.Map(N.Div(2)),
				)),
			)
		}
	})

	b.Run("idiomatic isPrime", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var result []int
			for _, value := range data {
				if isPrime(value) {
					result = append(result, value/2)
				}
			}
			benchResult = result
		}
	})

	b.Run("functional isEven", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResult = F.Pipe1(
				data,
				A.FilterMap(F.Flow2(
					O.FromPredicate(isEven),
					O.Map(N.Div(2)),
				)),
			)
		}
	})

	b.Run("idiomatic isEven", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var result []int
			for _, value := range data {
				if isEven(value) {
					result = append(result, value/2)
				}
			}
			benchResult = result
		}
	})

	globalResult = benchResult
}
