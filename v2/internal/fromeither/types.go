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

package fromeither

import (
	ET "github.com/IBM/fp-go/v2/either"
)

// FromEither represents a type that can be constructed from an Either value.
//
// This interface provides a way to lift Either values into other monadic contexts,
// enabling interoperability between Either and other effect types.
//
// Type Parameters:
//   - E: The error type in the Either
//   - A: The success value type in the Either
//   - HKTA: The target higher-kinded type
type FromEither[E, A, HKTA any] interface {
	// FromEither converts an Either value into the target monadic type.
	FromEither(ET.Either[E, A]) HKTA
}
