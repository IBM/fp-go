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

package pair

import (
	"fmt"
	"log/slog"
)

const (
	pairFmtTemplate = "Pair[%T, %T](%v, %v)"
)

// String prints some debug info for the object
//

func (p Pair[L, R]) String() string {
	return fmt.Sprintf(pairFmtTemplate, p.l, p.r, p.l, p.r)
}

// Format implements fmt.Formatter for Pair.
// Supports all standard format verbs:
//   - %s, %v, %+v, %q, and all other verbs: uses String() representation
//

func (p Pair[L, R]) Format(f fmt.State, c rune) {
	switch c {
	case 'q':
		fmt.Fprintf(f, "%q", p.String())
	default:
		fmt.Fprint(f, p.String())
	}
}

// LogValue implements slog.LogValuer for Pair.
// Returns a slog.Value that represents the Pair for structured logging.
// Returns a group value with "head" and "tail" keys.
//
// Example:
//
//	logger := slog.Default()
//	p := pair.MakePair("key", 42)
//	logger.Info("pair value", "data", p)
//	// Logs: {"msg":"pair value","data":{"head":"key","tail":42}}
//

func (p Pair[L, R]) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("head", p.l),
		slog.Any("tail", p.r),
	)
}
