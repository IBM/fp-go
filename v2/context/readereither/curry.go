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

package readereither

import (
	"context"

	"github.com/IBM/fp-go/v2/readereither"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[A any](f func(context.Context) (A, error)) ReaderEither[A] {
	return readereither.Curry0(f)
}

func Curry1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderEither[A] {
	return readereither.Curry1(f)
}

func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) func(T2) ReaderEither[A] {
	return readereither.Curry2(f)
}

func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderEither[A] {
	return readereither.Curry3(f)
}

func Uncurry1[T1, A any](f func(T1) ReaderEither[A]) func(context.Context, T1) (A, error) {
	return readereither.Uncurry1(f)
}

func Uncurry2[T1, T2, A any](f func(T1) func(T2) ReaderEither[A]) func(context.Context, T1, T2) (A, error) {
	return readereither.Uncurry2(f)
}

func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderEither[A]) func(context.Context, T1, T2, T3) (A, error) {
	return readereither.Uncurry3(f)
}
