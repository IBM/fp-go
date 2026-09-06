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

package effect

import (
	"context"
	"fmt"
	"slices"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
)

// exFilterConfig is the dependency used by the filter examples.
type exFilterConfig struct{}

// isEven is a point-free predicate: take the remainder modulo 2, test for zero.
var isEven = F.Flow2(N.Mod(2), P.IsZero[int]())

// ExampleFilterArray demonstrates keeping only the elements of a slice that
// satisfy a predicate, inside an effect.
func ExampleFilterArray() {
	eff := F.Pipe1(
		Succeed[exFilterConfig]([]int{1, 2, 3, 4, 5, 6}),
		FilterArray[exFilterConfig](isEven),
	)

	fmt.Println(RunSync(Provide[[]int](exFilterConfig{})(eff))(context.Background()))
	// Output:
	// [2 4 6] <nil>
}

// ExampleFilterArray_errorsPropagate demonstrates that a failed effect is passed
// through untouched — the predicate is never applied.
func ExampleFilterArray_errorsPropagate() {
	eff := F.Pipe1(
		Fail[exFilterConfig, []int](fmt.Errorf("upstream failed")),
		FilterArray[exFilterConfig](isEven),
	)

	fmt.Println(RunSync(Provide[[]int](exFilterConfig{})(eff))(context.Background()))
	// Output:
	// [] upstream failed
}

// ExampleFilterMapArray demonstrates filtering and mapping in a single pass:
// the Kleisli arrow returns None for elements to drop and Some for the mapped
// value of elements to keep.
func ExampleFilterMapArray() {
	// Point-free: keep positive numbers, then render them.
	positiveLabel := F.Flow2(
		O.FromPredicate(N.MoreThan(0)),
		O.Map(S.Format[int]("n=%d")),
	)

	eff := F.Pipe1(
		Succeed[exFilterConfig]([]int{-1, 2, -3, 4}),
		FilterMapArray[exFilterConfig](positiveLabel),
	)

	fmt.Println(RunSync(Provide[[]string](exFilterConfig{})(eff))(context.Background()))
	// Output:
	// [n=2 n=4] <nil>
}

// ExampleFilterIter demonstrates filtering a lazy iterator sequence.
func ExampleFilterIter() {
	eff := F.Pipe1(
		Succeed[exFilterConfig](iter.From(1, 2, 3, 4, 5, 6)),
		FilterIter[exFilterConfig](isEven),
	)

	seq, err := RunSync(Provide[Seq[int]](exFilterConfig{})(eff))(context.Background())

	fmt.Println(slices.Collect(seq), err)
	// Output:
	// [2 4 6] <nil>
}

// ExampleFilterMapIter demonstrates filtering and mapping a lazy iterator in a
// single pass.
func ExampleFilterMapIter() {
	positiveLabel := F.Flow2(
		O.FromPredicate(N.MoreThan(0)),
		O.Map(S.Format[int]("n=%d")),
	)

	eff := F.Pipe1(
		Succeed[exFilterConfig](iter.From(-1, 2, -3, 4)),
		FilterMapIter[exFilterConfig](positiveLabel),
	)

	seq, err := RunSync(Provide[Seq[string]](exFilterConfig{})(eff))(context.Background())

	fmt.Println(slices.Collect(seq), err)
	// Output:
	// [n=2 n=4] <nil>
}

// ExampleFilter demonstrates the container-agnostic form.  FilterArray and
// FilterIter are the two ready-made specialisations; Filter lets you supply the
// filtering function for any other container.
func ExampleFilter() {
	// A.Filter is already `func(Predicate[A]) Endomorphism[[]A]`, so it can be
	// handed to Filter directly — no wrapper lambda needed.
	filterSlice := Filter[exFilterConfig](A.Filter[int])

	eff := F.Pipe1(
		Succeed[exFilterConfig]([]int{1, 2, 3, 4, 5}),
		filterSlice(N.MoreThan(2)),
	)

	fmt.Println(RunSync(Provide[[]int](exFilterConfig{})(eff))(context.Background()))
	// Output:
	// [3 4 5] <nil>
}

// ExampleFilterMap demonstrates the container-agnostic filter-and-map form.
func ExampleFilterMap() {
	// A.FilterMap already has the shape FilterMap expects.
	filterMapSlice := FilterMap[exFilterConfig](A.FilterMap[int, string])

	eff := F.Pipe1(
		Succeed[exFilterConfig]([]int{-1, 2, -3, 4}),
		filterMapSlice(F.Flow2(
			O.FromPredicate(N.MoreThan(0)),
			O.Map(S.Format[int]("n=%d")),
		)),
	)

	fmt.Println(RunSync(Provide[[]string](exFilterConfig{})(eff))(context.Background()))
	// Output:
	// [n=2 n=4] <nil>
}
