// Copyright (c) 2024 - 2025 IBM Corp.
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
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/pair"
)

// SequencePair is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.
// The function takes a [Pair] of higher higher kinded types and returns a higher kinded type of a [Pair] with the resolved values.
func SequencePair[
	MAP ~func(func(T1) func(T2) P.Pair[T1, T2]) func(HKT_T1) HKT_F_T2,
	AP1 ~func(HKT_T2) func(HKT_F_T2) HKT_PAIR,
	T1,
	T2,
	HKT_T1, // HKT[T1]
	HKT_T2, // HKT[T2]
	HKT_F_T2, // HKT[func(T2) P.Pair[T1, T2]]
	HKT_PAIR any, // HKT[Pair[T1, T2]]
](
	fmap MAP,
	fap1 AP1,
	t P.Pair[HKT_T1, HKT_T2],
) HKT_PAIR {
	return F.Pipe2(
		P.Head(t),
		fmap(F.Curry2(P.MakePair[T1, T2])),
		fap1(P.Tail(t)),
	)
}

// TraversePair is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.
// The function takes a [Pair] of base types and 2 functions that transform these based types into higher higher kinded types. It returns a higher kinded type of a [Pair] with the resolved values.
func TraversePair[
	MAP ~func(func(T1) func(T2) P.Pair[T1, T2]) func(HKT_T1) HKT_F_T2,
	AP1 ~func(HKT_T2) func(HKT_F_T2) HKT_PAIR,
	F1 ~func(A1) HKT_T1,
	F2 ~func(A2) HKT_T2,
	A1, T1,
	A2, T2,
	HKT_T1, // HKT[T1]
	HKT_T2, // HKT[T2]
	HKT_F_T2, // HKT[func(T2) P.Pair[T1, T2]]
	HKT_PAIR any, // HKT[Pair[T1, T2]]
](
	fmap MAP,
	fap1 AP1,
	f1 F1,
	f2 F2,
	t P.Pair[A1, A2],
) HKT_PAIR {
	return F.Pipe2(
		f1(P.Head(t)),
		fmap(F.Curry2(P.MakePair[T1, T2])),
		fap1(f2(P.Tail(t))),
	)
}
