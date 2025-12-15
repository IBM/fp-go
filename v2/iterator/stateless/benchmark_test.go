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

package stateless

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
)

func BenchmarkMulti(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		single()
	}
}

func single() int64 {

	length := 10000
	nums := make([]int, 0, length)
	for i := 0; i < length; i++ {
		nums = append(nums, i+1)
	}

	return F.Pipe6(
		nums,
		FromArray[int],
		Filter(func(n int) bool {
			return n%2 == 0
		}),
		Map(func(t int) int64 {
			return int64(t)
		}),
		Filter(func(t int64) bool {
			n := t
			for n/10 != 0 {
				if n%10 == 4 {
					return false
				}
				n /= 10
			}
			return true
		}),
		Map(func(t int64) int {
			return int(t)
		}),
		Reduce(func(n int64, r int) int64 {
			return n + int64(r)
		}, int64(0)),
	)
}
