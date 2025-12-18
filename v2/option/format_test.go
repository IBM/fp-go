// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("Some value", func(t *testing.T) {
		opt := Some(42)
		result := opt.String()
		assert.Equal(t, "Some[int](42)", result)
	})

	t.Run("None value", func(t *testing.T) {
		opt := None[int]()
		result := opt.String()
		assert.Equal(t, "None[int]", result)
	})

	t.Run("Some with string", func(t *testing.T) {
		opt := Some("hello")
		result := opt.String()
		assert.Equal(t, "Some[string](hello)", result)
	})

	t.Run("None with string", func(t *testing.T) {
		opt := None[string]()
		result := opt.String()
		assert.Equal(t, "None[string]", result)
	})
}

func TestGoString(t *testing.T) {
	t.Run("Some value", func(t *testing.T) {
		opt := Some(42)
		result := opt.GoString()
		assert.Contains(t, result, "option.Some")
		assert.Contains(t, result, "42")
	})

	t.Run("None value", func(t *testing.T) {
		opt := None[int]()
		result := opt.GoString()
		assert.Contains(t, result, "option.None")
		assert.Contains(t, result, "int")
	})

	t.Run("Some with struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}
		opt := Some(TestStruct{Name: "Alice", Age: 30})
		result := opt.GoString()
		assert.Contains(t, result, "option.Some")
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
	})

	t.Run("None with custom type", func(t *testing.T) {
		opt := None[string]()
		result := opt.GoString()
		assert.Contains(t, result, "option.None")
		assert.Contains(t, result, "string")
	})
}

func TestFormatInterface(t *testing.T) {
	t.Run("Some value with %s", func(t *testing.T) {
		opt := Some(42)
		result := fmt.Sprintf("%s", opt)
		assert.Equal(t, "Some[int](42)", result)
	})

	t.Run("None value with %s", func(t *testing.T) {
		opt := None[int]()
		result := fmt.Sprintf("%s", opt)
		assert.Equal(t, "None[int]", result)
	})

	t.Run("Some value with %v", func(t *testing.T) {
		opt := Some(42)
		result := fmt.Sprintf("%v", opt)
		assert.Equal(t, "Some[int](42)", result)
	})

	t.Run("None value with %v", func(t *testing.T) {
		opt := None[string]()
		result := fmt.Sprintf("%v", opt)
		assert.Equal(t, "None[string]", result)
	})

	t.Run("Some value with %+v", func(t *testing.T) {
		opt := Some(42)
		result := fmt.Sprintf("%+v", opt)
		assert.Contains(t, result, "Some")
		assert.Contains(t, result, "42")
	})

	t.Run("Some value with %#v (GoString)", func(t *testing.T) {
		opt := Some(42)
		result := fmt.Sprintf("%#v", opt)
		assert.Contains(t, result, "option.Some")
		assert.Contains(t, result, "42")
	})

	t.Run("None value with %#v (GoString)", func(t *testing.T) {
		opt := None[int]()
		result := fmt.Sprintf("%#v", opt)
		assert.Contains(t, result, "option.None")
		assert.Contains(t, result, "int")
	})

	t.Run("Some value with %q", func(t *testing.T) {
		opt := Some("hello")
		result := fmt.Sprintf("%q", opt)
		assert.Contains(t, result, "Some")
	})

	t.Run("Some value with %T", func(t *testing.T) {
		opt := Some(42)
		result := fmt.Sprintf("%T", opt)
		assert.Contains(t, result, "option.Option")
	})
}

func TestLogValue(t *testing.T) {
	t.Run("Some value", func(t *testing.T) {
		opt := Some(42)
		logValue := opt.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "some", attrs[0].Key)
		assert.Equal(t, int64(42), attrs[0].Value.Any())
	})

	t.Run("None value", func(t *testing.T) {
		opt := None[int]()
		logValue := opt.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "none", attrs[0].Key)
		// Value should be struct{}{}
		assert.Equal(t, struct{}{}, attrs[0].Value.Any())
	})

	t.Run("Some with string", func(t *testing.T) {
		opt := Some("success")
		logValue := opt.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "some", attrs[0].Key)
		assert.Equal(t, "success", attrs[0].Value.Any())
	})

	t.Run("None with string", func(t *testing.T) {
		opt := None[string]()
		logValue := opt.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 1)
		assert.Equal(t, "none", attrs[0].Key)
		assert.Equal(t, struct{}{}, attrs[0].Value.Any())
	})

	t.Run("Integration with slog - Some", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		opt := Some(42)
		logger.Info("test message", "result", opt)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "result")
		assert.Contains(t, output, "some")
		assert.Contains(t, output, "42")
	})

	t.Run("Integration with slog - None", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		opt := None[int]()
		logger.Info("test message", "result", opt)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "result")
		assert.Contains(t, output, "none")
	})
}

func TestFormatComprehensive(t *testing.T) {
	t.Run("All format verbs for Some", func(t *testing.T) {
		opt := Some(42)

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"Some", "42"}},
			{"%v", []string{"Some", "42"}},
			{"%+v", []string{"Some", "42"}},
			{"%#v", []string{"option.Some", "42"}},
			{"%T", []string{"option.Option"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, opt)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})

	t.Run("All format verbs for None", func(t *testing.T) {
		opt := None[int]()

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"None", "int"}},
			{"%v", []string{"None", "int"}},
			{"%+v", []string{"None", "int"}},
			{"%#v", []string{"option.None", "int"}},
			{"%T", []string{"option.Option"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, opt)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})
}

func TestInterfaceImplementations(t *testing.T) {
	t.Run("fmt.Stringer interface", func(t *testing.T) {
		var _ fmt.Stringer = Some(42)
		var _ fmt.Stringer = None[int]()
	})

	t.Run("fmt.GoStringer interface", func(t *testing.T) {
		var _ fmt.GoStringer = Some(42)
		var _ fmt.GoStringer = None[int]()
	})

	t.Run("fmt.Formatter interface", func(t *testing.T) {
		var _ fmt.Formatter = Some(42)
		var _ fmt.Formatter = None[int]()
	})

	t.Run("slog.LogValuer interface", func(t *testing.T) {
		var _ slog.LogValuer = Some(42)
		var _ slog.LogValuer = None[int]()
	})
}
