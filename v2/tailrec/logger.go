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

package tailrec

import (
	"log/slog"
)

// LogValue implements the slog.LogValuer interface for Trampoline.
//
// This method allows Trampoline values to be logged using Go's structured logging
// (log/slog) with proper representation of their state:
//   - When Landed is true: returns a group with a single "landed" attribute containing the Land value
//   - When Landed is false: returns a group with a single "bouncing" attribute containing the Bounce value
//
// The implementation ensures that Trampoline values are logged in a structured,
// readable format that clearly shows the current state of the tail-recursive computation.
//
// Example usage:
//
//	trampoline := tailrec.Bounce[int](42)
//	slog.Info("Processing", "state", trampoline)
//	// Logs: {"level":"info","msg":"Processing","state":{"bouncing":42}}
//
//	result := tailrec.Land[int](100)
//	slog.Info("Complete", "result", result)
//	// Logs: {"level":"info","msg":"Complete","result":{"landed":100}}
//
// This is particularly useful for debugging tail-recursive computations and
// understanding the flow of recursive algorithms at runtime.
func (t Trampoline[B, L]) LogValue() slog.Value {
	if t.Landed {
		return slog.GroupValue(
			slog.Any("landed", t.Land),
		)
	}
	return slog.GroupValue(
		slog.Any("bouncing", t.Bounce),
	)
}
