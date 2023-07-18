package generic

import (
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
	RR "github.com/IBM/fp-go/internal/record"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func() ET.Either[E, B], GBS ~func() ET.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() ET.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GB, GBS, func() ET.Either[E, func(B) BBS], E, B, BBS],

		tas,
		f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func() ET.Either[E, B], GBS ~func() ET.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() ET.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GB, GBS, func() ET.Either[E, func(B) BBS], E, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func() ET.Either[E, A], GAS ~func() ET.Either[E, AAS], AAS ~[]A, GAAS ~[]GA, E, A any](tas GAAS) GAS {
	return MonadTraverseArray[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseRecord transforms an array
func MonadTraverseRecord[GB ~func() ET.Either[E, B], GBS ~func() ET.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() ET.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GB, GBS, func() ET.Either[E, func(B) BBS], E, B, BBS],

		tas,
		f,
	)
}

// TraverseRecord transforms an array
func TraverseRecord[GB ~func() ET.Either[E, B], GBS ~func() ET.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() ET.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GB, GBS, func() ET.Either[E, func(B) BBS], E, B, BBS],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[GA ~func() ET.Either[E, A], GAS ~func() ET.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, E, A any](tas GAAS) GAS {
	return MonadTraverseRecord[GA, GAS](tas, F.Identity[GA])
}
