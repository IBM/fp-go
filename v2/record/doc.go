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

// Package record contains monadic operations for maps as well as a rich set of utility functions.
//
// # Nil Map Handling
//
// Throughout this package, nil maps are treated identically to empty maps, consistent with
// Go's native map behavior. This means:
//   - A nil map has length 0 (len(nil) == 0)
//   - Iterating over a nil map produces zero iterations
//   - Lookup operations on nil maps return the zero value and false
//   - All record operations handle nil maps safely without panics
//
// Most transformation functions (Map, Filter, etc.) return non-nil empty maps when given
// nil input, while query functions (IsEmpty, Size, etc.) treat nil as empty.
//
// Example:
//
//	var nilMap record.Record[string, int]  // nil map
//	emptyMap := record.Record[string, int]{}  // empty but non-nil map
//
//	record.IsEmpty(nilMap)   // true
//	record.IsEmpty(emptyMap) // true
//	record.Size(nilMap)      // 0
//	record.Size(emptyMap)    // 0
//
//	// Transformations return non-nil empty maps
//	result := record.Map(func(x int) int { return x * 2 })(nilMap)
//	// result is non-nil but empty: map[string]int{}
package record
