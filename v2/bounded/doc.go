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

// Package bounded provides types and functions for working with bounded ordered types.
//
// A Bounded type extends Ord with minimum (Bottom) and maximum (Top) values,
// representing types that have well-defined lower and upper bounds.
//
// # Bounded Interface
//
// The Bounded interface combines ordering (Ord) with boundary values:
//
//	type Bounded[T any] interface {
//	    ord.Ord[T]
//	    Top() T    // Maximum value
//	    Bottom() T // Minimum value
//	}
//
// # Creating Bounded Instances
//
// Use MakeBounded to create a Bounded instance from an Ord and boundary values:
//
//	import (
//	    "github.com/IBM/fp-go/v2/bounded"
//	    "github.com/IBM/fp-go/v2/ord"
//	)
//
//	// Bounded integers from 0 to 100
//	boundedInt := bounded.MakeBounded(
//	    ord.FromStrictCompare[int](),
//	    100, // top
//	    0,   // bottom
//	)
//
//	top := boundedInt.Top()       // 100
//	bottom := boundedInt.Bottom() // 0
//
// # Clamping Values
//
// The Clamp function restricts values to stay within the bounds:
//
//	import (
//	    "github.com/IBM/fp-go/v2/bounded"
//	    "github.com/IBM/fp-go/v2/ord"
//	)
//
//	// Create bounded type for percentages (0-100)
//	percentage := bounded.MakeBounded(
//	    ord.FromStrictCompare[int](),
//	    100, // top
//	    0,   // bottom
//	)
//
//	clamp := bounded.Clamp(percentage)
//
//	result1 := clamp(50)   // 50 (within bounds)
//	result2 := clamp(150)  // 100 (clamped to top)
//	result3 := clamp(-10)  // 0 (clamped to bottom)
//
// # Reversing Bounds
//
// The Reverse function swaps the ordering and bounds:
//
//	import (
//	    "github.com/IBM/fp-go/v2/bounded"
//	    "github.com/IBM/fp-go/v2/ord"
//	)
//
//	original := bounded.MakeBounded(
//	    ord.FromStrictCompare[int](),
//	    100, // top
//	    0,   // bottom
//	)
//
//	reversed := bounded.Reverse(original)
//
//	// In reversed, ordering is flipped and bounds are swapped
//	// Compare(10, 20) returns 1 instead of -1
//	// Top() returns 0 and Bottom() returns 100
//
// # Use Cases
//
// Bounded types are useful for:
//
//   - Representing ranges with well-defined limits (e.g., percentages, grades)
//   - Implementing safe arithmetic that stays within bounds
//   - Validating input values against constraints
//   - Creating domain-specific types with natural boundaries
//
// # Example - Temperature Range
//
//	import (
//	    "github.com/IBM/fp-go/v2/bounded"
//	    "github.com/IBM/fp-go/v2/ord"
//	)
//
//	// Celsius temperature range for a thermostat
//	thermostat := bounded.MakeBounded(
//	    ord.FromStrictCompare[float64](),
//	    30.0, // max temperature
//	    15.0, // min temperature
//	)
//
//	clampTemp := bounded.Clamp(thermostat)
//
//	// User tries to set temperature
//	desired := 35.0
//	actual := clampTemp(desired) // 30.0 (clamped to maximum)
//
// # Example - Bounded Characters
//
//	import (
//	    "github.com/IBM/fp-go/v2/bounded"
//	    "github.com/IBM/fp-go/v2/ord"
//	)
//
//	// Lowercase letters only
//	lowercase := bounded.MakeBounded(
//	    ord.FromStrictCompare[rune](),
//	    'z', // top
//	    'a', // bottom
//	)
//
//	clampChar := bounded.Clamp(lowercase)
//
//	result1 := clampChar('m')  // 'm' (within bounds)
//	result2 := clampChar('A')  // 'a' (clamped to bottom)
//	result3 := clampChar('~')  // 'z' (clamped to top)
//
// # Laws
//
// Bounded instances must satisfy the Ord laws plus:
//
//   - Bottom is less than or equal to all values: Compare(Bottom(), x) <= 0
//   - Top is greater than or equal to all values: Compare(Top(), x) >= 0
//   - Bottom <= Top: Compare(Bottom(), Top()) <= 0
package bounded
