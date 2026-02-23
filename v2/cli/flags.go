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

package cli

import (
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	C "github.com/urfave/cli/v3"
)

// StringFlagPrism creates a Prism for extracting a StringFlag from a Flag.
// This provides a type-safe way to work with string flags, handling type
// mismatches gracefully through the Option type.
//
// The prism's GetOption attempts to cast a Flag to *C.StringFlag.
// If the cast succeeds, it returns Some(*C.StringFlag); if it fails, it returns None.
//
// The prism's ReverseGet converts a *C.StringFlag back to a Flag.
//
// # Returns
//
//   - A Prism[C.Flag, *C.StringFlag] for safe StringFlag extraction
//
// # Example Usage
//
//	prism := StringFlagPrism()
//
//	// Extract StringFlag from Flag
//	var flag C.Flag = &C.StringFlag{Name: "input", Value: "default"}
//	result := prism.GetOption(flag)  // Some(*C.StringFlag{...})
//
//	// Type mismatch returns None
//	var intFlag C.Flag = &C.IntFlag{Name: "count"}
//	result = prism.GetOption(intFlag)  // None[*C.StringFlag]()
//
//	// Convert back to Flag
//	strFlag := &C.StringFlag{Name: "output"}
//	flag = prism.ReverseGet(strFlag)
func StringFlagPrism() P.Prism[C.Flag, *C.StringFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.StringFlag] {
			if sf, ok := flag.(*C.StringFlag); ok {
				return O.Some(sf)
			}
			return O.None[*C.StringFlag]()
		},
		func(f *C.StringFlag) C.Flag { return f },
	)
}

// IntFlagPrism creates a Prism for extracting an IntFlag from a Flag.
// This provides a type-safe way to work with integer flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.IntFlag] for safe IntFlag extraction
//
// # Example Usage
//
//	prism := IntFlagPrism()
//
//	// Extract IntFlag from Flag
//	var flag C.Flag = &C.IntFlag{Name: "count", Value: 10}
//	result := prism.GetOption(flag)  // Some(*C.IntFlag{...})
func IntFlagPrism() P.Prism[C.Flag, *C.IntFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.IntFlag] {
			if f, ok := flag.(*C.IntFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.IntFlag]()
		},
		func(f *C.IntFlag) C.Flag { return f },
	)
}

// BoolFlagPrism creates a Prism for extracting a BoolFlag from a Flag.
// This provides a type-safe way to work with boolean flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.BoolFlag] for safe BoolFlag extraction
//
// # Example Usage
//
//	prism := BoolFlagPrism()
//
//	// Extract BoolFlag from Flag
//	var flag C.Flag = &C.BoolFlag{Name: "verbose", Value: true}
//	result := prism.GetOption(flag)  // Some(*C.BoolFlag{...})
func BoolFlagPrism() P.Prism[C.Flag, *C.BoolFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.BoolFlag] {
			if f, ok := flag.(*C.BoolFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.BoolFlag]()
		},
		func(f *C.BoolFlag) C.Flag { return f },
	)
}

// Float64FlagPrism creates a Prism for extracting a Float64Flag from a Flag.
// This provides a type-safe way to work with float64 flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.Float64Flag] for safe Float64Flag extraction
//
// # Example Usage
//
//	prism := Float64FlagPrism()
//
//	// Extract Float64Flag from Flag
//	var flag C.Flag = &C.Float64Flag{Name: "ratio", Value: 0.5}
//	result := prism.GetOption(flag)  // Some(*C.Float64Flag{...})
func Float64FlagPrism() P.Prism[C.Flag, *C.Float64Flag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.Float64Flag] {
			if f, ok := flag.(*C.Float64Flag); ok {
				return O.Some(f)
			}
			return O.None[*C.Float64Flag]()
		},
		func(f *C.Float64Flag) C.Flag { return f },
	)
}

// DurationFlagPrism creates a Prism for extracting a DurationFlag from a Flag.
// This provides a type-safe way to work with duration flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.DurationFlag] for safe DurationFlag extraction
//
// # Example Usage
//
//	prism := DurationFlagPrism()
//
//	// Extract DurationFlag from Flag
//	var flag C.Flag = &C.DurationFlag{Name: "timeout", Value: 30 * time.Second}
//	result := prism.GetOption(flag)  // Some(*C.DurationFlag{...})
func DurationFlagPrism() P.Prism[C.Flag, *C.DurationFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.DurationFlag] {
			if f, ok := flag.(*C.DurationFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.DurationFlag]()
		},
		func(f *C.DurationFlag) C.Flag { return f },
	)
}

// TimestampFlagPrism creates a Prism for extracting a TimestampFlag from a Flag.
// This provides a type-safe way to work with timestamp flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.TimestampFlag] for safe TimestampFlag extraction
//
// # Example Usage
//
//	prism := TimestampFlagPrism()
//
//	// Extract TimestampFlag from Flag
//	var flag C.Flag = &C.TimestampFlag{Name: "created"}
//	result := prism.GetOption(flag)  // Some(*C.TimestampFlag{...})
func TimestampFlagPrism() P.Prism[C.Flag, *C.TimestampFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.TimestampFlag] {
			if f, ok := flag.(*C.TimestampFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.TimestampFlag]()
		},
		func(f *C.TimestampFlag) C.Flag { return f },
	)
}

// StringSliceFlagPrism creates a Prism for extracting a StringSliceFlag from a Flag.
// This provides a type-safe way to work with string slice flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.StringSliceFlag] for safe StringSliceFlag extraction
//
// # Example Usage
//
//	prism := StringSliceFlagPrism()
//
//	// Extract StringSliceFlag from Flag
//	var flag C.Flag = &C.StringSliceFlag{Name: "tags"}
//	result := prism.GetOption(flag)  // Some(*C.StringSliceFlag{...})
func StringSliceFlagPrism() P.Prism[C.Flag, *C.StringSliceFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.StringSliceFlag] {
			if f, ok := flag.(*C.StringSliceFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.StringSliceFlag]()
		},
		func(f *C.StringSliceFlag) C.Flag { return f },
	)
}

// IntSliceFlagPrism creates a Prism for extracting an IntSliceFlag from a Flag.
// This provides a type-safe way to work with int slice flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.IntSliceFlag] for safe IntSliceFlag extraction
//
// # Example Usage
//
//	prism := IntSliceFlagPrism()
//
//	// Extract IntSliceFlag from Flag
//	var flag C.Flag = &C.IntSliceFlag{Name: "ports"}
//	result := prism.GetOption(flag)  // Some(*C.IntSliceFlag{...})
func IntSliceFlagPrism() P.Prism[C.Flag, *C.IntSliceFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.IntSliceFlag] {
			if f, ok := flag.(*C.IntSliceFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.IntSliceFlag]()
		},
		func(f *C.IntSliceFlag) C.Flag { return f },
	)
}

// Float64SliceFlagPrism creates a Prism for extracting a Float64SliceFlag from a Flag.
// This provides a type-safe way to work with float64 slice flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.Float64SliceFlag] for safe Float64SliceFlag extraction
//
// # Example Usage
//
//	prism := Float64SliceFlagPrism()
//
//	// Extract Float64SliceFlag from Flag
//	var flag C.Flag = &C.Float64SliceFlag{Name: "ratios"}
//	result := prism.GetOption(flag)  // Some(*C.Float64SliceFlag{...})
func Float64SliceFlagPrism() P.Prism[C.Flag, *C.Float64SliceFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.Float64SliceFlag] {
			if f, ok := flag.(*C.Float64SliceFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.Float64SliceFlag]()
		},
		func(f *C.Float64SliceFlag) C.Flag { return f },
	)
}

// UintFlagPrism creates a Prism for extracting a UintFlag from a Flag.
// This provides a type-safe way to work with unsigned integer flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.UintFlag] for safe UintFlag extraction
//
// # Example Usage
//
//	prism := UintFlagPrism()
//
//	// Extract UintFlag from Flag
//	var flag C.Flag = &C.UintFlag{Name: "workers", Value: 4}
//	result := prism.GetOption(flag)  // Some(*C.UintFlag{...})
func UintFlagPrism() P.Prism[C.Flag, *C.UintFlag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.UintFlag] {
			if f, ok := flag.(*C.UintFlag); ok {
				return O.Some(f)
			}
			return O.None[*C.UintFlag]()
		},
		func(f *C.UintFlag) C.Flag { return f },
	)
}

// Uint64FlagPrism creates a Prism for extracting a Uint64Flag from a Flag.
// This provides a type-safe way to work with uint64 flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.Uint64Flag] for safe Uint64Flag extraction
//
// # Example Usage
//
//	prism := Uint64FlagPrism()
//
//	// Extract Uint64Flag from Flag
//	var flag C.Flag = &C.Uint64Flag{Name: "size"}
//	result := prism.GetOption(flag)  // Some(*C.Uint64Flag{...})
func Uint64FlagPrism() P.Prism[C.Flag, *C.Uint64Flag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.Uint64Flag] {
			if f, ok := flag.(*C.Uint64Flag); ok {
				return O.Some(f)
			}
			return O.None[*C.Uint64Flag]()
		},
		func(f *C.Uint64Flag) C.Flag { return f },
	)
}

// Int64FlagPrism creates a Prism for extracting an Int64Flag from a Flag.
// This provides a type-safe way to work with int64 flags, handling type
// mismatches gracefully through the Option type.
//
// # Returns
//
//   - A Prism[C.Flag, *C.Int64Flag] for safe Int64Flag extraction
//
// # Example Usage
//
//	prism := Int64FlagPrism()
//
//	// Extract Int64Flag from Flag
//	var flag C.Flag = &C.Int64Flag{Name: "offset"}
//	result := prism.GetOption(flag)  // Some(*C.Int64Flag{...})
func Int64FlagPrism() P.Prism[C.Flag, *C.Int64Flag] {
	return P.MakePrism(
		func(flag C.Flag) O.Option[*C.Int64Flag] {
			if f, ok := flag.(*C.Int64Flag); ok {
				return O.Some(f)
			}
			return O.None[*C.Int64Flag]()
		},
		func(f *C.Int64Flag) C.Flag { return f },
	)
}

// Made with Bob
