//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/pair"
)

// Paired converts a [ReaderReaderIOResult] into a single-argument function that accepts
// a [Pair] bundling the context.Context (head) and the outer environment R (tail).
//
// Type structure:
//
//	Paired: (R -> context.Context -> IO[Either[error, A]]) -> Pair[context.Context, R] -> IO[Either[error, A]]
//
// This is useful when you need to treat a two-environment computation as a single-argument
// function, for example when mapping over a collection of (context.Context, R) pairs or
// composing with higher-order functions that expect a single-argument function.
//
// The pair places R in the tail and context.Context in the head. This follows the [pair]
// package convention where the tail is the primary value (the one [pair.Map] and other
// functor operations act on) and the head carries auxiliary data. Here R is the primary
// environment the computation transforms, while context.Context is auxiliary context
// threaded through unchanged.
//
//	p := pair.MakePair[context.Context, R](ctx, r)
//	result := Paired(f)(p)  // equivalent to f(r)(ctx)()
//
// Example:
//
//	fetch := func(cfg AppConfig) readerioresult.ReaderIOResult[string] {
//	    return func(ctx context.Context) ioresult.IOResult[string] {
//	        return ioresult.Of("hello")
//	    }
//	}
//	paired := Paired(fetch)
//	p := pair.MakePair[context.Context, AppConfig](ctx, cfg)
//	res := paired(p)()  // Either[error, string]
func Paired[R, A any](f ReaderReaderIOResult[R, A]) ioresult.Kleisli[Pair[context.Context, R], A] {
	return func(t Pair[context.Context, R]) IOResult[A] {
		return f(pair.Tail(t))(pair.Head(t))
	}
}
