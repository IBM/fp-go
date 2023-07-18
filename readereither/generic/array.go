package generic

import (
	ET "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	RA "github.com/ibm/fp-go/internal/array"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func(E) ET.Either[L, B], GBS ~func(E) ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, L, E, A, B any](ma AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, L, E, BBS],
		Map[GBS, func(E) ET.Either[L, func(B) BBS], L, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) ET.Either[L, func(B) BBS], L, E, B, BBS],

		ma, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(E) ET.Either[L, B], GBS ~func(E) ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, L, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, L, E, BBS],
		Map[GBS, func(E) ET.Either[L, func(B) BBS], L, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) ET.Either[L, func(B) BBS], L, E, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(E) ET.Either[L, A], GAS ~func(E) ET.Either[L, AAS], AAS ~[]A, GAAS ~[]GA, L, E, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}
