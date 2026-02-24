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

	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

// Provide supplies a context to an effect, converting it to a Thunk.
// This is the first step in running an effect - it eliminates the context dependency
// by providing the required context value.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value
//
// # Parameters
//
//   - c: The context value to provide to the effect
//
// # Returns
//
//   - func(Effect[C, A]) ReaderIOResult[A]: A function that converts an effect to a thunk
//
// # Example
//
//	ctx := MyContext{APIKey: "secret"}
//	eff := effect.Of[MyContext](42)
//	thunk := effect.Provide[MyContext, int](ctx)(eff)
//	// thunk is now a ReaderIOResult[int] that can be run
func Provide[C, A any](c C) func(Effect[C, A]) ReaderIOResult[A] {
	return readerreaderioresult.Read[A](c)
}

// RunSync executes a Thunk synchronously, converting it to a standard Go function.
// This is the final step in running an effect - it executes the IO operations
// and returns the result as a standard (value, error) tuple.
//
// # Type Parameters
//
//   - A: The type of the success value
//
// # Parameters
//
//   - fa: The thunk to execute
//
// # Returns
//
//   - readerresult.ReaderResult[A]: A function that takes a context.Context and returns (A, error)
//
// # Example
//
//	ctx := MyContext{APIKey: "secret"}
//	eff := effect.Of[MyContext](42)
//	thunk := effect.Provide[MyContext, int](ctx)(eff)
//	readerResult := effect.RunSync(thunk)
//	value, err := readerResult(context.Background())
//	// value == 42, err == nil
//
// # Complete Example
//
//	// Typical usage pattern:
//	result, err := effect.RunSync(
//		effect.Provide[MyContext, string](myContext)(myEffect),
//	)(context.Background())
func RunSync[A any](fa ReaderIOResult[A]) readerresult.ReaderResult[A] {
	return func(ctx context.Context) (A, error) {
		return result.Unwrap(fa(ctx)())
	}
}
