// Copyright (c) 2024 - 2025 IBM Corp.
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

package fromio

// FromIO represents a type that can be constructed from an IO computation.
//
// This interface provides a way to lift IO computations into other monadic contexts,
// enabling interoperability between IO and other effect types.
//
// Type Parameters:
//   - A: The value type produced by the IO computation
//   - GA: The IO type (constrained to func() A)
//   - HKTA: The target higher-kinded type
type FromIO[A, GA ~func() A, HKTA any] interface {
	// FromIO converts an IO computation into the target monadic type.
	FromIO(GA) HKTA
}
