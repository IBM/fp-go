//go:build go1.27

package common

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
	return IsoComposeIso[S](ab)(p)
}

// ComposeLens returns a new lens by composing this isomorphism (S ↔ A) with an
// inner lens (A → B), producing a Lens[S, B].
//
// Internally the isomorphism is first converted to a Lens[S, A] via IsoAsLens,
// then composed with the provided lens using LensComposeLens. The result focuses
// on a value of type B that is reached by first applying the isomorphism's Get
// direction and then the inner lens's Get.
//
// This is the method-receiver form of the package-level IsoComposeLens function,
// available only on Go 1.27 and later because Go did not support type parameters
// on methods before that version. On earlier toolchains use the equivalent free
// function:
//
//	// method form (go1.27+)
//	streetLens := addressIso.ComposeLens(streetFieldLens)
//
//	// free-function form (all versions)
//	streetLens := IsoComposeLens[S](streetFieldLens)(addressIso)
//
// The resulting lens satisfies all three lens laws whenever both the isomorphism
// and the inner lens individually satisfy them.
//
// Type Parameters:
//   - B: the focus type of the inner lens and of the resulting lens
//
// Parameters:
//   - ab: the inner lens from A to B
//
// Returns:
//   - Lens[S, B]: a new lens from S directly to B
//
// See Also:
//   - Compose: the equivalent method for composing with an Iso[A, B] instead of a Lens[A, B]
//   - IsoComposeLens: the equivalent package-level function
func (p Iso[S, A]) ComposeLens[B any](ab Lens[A, B]) Lens[S, B] {
	return IsoComposeLens[S](ab)(p)
}
