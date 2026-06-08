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

package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
)

// FromIdiomatic lifts a [KleisliI] function into a [Kleisli] arrow.
//
// A [KleisliI] is a standard Go function of the form
//
//	func(A) func(context.Context, R) (B, error)
//
// FromIdiomatic adapts it to the functional style [Kleisli][R, A, B] so it
// can be composed with the other combinators in this package.
func FromIdiomatic[R, A, B any](f KleisliI[R, A, B]) Kleisli[R, A, B] {
	return func(a A) ReaderReaderIOResult[R, B] {
		fa := f(a)
		return func(r R) ReaderIOResult[context.Context, B] {
			return func(ctx context.Context) IOResult[B] {
				return func() Result[B] {
					return result.TryCatchError(fa(ctx, r))
				}
			}
		}
	}
}

// MonadChainI sequences a [ReaderReaderIOResult] with an idiomatic Kleisli function.
// The success value of fa is passed to f; errors short-circuit the chain.
func MonadChainI[R, A, B any](fa ReaderReaderIOResult[R, A], f KleisliI[R, A, B]) ReaderReaderIOResult[R, B] {
	return MonadChain(fa, FromIdiomatic(f))
}

// MonadChainFirstI sequences fa with an idiomatic Kleisli function f for its side
// effects, but returns the original value of fa (not the result of f).
// If fa or f fails, the error is propagated.
func MonadChainFirstI[R, A, B any](fa ReaderReaderIOResult[R, A], f KleisliI[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirst(fa, FromIdiomatic(f))
}

// MonadTapI runs the idiomatic function f as a side-effect on the value of fa,
// discarding f's result and returning the original value of fa unchanged.
// Equivalent to [MonadChainFirstI].
func MonadTapI[R, A, B any](fa ReaderReaderIOResult[R, A], f KleisliI[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadTap(fa, FromIdiomatic(f))
}

// ChainI returns an [Operator] that sequences a computation with the idiomatic
// Kleisli function f, passing the success value forward.
// It is the curried form of [MonadChainI].
func ChainI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, B] {
	return Chain(FromIdiomatic(f))
}

// ChainFirstI returns an [Operator] that runs the idiomatic function f for its
// side effects and then returns the original monadic value unchanged.
// It is the curried form of [MonadChainFirstI].
func ChainFirstI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, A] {
	return ChainFirst(FromIdiomatic(f))
}

// TapI returns an [Operator] that executes f as a side-effect, returning the
// original value unchanged. Equivalent to [ChainFirstI].
// It is the curried form of [MonadTapI].
func TapI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, A] {
	return Tap(FromIdiomatic(f))
}

// MonadChainLeftI recovers from an error in fa by running the idiomatic function f
// on the error value. If fa succeeds its value is returned unchanged; if fa fails,
// f is called with the error and its result replaces the failed computation.
func MonadChainLeftI[R, A any](fa ReaderReaderIOResult[R, A], f KleisliI[R, error, A]) ReaderReaderIOResult[R, A] {
	return MonadChainLeft(fa, FromIdiomatic(f))
}

// ChainLeftI returns a function that recovers from errors using the idiomatic
// function f. If the input computation succeeds, its value passes through; if it fails,
// f is called with the error to produce a recovery computation.
// It is the curried form of [MonadChainLeftI].
func ChainLeftI[R, A any](f KleisliI[R, error, A]) func(ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return ChainLeft(FromIdiomatic(f))
}

// RetryingI retries the idiomatic action according to policy until check returns
// false or the policy is exhausted. It is a convenience wrapper around [Retrying] that
// accepts an idiomatic action instead of a [Kleisli].
func RetryingI[R, A any](
	policy retry.RetryPolicy,
	action KleisliI[R, retry.RetryStatus, A],
	check Predicate[Result[A]],
) ReaderReaderIOResult[R, A] {
	return Retrying(policy, FromIdiomatic(action), check)
}

// TraverseArrayI maps each element of a slice through the idiomatic function f
// and collects the results into a [ReaderReaderIOResult] holding a slice.
// The first error encountered short-circuits the traversal.
func TraverseArrayI[R, A, B any](f KleisliI[R, A, B]) Kleisli[R, []A, []B] {
	return TraverseArray(FromIdiomatic(f))
}

// BindI is the idiomatic-function variant of [Bind].
// It attaches the result of the idiomatic Kleisli f to a do-notation context using setter.
func BindI[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f KleisliI[R, S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, FromIdiomatic(f))
}

// BindIL is the lens-based, idiomatic-function variant of [BindL].
// It reads the focused field T from the context S via lens, passes it to f,
// and writes the result back through the same lens.
func BindIL[R, S, T any](
	lens Lens[S, T],
	f KleisliI[R, T, T],
) Operator[R, S, S] {
	return BindL(lens, FromIdiomatic(f))
}
