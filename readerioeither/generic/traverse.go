package generic

import (
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
	RR "github.com/IBM/fp-go/internal/record"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() ET.Either[L, B], GIOBS ~func() ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, E, L, A, B any](ma AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, GIOBS, E, L, BBS],
		Map[GBS, func(E) func() ET.Either[L, func(B) BBS], GIOBS, func() ET.Either[L, func(B) BBS], E, L, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() ET.Either[L, func(B) BBS], GIOB, GIOBS, func() ET.Either[L, func(B) BBS], E, L, B, BBS],

		ma, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() ET.Either[L, B], GIOBS ~func() ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, E, L, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, GIOBS, E, L, BBS],
		Map[GBS, func(E) func() ET.Either[L, func(B) BBS], GIOBS, func() ET.Either[L, func(B) BBS], E, L, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() ET.Either[L, func(B) BBS], GIOB, GIOBS, func() ET.Either[L, func(B) BBS], E, L, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(E) GIOA, GAS ~func(E) GIOAS, GIOA ~func() ET.Either[L, A], GIOAS ~func() ET.Either[L, AAS], AAS ~[]A, GAAS ~[]GA, E, L, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}

// MonadTraverseRecord transforms an array
func MonadTraverseRecord[GB ~func(C) GIOB, GBS ~func(C) GIOBS, GIOB ~func() ET.Either[E, B], GIOBS ~func() ET.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, C, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse[AAS](
		Of[GBS, GIOBS, C, E, BBS],
		Map[GBS, func(C) func() ET.Either[E, func(B) BBS], GIOBS, func() ET.Either[E, func(B) BBS], C, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(C) func() ET.Either[E, func(B) BBS], GIOB, GIOBS, func() ET.Either[E, func(B) BBS], C, E, B, BBS],

		tas,
		f,
	)
}

// TraverseRecord transforms an array
func TraverseRecord[GB ~func(C) GIOB, GBS ~func(C) GIOBS, GIOB ~func() ET.Either[E, B], GIOBS ~func() ET.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, C, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, GIOBS, C, E, BBS],
		Map[GBS, func(C) func() ET.Either[E, func(B) BBS], GIOBS, func() ET.Either[E, func(B) BBS], C, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(C) func() ET.Either[E, func(B) BBS], GIOB, GIOBS, func() ET.Either[E, func(B) BBS], C, E, B, BBS],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[GA ~func(C) GIOA, GAS ~func(C) GIOAS, GIOA ~func() ET.Either[E, A], GIOAS ~func() ET.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, C, E, A any](tas GAAS) GAS {
	return MonadTraverseRecord[GA, GAS](tas, F.Identity[GA])
}
