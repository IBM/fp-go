package generic

import (
	"context"

	E "github.com/IBM/fp-go/either"
	RE "github.com/IBM/fp-go/readerioeither/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[
	GRT ~func(context.Context) GIOT,
	GRA ~func(context.Context) GIOA,

	GIOA ~func() E.Either[error, A],
	GIOT ~func() E.Either[error, T.Tuple1[A]],

	A any](a GRA) GRT {
	return RE.SequenceT1[
		GRA,
		GRT,
	](a)
}

func SequenceT2[
	GRT ~func(context.Context) GIOT,
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOT ~func() E.Either[error, T.Tuple2[A, B]],

	A, B any](a GRA, b GRB) GRT {
	return RE.SequenceT2[
		GRA,
		GRB,
		GRT,
	](a, b)
}

func SequenceT3[
	GRT ~func(context.Context) GIOT,
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GRC ~func(context.Context) GIOC,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOC ~func() E.Either[error, C],
	GIOT ~func() E.Either[error, T.Tuple3[A, B, C]],

	A, B, C any](a GRA, b GRB, c GRC) GRT {
	return RE.SequenceT3[
		GRA,
		GRB,
		GRC,
		GRT,
	](a, b, c)
}

func SequenceT4[
	GRT ~func(context.Context) GIOT,
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GRC ~func(context.Context) GIOC,
	GRD ~func(context.Context) GIOD,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOC ~func() E.Either[error, C],
	GIOD ~func() E.Either[error, D],
	GIOT ~func() E.Either[error, T.Tuple4[A, B, C, D]],

	A, B, C, D any](a GRA, b GRB, c GRC, d GRD) GRT {
	return RE.SequenceT4[
		GRA,
		GRB,
		GRC,
		GRD,
		GRT,
	](a, b, c, d)
}
