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

package fromoption

import (
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Option represents an optional value that may or may not be present.
	Option[T any] = option.Option[T]

	// FromOption represents a type that can be constructed from an Option value.
	//
	// This interface provides a way to lift Option values into other monadic contexts,
	// enabling interoperability between Option and other effect types.
	//
	// Type Parameters:
	//   - A: The value type in the Option
	//   - HKTA: The target higher-kinded type
	FromOption[A, HKTA any] interface {
		// FromEither converts an Option value into the target monadic type.
		// Note: The method name should probably be FromOption, but is FromEither for compatibility.
		FromEither(Option[A]) HKTA
	}
)
