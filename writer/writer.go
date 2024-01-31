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

package writer

import (
	EM "github.com/IBM/fp-go/endomorphism"
	IO "github.com/IBM/fp-go/io"
	M "github.com/IBM/fp-go/monoid"
	S "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
	G "github.com/IBM/fp-go/writer/generic"
)

type Writer[W, A any] IO.IO[T.Tuple2[A, W]]

func Of[A, W any](m M.Monoid[W]) func(A) Writer[W, A] {
	return G.Of[Writer[W, A]](m)
}

// Listen modifies the result to include the changes to the accumulator
func Listen[W, A any](fa Writer[W, A]) Writer[W, T.Tuple2[A, W]] {
	return G.Listen[Writer[W, A], Writer[W, T.Tuple2[A, W]], W, A](fa)
}

// Pass applies the returned function to the accumulator
func Pass[W, A any](fa Writer[W, T.Tuple2[A, EM.Endomorphism[W]]]) Writer[W, A] {
	return G.Pass[Writer[W, T.Tuple2[A, EM.Endomorphism[W]]], Writer[W, A]](fa)
}

func MonadMap[FCT ~func(A) B, W, A, B any](fa Writer[W, A], f FCT) Writer[W, B] {
	return G.MonadMap[Writer[W, B], Writer[W, A]](fa, f)
}

func Map[W any, FCT ~func(A) B, A, B any](f FCT) func(Writer[W, A]) Writer[W, B] {
	return G.Map[Writer[W, B], Writer[W, A]](f)
}

func MonadChain[FCT ~func(A) Writer[W, B], W, A, B any](s S.Semigroup[W]) func(Writer[W, A], FCT) Writer[W, B] {
	return G.MonadChain[Writer[W, B], Writer[W, A], FCT](s)
}

func Chain[A, B, W any](s S.Semigroup[W]) func(func(A) Writer[W, B]) func(Writer[W, A]) Writer[W, B] {
	return G.Chain[Writer[W, B], Writer[W, A], func(A) Writer[W, B]](s)
}

func MonadAp[B, A, W any](s S.Semigroup[W]) func(Writer[W, func(A) B], Writer[W, A]) Writer[W, B] {
	return G.MonadAp[Writer[W, B], Writer[W, func(A) B], Writer[W, A]](s)
}

func Ap[B, A, W any](s S.Semigroup[W]) func(Writer[W, A]) func(Writer[W, func(A) B]) Writer[W, B] {
	return G.Ap[Writer[W, B], Writer[W, func(A) B], Writer[W, A]](s)
}

func MonadChainFirst[FCT ~func(A) Writer[W, B], W, A, B any](s S.Semigroup[W]) func(Writer[W, A], FCT) Writer[W, A] {
	return G.MonadChainFirst[Writer[W, B], Writer[W, A], FCT](s)
}

func ChainFirst[FCT ~func(A) Writer[W, B], W, A, B any](s S.Semigroup[W]) func(FCT) func(Writer[W, A]) Writer[W, A] {
	return G.ChainFirst[Writer[W, B], Writer[W, A], FCT](s)
}

func Flatten[W, A any](s S.Semigroup[W]) func(Writer[W, Writer[W, A]]) Writer[W, A] {
	return G.Flatten[Writer[W, Writer[W, A]], Writer[W, A]](s)
}

// Execute extracts the accumulator
func Execute[W, A any](fa Writer[W, A]) W {
	return G.Execute(fa)
}

// Evaluate extracts the value
func Evaluate[W, A any](fa Writer[W, A]) A {
	return G.Evaluate(fa)
}

// MonadCensor modifies the final accumulator value by applying a function
func MonadCensor[FCT ~func(W) W, W, A any](fa Writer[W, A], f FCT) Writer[W, A] {
	return G.MonadCensor[Writer[W, A]](fa, f)
}

// Censor modifies the final accumulator value by applying a function
func Censor[FCT ~func(W) W, W, A any](f FCT) func(Writer[W, A]) Writer[W, A] {
	return G.Censor[Writer[W, A]](f)
}

// MonadListens projects a value from modifications made to the accumulator during an action
func MonadListens[FCT ~func(W) B, W, A, B any](fa Writer[W, A], f FCT) Writer[W, T.Tuple2[A, B]] {
	return G.MonadListens[Writer[W, A], Writer[W, T.Tuple2[A, B]]](fa, f)
}

// Listens projects a value from modifications made to the accumulator during an action
func Listens[FCT ~func(W) B, W, A, B any](f FCT) func(Writer[W, A]) Writer[W, T.Tuple2[A, B]] {
	return G.Listens[Writer[W, A], Writer[W, T.Tuple2[A, B]]](f)
}
