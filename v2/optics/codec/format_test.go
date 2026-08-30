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
		codec := String()
		result := codec.String()
		assert.Equal(t, "string", result)
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int()
		result := codec.String()
		assert.Equal(t, "int", result)
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool()
		result := codec.String()
		assert.Equal(t, "bool", result)
	})
}

// TestTypeImplFormat tests the Format() method implementation
func TestTypeImplFormat(t *testing.T) {
	t.Run("String codec with %s", func(t *testing.T) {
		codec := String()
		result := fmt.Sprintf("%s", codec)
		assert.Equal(t, "string", result)
	})

	t.Run("String codec with %v", func(t *testing.T) {
		codec := String()
		result := fmt.Sprintf("%v", codec)
		assert.Equal(t, "string", result)
	})

	t.Run("String codec with %q", func(t *testing.T) {
		codec := String()
		result := fmt.Sprintf("%q", codec)
		assert.Equal(t, `"string"`, result)
	})

	t.Run("Int codec with %s", func(t *testing.T) {
		codec := Int()
		result := fmt.Sprintf("%s", codec)
		assert.Equal(t, "int", result)
	})

	t.Run("Int codec with %#v", func(t *testing.T) {
		codec := Int()
		result := fmt.Sprintf("%#v", codec)
		assert.Equal(t, "int", result)
	})
}

// TestTypeImplFormatWithPrintf tests that %#v uses String
func TestTypeImplFormatWithPrintf(t *testing.T) {
	stringCodec := String()

	// %#v falls back to String()
	result := fmt.Sprintf("%#v", stringCodec)
	assert.Equal(t, "string", result)
}

// TestComplexTypeFormatting tests formatting of more complex types
func TestComplexTypeFormatting(t *testing.T) {
	// Create an array codec
	arrayCodec := Array(Int())

	// Test String()
	name := arrayCodec.String()
	assert.Equal(t, "Array[int]", name)

	// Test Format with %s
	formatted := fmt.Sprintf("%s", arrayCodec)
	assert.Equal(t, "Array[int]", formatted)

	// Test %#v falls back to String
	assert.Equal(t, "Array[int]", fmt.Sprintf("%#v", arrayCodec))
}

// TestFormatterInterface verifies that Type implements fmt.Formatter
func TestFormatterInterface(t *testing.T) {
	var _ fmt.Formatter = Type[int, int, any]{}
}

// TestStringerInterface verifies that Type implements fmt.Stringer
func TestStringerInterface(t *testing.T) {
	var _ fmt.Stringer = Type[int, int, any]{}
}

// TestLogValuerInterface verifies that Type implements slog.LogValuer
func TestLogValuerInterface(t *testing.T) {
	var _ slog.LogValuer = Type[int, int, any]{}
}

// TestTypeImplLogValue tests the LogValue() method implementation
func TestTypeImplLogValue(t *testing.T) {
	t.Run("String codec", func(t *testing.T) {
		codec := String()
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "string", logValue.String())
	})

	t.Run("Int codec", func(t *testing.T) {
		codec := Int()
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "int", logValue.String())
	})

	t.Run("Bool codec", func(t *testing.T) {
		codec := Bool()
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "bool", logValue.String())
	})

	t.Run("Array codec", func(t *testing.T) {
		codec := Array(Int())
		logValue := codec.LogValue()
		assert.Equal(t, slog.KindString, logValue.Kind())
		assert.Equal(t, "Array[int]", logValue.String())
	})
}

// TestFormattableInterface verifies that Type implements formatting.Formattable
func TestFormattableInterface(t *testing.T) {
	var _ Formattable = Type[int, int, any]{}
}
