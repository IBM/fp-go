package codec

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTypeImplStringer tests the String() method implementation
func TestTypeImplStringer(t *testing.T) {
	t.Run("String codec", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		result := codec.String()
		assert.Equal(t, "string", result)
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		result := codec.String()
		assert.Equal(t, "int", result)
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool().(*typeImpl[bool, bool, any])
		result := codec.String()
		assert.Equal(t, "bool", result)
	})
}

// TestTypeImplFormat tests the Format() method implementation
func TestTypeImplFormat(t *testing.T) {
	t.Run("String codec with %s", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		result := fmt.Sprintf("%s", codec)
		assert.Equal(t, "string", result)
	})

	t.Run("String codec with %v", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		result := fmt.Sprintf("%v", codec)
		assert.Equal(t, "string", result)
	})

	t.Run("String codec with %q", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		result := fmt.Sprintf("%q", codec)
		assert.Equal(t, `"string"`, result)
	})

	t.Run("Int codec with %s", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		result := fmt.Sprintf("%s", codec)
		assert.Equal(t, "int", result)
	})

	t.Run("Int codec with %#v", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		result := fmt.Sprintf("%#v", codec)
		assert.Equal(t, `codec.Type[int, int, interface {}]{name: "int"}`, result)
	})
}

// TestTypeImplGoString tests the GoString() method implementation
func TestTypeImplGoString(t *testing.T) {
	t.Run("String codec", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		result := codec.GoString()
		assert.Equal(t, `codec.Type[string, string, interface {}]{name: "string"}`, result)
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		result := codec.GoString()
		assert.Equal(t, `codec.Type[int, int, interface {}]{name: "int"}`, result)
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool().(*typeImpl[bool, bool, any])
		result := codec.GoString()
		assert.Equal(t, `codec.Type[bool, bool, interface {}]{name: "bool"}`, result)
	})
}

// TestTypeImplFormatWithPrintf tests that %#v uses GoString
func TestTypeImplFormatWithPrintf(t *testing.T) {
	stringCodec := String().(*typeImpl[string, string, any])

	// Test that %#v calls GoString
	result := fmt.Sprintf("%#v", stringCodec)
	assert.Equal(t, `codec.Type[string, string, interface {}]{name: "string"}`, result)
}

// TestComplexTypeFormatting tests formatting of more complex types
func TestComplexTypeFormatting(t *testing.T) {
	// Create an array codec
	arrayCodec := Array(Int()).(*typeImpl[[]int, []int, any])

	// Test String()
	name := arrayCodec.String()
	assert.Equal(t, "Array[int]", name)

	// Test Format with %s
	formatted := fmt.Sprintf("%s", arrayCodec)
	assert.Equal(t, "Array[int]", formatted)

	// Test GoString
	goString := arrayCodec.GoString()
	// Just verify it's not empty
	assert.NotEmpty(t, goString)
}

// TestFormatterInterface verifies that typeImpl implements fmt.Formatter
func TestFormatterInterface(t *testing.T) {
	var _ fmt.Formatter = (*typeImpl[int, int, any])(nil)
}

// TestStringerInterface verifies that typeImpl implements fmt.Stringer
func TestStringerInterface(t *testing.T) {
	var _ fmt.Stringer = (*typeImpl[int, int, any])(nil)
}

// TestGoStringerInterface verifies that typeImpl implements fmt.GoStringer
func TestGoStringerInterface(t *testing.T) {
	var _ fmt.GoStringer = (*typeImpl[int, int, any])(nil)
}

// TestLogValuerInterface verifies that typeImpl implements slog.LogValuer
func TestLogValuerInterface(t *testing.T) {
	var _ slog.LogValuer = (*typeImpl[int, int, any])(nil)
}

// TestTypeImplLogValue tests the LogValue() method implementation
func TestTypeImplLogValue(t *testing.T) {
	t.Run("String codec", func(t *testing.T) {
		codec := String().(*typeImpl[string, string, any])
		logValue := codec.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())

		// Extract attributes from the group
		attrs := logValue.Group()
		assert.Len(t, attrs, 4)

		// Check that we have the expected attributes
		attrMap := make(map[string]string)
		for _, attr := range attrs {
			attrMap[attr.Key] = attr.Value.String()
		}

		assert.Equal(t, "string", attrMap["name"])
		assert.Equal(t, "string", attrMap["type_a"])
		assert.Equal(t, "string", attrMap["type_o"])
		assert.Contains(t, attrMap["type_i"], "interface")
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		logValue := codec.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())

		attrs := logValue.Group()
		assert.Len(t, attrs, 4)

		attrMap := make(map[string]string)
		for _, attr := range attrs {
			attrMap[attr.Key] = attr.Value.String()
		}

		assert.Equal(t, "int", attrMap["name"])
		assert.Equal(t, "int", attrMap["type_a"])
		assert.Equal(t, "int", attrMap["type_o"])
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool().(*typeImpl[bool, bool, any])
		logValue := codec.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())

		attrs := logValue.Group()
		assert.Len(t, attrs, 4)

		attrMap := make(map[string]string)
		for _, attr := range attrs {
			attrMap[attr.Key] = attr.Value.String()
		}

		assert.Equal(t, "bool", attrMap["name"])
		assert.Equal(t, "bool", attrMap["type_a"])
	})

	t.Run("Array codec", func(t *testing.T) {
		codec := Array(Int()).(*typeImpl[[]int, []int, any])
		logValue := codec.LogValue()

		assert.Equal(t, slog.KindGroup, logValue.Kind())

		attrs := logValue.Group()
		assert.Len(t, attrs, 4)

		attrMap := make(map[string]string)
		for _, attr := range attrs {
			attrMap[attr.Key] = attr.Value.String()
		}

		assert.Equal(t, "Array[int]", attrMap["name"])
	})
}

// TestFormattableInterface verifies that typeImpl implements formatting.Formattable
func TestFormattableInterface(t *testing.T) {
	var _ Formattable = (*typeImpl[int, int, any])(nil)
}
