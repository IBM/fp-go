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

package either

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("Right value", func(t *testing.T) {
		e := Right[error](42)
		result := e.String()
		assert.Equal(t, "Right[int](42)", result)
	})

	t.Run("Left value", func(t *testing.T) {
		e := Left[int](errors.New("test error"))
		result := e.String()
		assert.Contains(t, result, "Left[*errors.errorString]")
		assert.Contains(t, result, "test error")
	})

	t.Run("Right with string", func(t *testing.T) {
		e := Right[error]("hello")
		result := e.String()
		assert.Equal(t, "Right[string](hello)", result)
	})

	t.Run("Left with string", func(t *testing.T) {
		e := Left[int]("error message")
		result := e.String()
		assert.Equal(t, "Left[string](error message)", result)
	})

	t.Run("Right zero value", func(t *testing.T) {
		e := Right[string](0)
		result := e.String()
		assert.Equal(t, "Right[int](0)", result)
	})

	t.Run("Left zero value", func(t *testing.T) {
		e := Left[int]("")
		result := e.String()
		assert.Equal(t, "Left[string]()", result)
	})

	t.Run("Right with bool true", func(t *testing.T) {
		e := Right[error](true)
		result := e.String()
		assert.Equal(t, "Right[bool](true)", result)
	})

	t.Run("Right with bool false", func(t *testing.T) {
		e := Right[error](false)
		result := e.String()
		assert.Equal(t, "Right[bool](false)", result)
	})

	t.Run("Right with nil pointer", func(t *testing.T) {
		var p *int
		e := Right[error](p)
		result := e.String()
		assert.Equal(t, "Right[*int](<nil>)", result)
	})

	t.Run("Left with nil pointer", func(t *testing.T) {
		var p *string
		e := Left[int](p)
		result := e.String()
		assert.Equal(t, "Left[*string](<nil>)", result)
	})

	t.Run("String() equals fmt.Sprint()", func(t *testing.T) {
		right := Right[error](42)
		assert.Equal(t, right.String(), fmt.Sprint(right))

		left := Left[int]("err")
		assert.Equal(t, left.String(), fmt.Sprint(left))
	})

	t.Run("String() equals fmt.Sprintf %v", func(t *testing.T) {
		right := Right[error](42)
		assert.Equal(t, right.String(), fmt.Sprintf("%v", right))

		left := Left[int]("err")
		assert.Equal(t, left.String(), fmt.Sprintf("%v", left))
	})
}

func TestFormatInterface(t *testing.T) {
	t.Run("Right value with %s", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%s", e)
		assert.Equal(t, "Right[int](42)", result)
	})

	t.Run("Left value with %s", func(t *testing.T) {
		e := Left[int](errors.New("test error"))
		result := fmt.Sprintf("%s", e)
		assert.Contains(t, result, "Left")
		assert.Contains(t, result, "test error")
	})

	t.Run("Right value with %v", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%v", e)
		assert.Equal(t, "Right[int](42)", result)
	})

	t.Run("Left value with %v", func(t *testing.T) {
		e := Left[int]("error")
		result := fmt.Sprintf("%v", e)
		assert.Equal(t, "Left[string](error)", result)
	})

	t.Run("Right value with %+v", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%+v", e)
		assert.Contains(t, result, "Right")
		assert.Contains(t, result, "42")
	})

	t.Run("Right value with %#v", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%#v", e)
		assert.Equal(t, "Right[int](42)", result)
	})

	t.Run("Left value with %#v", func(t *testing.T) {
		e := Left[int]("error")
		result := fmt.Sprintf("%#v", e)
		assert.Equal(t, "Left[string](error)", result)
	})

	t.Run("Right value with %q", func(t *testing.T) {
		e := Right[error]("hello")
		result := fmt.Sprintf("%q", e)
		// %q quotes the String() representation
		assert.Equal(t, fmt.Sprintf("%q", e.String()), result)
		assert.Contains(t, result, "Right")
	})

	t.Run("Left value with %q", func(t *testing.T) {
		e := Left[int]("error")
		result := fmt.Sprintf("%q", e)
		assert.Equal(t, fmt.Sprintf("%q", e.String()), result)
		assert.Contains(t, result, "Left")
	})

	t.Run("Left value with %+v", func(t *testing.T) {
		e := Left[int]("error")
		result := fmt.Sprintf("%+v", e)
		assert.Equal(t, "Left[string](error)", result)
	})

	t.Run("Right value with %d (default fallthrough)", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%d", e)
		assert.Equal(t, "Right[int](42)", result)
	})

	t.Run("Left value with %x (default fallthrough)", func(t *testing.T) {
		e := Left[int]("err")
		result := fmt.Sprintf("%x", e)
		assert.Equal(t, "Left[string](err)", result)
	})

	t.Run("Right value with %T", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%T", e)
		assert.Contains(t, result, "either.Either")
	})
}

func TestLogValue(t *testing.T) {
	t.Run("Right value", func(t *testing.T) {
		e := Right[error](42)
		logValue := e.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "right", attrs[0].Key)
		assert.Equal(t, int64(42), attrs[0].Value.Any())
	})

	t.Run("Left value", func(t *testing.T) {
		e := Left[int](errors.New("test error"))
		logValue := e.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "left", attrs[0].Key)
		assert.NotNil(t, attrs[0].Value.Any())
	})

	t.Run("Right with string", func(t *testing.T) {
		e := Right[error]("success")
		logValue := e.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "right", attrs[0].Key)
		assert.Equal(t, "success", attrs[0].Value.Any())
	})

	t.Run("Left with string", func(t *testing.T) {
		e := Left[int]("error message")
		logValue := e.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "left", attrs[0].Key)
		assert.Equal(t, "error message", attrs[0].Value.Any())
	})

	t.Run("Integration with slog - Right", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		e := Right[error](42)
		logger.Info("test message", "result", e)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "result")
		assert.Contains(t, output, "right")
		assert.Contains(t, output, "42")
	})

	t.Run("Integration with slog - Left", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		e := Left[int]("error occurred")
		logger.Info("test message", "result", e)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "result")
		assert.Contains(t, output, "left")
		assert.Contains(t, output, "error occurred")
	})

	t.Run("Right with bool", func(t *testing.T) {
		e := Right[error](true)
		logValue := e.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "right", attrs[0].Key)
		assert.Equal(t, true, attrs[0].Value.Any())
	})

	t.Run("Right with float64", func(t *testing.T) {
		e := Right[error](3.14)
		logValue := e.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "right", attrs[0].Key)
		assert.Equal(t, 3.14, attrs[0].Value.Any())
	})

	t.Run("Left with error", func(t *testing.T) {
		err := errors.New("wrapped error")
		e := Left[int](err)
		logValue := e.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "left", attrs[0].Key)
		assert.Equal(t, err, attrs[0].Value.Any())
	})

	t.Run("Left with nil pointer stores nil", func(t *testing.T) {
		var p *string
		e := Left[int](p)
		logValue := e.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "left", attrs[0].Key)
		// nil pointer stored as any becomes (*string)(nil)
		assert.Nil(t, attrs[0].Value.Any())
	})

	t.Run("group always has exactly one attribute", func(t *testing.T) {
		for range 10 {
			right := Right[error](42)
			assert.Len(t, right.LogValue().Group(), 1)

			left := Left[int]("e")
			assert.Len(t, left.LogValue().Group(), 1)
		}
	})
}

func TestFormatComprehensive(t *testing.T) {
	t.Run("All format verbs for Right", func(t *testing.T) {
		e := Right[error](42)

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"Right", "42"}},
			{"%v", []string{"Right", "42"}},
			{"%+v", []string{"Right", "42"}},
			{"%#v", []string{"Right", "42"}},
			{"%T", []string{"either.Either"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, e)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})

	t.Run("All format verbs for Left", func(t *testing.T) {
		e := Left[int]("error")

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"Left", "error"}},
			{"%v", []string{"Left", "error"}},
			{"%+v", []string{"Left", "error"}},
			{"%#v", []string{"Left", "error"}},
			{"%T", []string{"either.Either"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, e)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})
}

func TestInterfaceImplementations(t *testing.T) {
	t.Run("fmt.Stringer interface", func(t *testing.T) {
		var _ fmt.Stringer = Right[error](42)
		var _ fmt.Stringer = Left[int](errors.New("error"))
	})

	t.Run("fmt.Formatter interface", func(t *testing.T) {
		var _ fmt.Formatter = Right[error](42)
		var _ fmt.Formatter = Left[int](errors.New("error"))
	})

	t.Run("slog.LogValuer interface", func(t *testing.T) {
		var _ slog.LogValuer = Right[error](42)
		var _ slog.LogValuer = Left[int](errors.New("error"))
	})
}

// TestStringNestedEither verifies that nested Either values are rendered via
// the inner value's String() representation.
func TestStringNestedEither(t *testing.T) {
	t.Run("Right wrapping Right", func(t *testing.T) {
		inner := Right[error](99)
		outer := Right[error](inner)
		result := outer.String()
		assert.Contains(t, result, "Right")
		// inner.String() appears as the %v of the inner value
		assert.Contains(t, result, inner.String())
	})

	t.Run("Right wrapping Left", func(t *testing.T) {
		inner := Left[int]("inner error")
		outer := Right[error](inner)
		result := outer.String()
		assert.Contains(t, result, "Right")
		assert.Contains(t, result, inner.String())
	})
}

// TestFormatQExact pins the exact output of the %q verb, which quotes the
// String() representation of the Either value.
func TestFormatQExact(t *testing.T) {
	tests := []struct {
		name string
		e    Either[error, string]
	}{
		{"Right", Right[error]("hello")},
		{"Left", Left[string](errors.New("oops"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := fmt.Sprintf("%q", tt.e.String())
			got := fmt.Sprintf("%q", tt.e)
			assert.Equal(t, want, got)
		})
	}
}

// TestLogValueKindIsAlwaysGroup ensures LogValue always returns a group
// regardless of the value type or nil-ness of the payload.
func TestLogValueKindIsAlwaysGroup(t *testing.T) {
	cases := []struct {
		name string
		val  slog.Value
	}{
		{"Right int", Right[error](0).LogValue()},
		{"Right string", Right[error]("").LogValue()},
		{"Right bool", Right[error](false).LogValue()},
		{"Left string", Left[int]("").LogValue()},
		{"Left error", Left[int](errors.New("e")).LogValue()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, slog.KindGroup, tc.val.Kind())
		})
	}
}
