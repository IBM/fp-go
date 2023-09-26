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

// package either implements the Either monad
//
// A data type that can be of either of two types but not both. This is
// typically used to carry an error or a return value
package either

import (
	E "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	FC "github.com/IBM/fp-go/internal/functor"
	O "github.com/IBM/fp-go/option"
)

// Of is equivalent to [Right]
func Of[E, A any](value A) Either[E, A] {
	return F.Pipe1(value, Right[E, A])
}

func FromIO[E, IO ~func() A, A any](f IO) Either[E, A] {
	return F.Pipe1(f(), Right[E, A])
}

func MonadAp[B, E, A any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	return MonadFold(fab, Left[B, E], func(ab func(A) B) Either[E, B] {
		return MonadFold(fa, Left[B, E], F.Flow2(ab, Right[E, B]))
	})
}

func Ap[B, E, A any](fa Either[E, A]) func(fab Either[E, func(a A) B]) Either[E, B] {
	return F.Bind2nd(MonadAp[B, E, A], fa)
}

func MonadMap[E, A, B any](fa Either[E, A], f func(a A) B) Either[E, B] {
	return MonadChain(fa, F.Flow2(f, Right[E, B]))
}

func MonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B] {
	return MonadFold(fa, F.Flow2(f, Left[B, E2]), F.Flow2(g, Right[E2, B]))
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B] {
	return Fold(F.Flow2(f, Left[B, E2]), F.Flow2(g, Right[E2, B]))
}

func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	return MonadMap(fa, F.Constant1[A](b))
}

func MapTo[E, A, B any](b B) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadMapTo[E, A, B], b)
}

func MonadMapLeft[E, A, B any](fa Either[E, A], f func(E) B) Either[B, A] {
	return MonadFold(fa, F.Flow2(f, Left[A, B]), Right[B, A])
}

func Map[E, A, B any](f func(a A) B) func(fa Either[E, A]) Either[E, B] {
	return Chain(F.Flow2(f, Right[E, B]))
}

func MapLeft[E, A, B any](f func(E) B) func(fa Either[E, A]) Either[B, A] {
	return F.Bind2nd(MonadMapLeft[E, A, B], f)
}

func MonadChain[E, A, B any](fa Either[E, A], f func(a A) Either[E, B]) Either[E, B] {
	return MonadFold(fa, Left[B, E], f)
}

func MonadChainFirst[E, A, B any](ma Either[E, A], f func(a A) Either[E, B]) Either[E, A] {
	return MonadChain(ma, func(a A) Either[E, A] {
		return MonadMap(f(a), F.Constant1[B](a))
	})
}

func MonadChainTo[E, A, B any](ma Either[E, A], mb Either[E, B]) Either[E, B] {
	return mb
}

func MonadChainOptionK[A, B, E any](onNone func() E, ma Either[E, A], f func(A) O.Option[B]) Either[E, B] {
	return MonadChain(ma, F.Flow2(f, FromOption[B](onNone)))
}

func ChainOptionK[A, B, E any](onNone func() E) func(func(A) O.Option[B]) func(Either[E, A]) Either[E, B] {
	from := FromOption[B](onNone)
	return func(f func(A) O.Option[B]) func(Either[E, A]) Either[E, B] {
		return Chain(F.Flow2(f, from))
	}
}

func ChainTo[E, A, B any](mb Either[E, B]) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadChainTo[E, A, B], mb)
}

func Chain[E, A, B any](f func(a A) Either[E, B]) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadChain[E, A, B], f)
}

func ChainFirst[E, A, B any](f func(a A) Either[E, B]) func(Either[E, A]) Either[E, A] {
	return F.Bind2nd(MonadChainFirst[E, A, B], f)
}

func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A] {
	return MonadChain(mma, F.Identity[Either[E, A]])
}

func TryCatch[FA ~func() (A, error), FE func(error) E, E, A any](f FA, onThrow FE) Either[E, A] {
	val, err := f()
	if err != nil {
		return F.Pipe2(err, onThrow, Left[A, E])
	}
	return F.Pipe1(val, Right[E, A])
}

func TryCatchError[F ~func() (A, error), A any](f F) Either[error, A] {
	return TryCatch(f, E.IdentityError)
}

func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2]) Either[E, R] {
		return MonadSequence2(e1, e2, f)
	}
}

func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(Either[E, T1], Either[E, T2], Either[E, T3]) Either[E, R] {
	return func(e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3]) Either[E, R] {
		return MonadSequence3(e1, e2, e3, f)
	}
}

func FromOption[A, E any](onNone func() E) func(O.Option[A]) Either[E, A] {
	return O.Fold(F.Nullary2(onNone, Left[A, E]), Right[E, A])
}

func ToOption[E, A any](ma Either[E, A]) O.Option[A] {
	return MonadFold(ma, F.Ignore1of1[E](O.None[A]), O.Some[A])
}

func FromError[A any](f func(a A) error) func(A) Either[error, A] {
	return func(a A) Either[error, A] {
		return TryCatchError(func() (A, error) {
			return a, f(a)
		})
	}
}

func ToError[A any](e Either[error, A]) error {
	return MonadFold(e, E.IdentityError, F.Constant1[A, error](nil))
}

func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B {
	return func(ma Either[E, A]) B {
		return MonadFold(ma, onLeft, onRight)
	}
}

// UnwrapError converts an Either into the idiomatic tuple
func UnwrapError[A any](ma Either[error, A]) (A, error) {
	return Unwrap[error](ma)
}

func FromPredicate[E, A any](pred func(A) bool, onFalse func(A) E) func(A) Either[E, A] {
	return func(a A) Either[E, A] {
		if pred(a) {
			return Right[E](a)
		}
		return Left[A, E](onFalse(a))
	}
}

func FromNillable[A, E any](e E) func(*A) Either[E, *A] {
	return FromPredicate(F.IsNonNil[A], F.Constant1[*A](e))
}

func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A {
	return Fold(onLeft, F.Identity[A])
}

func Reduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B {
	return Fold(
		F.Constant1[E](initial),
		F.Bind1st(f, initial),
	)
}

func AltW[E, E1, A any](that func() Either[E1, A]) func(Either[E, A]) Either[E1, A] {
	return Fold(F.Ignore1of1[E](that), Right[E1, A])
}

func Alt[E, A any](that func() Either[E, A]) func(Either[E, A]) Either[E, A] {
	return AltW[E](that)
}

func OrElse[E, A any](onLeft func(e E) Either[E, A]) func(Either[E, A]) Either[E, A] {
	return Fold(onLeft, Of[E, A])
}

func ToType[A, E any](onError func(any) E) func(any) Either[E, A] {
	return func(value any) Either[E, A] {
		return F.Pipe2(
			value,
			O.ToType[A],
			O.Fold(F.Nullary3(F.Constant(value), onError, Left[A, E]), Right[E, A]),
		)
	}
}

func Memoize[E, A any](val Either[E, A]) Either[E, A] {
	return val
}

func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[R, E], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[R, E], func(t2 T2) Either[E, R] {
			return f(t1, t2)
		})
	})
}

func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[R, E], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[R, E], func(t2 T2) Either[E, R] {
			return MonadFold(e3, Left[R, E], func(t3 T3) Either[E, R] {
				return f(t1, t2, t3)
			})
		})
	})
}

// Swap changes the order of type parameters
func Swap[E, A any](val Either[E, A]) Either[A, E] {
	return MonadFold(val, Right[A, E], Left[E, A])
}

func MonadFlap[E, A, B any](fab Either[E, func(A) B], a A) Either[E, B] {
	return FC.MonadFlap(MonadMap[E, func(A) B, B], fab, a)
}

func Flap[E, A, B any](a A) func(Either[E, func(A) B]) Either[E, B] {
	return F.Bind2nd(MonadFlap[E, A, B], a)
}
