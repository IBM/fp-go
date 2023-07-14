package either

import (
	F "github.com/ibm/fp-go/function"
	RR "github.com/ibm/fp-go/internal/record"
)

// TraverseRecord transforms a record of options into an option of a record
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, E, A, B any](f func(A) Either[E, B]) func(GA) Either[E, GB] {
	return RR.Traverse[GA](
		Of[E, GB],
		MonadMap[E, GB, func(B) GB],
		MonadAp[GB, E, B],
		f,
	)
}

// TraverseRecord transforms a record of eithers into an either of a record
func TraverseRecord[K comparable, E, A, B any](f func(A) Either[E, B]) func(map[K]A) Either[E, map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Either[E, A], K comparable, E, A any](ma GOA) Either[E, GA] {
	return TraverseRecordG[GOA, GA](F.Identity[Either[E, A]])(ma)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]Either[E, A]) Either[E, map[K]A] {
	return SequenceRecordG[map[K]A](ma)
}
