// Copyright (c) 2023 IBM Corp.
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

package generic

import (
	F "github.com/IBM/fp-go/function"
	N "github.com/IBM/fp-go/number"
	T "github.com/IBM/fp-go/tuple"
)

// ZipWith applies a function to pairs of elements at the same index in two arrays, collecting the results in a new array. If one
// input array is short, excess elements of the longer array are discarded.
func ZipWith[AS ~[]A, BS ~[]B, CS ~[]C, FCT ~func(A, B) C, A, B, C any](fa AS, fb BS, f FCT) CS {
	l := N.Min(len(fa), len(fb))
	res := make(CS, l)
	for i := l - 1; i >= 0; i-- {
		res[i] = f(fa[i], fb[i])
	}
	return res
}

// Zip takes two arrays and returns an array of corresponding pairs. If one input array is short, excess elements of the
// longer array are discarded
func Zip[AS ~[]A, BS ~[]B, CS ~[]T.Tuple2[A, B], A, B any](fb BS) func(AS) CS {
	return F.Bind23of3(ZipWith[AS, BS, CS, func(A, B) T.Tuple2[A, B]])(fb, T.MakeTuple2[A, B])
}

// Unzip is the function is reverse of [Zip]. Takes an array of pairs and return two corresponding arrays
func Unzip[AS ~[]A, BS ~[]B, CS ~[]T.Tuple2[A, B], A, B any](cs CS) T.Tuple2[AS, BS] {
	l := len(cs)
	as := make(AS, l)
	bs := make(BS, l)
	for i := l - 1; i >= 0; i-- {
		t := cs[i]
		as[i] = t.F1
		bs[i] = t.F2
	}
	return T.MakeTuple2(as, bs)
}
