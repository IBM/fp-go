package generic

import (
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
)

// TraverseArray transforms an array
func MonadTraverseArray[GB ~func(R) B, GBS ~func(R) BBS, AAS ~[]A, BBS ~[]B, R, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, R, BBS],
		Map[GBS, func(R) func(B) BBS, R, BBS, func(B) BBS],
		Ap[GB, GBS, func(R) func(B) BBS, R, B, BBS],
		tas, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(R) B, GBS ~func(R) BBS, AAS ~[]A, BBS ~[]B, R, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, R, BBS],
		Map[GBS, func(R) func(B) BBS, R, BBS, func(B) BBS],
		Ap[GB, GBS, func(R) func(B) BBS, R, B, BBS],
		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(R) A, GAS ~func(R) AAS, AAS ~[]A, GAAS ~[]GA, R, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}
