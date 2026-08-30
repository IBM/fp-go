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
