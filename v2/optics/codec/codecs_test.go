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
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURL(t *testing.T) {
	urlCodec := URL()

	getOrElseNull := either.GetOrElse(reader.Of[validation.Errors, *url.URL](nil))

	t.Run("decodes valid HTTP URL", func(t *testing.T) {
		result := urlCodec.Decode("https://example.com/path?query=value")

		assert.True(t, either.IsRight(result), "should successfully decode valid URL")

		parsedURL := getOrElseNull(result)

		require.NotNil(t, parsedURL)
		assert.Equal(t, "https", parsedURL.Scheme)
		assert.Equal(t, "example.com", parsedURL.Host)
		assert.Equal(t, "/path", parsedURL.Path)
		assert.Equal(t, "query=value", parsedURL.RawQuery)
	})

	t.Run("decodes valid HTTP URL without path", func(t *testing.T) {
		result := urlCodec.Decode("https://example.com")

		assert.True(t, either.IsRight(result))

		parsedURL := getOrElseNull(result)

		require.NotNil(t, parsedURL)
		assert.Equal(t, "https", parsedURL.Scheme)
		assert.Equal(t, "example.com", parsedURL.Host)
	})

	t.Run("decodes URL with port", func(t *testing.T) {
		result := urlCodec.Decode("http://localhost:8080/api")

		assert.True(t, either.IsRight(result))

		parsedURL := getOrElseNull(result)

		require.NotNil(t, parsedURL)
		assert.Equal(t, "http", parsedURL.Scheme)
		assert.Equal(t, "localhost:8080", parsedURL.Host)
		assert.Equal(t, "/api", parsedURL.Path)
	})

	t.Run("decodes URL with fragment", func(t *testing.T) {
		result := urlCodec.Decode("https://example.com/page#section")

		assert.True(t, either.IsRight(result))

		parsedURL := getOrElseNull(result)

		require.NotNil(t, parsedURL)
		assert.Equal(t, "section", parsedURL.Fragment)
	})

	t.Run("decodes relative URL", func(t *testing.T) {
		result := urlCodec.Decode("/path/to/resource")

		assert.True(t, either.IsRight(result))

		parsedURL := getOrElseNull(result)

		require.NotNil(t, parsedURL)
		assert.Equal(t, "/path/to/resource", parsedURL.Path)
	})

	t.Run("fails to decode invalid URL", func(t *testing.T) {
		result := urlCodec.Decode("not a valid url ://")

		assert.True(t, either.IsLeft(result), "should fail to decode invalid URL")

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(*url.URL) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		assert.NotEmpty(t, errors)
	})

	t.Run("fails to decode URL with invalid characters", func(t *testing.T) {
		result := urlCodec.Decode("http://example.com/path with spaces")

		// Note: url.Parse actually handles spaces, so let's test a truly invalid URL
		result = urlCodec.Decode("ht!tp://invalid")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes URL to string", func(t *testing.T) {
		parsedURL, err := url.Parse("https://example.com/path?query=value")
		require.NoError(t, err)

		encoded := urlCodec.Encode(parsedURL)

		assert.Equal(t, "https://example.com/path?query=value", encoded)
	})

	t.Run("encodes URL with fragment", func(t *testing.T) {
		parsedURL, err := url.Parse("https://example.com/page#section")
		require.NoError(t, err)

		encoded := urlCodec.Encode(parsedURL)

		assert.Equal(t, "https://example.com/page#section", encoded)
	})

	t.Run("round-trip encoding and decoding", func(t *testing.T) {
		original := "https://example.com/path?key=value&foo=bar#fragment"

		// Decode
		decodeResult := urlCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		parsedURL := getOrElseNull(decodeResult)

		// Encode
		encoded := urlCodec.Encode(parsedURL)

		assert.Equal(t, original, encoded)
	})

	t.Run("codec has correct name", func(t *testing.T) {
		assert.Equal(t, "URL", urlCodec.Name())
	})
}

func TestDate(t *testing.T) {

	getOrElseNull := either.GetOrElse(reader.Of[validation.Errors, time.Time](time.Time{}))

	t.Run("ISO 8601 date format", func(t *testing.T) {
		dateCodec := Date("2006-01-02")

		t.Run("decodes valid date", func(t *testing.T) {
			result := dateCodec.Decode("2024-03-15")

			assert.True(t, either.IsRight(result))

			parsedDate := getOrElseNull(result)

			assert.Equal(t, 2024, parsedDate.Year())
			assert.Equal(t, time.March, parsedDate.Month())
			assert.Equal(t, 15, parsedDate.Day())
		})

		t.Run("fails to decode invalid date format", func(t *testing.T) {
			result := dateCodec.Decode("15-03-2024")

			assert.True(t, either.IsLeft(result))

			errors := either.MonadFold(result,
				F.Identity[validation.Errors],
				func(time.Time) validation.Errors { return nil },
			)

			require.NotNil(t, errors)
			assert.NotEmpty(t, errors)
		})

		t.Run("fails to decode invalid date", func(t *testing.T) {
			result := dateCodec.Decode("2024-13-45")

			assert.True(t, either.IsLeft(result))
		})

		t.Run("fails to decode non-date string", func(t *testing.T) {
			result := dateCodec.Decode("not a date")

			assert.True(t, either.IsLeft(result))
		})

		t.Run("encodes date to string", func(t *testing.T) {
			date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

			encoded := dateCodec.Encode(date)

			assert.Equal(t, "2024-03-15", encoded)
		})

		t.Run("round-trip encoding and decoding", func(t *testing.T) {
			original := "2024-12-25"

			// Decode
			decodeResult := dateCodec.Decode(original)
			require.True(t, either.IsRight(decodeResult))

			parsedDate := getOrElseNull(decodeResult)

			// Encode
			encoded := dateCodec.Encode(parsedDate)

			assert.Equal(t, original, encoded)
		})
	})

	t.Run("RFC3339 timestamp format", func(t *testing.T) {
		timestampCodec := Date(time.RFC3339)

		t.Run("decodes valid RFC3339 timestamp", func(t *testing.T) {
			result := timestampCodec.Decode("2024-03-15T10:30:00Z")

			assert.True(t, either.IsRight(result))

			parsedTime := getOrElseNull(result)

			assert.Equal(t, 2024, parsedTime.Year())
			assert.Equal(t, time.March, parsedTime.Month())
			assert.Equal(t, 15, parsedTime.Day())
			assert.Equal(t, 10, parsedTime.Hour())
			assert.Equal(t, 30, parsedTime.Minute())
			assert.Equal(t, 0, parsedTime.Second())
		})

		t.Run("decodes RFC3339 with timezone offset", func(t *testing.T) {
			result := timestampCodec.Decode("2024-03-15T10:30:00+01:00")

			assert.True(t, either.IsRight(result))

			parsedTime := getOrElseNull(result)

			assert.Equal(t, 2024, parsedTime.Year())
			assert.Equal(t, time.March, parsedTime.Month())
			assert.Equal(t, 15, parsedTime.Day())
		})

		t.Run("fails to decode invalid RFC3339", func(t *testing.T) {
			result := timestampCodec.Decode("2024-03-15 10:30:00")

			assert.True(t, either.IsLeft(result))
		})

		t.Run("encodes timestamp to RFC3339 string", func(t *testing.T) {
			timestamp := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)

			encoded := timestampCodec.Encode(timestamp)

			assert.Equal(t, "2024-03-15T10:30:00Z", encoded)
		})

		t.Run("round-trip encoding and decoding", func(t *testing.T) {
			original := "2024-12-25T15:45:30Z"

			// Decode
			decodeResult := timestampCodec.Decode(original)
			require.True(t, either.IsRight(decodeResult))

			parsedTime := getOrElseNull(decodeResult)

			// Encode
			encoded := timestampCodec.Encode(parsedTime)

			assert.Equal(t, original, encoded)
		})
	})

	t.Run("custom date format", func(t *testing.T) {
		customCodec := Date("02/01/2006")

		t.Run("decodes custom format", func(t *testing.T) {
			result := customCodec.Decode("15/03/2024")

			assert.True(t, either.IsRight(result))

			parsedDate := getOrElseNull(result)

			assert.Equal(t, 2024, parsedDate.Year())
			assert.Equal(t, time.March, parsedDate.Month())
			assert.Equal(t, 15, parsedDate.Day())
		})

		t.Run("encodes to custom format", func(t *testing.T) {
			date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

			encoded := customCodec.Encode(date)

			assert.Equal(t, "15/03/2024", encoded)
		})
	})

	t.Run("codec has correct name", func(t *testing.T) {
		dateCodec := Date("2006-01-02")
		assert.Equal(t, "Date", dateCodec.Name())
	})
}

func TestValidateFromParser(t *testing.T) {
	t.Run("successful parsing", func(t *testing.T) {
		// Create a simple parser that always succeeds
		parser := func(s string) (int, error) {
			return 42, nil
		}

		validator := validateFromParser(parser)
		decode := validator("test")

		// Execute with empty context
		result := decode(validation.Context{})

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 42, value)
	})

	t.Run("failed parsing", func(t *testing.T) {
		// Create a parser that always fails
		parser := func(s string) (int, error) {
			return 0, assert.AnError
		}

		validator := validateFromParser(parser)
		decode := validator("test")

		// Execute with empty context
		result := decode(validation.Context{})

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		assert.NotEmpty(t, errors)

		// Check that the error contains the input value
		if len(errors) > 0 {
			assert.Equal(t, "test", errors[0].Value)
		}
	})

	t.Run("parser with context", func(t *testing.T) {
		parser := func(s string) (string, error) {
			if s == "" {
				return "", assert.AnError
			}
			return s, nil
		}

		validator := validateFromParser(parser)

		// Test with context
		ctx := validation.Context{
			{Key: "field", Type: "string"},
		}

		decode := validator("")
		result := decode(ctx)

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(string) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		assert.NotEmpty(t, errors)

		// Verify context is preserved
		if len(errors) > 0 {
			assert.Equal(t, ctx, errors[0].Context)
		}
	})
}

func TestRegex(t *testing.T) {

	getOrElseNull := either.GetOrElse(reader.Of[validation.Errors](prism.Match{}))

	t.Run("simple number pattern", func(t *testing.T) {
		numberRegex := regexp.MustCompile(`\d+`)
		regexCodec := Regex(numberRegex)

		t.Run("decodes string with number", func(t *testing.T) {
			result := regexCodec.Decode("Price: 42 dollars")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "Price: ", match.Before)
			assert.Equal(t, []string{"42"}, match.Groups)
			assert.Equal(t, " dollars", match.After)
		})

		t.Run("decodes number at start", func(t *testing.T) {
			result := regexCodec.Decode("123 items")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "", match.Before)
			assert.Equal(t, []string{"123"}, match.Groups)
			assert.Equal(t, " items", match.After)
		})

		t.Run("decodes number at end", func(t *testing.T) {
			result := regexCodec.Decode("Total: 999")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "Total: ", match.Before)
			assert.Equal(t, []string{"999"}, match.Groups)
			assert.Equal(t, "", match.After)
		})

		t.Run("fails to decode string without number", func(t *testing.T) {
			result := regexCodec.Decode("no numbers here")

			assert.True(t, either.IsLeft(result))

			errors := either.MonadFold(result,
				F.Identity[validation.Errors],
				func(prism.Match) validation.Errors { return nil },
			)

			require.NotNil(t, errors)
			assert.NotEmpty(t, errors)
		})

		t.Run("encodes Match to string", func(t *testing.T) {
			match := prism.Match{
				Before: "Price: ",
				Groups: []string{"42"},
				After:  " dollars",
			}

			encoded := regexCodec.Encode(match)

			assert.Equal(t, "Price: 42 dollars", encoded)
		})

		t.Run("round-trip encoding and decoding", func(t *testing.T) {
			original := "Count: 789 items"

			// Decode
			decodeResult := regexCodec.Decode(original)
			require.True(t, either.IsRight(decodeResult))

			match := getOrElseNull(decodeResult)

			// Encode
			encoded := regexCodec.Encode(match)

			assert.Equal(t, original, encoded)
		})
	})

	t.Run("pattern with capture groups", func(t *testing.T) {
		// Pattern to match word followed by number
		wordNumberRegex := regexp.MustCompile(`(\w+)(\d+)`)
		regexCodec := Regex(wordNumberRegex)

		t.Run("decodes with capture groups", func(t *testing.T) {
			result := regexCodec.Decode("item42")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "", match.Before)
			// Groups contains the full match and capture groups
			require.NotEmpty(t, match.Groups)
			assert.Equal(t, "item42", match.Groups[0])
			// Verify we have capture groups
			if len(match.Groups) > 1 {
				assert.Contains(t, match.Groups[1], "item")
				assert.Contains(t, match.Groups[len(match.Groups)-1], "2")
			}
			assert.Equal(t, "", match.After)
		})
	})

	t.Run("codec name contains pattern info", func(t *testing.T) {
		numberRegex := regexp.MustCompile(`\d+`)
		regexCodec := Regex(numberRegex)
		// The name is generated by FromRefinement and includes the pattern
		assert.Contains(t, regexCodec.Name(), "FromRefinement")
	})
}

func TestRegexNamed(t *testing.T) {

	getOrElseNull := either.GetOrElse(reader.Of[validation.Errors](prism.NamedMatch{}))

	t.Run("email pattern with named groups", func(t *testing.T) {
		emailRegex := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		emailCodec := RegexNamed(emailRegex)

		t.Run("decodes valid email", func(t *testing.T) {
			result := emailCodec.Decode("john@example.com")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "", match.Before)
			assert.Equal(t, "john@example.com", match.Full)
			assert.Equal(t, "", match.After)
			require.NotNil(t, match.Groups)
			assert.Equal(t, "john", match.Groups["user"])
			assert.Equal(t, "example.com", match.Groups["domain"])
		})

		t.Run("decodes email with surrounding text", func(t *testing.T) {
			result := emailCodec.Decode("Contact: alice@test.org for info")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "Contact: ", match.Before)
			assert.Equal(t, "alice@test.org", match.Full)
			assert.Equal(t, " for info", match.After)
			assert.Equal(t, "alice", match.Groups["user"])
			assert.Equal(t, "test.org", match.Groups["domain"])
		})

		t.Run("fails to decode invalid email", func(t *testing.T) {
			result := emailCodec.Decode("not-an-email")

			assert.True(t, either.IsLeft(result))

			errors := either.MonadFold(result,
				F.Identity[validation.Errors],
				func(prism.NamedMatch) validation.Errors { return nil },
			)

			require.NotNil(t, errors)
			assert.NotEmpty(t, errors)
		})

		t.Run("encodes NamedMatch to string", func(t *testing.T) {
			match := prism.NamedMatch{
				Before: "Email: ",
				Groups: map[string]string{"user": "bob", "domain": "example.com"},
				Full:   "bob@example.com",
				After:  "",
			}

			encoded := emailCodec.Encode(match)

			assert.Equal(t, "Email: bob@example.com", encoded)
		})

		t.Run("round-trip encoding and decoding", func(t *testing.T) {
			original := "Contact: support@company.io"

			// Decode
			decodeResult := emailCodec.Decode(original)
			require.True(t, either.IsRight(decodeResult))

			match := getOrElseNull(decodeResult)

			// Encode
			encoded := emailCodec.Encode(match)

			assert.Equal(t, original, encoded)
		})
	})

	t.Run("phone pattern with named groups", func(t *testing.T) {
		phoneRegex := regexp.MustCompile(`(?P<area>\d{3})-(?P<prefix>\d{3})-(?P<line>\d{4})`)
		phoneCodec := RegexNamed(phoneRegex)

		t.Run("decodes valid phone number", func(t *testing.T) {
			result := phoneCodec.Decode("555-123-4567")

			assert.True(t, either.IsRight(result))

			match := getOrElseNull(result)

			assert.Equal(t, "555-123-4567", match.Full)
			assert.Equal(t, "555", match.Groups["area"])
			assert.Equal(t, "123", match.Groups["prefix"])
			assert.Equal(t, "4567", match.Groups["line"])
		})

		t.Run("fails to decode invalid phone format", func(t *testing.T) {
			result := phoneCodec.Decode("123-45-6789")

			assert.True(t, either.IsLeft(result))
		})
	})

	t.Run("codec name contains refinement info", func(t *testing.T) {
		emailRegex := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+\.\w+)`)
		emailCodec := RegexNamed(emailRegex)
		// The name is generated by FromRefinement
		assert.Contains(t, emailCodec.Name(), "FromRefinement")
	})
}

func TestIntFromString(t *testing.T) {
	intCodec := IntFromString()

	t.Run("decodes positive integer", func(t *testing.T) {
		result := intCodec.Decode("42")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 42, value)
	})

	t.Run("decodes negative integer", func(t *testing.T) {
		result := intCodec.Decode("-123")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, -123, value)
	})

	t.Run("decodes zero", func(t *testing.T) {
		result := intCodec.Decode("0")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int { return -1 },
			F.Identity[int],
		)

		assert.Equal(t, 0, value)
	})

	t.Run("decodes integer with plus sign", func(t *testing.T) {
		result := intCodec.Decode("+456")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 456, value)
	})

	t.Run("fails to decode floating point", func(t *testing.T) {
		result := intCodec.Decode("3.14")

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		assert.NotEmpty(t, errors)
	})

	t.Run("fails to decode non-numeric string", func(t *testing.T) {
		result := intCodec.Decode("not a number")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode empty string", func(t *testing.T) {
		result := intCodec.Decode("")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode hexadecimal", func(t *testing.T) {
		result := intCodec.Decode("0xFF")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode with whitespace", func(t *testing.T) {
		result := intCodec.Decode(" 42 ")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes positive integer", func(t *testing.T) {
		encoded := intCodec.Encode(42)

		assert.Equal(t, "42", encoded)
	})

	t.Run("encodes negative integer", func(t *testing.T) {
		encoded := intCodec.Encode(-123)

		assert.Equal(t, "-123", encoded)
	})

	t.Run("encodes zero", func(t *testing.T) {
		encoded := intCodec.Encode(0)

		assert.Equal(t, "0", encoded)
	})

	t.Run("round-trip encoding and decoding", func(t *testing.T) {
		original := "9876"

		// Decode
		decodeResult := intCodec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		value := either.MonadFold(decodeResult,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		// Encode
		encoded := intCodec.Encode(value)

		assert.Equal(t, original, encoded)
	})

	t.Run("codec has correct name", func(t *testing.T) {
		assert.Equal(t, "IntFromString", intCodec.Name())
	})
}

func TestInt64FromString(t *testing.T) {
	int64Codec := Int64FromString()

	t.Run("decodes positive int64", func(t *testing.T) {
		result := int64Codec.Decode("9223372036854775807")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int64 { return 0 },
			F.Identity[int64],
		)

		assert.Equal(t, int64(9223372036854775807), value)
	})

	t.Run("decodes negative int64", func(t *testing.T) {
		result := int64Codec.Decode("-9223372036854775808")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int64 { return 0 },
			F.Identity[int64],
		)

		assert.Equal(t, int64(-9223372036854775808), value)
	})

	t.Run("decodes zero", func(t *testing.T) {
		result := int64Codec.Decode("0")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int64 { return -1 },
			F.Identity[int64],
		)

		assert.Equal(t, int64(0), value)
	})

	t.Run("decodes small int64", func(t *testing.T) {
		result := int64Codec.Decode("42")

		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) int64 { return 0 },
			F.Identity[int64],
		)

		assert.Equal(t, int64(42), value)
	})

	t.Run("fails to decode out of range positive", func(t *testing.T) {
		result := int64Codec.Decode("9223372036854775808")

		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int64) validation.Errors { return nil },
		)

		require.NotNil(t, errors)
		assert.NotEmpty(t, errors)
	})

	t.Run("fails to decode out of range negative", func(t *testing.T) {
		result := int64Codec.Decode("-9223372036854775809")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode floating point", func(t *testing.T) {
		result := int64Codec.Decode("3.14")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode non-numeric string", func(t *testing.T) {
		result := int64Codec.Decode("not a number")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode empty string", func(t *testing.T) {
		result := int64Codec.Decode("")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes positive int64", func(t *testing.T) {
		encoded := int64Codec.Encode(9223372036854775807)

		assert.Equal(t, "9223372036854775807", encoded)
	})

	t.Run("encodes negative int64", func(t *testing.T) {
		encoded := int64Codec.Encode(-9223372036854775808)

		assert.Equal(t, "-9223372036854775808", encoded)
	})

	t.Run("encodes zero", func(t *testing.T) {
		encoded := int64Codec.Encode(0)

		assert.Equal(t, "0", encoded)
	})

	t.Run("encodes small int64", func(t *testing.T) {
		encoded := int64Codec.Encode(42)

		assert.Equal(t, "42", encoded)
	})

	t.Run("round-trip encoding and decoding", func(t *testing.T) {
		original := "1234567890123456"

		// Decode
		decodeResult := int64Codec.Decode(original)
		require.True(t, either.IsRight(decodeResult))

		value := either.MonadFold(decodeResult,
			func(validation.Errors) int64 { return 0 },
			F.Identity[int64],
		)

		// Encode
		encoded := int64Codec.Encode(value)

		assert.Equal(t, original, encoded)
	})

	t.Run("codec has correct name", func(t *testing.T) {
		assert.Equal(t, "Int64FromString", int64Codec.Name())
	})
}
