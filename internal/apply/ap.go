package apply

import (
	F "github.com/IBM/fp-go/function"
)

// HKTFGA  = HKT<F, HKT<G, A>>
// HKTFGB  = HKT<F, HKT<G, B>>
// HKTFGAB = HKT<F, HKT<G, (a: A) => B>>

// HKTGA   = HKT<G, A>
// HKTGB   = HKT<G, B>
// HKTGAB  = HKT<G, (a: A) => B>
func MonadAp[HKTGA, HKTGB, HKTGAB, HKTFGAB, HKTFGGAB, HKTFGA, HKTFGB any](
	fap func(HKTFGGAB, HKTFGA) HKTFGB,
	fmap func(HKTFGAB, func(HKTGAB) func(HKTGA) HKTGB) HKTFGGAB,
	gap func(HKTGAB, HKTGA) HKTGB,

	fab HKTFGAB,
	fa HKTFGA) HKTFGB {

	return fap(fmap(fab, F.Bind1st(F.Bind1st[HKTGAB, HKTGA, HKTGB], gap)), fa)
}

// export function ap<F, G>(
// 	F: Apply<F>,
// 	G: Apply<G>
//   ): <A>(fa: HKT<F, HKT<G, A>>) => <B>(fab: HKT<F, HKT<G, (a: A) => B>>) => HKT<F, HKT<G, B>> {
// 	return <A>(fa: HKT<F, HKT<G, A>>) => <B>(fab: HKT<F, HKT<G, (a: A) => B>>): HKT<F, HKT<G, B>> =>
// 	  F.ap(
// 		F.map(fab, (gab) => (ga: HKT<G, A>) => G.ap(gab, ga)),
// 		fa
// 	  )
//   }

//  function apFirst<F>(A: Apply<F>): <B>(second: HKT<F, B>) => <A>(first: HKT<F, A>) => HKT<F, A> {
// 	return (second) => (first) =>
// 	  A.ap(
// 		A.map(first, (a) => () => a),
// 		second
// 	  )
//   }

// Functor<F>.map: <A, () => A>(fa: HKT<F, A>, f: (a: A) => () => A) => HKT<F, () => A>

// Apply<F>.ap: <B, A>(fab: HKT<F, (a: B) => A>, fa: HKT<F, B>) => HKT<F, A>

func MonadApFirst[HKTGA, HKTGB, HKTGBA, A, B any](
	fap func(HKTGBA, HKTGB) HKTGA,
	fmap func(HKTGA, func(A) func(B) A) HKTGBA,

	first HKTGA,
	second HKTGB,
) HKTGA {
	return fap(
		fmap(first, F.Constant1[B, A]),
		second,
	)
}

func ApFirst[HKTGA, HKTGB, HKTGBA, A, B any](
	fap func(HKTGBA, HKTGB) HKTGA,
	fmap func(HKTGA, func(A) func(B) A) HKTGBA,

	second HKTGB,
) func(HKTGA) HKTGA {
	return func(first HKTGA) HKTGA {
		return MonadApFirst(fap, fmap, first, second)
	}
}

func MonadApSecond[HKTGA, HKTGB, HKTGBB, A, B any](
	fap func(HKTGBB, HKTGB) HKTGB,
	fmap func(HKTGA, func(A) func(B) B) HKTGBB,

	first HKTGA,
	second HKTGB,
) HKTGB {
	return fap(
		fmap(first, F.Constant1[A](F.Identity[B])),
		second,
	)
}

func ApSecond[HKTGA, HKTGB, HKTGBB, A, B any](
	fap func(HKTGBB, HKTGB) HKTGB,
	fmap func(HKTGA, func(A) func(B) B) HKTGBB,

	second HKTGB,
) func(HKTGA) HKTGB {
	return func(first HKTGA) HKTGB {
		return MonadApSecond(fap, fmap, first, second)
	}
}
