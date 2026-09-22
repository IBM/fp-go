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

package prism

import (
	"uuid"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

// ParseUUID creates a prism for parsing and formatting UUIDs.
// It provides a safe way to work with UUID strings, handling parsing
// errors gracefully through the Option type.
//
// The prism's GetOption attempts to parse a string into a uuid.UUID.
// If parsing succeeds, it returns Some(uuid.UUID); if it fails (e.g. the string
// is not a valid UUID), it returns None.
//
// The prism's ReverseGet always succeeds, converting a uuid.UUID back to its
// canonical lower case, hyphenated string representation.
//
// GetOption accepts every textual form understood by uuid.Parse, i.e. the plain
// hyphenated form, the brace-enclosed form, the urn:uuid: prefixed form and the
// unhyphenated form, with hexadecimal digits in either case.
//
// This function requires Go 1.27 or later because it builds on the uuid package
// of the standard library, which was added in that release.
//
// Prism laws:
//
// The ReverseGet-GetOption law holds unconditionally: for every uuid.UUID u,
// GetOption(ReverseGet(u)) equals Some(u).
//
// The GetOption-ReverseGet law only holds for canonical input: a string s round
// trips to itself exactly when it is lower case and hyphenated, because
// ReverseGet always produces that form. Braced, urn-prefixed, unhyphenated or
// upper case input is normalized rather than preserved.
//
// Returns:
//   - A Prism[string, uuid.UUID] that safely handles UUID parsing/formatting
//
// Common use cases:
//   - Validating and parsing UUID configuration or environment values
//   - Extracting resource identifiers from request paths or headers
//   - Normalizing UUID strings to their canonical representation
//   - Building a validating string-to-UUID codec via codec.FromPrism
//
// See Also:
//   - ParseURL: the analogous prism for URL strings
//   - ParseInt: the analogous prism for integer strings
//   - MakePrismWithName: the constructor used to build this prism
func ParseUUID() Prism[string, uuid.UUID] {
	return MakePrismWithName(
		F.Flow2(
			result.Eitherize1(uuid.Parse),
			result.ToOption,
		),
		uuid.UUID.String,
		"PrismParseUUID",
	)
}
