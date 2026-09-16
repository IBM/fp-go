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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/function"
	IORI "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	OI "github.com/IBM/fp-go/v2/idiomatic/option"
	RIORI "github.com/IBM/fp-go/v2/idiomatic/readerioresult"
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

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

// fromReaderIOResultKleisliI converts an idiomatic Kleisli arrow
// `func(A) func(R) func() (B, error)` into the functional [Kleisli].
func fromReaderIOResultKleisliI[R, A, B any](f RIORI.Kleisli[R, A, B]) Kleisli[R, A, B] {
	return function.Flow2(f, FromReaderIOResultI[R, B])
}

// FromIdiomatic lifts a plain Go function of the form `func(R) func() (A, error)` into a
// [ReaderIOResult]. The environment R is forwarded to f, and the deferred `(A, error)`
// pair produced by the returned thunk is converted into a [Result].
//
// This is the primary entry-point for wrapping existing Go code that follows the
// idiomatic (value, error) return convention but needs to stay lazy.
//
// Example:
//
//	readFile := func(cfg Config) func() ([]byte, error) {
//	    return func() ([]byte, error) { return os.ReadFile(cfg.Path) }
//	}
//	rio := readerioresult.FromIdiomatic(readFile)
//
// See Also:
//   - [FromReaderIOResultI]: identical conversion, named after the RIORI.ReaderIOResult alias
//   - [FromResultI]: lifts an already-evaluated (A, error) pair
//
//go:inline
func FromIdiomatic[R, A any](f func(R) func() (A, error)) ReaderIOResult[R, A] {
	return FromReaderIOResultI[R, A](f)
}

// FromResultI lifts an idiomatic Go (value, error) pair into a [ReaderIOResult] that
// ignores the environment. This is the idiomatic counterpart of [FromResult].
// If err is non-nil the computation always fails with that error, otherwise it always
// succeeds with a.
//
// Example:
//
//	value, err := strconv.Atoi("42")
//	rio := readerioresult.FromResultI[Config](value, err)
//
//go:inline
func FromResultI[R, A any](a A, err error) ReaderIOResult[R, A] {
	return FromResult[R](result.TryCatchError(a, err))
}

// FromIOResultI converts an idiomatic IOResult (a `func() (A, error)`) into a
// [ReaderIOResult] that ignores the environment.
//
// Example:
//
//	readStdin := func() (string, error) { ... }
//	rio := readerioresult.FromIOResultI[Config](readStdin)
//
//go:inline
func FromIOResultI[R, A any](mr IORI.IOResult[A]) ReaderIOResult[R, A] {
	return FromIOResult[R](ioresult.TryCatchError(mr))
}

// FromReaderResultI converts an idiomatic ReaderResult (a `func(R) (A, error)`) into a
// [ReaderIOResult]. The resulting computation is pure with respect to IO: f is invoked
// when the IO action is executed.
//
// Example:
//
//	getUserID := func(cfg Config) (int, error) {
//	    if cfg.Valid {
//	        return 42, nil
//	    }
//	    return 0, errors.New("invalid config")
//	}
//	rio := readerioresult.FromReaderResultI(getUserID)
//
//go:inline
func FromReaderResultI[R, A any](rr RRI.ReaderResult[R, A]) ReaderIOResult[R, A] {
	return FromReaderEither(readerresult.FromReaderResultI(rr))
}

// FromReaderIOResultI converts an idiomatic ReaderIOResult (a `func(R) func() (A, error)`)
// into the functional [ReaderIOResult] (a `func(R) func() Result[A]`).
// This bridges Go's native error handling and the functional style.
//
// Example:
//
//	fetch := func(cfg Config) func() (*http.Response, error) {
//	    return func() (*http.Response, error) { return http.Get(cfg.URL) }
//	}
//	rio := readerioresult.FromReaderIOResultI(fetch)
func FromReaderIOResultI[R, A any](rr RIORI.ReaderIOResult[R, A]) ReaderIOResult[R, A] {
	return function.Flow2(rr, ioresult.TryCatchError[A])
}

// MonadChainI sequences a [ReaderIOResult] with an idiomatic Kleisli arrow
// `func(A) func(R) func() (B, error)`. The success value of ma is passed to f;
// errors short-circuit the chain.
//
// Example:
//
//	load := func(id int) func(Config) func() (User, error) { ... }
//	rio := readerioresult.MonadChainI(readerioresult.Of[Config](42), load)
//
//go:inline
func MonadChainI[R, A, B any](ma ReaderIOResult[R, A], f RIORI.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return MonadChain(ma, fromReaderIOResultKleisliI(f))
}

// ChainI is the curried version of [MonadChainI].
// It lifts an idiomatic Kleisli arrow into a [ReaderIOResult] [Operator].
//
// Example:
//
//	load := func(id int) func(Config) func() (User, error) { ... }
//	rio := F.Pipe1(readerioresult.Of[Config](42), readerioresult.ChainI(load))
//
//go:inline
func ChainI[R, A, B any](f RIORI.Kleisli[R, A, B]) Operator[R, A, B] {
	return Chain(fromReaderIOResultKleisliI(f))
}

// MonadChainFirstI sequences ma with an idiomatic Kleisli arrow for its side effects
// but returns the original value of ma. If ma or f fails the error is propagated.
//
//go:inline
func MonadChainFirstI[R, A, B any](ma ReaderIOResult[R, A], f RIORI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirst(ma, fromReaderIOResultKleisliI(f))
}

// ChainFirstI is the curried version of [MonadChainFirstI].
//
//go:inline
func ChainFirstI[R, A, B any](f RIORI.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirst(fromReaderIOResultKleisliI(f))
}

// MonadTapI is an alias for [MonadChainFirstI], emphasizing the side-effect nature
// of the operation.
//
//go:inline
func MonadTapI[R, A, B any](ma ReaderIOResult[R, A], f RIORI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadTap(ma, fromReaderIOResultKleisliI(f))
}

// TapI is an alias for [ChainFirstI], emphasizing the side-effect nature of the operation.
//
//go:inline
func TapI[R, A, B any](f RIORI.Kleisli[R, A, B]) Operator[R, A, A] {
	return Tap(fromReaderIOResultKleisliI(f))
}

// MonadChainEitherIK is the idiomatic version of [MonadChainEitherK].
// It chains a [ReaderIOResult] with a plain Go function returning (B, error).
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	rio := readerioresult.MonadChainEitherIK(readerioresult.Of[Config]("42"), parse)
//
//go:inline
func MonadChainEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, B] {
	return MonadChainEitherK(ma, fromResultKleisliI(f))
}

// MonadChainResultIK is an alias for [MonadChainEitherIK] with more explicit naming.
//
//go:inline
func MonadChainResultIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, B] {
	return MonadChainEitherK(ma, fromResultKleisliI(f))
}

// ChainEitherIK is the curried version of [MonadChainEitherIK].
// It lifts a plain Go function returning (B, error) into a [ReaderIOResult] [Operator].
//
// Example:
//
//	parse := func(s string) (int, error) { return strconv.Atoi(s) }
//	rio := F.Pipe1(readerioresult.Of[Config]("42"), readerioresult.ChainEitherIK[Config](parse))
//
//go:inline
func ChainEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainEitherK[R](fromResultKleisliI(f))
}

// ChainResultIK is an alias for [ChainEitherIK] with more explicit naming.
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
func MonadChainFirstEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, fromResultKleisliI(f))
}

// MonadChainFirstResultIK is an alias for [MonadChainFirstEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstResultIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, fromResultKleisliI(f))
}

// MonadTapEitherIK is an alias for [MonadChainFirstEitherIK].
//
//go:inline
func MonadTapEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, A] {
	return MonadTapEitherK(ma, fromResultKleisliI(f))
}

// MonadTapResultIK is an alias for [MonadChainFirstResultIK].
//
//go:inline
func MonadTapResultIK[R, A, B any](ma ReaderIOResult[R, A], f RI.Kleisli[A, B]) ReaderIOResult[R, A] {
	return MonadTapEitherK(ma, fromResultKleisliI(f))
}

// ChainFirstEitherIK is the curried version of [MonadChainFirstEitherIK].
//
// Example:
//
//	validate := func(x int) (string, error) { ... }
//	rio := F.Pipe1(readerioresult.Of[Config](42), readerioresult.ChainFirstEitherIK[Config](validate))
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

// MonadChainIOEitherIK is the idiomatic version of [MonadChainIOEitherK].
// It chains a plain Go function returning a deferred `func() (B, error)`.
//
// Example:
//
//	read := func(path string) func() ([]byte, error) {
//	    return func() ([]byte, error) { return os.ReadFile(path) }
//	}
//	rio := readerioresult.MonadChainIOEitherIK(readerioresult.Of[Config]("f.txt"), read)
//
//go:inline
func MonadChainIOEitherIK[R, A, B any](ma ReaderIOResult[R, A], f IORI.Kleisli[A, B]) ReaderIOResult[R, B] {
	return MonadChainIOEitherK(ma, fromIOResultKleisliI(f))
}

// MonadChainIOResultIK is an alias for [MonadChainIOEitherIK] with more explicit naming.
//
//go:inline
func MonadChainIOResultIK[R, A, B any](ma ReaderIOResult[R, A], f IORI.Kleisli[A, B]) ReaderIOResult[R, B] {
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

// ChainFirstIOEitherIK is the idiomatic version of [ChainFirstIOEitherK].
// It runs a deferred `func() (B, error)` for its side effects and keeps the original value.
//
//go:inline
func ChainFirstIOEitherIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOEitherK[R](fromIOResultKleisliI(f))
}

// ChainFirstIOResultIK is an alias for [ChainFirstIOEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstIOResultIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOEitherK[R](fromIOResultKleisliI(f))
}

// TapIOEitherIK is an alias for [ChainFirstIOEitherIK].
//
//go:inline
func TapIOEitherIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return TapIOEitherK[R](fromIOResultKleisliI(f))
}

// TapIOResultIK is an alias for [ChainFirstIOResultIK].
//
//go:inline
func TapIOResultIK[R, A, B any](f IORI.Kleisli[A, B]) Operator[R, A, A] {
	return TapIOEitherK[R](fromIOResultKleisliI(f))
}

// MonadChainReaderEitherIK is the idiomatic version of [MonadChainReaderEitherK].
// It chains a plain Go function returning `func(R) (B, error)`, i.e. a computation that
// reads the environment and may fail but performs no IO.
//
// Example:
//
//	lookup := func(key string) func(Config) (string, error) {
//	    return func(cfg Config) (string, error) { ... }
//	}
//	rio := readerioresult.MonadChainReaderEitherIK(readerioresult.Of[Config]("host"), lookup)
//
//go:inline
func MonadChainReaderEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return MonadChainReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadChainReaderResultIK is an alias for [MonadChainReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainReaderResultIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return MonadChainReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// ChainReaderEitherIK is the curried version of [MonadChainReaderEitherIK].
//
// Example:
//
//	lookup := func(key string) func(Config) (string, error) { ... }
//	rio := F.Pipe1(readerioresult.Of[Config]("host"), readerioresult.ChainReaderEitherIK(lookup))
//
//go:inline
func ChainReaderEitherIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return ChainReaderEitherK(fromReaderResultKleisliI(f))
}

// ChainReaderResultIK is an alias for [ChainReaderEitherIK] with more explicit naming.
// It is the idiomatic counterpart of [ChainReaderResultK].
//
//go:inline
func ChainReaderResultIK[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return ChainReaderEitherK(fromReaderResultKleisliI(f))
}

// MonadChainFirstReaderEitherIK is the idiomatic version of [MonadChainFirstReaderEitherK].
//
//go:inline
func MonadChainFirstReaderEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadChainFirstReaderResultIK is an alias for [MonadChainFirstReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstReaderResultIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadTapReaderEitherIK is an alias for [MonadChainFirstReaderEitherIK].
//
//go:inline
func MonadTapReaderEitherIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return MonadTapReaderEitherK(ma, fromReaderResultKleisliI(f))
}

// MonadTapReaderResultIK is an alias for [MonadChainFirstReaderResultIK].
//
//go:inline
func MonadTapReaderResultIK[R, A, B any](ma ReaderIOResult[R, A], f RRI.Kleisli[R, A, B]) ReaderIOResult[R, A] {
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

// ChainOptionIK is the idiomatic version of [ChainOptionK].
// It chains a plain Go function returning the comma-ok pair (B, bool), converting a
// false flag into an error produced by onNone.
//
// Example:
//
//	lookup := func(k string) (string, bool) { v, ok := m[k]; return v, ok }
//	notFound := func() error { return errors.New("not found") }
//	chain := readerioresult.ChainOptionIK[Config, string, string](notFound)
//	rio := F.Pipe1(readerioresult.Of[Config]("host"), chain(lookup))
//
//go:inline
func ChainOptionIK[R, A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[R, A, B] {
	return function.Flow2(fromOptionKleisliI[A, B], ChainOptionK[R, A, B](onNone))
}
