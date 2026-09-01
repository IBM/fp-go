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
// Returns a structured log value representing the current trampoline state.
// The exact structure of the returned slog.Value is not a stable contract and may change across versions.
//
// Example:
//
//	trampoline := tailrec.Bounce[int](42)
//	slog.Info("Processing", "state", trampoline)
//	// Logs: {"level":"info","msg":"Processing","state":{"bouncing":42}}
//
//	result := tailrec.Land[int](100)
//	slog.Info("Complete", "result", result)
//	// Logs: {"level":"info","msg":"Complete","result":{"landed":100}}
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
