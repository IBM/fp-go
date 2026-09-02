//go:build go1.27

package common

// Compose returns a new lens that focuses deeper into the structure by
// chaining this lens (S → A) with an inner lens (A → B), producing a
// composed lens (S → B).
//
// This is the method-receiver form of the package-level [IsoComposeIso] function,
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
	return LensComposeLens[S](ab)(l)
}

// ComposeIso returns a new lens that focuses on a value of type B by composing this lens
// (S → A) with an isomorphism (A ↔ B), producing a lens (S → B).
//
// This is the method-receiver form of the package-level LensComposeIso function,
// available only on Go 1.27 and later because Go did not support type parameters on
// methods before that version. On earlier toolchains use the equivalent free function:
//
//	// method form (go1.27+)
//	fahrenheitLens := readingLens.ComposeIso(celsiusToFahrenheit)
//
//	// free-function form (all versions)
//	fahrenheitLens := LensComposeIso[Thermometer](celsiusToFahrenheit)(readingLens)
//
// The resulting lens satisfies all three lens laws whenever both the outer lens and the
// isomorphism individually satisfy them.
//
// Type Parameters:
//   - B: the focus type of the resulting lens — the target type of the isomorphism
//
// Parameters:
//   - ab: the isomorphism from A to B
//
// Returns:
//   - Lens[S, B]: a new lens from S directly to B
//
// See Also:
//   - Compose: the equivalent method for composing with a Lens[A, B] instead of an Iso[A, B]
//   - LensComposeIso: the equivalent package-level function
func (l Lens[S, A]) ComposeIso[B any](ab Iso[A, B]) Lens[S, B] {
	return LensComposeIso[S](ab)(l)
}

// ComposePrism returns a new Optional that focuses on a value of type B by
// composing this lens (S → A) with a prism (A → B), producing an Optional
// (S → B).
//
// This is the method-receiver form of the package-level LensComposePrism
// function, available only on Go 1.27 and later because Go did not support
// type parameters on methods before that version. On earlier toolchains use
// the equivalent free function:
//
//	// method form (go1.27+)
//	opt := myLens.ComposePrism(myPrism)
//
//	// free-function form (all versions)
//	opt := LensComposePrism[S](myPrism)(myLens)
//
// The result is an Optional because the prism may not match the focused value.
// GetOption applies the lens first, then the prism; Set updates the structure
// only when GetOption would return Some — matching the Optional no-op law.
//
// The composed Optional satisfies the standard Optional laws (GetSet, SetGet,
// SetSet) whenever the constituent lens and prism individually satisfy the lens
// and prism laws respectively.
//
// Type Parameters:
//   - B: the focus type of the prism and of the resulting Optional
//
// Parameters:
//   - ab: the prism from A to B
//
// Returns:
//   - Optional[S, B]: a new Optional from S directly to B
//
// See Also:
//   - Compose: method for composing with a Lens[A, B] instead of a Prism[A, B]
//   - ComposeIso: method for composing with an Iso[A, B] instead of a Prism[A, B]
//   - LensComposePrism: the equivalent package-level function
func (l Lens[S, A]) ComposePrism[B any](ab Prism[A, B]) Optional[S, B] {
	return LensComposePrism[S](ab)(l)
}
