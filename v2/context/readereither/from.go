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
	"context"

	RE "github.com/IBM/fp-go/v2/readereither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[A any](f func(context.Context) (A, error)) func() ReaderEither[A] {
	return RE.From0[ReaderEither[A]](f)
}

func From1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderEither[A] {
	return RE.From1[ReaderEither[A]](f)
}

func From2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1, T2) ReaderEither[A] {
	return RE.From2[ReaderEither[A]](f)
}

func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderEither[A] {
	return RE.From3[ReaderEither[A]](f)
}
