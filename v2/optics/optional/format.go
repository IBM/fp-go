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

package optional

import (
	"fmt"
	"log/slog"

	"github.com/IBM/fp-go/v2/internal/formatting"
)

// String returns the name of the optional for debugging and display purposes.
//
// Example:
//
//	fieldOptional := optional.MakeOptionalWithName(..., "Person.Email")
//	fmt.Println(fieldOptional)  // Prints: "Person.Email"
func (o Optional[S, T]) String() string {
	return o.name
}

// Format implements fmt.Formatter for Optional.
// Supports all standard format verbs:
//   - %s, %v, %+v: uses String() representation (optional name)
//   - %#v: uses GoString() representation
//   - %q: quoted String() representation
//   - other verbs: uses String() representation
//
// Example:
//
//	fieldOptional := optional.MakeOptionalWithName(..., "Person.Email")
//	fmt.Printf("%s", fieldOptional)   // "Person.Email"
//	fmt.Printf("%v", fieldOptional)   // "Person.Email"
//	fmt.Printf("%#v", fieldOptional)  // "optional.Optional[Person, string]{name: \"Person.Email\"}"
//
//go:noinline
func (o Optional[S, T]) Format(f fmt.State, c rune) {
	formatting.FmtString(o, f, c)
}

// GoString implements fmt.GoStringer for Optional.
// Returns a Go-syntax representation of the Optional value.
//
// Example:
//
//	fieldOptional := optional.MakeOptionalWithName(..., "Person.Email")
//	fieldOptional.GoString() // "optional.Optional[Person, string]{name: \"Person.Email\"}"
//
//go:noinline
func (o Optional[S, T]) GoString() string {
	return fmt.Sprintf("optional.Optional[%s, %s]{name: %q}",
		formatting.TypeInfo(new(S)),
		formatting.TypeInfo(new(T)),
		o.name,
	)
}

// LogValue implements slog.LogValuer for Optional.
// Returns a slog.Value that represents the Optional for structured logging.
// Logs the optional name as a string value.
//
// Example:
//
//	logger := slog.Default()
//	fieldOptional := optional.MakeOptionalWithName(..., "Person.Email")
//	logger.Info("using optional", "optional", fieldOptional)
//	// Logs: {"msg":"using optional","optional":"Person.Email"}
//
//go:noinline
func (o Optional[S, T]) LogValue() slog.Value {
	return slog.StringValue(o.name)
}
