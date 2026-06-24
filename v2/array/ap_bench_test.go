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

package array

import (
	"testing"
)

var apSink []int

// BenchmarkAp exercises the array applicative (cartesian product of a slice of
// functions over a slice of values).
func BenchmarkAp(b *testing.B) {
	fab := []func(int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x * 2 },
		func(x int) int { return x - 1 },
	}
	fa := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apSink = MonadAp[int, int](fab, fa)
	}
}
