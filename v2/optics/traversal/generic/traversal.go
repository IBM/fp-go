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

package generic

import (
	AR "github.com/IBM/fp-go/v2/array/generic"
	C "github.com/IBM/fp-go/v2/constant"
	F "github.com/IBM/fp-go/v2/function"
)

type (
	Traversal[S, A, HKTS, HKTA any] = func(func(A) HKTA) func(S) HKTS
)

func Compose[
	TAB ~func(func(B) HKTB) func(A) HKTA,
	TSA ~func(func(A) HKTA) func(S) HKTS,
	TSB ~func(func(B) HKTB) func(S) HKTS,
	S, A, B, HKTS, HKTA, HKTB any](ab TAB) func(TSA) TSB {
	return func(sa TSA) TSB {
		return F.Flow2(ab, sa)
	}
}

func FromTraversable[
	TAB ~func(func(A) HKTFA) func(HKTTA) HKTAA,
	A,
	HKTTA,
	HKTFA,
	HKTAA any](
	traverseF func(HKTTA, func(A) HKTFA) HKTAA,
) TAB {
	return F.Bind1st(F.Bind2nd[HKTTA, func(A) HKTFA, HKTAA], traverseF)
}

// FoldMap maps each target to a `Monoid` and combines the result
func FoldMap[M, S, A any](f func(A) M) func(sa Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
	return func(sa Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
		return F.Flow2(
			F.Pipe1(
				F.Flow2(f, C.Make[M, A]),
				sa,
			),
			C.Unwrap[M, S],
		)
	}
}

// Fold maps each target to a `Monoid` and combines the result
func Fold[S, A any](sa Traversal[S, A, C.Const[A, S], C.Const[A, A]]) func(S) A {
	return FoldMap[A, S](F.Identity[A])(sa)
}

// GetAll gets all the targets of a traversal
func GetAll[GA ~[]A, S, A any](s S) func(sa Traversal[S, A, C.Const[GA, S], C.Const[GA, A]]) GA {
	fmap := FoldMap[GA, S](AR.Of[GA, A])
	return func(sa Traversal[S, A, C.Const[GA, S], C.Const[GA, A]]) GA {
		return fmap(sa)(s)
	}
}
