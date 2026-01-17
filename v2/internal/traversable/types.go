package traversable

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	Pointed[A, HKT_A any]                       = pointed.Pointed[A, HKT_A]
	Functor[A, B, HKT_A, HKT_B any]             = functor.Functor[A, B, HKT_A, HKT_B]
	Applicative[A, B, HKT_A, HKT_B, HKT_AB any] = applicative.Applicative[A, B, HKT_A, HKT_B, HKT_AB]

	TraverseType[A, B, HKT_T_A, HKT_T_B, HKT_F_B, HKT_F_T_B, HKT_F_T_A_B any] = func(
		// applicative F
		f_of pointed.OfType[HKT_T_B, HKT_F_T_B],
		f_map functor.MapType[HKT_T_B, func(B) HKT_T_B, HKT_F_T_B, HKT_F_T_A_B],
		f_ap apply.ApType[HKT_F_B, HKT_F_T_B, HKT_F_T_A_B],

	) func(func(A) HKT_F_B) func(HKT_T_A) HKT_F_T_B

	SequenceType[
		HKT_T_F_A,
		HKT_F_T_A any] = func(
		// applicative F
		f_of pointed.OfType[HKT_T_F_A, HKT_F_T_A],
		f_map functor.MapType[HKT_T_F_A, func(HKT_T_F_A) HKT_T_F_A, HKT_F_T_A, HKT_T_F_A],
		f_ap apply.ApType[HKT_T_F_A, HKT_F_T_A, HKT_T_F_A],
	) func(HKT_T_F_A) HKT_F_T_A
)

func ComposeTraverse[
	A,
	B,
	HKT_F_B,
	HKT_G_A,
	HKT_G_B,
	HKT_T_G_A,
	HKT_T_G_B,
	HKT_F_G_B,
	HKT_F_T_G_B,
	HKT_F_T_A_B any](
	t TraverseType[HKT_G_A, HKT_G_B, HKT_T_G_A, HKT_T_G_B, HKT_F_G_B, HKT_F_T_G_B, HKT_F_T_A_B],
	g TraverseType[A, B, HKT_G_A, HKT_G_B, HKT_F_B, HKT_F_G_B, HKT_F_T_A_B],
) func(
	// applicative F
	f_of pointed.OfType[HKT_T_G_B, HKT_F_T_G_B],
	f_map functor.MapType[HKT_T_G_B, func(HKT_G_B) HKT_T_G_B, HKT_F_T_G_B, HKT_F_T_A_B],
	f_ap apply.ApType[HKT_F_G_B, HKT_F_T_G_B, HKT_F_T_A_B],

	// applicative G
	g_of pointed.OfType[HKT_G_B, HKT_F_G_B],
	g_map functor.MapType[HKT_G_B, func(B) HKT_G_B, HKT_F_G_B, HKT_F_T_A_B],
	g_ap apply.ApType[HKT_F_B, HKT_F_G_B, HKT_F_T_A_B],
) func(func(A) HKT_F_B) func(HKT_T_G_A) HKT_F_T_G_B {

	return func(
		// applicative F
		f_of pointed.OfType[HKT_T_G_B, HKT_F_T_G_B],
		f_map functor.MapType[HKT_T_G_B, func(HKT_G_B) HKT_T_G_B, HKT_F_T_G_B, HKT_F_T_A_B],
		f_ap apply.ApType[HKT_F_G_B, HKT_F_T_G_B, HKT_F_T_A_B],

		// applicative G
		g_of pointed.OfType[HKT_G_B, HKT_F_G_B],
		g_map functor.MapType[HKT_G_B, func(B) HKT_G_B, HKT_F_G_B, HKT_F_T_A_B],
		g_ap apply.ApType[HKT_F_B, HKT_F_G_B, HKT_F_T_A_B],

	) func(func(A) HKT_F_B) func(HKT_T_G_A) HKT_F_T_G_B {

		return F.Flow2(
			g(g_of, g_map, g_ap),
			t(f_of, f_map, f_ap),
		)
	}

}

// func ComposeSequence[
// 	HKT_F_G_A,
// 	HKT_G_F_A,
// 	HKT_T_G_F_A,
// 	HKT_F_T_G_A,
// 	HKT_T_F_G_A any](

// 	t SequenceType[HKT_T_F_G_A, HKT_F_T_G_A],
// 	t_map functor.MapType[HKT_G_F_A, HKT_F_G_A, HKT_T_G_F_A, HKT_T_F_G_A],

// 	g SequenceType[HKT_G_F_A, HKT_F_G_A],
// ) func(
// // applicative F
// ) func(HKT_T_G_F_A) HKT_F_T_G_A {

// 	return func() func(HKT_T_G_F_A) HKT_F_T_G_A {
// 		return F.Flow2(
// 			t_map(g()),
// 			t(),
// 		)
// 	}
// }

func SequenceFromTraverse[
	A, HKT_T_A, HKT_F_B, HKT_F_T_B any](
	t TraverseType[HKT_T_A, HKT_T_A, HKT_T_A, HKT_T_A, HKT_T_A, HKT_F_T_B, HKT_T_A],
) SequenceType[HKT_T_A, HKT_F_T_B] {

	return func(
		// applicative F
		f_of pointed.OfType[HKT_T_A, HKT_F_T_B],
		f_map functor.MapType[HKT_T_A, func(HKT_T_A) HKT_T_A, HKT_F_T_B, HKT_T_A],
		f_ap apply.ApType[HKT_T_A, HKT_F_T_B, HKT_T_A],
	) func(HKT_T_A) HKT_F_T_B {
		return F.Pipe1(
			F.Identity[HKT_T_A],
			t(f_of, f_map, f_ap),
		)
	}
}
