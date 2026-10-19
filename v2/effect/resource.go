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

package effect

import (
	"context"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

// Bracket acquires a resource, uses it and releases it, also when the use fails.
//
// This is the effect counterpart of try/finally and of acquireUseRelease in effect-ts.
// The release function receives the acquired resource and the outcome of use, so it can
// decide, for example, whether to commit or roll back a transaction. If acquire fails,
// neither use nor release run. The outcome of use is returned; a failure of release
// does not replace it.
//
// The release runs uninterruptibly (see Uninterruptible), like a finalizer in effect-ts:
// it keeps the context values but is not cancelled, so the resource is released even if
// the context was cancelled or timed out while the resource was in use. A release that
// must be bounded can apply WithTimeout itself.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the acquired resource
//   - B: The type of the value produced by use
//   - ANY: The value produced by release, which is discarded
//
// Parameters:
//   - acquire: The effect that acquires the resource
//   - use: The function that uses the resource
//   - release: The function that releases the resource, given the outcome of use
//
// Returns:
//   - Effect[C, B]: An effect that produces the outcome of use
//
// See Also:
//   - WithResource: Bracket for releases that do not need the outcome of use
func Bracket[C, A, B, ANY any](
	acquire Effect[C, A],
	use Kleisli[C, A, B],
	release func(A, Result[B]) Effect[C, ANY],
) Effect[C, B] {
	return readerreaderioresult.Bracket(acquire, use, release)
}

// Uninterruptible runs an effect with a context that cannot be cancelled.
//
// The effect receives context.WithoutCancel of the caller's context. It keeps all context
// values (logger, request ID, ...), but it has no deadline and is never cancelled, so
// neither a cancellation nor a timeout of the caller interrupts it. This corresponds to
// Effect.uninterruptible in effect-ts. Bracket, WithResource and WithLock run their
// release this way.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - eff: The effect to protect from cancellation
//
// Returns:
//   - Effect[C, A]: The effect, run with a context that is never cancelled
//
// See Also:
//   - Bracket: Runs its release uninterruptibly
//   - WithTimeout: Bounds an effect, also inside Uninterruptible
//
//go:inline
func Uninterruptible[C, A any](eff Effect[C, A]) Effect[C, A] {
	return readerreaderioresult.Uninterruptible(eff)
}

// WithResource builds a scoped use of a resource from the functions that create and
// release it.
//
// The result is a function that accepts the code that uses the resource. The resource
// is created before that code runs and released after it completes or fails, following
// the semantics of Bracket.
//
// Type Parameters:
//   - A: The type of the value produced by the code that uses the resource
//   - C: The dependencies required by the effects
//   - R: The type of the resource
//   - ANY: The value produced by onRelease, which is discarded
//
// Parameters:
//   - onCreate: The effect that creates the resource
//   - onRelease: The function that releases the resource
//
// Returns:
//   - Kleisli[C, Kleisli[C, R, A], A]: A function from the use of the resource to an effect
//
// See Also:
//   - Bracket: The general form, with access to the outcome of use
func WithResource[A, C, R, ANY any](onCreate Effect[C, R], onRelease Kleisli[C, R, ANY]) Kleisli[C, Kleisli[C, R, A], A] {
	release := func(r R, _ Result[A]) Effect[C, ANY] {
		return onRelease(r)
	}
	return func(use Kleisli[C, R, A]) Effect[C, A] {
		return Bracket(onCreate, use, release)
	}
}

// WithLock runs an effect while holding a lock.
//
// The lock effect acquires the lock and produces the function that releases it. The
// release function is called after the guarded effect completes or fails.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the value produced by the guarded effect
//
// Parameters:
//   - lock: The effect that acquires the lock and returns its release function
//
// Returns:
//   - Operator[C, A, A]: A function that guards an effect with the lock
//
// See Also:
//   - WithSemaphore: Limits concurrency instead of enforcing mutual exclusion
func WithLock[C, A any](lock Effect[C, context.CancelFunc]) Operator[C, A, A] {
	return function.Flow2(
		function.Constant1[context.CancelFunc, Effect[C, A]],
		WithResource[A](lock, function.Flow2(
			io.FromImpure[context.CancelFunc],
			FromIO[C, Void],
		)),
	)
}

// WithSemaphore runs an effect while holding one slot of a counting semaphore.
//
// The semaphore is a buffered channel whose capacity is the number of slots, for example
// make(chan Void, 3). Effects guarded by the same semaphore run at most that many at a
// time. Waiting for a slot honours the context: if the context is cancelled while waiting,
// the effect does not run and fails with the cancellation cause.
//
// This corresponds to Semaphore.withPermits(1) in effect-ts.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - sem: A buffered channel used as a counting semaphore
//
// Returns:
//   - Operator[C, A, A]: A function that guards an effect with the semaphore
//
// See Also:
//   - WithLock: Mutual exclusion with an arbitrary lock
//   - TraverseArrayParN: Bounded parallel traversal of an array
func WithSemaphore[C, A any](sem chan Void) Operator[C, A, A] {
	return reader.Map[C](thunk.WithSemaphore[A](sem))
}
