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

package array

import (
	G "github.com/IBM/fp-go/array/generic"
	T "github.com/IBM/fp-go/tuple"
)

// ZipWith applies a function to pairs of elements at the same index in two arrays, collecting the results in a new array. If one
// input array is short, excess elements of the longer array are discarded.
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa []A, fb []B, f FCT) []C {
	return G.ZipWith[[]A, []B, []C, FCT](fa, fb, f)
}

// Zip takes two arrays and returns an array of corresponding pairs. If one input array is short, excess elements of the
// longer array are discarded
func Zip[A, B any](fb []B) func([]A) []T.Tuple2[A, B] {
	return G.Zip[[]A, []B, []T.Tuple2[A, B]](fb)
}

// Unzip is the function is reverse of [Zip]. Takes an array of pairs and return two corresponding arrays
func Unzip[A, B any](cs []T.Tuple2[A, B]) T.Tuple2[[]A, []B] {
	return G.Unzip[[]A, []B, []T.Tuple2[A, B]](cs)
}
