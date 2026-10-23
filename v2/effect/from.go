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
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// FromReaderIO lifts a ReaderIO into an Effect.
//
// The ReaderIO reads the effect's context C and performs an IO action that
// cannot fail, so the resulting Effect always succeeds with the produced value.
// The action does not see the runtime context.Context.
//
// Type Parameters:
//   - C: The context type read by the ReaderIO
//   - A: The type of the value produced by the ReaderIO
//
// Parameters:
//   - ma: The ReaderIO[C, A] to lift
//
// Returns:
//   - Effect[C, A]: An effect that always succeeds with the value produced by ma
//
// See Also:
//   - FromReader: Lifts a context-dependent computation without IO
//   - FromReaderIOResult: Lifts a context-dependent IO computation that may fail
//   - ChainReaderIOK: Chains a ReaderIO-returning function over an existing Effect
//
//go:inline
func FromReaderIO[C, A any](ma ReaderIO[C, A]) Effect[C, A] {
	return readerreaderioresult.FromReaderIO(ma)
}

// FromReaderIOResult lifts a ReaderIOResult from the readerioresult package
// into an Effect.
//
// The ReaderIOResult reads the effect's context C and performs an IO action
// that may fail. A successful result becomes a successful Effect, an error a
// failed Effect. The action does not see the runtime context.Context; use
// FromReaderThunk for computations that need it.
//
// Note that the parameter is a readerioresult.ReaderIOResult[C, A]
// (func(C) IOResult[A]), not the context-specific ReaderIOResult[A] of this
// package, which is the same as Thunk[A].
//
// Type Parameters:
//   - C: The context type read by the ReaderIOResult
//   - A: The type of the success value
//
// Parameters:
//   - ma: The readerioresult.ReaderIOResult[C, A] to lift
//
// Returns:
//   - Effect[C, A]: An effect that succeeds or fails according to ma
//
// See Also:
//   - FromReaderResult: Lifts a context-dependent computation that may fail without IO
//   - FromReaderThunk: Lifts a context-dependent computation that also sees context.Context
//   - FromIOResult: Lifts an IOResult that does not depend on C
//
//go:inline
func FromReaderIOResult[C, A any](ma RIOR.ReaderIOResult[C, A]) Effect[C, A] {
	return readerreaderioresult.FromReaderIOResult(ma)
}

// FromReaderThunk lifts a function from the context C to a Thunk into an Effect.
//
// A Reader[C, Thunk[A]] has exactly the shape of an Effect[C, A], so this is a
// type-level conversion that documents intent at the call site. The Thunk sees
// both the effect's context C and the runtime context.Context.
//
// Type Parameters:
//   - C: The context type read by the function
//   - A: The type of the success value
//
// Parameters:
//   - ma: A function from C to Thunk[A]
//
// Returns:
//   - Effect[C, A]: An effect that runs the Thunk produced for its context
//
// See Also:
//   - FromThunk: Lifts a Thunk that does not depend on C
//   - FromReaderIOResult: Similar but the action does not see context.Context
//
//go:inline
func FromReaderThunk[C, A any](ma Reader[C, Thunk[A]]) Effect[C, A] {
	return ma
}
