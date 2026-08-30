//go:build go1.27

// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"github.com/IBM/fp-go/v2/optics/optional"
)

// Pipe returns a new Type that refines or transforms the decoded value by
// chaining this codec (I → A) with a second codec (A → B), producing a
// composed codec (I → B).
//
// This is the method-receiver form of the package-level Pipe function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version. On earlier toolchains, use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.Pipe(refinement)
//
//	// free-function form (all versions)
//	composed := codec.Pipe[O, I](refinement)(base)
//
// Decoding applies both codecs in sequence: the receiver decodes the raw input
// I to an intermediate value A, then the argument codec validates and converts
// that intermediate value to the final type B. If either step fails, the
// composed codec propagates the validation errors.
//
// Encoding applies both codecs in reverse: ab encodes B back to A, then the
// receiver encodes A back to O.  Encoding never fails.
//
// The name of the composed codec is "Pipe(<receiver-name>, <ab-name>)".
//
// Type Parameters:
//   - B: the final decoded type produced by the composed codec
//
// Parameters:
//   - ab: the second codec that accepts the receiver's decoded type A as input
//
// Returns:
//   - Type[B, O, I]: a new codec from I directly to B
//
// See Also:
//   - Pipe: the equivalent package-level function
func (t Type[A, O, I]) Pipe[B any](ab Type[B, A, A]) Type[B, O, I] {
	return Pipe[O, I](ab)(t)
}

// PipeIso composes the receiver codec (I → A) with an isomorphism Iso[A, B],
// producing a new codec (I → B).
//
// This is the method-receiver form of the package-level PipeIso function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.PipeIso(myIso)
//
//	// free-function form (all versions)
//	composed := codec.PipeIso[O, I](myIso)(base)
//
// Because an isomorphism is a lossless transformation, decoding always
// succeeds once the receiver codec has decoded the raw input I to A
// successfully.  No additional validation is performed by the isomorphism
// step.
//
// Type Parameters:
//   - B: the final decoded type produced by the composed codec
//
// Parameters:
//   - ab: the Iso[A, B] to apply after the receiver has decoded to A
//
// Returns:
//   - Type[B, O, I]: a new codec from I directly to B via the isomorphism
//
// See Also:
//   - PipeIso: the equivalent package-level function
//   - PipeRefinement: method form for partial (potentially failing) refinements
func (t Type[A, O, I]) PipeIso[B any](ab Iso[A, B]) Type[B, O, I] {
	return PipeIso[O, I](ab)(t)
}

// PipeRefinement composes the receiver codec (I → A) with a Refinement[A, B],
// producing a new codec (I → B).
//
// This is the method-receiver form of the package-level PipeRefinement function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.PipeRefinement(myRefinement)
//
//	// free-function form (all versions)
//	composed := codec.PipeRefinement[O, I](myRefinement)(base)
//
// Unlike PipeIso, decoding may fail: if the decoded A does not satisfy the
// refinement's predicate (GetOption returns None) a validation error is
// propagated and decoding fails.
//
// Type Parameters:
//   - B: the refined type produced when the predicate is satisfied
//
// Parameters:
//   - ab: the Refinement[A, B] that validates and narrows values of type A to B
//
// Returns:
//   - Type[B, O, I]: a new codec from I directly to B, failing when the
//     refinement predicate is not satisfied
//
// See Also:
//   - PipeRefinement: the equivalent package-level function
//   - PipeIso: method form for lossless isomorphism transformations
func (t Type[A, O, I]) PipeRefinement[B any](ab Refinement[A, B]) Type[B, O, I] {
	return PipeRefinement[O, I](ab)(t)
}

// PipePrism composes the receiver codec (I → A) with a Prism[A, B],
// producing a new codec (I → B).
//
// This is the method-receiver form of the package-level PipePrism function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.PipePrism(myPrism)
//
//	// free-function form (all versions)
//	composed := codec.PipePrism[O, I](myPrism)(base)
//
// PipePrism is an alias for PipeRefinement on the method receiver.  Unlike
// PipeIso, decoding may fail: if the decoded A does not match the prism's
// focus (GetOption returns None) a validation error is propagated.
//
// Type Parameters:
//   - B: the focus type produced when the prism matches
//
// Parameters:
//   - ab: the Prism[A, B] applied after the receiver has decoded to A
//
// Returns:
//   - Type[B, O, I]: a new codec from I directly to B, failing when the
//     prism does not match
//
// See Also:
//   - PipePrism: the equivalent package-level function
//   - PipeRefinement: method form using the Refinement alias
//   - PipeIso: method form for lossless isomorphism transformations
func (t Type[A, O, I]) PipePrism[B any](ab Prism[A, B]) Type[B, O, I] {
	return PipeRefinement[O, I](ab)(t)
}

// Alt adds a fallback codec to the receiver, returning a new codec that tries
// the receiver first and falls back to second when decoding fails.
//
// This is the method-receiver form of the package-level Alt function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.Alt(fallback)
//
//	// free-function form (all versions)
//	composed := codec.Alt(fallback)(base)
//
// Decoding tries the receiver; if it succeeds the result is returned and
// second is never evaluated.  If it fails, second is evaluated and its
// result is returned.  When both fail their errors are accumulated.
// Encoding always uses the receiver's encoder.
//
// The name of the composed codec is "Alt[<receiver-name>]".
//
// Parameters:
//   - second: a lazy Type[A, O, I] evaluated only when the receiver fails
//
// Returns:
//   - Type[A, O, I]: a new codec with fallback behaviour
//
// See Also:
//   - Alt: the equivalent package-level function
//   - AltW: the package-level widening alternative (cannot be a method due to
//     Go's generic instantiation constraints)
func (t Type[A, O, I]) Alt(second Lazy[Type[A, O, I]]) Type[A, O, I] {
	return MonadAlt(t, second)
}

// ApSL sequences a required struct field into the receiver codec using a lens.
//
// This is the method-receiver form of the package-level ApSL function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.ApSL(monoid, fieldLens, fieldCodec)
//
//	// free-function form (all versions)
//	composed := codec.ApSL(monoid, fieldLens, fieldCodec)(base)
//
// Decoding validates the field T from input I using fa and sets it into the
// decoded S via the lens; errors from both the base codec and the field codec
// are accumulated (applicative semantics).  Encoding extracts T from S via
// the lens, encodes it, and combines it with the base encoding using the
// monoid.
//
// Type Parameters:
//   - T: the field type focused by the lens
//
// Parameters:
//   - m: a Monoid[O] for combining encoded outputs
//   - l: a Lens[A, T] that focuses on the field within A
//   - fa: a Type[T, O, I] codec for the field
//
// Returns:
//   - Type[A, O, I]: the receiver codec extended with the field
//
// See Also:
//   - ApSL: the equivalent package-level function
//   - ApSO: method form for optional fields
//   - Bind: method form for context-dependent (Kleisli) fields
func (t Type[A, O, I]) ApSL[T any](m Monoid[O], l Lens[A, T], fa Type[T, O, I]) Type[A, O, I] {
	return ApSL(m, l, fa)(t)
}

// ApSO sequences an optional struct field into the receiver codec using an optional.
//
// This is the method-receiver form of the package-level ApSO function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.ApSO(monoid, fieldOpt, fieldCodec)
//
//	// free-function form (all versions)
//	composed := codec.ApSO(monoid, fieldOpt, fieldCodec)(base)
//
// Unlike ApSL, the field focus may be absent.  Decoding always sets the field
// when present in the input; encoding only emits the field when the optional's
// GetOption returns Some.  Errors from both the base codec and the field codec
// are accumulated (applicative semantics).
//
// Type Parameters:
//   - T: the optional field type focused by the optional
//
// Parameters:
//   - m: a Monoid[O] for combining encoded outputs
//   - o: an optional.Optional[A, T] that focuses on the field within A
//   - fa: a Type[T, O, I] codec for the field
//
// Returns:
//   - Type[A, O, I]: the receiver codec extended with the optional field
//
// See Also:
//   - ApSO: the equivalent package-level function
//   - ApSL: method form for required fields via Lens
//   - Bind: method form for context-dependent (Kleisli) fields
func (t Type[A, O, I]) ApSO[T any](m Monoid[O], o optional.Optional[A, T], fa Type[T, O, I]) Type[A, O, I] {
	return ApSO(m, o, fa)(t)
}

// Bind sequences a context-dependent field into the receiver codec using a lens
// and a Kleisli arrow.
//
// This is the method-receiver form of the package-level Bind function,
// available only on Go 1.27 and later because Go did not support type
// parameters on methods before that version.  On earlier toolchains use the
// equivalent free function instead:
//
//	// method form (go1.27+)
//	composed := base.Bind(monoid, fieldLens, kleisli)
//
//	// free-function form (all versions)
//	composed := codec.Bind(monoid, fieldLens, kleisli)(base)
//
// Unlike ApSL, the field codec is determined by the already-decoded struct
// value (monadic sequencing): if the base codec fails the Kleisli arrow is
// never evaluated (fail-fast).  Encoding evaluates the Kleisli arrow on the
// struct value, encodes the field, and combines the result using the monoid.
//
// Type Parameters:
//   - T: the field type focused by the lens
//
// Parameters:
//   - m: a Monoid[O] for combining encoded outputs
//   - l: a Lens[A, T] that focuses on the field within A
//   - f: a Kleisli[A, T, O, I] — a function from the decoded A to Type[T, O, I]
//
// Returns:
//   - Type[A, O, I]: the receiver codec extended with the context-dependent field
//
// See Also:
//   - Bind: the equivalent package-level function
//   - ApSL: method form for fixed (non-context-dependent) required fields
//   - ApSO: method form for optional fields
//
// Note: the parameter is spelled Reader[A, Type[T, O, I]] rather than the
// equivalent Kleisli[A, T, O, I]. The two are identical by definition, but a
// method on Type may not name a generic alias whose own body instantiates
// Type: importing such a package self-deadlocks the go1.27 type checker while
// it unpacks Type (cmd/compile/internal/types2.(*Named).unpack).
func (t Type[A, O, I]) Bind[T any](m Monoid[O], l Lens[A, T], f Reader[A, Type[T, O, I]]) Type[A, O, I] {
	return Bind(m, l, f)(t)
}
