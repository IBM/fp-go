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

package readereither

import (
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[R, A any](f func(R) (A, error)) func() ReaderEither[R, error, A] {
	return G.From0[ReaderEither[R, error, A]](f)
}

func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderEither[R, error, A] {
	return G.From1[ReaderEither[R, error, A]](f)
}

func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderEither[R, error, A] {
	return G.From2[ReaderEither[R, error, A]](f)
}

func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderEither[R, error, A] {
	return G.From3[ReaderEither[R, error, A]](f)
}
