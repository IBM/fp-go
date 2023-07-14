package eithert

import (
	ET "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/apply"
	FC "github.com/ibm/fp-go/internal/functor"
)

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadMap[E, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E, A]) ET.Either[E, B]) HKTFB, fa HKTFA, f func(A) B) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return FC.MonadMap(fmap, ET.MonadMap[E, A, B], fa, f)
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadBiMap[E1, E2, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E1, A]) ET.Either[E2, B]) HKTFB, fa HKTFA, f func(E1) E2, g func(A) B) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return fmap(fa, ET.BiMap(f, g))
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func BiMap[E1, E2, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E1, A]) ET.Either[E2, B]) HKTFB, f func(E1) E2, g func(A) B) func(HKTFA) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return F.Bind2nd(fmap, ET.BiMap(f, g))
}

// HKTFA = HKT<F, Either<E, A>>
// HKTFB = HKT<F, Either<E, B>>
func MonadChain[E, A, B, HKTFA, HKTFB any](
	fchain func(HKTFA, func(ET.Either[E, A]) HKTFB) HKTFB,
	fof func(ET.Either[E, B]) HKTFB,
	ma HKTFA,
	f func(A) HKTFB) HKTFB {
	// dispatch to the even more generic implementation
	return fchain(ma, ET.Fold(F.Flow2(ET.Left[B, E], fof), f))
}

// func(fa func(R) T.Task[ET.Either[E, func(A) B]], f func(ET.Either[E, func(A) B]) func(ET.Either[E, A]) ET.Either[E, B]) GEFAB

// HKTFA   = HKT[Either[E, A]]
// HKTFB   = HKT[Either[E, B]]
// HKTFAB  = HKT[Either[E, func(A)B]]
func MonadAp[E, A, B, HKTFAB, HKTFGAB, HKTFA, HKTFB any](
	fap func(HKTFGAB, HKTFA) HKTFB,
	fmap func(HKTFAB, func(ET.Either[E, func(A) B]) func(ET.Either[E, A]) ET.Either[E, B]) HKTFGAB,
	fab HKTFAB,
	fa HKTFA) HKTFB {
	// HKTGA  = ET.Either[E, A]
	// HKTGB  = ET.Either[E, B]
	// HKTGAB = ET.Either[E, func(a A) B]
	return apply.MonadAp(fap, fmap, ET.MonadAp[B, E, A], fab, fa)
}

func Right[E, A, HKTA any](fof func(ET.Either[E, A]) HKTA, a A) HKTA {
	return F.Pipe2(a, ET.Right[E, A], fof)
}

func Left[E, A, HKTA any](fof func(ET.Either[E, A]) HKTA, e E) HKTA {
	return F.Pipe2(e, ET.Left[A, E], fof)
}

// HKTA  = HKT[A]
// HKTEA = HKT[Either[E, A]]
func RightF[E, A, HKTA, HKTEA any](fmap func(HKTA, func(A) ET.Either[E, A]) HKTEA, fa HKTA) HKTEA {
	return fmap(fa, ET.Right[E, A])
}

// HKTE  = HKT[E]
// HKTEA = HKT[Either[E, A]]
func LeftF[E, A, HKTE, HKTEA any](fmap func(HKTE, func(E) ET.Either[E, A]) HKTEA, fe HKTE) HKTEA {
	return fmap(fe, ET.Left[A, E])
}

func FoldE[E, A, HKTEA, HKTB any](mchain func(HKTEA, func(ET.Either[E, A]) HKTB) HKTB, ma HKTEA, onLeft func(E) HKTB, onRight func(A) HKTB) HKTB {
	return mchain(ma, ET.Fold(onLeft, onRight))
}

func MatchE[E, A, HKTEA, HKTB any](mchain func(HKTEA, func(ET.Either[E, A]) HKTB) HKTB, onLeft func(E) HKTB, onRight func(A) HKTB) func(HKTEA) HKTB {
	return F.Bind2nd(mchain, ET.Fold(onLeft, onRight))
}

func GetOrElse[E, A, HKTEA, HKTA any](mchain func(HKTEA, func(ET.Either[E, A]) HKTA) HKTA, mof func(A) HKTA, onLeft func(E) HKTA) func(HKTEA) HKTA {
	return MatchE(mchain, onLeft, mof)
}

func OrElse[E1, E2, A, HKTE1A, HKTE2A any](mchain func(HKTE1A, func(ET.Either[E1, A]) HKTE2A) HKTE2A, mof func(ET.Either[E2, A]) HKTE2A, onLeft func(E1) HKTE2A) func(HKTE1A) HKTE2A {
	return MatchE(mchain, onLeft, F.Flow2(ET.Right[E2, A], mof))
}

func OrLeft[E1, E2, A, HKTE1A, HKTE2, HKTE2A any](
	mchain func(HKTE1A, func(ET.Either[E1, A]) HKTE2A) HKTE2A,
	mmap func(HKTE2, func(E2) ET.Either[E2, A]) HKTE2A,
	mof func(ET.Either[E2, A]) HKTE2A,
	onLeft func(E1) HKTE2) func(HKTE1A) HKTE2A {

	return F.Bind2nd(mchain, ET.Fold(F.Flow2(onLeft, F.Bind2nd(mmap, ET.Left[A, E2])), F.Flow2(ET.Right[E2, A], mof)))
}

func MonadMapLeft[E, A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(ET.Either[E, A]) ET.Either[B, A]) HKTFB, fa HKTFA, f func(E) B) HKTFB {
	return FC.MonadMap(fmap, ET.MonadMapLeft[E, A, B], fa, f)
}
