package generic

import (
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() B, GIOBS ~func() BBS, AAS ~[]A, BBS ~[]B, E, A, B any](ma AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, GIOBS, E, BBS],
		Map[GBS, func(E) func() func(B) BBS, GIOBS, func() func(B) BBS, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() func(B) BBS, GIOB, GIOBS, func() func(B) BBS, E, B, BBS],

		ma, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() B, GIOBS ~func() BBS, AAS ~[]A, BBS ~[]B, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, GIOBS, E, BBS],
		Map[GBS, func(E) func() func(B) BBS, GIOBS, func() func(B) BBS, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() func(B) BBS, GIOB, GIOBS, func() func(B) BBS, E, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(E) GIOA, GAS ~func(E) GIOAS, GIOA ~func() A, GIOAS ~func() AAS, AAS ~[]A, GAAS ~[]GA, E, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}
