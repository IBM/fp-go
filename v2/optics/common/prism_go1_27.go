//go:build go1.27

package common

// Compose returns a new prism that focuses deeper into a sum type by chaining
// this prism (S → A) with an inner prism (A → B), producing a composed prism
// (S → B).
//
// This is the method-receiver form of the package-level [Compose] function,
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
