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

package formatting

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockFormattable is a test implementation of the Formattable interface
type mockFormattable struct {
	stringValue   string
	goStringValue string
}

func (m mockFormattable) String() string {
	return m.stringValue
}

func (m mockFormattable) GoString() string {
	return m.goStringValue
}

func (m mockFormattable) Format(f fmt.State, verb rune) {
	FmtString(m, f, verb)
}

func (m mockFormattable) LogValue() slog.Value {
	return slog.StringValue(m.stringValue)
}

func TestFmtString(t *testing.T) {
	t.Run("format with %v verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%v", mock)
		assert.Equal(t, "test value", result, "Should use String() for %v")
	})

	t.Run("format with %+v verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%+v", mock)
		assert.Equal(t, "test value", result, "Should use String() for %+v")
	})

	t.Run("format with %#v verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%#v", mock)
		assert.Equal(t, "test.GoString", result, "Should use GoString() for %#v")
	})

	t.Run("format with %s verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%s", mock)
		assert.Equal(t, "test value", result, "Should use String() for %s")
	})

	t.Run("format with %q verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%q", mock)
		assert.Equal(t, `"test value"`, result, "Should use quoted String() for %q")
	})

	t.Run("format with unsupported verb", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test value",
			goStringValue: "test.GoString",
		}
		// Using %d which is not a typical string verb
		result := fmt.Sprintf("%d", mock)
		assert.Equal(t, "test value", result, "Should use String() for unsupported verbs")
	})

	t.Run("format with special characters in string", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test\nvalue\twith\rspecial",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%s", mock)
		assert.Equal(t, "test\nvalue\twith\rspecial", result)
	})

	t.Run("format with empty string", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "",
			goStringValue: "",
		}
		result := fmt.Sprintf("%s", mock)
		assert.Equal(t, "", result)
	})

	t.Run("format with unicode characters", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "Hello 世界 🌍",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%s", mock)
		assert.Equal(t, "Hello 世界 🌍", result)
	})

	t.Run("format with %q and special characters", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test\nvalue",
			goStringValue: "test.GoString",
		}
		result := fmt.Sprintf("%q", mock)
		assert.Equal(t, `"test\nvalue"`, result, "Should properly escape special characters in quoted format")
	})
}

func TestTypeInfo(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		tests := []struct {
			name     string
			value    any
			expected string
		}{
			{"int", 42, "int"},
			{"string", "hello", "string"},
			{"bool", true, "bool"},
			{"float64", 3.14, "float64"},
			{"float32", float32(3.14), "float32"},
			{"int64", int64(42), "int64"},
			{"int32", int32(42), "int32"},
			{"uint", uint(42), "uint"},
			{"byte", byte(42), "uint8"},
			{"rune", rune('a'), "int32"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := TypeInfo(tt.value)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("pointer types", func(t *testing.T) {
		var intPtr *int
		result := TypeInfo(intPtr)
		assert.Equal(t, "int", result, "Should remove leading * from pointer type")

		var strPtr *string
		result = TypeInfo(strPtr)
		assert.Equal(t, "string", result, "Should remove leading * from pointer type")
	})

	t.Run("slice types", func(t *testing.T) {
		result := TypeInfo([]int{1, 2, 3})
		assert.Equal(t, "[]int", result)

		result = TypeInfo([]string{"a", "b"})
		assert.Equal(t, "[]string", result)

		result = TypeInfo([][]int{{1, 2}, {3, 4}})
		assert.Equal(t, "[][]int", result)
	})

	t.Run("map types", func(t *testing.T) {
		result := TypeInfo(map[string]int{"a": 1})
		assert.Equal(t, "map[string]int", result)

		result = TypeInfo(map[int]string{1: "a"})
		assert.Equal(t, "map[int]string", result)
	})

	t.Run("struct types", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}
		result := TypeInfo(TestStruct{})
		assert.Equal(t, "formatting.TestStruct", result)

		result = TypeInfo(&TestStruct{})
		assert.Equal(t, "formatting.TestStruct", result, "Should remove leading * from pointer to struct")
	})

	t.Run("interface types", func(t *testing.T) {
		var err error = fmt.Errorf("test error")
		result := TypeInfo(err)
		assert.Contains(t, result, "errors", "Should contain package name")
		assert.NotContains(t, result, "*", "Should not contain pointer prefix")
	})

	t.Run("channel types", func(t *testing.T) {
		ch := make(chan int)
		result := TypeInfo(ch)
		assert.Equal(t, "chan int", result)

		ch2 := make(chan string, 10)
		result = TypeInfo(ch2)
		assert.Equal(t, "chan string", result)
	})

	t.Run("function types", func(t *testing.T) {
		fn := func(int) string { return "" }
		result := TypeInfo(fn)
		assert.Equal(t, "func(int) string", result)
	})

	t.Run("array types", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		result := TypeInfo(arr)
		assert.Equal(t, "[3]int", result)
	})

	t.Run("complex types", func(t *testing.T) {
		type ComplexStruct struct {
			Data map[string][]int
		}
		result := TypeInfo(ComplexStruct{})
		assert.Equal(t, "formatting.ComplexStruct", result)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		result := TypeInfo(ptr)
		assert.Equal(t, "int", result, "Should handle nil pointer correctly")
	})
}

func TestTypeInfoWithCustomTypes(t *testing.T) {
	t.Run("custom type with methods", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "test",
			goStringValue: "test.GoString",
		}
		result := TypeInfo(mock)
		assert.Equal(t, "formatting.mockFormattable", result)
	})

	t.Run("pointer to custom type", func(t *testing.T) {
		mock := &mockFormattable{
			stringValue:   "test",
			goStringValue: "test.GoString",
		}
		result := TypeInfo(mock)
		assert.Equal(t, "formatting.mockFormattable", result, "Should remove pointer prefix")
	})
}

func TestFmtStringIntegration(t *testing.T) {
	t.Run("integration with fmt.Printf", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "integration test",
			goStringValue: "mock.GoString",
		}

		// Test various format combinations
		tests := []struct {
			format   string
			expected string
		}{
			{"%v", "integration test"},
			{"%+v", "integration test"},
			{"%#v", "mock.GoString"},
			{"%s", "integration test"},
			{"%q", `"integration test"`},
		}

		for _, tt := range tests {
			t.Run(tt.format, func(t *testing.T) {
				result := fmt.Sprintf(tt.format, mock)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("integration with fmt.Fprintf", func(t *testing.T) {
		mock := mockFormattable{
			stringValue:   "buffer test",
			goStringValue: "mock.GoString",
		}

		var buf []byte
		n, err := fmt.Fprintf((*mockWriter)(&buf), "%s", mock)
		assert.NoError(t, err)
		assert.Greater(t, n, 0)
		assert.Equal(t, "buffer test", string(buf))
	})
}

// mockWriter is a simple writer for testing fmt.Fprintf
type mockWriter []byte

func (m *mockWriter) Write(p []byte) (n int, err error) {
	*m = append(*m, p...)
	return len(p), nil
}

func BenchmarkFmtString(b *testing.B) {
	mock := mockFormattable{
		stringValue:   "benchmark test value",
		goStringValue: "mock.GoString",
	}

	b.Run("format with %v", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", mock)
		}
	})

	b.Run("format with %#v", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%#v", mock)
		}
	})

	b.Run("format with %s", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s", mock)
		}
	})

	b.Run("format with %q", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%q", mock)
		}
	})
}

func BenchmarkTypeInfo(b *testing.B) {
	values := []any{
		42,
		"string",
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		mockFormattable{},
	}

	for _, v := range values {
		b.Run(fmt.Sprintf("%T", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = TypeInfo(v)
			}
		})
	}
}
