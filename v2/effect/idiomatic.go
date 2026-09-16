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
	"context"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	IORI "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	OI "github.com/IBM/fp-go/v2/idiomatic/option"
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
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
	return readerreaderioresult.FromIdiomatic(f)
}

// MonadChainI sequences a [Effect] with an idiomatic Kleisli function.
// The success value of fa is passed to f; errors short-circuit the chain.
func MonadChainI[R, A, B any](fa Effect[R, A], f KleisliI[R, A, B]) Effect[R, B] {
	return readerreaderioresult.MonadChainI(fa, f)
}

// MonadChainFirstI sequences fa with an idiomatic Kleisli function f for its side
// effects, but returns the original value of fa (not the result of f).
// If fa or f fails, the error is propagated.
func MonadChainFirstI[R, A, B any](fa Effect[R, A], f KleisliI[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstI(fa, f)
}

// MonadTapI runs the idiomatic function f as a side-effect on the value of fa,
// discarding f's result and returning the original value of fa unchanged.
// Equivalent to [MonadChainFirstI].
func MonadTapI[R, A, B any](fa Effect[R, A], f KleisliI[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapI(fa, f)
}

// ChainI returns an [Operator] that sequences a computation with the idiomatic
// Kleisli function f, passing the success value forward.
// It is the curried form of [MonadChainI].
func ChainI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, B] {
	return readerreaderioresult.ChainI(f)
}

// ChainFirstI returns an [Operator] that runs the idiomatic function f for its
// side effects and then returns the original monadic value unchanged.
// It is the curried form of [MonadChainFirstI].
func ChainFirstI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstI(f)
}

// TapI returns an [Operator] that executes f as a side-effect, returning the
// original value unchanged. Equivalent to [ChainFirstI].
// It is the curried form of [MonadTapI].
func TapI[R, A, B any](f KleisliI[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.TapI(f)
}

// MonadChainLeftI recovers from an error in fa by running the idiomatic function f
// on the error value. If fa succeeds its value is returned unchanged; if fa fails,
// f is called with the error and its result replaces the failed computation.
func MonadChainLeftI[R, A any](fa Effect[R, A], f KleisliI[R, error, A]) Effect[R, A] {
	return readerreaderioresult.MonadChainLeftI(fa, f)
}

// ChainLeftI returns a function that recovers from errors using the idiomatic
// function f. If the input computation succeeds, its value passes through; if it fails,
// f is called with the error to produce a recovery computation.
// It is the curried form of [MonadChainLeftI].
func ChainLeftI[R, A any](f KleisliI[R, error, A]) Operator[R, A, A] {
	return readerreaderioresult.ChainLeftI(f)
}

// RetryingI retries the idiomatic action according to policy until check returns
// false or the policy is exhausted. It is a convenience wrapper around [Retrying] that
// accepts an idiomatic action instead of a [Kleisli].
func RetryingI[R, A any](
	policy retry.RetryPolicy,
	action KleisliI[R, retry.RetryStatus, A],
	check Predicate[Result[A]],
) Effect[R, A] {
	return readerreaderioresult.RetryingI(policy, action, check)
}

// TraverseArrayI maps each element of a slice through the idiomatic function f
// and collects the results into a [Effect] holding a slice.
// The first error encountered short-circuits the traversal.
//
// For an empty input array the resulting array is nil, the canonical empty array.
func TraverseArrayI[R, A, B any](f KleisliI[R, A, B]) Kleisli[R, []A, []B] {
	return readerreaderioresult.TraverseArrayI(f)
}

// BindI is the idiomatic-function variant of [Bind].
// It attaches the result of the idiomatic Kleisli f to a do-notation context using setter.
func BindI[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f KleisliI[R, S1, T],
) Operator[R, S1, S2] {
	return readerreaderioresult.BindI(setter, f)
}

// BindIL is the lens-based, idiomatic-function variant of [BindL].
// It reads the focused field T from the context S via lens, passes it to f,
// and writes the result back through the same lens.
func BindIL[R, S, T any](
	lens Lens[S, T],
	f KleisliI[R, T, T],
) Operator[R, S, S] {
	return readerreaderioresult.BindIL(lens, f)
}

// ---------------------------------------------------------------------------
// Adapters for the sub-monads
//
// [FromIdiomatic] and the functions above bridge the whole stack at once: a
// plain Go method that takes a context.Context plus the environment R and
// returns (value, error). The functions below bridge one layer at a time, for
// code that only needs part of the stack - a pure fallible computation
// (Result), a deferred one (IOResult), one that reads the environment
// (ReaderResult), or a context-aware effectful one (Thunk).
// ---------------------------------------------------------------------------

// FromResultI lifts an idiomatic Go (value, error) pair into an [Effect] that
// ignores both the environment and the context.
// If err is non-nil the effect always fails with that error, otherwise it always
// succeeds with a.
//
// Example:
//
//	value, err := strconv.Atoi("42")
//	eff := effect.FromResultI[Config](value, err)
//
//go:inline
func FromResultI[C, A any](a A, err error) Effect[C, A] {
	return readerreaderioresult.FromResultI[C](a, err)
}

// FromIOResultI converts an idiomatic IOResult (a `func() (A, error)`) into an
// [Effect] that ignores both the environment and the context.
//
//go:inline
func FromIOResultI[C, A any](mr IORI.IOResult[A]) Effect[C, A] {
	return readerreaderioresult.FromIOResultI[C](mr)
}

// FromReaderResultI converts an idiomatic ReaderResult (a `func(R) (A, error)`)
// into an [Effect]. The computation reads the environment but neither sees the
// context nor performs IO.
//
// Example:
//
//	endpoint := func(cfg Config) (string, error) {
//	    if cfg.URL == "" {
//	        return "", errors.New("no url configured")
//	    }
//	    return cfg.URL, nil
//	}
//	eff := effect.FromReaderResultI(endpoint)
//
//go:inline
func FromReaderResultI[R, A any](rr RRI.ReaderResult[R, A]) Effect[R, A] {
	return readerreaderioresult.FromReaderResultI(rr)
}

// FromThunkI converts an idiomatic Thunk (a `func(context.Context) func() (A, error)`)
// into an [Effect]. The computation sees the context and performs IO, but does not
// read the environment.
//
// Example:
//
//	get := func(ctx context.Context) func() (*http.Response, error) {
//	    return func() (*http.Response, error) { ... }
//	}
//	eff := effect.FromThunkI[Config](get)
//
//go:inline
func FromThunkI[C, A any](f func(context.Context) func() (A, error)) Effect[C, A] {
	return FromThunk[C](thunk.FromIdiomatic(f))
}

// fromThunkKleisliI converts an idiomatic Thunk Kleisli arrow
// `func(A) func(context.Context) func() (B, error)` into `thunk.Kleisli[A, B]`.
func fromThunkKleisliI[A, B any](f thunk.KleisliI[A, B]) thunk.Kleisli[A, B] {
	return function.Flow2(f, thunk.FromIdiomatic[B])
}

// MonadChainResultIK is the idiomatic version of [MonadChainResultK].
// It chains an [Effect] with a plain Go function returning (B, error).
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	eff := effect.MonadChainResultIK(effect.Of[Config]("42"), parse)
//
//go:inline
func MonadChainResultIK[C, A, B any](ma Effect[C, A], f RI.Kleisli[A, B]) Effect[C, B] {
	return readerreaderioresult.MonadChainResultIK(ma, f)
}

// ChainResultIK is the curried version of [MonadChainResultIK].
// It is the idiomatic counterpart of [ChainResultK]: it lifts a plain Go function
// returning (B, error) into an [Operator].
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	eff := F.Pipe1(effect.Of[Config]("42"), effect.ChainResultIK[Config](parse))
//
//go:inline
func ChainResultIK[C, A, B any](f RI.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainResultIK[C](f)
}

// MonadChainFirstResultIK runs an idiomatic (B, error) function for its side
// effects and keeps the original value.
//
//go:inline
func MonadChainFirstResultIK[C, A, B any](ma Effect[C, A], f RI.Kleisli[A, B]) Effect[C, A] {
	return readerreaderioresult.MonadChainFirstResultIK(ma, f)
}

// MonadTapResultIK is an alias for [MonadChainFirstResultIK].
//
//go:inline
func MonadTapResultIK[C, A, B any](ma Effect[C, A], f RI.Kleisli[A, B]) Effect[C, A] {
	return readerreaderioresult.MonadTapResultIK(ma, f)
}

// ChainFirstResultIK is the curried version of [MonadChainFirstResultIK].
// It is the shape a validation step has: run the check, keep the value.
//
// Example:
//
//	validate := func(s string) (int, error) { return strconv.Atoi(s) }
//	eff := F.Pipe1(effect.Of[Config]("42"), effect.ChainFirstResultIK[Config](validate))
//
//go:inline
func ChainFirstResultIK[C, A, B any](f RI.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirstResultIK[C](f)
}

// TapResultIK is an alias for [ChainFirstResultIK].
//
//go:inline
func TapResultIK[C, A, B any](f RI.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.TapResultIK[C](f)
}

// MonadChainReaderResultIK is the idiomatic version of chaining a computation
// that reads the environment and may fail but performs no IO.
//
//go:inline
func MonadChainReaderResultIK[R, A, B any](ma Effect[R, A], f RRI.Kleisli[R, A, B]) Effect[R, B] {
	return readerreaderioresult.MonadChainReaderResultIK(ma, f)
}

// ChainReaderResultIK is the curried version of [MonadChainReaderResultIK].
// It lifts a plain Go function returning `func(R) (B, error)` into an [Operator].
//
// Example:
//
//	lookup := func(key string) func(Config) (string, error) { ... }
//	eff := F.Pipe1(effect.Of[Config]("host"), effect.ChainReaderResultIK(lookup))
//
//go:inline
func ChainReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return readerreaderioresult.ChainReaderResultIK(f)
}

// MonadChainFirstReaderResultIK runs an idiomatic environment-reading function
// for its side effects and keeps the original value.
//
//go:inline
func MonadChainFirstReaderResultIK[R, A, B any](ma Effect[R, A], f RRI.Kleisli[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstReaderResultIK(ma, f)
}

// MonadTapReaderResultIK is an alias for [MonadChainFirstReaderResultIK].
//
//go:inline
func MonadTapReaderResultIK[R, A, B any](ma Effect[R, A], f RRI.Kleisli[R, A, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapReaderResultIK(ma, f)
}

// ChainFirstReaderResultIK is the curried version of [MonadChainFirstReaderResultIK].
//
//go:inline
func ChainFirstReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstReaderResultIK(f)
}

// TapReaderResultIK is an alias for [ChainFirstReaderResultIK].
//
//go:inline
func TapReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return readerreaderioresult.TapReaderResultIK(f)
}

// MonadChainIOResultIK is the idiomatic version of chaining a deferred
// `func() (B, error)` that neither reads the environment nor sees the context.
//
//go:inline
func MonadChainIOResultIK[C, A, B any](ma Effect[C, A], f IORI.Kleisli[A, B]) Effect[C, B] {
	return readerreaderioresult.MonadChainIOResultIK(ma, f)
}

// ChainIOResultIK is the curried version of [MonadChainIOResultIK].
//
// Example:
//
//	read := func(path string) func() ([]byte, error) {
//	    return func() ([]byte, error) { return os.ReadFile(path) }
//	}
//	eff := F.Pipe1(effect.Of[Config]("f.txt"), effect.ChainIOResultIK[Config](read))
//
//go:inline
func ChainIOResultIK[C, A, B any](f IORI.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainIOResultIK[C](f)
}

// MonadChainThunkIK is the idiomatic version of chaining a [Thunk]-returning
// function: it sees the context and may fail, but does not read the environment.
//
//go:inline
func MonadChainThunkIK[C, A, B any](ma Effect[C, A], f thunk.KleisliI[A, B]) Effect[C, B] {
	return MonadChainI(ma, func(a A) func(context.Context, C) (B, error) {
		mb := f(a)
		return func(ctx context.Context, _ C) (B, error) {
			return mb(ctx)()
		}
	})
}

// ChainThunkIK is the idiomatic version of [ChainThunkK].
// It chains a plain Go function of the form `func(A) func(context.Context) func() (B, error)`,
// the shape of a context-aware call that may fail.
//
// Example:
//
//	fetch := func(url string) func(context.Context) func() ([]byte, error) { ... }
//	eff := F.Pipe1(effect.Of[Config]("https://example.com"), effect.ChainThunkIK[Config](fetch))
//
//go:inline
func ChainThunkIK[C, A, B any](f thunk.KleisliI[A, B]) Operator[C, A, B] {
	return ChainThunkK[C](fromThunkKleisliI(f))
}

// ChainFirstThunkIK is the idiomatic version of [ChainFirstThunkK].
// It runs the context-aware, fallible function for its side effects and keeps the
// original value.
//
//go:inline
func ChainFirstThunkIK[C, A, B any](f thunk.KleisliI[A, B]) Operator[C, A, A] {
	return ChainFirstThunkK[C, A](fromThunkKleisliI(f))
}

// TapThunkIK is an alias for [ChainFirstThunkIK].
//
//go:inline
func TapThunkIK[C, A, B any](f thunk.KleisliI[A, B]) Operator[C, A, A] {
	return TapThunkK[C, A](fromThunkKleisliI(f))
}

// ChainFirstLeftThunkIK is the idiomatic version of [ChainFirstLeftThunkK].
// It runs a context-aware handler on the error path and preserves the original
// value or error.
//
// Example:
//
//	report := func(err error) func(context.Context) func() (string, error) { ... }
//	eff := F.Pipe1(failing, effect.ChainFirstLeftThunkIK[Config, int](report))
//
//go:inline
func ChainFirstLeftThunkIK[C, A, B any](f thunk.KleisliI[error, B]) Operator[C, A, A] {
	return ChainFirstLeftThunkK[C, A](fromThunkKleisliI(f))
}

// TapLeftThunkIK is an alias for [ChainFirstLeftThunkIK].
//
//go:inline
func TapLeftThunkIK[C, A, B any](f thunk.KleisliI[error, B]) Operator[C, A, A] {
	return TapLeftThunkK[C, A](fromThunkKleisliI(f))
}

// ChainOptionIK chains a plain Go function returning the comma-ok pair (B, bool),
// converting a false flag into an error produced by onNone.
//
// Example:
//
//	lookup := func(k string) (string, bool) { v, ok := m[k]; return v, ok }
//	notFound := func() error { return errors.New("not found") }
//	chain := effect.ChainOptionIK[Config, string, string](notFound)
//	eff := F.Pipe1(effect.Of[Config]("host"), chain(lookup))
//
//go:inline
func ChainOptionIK[C, A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainOptionIK[C, A, B](onNone)
}

// ChainFirstIOResultIK runs an idiomatic, deferred `func(A) func() (B, error)` for
// its side effects and keeps the original value.
//
//go:inline
func ChainFirstIOResultIK[C, A, B any](f IORI.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirstIOResultIK[C, A, B](f)
}

// TapIOResultIK is an alias for [ChainFirstIOResultIK].
//
//go:inline
func TapIOResultIK[C, A, B any](f IORI.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.TapIOResultIK[C, A, B](f)
}
