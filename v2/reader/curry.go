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

package reader

import (
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a reader with the context as the last parameter, which
// is a equivalent to a function returning a reader of that context
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[R, A any](f func(R) A) Reader[R, A] {
	return G.Curry0[Reader[R, A]](f)
}

func Curry1[R, T1, A any](f func(R, T1) A) func(T1) Reader[R, A] {
	return G.Curry1[Reader[R, A]](f)
}

func Curry2[R, T1, T2, A any](f func(R, T1, T2) A) func(T1) func(T2) Reader[R, A] {
	return G.Curry2[Reader[R, A]](f)
}

func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1) func(T2) func(T3) Reader[R, A] {
	return G.Curry3[Reader[R, A]](f)
}

func Curry4[R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1) func(T2) func(T3) func(T4) Reader[R, A] {
	return G.Curry4[Reader[R, A]](f)
}

func Uncurry0[R, A any](f Reader[R, A]) func(R) A {
	return G.Uncurry0(f)
}

func Uncurry1[R, T1, A any](f func(T1) Reader[R, A]) func(R, T1) A {
	return G.Uncurry1(f)
}

func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) Reader[R, A]) func(R, T1, T2) A {
	return G.Uncurry2(f)
}

func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) Reader[R, A]) func(R, T1, T2, T3) A {
	return G.Uncurry3(f)
}

func Uncurry4[R, T1, T2, T3, T4, A any](f func(T1) func(T2) func(T3) func(T4) Reader[R, A]) func(R, T1, T2, T3, T4) A {
	return G.Uncurry4(f)
}
