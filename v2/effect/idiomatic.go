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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/retry"
)

// FromIdiomatic lifts a [KleisliIdiomatic] function into a [Kleisli] arrow.
//
// A [KleisliIdiomatic] is a standard Go function of the form
//
//	func(A) func(context.Context, R) (B, error)
//
// FromIdiomatic adapts it to the functional style [Kleisli][R, A, B] so it
// can be composed with the other combinators in this package.
func FromIdiomatic[R, A, B any](f KleisliIdiomatic[R, A, B]) Kleisli[R, A, B] {
	return readerreaderioresult.FromIdiomatic(f)
}

// MonadChainIdiomatic sequences a [Effect] with an idiomatic Kleisli function.
// The success value of fa is passed to f; errors short-circuit the chain.
func MonadChainIdiomatic[R, A, B any](fa Effect[R, A], f KleisliIdiomatic[R, A, B]) Effect[R, B] {
	return readerreaderioresult.MonadChainIdiomatic(fa, f)
}

// MonadChainFirstIdiomatic sequences fa with an idiomatic Kleisli function f for its side
// effects, but returns the original value of fa (not the result of f).
// If fa or f fails, the error is propagated.
func MonadChainFirstIdiomatic[R, A, B any](fa Effect[R, A], f KleisliIdiomatic[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstIdiomatic(fa, f)
}

// MonadTapIdiomatic runs the idiomatic function f as a side-effect on the value of fa,
// discarding f's result and returning the original value of fa unchanged.
// Equivalent to [MonadChainFirstIdiomatic].
func MonadTapIdiomatic[R, A, B any](fa Effect[R, A], f KleisliIdiomatic[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapIdiomatic(fa, f)
}

// ChainIdiomatic returns an [Operator] that sequences a computation with the idiomatic
// Kleisli function f, passing the success value forward.
// It is the curried form of [MonadChainIdiomatic].
func ChainIdiomatic[R, A, B any](f KleisliIdiomatic[R, A, B]) Operator[R, A, B] {
	return readerreaderioresult.ChainIdiomatic(f)
}

// ChainFirstIdiomatic returns an [Operator] that runs the idiomatic function f for its
// side effects and then returns the original monadic value unchanged.
// It is the curried form of [MonadChainFirstIdiomatic].
func ChainFirstIdiomatic[R, A, B any](f KleisliIdiomatic[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstIdiomatic(f)
}

// TapIdiomatic returns an [Operator] that executes f as a side-effect, returning the
// original value unchanged. Equivalent to [ChainFirstIdiomatic].
// It is the curried form of [MonadTapIdiomatic].
func TapIdiomatic[R, A, B any](f KleisliIdiomatic[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.TapIdiomatic(f)
}

// MonadChainLeftIdiomatic recovers from an error in fa by running the idiomatic function f
// on the error value. If fa succeeds its value is returned unchanged; if fa fails,
// f is called with the error and its result replaces the failed computation.
func MonadChainLeftIdiomatic[R, A any](fa Effect[R, A], f KleisliIdiomatic[R, error, A]) Effect[R, A] {
	return readerreaderioresult.MonadChainLeftIdiomatic(fa, f)
}

// ChainLeftIdiomatic returns a function that recovers from errors using the idiomatic
// function f. If the input computation succeeds, its value passes through; if it fails,
// f is called with the error to produce a recovery computation.
// It is the curried form of [MonadChainLeftIdiomatic].
func ChainLeftIdiomatic[R, A any](f KleisliIdiomatic[R, error, A]) Operator[R, A, A] {
	return readerreaderioresult.ChainLeftIdiomatic(f)
}

// RetryingIdiomatic retries the idiomatic action according to policy until check returns
// false or the policy is exhausted. It is a convenience wrapper around [Retrying] that
// accepts an idiomatic action instead of a [Kleisli].
func RetryingIdiomatic[R, A any](
	policy retry.RetryPolicy,
	action KleisliIdiomatic[R, retry.RetryStatus, A],
	check Predicate[Result[A]],
) Effect[R, A] {
	return readerreaderioresult.RetryingIdiomatic(policy, action, check)
}

// TraverseArrayIdiomatic maps each element of a slice through the idiomatic function f
// and collects the results into a [Effect] holding a slice.
// The first error encountered short-circuits the traversal.
func TraverseArrayIdiomatic[R, A, B any](f KleisliIdiomatic[R, A, B]) Kleisli[R, []A, []B] {
	return readerreaderioresult.TraverseArrayIdiomatic(f)
}

// BindIdiomatic is the idiomatic-function variant of [Bind].
// It attaches the result of the idiomatic Kleisli f to a do-notation context using setter.
func BindIdiomatic[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f KleisliIdiomatic[R, S1, T],
) Operator[R, S1, S2] {
	return readerreaderioresult.BindIdiomatic(setter, f)
}

// BindLIdiomatic is the lens-based, idiomatic-function variant of [BindL].
// It reads the focused field T from the context S via lens, passes it to f,
// and writes the result back through the same lens.
func BindLIdiomatic[R, S, T any](
	lens Lens[S, T],
	f KleisliIdiomatic[R, T, T],
) Operator[R, S, S] {
	return readerreaderioresult.BindLIdiomatic(lens, f)
}
