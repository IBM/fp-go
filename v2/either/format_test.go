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
}

func TestGoString(t *testing.T) {
	t.Run("Right value", func(t *testing.T) {
		e := Right[error](42)
		result := e.GoString()
		assert.Contains(t, result, "either.Right")
		assert.Contains(t, result, "42")
	})

	t.Run("Left value", func(t *testing.T) {
		e := Left[int](errors.New("test error"))
		result := e.GoString()
		assert.Contains(t, result, "either.Left")
		assert.Contains(t, result, "test error")
	})

	t.Run("Right with struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}
		e := Right[error](TestStruct{Name: "Alice", Age: 30})
		result := e.GoString()
		assert.Contains(t, result, "either.Right")
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
	})

	t.Run("Left with custom error", func(t *testing.T) {
		e := Left[string]("custom error")
		result := e.GoString()
		assert.Contains(t, result, "either.Left")
		assert.Contains(t, result, "custom error")
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

	t.Run("Right value with %#v (GoString)", func(t *testing.T) {
		e := Right[error](42)
		result := fmt.Sprintf("%#v", e)
		assert.Contains(t, result, "either.Right")
		assert.Contains(t, result, "42")
	})

	t.Run("Left value with %#v (GoString)", func(t *testing.T) {
		e := Left[int]("error")
		result := fmt.Sprintf("%#v", e)
		assert.Contains(t, result, "either.Left")
		assert.Contains(t, result, "error")
	})

	t.Run("Right value with %q", func(t *testing.T) {
		e := Right[error]("hello")
		result := fmt.Sprintf("%q", e)
		// Should use String() representation
		assert.Contains(t, result, "Right")
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
			{"%#v", []string{"either.Right", "42"}},
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
			{"%#v", []string{"either.Left", "error"}},
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

	t.Run("fmt.GoStringer interface", func(t *testing.T) {
		var _ fmt.GoStringer = Right[error](42)
		var _ fmt.GoStringer = Left[int](errors.New("error"))
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
