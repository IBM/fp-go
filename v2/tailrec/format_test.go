package tailrec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStringerInterface verifies fmt.Stringer implementation
func TestStringerInterface(t *testing.T) {
	t.Run("Bounce String", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := tramp.String()
		assert.Equal(t, "Bounce(42)", result)
	})

	t.Run("Land String", func(t *testing.T) {
		tramp := Land[int]("done")
		result := tramp.String()
		assert.Equal(t, "Land(done)", result)
	})

	t.Run("fmt.Sprint uses String", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := fmt.Sprint(tramp)
		assert.Equal(t, "Bounce(42)", result)
	})
}

// TestFormatterInterface verifies fmt.Formatter implementation
func TestFormatterInterface(t *testing.T) {
	t.Run("default %v format", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := fmt.Sprintf("%v", tramp)
		assert.Equal(t, "Bounce(42)", result)
	})

	t.Run("detailed %+v format for Bounce", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := fmt.Sprintf("%+v", tramp)
		assert.Contains(t, result, "Trampoline[Bounce]")
		assert.Contains(t, result, "Bounce: 42")
		assert.Contains(t, result, "Landed: false")
	})

	t.Run("detailed %+v format for Land", func(t *testing.T) {
		tramp := Land[int]("done")
		result := fmt.Sprintf("%+v", tramp)
		assert.Contains(t, result, "Trampoline[Land]")
		assert.Contains(t, result, "Land: done")
		assert.Contains(t, result, "Landed: true")
	})

	t.Run("%#v format delegates to GoString", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := fmt.Sprintf("%#v", tramp)
		assert.Contains(t, result, "tailrec.Bounce")
	})

	t.Run("%s format", func(t *testing.T) {
		tramp := Land[int]("result")
		result := fmt.Sprintf("%s", tramp)
		assert.Equal(t, "Land(result)", result)
	})

	t.Run("%q format", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := fmt.Sprintf("%q", tramp)
		assert.Equal(t, "\"Bounce(42)\"", result)
	})
}

// TestGoStringerInterface verifies fmt.GoStringer implementation
func TestGoStringerInterface(t *testing.T) {
	t.Run("Bounce GoString", func(t *testing.T) {
		tramp := Bounce[string](42)
		result := tramp.GoString()
		assert.Contains(t, result, "tailrec.Bounce")
		assert.Contains(t, result, "string")
		assert.Contains(t, result, "42")
	})

	t.Run("Land GoString", func(t *testing.T) {
		tramp := Land[int]("done")
		result := tramp.GoString()
		assert.Contains(t, result, "tailrec.Land")
		assert.Contains(t, result, "int")
		assert.Contains(t, result, "done")
	})

	t.Run("fmt with %#v uses GoString", func(t *testing.T) {
		tramp := Land[int]("result")
		result := fmt.Sprintf("%#v", tramp)
		assert.Contains(t, result, "tailrec.Land")
	})
}
