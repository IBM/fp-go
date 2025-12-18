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
	"fmt"
	"log/slog"

	"github.com/IBM/fp-go/v2/internal/formatting"
)

// String returns the name of the prism for debugging and display purposes.
//
// Example:
//
//	successPrism := prism.MakePrismWithName(..., "Result.Success")
//	fmt.Println(successPrism)  // Prints: "Result.Success"
func (p Prism[S, T]) String() string {
	return p.name
}

// Format implements fmt.Formatter for Prism.
// Supports all standard format verbs:
//   - %s, %v, %+v: uses String() representation (prism name)
//   - %#v: uses GoString() representation
//   - %q: quoted String() representation
//   - other verbs: uses String() representation
//
// Example:
//
//	successPrism := prism.MakePrismWithName(..., "Result.Success")
//	fmt.Printf("%s", successPrism)   // "Result.Success"
//	fmt.Printf("%v", successPrism)   // "Result.Success"
//	fmt.Printf("%#v", successPrism)  // "prism.Prism[Result, int]{name: \"Result.Success\"}"
//
//go:noinline
func (p Prism[S, T]) Format(f fmt.State, c rune) {
	formatting.FmtString(p, f, c)
}

// GoString implements fmt.GoStringer for Prism.
// Returns a Go-syntax representation of the Prism value.
//
// Example:
//
//	successPrism := prism.MakePrismWithName(..., "Result.Success")
//	successPrism.GoString() // "prism.Prism[Result, int]{name: \"Result.Success\"}"
//
//go:noinline
func (p Prism[S, T]) GoString() string {
	return fmt.Sprintf("prism.Prism[%s, %s]{name: %q}",
		formatting.TypeInfo(new(S)),
		formatting.TypeInfo(new(T)),
		p.name,
	)
}

// LogValue implements slog.LogValuer for Prism.
// Returns a slog.Value that represents the Prism for structured logging.
// Logs the prism name as a string value.
//
// Example:
//
//	logger := slog.Default()
//	successPrism := prism.MakePrismWithName(..., "Result.Success")
//	logger.Info("using prism", "prism", successPrism)
//	// Logs: {"msg":"using prism","prism":"Result.Success"}
//
//go:noinline
func (p Prism[S, T]) LogValue() slog.Value {
	return slog.StringValue(p.name)
}
