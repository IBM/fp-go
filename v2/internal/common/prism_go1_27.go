//go:build go1.27

package common

// Compose returns a new prism that focuses deeper into a sum type by chaining
// this prism (S → A) with an inner prism (A → B), producing a composed prism
// (S → B).
//
// This is the method-receiver form of the package-level [OptionalComposeOptional] function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := outerPrism.Compose(innerPrism)
//
//	// free-function form (all versions)
//	composed := prism.Compose[S](innerPrism)(outerPrism)
//
// GetOption of the composed prism chains the two GetOption functions: it first
// applies the outer prism to obtain an Option[A], then chains the inner prism
// through that option via O.Chain, returning None whenever either prism fails
// to match.
//
// ReverseGet of the composed prism pipes the value back through both
// ReverseGet functions in order: inner first, then outer.
//
// The composed prism satisfies both prism laws whenever both constituent
// prisms individually satisfy them.
//
// Type Parameters:
//   - B: the focus type of the inner prism and of the resulting prism
//
// Parameters:
//   - ab: the inner prism from A to B
//
// Returns:
//   - Prism[S, B]: a new prism from S directly to B
//
// See Also:
//   - Compose: the equivalent package-level function
func (p Prism[S, A]) Compose[B any](ab Prism[A, B]) Prism[S, B] {
	return PrismComposePrism[S](ab)(p)
}

// ComposeLens returns a new Optional that focuses on a value of type B by
// composing this prism (S → A) with a lens (A → B), producing an Optional
// (S → B).
//
// Because a Lens always succeeds, the resulting Optional is None only when this
// prism does not match the source value. Set is a no-op in that same case, and
// the original S is returned unchanged — consistent with the Optional no-op law.
//
// This is the method-receiver form of the package-level PrismComposeLens
// function, available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	opt := myPrism.ComposeLens(myLens)
//
//	// free-function form (all versions)
//	opt := PrismComposeLens[S](myLens)(myPrism)
//
// GetOption of the composed optional applies this prism's GetOption to obtain
// an Option[A]; when Some, it applies the lens's Get to retrieve B, returning
// None whenever this prism does not match.
//
// Set of the composed optional checks whether this prism matches before writing:
// when it does, it uses the lens's Set to embed the new B back into A, then
// reconstructs S via the prism's ReverseGet. When this prism does not match,
// Set is a no-op and the original S is returned unchanged.
//
// The composed Optional satisfies all three optional laws (GetSet, SetGet,
// SetSet) whenever this prism satisfies the prism laws (the lens laws are
// always satisfied by a well-formed Lens).
//
// Type Parameters:
//   - B: the focus type of the inner lens and of the resulting Optional
//
// Parameters:
//   - ab: the inner lens from A to B
//
// Returns:
//   - Optional[S, B]: a new Optional from S directly to B
//
// See Also:
//   - Compose: method for composing with a Prism[A, B] instead of a Lens
//   - PrismComposeLens: the equivalent package-level function
func (p Prism[S, A]) ComposeLens[B any](ab Lens[A, B]) Optional[S, B] {
	return PrismComposeLens[S](ab)(p)
}

// ComposeIso returns a new prism by composing this prism (S → A) with an
// isomorphism (A ↔ B), producing a Prism[S, B].
//
// Because an Iso is total in both directions, the resulting Prism is None only
// when this prism does not match the source value — the iso never introduces
// additional partiality.
//
// This is the method-receiver form of the package-level PrismComposeIso function,
// available only on Go 1.27 and later because Go did not support type parameters
// on methods before that version. On earlier toolchains, use the equivalent free
// function instead:
//
//	// method form (go1.27+)
//	composed := myPrism.ComposeIso(myIso)
//
//	// free-function form (all versions)
//	composed := PrismComposeIso[S](myIso)(myPrism)
//
// GetOption of the composed prism applies this prism's GetOption to obtain
// Option[A], then maps the iso's forward function over the option to yield
// Option[B]; None is returned whenever this prism does not match.
//
// ReverseGet of the composed prism converts B → A via the iso's ReverseGet,
// then converts A → S via this prism's ReverseGet.
//
// The composed prism satisfies both prism laws whenever this prism and the
// isomorphism individually satisfy their respective laws.
//
// Type Parameters:
//   - B: the target type of the isomorphism and the focus type of the resulting prism
//
// Parameters:
//   - ab: the isomorphism from A to B
//
// Returns:
//   - Prism[S, B]: a new prism from S directly to B
//
// See Also:
//   - Compose: method for composing with a Prism[A, B] instead of an Iso
//   - ComposeLens: method for composing with a Lens[A, B] (yields an Optional, not a Prism)
//   - PrismComposeIso: the equivalent package-level function
func (p Prism[S, A]) ComposeIso[B any](ab Iso[A, B]) Prism[S, B] {
	return PrismComposeIso[S](ab)(p)
}

// ComposeOptional returns a new Optional by composing this prism (S → A) with
// an Optional[A, B], producing an Optional[S, B].
//
// Both steps are partial: the Prism may not match the source variant, and the
// inner Optional may not match the intermediate value A.  GetOption returns None
// whenever either step fails to match; Set is a no-op whenever GetOption would
// return None (satisfying the GetSet law).
//
// This is the method-receiver form of the package-level PrismComposeOptional
// function, available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	opt := myPrism.ComposeOptional(myOptional)
//
//	// free-function form (all versions)
//	opt := PrismComposeOptional[S](myOptional)(myPrism)
//
// GetOption of the resulting Optional first applies this prism's GetOption to
// obtain Option[A]; when Some(a), it then applies the inner optional's GetOption
// to a to obtain Option[B].  None is returned if either step yields None.
//
// Set of the resulting Optional applies both optics' Set paths in sequence: the
// inner optional's Set embeds B back into A, and this prism's ReverseGet then
// reconstructs S from the updated A.  If the prism does not match, Set is a
// no-op and the original S is returned unchanged.
//
// The resulting Optional satisfies all three standard Optional laws whenever both
// constituent optics individually satisfy their laws:
//
//  1. GetSet: GetOption(s) = None => Set(b)(s) = s
//  2. SetGet: GetOption(s) = Some(_) => GetOption(Set(b)(s)) = Some(b)
//  3. SetSet: Set(c)(Set(b)(s)) = Set(c)(s)
//
// Type Parameters:
//   - B: the focus type of the inner Optional and of the resulting Optional
//
// Parameters:
//   - ab: the inner Optional[A, B]
//
// Returns:
//   - Optional[S, B]: a new Optional from S directly to B
//
// See Also:
//   - ComposeLens: method for composing with a Lens[A, B] (always succeeds for the inner step)
//   - Compose: method for composing with a Prism[A, B] (yields a Prism, not an Optional)
//   - PrismComposeOptional: the equivalent package-level function
func (p Prism[S, A]) ComposeOptional[B any](ab Optional[A, B]) Optional[S, B] {
	return PrismComposeOptional[S](ab)(p)
}
