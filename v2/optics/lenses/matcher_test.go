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
	"testing"

	__prism "github.com/IBM/fp-go/v2/optics/prism"
	__option "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestMatchLenses_Before tests the Before lens for Match
func TestMatchLenses_Before(t *testing.T) {
	lenses := MakeMatchLenses()
	match := __prism.Match{
		Before: "prefix ",
		Groups: []string{"match", "group1"},
		After:  " suffix",
	}

	// Test Get
	before := lenses.Before.Get(match)
	assert.Equal(t, "prefix ", before)

	// Test Set (curried)
	updated := lenses.Before.Set("new prefix ")(match)
	assert.Equal(t, "new prefix ", updated.Before)
	assert.Equal(t, match.Groups, updated.Groups) // Other fields unchanged
	assert.Equal(t, match.After, updated.After)
	assert.Equal(t, "prefix ", match.Before) // Original unchanged
}

// TestMatchLenses_Groups tests the Groups lens for Match
func TestMatchLenses_Groups(t *testing.T) {
	lenses := MakeMatchLenses()
	match := __prism.Match{
		Before: "prefix ",
		Groups: []string{"match", "group1"},
		After:  " suffix",
	}

	// Test Get
	groups := lenses.Groups.Get(match)
	assert.Equal(t, []string{"match", "group1"}, groups)

	// Test Set (curried)
	newGroups := []string{"new", "groups", "here"}
	updated := lenses.Groups.Set(newGroups)(match)
	assert.Equal(t, newGroups, updated.Groups)
	assert.Equal(t, match.Before, updated.Before)
	assert.Equal(t, match.After, updated.After)
}

// TestMatchLenses_After tests the After lens for Match
func TestMatchLenses_After(t *testing.T) {
	lenses := MakeMatchLenses()
	match := __prism.Match{
		Before: "prefix ",
		Groups: []string{"match"},
		After:  " suffix",
	}

	// Test Get
	after := lenses.After.Get(match)
	assert.Equal(t, " suffix", after)

	// Test Set (curried)
	updated := lenses.After.Set(" new suffix")(match)
	assert.Equal(t, " new suffix", updated.After)
	assert.Equal(t, match.Before, updated.Before)
	assert.Equal(t, match.Groups, updated.Groups)
}

// TestMatchLenses_BeforeO tests the optional Before lens
func TestMatchLenses_BeforeO(t *testing.T) {
	lenses := MakeMatchLenses()

	t.Run("non-empty Before", func(t *testing.T) {
		match := __prism.Match{Before: "prefix ", Groups: []string{}, After: ""}
		opt := lenses.BeforeO.Get(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "prefix ", value)
	})

	t.Run("empty Before", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: ""}
		opt := lenses.BeforeO.Get(match)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("set Some", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: ""}
		updated := lenses.BeforeO.Set(__option.Some("new "))(match)
		assert.Equal(t, "new ", updated.Before)
	})

	t.Run("set None", func(t *testing.T) {
		match := __prism.Match{Before: "prefix ", Groups: []string{}, After: ""}
		updated := lenses.BeforeO.Set(__option.None[string]())(match)
		assert.Equal(t, "", updated.Before)
	})
}

// TestMatchLenses_AfterO tests the optional After lens
func TestMatchLenses_AfterO(t *testing.T) {
	lenses := MakeMatchLenses()

	t.Run("non-empty After", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: " suffix"}
		opt := lenses.AfterO.Get(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, " suffix", value)
	})

	t.Run("empty After", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: ""}
		opt := lenses.AfterO.Get(match)
		assert.True(t, __option.IsNone(opt))
	})
}

// TestMatchRefLenses_Before tests the Before lens for *Match
func TestMatchRefLenses_Before(t *testing.T) {
	lenses := MakeMatchRefLenses()
	match := &__prism.Match{
		Before: "prefix ",
		Groups: []string{"match"},
		After:  " suffix",
	}

	// Test Get
	before := lenses.Before.Get(match)
	assert.Equal(t, "prefix ", before)

	// Test Set (creates copy with MakeLensStrictWithName, curried)
	result := lenses.Before.Set("new prefix ")(match)
	assert.Equal(t, "new prefix ", result.Before)
	assert.Equal(t, "prefix ", match.Before) // Original unchanged
	assert.NotSame(t, match, result)         // Returns new pointer
}

// TestMatchRefLenses_Groups tests the Groups lens for *Match
func TestMatchRefLenses_Groups(t *testing.T) {
	lenses := MakeMatchRefLenses()
	match := &__prism.Match{
		Before: "prefix ",
		Groups: []string{"match", "group1"},
		After:  " suffix",
	}

	// Test Get
	groups := lenses.Groups.Get(match)
	assert.Equal(t, []string{"match", "group1"}, groups)

	// Test Set (creates copy with MakeLensRefWithName, curried)
	newGroups := []string{"new", "groups"}
	result := lenses.Groups.Set(newGroups)(match)
	assert.Equal(t, newGroups, result.Groups)
	assert.Equal(t, []string{"match", "group1"}, match.Groups) // Original unchanged
	assert.NotSame(t, match, result)
}

// TestMatchRefLenses_After tests the After lens for *Match
func TestMatchRefLenses_After(t *testing.T) {
	lenses := MakeMatchRefLenses()
	match := &__prism.Match{
		Before: "prefix ",
		Groups: []string{"match"},
		After:  " suffix",
	}

	// Test Get
	after := lenses.After.Get(match)
	assert.Equal(t, " suffix", after)

	// Test Set (creates copy with MakeLensStrictWithName, curried)
	result := lenses.After.Set(" new suffix")(match)
	assert.Equal(t, " new suffix", result.After)
	assert.Equal(t, " suffix", match.After) // Original unchanged
	assert.NotSame(t, match, result)
}

// TestMatchPrisms_Before tests the Before prism
func TestMatchPrisms_Before(t *testing.T) {
	prisms := MakeMatchPrisms()

	t.Run("GetOption with non-empty Before", func(t *testing.T) {
		match := __prism.Match{Before: "prefix ", Groups: []string{}, After: ""}
		opt := prisms.Before.GetOption(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "prefix ", value)
	})

	t.Run("GetOption with empty Before", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: ""}
		opt := prisms.Before.GetOption(match)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("ReverseGet", func(t *testing.T) {
		match := prisms.Before.ReverseGet("test ")
		assert.Equal(t, "test ", match.Before)
		assert.Nil(t, match.Groups)
		assert.Equal(t, "", match.After)
	})
}

// TestMatchPrisms_Groups tests the Groups prism
func TestMatchPrisms_Groups(t *testing.T) {
	prisms := MakeMatchPrisms()

	t.Run("GetOption always returns Some", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{"a", "b"}, After: ""}
		opt := prisms.Groups.GetOption(match)
		assert.True(t, __option.IsSome(opt))
		groups := __option.GetOrElse(func() []string { return nil })(opt)
		assert.Equal(t, []string{"a", "b"}, groups)
	})

	t.Run("GetOption with nil Groups", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: nil, After: ""}
		opt := prisms.Groups.GetOption(match)
		assert.True(t, __option.IsSome(opt))
	})

	t.Run("ReverseGet", func(t *testing.T) {
		groups := []string{"test", "groups"}
		match := prisms.Groups.ReverseGet(groups)
		assert.Equal(t, groups, match.Groups)
		assert.Equal(t, "", match.Before)
		assert.Equal(t, "", match.After)
	})
}

// TestMatchPrisms_After tests the After prism
func TestMatchPrisms_After(t *testing.T) {
	prisms := MakeMatchPrisms()

	t.Run("GetOption with non-empty After", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: " suffix"}
		opt := prisms.After.GetOption(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, " suffix", value)
	})

	t.Run("GetOption with empty After", func(t *testing.T) {
		match := __prism.Match{Before: "", Groups: []string{}, After: ""}
		opt := prisms.After.GetOption(match)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("ReverseGet", func(t *testing.T) {
		match := prisms.After.ReverseGet(" test")
		assert.Equal(t, " test", match.After)
		assert.Nil(t, match.Groups)
		assert.Equal(t, "", match.Before)
	})
}

// TestNamedMatchLenses_Before tests the Before lens for NamedMatch
func TestNamedMatchLenses_Before(t *testing.T) {
	lenses := MakeNamedMatchLenses()
	match := __prism.NamedMatch{
		Before: "Email: ",
		Groups: map[string]string{"user": "john", "domain": "example.com"},
		Full:   "john@example.com",
		After:  "",
	}

	// Test Get
	before := lenses.Before.Get(match)
	assert.Equal(t, "Email: ", before)

	// Test Set (curried)
	updated := lenses.Before.Set("Contact: ")(match)
	assert.Equal(t, "Contact: ", updated.Before)
	assert.Equal(t, match.Groups, updated.Groups)
	assert.Equal(t, match.Full, updated.Full)
	assert.Equal(t, match.After, updated.After)
	assert.Equal(t, "Email: ", match.Before) // Original unchanged
}

// TestNamedMatchLenses_Groups tests the Groups lens for NamedMatch
func TestNamedMatchLenses_Groups(t *testing.T) {
	lenses := MakeNamedMatchLenses()
	match := __prism.NamedMatch{
		Before: "Email: ",
		Groups: map[string]string{"user": "john", "domain": "example.com"},
		Full:   "john@example.com",
		After:  "",
	}

	// Test Get
	groups := lenses.Groups.Get(match)
	assert.Equal(t, map[string]string{"user": "john", "domain": "example.com"}, groups)

	// Test Set (curried)
	newGroups := map[string]string{"user": "alice", "domain": "test.org"}
	updated := lenses.Groups.Set(newGroups)(match)
	assert.Equal(t, newGroups, updated.Groups)
	assert.Equal(t, match.Before, updated.Before)
}

// TestNamedMatchLenses_Full tests the Full lens for NamedMatch
func TestNamedMatchLenses_Full(t *testing.T) {
	lenses := MakeNamedMatchLenses()
	match := __prism.NamedMatch{
		Before: "Email: ",
		Groups: map[string]string{"user": "john"},
		Full:   "john@example.com",
		After:  "",
	}

	// Test Get
	full := lenses.Full.Get(match)
	assert.Equal(t, "john@example.com", full)

	// Test Set (curried)
	updated := lenses.Full.Set("alice@test.org")(match)
	assert.Equal(t, "alice@test.org", updated.Full)
	assert.Equal(t, match.Before, updated.Before)
}

// TestNamedMatchLenses_After tests the After lens for NamedMatch
func TestNamedMatchLenses_After(t *testing.T) {
	lenses := MakeNamedMatchLenses()
	match := __prism.NamedMatch{
		Before: "Email: ",
		Groups: map[string]string{"user": "john"},
		Full:   "john@example.com",
		After:  " for contact",
	}

	// Test Get
	after := lenses.After.Get(match)
	assert.Equal(t, " for contact", after)

	// Test Set (curried)
	updated := lenses.After.Set(" for info")(match)
	assert.Equal(t, " for info", updated.After)
	assert.Equal(t, match.Before, updated.Before)
}

// TestNamedMatchLenses_Optional tests optional lenses for NamedMatch
func TestNamedMatchLenses_Optional(t *testing.T) {
	lenses := MakeNamedMatchLenses()

	t.Run("BeforeO with non-empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "prefix ", Groups: nil, Full: "", After: ""}
		opt := lenses.BeforeO.Get(match)
		assert.True(t, __option.IsSome(opt))
	})

	t.Run("BeforeO with empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt := lenses.BeforeO.Get(match)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("FullO with non-empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "test@example.com", After: ""}
		opt := lenses.FullO.Get(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "test@example.com", value)
	})

	t.Run("FullO with empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt := lenses.FullO.Get(match)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("AfterO with non-empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: " suffix"}
		opt := lenses.AfterO.Get(match)
		assert.True(t, __option.IsSome(opt))
	})

	t.Run("AfterO with empty", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt := lenses.AfterO.Get(match)
		assert.True(t, __option.IsNone(opt))
	})
}

// TestNamedMatchRefLenses_Immutability tests that reference lenses create copies
func TestNamedMatchRefLenses_Immutability(t *testing.T) {
	lenses := MakeNamedMatchRefLenses()
	match := &__prism.NamedMatch{
		Before: "Email: ",
		Groups: map[string]string{"user": "john"},
		Full:   "john@example.com",
		After:  "",
	}

	// Test Before (creates copy, curried)
	updated1 := lenses.Before.Set("Contact: ")(match)
	assert.Equal(t, "Contact: ", updated1.Before)
	assert.Equal(t, "Email: ", match.Before) // Original unchanged

	// Test Groups (creates copy, curried)
	newGroups := map[string]string{"user": "alice"}
	updated2 := lenses.Groups.Set(newGroups)(match)
	assert.Equal(t, newGroups, updated2.Groups)
	assert.Equal(t, map[string]string{"user": "john"}, match.Groups) // Original unchanged

	// Test Full (creates copy, curried)
	updated3 := lenses.Full.Set("alice@test.org")(match)
	assert.Equal(t, "alice@test.org", updated3.Full)
	assert.Equal(t, "john@example.com", match.Full) // Original unchanged

	// Test After (creates copy, curried)
	updated4 := lenses.After.Set(" for info")(match)
	assert.Equal(t, " for info", updated4.After)
	assert.Equal(t, "", match.After) // Original unchanged
}

// TestNamedMatchPrisms tests prisms for NamedMatch
func TestNamedMatchPrisms(t *testing.T) {
	prisms := MakeNamedMatchPrisms()

	t.Run("Before prism", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "prefix ", Groups: nil, Full: "", After: ""}
		opt := prisms.Before.GetOption(match)
		assert.True(t, __option.IsSome(opt))

		emptyMatch := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt = prisms.Before.GetOption(emptyMatch)
		assert.True(t, __option.IsNone(opt))

		constructed := prisms.Before.ReverseGet("test ")
		assert.Equal(t, "test ", constructed.Before)
	})

	t.Run("Groups prism always returns Some", func(t *testing.T) {
		match := __prism.NamedMatch{
			Before: "",
			Groups: map[string]string{"key": "value"},
			Full:   "",
			After:  "",
		}
		opt := prisms.Groups.GetOption(match)
		assert.True(t, __option.IsSome(opt))

		nilMatch := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt = prisms.Groups.GetOption(nilMatch)
		assert.True(t, __option.IsSome(opt))
	})

	t.Run("Full prism", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "test@example.com", After: ""}
		opt := prisms.Full.GetOption(match)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "test@example.com", value)

		emptyMatch := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt = prisms.Full.GetOption(emptyMatch)
		assert.True(t, __option.IsNone(opt))

		constructed := prisms.Full.ReverseGet("alice@test.org")
		assert.Equal(t, "alice@test.org", constructed.Full)
	})

	t.Run("After prism", func(t *testing.T) {
		match := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: " suffix"}
		opt := prisms.After.GetOption(match)
		assert.True(t, __option.IsSome(opt))

		emptyMatch := __prism.NamedMatch{Before: "", Groups: nil, Full: "", After: ""}
		opt = prisms.After.GetOption(emptyMatch)
		assert.True(t, __option.IsNone(opt))

		constructed := prisms.After.ReverseGet(" test")
		assert.Equal(t, " test", constructed.After)
	})
}

// TestMatchLenses_Immutability verifies that value lenses don't mutate originals
func TestMatchLenses_Immutability(t *testing.T) {
	lenses := MakeMatchLenses()
	original := __prism.Match{
		Before: "original ",
		Groups: []string{"group1", "group2"},
		After:  " original",
	}

	// Make a copy to compare later
	originalBefore := original.Before
	originalGroups := make([]string, len(original.Groups))
	copy(originalGroups, original.Groups)
	originalAfter := original.After

	// Perform multiple updates (curried)
	updated1 := lenses.Before.Set("updated ")(original)
	updated2 := lenses.Groups.Set([]string{"new"})(updated1)
	updated3 := lenses.After.Set(" updated")(updated2)

	// Verify original is unchanged
	assert.Equal(t, originalBefore, original.Before)
	assert.Equal(t, originalGroups, original.Groups)
	assert.Equal(t, originalAfter, original.After)

	// Verify updates worked
	assert.Equal(t, "updated ", updated3.Before)
	assert.Equal(t, []string{"new"}, updated3.Groups)
	assert.Equal(t, " updated", updated3.After)
}

// TestNamedMatchLenses_Immutability verifies that value lenses don't mutate originals
func TestNamedMatchLenses_Immutability(t *testing.T) {
	lenses := MakeNamedMatchLenses()
	original := __prism.NamedMatch{
		Before: "original ",
		Groups: map[string]string{"key": "value"},
		Full:   "original@example.com",
		After:  " original",
	}

	// Make copies to compare later
	originalBefore := original.Before
	originalFull := original.Full
	originalAfter := original.After

	// Perform multiple updates (curried)
	updated1 := lenses.Before.Set("updated ")(original)
	updated2 := lenses.Full.Set("updated@test.org")(updated1)
	updated3 := lenses.After.Set(" updated")(updated2)

	// Verify original is unchanged
	assert.Equal(t, originalBefore, original.Before)
	assert.Equal(t, originalFull, original.Full)
	assert.Equal(t, originalAfter, original.After)

	// Verify updates worked
	assert.Equal(t, "updated ", updated3.Before)
	assert.Equal(t, "updated@test.org", updated3.Full)
	assert.Equal(t, " updated", updated3.After)
}
