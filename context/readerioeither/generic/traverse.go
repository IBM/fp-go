package generic

import (
	"context"

	E "github.com/IBM/fp-go/either"
	RE "github.com/IBM/fp-go/readerioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[
	AS ~[]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~[]B,
	A, B any](f func(A) GRB) func(AS) GRBS {
	return RE.TraverseArray[GRB, GRBS, GIOB, GIOBS, AS](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[
	AS ~[]A,
	GAS ~[]GRA,
	GRAS ~func(context.Context) GIOAS,
	GRA ~func(context.Context) GIOA,
	GIOAS ~func() E.Either[error, AS],
	GIOA ~func() E.Either[error, A],
	A any](ma GAS) GRAS {
	return RE.SequenceArray[GRA, GRAS](ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable,
	AS ~map[K]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~map[K]B,

	A, B any](f func(A) GRB) func(AS) GRBS {
	return RE.TraverseRecord[GRB, GRBS, GIOB, GIOBS, AS](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable,
	AS ~map[K]A,
	GAS ~map[K]GRA,
	GRAS ~func(context.Context) GIOAS,
	GRA ~func(context.Context) GIOA,
	GIOAS ~func() E.Either[error, AS],
	GIOA ~func() E.Either[error, A],
	A any](ma GAS) GRAS {
	return RE.SequenceRecord[GRA, GRAS](ma)
}
