//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package reader

import (
	"context"

	R "github.com/ibm/fp-go/reader/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[A any](f func(context.Context) A) Reader[A] {
	return R.Curry0[Reader[A]](f)
}

func Curry1[T1, A any](f func(context.Context, T1) A) func(T1) Reader[A] {
	return R.Curry1[Reader[A]](f)
}

func Curry2[T1, T2, A any](f func(context.Context, T1, T2) A) func(T1) func(T2) Reader[A] {
	return R.Curry2[Reader[A]](f)
}

func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) A) func(T1) func(T2) func(T3) Reader[A] {
	return R.Curry3[Reader[A]](f)
}

func Uncurry1[T1, A any](f func(T1) Reader[A]) func(context.Context, T1) A {
	return R.Uncurry1(f)
}

func Uncurry2[T1, T2, A any](f func(T1) func(T2) Reader[A]) func(context.Context, T1, T2) A {
	return R.Uncurry2(f)
}

func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) func(T3) Reader[A]) func(context.Context, T1, T2, T3) A {
	return R.Uncurry3(f)
}
