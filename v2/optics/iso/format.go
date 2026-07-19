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

package iso

import (
	"fmt"
	"log/slog"
)

// The methods below are defined on the non-generic isoTag type rather than
// directly on Iso[S, A]. Because Go instantiates a separate copy of every
// generic method for each distinct type-argument combination, placing these
// methods on the embedded non-generic isoTag avoids that code bloat: a
// single compiled copy is shared across all Iso[S, A] instantiations.

// String returns a string representation of the isomorphism.
//
// Example:
//
//	tempIso := iso.MakeIso(...)
//	fmt.Println(tempIso)  // Prints: "Iso"
func (isoTag) String() string {
	return "Iso"
}

// Format implements fmt.Formatter for Iso.
// Supports all standard format verbs:
//   - %s, %v, %+v, %q, and all other verbs: uses String() representation
//
//go:noinline
func (isoTag) Format(f fmt.State, c rune) {
	switch c {
	case 'q':
		fmt.Fprintf(f, "%q", "Iso")
	default:
		fmt.Fprint(f, "Iso")
	}
}

// LogValue implements slog.LogValuer for Iso.
// Returns a slog.Value that represents the Iso for structured logging.
// Logs the type information as a string value.
//
// Example:
//
//	logger := slog.Default()
//	tempIso := iso.MakeIso(...)
//	logger.Info("using iso", "iso", tempIso)
//	// Logs: {"msg":"using iso","iso":"Iso"}
//
//go:noinline
func (isoTag) LogValue() slog.Value {
	return slog.StringValue("Iso")
}
