//go:build go1.27

package iso

// Compose returns a new isomorphism that focuses deeper by chaining this
// isomorphism (S → A) with an inner isomorphism (A → B), producing a composed
// isomorphism (S → B).
//
// This is the method-receiver form of the package-level [Compose] function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := outerIso.Compose(innerIso)
//
//	// free-function form (all versions)
//	composed := iso.Compose[S](innerIso)(outerIso)
//
// Get of the composed isomorphism applies the outer Get followed by the inner
// Get: composed.Get(s) = ab.Get(sa.Get(s)).
//
// ReverseGet of the composed isomorphism applies the inner ReverseGet followed
// by the outer ReverseGet: composed.ReverseGet(b) = sa.ReverseGet(ab.ReverseGet(b)).
//
// The composed isomorphism satisfies both iso round-trip laws whenever both
// constituent isomorphisms individually satisfy them:
//
//	composed.ReverseGet(composed.Get(s)) == s    for all s: S
//	composed.Get(composed.ReverseGet(b)) == b    for all b: B
//
// Type Parameters:
//   - B: the target type of the inner isomorphism and of the resulting isomorphism
//
// Parameters:
//   - ab: the inner isomorphism from A to B
//
// Returns:
//   - Iso[S, B]: a new isomorphism from S directly to B
//
// See Also:
//   - Compose: the equivalent package-level function
//   - Reverse: swaps the direction of an isomorphism
func (p Iso[S, A]) Compose[B any](ab Iso[A, B]) Iso[S, B] {
	return Compose[S](ab)(p)
}
