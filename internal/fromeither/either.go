package fromeither

import (
	ET "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	C "github.com/ibm/fp-go/internal/chain"
	O "github.com/ibm/fp-go/option"
)

func FromOption[E, A, HKTEA any](fromEither func(ET.Either[E, A]) HKTEA, onNone func() E) func(ma O.Option[A]) HKTEA {
	return F.Flow2(ET.FromOption[E, A](onNone), fromEither)
}

func FromPredicate[E, A, HKTEA any](fromEither func(ET.Either[E, A]) HKTEA, pred func(A) bool, onFalse func(A) E) func(A) HKTEA {
	return F.Flow2(ET.FromPredicate(pred, onFalse), fromEither)
}

func MonadFromOption[E, A, HKTEA any](
	fromEither func(ET.Either[E, A]) HKTEA,
	onNone func() E,
	ma O.Option[A],
) HKTEA {
	return F.Pipe1(
		O.MonadFold(
			ma,
			F.Nullary2(onNone, ET.Left[A, E]),
			ET.Right[E, A],
		),
		fromEither,
	)
}

func FromOptionK[E, A, B, HKTEB any](
	fromEither func(ET.Either[E, B]) HKTEB,
	onNone func() E) func(f func(A) O.Option[B]) func(A) HKTEB {
	// helper
	return F.Bind2nd(F.Flow2[func(A) O.Option[B], func(O.Option[B]) HKTEB, A, O.Option[B], HKTEB], FromOption(fromEither, onNone))
}

func MonadChainEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEB) HKTEB,
	fromEither func(ET.Either[E, B]) HKTEB,
	ma HKTEA,
	f func(A) ET.Either[E, B]) HKTEB {
	return mchain(ma, F.Flow2(f, fromEither))
}

func ChainOptionK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEB) HKTEB,
	fromEither func(ET.Either[E, B]) HKTEB,
	onNone func() E,
) func(f func(A) O.Option[B]) func(ma HKTEA) HKTEB {
	return F.Flow2(FromOptionK[E, A](fromEither, onNone), F.Bind1st(F.Bind2nd[HKTEA, func(A) HKTEB, HKTEB], mchain))
}

func MonadChainFirstEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEA) HKTEA,
	mmap func(HKTEB, func(B) A) HKTEA,
	fromEither func(ET.Either[E, B]) HKTEB,
	ma HKTEA,
	f func(A) ET.Either[E, B]) HKTEA {
	return C.MonadChainFirst(mchain, mmap, ma, F.Flow2(f, fromEither))
}

func ChainFirstEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEA) HKTEA,
	mmap func(HKTEB, func(B) A) HKTEA,
	fromEither func(ET.Either[E, B]) HKTEB,
	f func(A) ET.Either[E, B]) func(HKTEA) HKTEA {
	return C.ChainFirst(mchain, mmap, F.Flow2(f, fromEither))
}
