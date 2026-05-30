package traversable

import F "github.com/IBM/fp-go/v2/function"

// ComposeTraversables combines two [Traversable] instances into a single [Traversable] for their
// composition T[G[_]], where T is the outer traversable and G is the inner traversable.
//
// This implements the standard result that the composition of two traversable functors is
// itself traversable: given a way to traverse T and a way to traverse G, we can traverse T[G[_]]
// by first using G's traverse to lift the mapping function inside the inner container, then
// using T's traverse to sequence those effects across the outer container.
//
// Type parameters (using Array as T and Option as G for concreteness):
//
//   - A          — the item type inside the inner container, e.g. string
//   - HKT_F_B   — the applicative effect applied to B, e.g. IO[int]
//   - HKT_T_A   — the inner container G applied to A, e.g. Option[string]
//   - HKT_F_T_B — the applicative effect applied to G[B], e.g. IO[Option[int]]
//   - HKT_T_G_A — the outer container T applied to G[A], e.g. []Option[string]
//   - HKT_F_T_G_B — the applicative effect applied to T[G[B]], e.g. IO[[]Option[int]]
//
// Parameters:
//
//   - outer — traverse for T with G[A] as element type:
//     (Option[A] → IO[Option[B]]) → []Option[A] → IO[[]Option[B]]
//   - inner — traverse for G over items of type A:
//     (A → IO[B]) → Option[A] → IO[Option[B]]
//
// Returns a [Traversable] for the composed structure T[G[_]]:
//
//	(A → IO[B]) → []Option[A] → IO[[]Option[B]]
//
// Example (Array outer, Option inner, IO effect):
//
//	traverse := ComposeTraversables(
//	    array.MakeTraversable[Option[string], Option[int], IO[Option[int]], IO[[]Option[int]]](),
//	    option.MakeTraversable[string, int, IO[int], IO[Option[int]]](),
//	)
//	// traverse : func(func(string) IO[int]) func([]Option[string]) IO[[]Option[int]]
func ComposeTraversables[A, HKT_F_B, HKT_T_A, HKT_F_T_B, HKT_T_G_A, HKT_F_T_G_B any](
	outer Traversable[HKT_T_A, HKT_F_T_B, HKT_T_G_A, HKT_F_T_G_B],
	inner Traversable[A, HKT_F_B, HKT_T_A, HKT_F_T_B],
) Traversable[A, HKT_F_B, HKT_T_G_A, HKT_F_T_G_B] {
	return F.Flow2(inner, outer)
}
