package codec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Run("decodes valid string", func(t *testing.T) {
		stringType := String()
		result := stringType.Decode("hello")

		assert.Equal(t, validation.Of("hello"), result)
	})

	t.Run("fails to decode non-string", func(t *testing.T) {
		stringType := String()
		result := stringType.Decode(123)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode bool as string", func(t *testing.T) {
		stringType := String()
		result := stringType.Decode(true)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes string identity", func(t *testing.T) {
		stringType := String()
		encoded := stringType.Encode("world")

		assert.Equal(t, "world", encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		stringType := String()
		assert.Equal(t, "string", stringType.Name())
	})

	t.Run("decodes empty string", func(t *testing.T) {
		stringType := String()
		result := stringType.Decode("")

		assert.Equal(t, validation.Of(""), result)
	})
}

func TestInt(t *testing.T) {
	t.Run("decodes valid int", func(t *testing.T) {
		intType := Int()
		result := intType.Decode(42)

		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("fails to decode string as int", func(t *testing.T) {
		intType := Int()
		result := intType.Decode("42")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode float as int", func(t *testing.T) {
		intType := Int()
		result := intType.Decode(42.5)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes int identity", func(t *testing.T) {
		intType := Int()
		encoded := intType.Encode(100)

		assert.Equal(t, 100, encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		intType := Int()
		assert.Equal(t, "int", intType.Name())
	})

	t.Run("decodes negative int", func(t *testing.T) {
		intType := Int()
		result := intType.Decode(-42)

		assert.Equal(t, validation.Of(-42), result)
	})

	t.Run("decodes zero", func(t *testing.T) {
		intType := Int()
		result := intType.Decode(0)

		assert.Equal(t, validation.Of(0), result)
	})
}

func TestBool(t *testing.T) {
	t.Run("decodes true", func(t *testing.T) {
		boolType := Bool()
		result := boolType.Decode(true)

		assert.Equal(t, validation.Of(true), result)
	})

	t.Run("decodes false", func(t *testing.T) {
		boolType := Bool()
		result := boolType.Decode(false)

		assert.Equal(t, validation.Of(false), result)
	})

	t.Run("fails to decode int as bool", func(t *testing.T) {
		boolType := Bool()
		result := boolType.Decode(1)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails to decode string as bool", func(t *testing.T) {
		boolType := Bool()
		result := boolType.Decode("true")

		assert.True(t, either.IsLeft(result))
	})

	t.Run("encodes bool identity", func(t *testing.T) {
		boolType := Bool()
		encodedTrue := boolType.Encode(true)
		encodedFalse := boolType.Encode(false)

		assert.Equal(t, true, encodedTrue)
		assert.Equal(t, false, encodedFalse)
	})

	t.Run("has correct name", func(t *testing.T) {
		boolType := Bool()
		assert.Equal(t, "bool", boolType.Name())
	})
}

func TestArray(t *testing.T) {
	t.Run("decodes valid int array", func(t *testing.T) {
		intArray := Array(Int())
		result := intArray.Decode([]int{1, 2, 3})

		assert.Equal(t, validation.Of([]int{1, 2, 3}), result)
	})

	t.Run("decodes valid string array", func(t *testing.T) {
		stringArray := Array(String())
		result := stringArray.Decode([]string{"a", "b", "c"})

		assert.Equal(t, validation.Of([]string{"a", "b", "c"}), result)
	})

	t.Run("decodes empty array", func(t *testing.T) {
		intArray := Array(Int())
		result := intArray.Decode([]int{})

		assert.Equal(t, validation.Of([]int{}), result)
	})

	t.Run("fails when array contains invalid element", func(t *testing.T) {
		intArray := Array(Int())
		result := intArray.Decode([]any{1, "invalid", 3})

		assert.True(t, either.IsLeft(result))
	})

	t.Run("decodes string as character array", func(t *testing.T) {
		// When decoding a string with a rune/byte array type
		intArray := Array(Int())
		result := intArray.Decode("hello")

		// String is iterable, so it should try to decode each character
		assert.True(t, either.IsLeft(result)) // Will fail because characters aren't ints
	})

	t.Run("encodes array by mapping encode function", func(t *testing.T) {
		intArray := Array(Int())
		encoded := intArray.Encode([]int{1, 2, 3})

		assert.Equal(t, []int{1, 2, 3}, encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		intArray := Array(Int())
		assert.Equal(t, "Array[int]", intArray.Name())

		stringArray := Array(String())
		assert.Equal(t, "Array[string]", stringArray.Name())
	})

	t.Run("nested arrays", func(t *testing.T) {
		nestedArray := Array(Array(Int()))
		result := nestedArray.Decode([][]int{{1, 2}, {3, 4}})

		assert.Equal(t, validation.Of([][]int{{1, 2}, {3, 4}}), result)
	})

	t.Run("fails to decode non-iterable", func(t *testing.T) {
		intArray := Array(Int())
		result := intArray.Decode(42)

		assert.True(t, either.IsLeft(result))
	})

	t.Run("decodes array of bools", func(t *testing.T) {
		boolArray := Array(Bool())
		result := boolArray.Decode([]bool{true, false, true})

		assert.Equal(t, validation.Of([]bool{true, false, true}), result)
	})

	t.Run("collects multiple validation errors", func(t *testing.T) {
		intArray := Array(Int())
		result := intArray.Decode([]any{1, "bad", 3, "also bad", 5})

		assert.True(t, either.IsLeft(result))
		// Should collect errors for both invalid elements
		errs := either.MonadFold(result,
			F.Identity[validation.Errors],
			func([]int) validation.Errors { return nil },
		)
		require.NotNil(t, errs)
		// Should have at least 2 errors (for "bad" and "also bad")
		assert.GreaterOrEqual(t, len(errs), 2)
	})
}

func TestArrayEncoding(t *testing.T) {
	t.Run("encodes empty array", func(t *testing.T) {
		intArray := Array(Int())
		encoded := intArray.Encode([]int{})

		assert.Equal(t, []int{}, encoded)
	})

	t.Run("encodes single element array", func(t *testing.T) {
		stringArray := Array(String())
		encoded := stringArray.Encode([]string{"single"})

		assert.Equal(t, []string{"single"}, encoded)
	})

	t.Run("encodes nested arrays", func(t *testing.T) {
		nestedArray := Array(Array(String()))
		encoded := nestedArray.Encode([][]string{{"a", "b"}, {"c", "d"}})

		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, encoded)
	})
}

func TestIntegration(t *testing.T) {
	t.Run("string type validates and encodes", func(t *testing.T) {
		stringType := String()

		// Decode
		decoded := stringType.Decode("test")
		assert.True(t, either.IsRight(decoded))

		// Encode
		value := either.MonadFold(decoded,
			func(validation.Errors) string { return "" },
			F.Identity[string],
		)
		encoded := stringType.Encode(value)
		assert.Equal(t, "test", encoded)
	})

	t.Run("array of arrays of ints", func(t *testing.T) {
		arrayType := Array(Array(Int()))

		input := [][]int{{1, 2}, {3, 4, 5}, {6}}
		decoded := arrayType.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) [][]int { return nil },
			F.Identity[[][]int],
		)
		encoded := arrayType.Encode(value)
		assert.Equal(t, input, encoded)
	})
}

func TestTranscodeArray(t *testing.T) {
	t.Run("decodes valid int array from int slice", func(t *testing.T) {
		intTranscode := TranscodeArray(Int())
		result := intTranscode.Decode([]any{1, 2, 3})

		assert.Equal(t, validation.Of([]int{1, 2, 3}), result)
	})

	t.Run("decodes valid string array from string slice", func(t *testing.T) {
		stringTranscode := TranscodeArray(String())
		result := stringTranscode.Decode([]any{"a", "b", "c"})

		assert.Equal(t, validation.Of([]string{"a", "b", "c"}), result)
	})

	t.Run("decodes empty array", func(t *testing.T) {
		intTranscode := TranscodeArray(Int())
		result := intTranscode.Decode([]any{})

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) []int { return nil },
			F.Identity[[]int],
		)
		assert.Empty(t, value)
	})

	t.Run("encodes array by mapping encode function", func(t *testing.T) {
		intTranscode := TranscodeArray(Int())
		encoded := intTranscode.Encode([]int{1, 2, 3})

		assert.Equal(t, []int{1, 2, 3}, encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		intTranscode := TranscodeArray(Int())
		assert.Equal(t, "Array[int]", intTranscode.Name())

		stringTranscode := TranscodeArray(String())
		assert.Equal(t, "Array[string]", stringTranscode.Name())
	})

	t.Run("nested arrays", func(t *testing.T) {
		nestedTranscode := TranscodeArray(TranscodeArray(Int()))
		result := nestedTranscode.Decode([][]any{{1, 2}, {3, 4}})

		assert.Equal(t, validation.Of([][]int{{1, 2}, {3, 4}}), result)
	})

	t.Run("decodes array of bools", func(t *testing.T) {
		boolTranscode := TranscodeArray(Bool())
		result := boolTranscode.Decode([]any{true, false, true})

		assert.Equal(t, validation.Of([]bool{true, false, true}), result)
	})

	t.Run("encodes empty array", func(t *testing.T) {
		intTranscode := TranscodeArray(Int())
		encoded := intTranscode.Encode([]int{})

		assert.Equal(t, []int{}, encoded)
	})

	t.Run("encodes single element array", func(t *testing.T) {
		stringTranscode := TranscodeArray(String())
		encoded := stringTranscode.Encode([]string{"single"})

		assert.Equal(t, []string{"single"}, encoded)
	})

	t.Run("encodes nested arrays", func(t *testing.T) {
		nestedTranscode := TranscodeArray(TranscodeArray(String()))
		encoded := nestedTranscode.Encode([][]string{{"a", "b"}, {"c", "d"}})

		assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}}, encoded)
	})
}

func TestTranscodeArrayWithTransformation(t *testing.T) {
	// Create a custom type that transforms strings to ints
	stringToInt := MakeType(
		"StringToInt",
		func(u any) either.Either[error, int] {
			s, ok := u.(string)
			if !ok {
				return either.Left[int](assert.AnError)
			}
			// Simple conversion: length of string
			return either.Of[error](len(s))
		},
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				// Transform string to its length
				return validation.Success(len(s))
			}
		},
		func(i int) string {
			// Encode int back to string representation
			return string(rune('0' + i))
		},
	)

	t.Run("transforms string slice to int slice", func(t *testing.T) {
		arrayTranscode := TranscodeArray(stringToInt)
		result := arrayTranscode.Decode([]string{"a", "bb", "ccc"})

		assert.Equal(t, validation.Of([]int{1, 2, 3}), result)
	})

	t.Run("encodes int slice to string slice", func(t *testing.T) {
		arrayTranscode := TranscodeArray(stringToInt)
		encoded := arrayTranscode.Encode([]int{1, 2, 3})

		assert.Equal(t, []string{"1", "2", "3"}, encoded)
	})

	t.Run("handles empty transformation", func(t *testing.T) {
		arrayTranscode := TranscodeArray(stringToInt)
		result := arrayTranscode.Decode([]string{})

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) []int { return nil },
			F.Identity[[]int],
		)
		assert.Empty(t, value)
	})
}

func TestTranscodeArrayValidation(t *testing.T) {
	// Create a type that only accepts positive integers
	positiveInt := MakeType(
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

	t.Run("validates all elements successfully", func(t *testing.T) {
		arrayTranscode := TranscodeArray(positiveInt)
		result := arrayTranscode.Decode([]int{1, 2, 3, 4, 5})

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) []int { return nil },
			F.Identity[[]int],
		)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, value)
	})

	t.Run("fails when element validation fails", func(t *testing.T) {
		arrayTranscode := TranscodeArray(positiveInt)
		result := arrayTranscode.Decode([]int{1, -2, 3})

		assert.True(t, either.IsLeft(result))
	})

	t.Run("collects multiple validation errors", func(t *testing.T) {
		arrayTranscode := TranscodeArray(positiveInt)
		result := arrayTranscode.Decode([]int{1, -2, 3, -4, 5})

		assert.True(t, either.IsLeft(result))
		errs := either.MonadFold(result,
			F.Identity[validation.Errors],
			func([]int) validation.Errors { return nil },
		)
		require.NotNil(t, errs)
		// Should have at least 2 errors (for -2 and -4)
		assert.GreaterOrEqual(t, len(errs), 2)
	})

	t.Run("fails when all elements are invalid", func(t *testing.T) {
		arrayTranscode := TranscodeArray(positiveInt)
		result := arrayTranscode.Decode([]int{-1, -2, -3})

		assert.True(t, either.IsLeft(result))
		errs := either.MonadFold(result,
			F.Identity[validation.Errors],
			func([]int) validation.Errors { return nil },
		)
		require.NotNil(t, errs)
		assert.GreaterOrEqual(t, len(errs), 3)
	})
}

func TestTranscodeArrayIntegration(t *testing.T) {
	t.Run("round trip with identity codec", func(t *testing.T) {
		stringTranscode := TranscodeArray(String())

		input := []any{"hello", "world", "test"}
		decoded := stringTranscode.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) []string { return nil },
			F.Identity[[]string],
		)
		encoded := stringTranscode.Encode(value)
		assert.Equal(t, []string{"hello", "world", "test"}, encoded)
	})

	t.Run("nested arrays round trip", func(t *testing.T) {
		nestedTranscode := TranscodeArray(TranscodeArray(Int()))

		input := [][]any{{1, 2}, {3, 4, 5}, {6}}
		decoded := nestedTranscode.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) [][]int { return nil },
			F.Identity[[][]int],
		)
		encoded := nestedTranscode.Encode(value)
		assert.Equal(t, [][]int{{1, 2}, {3, 4, 5}, {6}}, encoded)
	})

	t.Run("deeply nested arrays", func(t *testing.T) {
		deeplyNested := TranscodeArray(TranscodeArray(TranscodeArray(Bool())))

		input := [][][]any{
			{{true, false}, {true}},
			{{false, false, true}},
		}
		decoded := deeplyNested.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) [][][]bool { return nil },
			F.Identity[[][][]bool],
		)
		encoded := deeplyNested.Encode(value)
		assert.Equal(t, [][][]bool{
			{{true, false}, {true}},
			{{false, false, true}},
		}, encoded)
	})
}

func TestTranscodeEither(t *testing.T) {
	t.Run("decodes Left value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Left[any, any]("error"))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsLeft(value))
		leftValue := either.MonadFold(value,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "error", leftValue)
	})

	t.Run("decodes Right value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Right[any, any](42))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsRight(value))
		rightValue := either.MonadFold(value,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, rightValue)
	})

	t.Run("encodes Left value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		encoded := eitherCodec.Encode(either.Left[int]("error"))

		assert.True(t, either.IsLeft(encoded))
		leftValue := either.MonadFold(encoded,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "error", leftValue)
	})

	t.Run("encodes Right value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		encoded := eitherCodec.Encode(either.Right[string](42))

		assert.True(t, either.IsRight(encoded))
		rightValue := either.MonadFold(encoded,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, rightValue)
	})

	t.Run("has correct name", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		assert.Equal(t, "Either[string, int]", eitherCodec.Name())
	})

	t.Run("fails when Left value is invalid", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Left[any, any](123)) // Should be string

		assert.True(t, either.IsLeft(result))
	})

	t.Run("fails when Right value is invalid", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Right[any, any]("not an int"))

		assert.True(t, either.IsLeft(result))
	})
}

func TestTranscodeEitherValidation(t *testing.T) {
	t.Run("validates Left value with context", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Left[any, any](123)) // Invalid: should be string

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(either.Either[string, int]) validation.Errors { return nil },
		)
		assert.NotEmpty(t, errors)
		// Verify error contains type information
		assert.Contains(t, fmt.Sprintf("%v", errors[0]), "string")
	})

	t.Run("validates Right value with context", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())
		result := eitherCodec.Decode(either.Right[any, any]("not a number"))

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(either.Either[string, int]) validation.Errors { return nil },
		)
		assert.NotEmpty(t, errors)
		// Verify error contains type information
		assert.Contains(t, fmt.Sprintf("%v", errors[0]), "int")
	})

	t.Run("preserves Either structure on validation failure", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())

		// Left with wrong type
		leftResult := eitherCodec.Decode(either.Left[any, any]([]int{1, 2, 3}))
		assert.True(t, either.IsLeft(leftResult))

		// Right with wrong type
		rightResult := eitherCodec.Decode(either.Right[any, any](true))
		assert.True(t, either.IsLeft(rightResult))
	})

	t.Run("validates with custom codec that can fail", func(t *testing.T) {
		// Create a codec that only accepts positive integers
		positiveInt := MakeType(
			"PositiveInt",
			func(u any) either.Either[error, int] {
				i, ok := u.(int)
				if !ok || i <= 0 {
					return either.Left[int](fmt.Errorf("not a positive integer"))
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

		eitherCodec := TranscodeEither(String(), positiveInt)

		// Valid positive integer
		validResult := eitherCodec.Decode(either.Right[any](42))
		assert.True(t, either.IsRight(validResult))

		// Invalid: zero
		zeroResult := eitherCodec.Decode(either.Right[any](0))
		assert.True(t, either.IsLeft(zeroResult))

		// Invalid: negative
		negativeResult := eitherCodec.Decode(either.Right[any](-5))
		assert.True(t, either.IsLeft(negativeResult))
	})

	t.Run("validates both branches independently", func(t *testing.T) {
		// Create codecs with specific validation rules
		nonEmptyString := MakeType(
			"NonEmptyString",
			func(u any) either.Either[error, string] {
				s, ok := u.(string)
				if !ok || len(s) == 0 {
					return either.Left[string](fmt.Errorf("not a non-empty string"))
				}
				return either.Of[error](s)
			},
			func(s string) Decode[Context, string] {
				return func(c Context) Validation[string] {
					if len(s) == 0 {
						return validation.FailureWithMessage[string](s, "must not be empty")(c)
					}
					return validation.Success(s)
				}
			},
			F.Identity[string],
		)

		evenInt := MakeType(
			"EvenInt",
			func(u any) either.Either[error, int] {
				i, ok := u.(int)
				if !ok || i%2 != 0 {
					return either.Left[int](fmt.Errorf("not an even integer"))
				}
				return either.Of[error](i)
			},
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					if i%2 != 0 {
						return validation.FailureWithMessage[int](i, "must be even")(c)
					}
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		eitherCodec := TranscodeEither(nonEmptyString, evenInt)

		// Valid Left: non-empty string
		validLeft := eitherCodec.Decode(either.Left[int]("hello"))
		assert.True(t, either.IsRight(validLeft))

		// Invalid Left: empty string
		invalidLeft := eitherCodec.Decode(either.Left[int](""))
		assert.True(t, either.IsLeft(invalidLeft))

		// Valid Right: even integer
		validRight := eitherCodec.Decode(either.Right[string](42))
		assert.True(t, either.IsRight(validRight))

		// Invalid Right: odd integer
		invalidRight := eitherCodec.Decode(either.Right[string](43))
		assert.True(t, either.IsLeft(invalidRight))
	})
}

func TestTranscodeEitherWithTransformation(t *testing.T) {
	// Create a codec that transforms strings to their lengths
	stringToLength := MakeType(
		"StringToLength",
		func(u any) either.Either[error, int] {
			s, ok := u.(string)
			if !ok {
				return either.Left[int](assert.AnError)
			}
			return either.Of[error](len(s))
		},
		func(s string) Decode[Context, int] {
			return func(c Context) Validation[int] {
				return validation.Success(len(s))
			}
		},
		func(i int) string {
			return fmt.Sprintf("len=%d", i)
		},
	)

	// Create a codec that doubles integers
	doubleInt := MakeType(
		"DoubleInt",
		func(u any) either.Either[error, int] {
			i, ok := u.(int)
			if !ok {
				return either.Left[int](assert.AnError)
			}
			return either.Of[error](i * 2)
		},
		func(i int) Decode[Context, int] {
			return func(c Context) Validation[int] {
				return validation.Success(i * 2)
			}
		},
		func(i int) int {
			return i / 2
		},
	)

	t.Run("transforms Left value", func(t *testing.T) {
		eitherCodec := TranscodeEither(stringToLength, doubleInt)
		result := eitherCodec.Decode(either.Left[int]("hello"))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[int, int] { return either.Left[int](0) },
			F.Identity[either.Either[int, int]],
		)
		assert.True(t, either.IsLeft(value))
		leftValue := either.MonadFold(value,
			F.Identity[int],
			func(int) int { return 0 },
		)
		assert.Equal(t, 5, leftValue) // "hello" has length 5
	})

	t.Run("transforms Right value", func(t *testing.T) {
		eitherCodec := TranscodeEither(stringToLength, doubleInt)
		result := eitherCodec.Decode(either.Right[string](10))

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[int, int] { return either.Right[int](0) },
			F.Identity[either.Either[int, int]],
		)
		assert.True(t, either.IsRight(value))
		rightValue := either.MonadFold(value,
			func(int) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 20, rightValue) // 10 * 2 = 20
	})

	t.Run("encodes Left value with transformation", func(t *testing.T) {
		eitherCodec := TranscodeEither(stringToLength, doubleInt)
		encoded := eitherCodec.Encode(either.Left[int](5))

		assert.True(t, either.IsLeft(encoded))
		leftValue := either.MonadFold(encoded,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "len=5", leftValue)
	})

	t.Run("encodes Right value with transformation", func(t *testing.T) {
		eitherCodec := TranscodeEither(stringToLength, doubleInt)
		encoded := eitherCodec.Encode(either.Right[int](20))

		assert.True(t, either.IsRight(encoded))
		rightValue := either.MonadFold(encoded,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 10, rightValue) // 20 / 2 = 10
	})
}

func TestTranscodeEitherNested(t *testing.T) {
	t.Run("nested Either values", func(t *testing.T) {
		innerEither := TranscodeEither(String(), Int())
		outerEither := TranscodeEither(Bool(), innerEither)

		// Test Left(bool)
		result := outerEither.Decode(either.Left[either.Either[any, any], any](true))
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[bool, either.Either[string, int]] {
				return either.Left[either.Either[string, int]](false)
			},
			F.Identity[either.Either[bool, either.Either[string, int]]],
		)
		assert.True(t, either.IsLeft(value))
	})

	t.Run("nested Either Right(Left(string))", func(t *testing.T) {
		innerEither := TranscodeEither(String(), Int())
		outerEither := TranscodeEither(Bool(), innerEither)

		result := outerEither.Decode(either.Right[any](either.Left[any, any]("error")))
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[bool, either.Either[string, int]] {
				return either.Left[either.Either[string, int]](false)
			},
			F.Identity[either.Either[bool, either.Either[string, int]]],
		)
		assert.True(t, either.IsRight(value))
		innerValue := either.MonadFold(value,
			func(bool) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsLeft(innerValue))
	})

	t.Run("nested Either Right(Right(int))", func(t *testing.T) {
		innerEither := TranscodeEither(String(), Int())
		outerEither := TranscodeEither(Bool(), innerEither)

		result := outerEither.Decode(either.Right[any](either.Right[any, any](42)))
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) either.Either[bool, either.Either[string, int]] {
				return either.Left[either.Either[string, int]](false)
			},
			F.Identity[either.Either[bool, either.Either[string, int]]],
		)
		assert.True(t, either.IsRight(value))
		innerValue := either.MonadFold(value,
			func(bool) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsRight(innerValue))
		finalValue := either.MonadFold(innerValue,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, finalValue)
	})
}

func TestTranscodeEitherIntegration(t *testing.T) {
	t.Run("round trip with Left value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())

		input := either.Left[any, any]("error message")
		decoded := eitherCodec.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		encoded := eitherCodec.Encode(value)

		assert.True(t, either.IsLeft(encoded))
		leftValue := either.MonadFold(encoded,
			F.Identity[string],
			func(int) string { return "" },
		)
		assert.Equal(t, "error message", leftValue)
	})

	t.Run("round trip with Right value", func(t *testing.T) {
		eitherCodec := TranscodeEither(String(), Int())

		input := either.Right[any, any](123)
		decoded := eitherCodec.Decode(input)
		assert.True(t, either.IsRight(decoded))

		value := either.MonadFold(decoded,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		encoded := eitherCodec.Encode(value)

		assert.True(t, either.IsRight(encoded))
		rightValue := either.MonadFold(encoded,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 123, rightValue)
	})

	t.Run("Either with arrays", func(t *testing.T) {
		stringArray := TranscodeArray(String())
		intArray := TranscodeArray(Int())
		eitherCodec := TranscodeEither(stringArray, intArray)

		// Test Left with array
		leftInput := either.Left[[]any]([]any{"a", "b", "c"})
		leftResult := eitherCodec.Decode(leftInput)
		assert.True(t, either.IsRight(leftResult))

		// Test Right with array
		rightInput := either.Right[[]any]([]any{1, 2, 3})
		rightResult := eitherCodec.Decode(rightInput)
		assert.True(t, either.IsRight(rightResult))
	})
}

func TestTypeToPrism(t *testing.T) {
	// Create a Type[string, string, string] for testing
	stringIdentity := MakeType(
		"StringIdentity",
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

	t.Run("converts Type to Prism with successful decode", func(t *testing.T) {
		prism := TypeToPrism(stringIdentity)

		// Test GetOption with valid value
		result := prism.GetOption("hello")
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "hello", value)
	})

	t.Run("ReverseGet encodes value", func(t *testing.T) {
		prism := TypeToPrism(stringIdentity)

		// Test ReverseGet
		encoded := prism.ReverseGet("world")
		assert.Equal(t, "world", encoded)
	})

	t.Run("preserves Type name", func(t *testing.T) {
		prism := TypeToPrism(stringIdentity)

		assert.Equal(t, "StringIdentity", prism.String())
	})

	t.Run("round trip with valid value", func(t *testing.T) {
		prism := TypeToPrism(stringIdentity)

		// Encode then decode
		original := "test"
		encoded := prism.ReverseGet(original)
		decoded := prism.GetOption(encoded)

		assert.True(t, option.IsSome(decoded))
		value := option.GetOrElse(F.Constant(""))(decoded)
		assert.Equal(t, original, value)
	})
}

func TestTypeToPrismWithValidation(t *testing.T) {
	// Create a type that only accepts positive integers
	positiveInt := MakeType(
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

	prism := TypeToPrism(positiveInt)

	t.Run("GetOption succeeds for valid value", func(t *testing.T) {
		result := prism.GetOption(42)
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(0))(result)
		assert.Equal(t, 42, value)
	})

	t.Run("GetOption fails for invalid value", func(t *testing.T) {
		result := prism.GetOption(-5)
		assert.True(t, option.IsNone(result))
	})

	t.Run("GetOption fails for zero", func(t *testing.T) {
		result := prism.GetOption(0)
		assert.True(t, option.IsNone(result))
	})

	t.Run("ReverseGet works for any value", func(t *testing.T) {
		// ReverseGet doesn't validate, it just encodes
		encoded := prism.ReverseGet(-10)
		assert.Equal(t, -10, encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		assert.Equal(t, "PositiveInt", prism.String())
	})
}

func TestTypeToPrismWithArrays(t *testing.T) {
	// Create identity int type for TranscodeArray
	intIdentity := MakeType(
		"IntIdentity",
		func(u any) either.Either[error, int] {
			i, ok := u.(int)
			if !ok {
				return either.Left[int](assert.AnError)
			}
			return either.Of[error](i)
		},
		func(i int) Decode[Context, int] {
			return func(c Context) Validation[int] {
				return validation.Success(i)
			}
		},
		F.Identity[int],
	)

	intArray := TranscodeArray(intIdentity)
	prism := TypeToPrism(intArray)

	t.Run("GetOption succeeds for valid array", func(t *testing.T) {
		result := prism.GetOption([]int{1, 2, 3})
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant[[]int](nil))(result)
		assert.Equal(t, []int{1, 2, 3}, value)
	})

	t.Run("ReverseGet encodes array", func(t *testing.T) {
		encoded := prism.ReverseGet([]int{10, 20, 30})
		assert.Equal(t, []int{10, 20, 30}, encoded)
	})

	t.Run("GetOption succeeds for empty array", func(t *testing.T) {
		result := prism.GetOption([]int{})
		assert.True(t, option.IsSome(result))
	})
}

func TestTypeToPrismWithEither(t *testing.T) {
	// Create identity types for TranscodeEither
	stringIdentity := MakeType(
		"StringIdentity",
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

	intIdentity := MakeType(
		"IntIdentity",
		func(u any) either.Either[error, int] {
			i, ok := u.(int)
			if !ok {
				return either.Left[int](assert.AnError)
			}
			return either.Of[error](i)
		},
		func(i int) Decode[Context, int] {
			return func(c Context) Validation[int] {
				return validation.Success(i)
			}
		},
		F.Identity[int],
	)

	eitherCodec := TranscodeEither(stringIdentity, intIdentity)
	prism := TypeToPrism(eitherCodec)

	t.Run("GetOption succeeds for Left value", func(t *testing.T) {
		result := prism.GetOption(either.Left[int]("error"))
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(either.Right[string](0)))(result)
		assert.True(t, either.IsLeft(value))
	})

	t.Run("GetOption succeeds for Right value", func(t *testing.T) {
		result := prism.GetOption(either.Right[string](42))
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(either.Left[int]("")))(result)
		assert.True(t, either.IsRight(value))
	})

	t.Run("ReverseGet encodes Left value", func(t *testing.T) {
		encoded := prism.ReverseGet(either.Left[int]("error"))
		assert.True(t, either.IsLeft(encoded))
	})

	t.Run("ReverseGet encodes Right value", func(t *testing.T) {
		encoded := prism.ReverseGet(either.Right[string](100))
		assert.True(t, either.IsRight(encoded))
	})
}

func TestTypeToPrismComposition(t *testing.T) {
	t.Run("compose two Type-based prisms", func(t *testing.T) {
		// Create a string identity type
		stringIdentity := MakeType(
			"StringIdentity",
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

		// Create a prism from the type
		stringPrism := TypeToPrism(stringIdentity)

		// Test basic functionality
		result := stringPrism.GetOption("hello")
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "hello", value)

		// Test ReverseGet
		encoded := stringPrism.ReverseGet("world")
		assert.Equal(t, "world", encoded)
	})
}

func TestTypeToPrismIntegration(t *testing.T) {
	t.Run("nested structures", func(t *testing.T) {
		// Create identity int type
		intIdentity := MakeType(
			"IntIdentity",
			func(u any) either.Either[error, int] {
				i, ok := u.(int)
				if !ok {
					return either.Left[int](assert.AnError)
				}
				return either.Of[error](i)
			},
			func(i int) Decode[Context, int] {
				return func(c Context) Validation[int] {
					return validation.Success(i)
				}
			},
			F.Identity[int],
		)

		// Array of arrays
		nestedArray := TranscodeArray(TranscodeArray(intIdentity))
		prism := TypeToPrism(nestedArray)

		result := prism.GetOption([][]int{{1, 2}, {3, 4}})
		assert.True(t, option.IsSome(result))

		value := option.GetOrElse(F.Constant[[][]int](nil))(result)
		assert.Equal(t, [][]int{{1, 2}, {3, 4}}, value)
	})
}

func TestId(t *testing.T) {
	t.Run("decodes string successfully", func(t *testing.T) {
		idCodec := Id[string]()
		result := idCodec.Decode("hello")

		assert.Equal(t, validation.Of("hello"), result)
	})

	t.Run("decodes int successfully", func(t *testing.T) {
		idCodec := Id[int]()
		result := idCodec.Decode(42)

		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("encodes with identity function", func(t *testing.T) {
		idCodec := Id[string]()
		encoded := idCodec.Encode("world")

		assert.Equal(t, "world", encoded)
	})

	t.Run("has correct name", func(t *testing.T) {
		stringId := Id[string]()
		assert.Equal(t, "string", stringId.Name())

		intId := Id[int]()
		assert.Equal(t, "int", intId.Name())
	})

	t.Run("round trip preserves value", func(t *testing.T) {
		idCodec := Id[int]()
		original := 100

		// Encode
		encoded := idCodec.Encode(original)
		assert.Equal(t, original, encoded)

		// Decode
		decoded := idCodec.Decode(encoded)
		assert.True(t, either.IsRight(decoded))
		value := either.MonadFold(decoded,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, original, value)
	})

	t.Run("works with bool", func(t *testing.T) {
		idCodec := Id[bool]()

		resultTrue := idCodec.Decode(true)
		assert.True(t, either.IsRight(resultTrue))

		resultFalse := idCodec.Decode(false)
		assert.True(t, either.IsRight(resultFalse))
	})

	t.Run("works with struct types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		idCodec := Id[Person]()
		person := Person{Name: "Alice", Age: 30}

		result := idCodec.Decode(person)
		assert.Equal(t, validation.Of(person), result)

		encoded := idCodec.Encode(person)
		assert.Equal(t, person, encoded)
	})
}

func TestIdWithTranscodeArray(t *testing.T) {
	t.Run("Id codec with TranscodeArray", func(t *testing.T) {
		intId := Id[int]()
		arrayCodec := TranscodeArray(intId)

		result := arrayCodec.Decode([]int{1, 2, 3, 4, 5})
		assert.Equal(t, validation.Of([]int{1, 2, 3, 4, 5}), result)
	})

	t.Run("Id codec encodes array with identity", func(t *testing.T) {
		stringId := Id[string]()
		arrayCodec := TranscodeArray(stringId)

		encoded := arrayCodec.Encode([]string{"a", "b", "c"})
		assert.Equal(t, []string{"a", "b", "c"}, encoded)
	})

	t.Run("nested arrays with Id", func(t *testing.T) {
		intId := Id[int]()
		nestedCodec := TranscodeArray(TranscodeArray(intId))

		input := [][]int{{1, 2}, {3, 4}, {5}}
		result := nestedCodec.Decode(input)
		assert.Equal(t, validation.Of(input), result)
	})
}

func TestIdWithTranscodeEither(t *testing.T) {
	t.Run("Id codec with TranscodeEither Left", func(t *testing.T) {
		stringId := Id[string]()
		intId := Id[int]()
		eitherCodec := TranscodeEither(stringId, intId)

		result := eitherCodec.Decode(either.Left[int]("error"))
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Left[int]("") },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsLeft(value))
	})

	t.Run("Id codec with TranscodeEither Right", func(t *testing.T) {
		stringId := Id[string]()
		intId := Id[int]()
		eitherCodec := TranscodeEither(stringId, intId)

		result := eitherCodec.Decode(either.Right[string](42))
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) either.Either[string, int] { return either.Right[string](0) },
			F.Identity[either.Either[string, int]],
		)
		assert.True(t, either.IsRight(value))

		intValue := either.MonadFold(value,
			func(string) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, intValue)
	})
}

func TestIdWithTypeToPrism(t *testing.T) {
	t.Run("Id codec converts to Prism", func(t *testing.T) {
		idCodec := Id[string]()
		prism := TypeToPrism(idCodec)

		// GetOption always succeeds
		result := prism.GetOption("test")
		assert.True(t, option.IsSome(result))
		value := option.GetOrElse(F.Constant(""))(result)
		assert.Equal(t, "test", value)

		// ReverseGet is identity
		encoded := prism.ReverseGet("hello")
		assert.Equal(t, "hello", encoded)
	})

	t.Run("Id prism round trip", func(t *testing.T) {
		idCodec := Id[int]()
		prism := TypeToPrism(idCodec)

		original := 123
		encoded := prism.ReverseGet(original)
		decoded := prism.GetOption(encoded)

		assert.True(t, option.IsSome(decoded))
		value := option.GetOrElse(F.Constant(0))(decoded)
		assert.Equal(t, original, value)
	})
}

func TestIdIntegration(t *testing.T) {
	t.Run("complex nested structure with Id", func(t *testing.T) {
		// Create a complex type: Either[[]string, []int]
		stringId := Id[string]()
		intId := Id[int]()

		stringArray := TranscodeArray(stringId)
		intArray := TranscodeArray(intId)

		eitherCodec := TranscodeEither(stringArray, intArray)

		// Test Left with string array
		leftInput := either.Left[[]int]([]string{"a", "b", "c"})
		leftResult := eitherCodec.Decode(leftInput)
		assert.True(t, either.IsRight(leftResult))

		// Test Right with int array
		rightInput := either.Right[[]string]([]int{1, 2, 3})
		rightResult := eitherCodec.Decode(rightInput)
		assert.True(t, either.IsRight(rightResult))
	})

	t.Run("Id preserves all values without validation", func(t *testing.T) {
		type ComplexStruct struct {
			Name   string
			Age    int
			Active bool
			Tags   []string
		}

		idCodec := Id[ComplexStruct]()

		original := ComplexStruct{
			Name:   "Test",
			Age:    25,
			Active: true,
			Tags:   []string{"tag1", "tag2"},
		}

		// Decode
		decoded := idCodec.Decode(original)
		assert.True(t, either.IsRight(decoded))

		decodedValue := either.MonadFold(decoded,
			func(validation.Errors) ComplexStruct { return ComplexStruct{} },
			F.Identity[ComplexStruct],
		)
		assert.Equal(t, original, decodedValue)

		// Encode
		encoded := idCodec.Encode(original)
		assert.Equal(t, original, encoded)
	})
}

// TestFromRefinement tests the FromRefinement function with basic refinements
func TestFromRefinement(t *testing.T) {
	// Create a refinement for positive integers
	positiveIntPrism := prism.MakePrismWithName(
		func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n)
			}
			return option.None[int]()
		},
		func(n int) int { return n },
		"PositiveInt",
	)

	codec := FromRefinement(positiveIntPrism)

	t.Run("Name", func(t *testing.T) {
		name := codec.Name()
		assert.Equal(t, "FromRefinement(PositiveInt)", name)
	})

	t.Run("DecodeValidPositive", func(t *testing.T) {
		result := codec.Decode(42)
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("DecodeInvalidZero", func(t *testing.T) {
		result := codec.Decode(0)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("DecodeInvalidNegative", func(t *testing.T) {
		result := codec.Decode(-5)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("Encode", func(t *testing.T) {
		encoded := codec.Encode(42)
		assert.Equal(t, 42, encoded)
	})

	t.Run("Is", func(t *testing.T) {
		// Is checks if value is of type int
		result := codec.Is(42)
		assert.True(t, either.IsRight(result))

		result = codec.Is("not an int")
		assert.True(t, either.IsLeft(result))
	})
}

// TestFromRefinementWithStrings tests FromRefinement with string refinements
func TestFromRefinementWithStrings(t *testing.T) {
	// Create a refinement for non-empty strings
	nonEmptyStringPrism := prism.MakePrismWithName(
		func(s string) option.Option[string] {
			if len(s) > 0 {
				return option.Some(s)
			}
			return option.None[string]()
		},
		func(s string) string { return s },
		"NonEmptyString",
	)

	codec := FromRefinement(nonEmptyStringPrism)

	t.Run("DecodeValidNonEmpty", func(t *testing.T) {
		result := codec.Decode("hello")
		assert.Equal(t, validation.Success("hello"), result)
	})

	t.Run("DecodeInvalidEmpty", func(t *testing.T) {
		result := codec.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("Encode", func(t *testing.T) {
		encoded := codec.Encode("world")
		assert.Equal(t, "world", encoded)
	})
}

// TestFromRefinementWithRange tests FromRefinement with bounded numeric types
func TestFromRefinementWithRange(t *testing.T) {
	// Create a refinement for integers in range [0, 100]
	boundedIntPrism := prism.MakePrismWithName(
		func(n int) option.Option[int] {
			if n >= 0 && n <= 100 {
				return option.Some(n)
			}
			return option.None[int]()
		},
		func(n int) int { return n },
		"BoundedInt[0,100]",
	)

	codec := FromRefinement(boundedIntPrism)

	t.Run("DecodeValidLowerBound", func(t *testing.T) {
		result := codec.Decode(0)
		assert.Equal(t, validation.Success(0), result)
	})

	t.Run("DecodeValidUpperBound", func(t *testing.T) {
		result := codec.Decode(100)
		assert.Equal(t, validation.Success(100), result)
	})

	t.Run("DecodeValidMidRange", func(t *testing.T) {
		result := codec.Decode(50)
		assert.Equal(t, validation.Success(50), result)
	})

	t.Run("DecodeInvalidBelowRange", func(t *testing.T) {
		result := codec.Decode(-1)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("DecodeInvalidAboveRange", func(t *testing.T) {
		result := codec.Decode(101)
		assert.True(t, either.IsLeft(result))
	})
}

// TestFromRefinementComposition tests composing FromRefinement with other codecs
func TestFromRefinementComposition(t *testing.T) {
	// Create a refinement for positive integers
	positiveIntPrism := prism.MakePrismWithName(
		func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n)
			}
			return option.None[int]()
		},
		func(n int) int { return n },
		"PositiveInt",
	)

	positiveCodec := FromRefinement(positiveIntPrism)

	// Compose with Int codec using Pipe
	composed := Pipe[int, any](positiveCodec)(Int())

	t.Run("ComposedDecodeValid", func(t *testing.T) {
		result := composed.Decode(42)
		assert.Equal(t, validation.Success(42), result)
	})

	t.Run("ComposedDecodeInvalidType", func(t *testing.T) {
		result := composed.Decode("not an int")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("ComposedDecodeInvalidRefinement", func(t *testing.T) {
		result := composed.Decode(-5)
		assert.True(t, either.IsLeft(result))
	})

	t.Run("ComposedEncode", func(t *testing.T) {
		encoded := composed.Encode(42)
		assert.Equal(t, 42, encoded)
	})
}

// TestFromRefinementWithCustomTypes tests FromRefinement with custom struct types
func TestFromRefinementWithCustomTypes(t *testing.T) {
	type Email struct {
		Value string
	}

	// Create a refinement that validates email format (simplified)
	emailPrism := prism.MakePrismWithName(
		func(s string) option.Option[Email] {
			// Simplified email validation
			if len(s) > 0 && strings.Contains(s, "@") && strings.Contains(s, ".") {
				return option.Some(Email{Value: s})
			}
			return option.None[Email]()
		},
		func(e Email) string { return e.Value },
		"Email",
	)

	codec := FromRefinement(emailPrism)

	t.Run("DecodeValidEmail", func(t *testing.T) {
		result := codec.Decode("user@example.com")
		assert.Equal(t, validation.Success(Email{Value: "user@example.com"}), result)
	})

	t.Run("DecodeInvalidEmailNoAt", func(t *testing.T) {
		result := codec.Decode("userexample.com")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("DecodeInvalidEmailNoDot", func(t *testing.T) {
		result := codec.Decode("user@examplecom")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("DecodeInvalidEmpty", func(t *testing.T) {
		result := codec.Decode("")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("Encode", func(t *testing.T) {
		email := Email{Value: "user@example.com"}
		encoded := codec.Encode(email)
		assert.Equal(t, "user@example.com", encoded)
	})
}

// TestFromRefinementValidationContext tests that validation context is properly maintained
func TestFromRefinementValidationContext(t *testing.T) {
	positiveIntPrism := prism.MakePrismWithName(
		func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n)
			}
			return option.None[int]()
		},
		func(n int) int { return n },
		"PositiveInt",
	)

	codec := FromRefinement(positiveIntPrism)

	t.Run("ValidationErrorContainsContext", func(t *testing.T) {
		result := codec.Decode(-5)
		assert.True(t, either.IsLeft(result))

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(int) validation.Errors { return nil },
		)
		require.NotEmpty(t, errors)

		// Check that error contains the value and context
		err := errors[0]
		assert.Equal(t, -5, err.Value)
	})
}

// TestEmpty_Success tests that Empty always succeeds during decoding
func TestEmpty_Success(t *testing.T) {
	t.Run("decodes any input to default value", func(t *testing.T) {
		defaultCodec := Empty[int, string, any](lazy.Of(pair.MakePair("default", 42)))

		// Test with various input types
		testCases := []struct {
			name  string
			input any
		}{
			{"string input", "anything"},
			{"int input", 123},
			{"nil input", nil},
			{"bool input", true},
			{"struct input", struct{ X int }{X: 10}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := defaultCodec.Decode(tc.input)

				assert.Equal(t, validation.Of(42), result)
			})
		}
	})

	t.Run("always returns same default value", func(t *testing.T) {
		defaultCodec := Empty[string, string, any](lazy.Of(pair.MakePair("output", "default")))

		result1 := defaultCodec.Decode(123)
		result2 := defaultCodec.Decode("different")
		result3 := defaultCodec.Decode(nil)

		assert.True(t, either.IsRight(result1))
		assert.True(t, either.IsRight(result2))
		assert.True(t, either.IsRight(result3))

		value1 := either.MonadFold(result1, func(validation.Errors) string { return "" }, F.Identity[string])
		value2 := either.MonadFold(result2, func(validation.Errors) string { return "" }, F.Identity[string])
		value3 := either.MonadFold(result3, func(validation.Errors) string { return "" }, F.Identity[string])

		assert.Equal(t, "default", value1)
		assert.Equal(t, "default", value2)
		assert.Equal(t, "default", value3)
	})
}

// TestEmpty_Encoding tests that Empty always uses default output during encoding
func TestEmpty_Encoding(t *testing.T) {
	t.Run("encodes any value to default output", func(t *testing.T) {
		defaultCodec := Empty[int, string, any](lazy.Of(pair.MakePair("default", 42)))

		// Test with various input values
		testCases := []struct {
			name  string
			input int
		}{
			{"zero value", 0},
			{"positive value", 100},
			{"negative value", -50},
			{"default value", 42},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				encoded := defaultCodec.Encode(tc.input)
				assert.Equal(t, "default", encoded)
			})
		}
	})

	t.Run("always returns same default output", func(t *testing.T) {
		defaultCodec := Empty[string, int, any](lazy.Of(pair.MakePair(999, "ignored")))

		encoded1 := defaultCodec.Encode("value1")
		encoded2 := defaultCodec.Encode("value2")
		encoded3 := defaultCodec.Encode("")

		assert.Equal(t, 999, encoded1)
		assert.Equal(t, 999, encoded2)
		assert.Equal(t, 999, encoded3)
	})
}

// TestEmpty_Name tests that Empty has correct name
func TestEmpty_Name(t *testing.T) {
	t.Run("has name 'Empty'", func(t *testing.T) {
		defaultCodec := Empty[int, int, any](lazy.Of(pair.MakePair(0, 0)))
		assert.Equal(t, "Empty", defaultCodec.Name())
	})
}

// TestEmpty_TypeChecking tests that Empty performs standard type checking
func TestEmpty_TypeChecking(t *testing.T) {
	t.Run("Is checks for correct type", func(t *testing.T) {
		defaultCodec := Empty[int, string, any](lazy.Of(pair.MakePair("default", 42)))

		// Should succeed for int
		result := defaultCodec.Is(100)
		assert.True(t, either.IsRight(result))

		// Should fail for non-int
		result = defaultCodec.Is("not an int")
		assert.True(t, either.IsLeft(result))
	})

	t.Run("Is checks for string type", func(t *testing.T) {
		defaultCodec := Empty[string, string, any](lazy.Of(pair.MakePair("out", "in")))

		// Should succeed for string
		result := defaultCodec.Is("hello")
		assert.True(t, either.IsRight(result))

		// Should fail for non-string
		result = defaultCodec.Is(123)
		assert.True(t, either.IsLeft(result))
	})
}

// TestEmpty_LazyEvaluation tests that the Pair parameter allows dynamic values
func TestEmpty_LazyEvaluation(t *testing.T) {
	t.Run("lazy pair allows dynamic values", func(t *testing.T) {
		counter := 0
		lazyPair := func() pair.Pair[int, int] {
			counter++
			return pair.MakePair(counter, counter*10)
		}

		defaultCodec := Empty[int, int, any](lazyPair)

		// Each decode can get a different value if the lazy function is dynamic
		result1 := defaultCodec.Decode("input1")
		value1 := either.MonadFold(result1,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		result2 := defaultCodec.Decode("input2")
		value2 := either.MonadFold(result2,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		// Values can be different if lazy function produces different results
		assert.True(t, value1 > 0)
		assert.True(t, value2 > 0)
	})
}

// TestEmpty_WithStructs tests Empty with struct types
func TestEmpty_WithStructs(t *testing.T) {
	type Config struct {
		Timeout int
		Retries int
	}

	t.Run("provides default struct value", func(t *testing.T) {
		defaultConfig := Config{Timeout: 30, Retries: 3}
		defaultCodec := Empty[Config, Config, any](lazy.Of(pair.MakePair(defaultConfig, defaultConfig)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) Config { return Config{} },
			F.Identity[Config],
		)
		assert.Equal(t, 30, value.Timeout)
		assert.Equal(t, 3, value.Retries)
	})

	t.Run("encodes to default struct", func(t *testing.T) {
		defaultConfig := Config{Timeout: 30, Retries: 3}
		inputConfig := Config{Timeout: 60, Retries: 5}

		defaultCodec := Empty[Config, Config, any](lazy.Of(pair.MakePair(defaultConfig, defaultConfig)))

		encoded := defaultCodec.Encode(inputConfig)
		assert.Equal(t, 30, encoded.Timeout)
		assert.Equal(t, 3, encoded.Retries)
	})
}

// TestEmpty_WithPointers tests Empty with pointer types
func TestEmpty_WithPointers(t *testing.T) {
	t.Run("provides default pointer value", func(t *testing.T) {
		defaultValue := 42
		defaultCodec := Empty[*int, *int, any](lazy.Of(pair.MakePair(&defaultValue, &defaultValue)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) *int { return nil },
			F.Identity[*int],
		)
		require.NotNil(t, value)
		assert.Equal(t, 42, *value)
	})

	t.Run("provides nil pointer as default", func(t *testing.T) {
		var nilPtr *int
		defaultCodec := Empty[*int, *int, any](lazy.Of(pair.MakePair(nilPtr, nilPtr)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) *int { return new(int) },
			F.Identity[*int],
		)
		assert.Nil(t, value)
	})
}

// TestEmpty_WithSlices tests Empty with slice types
func TestEmpty_WithSlices(t *testing.T) {
	t.Run("provides default slice value", func(t *testing.T) {
		defaultSlice := []int{1, 2, 3}
		defaultCodec := Empty[[]int, []int, any](lazy.Of(pair.MakePair(defaultSlice, defaultSlice)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) []int { return nil },
			F.Identity[[]int],
		)
		assert.Equal(t, []int{1, 2, 3}, value)
	})

	t.Run("provides empty slice as default", func(t *testing.T) {
		emptySlice := []int{}
		defaultCodec := Empty[[]int, []int, any](lazy.Of(pair.MakePair(emptySlice, emptySlice)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))

		value := either.MonadFold(result,
			func(validation.Errors) []int { return nil },
			F.Identity[[]int],
		)
		assert.Equal(t, []int{}, value)
	})
}

// TestEmpty_DifferentInputOutput tests Empty with different input and output types
func TestEmpty_DifferentInputOutput(t *testing.T) {
	t.Run("decodes to int, encodes to string", func(t *testing.T) {
		defaultCodec := Empty[int, string, any](lazy.Of(pair.MakePair("default-output", 42)))

		// Decode always returns 42
		result := defaultCodec.Decode("any input")
		assert.Equal(t, validation.Of(42), result)

		// Encode always returns "default-output"
		encoded := defaultCodec.Encode(100)
		assert.Equal(t, "default-output", encoded)
	})

	t.Run("decodes to string, encodes to int", func(t *testing.T) {
		defaultCodec := Empty[string, int, any](lazy.Of(pair.MakePair(999, "default-value")))

		// Decode always returns "default-value"
		result := defaultCodec.Decode(123)
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "default-value", value)

		// Encode always returns 999
		encoded := defaultCodec.Encode("any string")
		assert.Equal(t, 999, encoded)
	})
}

// TestEmpty_EdgeCases tests edge cases for Empty
func TestEmpty_EdgeCases(t *testing.T) {
	t.Run("with zero values", func(t *testing.T) {
		defaultCodec := Empty[int, int, any](lazy.Of(pair.MakePair(0, 0)))

		result := defaultCodec.Decode("anything")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) int { return -1 },
			F.Identity[int],
		)
		assert.Equal(t, 0, value)

		encoded := defaultCodec.Encode(100)
		assert.Equal(t, 0, encoded)
	})

	t.Run("with empty string", func(t *testing.T) {
		defaultCodec := Empty[string, string, any](lazy.Of(pair.MakePair("", "")))

		result := defaultCodec.Decode("non-empty")
		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(validation.Errors) string { return "error" },
			F.Identity[string],
		)
		assert.Equal(t, "", value)

		encoded := defaultCodec.Encode("non-empty")
		assert.Equal(t, "", encoded)
	})

	t.Run("with false boolean", func(t *testing.T) {
		defaultCodec := Empty[bool, bool, any](lazy.Of(pair.MakePair(false, false)))

		result := defaultCodec.Decode(true)
		assert.Equal(t, validation.Of(false), result)

		encoded := defaultCodec.Encode(true)
		assert.Equal(t, false, encoded)
	})
}

// TestEmpty_Integration tests Empty in composition scenarios
func TestEmpty_Integration(t *testing.T) {
	t.Run("composes with other codecs using Pipe", func(t *testing.T) {
		// Create a codec that always provides a default int
		defaultIntCodec := Empty[int, int, any](lazy.Of(pair.MakePair(42, 42)))

		// Create a refinement that only accepts positive integers
		positiveIntPrism := prism.MakePrismWithName(
			func(n int) option.Option[int] {
				if n > 0 {
					return option.Some(n)
				}
				return option.None[int]()
			},
			func(n int) int { return n },
			"PositiveInt",
		)

		positiveCodec := FromRefinement(positiveIntPrism)

		// Compose: always decode to 42, then validate it's positive
		composed := Pipe[int, any](positiveCodec)(defaultIntCodec)

		// Should succeed because 42 is positive
		result := composed.Decode("anything")
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("used as placeholder in generic contexts", func(t *testing.T) {
		// Empty can be used where a codec is required but not actually used
		unitCodec := Empty[Void, Void, any](
			lazy.Of(pair.MakePair(F.VOID, F.VOID)),
		)

		result := unitCodec.Decode("ignored")
		assert.Equal(t, validation.Of(F.VOID), result)

		encoded := unitCodec.Encode(F.VOID)
		assert.Equal(t, F.VOID, encoded)
	})
}

// TestEmpty_RoundTrip tests that Empty maintains consistency
func TestEmpty_RoundTrip(t *testing.T) {
	t.Run("decode then encode returns default output", func(t *testing.T) {
		defaultCodec := Empty[int, string, any](lazy.Of(pair.MakePair("output", 42)))

		// Decode
		result := defaultCodec.Decode("input")
		require.True(t, either.IsRight(result))

		decoded := either.MonadFold(result,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)

		// Encode
		encoded := defaultCodec.Encode(decoded)

		// Should get default output, not related to decoded value
		assert.Equal(t, "output", encoded)
	})

	t.Run("multiple round trips are consistent", func(t *testing.T) {
		defaultCodec := Empty[int, int, any](lazy.Of(pair.MakePair(100, 50)))

		// First round trip
		result1 := defaultCodec.Decode("input1")
		decoded1 := either.MonadFold(result1,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)
		encoded1 := defaultCodec.Encode(decoded1)

		// Second round trip
		result2 := defaultCodec.Decode("input2")
		decoded2 := either.MonadFold(result2,
			func(validation.Errors) int { return 0 },
			F.Identity[int],
		)
		encoded2 := defaultCodec.Encode(decoded2)

		// All decoded values should be the same
		assert.Equal(t, 50, decoded1)
		assert.Equal(t, 50, decoded2)

		// All encoded values should be the same
		assert.Equal(t, 100, encoded1)
		assert.Equal(t, 100, encoded2)
	})
}
