package codec

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestTypeToPrismBasic tests basic TypeToPrism functionality
func TestTypeToPrismBasic(t *testing.T) {
	// Create a simple string identity type
	stringType := Id[string]()

	prism := TypeToPrism(stringType)

	t.Run("GetOption returns Some for valid value", func(t *testing.T) {
		result := prism.GetOption("hello")
		assert.True(t, option.IsSome(result), "Expected Some for valid string")

		value := option.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "hello", value)
	})

	t.Run("ReverseGet encodes value correctly", func(t *testing.T) {
		encoded := prism.ReverseGet("world")
		assert.Equal(t, "world", encoded)
	})

	t.Run("Name is preserved from Type", func(t *testing.T) {
		assert.Equal(t, stringType.Name(), prism.String())
	})

	t.Run("Round trip preserves value", func(t *testing.T) {
		original := "test value"
		encoded := prism.ReverseGet(original)
		decoded := prism.GetOption(encoded)

		assert.True(t, option.IsSome(decoded))
		value := option.GetOrElse(F.Constant(""))(decoded)
		assert.Equal(t, original, value)
	})
}

// TestTypeToPrismValidationLogic tests TypeToPrism with validation logic
func TestTypeToPrismValidationLogic(t *testing.T) {
	// Create a type that validates positive integers
	positiveIntType := MakeType(
		"PositiveInt",
		func(u any) either.Either[error, int] {
			i, ok := u.(int)
			if !ok || i <= 0 {
				return either.Left[int](assert.AnError)
			}
			return either.Of[error](i)
		},
		func(i int) Decode[Context, int] {
			return func(c Context) Validation[int] {
				if i <= 0 {
					return validation.FailureWithMessage[int](i, "must be positive")(c)
				}
				return validation.Success(i)
			}
		},
		F.Identity[int],
	)

	prism := TypeToPrism(positiveIntType)

	t.Run("GetOption returns Some for valid positive integer", func(t *testing.T) {
		result := prism.GetOption(42)
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 42, value)
	})

	t.Run("GetOption returns None for negative integer", func(t *testing.T) {
		result := prism.GetOption(-5)
		assert.True(t, option.IsNone(result), "Expected None for negative integer")
	})

	t.Run("GetOption returns None for zero", func(t *testing.T) {
		result := prism.GetOption(0)
		assert.True(t, option.IsNone(result), "Expected None for zero")
	})

	t.Run("GetOption returns Some for boundary value", func(t *testing.T) {
		result := prism.GetOption(1)
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 1, value)
	})

	t.Run("ReverseGet does not validate", func(t *testing.T) {
		// ReverseGet should encode without validation
		encoded := prism.ReverseGet(-10)
		assert.Equal(t, -10, encoded, "ReverseGet should not validate")
	})

	t.Run("Name reflects validation purpose", func(t *testing.T) {
		assert.Equal(t, "PositiveInt", prism.String())
	})
}

// TestTypeToPrismWithComplexValidation tests more complex validation scenarios
func TestTypeToPrismWithComplexValidation(t *testing.T) {
	// Create a type that validates strings with length constraints
	boundedStringType := MakeType(
		"BoundedString",
		func(u any) either.Either[error, string] {
			s, ok := u.(string)
			if !ok {
				return either.Left[string](assert.AnError)
			}
			return either.Of[error](s)
		},
		func(s string) Decode[Context, string] {
			return func(c Context) Validation[string] {
				if len(s) < 3 {
					return validation.FailureWithMessage[string](s, "must be at least 3 characters")(c)
				}
				if len(s) > 10 {
					return validation.FailureWithMessage[string](s, "must be at most 10 characters")(c)
				}
				return validation.Success(s)
			}
		},
		F.Identity[string],
	)

	prism := TypeToPrism(boundedStringType)

	t.Run("GetOption returns Some for valid length", func(t *testing.T) {
		result := prism.GetOption("hello")
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "hello", value)
	})

	t.Run("GetOption returns None for too short string", func(t *testing.T) {
		result := prism.GetOption("ab")
		assert.True(t, option.IsNone(result))
	})

	t.Run("GetOption returns None for too long string", func(t *testing.T) {
		result := prism.GetOption("this is way too long")
		assert.True(t, option.IsNone(result))
	})

	t.Run("GetOption returns Some for minimum length", func(t *testing.T) {
		result := prism.GetOption("abc")
		assert.True(t, option.IsSome(result))
	})

	t.Run("GetOption returns Some for maximum length", func(t *testing.T) {
		result := prism.GetOption("1234567890")
		assert.True(t, option.IsSome(result))
	})
}

// TestTypeToPrismWithNumericTypes tests TypeToPrism with different numeric types
func TestTypeToPrismWithNumericTypes(t *testing.T) {
	t.Run("Float64 type", func(t *testing.T) {
		floatType := Id[float64]()

		prism := TypeToPrism(floatType)

		result := prism.GetOption(3.14)
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(0.0))(result)
		assert.Equal(t, 3.14, value)
	})

	t.Run("Int64 type", func(t *testing.T) {
		int64Type := Id[int64]()

		prism := TypeToPrism(int64Type)

		result := prism.GetOption(int64(9223372036854775807))
		assert.True(t, option.IsSome(result))
	})
}

// TestTypeToPrismWithBooleanType tests TypeToPrism with boolean type
func TestTypeToPrismWithBooleanType(t *testing.T) {
	boolType := Id[bool]()

	prism := TypeToPrism(boolType)

	t.Run("GetOption returns Some for true", func(t *testing.T) {
		result := prism.GetOption(true)
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(false))(result)
		assert.True(t, value)
	})

	t.Run("GetOption returns Some for false", func(t *testing.T) {
		result := prism.GetOption(false)
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant(true))(result)
		assert.False(t, value)
	})

	t.Run("ReverseGet preserves boolean values", func(t *testing.T) {
		assert.True(t, prism.ReverseGet(true))
		assert.False(t, prism.ReverseGet(false))
	})
}

// TestTypeToPrismEdgeCases tests edge cases and special scenarios
func TestTypeToPrismEdgeCases(t *testing.T) {
	t.Run("Empty string validation", func(t *testing.T) {
		nonEmptyStringType := MakeType(
			"NonEmptyString",
			func(u any) either.Either[error, string] {
				s, ok := u.(string)
				if !ok {
					return either.Left[string](assert.AnError)
				}
				return either.Of[error](s)
			},
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if s == "" {
						return validation.FailureWithMessage[string](s, "must not be empty")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		prism := TypeToPrism(nonEmptyStringType)

		emptyResult := prism.GetOption("")
		assert.True(t, option.IsNone(emptyResult), "Empty string should fail validation")

		nonEmptyResult := prism.GetOption("a")
		assert.True(t, option.IsSome(nonEmptyResult))
	})

	t.Run("Multiple validation failures", func(t *testing.T) {
		strictIntType := MakeType(
			"StrictInt",
			func(u any) either.Either[error, int] {
				i, ok := u.(int)
				if !ok {
					return either.Left[int](assert.AnError)
				}
				return either.Of[error](i)
			},
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if i < 0 {
						return validation.FailureWithMessage[int](i, "must be non-negative")(c)
					}
					if i > 100 {
						return validation.FailureWithMessage[int](i, "must be at most 100")(c)
					}
					if i%2 != 0 {
						return validation.FailureWithMessage[int](i, "must be even")(c)
					}
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		prism := TypeToPrism(strictIntType)

		// Valid value
		validResult := prism.GetOption(42)
		assert.True(t, option.IsSome(validResult))

		// Various invalid values
		assert.True(t, option.IsNone(prism.GetOption(-1)), "Negative should fail")
		assert.True(t, option.IsNone(prism.GetOption(101)), "Too large should fail")
		assert.True(t, option.IsNone(prism.GetOption(43)), "Odd should fail")
	})
}

// TestTypeToPrismNamePreservation tests that prism names are correctly preserved
func TestTypeToPrismNamePreservation(t *testing.T) {
	testCases := []struct {
		name     string
		typeName string
	}{
		{"Simple name", "SimpleType"},
		{"Descriptive name", "PositiveIntegerValidator"},
		{"With spaces", "Type With Spaces"},
		{"With special chars", "Type_With-Special.Chars"},
		{"Unicode name", "类型名称"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stringType := MakeType(
				tc.typeName,
				func(u any) either.Either[error, string] {
					s, ok := u.(string)
					if !ok {
						return either.Left[string](assert.AnError)
					}
					return either.Of[error](s)
				},
				func(s string) Decode[Context, string] {
					return func(c Context) Validation[string] {
						return validation.Success(s)
					}
				},
				F.Identity[string],
			)

			prism := TypeToPrism(stringType)
			assert.Equal(t, tc.typeName, prism.String())
		})
	}
}
