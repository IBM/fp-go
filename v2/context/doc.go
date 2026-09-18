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

// Package context contains versions of reader IO monads that work with a golang [context.Context].
//
// The sub-packages fix the Reader environment to [context.Context], so a computation
// is a function of the context that is only evaluated when a context is supplied:
//
//	reader.Reader[A]                               = func(context.Context) A
//	readerio.ReaderIO[A]                           = func(context.Context) func() A
//	readerresult.ReaderResult[A]                   = func(context.Context) Result[A]
//	readerioresult.ReaderIOResult[A]               = func(context.Context) func() Result[A]
//	statereaderioresult.StateReaderIOResult[S, A]  = func(S) func(context.Context) func() Result[Pair[S, A]]
//	idiomatic/context/readerresult.ReaderResult[A] = func(context.Context) (A, error)
//
// # Working with context.Context the fp-go way
//
// In imperative Go the context is threaded by hand: every function takes ctx as its
// first argument, values are read with ctx.Value(key).(T), and scopes are opened with
// ctx, cancel := context.WithTimeout(...); defer cancel(). In fp-go the context is the
// Reader environment instead. Pipelines are built without mentioning ctx at all, and
// the context is supplied exactly once, when the pipeline is run. Reading and scoping
// the context are expressed with a small set of operators that exist uniformly in
// readerio, readerresult, readerioresult, statereaderioresult and the idiomatic
// readerresult package.
//
// # Reading the context
//
//   - Ask() returns the whole context as the computation's value, and Asks(f) / FromReader(f)
//     derive a value from it with a pure function (the exact set varies by package).
//   - AskValue[V](key) reads a single typed value. It yields Some(v) if the key holds a V,
//     and None if the key is absent or holds a value of a different type. It never panics
//     and never fails, so the caller decides what "missing" means.
//
// Use a default for optional values, or turn the None into an error for required ones:
//
//	type ctxKey string
//	const userKey ctxKey = "user"
//
//	// optional: fall back to a default
//	userOrAnonymous := F.Pipe1(
//	    readerioresult.AskValue[string](userKey),
//	    readerioresult.Map(O.GetOrElse(F.Constant("anonymous"))),
//	)
//
//	// required: fail if absent
//	requireUser := F.Pipe1(
//	    readerioresult.AskValue[string](userKey),
//	    readerioresult.Chain(readerioresult.FromOption[string](F.Constant(errNoUser))),
//	)
//
// The plain Reader variant, reader.AskValue[V](key), has type func(context.Context) Option[V]
// and is the building block for callbacks that expect a Reader (for example logger lookups).
// Where a lens is more convenient, optics/lenses.AtContext[V](key) offers the same getter
// together with a setter.
//
// # Scoping the context
//
// Scoping operators run a computation with a derived context. The derived context is visible
// to the wrapped computation only, and its cancel function is always called once the
// computation completes, so no cancel function leaks and no defer is needed at the call site.
//
//   - WithValue[A](key, value) adds key → value (see [context.WithValue]).
//   - WithTimeout[A](d) bounds the computation by a relative duration (see [context.WithTimeout]).
//   - WithDeadline[A](t) bounds the computation by an absolute time (see [context.WithDeadline]).
//     If the parent context already has an earlier deadline, that one wins.
//   - Local(f) is the general form: f maps the context to a new context and its cancel
//     function. All of the above are implemented in terms of Local.
//
// Because they are ordinary operators they compose in pipelines; the outermost operator is
// applied first when the context flows inwards:
//
//	handler := F.Pipe3(
//	    processRequest,
//	    readerioresult.WithTimeout[Response](5*time.Second),
//	    readerioresult.WithValue[Response](requestIDKey, requestID),
//	    readerioresult.WithValue[Response](userKey, user),
//	)
//
//	res := handler(ctx)() // the only place a context is supplied
//
// To derive a context value from an effect (loading a token, generating an ID), use
// LocalIOK or LocalIOResultK in readerioresult, together with reader.WithValue and
// reader.NopCancel to construct the new context:
//
//	addRequestID := func(ctx context.Context) io.IO[readerioresult.ContextCancel] {
//	    return func() readerioresult.ContextCancel {
//	        return reader.NopCancel(reader.WithValue[string](requestIDKey)(uuid.NewString())(ctx))
//	    }
//	}
//	scoped := readerioresult.LocalIOK[Response](addRequestID)(processRequest)
//
// # Cancellation
//
// Computations in these packages observe cancellation through the context they receive.
// Operators such as Delay, retries and resource brackets check ctx.Done() and return
// [context.Cause] as the error. WithContext / WithContextK short-circuit a computation if
// its context is already cancelled before it starts.
//
// # Guidelines
//
//   - Use an unexported, package-specific key type (type ctxKey string) rather than plain
//     strings, so keys from different packages cannot collide. AskValue distinguishes keys
//     by type as well as by value.
//   - Store request-scoped data only (IDs, principals, loggers, deadlines). Dependencies
//     such as database handles or configuration are better modelled as an explicit
//     environment, e.g. with the effect package or readerreaderioresult.
//   - Prefer AskValue over Asks with a manual type assertion, and WithValue / WithTimeout /
//     WithDeadline over a hand-written Local, so cancel functions are always released.
//   - Supply the context once, at the edge of the program (an HTTP handler, main, a test
//     via t.Context()), rather than capturing it in closures.
//
// # Sub-packages
//
//   - [github.com/IBM/fp-go/v2/context/reader]: pure functions of the context, plus the
//     WithValue / AskValue / NopCancel building blocks
//   - [github.com/IBM/fp-go/v2/context/readerio]: context-dependent side effects that cannot fail
//   - [github.com/IBM/fp-go/v2/context/readerresult]: context-dependent computations that can fail
//   - [github.com/IBM/fp-go/v2/context/readerioresult]: context-dependent side effects that can fail;
//     the most commonly used variant
//   - [github.com/IBM/fp-go/v2/context/readerreaderioresult]: an additional typed environment on
//     top of the context
//   - [github.com/IBM/fp-go/v2/context/statereaderioresult]: threads explicit state alongside the context
//   - [github.com/IBM/fp-go/v2/idiomatic/context/readerresult]: the same model using Go's native
//     (value, error) return convention
package context
