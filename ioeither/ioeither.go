// Copyright (c) 2023 IBM Corp.
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

package ioeither

import (
	"time"

	ET "github.com/IBM/fp-go/either"
	I "github.com/IBM/fp-go/io"
	G "github.com/IBM/fp-go/ioeither/generic"
	IOO "github.com/IBM/fp-go/iooption"
	L "github.com/IBM/fp-go/lazy"
	O "github.com/IBM/fp-go/option"
)

// IOEither represents a synchronous computation that may fail
// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details
type IOEither[E, A any] I.IO[ET.Either[E, A]]

func MakeIO[E, A any](f IOEither[E, A]) IOEither[E, A] {
	return G.MakeIO(f)
}

func Left[A, E any](l E) IOEither[E, A] {
	return G.Left[IOEither[E, A]](l)
}

func Right[E, A any](r A) IOEither[E, A] {
	return G.Right[IOEither[E, A]](r)
}

func Of[E, A any](r A) IOEither[E, A] {
	return G.Of[IOEither[E, A]](r)
}

func MonadOf[E, A any](r A) IOEither[E, A] {
	return G.MonadOf[IOEither[E, A]](r)
}

func LeftIO[A, E any](ml I.IO[E]) IOEither[E, A] {
	return G.LeftIO[IOEither[E, A]](ml)
}

func RightIO[E, A any](mr I.IO[A]) IOEither[E, A] {
	return G.RightIO[IOEither[E, A]](mr)
}

func FromEither[E, A any](e ET.Either[E, A]) IOEither[E, A] {
	return G.FromEither[IOEither[E, A]](e)
}

func FromOption[A, E any](onNone func() E) func(o O.Option[A]) IOEither[E, A] {
	return G.FromOption[IOEither[E, A]](onNone)
}

func FromIOOption[A, E any](onNone func() E) func(o IOO.IOOption[A]) IOEither[E, A] {
	return G.FromIOOption[IOEither[E, A], IOO.IOOption[A]](onNone)
}

func ChainOptionK[A, B, E any](onNone func() E) func(func(A) O.Option[B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ChainOptionK[IOEither[E, A], IOEither[E, B]](onNone)
}

func MonadChainIOK[E, A, B any](ma IOEither[E, A], f func(A) I.IO[B]) IOEither[E, B] {
	return G.MonadChainIOK[IOEither[E, A], IOEither[E, B]](ma, f)
}

func ChainIOK[E, A, B any](f func(A) I.IO[B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ChainIOK[IOEither[E, A], IOEither[E, B]](f)
}

func ChainLazyK[E, A, B any](f func(A) L.Lazy[B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ChainIOK[IOEither[E, A], IOEither[E, B]](f)
}

// FromIO creates an [IOEither] from an [IO] instance, invoking [IO] for each invocation of [IOEither]
func FromIO[E, A any](mr I.IO[A]) IOEither[E, A] {
	return G.FromIO[IOEither[E, A]](mr)
}

// FromLazy creates an [IOEither] from a [Lazy] instance, invoking [Lazy] for each invocation of [IOEither]
func FromLazy[E, A any](mr L.Lazy[A]) IOEither[E, A] {
	return G.FromIO[IOEither[E, A]](mr)
}

func MonadMap[E, A, B any](fa IOEither[E, A], f func(A) B) IOEither[E, B] {
	return G.MonadMap[IOEither[E, A], IOEither[E, B]](fa, f)
}

func Map[E, A, B any](f func(A) B) func(IOEither[E, A]) IOEither[E, B] {
	return G.Map[IOEither[E, A], IOEither[E, B]](f)
}

func MonadMapTo[E, A, B any](fa IOEither[E, A], b B) IOEither[E, B] {
	return G.MonadMapTo[IOEither[E, A], IOEither[E, B]](fa, b)
}

func MapTo[E, A, B any](b B) func(IOEither[E, A]) IOEither[E, B] {
	return G.MapTo[IOEither[E, A], IOEither[E, B]](b)
}

func MonadChain[E, A, B any](fa IOEither[E, A], f func(A) IOEither[E, B]) IOEither[E, B] {
	return G.MonadChain(fa, f)
}

func Chain[E, A, B any](f func(A) IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.Chain[IOEither[E, A]](f)
}

func MonadChainEitherK[E, A, B any](ma IOEither[E, A], f func(A) ET.Either[E, B]) IOEither[E, B] {
	return G.MonadChainEitherK[IOEither[E, A], IOEither[E, B]](ma, f)
}

func ChainEitherK[E, A, B any](f func(A) ET.Either[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ChainEitherK[IOEither[E, A], IOEither[E, B]](f)
}

func MonadAp[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return G.MonadAp[IOEither[E, B]](mab, ma)
}

// Ap is an alias of [ApPar]
func Ap[B, E, A any](ma IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B] {
	return G.Ap[IOEither[E, B], IOEither[E, func(A) B]](ma)
}

func MonadApPar[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return G.MonadApPar[IOEither[E, B]](mab, ma)
}

// ApPar applies function and value in parallel
func ApPar[B, E, A any](ma IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B] {
	return G.ApPar[IOEither[E, B], IOEither[E, func(A) B]](ma)
}

func MonadApSeq[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return G.MonadApSeq[IOEither[E, B]](mab, ma)
}

// ApSeq applies function and value sequentially
func ApSeq[B, E, A any](ma IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B] {
	return G.ApSeq[IOEither[E, B], IOEither[E, func(A) B]](ma)
}

func Flatten[E, A any](mma IOEither[E, IOEither[E, A]]) IOEither[E, A] {
	return G.Flatten(mma)
}

func TryCatch[E, A any](f func() (A, error), onThrow func(error) E) IOEither[E, A] {
	return G.TryCatch[IOEither[E, A]](f, onThrow)
}

func TryCatchError[A any](f func() (A, error)) IOEither[error, A] {
	return G.TryCatchError[IOEither[error, A]](f)
}

func Memoize[E, A any](ma IOEither[E, A]) IOEither[E, A] {
	return G.Memoize(ma)
}

func MonadMapLeft[E1, E2, A any](fa IOEither[E1, A], f func(E1) E2) IOEither[E2, A] {
	return G.MonadMapLeft[IOEither[E1, A], IOEither[E2, A]](fa, f)
}

func MapLeft[A, E1, E2 any](f func(E1) E2) func(IOEither[E1, A]) IOEither[E2, A] {
	return G.MapLeft[IOEither[E1, A], IOEither[E2, A]](f)
}

func MonadBiMap[E1, E2, A, B any](fa IOEither[E1, A], f func(E1) E2, g func(A) B) IOEither[E2, B] {
	return G.MonadBiMap[IOEither[E1, A], IOEither[E2, B]](fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(IOEither[E1, A]) IOEither[E2, B] {
	return G.BiMap[IOEither[E1, A], IOEither[E2, B]](f, g)
}

// Fold converts an IOEither into an IO
func Fold[E, A, B any](onLeft func(E) I.IO[B], onRight func(A) I.IO[B]) func(IOEither[E, A]) I.IO[B] {
	return G.Fold[IOEither[E, A]](onLeft, onRight)
}

// GetOrElse extracts the value or maps the error
func GetOrElse[E, A any](onLeft func(E) I.IO[A]) func(IOEither[E, A]) I.IO[A] {
	return G.GetOrElse[IOEither[E, A]](onLeft)
}

// MonadChainTo composes to the second monad ignoring the return value of the first
func MonadChainTo[A, E, B any](fa IOEither[E, A], fb IOEither[E, B]) IOEither[E, B] {
	return G.MonadChainTo(fa, fb)
}

// ChainTo composes to the second [IOEither] monad ignoring the return value of the first
func ChainTo[A, E, B any](fb IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ChainTo[IOEither[E, A]](fb)
}

// MonadChainFirst runs the [IOEither] monad returned by the function but returns the result of the original monad
func MonadChainFirst[E, A, B any](ma IOEither[E, A], f func(A) IOEither[E, B]) IOEither[E, A] {
	return G.MonadChainFirst(ma, f)
}

// ChainFirst runs the [IOEither] monad returned by the function but returns the result of the original monad
func ChainFirst[E, A, B any](f func(A) IOEither[E, B]) func(IOEither[E, A]) IOEither[E, A] {
	return G.ChainFirst[IOEither[E, A]](f)
}

func MonadChainFirstEitherK[A, E, B any](ma IOEither[E, A], f func(A) ET.Either[E, B]) IOEither[E, A] {
	return G.MonadChainFirstEitherK[IOEither[E, A]](ma, f)
}

func ChainFirstEitherK[A, E, B any](f func(A) ET.Either[E, B]) func(ma IOEither[E, A]) IOEither[E, A] {
	return G.ChainFirstEitherK[IOEither[E, A]](f)
}

// MonadChainFirstIOK runs [IO] the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[E, A, B any](ma IOEither[E, A], f func(A) I.IO[B]) IOEither[E, A] {
	return G.MonadChainFirstIOK(ma, f)
}

// ChainFirstIOK runs the [IO] monad returned by the function but returns the result of the original monad
func ChainFirstIOK[E, A, B any](f func(A) I.IO[B]) func(IOEither[E, A]) IOEither[E, A] {
	return G.ChainFirstIOK[IOEither[E, A]](f)
}

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[A, E, R, ANY any](onCreate IOEither[E, R], onRelease func(R) IOEither[E, ANY]) func(func(R) IOEither[E, A]) IOEither[E, A] {
	return G.WithResource[IOEither[E, A]](onCreate, onRelease)
}

// Swap changes the order of type parameters
func Swap[E, A any](val IOEither[E, A]) IOEither[A, E] {
	return G.Swap[IOEither[E, A], IOEither[A, E]](val)
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[E any](f func()) IOEither[E, any] {
	return G.FromImpure[IOEither[E, any]](f)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[E, A any](gen L.Lazy[IOEither[E, A]]) IOEither[E, A] {
	return G.Defer[IOEither[E, A]](gen)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[E, A any](first IOEither[E, A], second L.Lazy[IOEither[E, A]]) IOEither[E, A] {
	return G.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor
func Alt[E, A any](second L.Lazy[IOEither[E, A]]) func(IOEither[E, A]) IOEither[E, A] {
	return G.Alt(second)
}

// OrElse returns the original IOEither if it is a Right, otherwise it applies the given function to the error and returns the result.
func OrElse[E, A any](onLeft func(E) IOEither[E, A]) func(IOEither[E, A]) IOEither[E, A] {
	return G.OrElse[IOEither[E, A]](onLeft)
}

func MonadFlap[E, B, A any](fab IOEither[E, func(A) B], a A) IOEither[E, B] {
	return G.MonadFlap[IOEither[E, func(A) B], IOEither[E, B]](fab, a)
}

func Flap[E, B, A any](a A) func(IOEither[E, func(A) B]) IOEither[E, B] {
	return G.Flap[IOEither[E, func(A) B], IOEither[E, B]](a)
}

// ToIOOption converts an [IOEither] to an [IOO.IOOption]
func ToIOOption[E, A any](ioe IOEither[E, A]) IOO.IOOption[A] {
	return G.ToIOOption[IOO.IOOption[A]](ioe)
}

// Delay creates an operation that passes in the value after some delay
func Delay[E, A any](delay time.Duration) func(IOEither[E, A]) IOEither[E, A] {
	return G.Delay[IOEither[E, A]](delay)
}

// After creates an operation that passes after the given [time.Time]
func After[E, A any](timestamp time.Time) func(IOEither[E, A]) IOEither[E, A] {
	return G.After[IOEither[E, A]](timestamp)
}
