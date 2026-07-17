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

package lens

import (
	"fmt"
	"log/slog"

	"github.com/IBM/fp-go/v2/internal/formatting"
)

// The methods below are defined on the non-generic lensName type rather than
// directly on Lens[S, A]. Because Go instantiates a separate copy of every
// generic method for each distinct type-argument combination, placing these
// methods on the embedded non-generic lensName avoids that code bloat: a
// single compiled copy is shared across all Lens[S, A] instantiations.

// String returns the name of the lens for debugging and display purposes.
func (l lensName) String() string {
	return l.n
}

// Format implements fmt.Formatter.
//
// Supports all standard format verbs:
//   - %s, %v, %+v: uses the String() representation (the lens name)
//   - %#v: uses the GoString() representation
//   - %q: quoted String() representation
//   - all other verbs: uses the String() representation
func (l lensName) Format(f fmt.State, c rune) {
	formatting.FmtString(l, f, c)
}

// GoString implements fmt.GoStringer.
//
// Returns a Go-syntax representation of the lens value, suitable for use with
// the %#v format verb.
func (l lensName) GoString() string {
	return fmt.Sprintf("lens.Lens[%s, %s]{name: %q}",
		l.s,
		l.a,
		l.n,
	)

}

// LogValue implements slog.LogValuer.
//
// Returns a slog.Value that represents the lens for structured logging.
// The lens name is logged as a string value.
func (l lensName) LogValue() slog.Value {
	return slog.StringValue(l.n)
}
