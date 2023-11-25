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

// package option implements the Option monad, a data type that can have a defined value or none
package option

import (
	F "github.com/IBM/fp-go/function"
	FC "github.com/IBM/fp-go/internal/functor"
)

func fromPredicate[A any](a A, pred func(A) bool) Option[A] {
	if pred(a) {
		return Some(a)
	}
	return None[A]()
}

func FromPredicate[A any](pred func(A) bool) func(A) Option[A] {
	return F.Bind2nd(fromPredicate[A], pred)
}

func FromNillable[A any](a *A) Option[*A] {
	return fromPredicate(a, F.IsNonNil[A])
}

func FromValidation[A, B any](f func(A) (B, bool)) func(A) Option[B] {
	return Optionize1(f)
}

// MonadAp is the applicative functor of Option
func MonadAp[B, A any](fab Option[func(A) B], fa Option[A]) Option[B] {
	return MonadFold(fab, None[B], func(ab func(A) B) Option[B] {
		return MonadFold(fa, None[B], F.Flow2(ab, Some[B]))
	})
}

// Ap is the applicative functor of Option
func Ap[B, A any](fa Option[A]) func(Option[func(A) B]) Option[B] {
	return F.Bind2nd(MonadAp[B, A], fa)
}

func MonadMap[A, B any](fa Option[A], f func(A) B) Option[B] {
	return MonadChain(fa, F.Flow2(f, Some[B]))
}

func Map[A, B any](f func(a A) B) func(Option[A]) Option[B] {
	return Chain(F.Flow2(f, Some[B]))
}

func MonadMapTo[A, B any](fa Option[A], b B) Option[B] {
	return MonadMap(fa, F.Constant1[A](b))
}

func MapTo[A, B any](b B) func(Option[A]) Option[B] {
	return F.Bind2nd(MonadMapTo[A, B], b)
}

func TryCatch[A any](f func() (A, error)) Option[A] {
	val, err := f()
	if err != nil {
		return None[A]()
	}
	return Some(val)
}

func Fold[A, B any](onNone func() B, onSome func(a A) B) func(ma Option[A]) B {
	return func(ma Option[A]) B {
		return MonadFold(ma, onNone, onSome)
	}
}

func MonadGetOrElse[A any](fa Option[A], onNone func() A) A {
	return MonadFold(fa, onNone, F.Identity[A])
}

func GetOrElse[A any](onNone func() A) func(Option[A]) A {
	return Fold(onNone, F.Identity[A])
}

func MonadChain[A, B any](fa Option[A], f func(A) Option[B]) Option[B] {
	return MonadFold(fa, None[B], f)
}

func Chain[A, B any](f func(A) Option[B]) func(Option[A]) Option[B] {
	return F.Bind2nd(MonadChain[A, B], f)
}

func MonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B] {
	return mb
}

func ChainTo[A, B any](mb Option[B]) func(Option[A]) Option[B] {
	return F.Bind2nd(MonadChainTo[A, B], mb)
}

func MonadChainFirst[A, B any](ma Option[A], f func(A) Option[B]) Option[A] {
	return MonadChain(ma, func(a A) Option[A] {
		return MonadMap(f(a), F.Constant1[B](a))
	})
}

func ChainFirst[A, B any](f func(A) Option[B]) func(Option[A]) Option[A] {
	return F.Bind2nd(MonadChainFirst[A, B], f)
}

func Flatten[A any](mma Option[Option[A]]) Option[A] {
	return MonadChain(mma, F.Identity[Option[A]])
}

func MonadAlt[A any](fa Option[A], that func() Option[A]) Option[A] {
	return MonadFold(fa, that, Of[A])
}

func Alt[A any](that func() Option[A]) func(Option[A]) Option[A] {
	return Fold(that, Of[A])
}

func MonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R] {
	return MonadFold(o1, None[R], func(t1 T1) Option[R] {
		return MonadFold(o2, None[R], func(t2 T2) Option[R] {
			return f(t1, t2)
		})
	})
}

func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R] {
	return func(o1 Option[T1], o2 Option[T2]) Option[R] {
		return MonadSequence2(o1, o2, f)
	}
}

func Reduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B {
	return Fold(F.Constant(initial), F.Bind1st(f, initial))
}

// Filter converts an optional onto itself if it is some and the predicate is true
func Filter[A any](pred func(A) bool) func(Option[A]) Option[A] {
	return Fold(None[A], F.Ternary(pred, Of[A], F.Ignore1of1[A](None[A])))
}

func MonadFlap[B, A any](fab Option[func(A) B], a A) Option[B] {
	return FC.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

func Flap[B, A any](a A) func(Option[func(A) B]) Option[B] {
	return F.Bind2nd(MonadFlap[B, A], a)
}
