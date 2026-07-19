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
		assert.Equal(t, "int", result)
	})
}

// TestTypeImplFormatWithPrintf tests that %#v uses String
func TestTypeImplFormatWithPrintf(t *testing.T) {
	stringCodec := String().(*typeImpl[string, string, any])

	// %#v falls back to String()
	result := fmt.Sprintf("%#v", stringCodec)
	assert.Equal(t, "string", result)
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

	// Test %#v falls back to String
	assert.Equal(t, "Array[int]", fmt.Sprintf("%#v", arrayCodec))
}

// TestFormatterInterface verifies that typeImpl implements fmt.Formatter
func TestFormatterInterface(t *testing.T) {
	var _ fmt.Formatter = (*typeImpl[int, int, any])(nil)
}

// TestStringerInterface verifies that typeImpl implements fmt.Stringer
func TestStringerInterface(t *testing.T) {
	var _ fmt.Stringer = (*typeImpl[int, int, any])(nil)
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
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "string", logValue.String())
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int().(*typeImpl[int, int, any])
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "int", logValue.String())
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool().(*typeImpl[bool, bool, any])
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "bool", logValue.String())
	})

	t.Run("Array codec", func(t *testing.T) {
		codec := Array(Int()).(*typeImpl[[]int, []int, any])
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "Array[int]", logValue.String())
	})
}

// TestFormattableInterface verifies that typeImpl implements formatting.Formattable
func TestFormattableInterface(t *testing.T) {
	var _ Formattable = (*typeImpl[int, int, any])(nil)
}
