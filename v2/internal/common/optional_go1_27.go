//go:build go1.27

package common

// Compose returns a new Optional that focuses deeper into a structure by
// chaining this optional (S → A) with an inner optional (A → B), producing a
// composed optional (S → B).
//
// This is the method-receiver form of the package-level [OptionalComposeOptional] function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := outerOptional.Compose(innerOptional)
//
//	// free-function form (all versions)
//	composed := optional.Compose[S](innerOptional)(outerOptional)
//
// GetOption of the composed optional chains the two GetOption functions: it
// first applies the outer optional to obtain an Option[A], then chains the
// inner optional through that option via O.Chain, returning None whenever
// either optional produces None.
//
// Set of the composed optional updates the deeply nested focus by threading
// the new value back through the inner optional's Set and then propagating
// the change upward via the outer optional's Set. When either optional
// produces None, Set is a no-op and the original S is returned unchanged,
// consistent with the Optional no-op law.
//
// The composed optional satisfies all three optional laws (GetSet, SetGet,
// SetSet) whenever both constituent optionals individually satisfy them.
//
// Type Parameters:
//   - B: the focus type of the inner optional and of the resulting optional
//
// Parameters:
//   - ab: the inner optional from A to B
//
// Returns:
//   - Optional[S, B]: a new optional from S directly to B
//
// See Also:
//   - Compose: the equivalent package-level function
//   - ComposeRef: variant for pointer-based outer structures
func (p Optional[S, A]) Compose[B any](ab Optional[A, B]) Optional[S, B] {
	return OptionalComposeOptional[S](ab)(p)
}

// ComposeLens returns a new Optional that focuses deeper into a structure by
// composing this optional (S → A) with a lens (A → B), producing an optional
// (S → B).
//
// Because a Lens always succeeds, the resulting Optional inherits its partiality
// exclusively from the outer optional: GetOption returns None only when this
// optional produces None; Set is a no-op only in that same case.
//
// This is the method-receiver form of the package-level OptionalComposeLens
// function, available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := outerOptional.ComposeLens(innerLens)
//
//	// free-function form (all versions)
//	composed := OptionalComposeLens[S](innerLens)(outerOptional)
//
// GetOption of the composed optional applies this optional's GetOption to obtain
// an Option[A], then – when Some – applies the lens's Get to retrieve B,
// returning None whenever this optional produces None.
//
// Set of the composed optional updates the deeply nested focus: it uses the
// lens's Set to embed the new B back into A, then propagates the updated A
// upward via this optional's Set. When this optional produces None, Set is a
// no-op and the original S is returned unchanged, consistent with the Optional
// no-op law.
//
// The composed optional satisfies all three optional laws (GetSet, SetGet,
// SetSet) whenever this optional individually satisfies them (the lens laws
// are always satisfied by a well-formed Lens).
//
// Type Parameters:
//   - B: the focus type of the inner lens and of the resulting optional
//
// Parameters:
//   - ab: the inner lens from A to B
//
// Returns:
//   - Optional[S, B]: a new optional from S directly to B
//
// See Also:
//   - Compose: variant that composes with an Optional[A, B] instead of a Lens
//   - OptionalComposeLens: the equivalent package-level function
func (l Optional[S, A]) ComposeLens[B any](ab Lens[A, B]) Optional[S, B] {
	return OptionalComposeLens[S](ab)(l)
}

// ComposePrism returns a new Optional that focuses deeper into a structure by
// composing this optional (S → A) with a prism (A → B), producing an optional
// (S → B).
//
// Because a Prism may not match its input, the resulting Optional is None
// whenever either this optional produces None or the prism does not match the
// focused A value.  Set is a no-op in both cases.
//
// This is the method-receiver form of the package-level OptionalComposePrism
// function, available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := outerOptional.ComposePrism(innerPrism)
//
//	// free-function form (all versions)
//	composed := OptionalComposePrism[S](innerPrism)(outerOptional)
//
// GetOption of the composed optional first applies this optional's GetOption to
// obtain an Option[A]; when Some, it applies the prism's GetOption to try to
// extract B from A.  None is returned whenever either step fails.
//
// Set of the composed optional checks whether GetOption would return Some before
// writing: when it would, it uses the prism's ReverseGet to reconstruct an A
// from the new B, then propagates the updated A upward via this optional's Set.
// When GetOption returns None (either because this optional misses or because
// the prism does not match), Set is a no-op and the original S is returned
// unchanged, consistent with the Optional no-op law.
//
// The composed optional satisfies all three optional laws (GetSet, SetGet,
// SetSet) whenever both this optional and the prism individually satisfy their
// respective laws.
//
// Type Parameters:
//   - B: the focus type of the prism and of the resulting optional
//
// Parameters:
//   - ab: the prism from A to B
//
// Returns:
//   - Optional[S, B]: a new optional from S directly to B
//
// See Also:
//   - Compose: variant that composes with an Optional[A, B] instead of a Prism
//   - ComposeLens: variant that composes with a Lens[A, B] instead of a Prism
//   - OptionalComposePrism: the equivalent package-level function
func (l Optional[S, A]) ComposePrism[B any](ab Prism[A, B]) Optional[S, B] {
	return OptionalComposePrism[S](ab)(l)
}
