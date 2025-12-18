// Copyright (c) 2024 - 2025 IBM Corp.
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

package pair

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("Pair with int and string", func(t *testing.T) {
		p := MakePair("hello", 42)
		result := p.String()
		assert.Equal(t, "Pair[string, int](hello, 42)", result)
	})

	t.Run("Pair with string and string", func(t *testing.T) {
		p := MakePair("key", "value")
		result := p.String()
		assert.Equal(t, "Pair[string, string](key, value)", result)
	})

	t.Run("Pair with error", func(t *testing.T) {
		p := MakePair(errors.New("test error"), 42)
		result := p.String()
		assert.Contains(t, result, "Pair[*errors.errorString, int]")
		assert.Contains(t, result, "test error")
	})

	t.Run("Pair with struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		p := MakePair(User{Name: "Alice", Age: 30}, "active")
		result := p.String()
		assert.Contains(t, result, "Pair")
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
	})
}

func TestGoString(t *testing.T) {
	t.Run("Pair with int and string", func(t *testing.T) {
		p := MakePair("hello", 42)
		result := p.GoString()
		assert.Contains(t, result, "pair.MakePair")
		assert.Contains(t, result, "string")
		assert.Contains(t, result, "int")
		assert.Contains(t, result, "hello")
		assert.Contains(t, result, "42")
	})

	t.Run("Pair with error", func(t *testing.T) {
		p := MakePair(errors.New("test error"), 42)
		result := p.GoString()
		assert.Contains(t, result, "pair.MakePair")
		assert.Contains(t, result, "test error")
	})

	t.Run("Pair with struct", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}
		p := MakePair(TestStruct{Name: "Bob", Age: 25}, 100)
		result := p.GoString()
		assert.Contains(t, result, "pair.MakePair")
		assert.Contains(t, result, "Bob")
		assert.Contains(t, result, "25")
		assert.Contains(t, result, "100")
	})
}

func TestFormatInterface(t *testing.T) {
	t.Run("Pair with %s", func(t *testing.T) {
		p := MakePair("key", 42)
		result := fmt.Sprintf("%s", p)
		assert.Equal(t, "Pair[string, int](key, 42)", result)
	})

	t.Run("Pair with %v", func(t *testing.T) {
		p := MakePair("key", 42)
		result := fmt.Sprintf("%v", p)
		assert.Equal(t, "Pair[string, int](key, 42)", result)
	})

	t.Run("Pair with %+v", func(t *testing.T) {
		p := MakePair("key", 42)
		result := fmt.Sprintf("%+v", p)
		assert.Contains(t, result, "Pair")
		assert.Contains(t, result, "key")
		assert.Contains(t, result, "42")
	})

	t.Run("Pair with %#v (GoString)", func(t *testing.T) {
		p := MakePair("key", 42)
		result := fmt.Sprintf("%#v", p)
		assert.Contains(t, result, "pair.MakePair")
		assert.Contains(t, result, "key")
		assert.Contains(t, result, "42")
	})

	t.Run("Pair with %q", func(t *testing.T) {
		p := MakePair("key", "value")
		result := fmt.Sprintf("%q", p)
		// Should use String() representation
		assert.Contains(t, result, "Pair")
	})

	t.Run("Pair with %T", func(t *testing.T) {
		p := MakePair("key", 42)
		result := fmt.Sprintf("%T", p)
		assert.Contains(t, result, "pair.Pair")
	})
}

func TestLogValue(t *testing.T) {
	t.Run("Pair with int and string", func(t *testing.T) {
		p := MakePair("key", 42)
		logValue := p.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 2)
		assert.Equal(t, "head", attrs[0].Key)
		assert.Equal(t, "key", attrs[0].Value.Any())
		assert.Equal(t, "tail", attrs[1].Key)
		assert.Equal(t, int64(42), attrs[1].Value.Any())
	})

	t.Run("Pair with error", func(t *testing.T) {
		p := MakePair(errors.New("test error"), "value")
		logValue := p.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 2)
		assert.Equal(t, "head", attrs[0].Key)
		assert.NotNil(t, attrs[0].Value.Any())
		assert.Equal(t, "tail", attrs[1].Key)
		assert.Equal(t, "value", attrs[1].Value.Any())
	})

	t.Run("Pair with strings", func(t *testing.T) {
		p := MakePair("first", "second")
		logValue := p.LogValue()

		// Should be a group value
		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract the group attributes
		attrs := logValue.Group()
		assert.Len(t, attrs, 2)
		assert.Equal(t, "head", attrs[0].Key)
		assert.Equal(t, "first", attrs[0].Value.Any())
		assert.Equal(t, "tail", attrs[1].Key)
		assert.Equal(t, "second", attrs[1].Value.Any())
	})

	t.Run("Integration with slog", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		p := MakePair("username", 42)
		logger.Info("test message", "data", p)

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "data")
		assert.Contains(t, output, "head")
		assert.Contains(t, output, "username")
		assert.Contains(t, output, "tail")
		assert.Contains(t, output, "42")
	})
}

func TestFormatComprehensive(t *testing.T) {
	t.Run("All format verbs", func(t *testing.T) {
		p := MakePair("key", 42)

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"Pair", "key", "42"}},
			{"%v", []string{"Pair", "key", "42"}},
			{"%+v", []string{"Pair", "key", "42"}},
			{"%#v", []string{"pair.MakePair", "key", "42"}},
			{"%T", []string{"pair.Pair"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, p)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})

	t.Run("Complex types", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}
		p := MakePair(Config{Host: "localhost", Port: 8080}, []string{"a", "b", "c"})

		tests := []struct {
			verb     string
			contains []string
		}{
			{"%s", []string{"Pair", "localhost", "8080"}},
			{"%v", []string{"Pair", "localhost", "8080"}},
			{"%#v", []string{"pair.MakePair", "localhost", "8080"}},
		}

		for _, tt := range tests {
			t.Run(tt.verb, func(t *testing.T) {
				result := fmt.Sprintf(tt.verb, p)
				for _, substr := range tt.contains {
					assert.Contains(t, result, substr, "Format %s should contain %s", tt.verb, substr)
				}
			})
		}
	})
}

func TestInterfaceImplementations(t *testing.T) {
	t.Run("fmt.Stringer interface", func(t *testing.T) {
		var _ fmt.Stringer = MakePair("key", 42)
	})

	t.Run("fmt.GoStringer interface", func(t *testing.T) {
		var _ fmt.GoStringer = MakePair("key", 42)
	})

	t.Run("fmt.Formatter interface", func(t *testing.T) {
		var _ fmt.Formatter = MakePair("key", 42)
	})

	t.Run("slog.LogValuer interface", func(t *testing.T) {
		var _ slog.LogValuer = MakePair("key", 42)
	})
}
