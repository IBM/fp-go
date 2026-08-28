//go:build go1.27

package lens

// Compose returns a new lens that focuses deeper into the structure by
// chaining this lens (S → A) with an inner lens (A → B), producing a
// composed lens (S → B).
//
// This is the method-receiver form of the package-level [Compose] function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	personStreetLens := addressLens.Compose(streetLens)
//
//	// free-function form (all versions)
//	personStreetLens := lens.Compose[Person](streetLens)(addressLens)
//
// The composed lens satisfies all three lens laws whenever both constituent
// lenses individually satisfy them.
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
//   - Compose: the equivalent package-level function
func (l Lens[S, A]) Compose[B any](ab Lens[A, B]) Lens[S, B] {
	return Compose[S](ab)(l)
}
