// Copyright (c) 2024 IBM Corp.
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

package codec

import (
	"encoding/json"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// URL
// ---------------------------------------------------------------------------

func TestURL_Decode_Success(t *testing.T) {
	t.Run("decodes valid absolute URL", func(t *testing.T) {
		c := URL()
		result := c.Decode("https://example.com/path?query=value")
		require.True(t, either.IsRight(result))
		u := either.MonadFold(result, func(validation.Errors) *url.URL { return nil }, func(u *url.URL) *url.URL { return u })
		assert.Equal(t, "https", u.Scheme)
		assert.Equal(t, "example.com", u.Host)
		assert.Equal(t, "/path", u.Path)
	})

	t.Run("decodes simple URL", func(t *testing.T) {
		c := URL()
		result := c.Decode("https://example.com")
		assert.True(t, either.IsRight(result))
	})

	t.Run("decodes URL with port", func(t *testing.T) {
		c := URL()
		result := c.Decode("http://localhost:8080/api")
		require.True(t, either.IsRight(result))
		u := either.MonadFold(result, func(validation.Errors) *url.URL { return nil }, func(u *url.URL) *url.URL { return u })
		assert.Equal(t, "localhost:8080", u.Host)
	})

	t.Run("decodes relative URL", func(t *testing.T) {
		c := URL()
		// url.Parse accepts relative URLs too
		result := c.Decode("/relative/path")
		assert.True(t, either.IsRight(result))
	})
}

func TestURL_Decode_Failure(t *testing.T) {
	t.Run("fails on URL with invalid characters", func(t *testing.T) {
		c := URL()
		// url.Parse is very permissive; use a string with a control character
		result := c.Decode("://\x00invalid")
		assert.True(t, either.IsLeft(result))
	})
}

func TestURL_Encode(t *testing.T) {
	t.Run("encodes URL to string", func(t *testing.T) {
		c := URL()
		u, err := url.Parse("https://example.com/path")
		require.NoError(t, err)
		encoded := c.Encode(u)
		assert.Equal(t, "https://example.com/path", encoded)
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := URL()
		original := "https://example.com/path?q=1"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		u := either.MonadFold(result, func(validation.Errors) *url.URL { return nil }, func(u *url.URL) *url.URL { return u })
		assert.Equal(t, original, c.Encode(u))
	})
}

func TestURL_Name(t *testing.T) {
	assert.Equal(t, "URL", URL().Name())
}

// ---------------------------------------------------------------------------
// Date
// ---------------------------------------------------------------------------

func TestDate_Decode_Success(t *testing.T) {
	layout := "2006-01-02"

	t.Run("decodes valid date string", func(t *testing.T) {
		c := Date(layout)
		result := c.Decode("2024-03-15")
		require.True(t, either.IsRight(result))
		got := either.MonadFold(result, func(validation.Errors) time.Time { return time.Time{} }, func(t time.Time) time.Time { return t })
		assert.Equal(t, 2024, got.Year())
		assert.Equal(t, time.March, got.Month())
		assert.Equal(t, 15, got.Day())
	})

	t.Run("decodes RFC3339 timestamp", func(t *testing.T) {
		c := Date(time.RFC3339)
		result := c.Decode("2024-03-15T10:30:00Z")
		assert.True(t, either.IsRight(result))
	})
}

func TestDate_Decode_Failure(t *testing.T) {
	layout := "2006-01-02"

	t.Run("fails on wrong format", func(t *testing.T) {
		c := Date(layout)
		result := c.Decode("15-03-2024")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on non-date string", func(t *testing.T) {
		c := Date(layout)
		result := c.Decode("not a date")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on empty string", func(t *testing.T) {
		c := Date(layout)
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})
}

func TestDate_Encode(t *testing.T) {
	layout := "2006-01-02"

	t.Run("encodes time.Time to string", func(t *testing.T) {
		c := Date(layout)
		tm := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		encoded := c.Encode(tm)
		assert.Equal(t, "2024-03-15", encoded)
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := Date(layout)
		original := "2024-03-15"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		tm := either.MonadFold(result, func(validation.Errors) time.Time { return time.Time{} }, func(t time.Time) time.Time { return t })
		assert.Equal(t, original, c.Encode(tm))
	})
}

func TestDate_Name(t *testing.T) {
	assert.Equal(t, "Date", Date("2006-01-02").Name())
}

// ---------------------------------------------------------------------------
// Regex
// ---------------------------------------------------------------------------

func TestRegex_Decode_Success(t *testing.T) {
	t.Run("decodes matching string", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		c := Regex(re)
		result := c.Decode("Price: 42 dollars")
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.Match { return prism.Match{} }, func(m prism.Match) prism.Match { return m })
		assert.Equal(t, "Price: ", m.Before)
		assert.Equal(t, "42", m.Groups[0])
		assert.Equal(t, " dollars", m.After)
	})

	t.Run("decodes with capture groups", func(t *testing.T) {
		re := regexp.MustCompile(`(\w+)@(\w+\.\w+)`)
		c := Regex(re)
		result := c.Decode("user@example.com")
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.Match { return prism.Match{} }, func(m prism.Match) prism.Match { return m })
		assert.Equal(t, "user@example.com", m.Groups[0])
		assert.Equal(t, "user", m.Groups[1])
		assert.Equal(t, "example.com", m.Groups[2])
	})
}

func TestRegex_Decode_Failure(t *testing.T) {
	t.Run("fails on non-matching string", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		c := Regex(re)
		result := c.Decode("no numbers here")
		assert.True(t, either.IsLeft(result))
	})
}

func TestRegex_Encode(t *testing.T) {
	t.Run("encodes Match back to original string", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		c := Regex(re)
		m := prism.Match{Before: "Price: ", Groups: []string{"42"}, After: " dollars"}
		encoded := c.Encode(m)
		assert.Equal(t, "Price: 42 dollars", encoded)
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		c := Regex(re)
		original := "Price: 42 dollars"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.Match { return prism.Match{} }, func(m prism.Match) prism.Match { return m })
		assert.Equal(t, original, c.Encode(m))
	})
}

// ---------------------------------------------------------------------------
// RegexNamed
// ---------------------------------------------------------------------------

func TestRegexNamed_Decode_Success(t *testing.T) {
	t.Run("decodes string with named capture groups", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		c := RegexNamed(re)
		result := c.Decode("john@example.com")
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.NamedMatch { return prism.NamedMatch{} }, func(m prism.NamedMatch) prism.NamedMatch { return m })
		assert.Equal(t, "john@example.com", m.Full)
		assert.Equal(t, "john", m.Groups["user"])
		assert.Equal(t, "example.com", m.Groups["domain"])
		assert.Equal(t, "", m.Before)
		assert.Equal(t, "", m.After)
	})

	t.Run("captures before and after text", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<num>\d+)`)
		c := RegexNamed(re)
		result := c.Decode("value: 42 end")
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.NamedMatch { return prism.NamedMatch{} }, func(m prism.NamedMatch) prism.NamedMatch { return m })
		assert.Equal(t, "value: ", m.Before)
		assert.Equal(t, "42", m.Groups["num"])
		assert.Equal(t, " end", m.After)
	})
}

func TestRegexNamed_Decode_Failure(t *testing.T) {
	t.Run("fails on non-matching string", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		c := RegexNamed(re)
		result := c.Decode("not-an-email")
		assert.True(t, either.IsLeft(result))
	})
}

func TestRegexNamed_Encode(t *testing.T) {
	t.Run("encodes NamedMatch back to original string", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		c := RegexNamed(re)
		m := prism.NamedMatch{
			Before: "",
			Groups: map[string]string{"user": "john", "domain": "example.com"},
			Full:   "john@example.com",
			After:  "",
		}
		encoded := c.Encode(m)
		assert.Equal(t, "john@example.com", encoded)
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		re := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		c := RegexNamed(re)
		original := "john@example.com"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		m := either.MonadFold(result, func(validation.Errors) prism.NamedMatch { return prism.NamedMatch{} }, func(m prism.NamedMatch) prism.NamedMatch { return m })
		assert.Equal(t, original, c.Encode(m))
	})
}

// ---------------------------------------------------------------------------
// IntFromString
// ---------------------------------------------------------------------------

func TestIntFromString_Decode_Success(t *testing.T) {
	t.Run("decodes positive integer", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("42")
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("decodes negative integer", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("-123")
		assert.Equal(t, validation.Success(-123), result)
	})

	t.Run("decodes zero", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("0")
		assert.Equal(t, validation.Success(0), result)
	})

	t.Run("decodes integer with leading plus sign", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("+7")
		assert.Equal(t, validation.Success(7), result)
	})
}

func TestIntFromString_Decode_Failure(t *testing.T) {
	t.Run("fails on non-numeric string", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("not a number")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on floating point", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("3.14")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on empty string", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on hex string", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("0xFF")
		assert.True(t, either.IsLeft(result))
	})
}

func TestIntFromString_Encode(t *testing.T) {
	t.Run("encodes int to string", func(t *testing.T) {
		c := IntFromString()
		assert.Equal(t, "42", c.Encode(42))
	})

	t.Run("encodes negative int to string", func(t *testing.T) {
		c := IntFromString()
		assert.Equal(t, "-123", c.Encode(-123))
	})

	t.Run("encodes zero to string", func(t *testing.T) {
		c := IntFromString()
		assert.Equal(t, "0", c.Encode(0))
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := IntFromString()
		result := c.Decode("42")
		require.True(t, either.IsRight(result))
		n := either.MonadFold(result, func(validation.Errors) int { return 0 }, func(n int) int { return n })
		assert.Equal(t, "42", c.Encode(n))
	})
}

func TestIntFromString_Name(t *testing.T) {
	assert.Equal(t, "IntFromString", IntFromString().Name())
}

// ---------------------------------------------------------------------------
// Int64FromString
// ---------------------------------------------------------------------------

func TestInt64FromString_Decode_Success(t *testing.T) {
	t.Run("decodes max int64", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("9223372036854775807")
		assert.Equal(t, validation.Success(int64(9223372036854775807)), result)
	})

	t.Run("decodes min int64", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("-9223372036854775808")
		assert.Equal(t, validation.Success(int64(-9223372036854775808)), result)
	})

	t.Run("decodes zero", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("0")
		assert.Equal(t, validation.Success(int64(0)), result)
	})

	t.Run("decodes regular int64", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("42")
		assert.Equal(t, validation.Success(int64(42)), result)
	})
}

func TestInt64FromString_Decode_Failure(t *testing.T) {
	t.Run("fails on non-numeric string", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("not a number")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on overflow", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("9223372036854775808") // max int64 + 1
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on floating point", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("3.14")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on empty string", func(t *testing.T) {
		c := Int64FromString()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})
}

func TestInt64FromString_Encode(t *testing.T) {
	t.Run("encodes int64 to string", func(t *testing.T) {
		c := Int64FromString()
		assert.Equal(t, "42", c.Encode(42))
	})

	t.Run("encodes max int64 to string", func(t *testing.T) {
		c := Int64FromString()
		assert.Equal(t, "9223372036854775807", c.Encode(9223372036854775807))
	})

	t.Run("encodes negative int64 to string", func(t *testing.T) {
		c := Int64FromString()
		assert.Equal(t, "-9223372036854775808", c.Encode(-9223372036854775808))
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := Int64FromString()
		original := "9223372036854775807"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		n := either.MonadFold(result, func(validation.Errors) int64 { return 0 }, func(n int64) int64 { return n })
		assert.Equal(t, original, c.Encode(n))
	})
}

func TestInt64FromString_Name(t *testing.T) {
	assert.Equal(t, "Int64FromString", Int64FromString().Name())
}

// ---------------------------------------------------------------------------
// BoolFromString
// ---------------------------------------------------------------------------

func TestBoolFromString_Decode_Success(t *testing.T) {
	t.Run("decodes 'true' string", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("true")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes 'false' string", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("false")
		assert.Equal(t, validation.Success(false), result)
	})

	t.Run("decodes '1' as true", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("1")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes '0' as false", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("0")
		assert.Equal(t, validation.Success(false), result)
	})

	t.Run("decodes 't' as true", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("t")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes 'f' as false", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("f")
		assert.Equal(t, validation.Success(false), result)
	})

	t.Run("decodes 'T' as true", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("T")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes 'F' as false", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("F")
		assert.Equal(t, validation.Success(false), result)
	})

	t.Run("decodes 'TRUE' as true", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("TRUE")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes 'FALSE' as false", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("FALSE")
		assert.Equal(t, validation.Success(false), result)
	})

	t.Run("decodes 'True' as true", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("True")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("decodes 'False' as false", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("False")
		assert.Equal(t, validation.Success(false), result)
	})
}

func TestBoolFromString_Decode_Failure(t *testing.T) {
	t.Run("fails on 'yes'", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("yes")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on 'no'", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("no")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on empty string", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on numeric string other than 0 or 1", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("2")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on arbitrary text", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("not a boolean")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on whitespace", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode(" ")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on 'true' with leading/trailing spaces", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode(" true ")
		assert.True(t, either.IsLeft(result))
	})
}

func TestBoolFromString_Encode(t *testing.T) {
	t.Run("encodes true to 'true'", func(t *testing.T) {
		c := BoolFromString()
		assert.Equal(t, "true", c.Encode(true))
	})

	t.Run("encodes false to 'false'", func(t *testing.T) {
		c := BoolFromString()
		assert.Equal(t, "false", c.Encode(false))
	})

	t.Run("round-trip: decode 'true' then encode", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("true")
		require.True(t, either.IsRight(result))
		b := either.MonadFold(result, func(validation.Errors) bool { return false }, func(b bool) bool { return b })
		assert.Equal(t, "true", c.Encode(b))
	})

	t.Run("round-trip: decode 'false' then encode", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("false")
		require.True(t, either.IsRight(result))
		b := either.MonadFold(result, func(validation.Errors) bool { return true }, func(b bool) bool { return b })
		assert.Equal(t, "false", c.Encode(b))
	})

	t.Run("round-trip: decode '1' encodes as 'true'", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("1")
		require.True(t, either.IsRight(result))
		b := either.MonadFold(result, func(validation.Errors) bool { return false }, func(b bool) bool { return b })
		// Note: strconv.FormatBool always returns "true" or "false", not "1" or "0"
		assert.Equal(t, "true", c.Encode(b))
	})

	t.Run("round-trip: decode '0' encodes as 'false'", func(t *testing.T) {
		c := BoolFromString()
		result := c.Decode("0")
		require.True(t, either.IsRight(result))
		b := either.MonadFold(result, func(validation.Errors) bool { return true }, func(b bool) bool { return b })
		assert.Equal(t, "false", c.Encode(b))
	})
}

func TestBoolFromString_EdgeCases(t *testing.T) {
	t.Run("case sensitivity variations", func(t *testing.T) {
		c := BoolFromString()
		cases := []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"True", true},
			{"TRUE", true},
			{"false", false},
			{"False", false},
			{"FALSE", false},
			{"t", true},
			{"T", true},
			{"f", false},
			{"F", false},
		}
		for _, tc := range cases {
			result := c.Decode(tc.input)
			require.True(t, either.IsRight(result), "expected success for %s", tc.input)
			b := either.MonadFold(result, func(validation.Errors) bool { return !tc.expected }, func(b bool) bool { return b })
			assert.Equal(t, tc.expected, b, "input: %s", tc.input)
		}
	})
}

func TestBoolFromString_Name(t *testing.T) {
	assert.Equal(t, "BoolFromString", BoolFromString().Name())
}

func TestBoolFromString_Integration(t *testing.T) {
	t.Run("decodes and encodes multiple boolean values", func(t *testing.T) {
		c := BoolFromString()
		cases := []struct {
			str string
			val bool
		}{
			{"true", true},
			{"false", false},
			{"1", true},
			{"0", false},
			{"T", true},
			{"F", false},
		}
		for _, tc := range cases {
			result := c.Decode(tc.str)
			require.True(t, either.IsRight(result), "expected success for %s", tc.str)
			b := either.MonadFold(result, func(validation.Errors) bool { return !tc.val }, func(b bool) bool { return b })
			assert.Equal(t, tc.val, b)
			// Note: encoding always produces "true" or "false", not the original input
			if tc.val {
				assert.Equal(t, "true", c.Encode(b))
			} else {
				assert.Equal(t, "false", c.Encode(b))
			}
		}
	})
}

// ---------------------------------------------------------------------------
// FromNonZero
// ---------------------------------------------------------------------------

func TestFromNonZero_Decode_Success(t *testing.T) {
	t.Run("int - decodes non-zero value", func(t *testing.T) {
		c := FromNonZero[int]()
		result := c.Decode(42)
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("int - decodes negative value", func(t *testing.T) {
		c := FromNonZero[int]()
		result := c.Decode(-5)
		assert.Equal(t, validation.Success(-5), result)
	})

	t.Run("string - decodes non-empty string", func(t *testing.T) {
		c := FromNonZero[string]()
		result := c.Decode("hello")
		assert.Equal(t, validation.Success("hello"), result)
	})

	t.Run("string - decodes whitespace string", func(t *testing.T) {
		c := FromNonZero[string]()
		result := c.Decode("   ")
		assert.Equal(t, validation.Success("   "), result)
	})

	t.Run("bool - decodes true", func(t *testing.T) {
		c := FromNonZero[bool]()
		result := c.Decode(true)
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("float64 - decodes non-zero value", func(t *testing.T) {
		c := FromNonZero[float64]()
		result := c.Decode(3.14)
		assert.Equal(t, validation.Success(3.14), result)
	})

	t.Run("float64 - decodes negative value", func(t *testing.T) {
		c := FromNonZero[float64]()
		result := c.Decode(-2.5)
		assert.Equal(t, validation.Success(-2.5), result)
	})

	t.Run("pointer - decodes non-nil pointer", func(t *testing.T) {
		c := FromNonZero[*int]()
		value := 42
		result := c.Decode(&value)
		assert.True(t, either.IsRight(result))
		ptr := either.MonadFold(result, func(validation.Errors) *int { return nil }, func(p *int) *int { return p })
		require.NotNil(t, ptr)
		assert.Equal(t, 42, *ptr)
	})
}

func TestFromNonZero_Decode_Failure(t *testing.T) {
	t.Run("int - fails on zero", func(t *testing.T) {
		c := FromNonZero[int]()
		result := c.Decode(0)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("string - fails on empty string", func(t *testing.T) {
		c := FromNonZero[string]()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("bool - fails on false", func(t *testing.T) {
		c := FromNonZero[bool]()
		result := c.Decode(false)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("float64 - fails on zero", func(t *testing.T) {
		c := FromNonZero[float64]()
		result := c.Decode(0.0)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("pointer - fails on nil", func(t *testing.T) {
		c := FromNonZero[*int]()
		result := c.Decode(nil)
		assert.True(t, either.IsLeft(result))
	})
}

func TestFromNonZero_Encode(t *testing.T) {
	t.Run("int - encodes value unchanged", func(t *testing.T) {
		c := FromNonZero[int]()
		assert.Equal(t, 42, c.Encode(42))
	})

	t.Run("string - encodes value unchanged", func(t *testing.T) {
		c := FromNonZero[string]()
		assert.Equal(t, "hello", c.Encode("hello"))
	})

	t.Run("bool - encodes value unchanged", func(t *testing.T) {
		c := FromNonZero[bool]()
		assert.Equal(t, true, c.Encode(true))
	})

	t.Run("float64 - encodes value unchanged", func(t *testing.T) {
		c := FromNonZero[float64]()
		assert.Equal(t, 3.14, c.Encode(3.14))
	})

	t.Run("pointer - encodes value unchanged", func(t *testing.T) {
		c := FromNonZero[*int]()
		value := 42
		ptr := &value
		assert.Equal(t, ptr, c.Encode(ptr))
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := FromNonZero[int]()
		original := 42
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		decoded := either.MonadFold(result, func(validation.Errors) int { return 0 }, func(n int) int { return n })
		assert.Equal(t, original, c.Encode(decoded))
	})
}

func TestFromNonZero_Name(t *testing.T) {
	t.Run("int codec name", func(t *testing.T) {
		c := FromNonZero[int]()
		assert.Contains(t, c.Name(), "FromRefinement")
		assert.Contains(t, c.Name(), "PrismFromNonZero")
	})

	t.Run("string codec name", func(t *testing.T) {
		c := FromNonZero[string]()
		assert.Contains(t, c.Name(), "FromRefinement")
		assert.Contains(t, c.Name(), "PrismFromNonZero")
	})
}

func TestFromNonZero_Integration(t *testing.T) {
	t.Run("validates multiple non-zero integers", func(t *testing.T) {
		c := FromNonZero[int]()
		values := []int{1, -1, 42, -100, 999}
		for _, v := range values {
			result := c.Decode(v)
			require.True(t, either.IsRight(result), "expected success for %d", v)
			decoded := either.MonadFold(result, func(validation.Errors) int { return 0 }, func(n int) int { return n })
			assert.Equal(t, v, decoded)
			assert.Equal(t, v, c.Encode(decoded))
		}
	})

	t.Run("rejects zero values", func(t *testing.T) {
		c := FromNonZero[int]()
		result := c.Decode(0)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("works with custom comparable types", func(t *testing.T) {
		type UserID string
		c := FromNonZero[UserID]()

		result := c.Decode(UserID("user123"))
		assert.Equal(t, validation.Success(UserID("user123")), result)

		result = c.Decode(UserID(""))
		assert.True(t, either.IsLeft(result))
	})
}

// ---------------------------------------------------------------------------
// NonEmptyString
// ---------------------------------------------------------------------------

func TestNonEmptyString_Decode_Success(t *testing.T) {
	t.Run("decodes non-empty string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("hello")
		assert.Equal(t, validation.Success("hello"), result)
	})

	t.Run("decodes single character", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("a")
		assert.Equal(t, validation.Success("a"), result)
	})

	t.Run("decodes whitespace string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("   ")
		assert.Equal(t, validation.Success("   "), result)
	})

	t.Run("decodes string with newlines", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("\n\t")
		assert.Equal(t, validation.Success("\n\t"), result)
	})

	t.Run("decodes unicode string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("你好")
		assert.Equal(t, validation.Success("你好"), result)
	})

	t.Run("decodes emoji string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("🎉")
		assert.Equal(t, validation.Success("🎉"), result)
	})

	t.Run("decodes multiline string", func(t *testing.T) {
		c := NonEmptyString()
		multiline := "line1\nline2\nline3"
		result := c.Decode(multiline)
		assert.Equal(t, validation.Success(multiline), result)
	})
}

func TestNonEmptyString_Decode_Failure(t *testing.T) {
	t.Run("fails on empty string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("error contains context", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("")
		require.True(t, either.IsLeft(result))
		errors := either.MonadFold(result, func(e validation.Errors) validation.Errors { return e }, func(string) validation.Errors { return nil })
		require.NotEmpty(t, errors)
	})
}

func TestNonEmptyString_Encode(t *testing.T) {
	t.Run("encodes string unchanged", func(t *testing.T) {
		c := NonEmptyString()
		assert.Equal(t, "hello", c.Encode("hello"))
	})

	t.Run("encodes unicode string unchanged", func(t *testing.T) {
		c := NonEmptyString()
		assert.Equal(t, "你好", c.Encode("你好"))
	})

	t.Run("encodes whitespace string unchanged", func(t *testing.T) {
		c := NonEmptyString()
		assert.Equal(t, "   ", c.Encode("   "))
	})

	t.Run("round-trip: decode then encode", func(t *testing.T) {
		c := NonEmptyString()
		original := "test string"
		result := c.Decode(original)
		require.True(t, either.IsRight(result))
		decoded := either.MonadFold(result, func(validation.Errors) string { return "" }, func(s string) string { return s })
		assert.Equal(t, original, c.Encode(decoded))
	})
}

func TestNonEmptyString_Name(t *testing.T) {
	c := NonEmptyString()
	assert.Equal(t, c.Name(), "NonEmptyString")
}

func TestNonEmptyString_Integration(t *testing.T) {
	t.Run("validates multiple non-empty strings", func(t *testing.T) {
		c := NonEmptyString()
		strings := []string{"a", "hello", "world", "test123", "  spaces  ", "🎉"}
		for _, s := range strings {
			result := c.Decode(s)
			require.True(t, either.IsRight(result), "expected success for %q", s)
			decoded := either.MonadFold(result, func(validation.Errors) string { return "" }, func(str string) string { return str })
			assert.Equal(t, s, decoded)
			assert.Equal(t, s, c.Encode(decoded))
		}
	})

	t.Run("rejects empty string", func(t *testing.T) {
		c := NonEmptyString()
		result := c.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("compose with IntFromString", func(t *testing.T) {
		// Create a codec that only parses non-empty strings to integers
		nonEmptyThenInt := Pipe[string, string](IntFromString())(NonEmptyString())

		// Valid non-empty string with integer
		result := nonEmptyThenInt.Decode("42")
		assert.Equal(t, validation.Success(42), result)

		// Empty string fails at NonEmptyString stage
		result = nonEmptyThenInt.Decode("")
		assert.True(t, either.IsLeft(result))

		// Non-empty but invalid integer fails at IntFromString stage
		result = nonEmptyThenInt.Decode("abc")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("use in validation pipeline", func(t *testing.T) {
		c := NonEmptyString()

		// Simulate validating user input
		inputs := []struct {
			value    string
			expected bool
		}{
			{"john_doe", true},
			{"", false},
			{"a", true},
			{"user@example.com", true},
		}

		for _, input := range inputs {
			result := c.Decode(input.value)
			if input.expected {
				assert.True(t, either.IsRight(result), "expected success for %q", input.value)
			} else {
				assert.True(t, either.IsLeft(result), "expected failure for %q", input.value)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// WithName
// ---------------------------------------------------------------------------

func TestWithName_BasicFunctionality(t *testing.T) {
	t.Run("renames codec without changing behavior", func(t *testing.T) {
		original := IntFromString()
		renamed := WithName[int, string, string]("CustomIntCodec")(original)

		// Name should be changed
		assert.Equal(t, "CustomIntCodec", renamed.Name())
		assert.NotEqual(t, original.Name(), renamed.Name())

		// Behavior should be unchanged
		result := renamed.Decode("42")
		assert.Equal(t, validation.Success(42), result)

		encoded := renamed.Encode(42)
		assert.Equal(t, "42", encoded)
	})

	t.Run("preserves validation logic", func(t *testing.T) {
		original := IntFromString()
		renamed := WithName[int, string, string]("MyInt")(original)

		// Valid input should succeed
		result := renamed.Decode("123")
		assert.True(t, either.IsRight(result))

		// Invalid input should fail
		result = renamed.Decode("not a number")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("preserves encoding logic", func(t *testing.T) {
		original := BoolFromString()
		renamed := WithName[bool, string, string]("CustomBool")(original)

		assert.Equal(t, "true", renamed.Encode(true))
		assert.Equal(t, "false", renamed.Encode(false))
	})
}

func TestWithName_WithComposedCodecs(t *testing.T) {
	t.Run("renames composed codec", func(t *testing.T) {
		// Create a composed codec
		composed := Pipe[string, string](IntFromString())(NonEmptyString())

		// Rename it
		renamed := WithName[int, string, string]("NonEmptyIntString")(composed)

		assert.Equal(t, "NonEmptyIntString", renamed.Name())

		// Behavior should be preserved
		result := renamed.Decode("42")
		assert.Equal(t, validation.Success(42), result)

		// Empty string should fail
		result = renamed.Decode("")
		assert.True(t, either.IsLeft(result))

		// Non-numeric should fail
		result = renamed.Decode("abc")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("works in pipeline with F.Pipe", func(t *testing.T) {
		codec := F.Pipe1(
			IntFromString(),
			WithName[int, string, string]("UserAge"),
		)

		assert.Equal(t, "UserAge", codec.Name())

		result := codec.Decode("25")
		assert.Equal(t, validation.Success(25), result)
	})
}

func TestWithName_PreservesTypeChecking(t *testing.T) {
	t.Run("preserves Is function", func(t *testing.T) {
		original := String()
		renamed := WithName[string, string, any]("CustomString")(original)

		// Should accept string
		result := renamed.Is("hello")
		assert.True(t, either.IsRight(result))

		// Should reject non-string
		result = renamed.Is(42)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("preserves complex type checking", func(t *testing.T) {
		original := Array(Int())
		renamed := WithName[[]int, []int, any]("IntArray")(original)

		// Should accept []int
		result := renamed.Is([]int{1, 2, 3})
		assert.True(t, either.IsRight(result))

		// Should reject []string
		result = renamed.Is([]string{"a", "b"})
		assert.True(t, either.IsLeft(result))
	})
}

func TestWithName_RoundTrip(t *testing.T) {
	t.Run("maintains round-trip property", func(t *testing.T) {
		original := Int64FromString()
		renamed := WithName[int64, string, string]("CustomInt64")(original)

		testValues := []string{"0", "42", "-100", "9223372036854775807"}
		for _, input := range testValues {
			result := renamed.Decode(input)
			require.True(t, either.IsRight(result), "expected success for %s", input)

			decoded := either.MonadFold(result, func(validation.Errors) int64 { return 0 }, func(n int64) int64 { return n })
			encoded := renamed.Encode(decoded)
			assert.Equal(t, input, encoded)
		}
	})
}

func TestWithName_ErrorMessages(t *testing.T) {
	t.Run("custom name appears in validation context", func(t *testing.T) {
		codec := WithName[int, string, string]("PositiveInteger")(IntFromString())

		result := codec.Decode("not a number")
		require.True(t, either.IsLeft(result))

		// The error context should reference the custom name
		errors := either.MonadFold(result, func(e validation.Errors) validation.Errors { return e }, func(int) validation.Errors { return nil })
		require.NotEmpty(t, errors)

		// Check that at least one error references our custom name
		found := false
		for _, err := range errors {
			if len(err.Context) > 0 {
				for _, ctx := range err.Context {
					if ctx.Type == "PositiveInteger" {
						found = true
						break
					}
				}
			}
		}
		assert.True(t, found, "expected custom name 'PositiveInteger' in error context")
	})
}

func TestWithName_MultipleRenames(t *testing.T) {
	t.Run("can rename multiple times", func(t *testing.T) {
		codec := IntFromString()

		renamed1 := WithName[int, string, string]("FirstName")(codec)
		assert.Equal(t, "FirstName", renamed1.Name())

		renamed2 := WithName[int, string, string]("SecondName")(renamed1)
		assert.Equal(t, "SecondName", renamed2.Name())

		// Behavior should still work
		result := renamed2.Decode("42")
		assert.Equal(t, validation.Success(42), result)
	})
}

func TestWithName_WithDifferentTypes(t *testing.T) {
	t.Run("works with string codec", func(t *testing.T) {
		codec := WithName[string, string, string]("Username")(NonEmptyString())
		assert.Equal(t, "Username", codec.Name())

		result := codec.Decode("john_doe")
		assert.Equal(t, validation.Success("john_doe"), result)
	})

	t.Run("works with bool codec", func(t *testing.T) {
		codec := WithName[bool, string, string]("IsActive")(BoolFromString())
		assert.Equal(t, "IsActive", codec.Name())

		result := codec.Decode("true")
		assert.Equal(t, validation.Success(true), result)
	})

	t.Run("works with URL codec", func(t *testing.T) {
		codec := WithName[*url.URL, string, string]("WebsiteURL")(URL())
		assert.Equal(t, "WebsiteURL", codec.Name())

		result := codec.Decode("https://example.com")
		assert.True(t, either.IsRight(result))
	})

	t.Run("works with array codec", func(t *testing.T) {
		codec := WithName[[]int, []int, any]("Numbers")(Array(Int()))
		assert.Equal(t, "Numbers", codec.Name())

		result := codec.Decode([]int{1, 2, 3})
		assert.Equal(t, validation.Success([]int{1, 2, 3}), result)
	})
}

func TestWithName_AsDecoderEncoder(t *testing.T) {
	t.Run("AsDecoder returns decoder interface", func(t *testing.T) {
		codec := WithName[int, string, string]("MyInt")(IntFromString())
		decoder := codec.AsDecoder()

		result := decoder.Decode("42")
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("AsEncoder returns encoder interface", func(t *testing.T) {
		codec := WithName[int, string, string]("MyInt")(IntFromString())
		encoder := codec.AsEncoder()

		encoded := encoder.Encode(42)
		assert.Equal(t, "42", encoded)
	})
}

func TestWithName_Integration(t *testing.T) {
	t.Run("domain-specific codec names", func(t *testing.T) {
		// Create domain-specific codecs with meaningful names
		emailCodec := WithName[string, string, string]("EmailAddress")(NonEmptyString())
		phoneCodec := WithName[string, string, string]("PhoneNumber")(NonEmptyString())
		ageCodec := WithName[int, string, string]("Age")(IntFromString())

		// Test email
		result := emailCodec.Decode("user@example.com")
		assert.True(t, either.IsRight(result))
		assert.Equal(t, "EmailAddress", emailCodec.Name())

		// Test phone
		result = phoneCodec.Decode("+1234567890")
		assert.True(t, either.IsRight(result))
		assert.Equal(t, "PhoneNumber", phoneCodec.Name())

		// Test age
		ageResult := ageCodec.Decode("25")
		assert.True(t, either.IsRight(ageResult))
		assert.Equal(t, "Age", ageCodec.Name())
	})

	t.Run("naming complex validation pipelines", func(t *testing.T) {
		// Create a complex codec and give it a clear name
		positiveIntCodec := F.Pipe2(
			NonEmptyString(),
			Pipe[string, string](IntFromString()),
			WithName[int, string, string]("PositiveIntegerFromString"),
		)

		assert.Equal(t, "PositiveIntegerFromString", positiveIntCodec.Name())

		result := positiveIntCodec.Decode("42")
		assert.True(t, either.IsRight(result))

		result = positiveIntCodec.Decode("")
		assert.True(t, either.IsLeft(result))
	})
}

// ---------------------------------------------------------------------------
// MarshalJSON
// ---------------------------------------------------------------------------

// jsonTestType is a helper type that implements json.Marshaler and json.Unmarshaler
// by delegating to the standard encoding/json package.
type jsonTestType struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (j *jsonTestType) MarshalJSON() ([]byte, error) {
	type alias jsonTestType
	return json.Marshal((*alias)(j))
}

func (j *jsonTestType) UnmarshalJSON(b []byte) error {
	type alias jsonTestType
	return json.Unmarshal(b, (*alias)(j))
}

func TestMarshalJSON_Decode_Success(t *testing.T) {
	t.Run("decodes valid JSON bytes", func(t *testing.T) {
		var instance jsonTestType
		c := MarshalJSON[jsonTestType](&instance, &instance)
		result := c.Decode([]byte(`{"name":"test","value":42}`))
		// The codec calls UnmarshalJSON on the shared instance and returns zero T.
		// Verify decoding succeeds (is Right).
		assert.True(t, either.IsRight(result))
		// The shared instance should have been populated by UnmarshalJSON.
		assert.Equal(t, "test", instance.Name)
		assert.Equal(t, 42, instance.Value)
	})
}

func TestMarshalJSON_Decode_Failure(t *testing.T) {
	t.Run("fails on invalid JSON", func(t *testing.T) {
		var instance jsonTestType
		c := MarshalJSON[jsonTestType](&instance, &instance)
		result := c.Decode([]byte(`not valid json`))
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on malformed JSON object", func(t *testing.T) {
		var instance jsonTestType
		c := MarshalJSON[jsonTestType](&instance, &instance)
		result := c.Decode([]byte(`{"name": `))
		assert.True(t, either.IsLeft(result))
	})
}

func TestMarshalJSON_Encode(t *testing.T) {
	t.Run("encodes using the provided marshaler", func(t *testing.T) {
		instance := jsonTestType{Name: "hello", Value: 99}
		c := MarshalJSON[jsonTestType](&instance, &instance)
		encoded := c.Encode(instance)
		// The encode uses the shared enc instance (not the argument), so it
		// encodes whatever state enc currently holds.
		assert.NotEmpty(t, encoded)
	})
}

func TestMarshalJSON_Name(t *testing.T) {
	var instance jsonTestType
	c := MarshalJSON[jsonTestType](&instance, &instance)
	assert.Equal(t, "UnmarshalJSON", c.Name())
}

// ---------------------------------------------------------------------------
// Integration: Pipe composition with codecs from codecs.go
// ---------------------------------------------------------------------------

func TestIntFromString_PipeComposition(t *testing.T) {
	// Build a codec that parses a string into a positive int by composing
	// IntFromString with a positive-int refinement via Pipe.
	positiveIntPrism := prism.MakePrismWithName(
		func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n)
			}
			return option.None[int]()
		},
		func(n int) int { return n },
		"PositiveInt",
	)
	positiveIntCodec := Pipe[string, string](
		FromRefinement(positiveIntPrism),
	)(IntFromString())

	t.Run("decodes positive integer string", func(t *testing.T) {
		result := positiveIntCodec.Decode("42")
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("fails on zero", func(t *testing.T) {
		result := positiveIntCodec.Decode("0")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on negative integer string", func(t *testing.T) {
		result := positiveIntCodec.Decode("-5")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails on non-numeric string", func(t *testing.T) {
		result := positiveIntCodec.Decode("abc")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes positive int back to string", func(t *testing.T) {
		encoded := positiveIntCodec.Encode(42)
		assert.Equal(t, "42", encoded)
	})
}

func TestDate_Integration(t *testing.T) {
	t.Run("multiple dates round-trip correctly", func(t *testing.T) {
		c := Date("2006-01-02")
		dates := []string{"2000-01-01", "2024-12-31", "1970-06-15"}
		for _, d := range dates {
			result := c.Decode(d)
			require.True(t, either.IsRight(result), "expected success for %s", d)
			tm := either.MonadFold(result, func(validation.Errors) time.Time { return time.Time{} }, func(t time.Time) time.Time { return t })
			assert.Equal(t, d, c.Encode(tm))
		}
	})
}

func TestURL_Integration(t *testing.T) {
	t.Run("multiple URLs round-trip correctly", func(t *testing.T) {
		c := URL()
		urls := []string{
			"https://example.com",
			"http://localhost:8080/api/v1",
			"ftp://files.example.org/pub/data",
		}
		for _, raw := range urls {
			result := c.Decode(raw)
			require.True(t, either.IsRight(result), "expected success for %s", raw)
			u := either.MonadFold(result, func(validation.Errors) *url.URL { return nil }, func(u *url.URL) *url.URL { return u })
			assert.Equal(t, raw, c.Encode(u))
		}
	})
}

func TestIntFromString_Integration(t *testing.T) {
	t.Run("decodes and encodes a range of integers", func(t *testing.T) {
		c := IntFromString()
		cases := []struct {
			str string
			val int
		}{
			{"0", 0},
			{"1", 1},
			{"-1", -1},
			{"2147483647", 2147483647},
		}
		for _, tc := range cases {
			result := c.Decode(tc.str)
			require.True(t, either.IsRight(result), "expected success for %s", tc.str)
			n := either.MonadFold(result, func(validation.Errors) int { return -999 }, func(n int) int { return n })
			assert.Equal(t, tc.val, n)
			assert.Equal(t, tc.str, c.Encode(n))
		}
	})
}
