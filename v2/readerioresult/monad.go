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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

// Pointed returns the pointed operations for [ReaderIOResult]
//
//go:inline
func Pointed[R, A any]() pointed.Pointed[A, ReaderIOResult[R, A]] {
	return RIOE.Pointed[R, error, A]()
}

// Functor returns the functor operations for [ReaderIOResult]
//
//go:inline
func Functor[R, A, B any]() functor.Functor[A, B, ReaderIOResult[R, A], ReaderIOResult[R, B]] {
	return RIOE.Functor[R, error, A, B]()
}

// Monad returns the monadic operations for [ReaderIOResult]
//
//go:inline
func Monad[R, A, B any]() monad.Monad[A, B, ReaderIOResult[R, A], ReaderIOResult[R, B], ReaderIOResult[R, func(A) B]] {
	return RIOE.Monad[R, error, A, B]()
}
