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
