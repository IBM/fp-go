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

package prism

import (
	"regexp"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestFromZero tests the FromZero prism with various comparable types
func TestFromZero(t *testing.T) {
	t.Run("int - match zero", func(t *testing.T) {
		prism := FromZero[int]()

		result := prism.GetOption(0)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("int - non-zero returns None", func(t *testing.T) {
		prism := FromZero[int]()

		result := prism.GetOption(42)
		assert.True(t, O.IsNone(result))
	})

	t.Run("string - match empty string", func(t *testing.T) {
		prism := FromZero[string]()

		result := prism.GetOption("")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "", O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("string - non-empty returns None", func(t *testing.T) {
		prism := FromZero[string]()

		result := prism.GetOption("hello")
		assert.True(t, O.IsNone(result))
	})

	t.Run("bool - match false", func(t *testing.T) {
		prism := FromZero[bool]()

		result := prism.GetOption(false)
		assert.True(t, O.IsSome(result))
		assert.False(t, O.GetOrElse(F.Constant(true))(result))
	})

	t.Run("bool - true returns None", func(t *testing.T) {
		prism := FromZero[bool]()

		result := prism.GetOption(true)
		assert.True(t, O.IsNone(result))
	})

	t.Run("float64 - match 0.0", func(t *testing.T) {
		prism := FromZero[float64]()

		result := prism.GetOption(0.0)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0.0, O.GetOrElse(F.Constant(-1.0))(result))
	})

	t.Run("float64 - non-zero returns None", func(t *testing.T) {
		prism := FromZero[float64]()

		result := prism.GetOption(3.14)
		assert.True(t, O.IsNone(result))
	})

	t.Run("pointer - match nil", func(t *testing.T) {
		prism := FromZero[*int]()

		var nilPtr *int
		result := prism.GetOption(nilPtr)
		assert.True(t, O.IsSome(result))
	})

	t.Run("pointer - non-nil returns None", func(t *testing.T) {
		prism := FromZero[*int]()

		value := 42
		result := prism.GetOption(&value)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get is identity", func(t *testing.T) {
		prism := FromZero[int]()

		assert.Equal(t, 0, prism.ReverseGet(0))
		assert.Equal(t, 42, prism.ReverseGet(42))
	})
}

// TestFromZeroWithSet tests using Set with FromZero prism
func TestFromZeroWithSet(t *testing.T) {
	t.Run("set on zero value", func(t *testing.T) {
		prism := FromZero[int]()

		setter := Set[int](100)
		result := setter(prism)(0)

		assert.Equal(t, 100, result)
	})

	t.Run("set on non-zero returns original", func(t *testing.T) {
		prism := FromZero[int]()

		setter := Set[int](100)
		result := setter(prism)(42)

		assert.Equal(t, 42, result)
	})
}

// TestFromZeroPrismLaws tests that FromZero satisfies prism laws
func TestFromZeroPrismLaws(t *testing.T) {
	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a) for zero", func(t *testing.T) {
		prism := FromZero[int]()

		reversed := prism.ReverseGet(0)
		extracted := prism.GetOption(reversed)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) == s", func(t *testing.T) {
		prism := FromZero[string]()

		extracted := prism.GetOption("")
		if O.IsSome(extracted) {
			value := O.GetOrElse(F.Constant("default"))(extracted)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, "", reconstructed)
		}
	})
}

// TestRegexMatcher tests the RegexMatcher prism
func TestRegexMatcher(t *testing.T) {
	t.Run("simple number match", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		prism := RegexMatcher(re)

		result := prism.GetOption("price: 42 dollars")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, "price: ", match.Before)
		assert.Equal(t, "42", match.FullMatch())
		assert.Equal(t, " dollars", match.After)
	})

	t.Run("no match returns None", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		prism := RegexMatcher(re)

		result := prism.GetOption("no numbers here")
		assert.True(t, O.IsNone(result))
	})

	t.Run("match with capture groups", func(t *testing.T) {
		re := regexp.MustCompile(`(\w+)@(\w+\.\w+)`)
		prism := RegexMatcher(re)

		result := prism.GetOption("contact: user@example.com")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, "contact: ", match.Before)
		assert.Equal(t, "user@example.com", match.FullMatch())
		assert.Equal(t, "user", match.Group(1))
		assert.Equal(t, "example.com", match.Group(2))
		assert.Equal(t, "", match.After)
	})

	t.Run("match at beginning", func(t *testing.T) {
		re := regexp.MustCompile(`^\d+`)
		prism := RegexMatcher(re)

		result := prism.GetOption("123 test")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, "", match.Before)
		assert.Equal(t, "123", match.FullMatch())
		assert.Equal(t, " test", match.After)
	})

	t.Run("match at end", func(t *testing.T) {
		re := regexp.MustCompile(`\d+$`)
		prism := RegexMatcher(re)

		result := prism.GetOption("test 123")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, "test ", match.Before)
		assert.Equal(t, "123", match.FullMatch())
		assert.Equal(t, "", match.After)
	})

	t.Run("reconstruct original string", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		prism := RegexMatcher(re)

		original := "price: 42 dollars"
		result := prism.GetOption(original)
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		reconstructed := match.Reconstruct()
		assert.Equal(t, original, reconstructed)
	})

	t.Run("reverse get reconstructs", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		prism := RegexMatcher(re)

		match := Match{
			Before: "price: ",
			Groups: []string{"42"},
			After:  " dollars",
		}

		reconstructed := prism.ReverseGet(match)
		assert.Equal(t, "price: 42 dollars", reconstructed)
	})

	t.Run("Group with invalid index returns empty", func(t *testing.T) {
		match := Match{
			Groups: []string{"full", "group1"},
		}

		assert.Equal(t, "full", match.Group(0))
		assert.Equal(t, "group1", match.Group(1))
		assert.Equal(t, "", match.Group(5))
	})

	t.Run("empty string match", func(t *testing.T) {
		re := regexp.MustCompile(`.*`)
		prism := RegexMatcher(re)

		result := prism.GetOption("")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, "", match.Before)
		assert.Equal(t, "", match.FullMatch())
		assert.Equal(t, "", match.After)
	})
}

// TestRegexMatcherPrismLaws tests that RegexMatcher satisfies prism laws
func TestRegexMatcherPrismLaws(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	prism := RegexMatcher(re)

	t.Run("law 1: GetOption(ReverseGet(match)) reconstructs", func(t *testing.T) {
		match := Match{
			Before: "test ",
			Groups: []string{"123"},
			After:  " end",
		}

		str := prism.ReverseGet(match)
		result := prism.GetOption(str)

		assert.True(t, O.IsSome(result))
		reconstructed := O.GetOrElse(F.Constant(Match{}))(result)
		assert.Equal(t, match.Before, reconstructed.Before)
		assert.Equal(t, match.Groups[0], reconstructed.Groups[0])
		assert.Equal(t, match.After, reconstructed.After)
	})

	t.Run("law 2: ReverseGet(GetOption(s)) == s for matching strings", func(t *testing.T) {
		original := "value: 42 units"
		extracted := prism.GetOption(original)

		if O.IsSome(extracted) {
			match := O.GetOrElse(F.Constant(Match{}))(extracted)
			reconstructed := prism.ReverseGet(match)
			assert.Equal(t, original, reconstructed)
		}
	})
}

// TestRegexNamedMatcher tests the RegexNamedMatcher prism
func TestRegexNamedMatcher(t *testing.T) {
	t.Run("email with named groups", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("contact: user@example.com")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, "contact: ", match.Before)
		assert.Equal(t, "user@example.com", match.Full)
		assert.Equal(t, "", match.After)
		assert.Equal(t, "user", match.Groups["user"])
		assert.Equal(t, "example.com", match.Groups["domain"])
	})

	t.Run("date with named groups", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("Date: 2024-03-15")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, "Date: ", match.Before)
		assert.Equal(t, "2024-03-15", match.Full)
		assert.Equal(t, "2024", match.Groups["year"])
		assert.Equal(t, "03", match.Groups["month"])
		assert.Equal(t, "15", match.Groups["day"])
	})

	t.Run("no match returns None", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<num>\d+)`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("no numbers")
		assert.True(t, O.IsNone(result))
	})

	t.Run("reconstruct original string", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		prism := RegexNamedMatcher(re)

		original := "email: admin@site.com here"
		result := prism.GetOption(original)
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		reconstructed := match.Reconstruct()
		assert.Equal(t, original, reconstructed)
	})

	t.Run("reverse get reconstructs", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<num>\d+)`)
		prism := RegexNamedMatcher(re)

		match := NamedMatch{
			Before: "value: ",
			Full:   "42",
			Groups: map[string]string{"num": "42"},
			After:  " end",
		}

		reconstructed := prism.ReverseGet(match)
		assert.Equal(t, "value: 42 end", reconstructed)
	})

	t.Run("unnamed groups not in map", func(t *testing.T) {
		// Mix of named and unnamed groups - use non-greedy match for clarity
		re := regexp.MustCompile(`(?P<name>[a-z]+)(\d+)`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("user123")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, "user123", match.Full)
		assert.Equal(t, "user", match.Groups["name"])
		// Only named groups should be in the map, not unnamed ones
		assert.Equal(t, 1, len(match.Groups))
	})

	t.Run("empty string match", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<all>.*)`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, "", match.Before)
		assert.Equal(t, "", match.Full)
		assert.Equal(t, "", match.After)
	})

	t.Run("multiple matches - only first", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<num>\d+)`)
		prism := RegexNamedMatcher(re)

		result := prism.GetOption("first 123 second 456")
		assert.True(t, O.IsSome(result))

		match := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, "first ", match.Before)
		assert.Equal(t, "123", match.Full)
		assert.Equal(t, " second 456", match.After)
		assert.Equal(t, "123", match.Groups["num"])
	})
}

// TestRegexNamedMatcherPrismLaws tests that RegexNamedMatcher satisfies prism laws
func TestRegexNamedMatcherPrismLaws(t *testing.T) {
	re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
	prism := RegexNamedMatcher(re)

	t.Run("law 1: GetOption(ReverseGet(match)) reconstructs", func(t *testing.T) {
		match := NamedMatch{
			Before: "email: ",
			Full:   "user@example.com",
			Groups: map[string]string{
				"user":   "user",
				"domain": "example.com",
			},
			After: "",
		}

		str := prism.ReverseGet(match)
		result := prism.GetOption(str)

		assert.True(t, O.IsSome(result))
		reconstructed := O.GetOrElse(F.Constant(NamedMatch{}))(result)
		assert.Equal(t, match.Before, reconstructed.Before)
		assert.Equal(t, match.Full, reconstructed.Full)
		assert.Equal(t, match.After, reconstructed.After)
	})

	t.Run("law 2: ReverseGet(GetOption(s)) == s for matching strings", func(t *testing.T) {
		original := "contact: admin@site.com"
		extracted := prism.GetOption(original)

		if O.IsSome(extracted) {
			match := O.GetOrElse(F.Constant(NamedMatch{}))(extracted)
			reconstructed := prism.ReverseGet(match)
			assert.Equal(t, original, reconstructed)
		}
	})
}

// TestRegexMatcherWithSet tests using Set with RegexMatcher
func TestRegexMatcherWithSet(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	prism := RegexMatcher(re)

	t.Run("set on matching string", func(t *testing.T) {
		original := "price: 42 dollars"

		newMatch := Match{
			Before: "price: ",
			Groups: []string{"100"},
			After:  " dollars",
		}

		setter := Set[string](newMatch)
		result := setter(prism)(original)

		assert.Equal(t, "price: 100 dollars", result)
	})

	t.Run("set on non-matching string returns original", func(t *testing.T) {
		original := "no numbers"

		newMatch := Match{
			Before: "",
			Groups: []string{"42"},
			After:  "",
		}

		setter := Set[string](newMatch)
		result := setter(prism)(original)

		assert.Equal(t, original, result)
	})
}

// TestRegexNamedMatcherWithSet tests using Set with RegexNamedMatcher
func TestRegexNamedMatcherWithSet(t *testing.T) {
	re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
	prism := RegexNamedMatcher(re)

	t.Run("set on matching string", func(t *testing.T) {
		original := "email: user@example.com"

		newMatch := NamedMatch{
			Before: "email: ",
			Full:   "admin@newsite.com",
			Groups: map[string]string{
				"user":   "admin",
				"domain": "newsite.com",
			},
			After: "",
		}

		setter := Set[string](newMatch)
		result := setter(prism)(original)

		assert.Equal(t, "email: admin@newsite.com", result)
	})

	t.Run("set on non-matching string returns original", func(t *testing.T) {
		original := "no email here"

		newMatch := NamedMatch{
			Before: "",
			Full:   "test@test.com",
			Groups: map[string]string{
				"user":   "test",
				"domain": "test.com",
			},
			After: "",
		}

		setter := Set[string](newMatch)
		result := setter(prism)(original)

		assert.Equal(t, original, result)
	})
}
