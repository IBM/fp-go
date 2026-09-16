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

	"github.com/IBM/fp-go/v2/function"
	IORI "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	OI "github.com/IBM/fp-go/v2/idiomatic/option"
	RIORI "github.com/IBM/fp-go/v2/idiomatic/readerioresult"
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/readerresult"
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
	return function.Flow2(f, fromReaderReaderIOResultI[R, B])
}

// fromReaderReaderIOResultI converts an idiomatic `func(context.Context, R) (A, error)`
// into a [ReaderReaderIOResult] by swapping the argument order and currying.
func fromReaderReaderIOResultI[R, A any](fa func(context.Context, R) (A, error)) ReaderReaderIOResult[R, A] {
	return function.Curry2(function.Swap(ioresult.Eitherize2(fa)))
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
//
// For an empty input array the resulting array is nil, the canonical empty array.
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

// ---------------------------------------------------------------------------
// Adapters for the sub-monads
//
// The functions below bridge the idiomatic (value, error) convention at each
// individual layer of the stack: a plain Result, an IOResult, a ReaderResult
// over the outer environment R, or the full ReaderIOResult. They complement
// [FromIdiomatic], which bridges the whole stack at once.
// ---------------------------------------------------------------------------

// fromResultKleisliI converts an idiomatic Kleisli arrow `func(A) (B, error)` into
// the functional `result.Kleisli[A, B]`.
func fromResultKleisliI[A, B any](f RI.Kleisli[A, B]) result.Kleisli[A, B] {
	return result.Eitherize1(f)
}

// fromOptionKleisliI converts an idiomatic Kleisli arrow `func(A) (B, bool)` into
// the functional `option.Kleisli[A, B]`.
func fromOptionKleisliI[A, B any](f OI.Kleisli[A, B]) option.Kleisli[A, B] {
	return option.Optionize1(f)
}

// fromIOResultKleisliI converts an idiomatic Kleisli arrow `func(A) func() (B, error)`
// into the functional `ioresult.Kleisli[A, B]`.
func fromIOResultKleisliI[A, B any](f IORI.Kleisli[A, B]) ioresult.Kleisli[A, B] {
	return function.Flow2(f, ioresult.TryCatchError[B])
}

// fromReaderResultKleisliI converts an idiomatic Kleisli arrow `func(A) func(R) (B, error)`
// into the functional `readerresult.Kleisli[R, A, B]`.
func fromReaderResultKleisliI[R, A, B any](f RRI.Kleisli[R, A, B]) readerresult.Kleisli[R, A, B] {
	return function.Flow2(f, readerresult.FromReaderResultI[R, B])
}

// FromResultI lifts an idiomatic Go (value, error) pair into a [ReaderReaderIOResult]
// that ignores both the environment and the context.
// If err is non-nil the computation always fails with that error, otherwise it
// always succeeds with a.
//
// Example:
//
//	value, err := strconv.Atoi("42")
//	eff := readerreaderioresult.FromResultI[Config](value, err)
//
//go:inline
func FromResultI[R, A any](a A, err error) ReaderReaderIOResult[R, A] {
	return FromResult[R](result.TryCatchError(a, err))
}

// FromIOResultI converts an idiomatic IOResult (a `func() (A, error)`) into a
// [ReaderReaderIOResult] that ignores both the environment and the context.
//
// Example:
//
//	readStdin := func() (string, error) { ... }
//	eff := readerreaderioresult.FromIOResultI[Config](readStdin)
//
//go:inline
func FromIOResultI[R, A any](mr IORI.IOResult[A]) ReaderReaderIOResult[R, A] {
	return FromIOResult[R](ioresult.TryCatchError(mr))
}

// FromReaderResultI converts an idiomatic ReaderResult (a `func(R) (A, error)`)
// into a [ReaderReaderIOResult]. The computation reads the outer environment but
// neither sees the context nor performs IO.
//
// Example:
//
//	endpoint := func(cfg Config) (string, error) {
//	    if cfg.URL == "" {
//	        return "", errors.New("no url configured")
//	    }
//	    return cfg.URL, nil
//	}
//	eff := readerreaderioresult.FromReaderResultI(endpoint)
//
//go:inline
func FromReaderResultI[R, A any](rr RRI.ReaderResult[R, A]) ReaderReaderIOResult[R, A] {
	return FromReaderResult(readerresult.FromReaderResultI(rr))
}

// FromReaderIOResultI converts an idiomatic ReaderIOResult (a
// `func(R) func() (A, error)`) into a [ReaderReaderIOResult]. The computation
// reads the outer environment and performs IO, but does not see the context.
//
//go:inline
func FromReaderIOResultI[R, A any](rr RIORI.ReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return FromReaderIOResult(readerioresult.FromReaderIOResultI(rr))
}

// MonadChainEitherIK is the idiomatic version of [MonadChainEitherK].
// It chains a [ReaderReaderIOResult] with a plain Go function returning (B, error).
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	eff := readerreaderioresult.MonadChainEitherIK(readerreaderioresult.Of[Config]("42"), parse)
//
//go:inline
func MonadChainEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainEitherK(ma, fromResultKleisliI(f))
}

// MonadChainResultIK is an alias for [MonadChainEitherIK] with more explicit naming.
//
//go:inline
func MonadChainResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainEitherK(ma, fromResultKleisliI(f))
}

// ChainEitherIK is the curried version of [MonadChainEitherIK].
// It lifts a plain Go function returning (B, error) into an [Operator].
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	eff := F.Pipe1(readerreaderioresult.Of[Config]("42"), readerreaderioresult.ChainEitherIK[Config](parse))
//
//go:inline
func ChainEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainEitherK[R](fromResultKleisliI(f))
}

// ChainResultIK is an alias for [ChainEitherIK] with more explicit naming.
// It is the idiomatic counterpart of [ChainResultK].
//
//go:inline
func ChainResultIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainEitherK[R](fromResultKleisliI(f))
}

// MonadChainFirstEitherIK is the idiomatic version of [MonadChainFirstEitherK].
// It runs a plain Go function returning (B, error) for its side effects and keeps
// the original value.
//
//go:inline
func MonadChainFirstEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, fromResultKleisliI(f))
}

// MonadChainFirstResultIK is an alias for [MonadChainFirstEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, fromResultKleisliI(f))
}

// MonadTapEitherIK is an alias for [MonadChainFirstEitherIK].
//
//go:inline
func MonadTapEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadTapEitherK(ma, fromResultKleisliI(f))
}

// MonadTapResultIK is an alias for [MonadChainFirstResultIK].
//
//go:inline
func MonadTapResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadTapEitherK(ma, fromResultKleisliI(f))
}

// ChainFirstEitherIK is the curried version of [MonadChainFirstEitherIK].
//
//go:inline
func ChainFirstEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstEitherK[R](fromResultKleisliI(f))
}

// ChainFirstResultIK is an alias for [ChainFirstEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstResultIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstEitherK[R](fromResultKleisliI(f))
}

// TapEitherIK is an alias for [ChainFirstEitherIK].
//
//go:inline
func TapEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, A] {
	return TapEitherK[R](fromResultKleisliI(f))
}

// TapResultIK is an alias for [ChainFirstResultIK].
//
//go:inline
func TapResultIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, A] {
	return TapEitherK[R](fromResultKleisliI(f))
}

// MonadChainReaderEitherIK is the idiomatic version of [MonadChainReaderEitherK].
// It chains a plain Go function returning `func(R) (B, error)`, i.e. a computation
// that reads the outer environment and may fail but performs no IO.
//
// Example:
//
//	lookup := func(key string) func(Config) (string, error) { ... }
//	eff := readerreaderioresult.MonadChainReaderEitherIK(readerreaderioresult.Of[Config]("host"), lookup)
//
//go:inline
func MonadChainReaderEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadChainReaderResultIK is an alias for [MonadChainReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainReaderResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// ChainReaderEitherIK is the curried version of [MonadChainReaderEitherIK].
//
// Example:
//
//	lookup := func(key string) func(Config) (string, error) { ... }
//	eff := F.Pipe1(readerreaderioresult.Of[Config]("host"), readerreaderioresult.ChainReaderEitherIK(lookup))
//
//go:inline
func ChainReaderEitherIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return ChainReaderEitherK(fromReaderResultKleisliI(f))
}

// ChainReaderResultIK is an alias for [ChainReaderEitherIK] with more explicit naming.
//
//go:inline
func ChainReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return ChainReaderEitherK(fromReaderResultKleisliI(f))
}

// MonadChainFirstReaderEitherIK is the idiomatic version of [MonadChainFirstReaderEitherK].
//
//go:inline
func MonadChainFirstReaderEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadChainFirstReaderResultIK is an alias for [MonadChainFirstReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstReaderResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadTapReaderEitherIK is an alias for [MonadChainFirstReaderEitherIK].
//
//go:inline
func MonadTapReaderEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadTapReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadTapReaderResultIK is an alias for [MonadChainFirstReaderResultIK].
//
//go:inline
func MonadTapReaderResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadTapReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// ChainFirstReaderEitherIK is the curried version of [MonadChainFirstReaderEitherIK].
//
//go:inline
func ChainFirstReaderEitherIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderEitherK(fromReaderResultKleisliI(f))
}

// ChainFirstReaderResultIK is an alias for [ChainFirstReaderEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderEitherK(fromReaderResultKleisliI(f))
}

// TapReaderEitherIK is an alias for [ChainFirstReaderEitherIK].
//
//go:inline
func TapReaderEitherIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return TapReaderEitherK(fromReaderResultKleisliI(f))
}

// TapReaderResultIK is an alias for [ChainFirstReaderResultIK].
//
//go:inline
func TapReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, A] {
	return TapReaderEitherK(fromReaderResultKleisliI(f))
}

// MonadChainIOEitherIK is the idiomatic version of [MonadChainIOEitherK].
// It chains a plain Go function returning a deferred `func() (B, error)`.
//
// Example:
//
//	read := func(path string) func() ([]byte, error) {
//	    return func() ([]byte, error) { return os.ReadFile(path) }
//	}
//	eff := readerreaderioresult.MonadChainIOEitherIK(readerreaderioresult.Of[Config]("f.txt"), read)
//
//go:inline
func MonadChainIOEitherIK[R, A, B any](ma ReaderReaderIOResult[R, A], f IORI.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainIOEitherK(ma, fromIOResultKleisliI(f))
}

// MonadChainIOResultIK is an alias for [MonadChainIOEitherIK] with more explicit naming.
//
//go:inline
func MonadChainIOResultIK[R, A, B any](ma ReaderReaderIOResult[R, A], f IORI.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return MonadChainIOEitherK(ma, fromIOResultKleisliI(f))
}

// ChainIOEitherIK is the curried version of [MonadChainIOEitherIK].
//
//go:inline
func ChainIOEitherIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainIOEitherK[R](fromIOResultKleisliI(f))
}

// ChainIOResultIK is an alias for [ChainIOEitherIK] with more explicit naming.
//
//go:inline
func ChainIOResultIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainIOEitherK[R](fromIOResultKleisliI(f))
}

// ChainOptionIK is the idiomatic version of [ChainOptionK].
// It chains a plain Go function returning the comma-ok pair (B, bool), converting a
// false flag into an error produced by onNone.
//
// Example:
//
//	lookup := func(k string) (string, bool) { v, ok := m[k]; return v, ok }
//	notFound := func() error { return errors.New("not found") }
//	chain := readerreaderioresult.ChainOptionIK[Config, string, string](notFound)
//	eff := F.Pipe1(readerreaderioresult.Of[Config]("host"), chain(lookup))
//
//go:inline
func ChainOptionIK[R, A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[R, A, B] {
	return function.Flow2(fromOptionKleisliI[A, B], ChainOptionK[R, A, B](onNone))
}

// ChainFirstIOEitherIK runs an idiomatic, deferred `func(A) func() (B, error)` for
// its side effects and keeps the original value.
//
// Note that this package has no functional ChainFirstIOEitherK counterpart; the
// operator is built from [ChainFirst] and [FromIOResult].
//
//go:inline
func ChainFirstIOEitherIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirst(function.Flow2(fromIOResultKleisliI(f), FromIOResult[R, B]))
}

// ChainFirstIOResultIK is an alias for [ChainFirstIOEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstIOResultIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOEitherIK[R, A, B](f)
}

// TapIOEitherIK is an alias for [ChainFirstIOEitherIK].
//
//go:inline
func TapIOEitherIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return Tap(function.Flow2(fromIOResultKleisliI(f), FromIOResult[R, B]))
}

// TapIOResultIK is an alias for [ChainFirstIOResultIK].
//
//go:inline
func TapIOResultIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return TapIOEitherIK[R, A, B](f)
}
