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

package readerreaderioeither

import (
	G "github.com/IBM/fp-go/v2/internal/bracket"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/readerio"
)

// Bracket safely acquires a resource, uses it, and releases it again, following the
// acquire-use-release pattern. All three steps run with access to both the outer (R) and
// the inner (C) environment.
//
// Semantics:
//
//  1. acquire is executed. If it fails, its error is returned and neither use nor
//     release is executed (there is nothing to release).
//  2. use is executed with the acquired resource.
//  3. release is executed with the resource and the outcome of use (an Either[E, B]).
//     It runs no matter whether use succeeded or failed.
//  4. If release succeeds, the outcome of use is returned unchanged (a failed use stays
//     failed). If release fails, its error is returned instead, even if use succeeded,
//     and it replaces the error of a failed use.
//
// The value produced by a successful release (of type ANY) is discarded.
//
// Bracket is lazy: no step runs until both environments have been supplied and the
// resulting IO is executed. Every execution acquires and releases a fresh resource.
//
// Bracket is implemented point-free on top of the generic bracket combinator: the
// resource is threaded through [MonadChain], and the outcome of use, which may be a
// Left, is fed into release via a chain over the full Either.
//
// Type Parameters:
//   - R: The outer environment type
//   - C: The inner environment type
//   - E: The error type
//   - A: The resource type
//   - B: The result type of use
//   - ANY: The (ignored) result type of release
//
// Parameters:
//   - acquire: Produces the resource
//   - use: Computes the result from the resource
//   - release: Frees the resource; it receives the resource and the outcome of use
//
// Returns:
//   - A ReaderReaderIOEither[R, C, E, B] that yields the result of use and guarantees
//     that release runs whenever acquire succeeded
//
// See Also:
//   - ExampleBracket: basic acquire/use/close pattern
//   - ExampleBracket_useFails: resource released even when use returns Left
//   - ExampleBracket_acquireFails: neither use nor release run when acquire returns Left
//   - ExampleBracket_transaction: release receives the outcome of use to commit or roll back
//
//go:inline
func Bracket[
	R, C, E, A, B, ANY any](
	acquire ReaderReaderIOEither[R, C, E, A],
	use Kleisli[R, C, E, A, B],
	release func(A, Either[E, B]) ReaderReaderIOEither[R, C, E, ANY],
) ReaderReaderIOEither[R, C, E, B] {
	return G.MonadBracket[
		ReaderReaderIOEither[R, C, E, A],
		ReaderReaderIOEither[R, C, E, B],
		ReaderReaderIOEither[R, C, E, ANY],
		Either[E, B],
		A,
		B,
	](
		FromEither[R, C, E, B],
		MonadChain[R, C, E, A, B],
		monadChainReaderReaderIO[R, C, E, B, B],
		MonadChain[R, C, E, ANY, B],

		acquire,
		use,
		release,
	)
}

// monadChainReaderReaderIO is the chain of the underlying ReaderReaderIO monad, i.e. it treats
// a ReaderReaderIOEither as a ReaderReaderIO of an Either and sequences it with a function that
// receives the complete outcome, the Either itself. Unlike [MonadChain] the continuation also
// runs for a Left, which is what allows Bracket to release a resource after a failed use.
//
//go:inline
func monadChainReaderReaderIO[R, C, E, A, B any](
	fa ReaderReaderIOEither[R, C, E, A],
	f func(Either[E, A]) ReaderReaderIOEither[R, C, E, B],
) ReaderReaderIOEither[R, C, E, B] {
	return readert.MonadChain(
		readerio.MonadChain[C, Either[E, A], Either[E, B]],
		fa,
		f,
	)
}
