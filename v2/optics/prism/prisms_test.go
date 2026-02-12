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
	"errors"
	"regexp"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
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

// TestFromNonZero tests the FromNonZero prism with various comparable types
func TestFromNonZero(t *testing.T) {
	t.Run("int - match non-zero", func(t *testing.T) {
		prism := FromNonZero[int]()

		result := prism.GetOption(42)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("int - zero returns None", func(t *testing.T) {
		prism := FromNonZero[int]()

		result := prism.GetOption(0)
		assert.True(t, O.IsNone(result))
	})

	t.Run("string - match non-empty string", func(t *testing.T) {
		prism := FromNonZero[string]()

		result := prism.GetOption("hello")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("string - empty returns None", func(t *testing.T) {
		prism := FromNonZero[string]()

		result := prism.GetOption("")
		assert.True(t, O.IsNone(result))
	})

	t.Run("bool - match true", func(t *testing.T) {
		prism := FromNonZero[bool]()

		result := prism.GetOption(true)
		assert.True(t, O.IsSome(result))
		assert.True(t, O.GetOrElse(F.Constant(false))(result))
	})

	t.Run("bool - false returns None", func(t *testing.T) {
		prism := FromNonZero[bool]()

		result := prism.GetOption(false)
		assert.True(t, O.IsNone(result))
	})

	t.Run("float64 - match non-zero", func(t *testing.T) {
		prism := FromNonZero[float64]()

		result := prism.GetOption(3.14)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 3.14, O.GetOrElse(F.Constant(-1.0))(result))
	})

	t.Run("float64 - zero returns None", func(t *testing.T) {
		prism := FromNonZero[float64]()

		result := prism.GetOption(0.0)
		assert.True(t, O.IsNone(result))
	})

	t.Run("pointer - match non-nil", func(t *testing.T) {
		prism := FromNonZero[*int]()

		value := 42
		result := prism.GetOption(&value)
		assert.True(t, O.IsSome(result))
	})

	t.Run("pointer - nil returns None", func(t *testing.T) {
		prism := FromNonZero[*int]()

		var nilPtr *int
		result := prism.GetOption(nilPtr)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get is identity", func(t *testing.T) {
		prism := FromNonZero[int]()

		assert.Equal(t, 0, prism.ReverseGet(0))
		assert.Equal(t, 42, prism.ReverseGet(42))
	})
}

// TestFromNonZeroWithSet tests using Set with FromNonZero prism
func TestFromNonZeroWithSet(t *testing.T) {
	t.Run("set on non-zero value", func(t *testing.T) {
		prism := FromNonZero[int]()

		setter := Set[int](100)
		result := setter(prism)(42)

		assert.Equal(t, 100, result)
	})

	t.Run("set on zero returns original", func(t *testing.T) {
		prism := FromNonZero[int]()

		setter := Set[int](100)
		result := setter(prism)(0)

		assert.Equal(t, 0, result)
	})
}

// TestParseInt tests the ParseInt prism
func TestParseInt(t *testing.T) {
	prism := ParseInt()

	t.Run("parse valid positive integer", func(t *testing.T) {
		result := prism.GetOption("42")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("parse valid negative integer", func(t *testing.T) {
		result := prism.GetOption("-123")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, -123, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("parse zero", func(t *testing.T) {
		result := prism.GetOption("0")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("parse invalid integer", func(t *testing.T) {
		result := prism.GetOption("not-a-number")
		assert.True(t, O.IsNone(result))
	})

	t.Run("parse float as integer fails", func(t *testing.T) {
		result := prism.GetOption("3.14")
		assert.True(t, O.IsNone(result))
	})

	t.Run("parse empty string fails", func(t *testing.T) {
		result := prism.GetOption("")
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get formats integer", func(t *testing.T) {
		assert.Equal(t, "42", prism.ReverseGet(42))
		assert.Equal(t, "-123", prism.ReverseGet(-123))
		assert.Equal(t, "0", prism.ReverseGet(0))
	})

	t.Run("round trip", func(t *testing.T) {
		original := "12345"
		result := prism.GetOption(original)
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(0))(result)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, original, reconstructed)
		}
	})
}

// TestParseInt64 tests the ParseInt64 prism
func TestParseInt64(t *testing.T) {
	prism := ParseInt64()

	t.Run("parse valid int64", func(t *testing.T) {
		result := prism.GetOption("9223372036854775807")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, int64(9223372036854775807), O.GetOrElse(F.Constant(int64(-1)))(result))
	})

	t.Run("parse negative int64", func(t *testing.T) {
		result := prism.GetOption("-9223372036854775808")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, int64(-9223372036854775808), O.GetOrElse(F.Constant(int64(0)))(result))
	})

	t.Run("parse invalid int64", func(t *testing.T) {
		result := prism.GetOption("not-a-number")
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get formats int64", func(t *testing.T) {
		assert.Equal(t, "42", prism.ReverseGet(int64(42)))
		assert.Equal(t, "9223372036854775807", prism.ReverseGet(int64(9223372036854775807)))
	})

	t.Run("round trip", func(t *testing.T) {
		original := "1234567890123456789"
		result := prism.GetOption(original)
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(int64(0)))(result)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, original, reconstructed)
		}
	})
}

// TestParseBool tests the ParseBool prism
func TestParseBool(t *testing.T) {
	prism := ParseBool()

	t.Run("parse true variations", func(t *testing.T) {
		trueValues := []string{"true", "True", "TRUE", "t", "T", "1"}
		for _, val := range trueValues {
			result := prism.GetOption(val)
			assert.True(t, O.IsSome(result), "Should parse: %s", val)
			assert.True(t, O.GetOrElse(F.Constant(false))(result), "Should be true: %s", val)
		}
	})

	t.Run("parse false variations", func(t *testing.T) {
		falseValues := []string{"false", "False", "FALSE", "f", "F", "0"}
		for _, val := range falseValues {
			result := prism.GetOption(val)
			assert.True(t, O.IsSome(result), "Should parse: %s", val)
			assert.False(t, O.GetOrElse(F.Constant(true))(result), "Should be false: %s", val)
		}
	})

	t.Run("parse invalid bool", func(t *testing.T) {
		invalidValues := []string{"maybe", "yes", "no", "2", ""}
		for _, val := range invalidValues {
			result := prism.GetOption(val)
			assert.True(t, O.IsNone(result), "Should not parse: %s", val)
		}
	})

	t.Run("reverse get formats bool", func(t *testing.T) {
		assert.Equal(t, "true", prism.ReverseGet(true))
		assert.Equal(t, "false", prism.ReverseGet(false))
	})

	t.Run("round trip with true", func(t *testing.T) {
		result := prism.GetOption("true")
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(false))(result)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, "true", reconstructed)
		}
	})

	t.Run("round trip with false", func(t *testing.T) {
		result := prism.GetOption("false")
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(true))(result)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, "false", reconstructed)
		}
	})
}

// TestParseFloat32 tests the ParseFloat32 prism
func TestParseFloat32(t *testing.T) {
	prism := ParseFloat32()

	t.Run("parse valid float32", func(t *testing.T) {
		result := prism.GetOption("3.14")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(float32(0)))(result)
		assert.InDelta(t, float32(3.14), value, 0.0001)
	})

	t.Run("parse negative float32", func(t *testing.T) {
		result := prism.GetOption("-2.71")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(float32(0)))(result)
		assert.InDelta(t, float32(-2.71), value, 0.0001)
	})

	t.Run("parse scientific notation", func(t *testing.T) {
		result := prism.GetOption("1.5e10")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(float32(0)))(result)
		assert.InDelta(t, float32(1.5e10), value, 1e6)
	})

	t.Run("parse integer as float", func(t *testing.T) {
		result := prism.GetOption("42")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(float32(0)))(result)
		assert.Equal(t, float32(42), value)
	})

	t.Run("parse invalid float", func(t *testing.T) {
		result := prism.GetOption("not-a-number")
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get formats float32", func(t *testing.T) {
		str := prism.ReverseGet(float32(3.14))
		assert.Contains(t, str, "3.14")
	})

	t.Run("round trip", func(t *testing.T) {
		original := "3.14159"
		result := prism.GetOption(original)
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(float32(0)))(result)
			reconstructed := prism.ReverseGet(value)
			// Parse both to compare as floats due to precision
			origFloat := F.Pipe1(original, prism.GetOption)
			reconFloat := F.Pipe1(reconstructed, prism.GetOption)
			if O.IsSome(origFloat) && O.IsSome(reconFloat) {
				assert.InDelta(t,
					O.GetOrElse(F.Constant(float32(0)))(origFloat),
					O.GetOrElse(F.Constant(float32(0)))(reconFloat),
					0.0001)
			}
		}
	})
}

// TestParseFloat64 tests the ParseFloat64 prism
func TestParseFloat64(t *testing.T) {
	prism := ParseFloat64()

	t.Run("parse valid float64", func(t *testing.T) {
		result := prism.GetOption("3.141592653589793")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(0.0))(result)
		assert.InDelta(t, 3.141592653589793, value, 1e-15)
	})

	t.Run("parse negative float64", func(t *testing.T) {
		result := prism.GetOption("-2.718281828459045")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(0.0))(result)
		assert.InDelta(t, -2.718281828459045, value, 1e-15)
	})

	t.Run("parse scientific notation", func(t *testing.T) {
		result := prism.GetOption("1.5e100")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(0.0))(result)
		assert.InDelta(t, 1.5e100, value, 1e85)
	})

	t.Run("parse integer as float", func(t *testing.T) {
		result := prism.GetOption("42")
		assert.True(t, O.IsSome(result))
		value := O.GetOrElse(F.Constant(0.0))(result)
		assert.Equal(t, 42.0, value)
	})

	t.Run("parse invalid float", func(t *testing.T) {
		result := prism.GetOption("not-a-number")
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get formats float64", func(t *testing.T) {
		str := prism.ReverseGet(3.141592653589793)
		assert.Contains(t, str, "3.14159")
	})

	t.Run("round trip", func(t *testing.T) {
		original := "3.141592653589793"
		result := prism.GetOption(original)
		if O.IsSome(result) {
			value := O.GetOrElse(F.Constant(0.0))(result)
			reconstructed := prism.ReverseGet(value)
			// Parse both to compare as floats
			origFloat := F.Pipe1(original, prism.GetOption)
			reconFloat := F.Pipe1(reconstructed, prism.GetOption)
			if O.IsSome(origFloat) && O.IsSome(reconFloat) {
				assert.InDelta(t,
					O.GetOrElse(F.Constant(0.0))(origFloat),
					O.GetOrElse(F.Constant(0.0))(reconFloat),
					1e-15)
			}
		}
	})
}

// TestParseIntWithSet tests using Set with ParseInt prism
func TestParseIntWithSet(t *testing.T) {
	prism := ParseInt()

	t.Run("set on valid integer string", func(t *testing.T) {
		setter := Set[string](100)
		result := setter(prism)("42")
		assert.Equal(t, "100", result)
	})

	t.Run("set on invalid string returns original", func(t *testing.T) {
		setter := Set[string](100)
		result := setter(prism)("not-a-number")
		assert.Equal(t, "not-a-number", result)
	})
}

// TestParseBoolWithSet tests using Set with ParseBool prism
func TestParseBoolWithSet(t *testing.T) {
	prism := ParseBool()

	t.Run("set on valid bool string", func(t *testing.T) {
		setter := Set[string](true)
		result := setter(prism)("false")
		assert.Equal(t, "true", result)
	})

	t.Run("set on invalid string returns original", func(t *testing.T) {
		setter := Set[string](true)
		result := setter(prism)("maybe")
		assert.Equal(t, "maybe", result)
	})
}

// TestFromOption tests the FromOption prism
func TestFromOption(t *testing.T) {
	t.Run("extract from Some", func(t *testing.T) {
		prism := FromOption[int]()

		someValue := O.Some(42)
		result := prism.GetOption(someValue)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("extract from None returns None", func(t *testing.T) {
		prism := FromOption[int]()

		noneValue := O.None[int]()
		result := prism.GetOption(noneValue)
		assert.True(t, O.IsNone(result))
	})

	t.Run("reverse get wraps in Some", func(t *testing.T) {
		prism := FromOption[int]()

		wrapped := prism.ReverseGet(100)
		assert.True(t, O.IsSome(wrapped))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(-1))(wrapped))
	})

	t.Run("works with string type", func(t *testing.T) {
		prism := FromOption[string]()

		someStr := O.Some("hello")
		result := prism.GetOption(someStr)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("works with pointer type", func(t *testing.T) {
		prism := FromOption[*int]()

		value := 42
		ptr := &value
		somePtr := O.Some(ptr)
		result := prism.GetOption(somePtr)
		assert.True(t, O.IsSome(result))
		extractedPtr := O.GetOrElse(F.Constant[*int](nil))(result)
		assert.NotNil(t, extractedPtr)
		assert.Equal(t, 42, *extractedPtr)
	})

	t.Run("works with struct type", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		prism := FromOption[Person]()

		person := Person{Name: "Alice", Age: 30}
		somePerson := O.Some(person)
		result := prism.GetOption(somePerson)
		assert.True(t, O.IsSome(result))
		extracted := O.GetOrElse(F.Constant(Person{}))(result)
		assert.Equal(t, "Alice", extracted.Name)
		assert.Equal(t, 30, extracted.Age)
	})

	t.Run("round trip with Some", func(t *testing.T) {
		prism := FromOption[int]()

		original := O.Some(42)
		extracted := prism.GetOption(original)
		if O.IsSome(extracted) {
			value := O.GetOrElse(F.Constant(0))(extracted)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, original, reconstructed)
		}
	})
}

// TestFromOptionWithSet tests using Set with FromOption prism
func TestFromOptionWithSet(t *testing.T) {
	t.Run("set on Some value", func(t *testing.T) {
		prism := FromOption[int]()

		someValue := O.Some(42)
		setter := Set[Option[int]](200)
		result := setter(prism)(someValue)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 200, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("set on None returns None", func(t *testing.T) {
		prism := FromOption[int]()

		noneValue := O.None[int]()
		setter := Set[Option[int]](200)
		result := setter(prism)(noneValue)

		assert.True(t, O.IsNone(result))
	})

	t.Run("set with string type", func(t *testing.T) {
		prism := FromOption[string]()

		someStr := O.Some("hello")
		setter := Set[Option[string]]("world")
		result := setter(prism)(someStr)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "world", O.GetOrElse(F.Constant(""))(result))
	})
}

// TestFromOptionPrismLaws tests that FromOption satisfies prism laws
func TestFromOptionPrismLaws(t *testing.T) {
	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		prism := FromOption[int]()

		value := 42
		wrapped := prism.ReverseGet(value)
		extracted := prism.GetOption(wrapped)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, value, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) == s for Some", func(t *testing.T) {
		prism := FromOption[string]()

		original := O.Some("test")
		extracted := prism.GetOption(original)

		if O.IsSome(extracted) {
			value := O.GetOrElse(F.Constant(""))(extracted)
			reconstructed := prism.ReverseGet(value)
			assert.Equal(t, original, reconstructed)
		}
	})

	t.Run("GetOption is identity for Option type", func(t *testing.T) {
		prism := FromOption[int]()

		someValue := O.Some(42)
		result := prism.GetOption(someValue)
		assert.Equal(t, someValue, result)

		noneValue := O.None[int]()
		result = prism.GetOption(noneValue)
		assert.Equal(t, noneValue, result)
	})
}

// TestFromOptionComposition tests composing FromOption with other prisms
func TestFromOptionComposition(t *testing.T) {
	t.Run("compose with ParseInt", func(t *testing.T) {
		// Create a prism that extracts int from Option[string] by parsing
		stringPrism := FromOption[string]()
		intPrism := ParseInt()

		// Compose: Option[string] -> string -> int
		composed := Compose[Option[string]](intPrism)(stringPrism)

		// Test with Some("42")
		someStr := O.Some("42")
		result := composed.GetOption(someStr)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))

		// Test with Some("invalid")
		invalidStr := O.Some("invalid")
		result = composed.GetOption(invalidStr)
		assert.True(t, O.IsNone(result))

		// Test with None
		noneStr := O.None[string]()
		result = composed.GetOption(noneStr)
		assert.True(t, O.IsNone(result))
	})

	t.Run("nested Options", func(t *testing.T) {
		// Extract int from Option[Option[int]]
		outerPrism := FromOption[Option[int]]()
		innerPrism := FromOption[int]()

		composed := Compose[Option[Option[int]]](innerPrism)(outerPrism)

		// Test with Some(Some(42))
		nested := O.Some(O.Some(42))
		result := composed.GetOption(nested)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))

		// Test with Some(None)
		someNone := O.Some(O.None[int]())
		result = composed.GetOption(someNone)
		assert.True(t, O.IsNone(result))

		// Test with None
		none := O.None[Option[int]]()
		result = composed.GetOption(none)
		assert.True(t, O.IsNone(result))
	})
}

// TestNonEmptyString tests the NonEmptyString prism
func TestNonEmptyString(t *testing.T) {
	t.Run("match non-empty string", func(t *testing.T) {
		prism := NonEmptyString()

		result := prism.GetOption("hello")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("empty string returns None", func(t *testing.T) {
		prism := NonEmptyString()

		result := prism.GetOption("")
		assert.True(t, O.IsNone(result))
	})

	t.Run("whitespace string is non-empty", func(t *testing.T) {
		prism := NonEmptyString()

		result := prism.GetOption("   ")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "   ", O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("single character string", func(t *testing.T) {
		prism := NonEmptyString()

		result := prism.GetOption("a")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "a", O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("multiline string", func(t *testing.T) {
		prism := NonEmptyString()

		multiline := "line1\nline2\nline3"
		result := prism.GetOption(multiline)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, multiline, O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("unicode string", func(t *testing.T) {
		prism := NonEmptyString()

		unicode := "Hello 世界 🌍"
		result := prism.GetOption(unicode)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, unicode, O.GetOrElse(F.Constant("default"))(result))
	})

	t.Run("reverse get is identity", func(t *testing.T) {
		prism := NonEmptyString()

		assert.Equal(t, "", prism.ReverseGet(""))
		assert.Equal(t, "hello", prism.ReverseGet("hello"))
		assert.Equal(t, "world", prism.ReverseGet("world"))
	})
}

// TestNonEmptyStringWithSet tests using Set with NonEmptyString prism
func TestNonEmptyStringWithSet(t *testing.T) {
	t.Run("set on non-empty string", func(t *testing.T) {
		prism := NonEmptyString()

		setter := Set[string]("updated")
		result := setter(prism)("original")

		assert.Equal(t, "updated", result)
	})

	t.Run("set on empty string returns original", func(t *testing.T) {
		prism := NonEmptyString()

		setter := Set[string]("updated")
		result := setter(prism)("")

		assert.Equal(t, "", result)
	})

	t.Run("set with empty value on non-empty string", func(t *testing.T) {
		prism := NonEmptyString()

		setter := Set[string]("")
		result := setter(prism)("original")

		assert.Equal(t, "", result)
	})
}

// TestNonEmptyStringPrismLaws tests that NonEmptyString satisfies prism laws
func TestNonEmptyStringPrismLaws(t *testing.T) {
	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		prism := NonEmptyString()

		// For any non-empty string a, GetOption(ReverseGet(a)) should return Some(a)
		testCases := []string{"hello", "world", "a", "test string", "123"}
		for _, testCase := range testCases {
			reversed := prism.ReverseGet(testCase)
			result := prism.GetOption(reversed)

			assert.True(t, O.IsSome(result), "Expected Some for: %s", testCase)
			assert.Equal(t, testCase, O.GetOrElse(F.Constant(""))(result))
		}
	})

	t.Run("law 2: if GetOption(s) == Some(a), then ReverseGet(a) == s", func(t *testing.T) {
		prism := NonEmptyString()

		// For any non-empty string s where GetOption(s) returns Some(a),
		// ReverseGet(a) should equal s
		testCases := []string{"hello", "world", "test", "   ", "123"}
		for _, testCase := range testCases {
			optResult := prism.GetOption(testCase)
			if O.IsSome(optResult) {
				extracted := O.GetOrElse(F.Constant(""))(optResult)
				reversed := prism.ReverseGet(extracted)
				assert.Equal(t, testCase, reversed)
			}
		}
	})

	t.Run("law 3: GetOption is idempotent", func(t *testing.T) {
		prism := NonEmptyString()

		testCases := []string{"hello", "", "world", "   "}
		for _, testCase := range testCases {
			result1 := prism.GetOption(testCase)
			result2 := prism.GetOption(testCase)

			assert.Equal(t, result1, result2, "GetOption should be idempotent for: %s", testCase)
		}
	})
}

// TestNonEmptyStringComposition tests composing NonEmptyString with other prisms
func TestNonEmptyStringComposition(t *testing.T) {
	t.Run("compose with ParseInt", func(t *testing.T) {
		// Create a prism that only parses non-empty strings to int
		nonEmptyPrism := NonEmptyString()
		intPrism := ParseInt()

		// Compose: string -> non-empty string -> int
		composed := Compose[string](intPrism)(nonEmptyPrism)

		// Test with valid non-empty string
		result := composed.GetOption("42")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(result))

		// Test with empty string
		result = composed.GetOption("")
		assert.True(t, O.IsNone(result))

		// Test with invalid non-empty string
		result = composed.GetOption("abc")
		assert.True(t, O.IsNone(result))
	})

	t.Run("compose with ParseFloat64", func(t *testing.T) {
		// Create a prism that only parses non-empty strings to float64
		nonEmptyPrism := NonEmptyString()
		floatPrism := ParseFloat64()

		composed := Compose[string](floatPrism)(nonEmptyPrism)

		// Test with valid non-empty string
		result := composed.GetOption("3.14")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 3.14, O.GetOrElse(F.Constant(-1.0))(result))

		// Test with empty string
		result = composed.GetOption("")
		assert.True(t, O.IsNone(result))

		// Test with invalid non-empty string
		result = composed.GetOption("not a number")
		assert.True(t, O.IsNone(result))
	})

	t.Run("compose with FromOption", func(t *testing.T) {
		// Create a prism that extracts non-empty strings from Option[string]
		optionPrism := FromOption[string]()
		nonEmptyPrism := NonEmptyString()

		composed := Compose[Option[string]](nonEmptyPrism)(optionPrism)

		// Test with Some(non-empty)
		someNonEmpty := O.Some("hello")
		result := composed.GetOption(someNonEmpty)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(result))

		// Test with Some(empty)
		someEmpty := O.Some("")
		result = composed.GetOption(someEmpty)
		assert.True(t, O.IsNone(result))

		// Test with None
		none := O.None[string]()
		result = composed.GetOption(none)
		assert.True(t, O.IsNone(result))
	})
}

// TestNonEmptyStringValidation tests NonEmptyString for validation scenarios
func TestNonEmptyStringValidation(t *testing.T) {
	t.Run("validate username", func(t *testing.T) {
		prism := NonEmptyString()

		// Valid username
		validUsername := "john_doe"
		result := prism.GetOption(validUsername)
		assert.True(t, O.IsSome(result))

		// Invalid empty username
		emptyUsername := ""
		result = prism.GetOption(emptyUsername)
		assert.True(t, O.IsNone(result))
	})

	t.Run("validate configuration value", func(t *testing.T) {
		prism := NonEmptyString()

		// Valid config value
		configValue := "production"
		result := prism.GetOption(configValue)
		assert.True(t, O.IsSome(result))

		// Invalid empty config
		emptyConfig := ""
		result = prism.GetOption(emptyConfig)
		assert.True(t, O.IsNone(result))
	})

	t.Run("filter non-empty strings from slice", func(t *testing.T) {
		prism := NonEmptyString()

		inputs := []string{"hello", "", "world", "", "test"}
		var nonEmpty []string

		for _, input := range inputs {
			if result := prism.GetOption(input); O.IsSome(result) {
				nonEmpty = append(nonEmpty, O.GetOrElse(F.Constant(""))(result))
			}
		}

		assert.Equal(t, []string{"hello", "world", "test"}, nonEmpty)
	})
}

// TestFromResult tests the FromResult prism with Result types
func TestFromResult(t *testing.T) {
	t.Run("extract from successful result", func(t *testing.T) {
		prism := FromResult[int]()

		success := result.Of[int](42)
		extracted := prism.GetOption(success)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("extract from error result", func(t *testing.T) {
		prism := FromResult[int]()

		failure := E.Left[int](errors.New("test error"))
		extracted := prism.GetOption(failure)

		assert.True(t, O.IsNone(extracted))
	})

	t.Run("ReverseGet wraps value in successful result", func(t *testing.T) {
		prism := FromResult[int]()

		wrapped := prism.ReverseGet(100)

		// Verify it's a successful result
		extracted := prism.GetOption(wrapped)
		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, 100, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("works with string type", func(t *testing.T) {
		prism := FromResult[string]()

		success := result.Of[string]("hello")
		extracted := prism.GetOption(success)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, "hello", O.GetOrElse(F.Constant(""))(extracted))
	})

	t.Run("works with struct type", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		prism := FromResult[Person]()

		person := Person{Name: "Alice", Age: 30}
		success := result.Of[Person](person)
		extracted := prism.GetOption(success)

		assert.True(t, O.IsSome(extracted))
		result := O.GetOrElse(F.Constant(Person{}))(extracted)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})
}

// TestFromResultWithSet tests using Set with FromResult prism
func TestFromResultWithSet(t *testing.T) {
	t.Run("set on successful result", func(t *testing.T) {
		prism := FromResult[int]()
		setter := Set[result.Result[int], int](200)

		success := result.Of[int](42)
		updated := setter(prism)(success)

		// Verify the value was updated
		extracted := prism.GetOption(updated)
		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, 200, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("set on error result leaves it unchanged", func(t *testing.T) {
		prism := FromResult[int]()
		setter := Set[result.Result[int], int](200)

		failure := E.Left[int](errors.New("test error"))
		updated := setter(prism)(failure)

		// Verify it's still an error
		extracted := prism.GetOption(updated)
		assert.True(t, O.IsNone(extracted))
	})
}

// TestFromResultPrismLaws tests that FromResult satisfies prism laws
func TestFromResultPrismLaws(t *testing.T) {
	prism := FromResult[int]()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		value := 42
		wrapped := prism.ReverseGet(value)
		extracted := prism.GetOption(wrapped)

		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, value, O.GetOrElse(F.Constant(-1))(extracted))
	})

	t.Run("law 2: ReverseGet is consistent", func(t *testing.T) {
		value := 42
		result1 := prism.ReverseGet(value)
		result2 := prism.ReverseGet(value)

		// Both should extract the same value
		extracted1 := prism.GetOption(result1)
		extracted2 := prism.GetOption(result2)

		val1 := O.GetOrElse(F.Constant(-1))(extracted1)
		val2 := O.GetOrElse(F.Constant(-1))(extracted2)
		assert.Equal(t, val1, val2)
	})
}

// TestFromResultComposition tests composing FromResult with other prisms
func TestFromResultComposition(t *testing.T) {
	t.Run("compose with predicate prism", func(t *testing.T) {
		// Create a prism that only matches positive numbers
		positivePrism := FromPredicate(func(n int) bool { return n > 0 })

		// Compose: Result[int] -> int -> positive int
		composed := Compose[result.Result[int]](positivePrism)(FromResult[int]())

		// Test with positive number
		success := result.Of[int](42)
		extracted := composed.GetOption(success)
		assert.True(t, O.IsSome(extracted))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(extracted))

		// Test with negative number
		negativeSuccess := result.Of[int](-5)
		extracted = composed.GetOption(negativeSuccess)
		assert.True(t, O.IsNone(extracted))

		// Test with error
		failure := E.Left[int](errors.New("test error"))
		extracted = composed.GetOption(failure)
		assert.True(t, O.IsNone(extracted))
	})
}

// TestParseJSON tests the ParseJSON prism with various JSON data
func TestParseJSON(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("parse valid JSON", func(t *testing.T) {
		prism := ParseJSON[Person]()

		jsonData := []byte(`{"name":"Alice","age":30}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "Alice", person.Name)
		assert.Equal(t, 30, person.Age)
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		prism := ParseJSON[Person]()

		invalidJSON := []byte(`{invalid json}`)
		parsed := prism.GetOption(invalidJSON)

		assert.True(t, O.IsNone(parsed))
	})

	t.Run("parse JSON with missing fields", func(t *testing.T) {
		prism := ParseJSON[Person]()

		// Missing age field - should use zero value
		jsonData := []byte(`{"name":"Bob"}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "Bob", person.Name)
		assert.Equal(t, 0, person.Age)
	})

	t.Run("parse JSON with extra fields", func(t *testing.T) {
		prism := ParseJSON[Person]()

		// Extra field should be ignored
		jsonData := []byte(`{"name":"Charlie","age":25,"extra":"ignored"}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "Charlie", person.Name)
		assert.Equal(t, 25, person.Age)
	})

	t.Run("ReverseGet marshals to JSON", func(t *testing.T) {
		prism := ParseJSON[Person]()

		person := Person{Name: "David", Age: 35}
		jsonBytes := prism.ReverseGet(person)

		// Parse it back to verify
		parsed := prism.GetOption(jsonBytes)
		assert.True(t, O.IsSome(parsed))
		result := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "David", result.Name)
		assert.Equal(t, 35, result.Age)
	})

	t.Run("works with primitive types", func(t *testing.T) {
		prism := ParseJSON[int]()

		jsonData := []byte(`42`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(parsed))
	})

	t.Run("works with arrays", func(t *testing.T) {
		prism := ParseJSON[[]string]()

		jsonData := []byte(`["hello","world","test"]`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		arr := O.GetOrElse(F.Constant([]string{}))(parsed)
		assert.Equal(t, []string{"hello", "world", "test"}, arr)
	})

	t.Run("works with maps", func(t *testing.T) {
		prism := ParseJSON[map[string]int]()

		jsonData := []byte(`{"a":1,"b":2,"c":3}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		m := O.GetOrElse(F.Constant(map[string]int{}))(parsed)
		assert.Equal(t, 1, m["a"])
		assert.Equal(t, 2, m["b"])
		assert.Equal(t, 3, m["c"])
	})

	t.Run("works with nested structures", func(t *testing.T) {
		type Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}
		type PersonWithAddress struct {
			Name    string  `json:"name"`
			Address Address `json:"address"`
		}

		prism := ParseJSON[PersonWithAddress]()

		jsonData := []byte(`{"name":"Eve","address":{"street":"123 Main St","city":"NYC"}}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(PersonWithAddress{}))(parsed)
		assert.Equal(t, "Eve", person.Name)
		assert.Equal(t, "123 Main St", person.Address.Street)
		assert.Equal(t, "NYC", person.Address.City)
	})

	t.Run("parse empty JSON object", func(t *testing.T) {
		prism := ParseJSON[Person]()

		jsonData := []byte(`{}`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "", person.Name)
		assert.Equal(t, 0, person.Age)
	})

	t.Run("parse null JSON", func(t *testing.T) {
		prism := ParseJSON[*Person]()

		jsonData := []byte(`null`)
		parsed := prism.GetOption(jsonData)

		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(&Person{}))(parsed)
		assert.Nil(t, person)
	})
}

// TestParseJSONWithSet tests using Set with ParseJSON prism
func TestParseJSONWithSet(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("set updates JSON data", func(t *testing.T) {
		prism := ParseJSON[Person]()

		originalJSON := []byte(`{"name":"Alice","age":30}`)
		newPerson := Person{Name: "Bob", Age: 25}

		setter := Set[[]byte, Person](newPerson)
		updatedJSON := setter(prism)(originalJSON)

		// Parse the updated JSON
		parsed := prism.GetOption(updatedJSON)
		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "Bob", person.Name)
		assert.Equal(t, 25, person.Age)
	})

	t.Run("set on invalid JSON returns original unchanged", func(t *testing.T) {
		prism := ParseJSON[Person]()

		invalidJSON := []byte(`{invalid}`)
		newPerson := Person{Name: "Charlie", Age: 35}

		setter := Set[[]byte, Person](newPerson)
		result := setter(prism)(invalidJSON)

		// Should return original unchanged since it couldn't be parsed
		assert.Equal(t, invalidJSON, result)
	})
}

// TestParseJSONPrismLaws tests that ParseJSON satisfies prism laws
func TestParseJSONPrismLaws(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	prism := ParseJSON[Person]()

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}
		jsonBytes := prism.ReverseGet(person)
		parsed := prism.GetOption(jsonBytes)

		assert.True(t, O.IsSome(parsed))
		result := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, person.Name, result.Name)
		assert.Equal(t, person.Age, result.Age)
	})

	t.Run("law 2: ReverseGet is consistent", func(t *testing.T) {
		person := Person{Name: "Bob", Age: 25}
		json1 := prism.ReverseGet(person)
		json2 := prism.ReverseGet(person)

		// Both should parse to the same value
		parsed1 := prism.GetOption(json1)
		parsed2 := prism.GetOption(json2)

		result1 := O.GetOrElse(F.Constant(Person{}))(parsed1)
		result2 := O.GetOrElse(F.Constant(Person{}))(parsed2)

		assert.Equal(t, result1.Name, result2.Name)
		assert.Equal(t, result1.Age, result2.Age)
	})
}

// TestParseJSONComposition tests composing ParseJSON with other prisms
func TestParseJSONComposition(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("compose with predicate prism", func(t *testing.T) {
		// Create a prism that only matches adults (age >= 18)
		adultPrism := FromPredicate(func(p Person) bool { return p.Age >= 18 })

		// Compose: []byte -> Person -> Adult
		composed := Compose[[]byte](adultPrism)(ParseJSON[Person]())

		// Test with adult
		adultJSON := []byte(`{"name":"Alice","age":30}`)
		parsed := composed.GetOption(adultJSON)
		assert.True(t, O.IsSome(parsed))
		person := O.GetOrElse(F.Constant(Person{}))(parsed)
		assert.Equal(t, "Alice", person.Name)

		// Test with minor
		minorJSON := []byte(`{"name":"Bob","age":15}`)
		parsed = composed.GetOption(minorJSON)
		assert.True(t, O.IsNone(parsed))

		// Test with invalid JSON
		invalidJSON := []byte(`{invalid}`)
		parsed = composed.GetOption(invalidJSON)
		assert.True(t, O.IsNone(parsed))
	})
}
