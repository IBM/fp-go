package iter

import (
	I "iter"

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	// Seq represents Go's standard library iterator type for single values.
	// It's an alias for iter.Seq[A] and provides interoperability with Go 1.23+ range-over-func.
	Seq[A any] = I.Seq[A]

	Endomorphism[A any] = endomorphism.Endomorphism[A]

	OfType[A, HKT_A any]             = pointed.OfType[A, HKT_A]
	MapType[A, B, HKT_A, HKT_B any]  = functor.MapType[A, B, HKT_A, HKT_B]
	ApType[HKT_A, HKT_B, HKT_AB any] = apply.ApType[HKT_A, HKT_B, HKT_AB]

	Kleisli[A, HKT_B any] = func(A) HKT_B
)
