package witherable

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/filterable"
	"github.com/IBM/fp-go/v2/internal/functor"
)

func Filter[A, HKT_G_A, HKT_F_HKT_G_A any](
	fmap functor.MapType[HKT_G_A, HKT_G_A, HKT_F_HKT_G_A, HKT_F_HKT_G_A],
	ffilter filterable.FilterType[A, HKT_G_A],
) func(func(A) bool) func(HKT_F_HKT_G_A) HKT_F_HKT_G_A {
	return function.Flow2(
		ffilter,
		fmap,
	)
}

func FilterMap[A, B, HKT_G_A, HKT_G_B, HKT_F_HKT_G_A, HKT_F_HKT_G_B any](
	fmap functor.MapType[HKT_G_A, HKT_G_B, HKT_F_HKT_G_A, HKT_F_HKT_G_B],
	ffilter filterable.FilterMapType[A, B, HKT_G_A, HKT_G_B],
) func(func(A) Option[B]) func(HKT_F_HKT_G_A) HKT_F_HKT_G_B {
	return function.Flow2(
		ffilter,
		fmap,
	)
}
