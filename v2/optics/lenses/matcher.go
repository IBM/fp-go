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

package lenses

import (
	__iso_option "github.com/IBM/fp-go/v2/optics/iso/option"
	__lens "github.com/IBM/fp-go/v2/optics/lens"
	__lens_option "github.com/IBM/fp-go/v2/optics/lens/option"
	__prism "github.com/IBM/fp-go/v2/optics/prism"
	__option "github.com/IBM/fp-go/v2/option"
)

// MatchLenses provides lenses for accessing and modifying fields of Match structures.
// Lenses enable functional updates to immutable data structures by providing
// composable getters and setters.
//
// This struct contains both mandatory lenses (for direct field access) and optional
// lenses (for fields that may be zero values, treating them as Option types).
//
// Fields:
//   - Before: Lens for the text before the match
//   - Groups: Lens for the capture groups array
//   - After: Lens for the text after the match
//   - BeforeO: Optional lens treating empty Before as None
//   - AfterO: Optional lens treating empty After as None
//
// Example:
//
//	lenses := MakeMatchLenses()
//	match := Match{Before: "hello ", Groups: []string{"world"}, After: "!"}
//
//	// Get a field value
//	before := lenses.Before.Get(match) // "hello "
//
//	// Set a field value (returns new Match)
//	updated := lenses.Before.Set(match, "hi ")
//	// updated.Before is now "hi "
type MatchLenses struct {
	// mandatory fields
	Before __lens.Lens[__prism.Match, string]
	Groups __lens.Lens[__prism.Match, []string]
	After  __lens.Lens[__prism.Match, string]
	// optional fields
	BeforeO __lens_option.LensO[__prism.Match, string]
	AfterO  __lens_option.LensO[__prism.Match, string]
}

// MatchRefLenses provides lenses for accessing and modifying fields of Match structures
// via pointers. This is useful when working with mutable references to Match values.
//
// In addition to standard lenses, this struct also includes prisms for each field,
// which provide optional access patterns (useful for validation or conditional updates).
//
// Fields:
//   - Before, Groups, After: Standard lenses for pointer-based access
//   - BeforeO, AfterO: Optional lenses treating zero values as None
//   - BeforeP, GroupsP, AfterP: Prisms for optional field access
//
// Example:
//
//	lenses := MakeMatchRefLenses()
//	match := &Match{Before: "hello ", Groups: []string{"world"}, After: "!"}
//
//	// Get a field value
//	before := lenses.Before.Get(match) // "hello "
//
//	// Set a field value (mutates the pointer)
//	lenses.Before.Set(match, "hi ")
//	// match.Before is now "hi "
type MatchRefLenses struct {
	// mandatory fields
	Before __lens.Lens[*__prism.Match, string]
	Groups __lens.Lens[*__prism.Match, []string]
	After  __lens.Lens[*__prism.Match, string]
	// optional fields
	BeforeO __lens_option.LensO[*__prism.Match, string]
	AfterO  __lens_option.LensO[*__prism.Match, string]
	// prisms
	BeforeP __prism.Prism[*__prism.Match, string]
	GroupsP __prism.Prism[*__prism.Match, []string]
	AfterP  __prism.Prism[*__prism.Match, string]
}

// MatchPrisms provides prisms for accessing fields of Match structures.
// Prisms enable safe access to fields that may not be present (zero values),
// returning Option types instead of direct values.
//
// Fields:
//   - Before: Prism for Before field (None if empty string)
//   - Groups: Prism for Groups field (always Some)
//   - After: Prism for After field (None if empty string)
//
// Example:
//
//	prisms := MakeMatchPrisms()
//	match := Match{Before: "", Groups: []string{"test"}, After: "!"}
//
//	// Try to get Before (returns None because it's empty)
//	beforeOpt := prisms.Before.GetOption(match) // None
//
//	// Get After (returns Some because it's non-empty)
//	afterOpt := prisms.After.GetOption(match) // Some("!")
type MatchPrisms struct {
	Before __prism.Prism[__prism.Match, string]
	Groups __prism.Prism[__prism.Match, []string]
	After  __prism.Prism[__prism.Match, string]
}

// MakeMatchLenses creates a new MatchLenses with lenses for all fields of Match.
// This function constructs both mandatory lenses (for direct field access) and
// optional lenses (for treating zero values as Option types).
//
// The returned lenses enable functional-style updates to Match structures,
// allowing you to get and set field values while maintaining immutability.
//
// Returns:
//   - A MatchLenses struct with lenses for Before, Groups, After fields,
//     plus optional lenses BeforeO and AfterO
//
// Example:
//
//	lenses := MakeMatchLenses()
//	match := Match{Before: "start ", Groups: []string{"middle"}, After: " end"}
//
//	// Get field values
//	before := lenses.Before.Get(match)        // "start "
//	groups := lenses.Groups.Get(match)        // []string{"middle"}
//
//	// Update a field (returns new Match)
//	updated := lenses.After.Set(match, " finish")
//	// updated is a new Match with After = " finish"
//
//	// Use optional lens (treats empty string as None)
//	emptyMatch := Match{Before: "", Groups: []string{}, After: ""}
//	beforeOpt := lenses.BeforeO.GetOption(emptyMatch) // None
func MakeMatchLenses() MatchLenses {
	// mandatory lenses
	lensBefore := __lens.MakeLensWithName(
		func(s __prism.Match) string { return s.Before },
		func(s __prism.Match, v string) __prism.Match { s.Before = v; return s },
		"Match.Before",
	)
	lensGroups := __lens.MakeLensWithName(
		func(s __prism.Match) []string { return s.Groups },
		func(s __prism.Match, v []string) __prism.Match { s.Groups = v; return s },
		"Match.Groups",
	)
	lensAfter := __lens.MakeLensWithName(
		func(s __prism.Match) string { return s.After },
		func(s __prism.Match, v string) __prism.Match { s.After = v; return s },
		"Match.After",
	)
	// optional lenses
	lensBeforeO := __lens_option.FromIso[__prism.Match](__iso_option.FromZero[string]())(lensBefore)
	lensAfterO := __lens_option.FromIso[__prism.Match](__iso_option.FromZero[string]())(lensAfter)
	return MatchLenses{
		// mandatory lenses
		Before: lensBefore,
		Groups: lensGroups,
		After:  lensAfter,
		// optional lenses
		BeforeO: lensBeforeO,
		AfterO:  lensAfterO,
	}
}

// MakeMatchRefLenses creates a new MatchRefLenses with lenses for all fields of *Match.
// This function constructs lenses that work with pointers to Match structures,
// enabling both immutable-style updates and direct mutations.
//
// The returned lenses include:
//   - Standard lenses for pointer-based field access
//   - Optional lenses for treating zero values as Option types
//   - Prisms for safe optional field access
//
// Returns:
//   - A MatchRefLenses struct with lenses and prisms for all Match fields
//
// Example:
//
//	lenses := MakeMatchRefLenses()
//	match := &Match{Before: "prefix ", Groups: []string{"data"}, After: " suffix"}
//
//	// Get field value
//	before := lenses.Before.Get(match) // "prefix "
//
//	// Set field value (mutates the pointer)
//	lenses.Before.Set(match, "new ")
//	// match.Before is now "new "
//
//	// Use prism for optional access
//	afterOpt := lenses.AfterP.GetOption(match) // Some(" suffix")
func MakeMatchRefLenses() MatchRefLenses {
	// mandatory lenses
	lensBefore := __lens.MakeLensStrictWithName(
		func(s *__prism.Match) string { return s.Before },
		func(s *__prism.Match, v string) *__prism.Match { s.Before = v; return s },
		"(*Match).Before",
	)
	lensGroups := __lens.MakeLensRefWithName(
		func(s *__prism.Match) []string { return s.Groups },
		func(s *__prism.Match, v []string) *__prism.Match { s.Groups = v; return s },
		"(*Match).Groups",
	)
	lensAfter := __lens.MakeLensStrictWithName(
		func(s *__prism.Match) string { return s.After },
		func(s *__prism.Match, v string) *__prism.Match { s.After = v; return s },
		"(*Match).After",
	)
	// optional lenses
	lensBeforeO := __lens_option.FromIso[*__prism.Match](__iso_option.FromZero[string]())(lensBefore)
	lensAfterO := __lens_option.FromIso[*__prism.Match](__iso_option.FromZero[string]())(lensAfter)
	return MatchRefLenses{
		// mandatory lenses
		Before: lensBefore,
		Groups: lensGroups,
		After:  lensAfter,
		// optional lenses
		BeforeO: lensBeforeO,
		AfterO:  lensAfterO,
	}
}

// MakeMatchPrisms creates a new MatchPrisms with prisms for all fields of Match.
// This function constructs prisms that provide safe optional access to Match fields,
// treating zero values (empty strings) as None.
//
// The returned prisms enable pattern matching on field presence:
//   - Before and After prisms return None for empty strings
//   - Groups prism always returns Some (even for empty slices)
//
// Returns:
//   - A MatchPrisms struct with prisms for Before, Groups, and After fields
//
// Example:
//
//	prisms := MakeMatchPrisms()
//	match := Match{Before: "", Groups: []string{"data"}, After: "!"}
//
//	// Try to get Before (returns None because it's empty)
//	beforeOpt := prisms.Before.GetOption(match) // None
//
//	// Get Groups (always returns Some)
//	groupsOpt := prisms.Groups.GetOption(match) // Some([]string{"data"})
//
//	// Get After (returns Some because it's non-empty)
//	afterOpt := prisms.After.GetOption(match) // Some("!")
//
//	// Construct a Match from a value using ReverseGet
//	newMatch := prisms.Before.ReverseGet("prefix ")
//	// newMatch is Match{Before: "prefix ", Groups: nil, After: ""}
func MakeMatchPrisms() MatchPrisms {
	_fromNonZeroBefore := __option.FromNonZero[string]()
	_prismBefore := __prism.MakePrismWithName(
		func(s __prism.Match) __option.Option[string] { return _fromNonZeroBefore(s.Before) },
		func(v string) __prism.Match {
			return __prism.Match{Before: v}
		},
		"Match.Before",
	)
	_prismGroups := __prism.MakePrismWithName(
		func(s __prism.Match) __option.Option[[]string] { return __option.Some(s.Groups) },
		func(v []string) __prism.Match {
			return __prism.Match{Groups: v}
		},
		"Match.Groups",
	)
	_fromNonZeroAfter := __option.FromNonZero[string]()
	_prismAfter := __prism.MakePrismWithName(
		func(s __prism.Match) __option.Option[string] { return _fromNonZeroAfter(s.After) },
		func(v string) __prism.Match {
			return __prism.Match{After: v}
		},
		"Match.After",
	)
	return MatchPrisms{
		Before: _prismBefore,
		Groups: _prismGroups,
		After:  _prismAfter,
	}
}

// NamedMatchLenses provides lenses for accessing and modifying fields of NamedMatch structures.
// NamedMatch represents regex matches with named capture groups, and these lenses enable
// functional updates to its fields.
//
// This struct contains both mandatory lenses (for direct field access) and optional
// lenses (for fields that may be zero values, treating them as Option types).
//
// Fields:
//   - Before: Lens for the text before the match
//   - Groups: Lens for the named capture groups map
//   - Full: Lens for the complete matched text
//   - After: Lens for the text after the match
//   - BeforeO, FullO, AfterO: Optional lenses treating empty strings as None
//
// Example:
//
//	lenses := MakeNamedMatchLenses()
//	match := NamedMatch{
//	    Before: "Email: ",
//	    Groups: map[string]string{"user": "john", "domain": "example.com"},
//	    Full: "john@example.com",
//	    After: "",
//	}
//
//	// Get a field value
//	full := lenses.Full.Get(match) // "john@example.com"
//
//	// Set a field value (returns new NamedMatch)
//	updated := lenses.Before.Set(match, "Contact: ")
//	// updated.Before is now "Contact: "
type NamedMatchLenses struct {
	// mandatory fields
	Before __lens.Lens[__prism.NamedMatch, string]
	Groups __lens.Lens[__prism.NamedMatch, map[string]string]
	Full   __lens.Lens[__prism.NamedMatch, string]
	After  __lens.Lens[__prism.NamedMatch, string]
	// optional fields
	BeforeO __lens_option.LensO[__prism.NamedMatch, string]
	FullO   __lens_option.LensO[__prism.NamedMatch, string]
	AfterO  __lens_option.LensO[__prism.NamedMatch, string]
}

// NamedMatchRefLenses provides lenses for accessing and modifying fields of NamedMatch
// structures via pointers. This is useful when working with mutable references to
// NamedMatch values.
//
// In addition to standard lenses, this struct also includes prisms for each field,
// which provide optional access patterns (useful for validation or conditional updates).
//
// Fields:
//   - Before, Groups, Full, After: Standard lenses for pointer-based access
//   - BeforeO, FullO, AfterO: Optional lenses treating zero values as None
//   - BeforeP, GroupsP, FullP, AfterP: Prisms for optional field access
//
// Example:
//
//	lenses := MakeNamedMatchRefLenses()
//	match := &NamedMatch{
//	    Before: "Email: ",
//	    Groups: map[string]string{"user": "alice", "domain": "test.org"},
//	    Full: "alice@test.org",
//	    After: " for info",
//	}
//
//	// Get a field value
//	full := lenses.Full.Get(match) // "alice@test.org"
//
//	// Set a field value (mutates the pointer)
//	lenses.After.Set(match, " for contact")
//	// match.After is now " for contact"
type NamedMatchRefLenses struct {
	// mandatory fields
	Before __lens.Lens[*__prism.NamedMatch, string]
	Groups __lens.Lens[*__prism.NamedMatch, map[string]string]
	Full   __lens.Lens[*__prism.NamedMatch, string]
	After  __lens.Lens[*__prism.NamedMatch, string]
	// optional fields
	BeforeO __lens_option.LensO[*__prism.NamedMatch, string]
	FullO   __lens_option.LensO[*__prism.NamedMatch, string]
	AfterO  __lens_option.LensO[*__prism.NamedMatch, string]
	// prisms
	BeforeP __prism.Prism[*__prism.NamedMatch, string]
	GroupsP __prism.Prism[*__prism.NamedMatch, map[string]string]
	FullP   __prism.Prism[*__prism.NamedMatch, string]
	AfterP  __prism.Prism[*__prism.NamedMatch, string]
}

// NamedMatchPrisms provides prisms for accessing fields of NamedMatch structures.
// Prisms enable safe access to fields that may not be present (zero values),
// returning Option types instead of direct values.
//
// Fields:
//   - Before: Prism for Before field (None if empty string)
//   - Groups: Prism for Groups field (always Some)
//   - Full: Prism for Full field (None if empty string)
//   - After: Prism for After field (None if empty string)
//
// Example:
//
//	prisms := MakeNamedMatchPrisms()
//	match := NamedMatch{
//	    Before: "",
//	    Groups: map[string]string{"user": "bob"},
//	    Full: "bob@example.com",
//	    After: "",
//	}
//
//	// Try to get Before (returns None because it's empty)
//	beforeOpt := prisms.Before.GetOption(match) // None
//
//	// Get Full (returns Some because it's non-empty)
//	fullOpt := prisms.Full.GetOption(match) // Some("bob@example.com")
type NamedMatchPrisms struct {
	Before __prism.Prism[__prism.NamedMatch, string]
	Groups __prism.Prism[__prism.NamedMatch, map[string]string]
	Full   __prism.Prism[__prism.NamedMatch, string]
	After  __prism.Prism[__prism.NamedMatch, string]
}

// MakeNamedMatchLenses creates a new NamedMatchLenses with lenses for all fields of NamedMatch.
// This function constructs both mandatory lenses (for direct field access) and
// optional lenses (for treating zero values as Option types).
//
// The returned lenses enable functional-style updates to NamedMatch structures,
// allowing you to get and set field values while maintaining immutability.
//
// Returns:
//   - A NamedMatchLenses struct with lenses for Before, Groups, Full, After fields,
//     plus optional lenses BeforeO, FullO, and AfterO
//
// Example:
//
//	lenses := MakeNamedMatchLenses()
//	match := NamedMatch{
//	    Before: "Email: ",
//	    Groups: map[string]string{"user": "john", "domain": "example.com"},
//	    Full: "john@example.com",
//	    After: "",
//	}
//
//	// Get field values
//	full := lenses.Full.Get(match)        // "john@example.com"
//	groups := lenses.Groups.Get(match)    // map[string]string{"user": "john", ...}
//
//	// Update a field (returns new NamedMatch)
//	updated := lenses.Before.Set(match, "Contact: ")
//	// updated is a new NamedMatch with Before = "Contact: "
//
//	// Use optional lens (treats empty string as None)
//	afterOpt := lenses.AfterO.GetOption(match) // None (because After is empty)
func MakeNamedMatchLenses() NamedMatchLenses {
	// mandatory lenses
	lensBefore := __lens.MakeLensWithName(
		func(s __prism.NamedMatch) string { return s.Before },
		func(s __prism.NamedMatch, v string) __prism.NamedMatch { s.Before = v; return s },
		"NamedMatch.Before",
	)
	lensGroups := __lens.MakeLensWithName(
		func(s __prism.NamedMatch) map[string]string { return s.Groups },
		func(s __prism.NamedMatch, v map[string]string) __prism.NamedMatch { s.Groups = v; return s },
		"NamedMatch.Groups",
	)
	lensFull := __lens.MakeLensWithName(
		func(s __prism.NamedMatch) string { return s.Full },
		func(s __prism.NamedMatch, v string) __prism.NamedMatch { s.Full = v; return s },
		"NamedMatch.Full",
	)
	lensAfter := __lens.MakeLensWithName(
		func(s __prism.NamedMatch) string { return s.After },
		func(s __prism.NamedMatch, v string) __prism.NamedMatch { s.After = v; return s },
		"NamedMatch.After",
	)
	// optional lenses
	lensBeforeO := __lens_option.FromIso[__prism.NamedMatch](__iso_option.FromZero[string]())(lensBefore)
	lensFullO := __lens_option.FromIso[__prism.NamedMatch](__iso_option.FromZero[string]())(lensFull)
	lensAfterO := __lens_option.FromIso[__prism.NamedMatch](__iso_option.FromZero[string]())(lensAfter)
	return NamedMatchLenses{
		// mandatory lenses
		Before: lensBefore,
		Groups: lensGroups,
		Full:   lensFull,
		After:  lensAfter,
		// optional lenses
		BeforeO: lensBeforeO,
		FullO:   lensFullO,
		AfterO:  lensAfterO,
	}
}

// MakeNamedMatchRefLenses creates a new NamedMatchRefLenses with lenses for all fields of *NamedMatch.
// This function constructs lenses that work with pointers to NamedMatch structures,
// enabling both immutable-style updates and direct mutations.
//
// The returned lenses include:
//   - Standard lenses for pointer-based field access
//   - Optional lenses for treating zero values as Option types
//   - Prisms for safe optional field access
//
// Returns:
//   - A NamedMatchRefLenses struct with lenses and prisms for all NamedMatch fields
//
// Example:
//
//	lenses := MakeNamedMatchRefLenses()
//	match := &NamedMatch{
//	    Before: "Email: ",
//	    Groups: map[string]string{"user": "alice", "domain": "test.org"},
//	    Full: "alice@test.org",
//	    After: "",
//	}
//
//	// Get field value
//	full := lenses.Full.Get(match) // "alice@test.org"
//
//	// Set field value (mutates the pointer)
//	lenses.Before.Set(match, "Contact: ")
//	// match.Before is now "Contact: "
//
//	// Use prism for optional access
//	fullOpt := lenses.FullP.GetOption(match) // Some("alice@test.org")
func MakeNamedMatchRefLenses() NamedMatchRefLenses {
	// mandatory lenses
	lensBefore := __lens.MakeLensStrictWithName(
		func(s *__prism.NamedMatch) string { return s.Before },
		func(s *__prism.NamedMatch, v string) *__prism.NamedMatch { s.Before = v; return s },
		"(*NamedMatch).Before",
	)
	lensGroups := __lens.MakeLensRefWithName(
		func(s *__prism.NamedMatch) map[string]string { return s.Groups },
		func(s *__prism.NamedMatch, v map[string]string) *__prism.NamedMatch { s.Groups = v; return s },
		"(*NamedMatch).Groups",
	)
	lensFull := __lens.MakeLensStrictWithName(
		func(s *__prism.NamedMatch) string { return s.Full },
		func(s *__prism.NamedMatch, v string) *__prism.NamedMatch { s.Full = v; return s },
		"(*NamedMatch).Full",
	)
	lensAfter := __lens.MakeLensStrictWithName(
		func(s *__prism.NamedMatch) string { return s.After },
		func(s *__prism.NamedMatch, v string) *__prism.NamedMatch { s.After = v; return s },
		"(*NamedMatch).After",
	)
	// optional lenses
	lensBeforeO := __lens_option.FromIso[*__prism.NamedMatch](__iso_option.FromZero[string]())(lensBefore)
	lensFullO := __lens_option.FromIso[*__prism.NamedMatch](__iso_option.FromZero[string]())(lensFull)
	lensAfterO := __lens_option.FromIso[*__prism.NamedMatch](__iso_option.FromZero[string]())(lensAfter)
	return NamedMatchRefLenses{
		// mandatory lenses
		Before: lensBefore,
		Groups: lensGroups,
		Full:   lensFull,
		After:  lensAfter,
		// optional lenses
		BeforeO: lensBeforeO,
		FullO:   lensFullO,
		AfterO:  lensAfterO,
	}
}

// MakeNamedMatchPrisms creates a new NamedMatchPrisms with prisms for all fields of NamedMatch.
// This function constructs prisms that provide safe optional access to NamedMatch fields,
// treating zero values (empty strings) as None.
//
// The returned prisms enable pattern matching on field presence:
//   - Before, Full, and After prisms return None for empty strings
//   - Groups prism always returns Some (even for nil or empty maps)
//
// Returns:
//   - A NamedMatchPrisms struct with prisms for Before, Groups, Full, and After fields
//
// Example:
//
//	prisms := MakeNamedMatchPrisms()
//	match := NamedMatch{
//	    Before: "",
//	    Groups: map[string]string{"user": "bob", "domain": "example.com"},
//	    Full: "bob@example.com",
//	    After: "",
//	}
//
//	// Try to get Before (returns None because it's empty)
//	beforeOpt := prisms.Before.GetOption(match) // None
//
//	// Get Groups (always returns Some)
//	groupsOpt := prisms.Groups.GetOption(match)
//	// Some(map[string]string{"user": "bob", "domain": "example.com"})
//
//	// Get Full (returns Some because it's non-empty)
//	fullOpt := prisms.Full.GetOption(match) // Some("bob@example.com")
//
//	// Construct a NamedMatch from a value using ReverseGet
//	newMatch := prisms.Full.ReverseGet("test@example.com")
//	// newMatch is NamedMatch{Before: "", Groups: nil, Full: "test@example.com", After: ""}
func MakeNamedMatchPrisms() NamedMatchPrisms {
	_fromNonZeroBefore := __option.FromNonZero[string]()
	_prismBefore := __prism.MakePrismWithName(
		func(s __prism.NamedMatch) __option.Option[string] { return _fromNonZeroBefore(s.Before) },
		func(v string) __prism.NamedMatch {
			return __prism.NamedMatch{Before: v}
		},
		"NamedMatch.Before",
	)
	_prismGroups := __prism.MakePrismWithName(
		func(s __prism.NamedMatch) __option.Option[map[string]string] { return __option.Some(s.Groups) },
		func(v map[string]string) __prism.NamedMatch {
			return __prism.NamedMatch{Groups: v}
		},
		"NamedMatch.Groups",
	)
	_fromNonZeroFull := __option.FromNonZero[string]()
	_prismFull := __prism.MakePrismWithName(
		func(s __prism.NamedMatch) __option.Option[string] { return _fromNonZeroFull(s.Full) },
		func(v string) __prism.NamedMatch {
			return __prism.NamedMatch{Full: v}
		},
		"NamedMatch.Full",
	)
	_fromNonZeroAfter := __option.FromNonZero[string]()
	_prismAfter := __prism.MakePrismWithName(
		func(s __prism.NamedMatch) __option.Option[string] { return _fromNonZeroAfter(s.After) },
		func(v string) __prism.NamedMatch {
			return __prism.NamedMatch{After: v}
		},
		"NamedMatch.After",
	)
	return NamedMatchPrisms{
		Before: _prismBefore,
		Groups: _prismGroups,
		Full:   _prismFull,
		After:  _prismAfter,
	}
}
