// package either implements the Either monad
//
// A data type that can be of either of two types but not both. This is
// typically used to carry an error or a return value
package either

import (
	E "github.com/ibm/fp-go/errors"
	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
)

func Of[E, A any](value A) Either[E, A] {
	return F.Pipe1(value, Right[E, A])
}

func FromIO[E, A any](f func() A) Either[E, A] {
	return F.Pipe1(f(), Right[E, A])
}

func MonadAp[E, A, B any](fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	return MonadFold(fab, Left[E, B], func(ab func(A) B) Either[E, B] {
		return MonadFold(fa, Left[E, B], F.Flow2(ab, Right[E, B]))
	})
}

func Ap[E, A, B any](fa Either[E, A]) func(fab Either[E, func(a A) B]) Either[E, B] {
	return F.Bind2nd(MonadAp[E, A, B], fa)
}

func MonadMap[E, A, B any](fa Either[E, A], f func(a A) B) Either[E, B] {
	return MonadChain(fa, F.Flow2(f, Right[E, B]))
}

func MonadBiMap[E1, E2, A, B any](fa Either[E1, A], f func(E1) E2, g func(a A) B) Either[E2, B] {
	return MonadFold(fa, F.Flow2(f, Left[E2, B]), F.Flow2(g, Right[E2, B]))
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(a A) B) func(Either[E1, A]) Either[E2, B] {
	return Fold(F.Flow2(f, Left[E2, B]), F.Flow2(g, Right[E2, B]))
}

func MonadMapTo[E, A, B any](fa Either[E, A], b B) Either[E, B] {
	return MonadMap(fa, F.Constant1[A](b))
}

func MapTo[E, A, B any](b B) func(Either[E, A]) Either[E, B] {
	return F.Bind2nd(MonadMapTo[E, A, B], b)
}

func MonadMapLeft[E, A, B any](fa Either[E, A], f func(E) B) Either[B, A] {
	return MonadFold(fa, F.Flow2(f, Left[B, A]), Right[B, A])
}

func Map[E, A, B any](f func(a A) B) func(fa Either[E, A]) Either[E, B] {
	return Chain(F.Flow2(f, Right[E, B]))
}

func MapLeft[E, A, B any](f func(E) B) func(fa Either[E, A]) Either[B, A] {
	return F.Bind2nd(MonadMapLeft[E, A, B], f)
}

func MonadChain[E, A, B any](fa Either[E, A], f func(a A) Either[E, B]) Either[E, B] {
	return MonadFold(fa, Left[E, B], f)
}

func MonadChainFirst[E, A, B any](ma Either[E, A], f func(a A) Either[E, B]) Either[E, A] {
	return MonadChain(ma, func(a A) Either[E, A] {
		return MonadMap(f(a), F.Constant1[B](a))
	})
}

func MonadChainTo[E, A, B any](ma Either[E, A], mb Either[E, B]) Either[E, B] {
	return mb
}

func MonadChainOptionK[E, A, B any](onNone func() E, ma Either[E, A], f func(A) O.Option[B]) Either[E, B] {
	return MonadChain(ma, F.Flow2(f, FromOption[E, B](onNone)))
}

func ChainOptionK[E, A, B any](onNone func() E) func(func(A) O.Option[B]) func(Either[E, A]) Either[E, B] {
	from := FromOption[E, B](onNone)
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

func TryCatch[E, A any](f func() (A, error), onThrow func(error) E) Either[E, A] {
	val, err := f()
	if err != nil {
		return F.Pipe2(err, onThrow, Left[E, A])
	}
	return F.Pipe1(val, Right[E, A])
}

func TryCatchErrorG[GA ~func() (A, error), A any](f GA) Either[error, A] {
	return TryCatch(f, E.IdentityError)
}

func TryCatchError[A any](f func() (A, error)) Either[error, A] {
	return TryCatchErrorG(f)
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

func FromOption[E, A any](onNone func() E) func(O.Option[A]) Either[E, A] {
	return O.Fold(F.Nullary2(onNone, Left[E, A]), Right[E, A])
}

func ToOption[E, A any]() func(Either[E, A]) O.Option[A] {
	return Fold(F.Ignore1[E](O.None[A]), O.Some[A])
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

func Eitherize0G[GA ~func() (R, error), GB ~func() Either[error, R], R any](f GA) GB {
	return F.Bind1(TryCatchErrorG[GA, R], f)
}

func Eitherize0[R any](f func() (R, error)) func() Either[error, R] {
	return Eitherize0G[func() (R, error), func() Either[error, R]](f)
}

func Uneitherize0G[GA ~func() Either[error, R], GB ~func() (R, error), R any](f GA) GB {
	return func() (R, error) {
		return UnwrapError(f())
	}
}

func Uneitherize0[R any](f func() Either[error, R]) func() (R, error) {
	return Uneitherize0G[func() Either[error, R], func() (R, error)](f)
}

func Eitherize1G[GA ~func(T1) (R, error), GB ~func(T1) Either[error, R], T1, R any](f GA) GB {
	return func(t1 T1) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1)
		})
	}
}

func Eitherize1[T1, R any](f func(T1) (R, error)) func(T1) Either[error, R] {
	return Eitherize1G[func(T1) (R, error), func(T1) Either[error, R]](f)
}

func Uneitherize1G[GA ~func(T1) Either[error, R], GB ~func(T1) (R, error), T1, R any](f GA) GB {
	return func(t1 T1) (R, error) {
		return UnwrapError(f(t1))
	}
}

func Uneitherize1[T1, R any](f func(T1) Either[error, R]) func(T1) (R, error) {
	return Uneitherize1G[func(T1) Either[error, R], func(T1) (R, error)](f)
}

func Eitherize2G[GA ~func(t1 T1, t2 T2) (R, error), GB ~func(t1 T1, t2 T2) Either[error, R], T1, T2, R any](f GA) GB {
	return func(t1 T1, t2 T2) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2)
		})
	}
}

func Eitherize2[T1, T2, R any](f func(t1 T1, t2 T2) (R, error)) func(t1 T1, t2 T2) Either[error, R] {
	return Eitherize2G[func(t1 T1, t2 T2) (R, error), func(t1 T1, t2 T2) Either[error, R]](f)
}

func Uneitherize2G[GA ~func(T1, T2) Either[error, R], GB ~func(T1, T2) (R, error), T1, T2, R any](f GA) GB {
	return func(t1 T1, t2 T2) (R, error) {
		return UnwrapError(f(t1, t2))
	}
}

func Uneitherize2[T1, T2, R any](f func(T1, T2) Either[error, R]) func(T1, T2) (R, error) {
	return Uneitherize2G[func(T1, T2) Either[error, R], func(T1, T2) (R, error)](f)
}

func Eitherize3G[GA ~func(t1 T1, t2 T2, t3 T3) (R, error), GB ~func(t1 T1, t2 T2, t3 T3) Either[error, R], T1, T2, T3, R any](f GA) GB {
	return func(t1 T1, t2 T2, t3 T3) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3)
		})
	}
}

func Eitherize3[T1, T2, T3, R any](f func(t1 T1, t2 T2, t3 T3) (R, error)) func(t1 T1, t2 T2, t3 T3) Either[error, R] {
	return Eitherize3G[func(t1 T1, t2 T2, t3 T3) (R, error), func(t1 T1, t2 T2, t3 T3) Either[error, R]](f)
}

func Uneitherize3G[GA ~func(T1, T2, T3) Either[error, R], GB ~func(T1, T2, T3) (R, error), T1, T2, T3, R any](f GA) GB {
	return func(t1 T1, t2 T2, t3 T3) (R, error) {
		return UnwrapError(f(t1, t2, t3))
	}
}

func Uneitherize3[T1, T2, T3, R any](f func(T1, T2, T3) Either[error, R]) func(T1, T2, T3) (R, error) {
	return Uneitherize3G[func(T1, T2, T3) Either[error, R], func(T1, T2, T3) (R, error)](f)
}

func Eitherize4G[GA ~func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error), GB ~func(t1 T1, t2 T2, t3 T3, t4 T4) Either[error, R], T1, T2, T3, T4, R any](f GA) GB {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) Either[error, R] {
		return TryCatchError(func() (R, error) {
			return f(t1, t2, t3, t4)
		})
	}
}

func Eitherize4[T1, T2, T3, T4, R any](f func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error)) func(t1 T1, t2 T2, t3 T3, t4 T4) Either[error, R] {
	return Eitherize4G[func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error), func(t1 T1, t2 T2, t3 T3, t4 T4) Either[error, R]](f)
}

func Uneitherize4G[GA ~func(T1, T2, T3, T4) Either[error, R], GB ~func(T1, T2, T3, T4) (R, error), T1, T2, T3, T4, R any](f GA) GB {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
		return UnwrapError(f(t1, t2, t3, t4))
	}
}

func Uneitherize4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) Either[error, R]) func(T1, T2, T3, T4) (R, error) {
	return Uneitherize4G[func(T1, T2, T3, T4) Either[error, R], func(T1, T2, T3, T4) (R, error)](f)
}

func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B {
	return func(ma Either[E, A]) B {
		return MonadFold(ma, onLeft, onRight)
	}
}

func UnwrapError[A any](ma Either[error, A]) (A, error) {
	return Unwrap[error](ma)
}

func FromPredicate[E, A any](pred func(A) bool, onFalse func(A) E) func(A) Either[E, A] {
	return func(a A) Either[E, A] {
		if pred(a) {
			return Right[E](a)
		}
		return Left[E, A](onFalse(a))
	}
}

func FromNillable[E, A any](e E) func(*A) Either[E, *A] {
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
	return Fold(F.Ignore1[E](that), Right[E1, A])
}

func Alt[E, A any](that func() Either[E, A]) func(Either[E, A]) Either[E, A] {
	return AltW[E](that)
}

func OrElse[E, A any](onLeft func(e E) Either[E, A]) func(Either[E, A]) Either[E, A] {
	return Fold(onLeft, Of[E, A])
}

func ToType[E, A any](onError func(any) E) func(any) Either[E, A] {
	return func(value any) Either[E, A] {
		return F.Pipe2(
			value,
			O.ToType[A],
			O.Fold(F.Nullary3(F.Constant(value), onError, Left[E, A]), Right[E, A]),
		)
	}
}

func Memoize[E, A any](val Either[E, A]) Either[E, A] {
	return val
}

func MonadSequence2[E, T1, T2, R any](e1 Either[E, T1], e2 Either[E, T2], f func(T1, T2) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[E, R], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[E, R], func(t2 T2) Either[E, R] {
			return f(t1, t2)
		})
	})
}

func MonadSequence3[E, T1, T2, T3, R any](e1 Either[E, T1], e2 Either[E, T2], e3 Either[E, T3], f func(T1, T2, T3) Either[E, R]) Either[E, R] {
	return MonadFold(e1, Left[E, R], func(t1 T1) Either[E, R] {
		return MonadFold(e2, Left[E, R], func(t2 T2) Either[E, R] {
			return MonadFold(e3, Left[E, R], func(t3 T3) Either[E, R] {
				return f(t1, t2, t3)
			})
		})
	})
}

// Swap changes the order of type parameters
func Swap[E, A any](val Either[E, A]) Either[A, E] {
	return MonadFold(val, Right[A, E], Left[A, E])
}
