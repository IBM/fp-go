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

package readerio

import (
	"github.com/IBM/fp-go/v2/reader"
)

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[F ~func(R) IO[A], R, A any](f func(R) IO[A]) func() ReaderIO[R, A] {
	return reader.From0(f)
}

func From1[F ~func(R, T1) IO[A], R, T1, A any](f func(R, T1) IO[A]) func(T1) ReaderIO[R, A] {
	return reader.From1(f)
}

func From2[F ~func(R, T1, T2) IO[A], R, T1, T2, A any](f func(R, T1, T2) IO[A]) func(T1, T2) ReaderIO[R, A] {
	return reader.From2(f)
}

func From3[F ~func(R, T1, T2, T3) IO[A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) IO[A]) func(T1, T2, T3) ReaderIO[R, A] {
	return reader.From3(f)
}
