package builder

import (
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakePersonValidate_ValidPerson tests validation of a valid person
func TestMakePersonValidate_ValidPerson(t *testing.T) {
	// Arrange
	validate := makePersonValidate()
	builder := MakePerson("Alice", 25)
	ctx := A.Of(validation.ContextEntry{Type: "Person", Actual: builder})

	// Act
	result := validate(builder)(ctx)

	// Assert
	assert.True(t, either.IsRight(result), "Expected validation to succeed")

	person, _ := either.Unwrap(result)
	require.NotNil(t, person, "Expected to unwrap person")

	assert.Equal(t, NonEmptyString("Alice"), person.Name)
	assert.Equal(t, AdultAge(25), person.Age)
}

// TestMakePersonValidate_EmptyName tests validation failure for empty name
func TestMakePersonValidate_EmptyName(t *testing.T) {
	// Arrange
	validate := makePersonValidate()
	builder := MakePerson("", 25)
	ctx := A.Of(validation.ContextEntry{Type: "Person", Actual: builder})

	// Act
	result := validate(builder)(ctx)

	// Assert
	assert.True(t, either.IsLeft(result), "Expected validation to fail for empty name")

	_, errors := either.Unwrap(result)
	assert.NotEmpty(t, errors, "Expected validation errors")
}

// TestMakePersonValidate_InvalidAge tests validation failure for age < 18
func TestMakePersonValidate_InvalidAge(t *testing.T) {
	// Arrange
	validate := makePersonValidate()
	builder := MakePerson("Bob", 15)
	ctx := A.Of(validation.ContextEntry{Type: "Person", Actual: builder})

	// Act
	result := validate(builder)(ctx)

	// Assert
	assert.True(t, either.IsLeft(result), "Expected validation to fail for age < 18")

	_, errors := either.Unwrap(result)
	assert.NotEmpty(t, errors, "Expected validation errors")
}

// TestMakePersonValidate_MultipleErrors tests validation with multiple errors
func TestMakePersonValidate_MultipleErrors(t *testing.T) {
	// Arrange
	validate := makePersonValidate()
	builder := MakePerson("", 10) // Both empty name and invalid age
	ctx := A.Of(validation.ContextEntry{Type: "Person", Actual: builder})

	// Act
	result := validate(builder)(ctx)

	// Assert
	assert.True(t, either.IsLeft(result), "Expected validation to fail")

	_, errors := either.Unwrap(result)
	assert.Len(t, errors, 2, "Expected two validation errors")
}

// TestMakePersonValidate_BoundaryAge tests validation at age boundary (18)
func TestMakePersonValidate_BoundaryAge(t *testing.T) {
	// Arrange
	validate := makePersonValidate()
	builder := MakePerson("Charlie", 18)
	ctx := A.Of(validation.ContextEntry{Type: "Person", Actual: builder})

	// Act
	result := validate(builder)(ctx)

	// Assert
	assert.True(t, either.IsRight(result), "Expected validation to succeed for age 18")

	person, _ := either.Unwrap(result)
	require.NotNil(t, person, "Expected to unwrap person")
	assert.Equal(t, AdultAge(18), person.Age)
}

// TestMakePersonCodec_Decode tests the codec's Decode method
func TestMakePersonCodec_Decode(t *testing.T) {
	// Arrange
	codec := makePersonCodec()
	builder := MakePerson("Diana", 30)

	// Act
	result := codec.Decode(builder)

	// Assert
	assert.True(t, either.IsRight(result), "Expected decode to succeed")

	person, _ := either.Unwrap(result)
	require.NotNil(t, person, "Expected to unwrap person")

	assert.Equal(t, NonEmptyString("Diana"), person.Name)
	assert.Equal(t, AdultAge(30), person.Age)
}

// TestMakePersonCodec_Decode_Invalid tests the codec's Decode method with invalid data
func TestMakePersonCodec_Decode_Invalid(t *testing.T) {
	// Arrange
	codec := makePersonCodec()
	builder := MakePerson("", 10) // Invalid name and age

	// Act
	result := codec.Decode(builder)

	// Assert
	assert.True(t, either.IsLeft(result), "Expected decode to fail")

	_, errors := either.Unwrap(result)
	assert.Len(t, errors, 2, "Expected two validation errors")
}

// TestMakePersonCodec_Encode tests the codec's Encode method
func TestMakePersonCodec_Encode(t *testing.T) {
	// Arrange
	codec := makePersonCodec()
	person := &Person{
		Name: NonEmptyString("Eve"),
		Age:  AdultAge(28),
	}

	// Act
	builder := codec.Encode(person)

	// Apply the builder to get a PartialPerson
	partial := builder(emptyPartialPerson)

	// Assert
	assert.Equal(t, "Eve", partial.name)
	assert.Equal(t, 28, partial.age)
}

// TestMakePersonCodec_RoundTrip tests encoding and decoding round-trip
func TestMakePersonCodec_RoundTrip(t *testing.T) {
	// Arrange
	codec := makePersonCodec()
	originalPerson := &Person{
		Name: NonEmptyString("Frank"),
		Age:  AdultAge(35),
	}

	// Act - Encode to builder
	builder := codec.Encode(originalPerson)

	// Decode back to person
	result := codec.Decode(builder)

	// Assert
	assert.True(t, either.IsRight(result), "Expected round-trip to succeed")

	decodedPerson, _ := either.Unwrap(result)
	require.NotNil(t, decodedPerson, "Expected to unwrap person")

	assert.Equal(t, originalPerson.Name, decodedPerson.Name)
	assert.Equal(t, originalPerson.Age, decodedPerson.Age)
}

// TestMakePersonCodec_Name tests the codec's Name method
func TestMakePersonCodec_Name(t *testing.T) {
	// Arrange
	codec := makePersonCodec()

	// Act
	name := codec.Name()

	// Assert
	assert.Equal(t, "Person", name)
}

// TestNameCodec_Validate tests the name codec validation
func TestNameCodec_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "valid name",
			input:     "Alice",
			wantValid: true,
		},
		{
			name:      "empty name",
			input:     "",
			wantValid: false,
		},
		{
			name:      "whitespace name",
			input:     "   ",
			wantValid: true, // Non-empty string, even if whitespace
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := A.Of(validation.ContextEntry{Type: "Name", Actual: tt.input})

			// Act
			result := nameCodec.Validate(tt.input)(ctx)

			// Assert
			if tt.wantValid {
				assert.True(t, either.IsRight(result), "Expected validation to succeed")
			} else {
				assert.True(t, either.IsLeft(result), "Expected validation to fail")
			}
		})
	}
}

// TestAgeCodec_Validate tests the age codec validation
func TestAgeCodec_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		wantValid bool
	}{
		{
			name:      "valid adult age",
			input:     25,
			wantValid: true,
		},
		{
			name:      "boundary age 18",
			input:     18,
			wantValid: true,
		},
		{
			name:      "minor age",
			input:     17,
			wantValid: false,
		},
		{
			name:      "zero age",
			input:     0,
			wantValid: false,
		},
		{
			name:      "negative age",
			input:     -5,
			wantValid: false,
		},
		{
			name:      "very old age",
			input:     120,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := A.Of(validation.ContextEntry{Type: "Age", Actual: tt.input})

			// Act
			result := ageCodec.Validate(tt.input)(ctx)

			// Assert
			if tt.wantValid {
				assert.True(t, either.IsRight(result), "Expected validation to succeed")
			} else {
				assert.True(t, either.IsLeft(result), "Expected validation to fail")
			}
		})
	}
}

// TestMakePersonCodec_WithComposedBuilders tests codec with composed builders
func TestMakePersonCodec_WithComposedBuilders(t *testing.T) {
	// Arrange
	codec := makePersonCodec()

	// Create a builder by composing individual field setters
	builder := endomorphism.Chain(
		WithAge(40),
	)(WithName("Grace"))

	// Act
	result := codec.Decode(builder)

	// Assert
	assert.True(t, either.IsRight(result), "Expected decode to succeed")

	person, _ := either.Unwrap(result)
	require.NotNil(t, person, "Expected to unwrap person")

	assert.Equal(t, NonEmptyString("Grace"), person.Name)
	assert.Equal(t, AdultAge(40), person.Age)
}

// TestMakePersonCodec_PartialBuilder tests codec with partial builder (missing fields)
func TestMakePersonCodec_PartialBuilder(t *testing.T) {
	// Arrange
	codec := makePersonCodec()

	// Create a builder that only sets name
	builder := WithName("Henry")

	// Act
	result := codec.Decode(builder)

	// Assert
	// Should fail because age is 0 (< 18)
	assert.True(t, either.IsLeft(result), "Expected decode to fail for missing age")
}
