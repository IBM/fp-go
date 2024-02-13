// Copyright (c) 2024 IBM Corp.
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

package writer

import (
	"github.com/IBM/fp-go/internal/applicative"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	M "github.com/IBM/fp-go/monoid"
	G "github.com/IBM/fp-go/writer/generic"
)

// Pointed implements the pointed operations for [Writer]
func Pointed[W, A any](m M.Monoid[W]) pointed.Pointed[A, Writer[W, A]] {
	return G.Pointed[Writer[W, A], W, A](m)
}

// Functor implements the pointed operations for [Writer]
func Functor[W, A, B any]() functor.Functor[A, B, Writer[W, A], Writer[W, B]] {
	return G.Functor[Writer[W, B], Writer[W, A], W, A, B]()
}

// Applicative implements the applicative operations for [Writer]
func Applicative[W, A, B any](m M.Monoid[W]) applicative.Applicative[A, B, Writer[W, A], Writer[W, B], Writer[W, func(A) B]] {
	return G.Applicative[Writer[W, B], Writer[W, func(A) B], Writer[W, A]](m)
}

// Monad implements the monadic operations for [Writer]
func Monad[W, A, B any](m M.Monoid[W]) monad.Monad[A, B, Writer[W, A], Writer[W, B], Writer[W, func(A) B]] {
	return G.Monad[Writer[W, B], Writer[W, func(A) B], Writer[W, A]](m)
}
