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

package traversal

import (
	C "github.com/IBM/fp-go/v2/constant"
	F "github.com/IBM/fp-go/v2/function"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// Id is the identity constructor of a traversal
func Id[S, A any]() G.Traversal[S, S, A, A] {
	return F.Identity[func(S) A]
}

// Modify applies a transformation function to a traversal
func Modify[S, A any](f func(A) A) func(sa G.Traversal[S, A, S, A]) func(S) S {
	return func(sa G.Traversal[S, A, S, A]) func(S) S {
		return sa(f)
	}
}

// Set sets a constant value for all values of the traversal
func Set[S, A any](a A) func(sa G.Traversal[S, A, S, A]) func(S) S {
	return Modify[S](F.Constant1[A](a))
}

// FoldMap maps each target to a `Monoid` and combines the result
func FoldMap[M, S, A any](f func(A) M) func(sa G.Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
	return G.FoldMap[M, S](f)
}

// Fold maps each target to a `Monoid` and combines the result
func Fold[S, A any](sa G.Traversal[S, A, C.Const[A, S], C.Const[A, A]]) func(S) A {
	return G.Fold(sa)
}

// GetAll gets all the targets of a traversal
func GetAll[S, A any](s S) func(sa G.Traversal[S, A, C.Const[[]A, S], C.Const[[]A, A]]) []A {
	return G.GetAll[[]A](s)
}

// Compose composes two traversables
func Compose[
	S, A, B, HKTS, HKTA, HKTB any](ab G.Traversal[A, B, HKTA, HKTB]) func(sa G.Traversal[S, A, HKTS, HKTA]) G.Traversal[S, B, HKTS, HKTB] {
	return G.Compose[
		G.Traversal[A, B, HKTA, HKTB],
		G.Traversal[S, A, HKTS, HKTA],
		G.Traversal[S, B, HKTS, HKTB]](ab)
}
