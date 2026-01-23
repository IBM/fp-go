package codec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/IBM/fp-go/v2/array"
	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	R "github.com/IBM/fp-go/v2/reflect"
	"github.com/IBM/fp-go/v2/result"
)

// typeImpl is the internal implementation of the Type interface.
// It combines encoding, decoding, validation, and type checking capabilities.
type typeImpl[A, O, I any] struct {
	name     string
	is       Reader[any, Result[A]]
	validate Validate[I, A]
	encode   Encode[A, O]
}

var emptyContext = A.Empty[validation.ContextEntry]()

// MakeType creates a new Type with the given name, type checker, validator, and encoder.
//
// Parameters:
//   - name: A descriptive name for this type (used in error messages)
//   - is: A function that checks if a value is of type A
//   - validate: A function that validates and decodes input I to type A
//   - encode: A function that encodes type A to output O
//
// Returns a Type[A, O, I] that can both encode and decode values.
func MakeType[A, O, I any](
	name string,
	is Reader[any, Result[A]],
	validate Validate[I, A],
	encode Encode[A, O],
) Type[A, O, I] {
	return &typeImpl[A, O, I]{
		name:     name,
		is:       is,
		validate: validate,
		encode:   encode,
	}
}

// Validate validates the input value in the context of a validation path.
// Returns a Reader that takes a Context and produces a Validation result.
func (t *typeImpl[A, O, I]) Validate(i I) Decode[Context, A] {
	return t.validate(i)
}

// Decode validates and decodes the input value, creating a new context with this type's name.
// This is a convenience method that calls Validate with a fresh context.
func (t *typeImpl[A, O, I]) Decode(i I) Validation[A] {
	return t.validate(i)(array.Of(validation.ContextEntry{Type: t.name, Actual: i}))
}

// Encode transforms a value of type A into the output format O.
func (t *typeImpl[A, O, I]) Encode(a A) O {
	return t.encode(a)
}

// AsDecoder returns this Type as a Decoder interface.
func (t *typeImpl[A, O, I]) AsDecoder() Decoder[I, A] {
	return t
}

// AsEncoder returns this Type as an Encoder interface.
func (t *typeImpl[A, O, I]) AsEncoder() Encoder[A, O] {
	return t
}

// Name returns the descriptive name of this type.
func (t *typeImpl[A, O, I]) Name() string {
	return t.name
}

func (t *typeImpl[A, O, I]) Is(i any) Result[A] {
	return t.is(i)
}

// Pipe composes two Types, creating a pipeline where:
//   - Decoding: I -> A -> B (decode with 'this', then validate with 'ab')
//   - Encoding: B -> A -> O (encode with 'ab', then encode with 'this')
//
// This allows building complex codecs from simpler ones.
//
// Example:
//
//	stringToInt := codec.MakeType(...)  // Type[int, string, string]
//	intToPositive := codec.MakeType(...) // Type[PositiveInt, int, int]
//	composed := codec.Pipe(intToPositive)(stringToInt) // Type[PositiveInt, string, string]
func Pipe[A, B, O, I any](ab Type[B, A, A]) func(Type[A, O, I]) Type[B, O, I] {
	return func(this Type[A, O, I]) Type[B, O, I] {
		return MakeType(
			fmt.Sprintf("Pipe(%s, %s)", this.Name(), ab.Name()),
			ab.Is,
			F.Flow2(
				this.Validate,
				readereither.Chain(ab.Validate),
			),
			F.Flow2(
				ab.Encode,
				this.Encode,
			),
		)
	}
}

// isNil checks if a value is nil, handling both typed and untyped nil values.
// It uses reflection to detect nil pointers, maps, slices, channels, functions, and interfaces.
func isNil(x any) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// isTypedNil checks if a value is nil and returns it as a typed nil pointer.
// Returns Some(nil) if the value is nil, None otherwise.
func isTypedNil[A any](x any) Result[*A] {
	if isNil(x) {
		return result.Of[*A](nil)
	}
	return result.Left[*A](errors.New("expecting nil"))
}

func validateFromIs[A, I any](
	is ReaderResult[I, A],
	msg string,
) Validate[I, A] {
	return func(i I) Decode[Context, A] {
		return F.Pipe2(
			i,
			is,
			result.Fold(
				validation.FailureWithError[A](F.ToAny(i), msg),
				F.Flow2(
					validation.Success[A],
					reader.Of[Context],
				),
			),
		)
	}
}

func isFromValidate[T, I any](val Validate[I, T]) ReaderResult[any, T] {
	invalidType := result.Left[T](errors.New("invalid input type"))
	return func(u any) Result[T] {
		i, ok := u.(I)
		if !ok {
			return invalidType
		}
		return validation.ToResult(val(i)(emptyContext))
	}
}

// MakeNilType creates a Type that validates nil values.
// It accepts any input and validates that it is nil, returning a typed nil pointer.
//
// Example:
//
//	nilType := codec.MakeNilType[string]()
//	result := nilType.Decode(nil)        // Success: Right((*string)(nil))
//	result := nilType.Decode("not nil")  // Failure: Left(errors)
func Nil[A any]() Type[*A, *A, any] {

	is := isTypedNil[A]

	return MakeType(
		"nil",
		is,
		validateFromIs(is, "nil"),
		F.Identity[*A],
	)
}

func MakeSimpleType[A any]() Type[A, A, any] {
	name := fmt.Sprintf("%T", *new(A))
	is := Is[A]()

	return MakeType(
		name,
		is,
		validateFromIs(is, name),
		F.Identity[A],
	)
}

// String creates a Type for string values.
// It validates that input is a string type and provides identity encoding/decoding.
// This is a simple type that accepts any input and validates it's a string.
//
// Returns:
//   - A Type[string, string, any] that can validate, decode, and encode string values
//
// Example:
//
//	stringType := codec.String()
//	result := stringType.Decode("hello")     // Success: Right("hello")
//	result := stringType.Decode(123)         // Failure: Left(validation errors)
//	encoded := stringType.Encode("world")    // Returns: "world"
func String() Type[string, string, any] {
	return MakeSimpleType[string]()
}

// Int creates a Type for int values.
// It validates that input is an int type and provides identity encoding/decoding.
// This is a simple type that accepts any input and validates it's an int.
//
// Returns:
//   - A Type[int, int, any] that can validate, decode, and encode int values
//
// Example:
//
//	intType := codec.Int()
//	result := intType.Decode(42)         // Success: Right(42)
//	result := intType.Decode("42")       // Failure: Left(validation errors)
//	encoded := intType.Encode(100)       // Returns: 100
func Int() Type[int, int, any] {
	return MakeSimpleType[int]()
}

// Bool creates a Type for bool values.
// It validates that input is a bool type and provides identity encoding/decoding.
// This is a simple type that accepts any input and validates it's a bool.
//
// Returns:
//   - A Type[bool, bool, any] that can validate, decode, and encode bool values
//
// Example:
//
//	boolType := codec.Bool()
//	result := boolType.Decode(true)      // Success: Right(true)
//	result := boolType.Decode(1)         // Failure: Left(validation errors)
//	encoded := boolType.Encode(false)    // Returns: false
func Bool() Type[bool, bool, any] {
	return MakeSimpleType[bool]()
}

func appendContext(key, typ string, actual any) Endomorphism[Context] {
	return A.Push(validation.ContextEntry{Key: key, Type: typ, Actual: actual})
}

type validationPair[T any] = Pair[validation.Errors, T]

func pairToValidation[T any](p validationPair[T]) Validation[T] {
	errors, value := pair.Unpack(p)
	if A.IsNonEmpty(errors) {
		return either.Left[T](errors)
	}
	return either.Of[validation.Errors](value)
}

func validateArrayFromArray[T, O, I any](item Type[T, O, I]) Validate[[]I, []T] {

	appendErrors := F.Flow2(
		A.Concat,
		pair.MapHead[[]T, validation.Errors],
	)

	appendValues := F.Flow2(
		A.Push,
		pair.MapTail[validation.Errors, []T],
	)

	itemName := item.Name()

	zero := pair.Zero[validation.Errors, []T]()

	return func(is []I) Decode[Context, []T] {

		return func(c Context) Validation[[]T] {

			return F.Pipe1(
				A.MonadReduceWithIndex(is, func(i int, p validationPair[[]T], v I) validationPair[[]T] {
					return either.MonadFold(
						item.Validate(v)(appendContext(strconv.Itoa(i), itemName, v)(c)),
						appendErrors,
						appendValues,
					)(p)
				}, zero),
				pairToValidation,
			)
		}
	}
}

func validateArray[T, O any](item Type[T, O, any]) Validate[any, []T] {

	appendErrors := F.Flow2(
		A.Concat,
		pair.MapHead[[]T, validation.Errors],
	)

	appendValues := F.Flow2(
		A.Push,
		pair.MapTail[validation.Errors, []T],
	)

	itemName := item.Name()

	zero := pair.Zero[validation.Errors, []T]()

	return func(i any) Decode[Context, []T] {

		res, ok := i.([]T)
		if ok {
			return reader.Of[Context](validation.Success(res))
		}

		val := reflect.ValueOf(i)
		if !val.IsValid() {
			return validation.FailureWithMessage[[]T](val, "invalid value")
		}
		kind := val.Kind()

		switch kind {
		case reflect.Array, reflect.Slice, reflect.String:

			return func(c Context) Validation[[]T] {

				return F.Pipe1(
					R.MonadReduceWithIndex(val, func(i int, p validationPair[[]T], v reflect.Value) validationPair[[]T] {
						vIface := v.Interface()
						return either.MonadFold(
							item.Validate(vIface)(appendContext(strconv.Itoa(i), itemName, vIface)(c)),
							appendErrors,
							appendValues,
						)(p)
					}, zero),
					pairToValidation,
				)
			}
		default:
			return validation.FailureWithMessage[[]T](val, fmt.Sprintf("type %s is not iterable", kind))
		}
	}
}

// Array creates a Type for array/slice values with elements of type T.
// It validates that input is an array, slice, or string, and validates each element
// using the provided item Type. During encoding, it maps the encode function over all elements.
//
// Type Parameters:
//   - T: The type of elements in the decoded array
//   - O: The type of elements in the encoded array
//
// Parameters:
//   - item: A Type[T, O, any] that defines how to validate/encode individual elements
//
// Returns:
//   - A Type[[]T, []O, any] that can validate, decode, and encode array values
//
// The function handles:
//   - Native Go slices of type []T (passed through directly)
//   - reflect.Array, reflect.Slice, reflect.String (validated element by element)
//   - Collects all validation errors from individual elements
//   - Provides detailed context for each element's position in error messages
//
// Example:
//
//	intArray := codec.Array(codec.Int())
//	result := intArray.Decode([]int{1, 2, 3})           // Success: Right([1, 2, 3])
//	result := intArray.Decode([]any{1, "2", 3})         // Failure: validation error at index 1
//	encoded := intArray.Encode([]int{1, 2, 3})          // Returns: []int{1, 2, 3}
//
//	stringArray := codec.Array(codec.String())
//	result := stringArray.Decode([]string{"a", "b"})    // Success: Right(["a", "b"])
//	result := stringArray.Decode("hello")               // Success: Right(["h", "e", "l", "l", "o"])
func Array[T, O any](item Type[T, O, any]) Type[[]T, []O, any] {

	validate := validateArray(item)
	is := isFromValidate(validate)
	name := fmt.Sprintf("Array[%s]", item.Name())

	return MakeType(
		name,
		is,
		validate,
		A.Map(item.Encode),
	)

}

// TranscodeArray creates a Type for array/slice values with strongly-typed input.
// Unlike Array which accepts any input type, TranscodeArray requires the input to be
// a slice of type []I, providing type safety at the input level.
//
// This function validates each element of the input slice using the provided item Type,
// transforming []I -> []T during decoding and []T -> []O during encoding.
//
// Type Parameters:
//   - T: The type of elements in the decoded array
//   - O: The type of elements in the encoded array
//   - I: The type of elements in the input array (must be a slice)
//
// Parameters:
//   - item: A Type[T, O, I] that defines how to validate/encode individual elements
//
// Returns:
//   - A Type[[]T, []O, []I] that can validate, decode, and encode array values
//
// The function:
//   - Requires input to be exactly []I (not any)
//   - Validates each element using the item Type's validation logic
//   - Collects all validation errors from individual elements
//   - Provides detailed context for each element's position in error messages
//   - Maps the encode function over all elements during encoding
//
// Example:
//
//	// Create a codec that transforms string slices to int slices
//	stringToInt := codec.MakeType[int, int, string](
//	    "StringToInt",
//	    func(s any) result.Result[int] { ... },
//	    func(s string) codec.Validate[int] { ... },
//	    func(i int) int { return i },
//	)
//	arrayCodec := codec.TranscodeArray(stringToInt)
//
//	// Decode: []string -> []int
//	result := arrayCodec.Decode([]string{"1", "2", "3"})  // Success: Right([1, 2, 3])
//	result := arrayCodec.Decode([]string{"1", "x", "3"})  // Failure: validation error at index 1
//
//	// Encode: []int -> []int
//	encoded := arrayCodec.Encode([]int{1, 2, 3})          // Returns: []int{1, 2, 3}
//
// Use TranscodeArray when:
//   - You need type-safe input validation ([]I instead of any)
//   - You're transforming between different slice element types
//   - You want compile-time guarantees about input types
//
// Use Array when:
//   - You need to accept various input types (any, reflect.Value, etc.)
//   - You're working with dynamic or unknown input types
func TranscodeArray[T, O, I any](item Type[T, O, I]) Type[[]T, []O, []I] {
	validate := validateArrayFromArray(item)
	is := isFromValidate(validate)
	name := fmt.Sprintf("Array[%s]", item.Name())

	return MakeType(
		name,
		is,
		validate,
		A.Map(item.Encode),
	)
}

func validateEitherFromEither[L, R, OL, OR, IL, IR any](
	leftItem Type[L, OL, IL],
	rightItem Type[R, OR, IR],
) Validate[either.Either[IL, IR], either.Either[L, R]] {

	// leftName := left.Name()
	// rightName := right.Name()

	return func(is either.Either[IL, IR]) Decode[Context, either.Either[L, R]] {

		return either.MonadFold(
			is,
			F.Flow2(
				leftItem.Validate,
				readereither.Map[Context, validation.Errors](either.Left[R, L]),
			),
			F.Flow2(
				rightItem.Validate,
				readereither.Map[Context, validation.Errors](either.Right[L, R]),
			),
		)

	}
}

// TranscodeEither creates a Type for Either values with strongly-typed left and right branches.
// It validates and transforms Either[IL, IR] to Either[L, R] during decoding, and
// Either[L, R] to Either[OL, OR] during encoding.
//
// This function is useful for handling sum types (discriminated unions) where a value can be
// one of two possible types. Each branch (Left and Right) is validated and transformed
// independently using its respective Type codec.
//
// Type Parameters:
//   - L: The type of the decoded Left value
//   - R: The type of the decoded Right value
//   - OL: The type of the encoded Left value
//   - OR: The type of the encoded Right value
//   - IL: The type of the input Left value
//   - IR: The type of the input Right value
//
// Parameters:
//   - leftItem: A Type[L, OL, IL] that defines how to validate/encode Left values
//   - rightItem: A Type[R, OR, IR] that defines how to validate/encode Right values
//
// Returns:
//   - A Type[Either[L, R], Either[OL, OR], Either[IL, IR]] that can validate, decode, and encode Either values
//
// The function:
//   - Validates Left values using leftItem's validation logic
//   - Validates Right values using rightItem's validation logic
//   - Preserves the Either structure (Left stays Left, Right stays Right)
//   - Provides context-aware error messages indicating which branch failed
//   - Transforms values through the respective codecs during encoding
//
// Example:
//
//	// Create a codec for Either[string, int]
//	stringCodec := codec.String()
//	intCodec := codec.Int()
//	eitherCodec := codec.TranscodeEither(stringCodec, intCodec)
//
//	// Decode Left value
//	leftResult := eitherCodec.Decode(either.Left[int]("error"))
//	// Success: Right(Either.Left("error"))
//
//	// Decode Right value
//	rightResult := eitherCodec.Decode(either.Right[string](42))
//	// Success: Right(Either.Right(42))
//
//	// Encode Left value
//	encodedLeft := eitherCodec.Encode(either.Left[int]("error"))
//	// Returns: Either.Left("error")
//
//	// Encode Right value
//	encodedRight := eitherCodec.Encode(either.Right[string](42))
//	// Returns: Either.Right(42)
//
// Use TranscodeEither when:
//   - You need to handle sum types or discriminated unions
//   - You want to validate and transform both branches of an Either independently
//   - You're working with error handling patterns (Left for errors, Right for success)
//   - You need type-safe transformations for both possible values
//
// Common patterns:
//   - Error handling: Either[Error, Value]
//   - Optional with reason: Either[Reason, Value]
//   - Validation results: Either[ValidationError, ValidatedData]
func TranscodeEither[L, R, OL, OR, IL, IR any](leftItem Type[L, OL, IL], rightItem Type[R, OR, IR]) Type[either.Either[L, R], either.Either[OL, OR], either.Either[IL, IR]] {
	validate := validateEitherFromEither(leftItem, rightItem)
	is := isFromValidate(validate)
	name := fmt.Sprintf("Either[%s, %s]", leftItem.Name(), rightItem.Name())

	return MakeType(
		name,
		is,
		validate,
		either.Fold(F.Flow2(
			leftItem.Encode,
			either.Left[OR, OL],
		), F.Flow2(
			rightItem.Encode,
			either.Right[OL, OR],
		)),
	)
}

func validateAlways[T any](is T) Decode[Context, T] {
	return reader.Of[Context](validation.Success(is))
}

// Id creates an identity Type codec that performs no transformation or validation.
//
// An identity codec is a Type[T, T, T] where:
//   - Decode: Always succeeds and returns the input value unchanged
//   - Encode: Returns the input value unchanged (identity function)
//   - Validation: Always succeeds without any checks
//
// This is useful as:
//   - A building block for more complex codecs
//   - A no-op codec when you need a Type but don't want any transformation
//   - A starting point for codec composition
//   - Testing and debugging codec pipelines
//
// Type Parameters:
//   - T: The type that passes through unchanged
//
// Returns:
//   - A Type[T, T, T] that performs identity operations on type T
//
// The codec:
//   - Name: Uses the type's string representation (e.g., "int", "string")
//   - Is: Checks if a value is of type T
//   - Validate: Always succeeds and returns the input value
//   - Encode: Identity function (returns input unchanged)
//
// Example:
//
//	// Create an identity codec for strings
//	stringId := codec.Id[string]()
//
//	// Decode always succeeds
//	result := stringId.Decode("hello")  // Success: Right("hello")
//
//	// Encode is identity
//	encoded := stringId.Encode("world")  // Returns: "world"
//
//	// Use in composition
//	arrayOfStrings := codec.TranscodeArray(stringId)
//	result := arrayOfStrings.Decode([]string{"a", "b", "c"})
//
// Use cases:
//   - When you need a Type but don't want any validation or transformation
//   - As a placeholder in generic code that requires a Type parameter
//   - Building blocks for TranscodeArray, TranscodeEither, etc.
//   - Testing codec composition without side effects
//
// Note: Unlike MakeSimpleType which validates the type, Id always succeeds
// in validation. It only checks the type during the Is operation.
func Id[T any]() Type[T, T, T] {
	return MakeType(
		fmt.Sprintf("%T", *new(T)),
		Is[T](),
		validateAlways[T],
		F.Identity[T],
	)
}

func validateFromRefinement[A, B any](refinement Refinement[A, B]) Validate[A, B] {

	return func(a A) Decode[Context, B] {

		return func(ctx Context) Validation[B] {
			return F.Pipe2(
				a,
				refinement.GetOption,
				either.FromOption[B](func() validation.Errors {
					return array.Of(&validation.ValidationError{
						Value:    a,
						Context:  ctx,
						Messsage: fmt.Sprintf("type cannot be refined: %s", refinement),
					})
				}),
			)
		}
	}
}

func isFromRefinement[A, B any](refinement Refinement[A, B]) ReaderResult[any, B] {

	isA := Is[A]()
	isB := Is[B]()

	err := fmt.Errorf("type cannot be refined: %s", refinement)

	isAtoB := F.Flow2(
		isA,
		result.ChainOptionK[A, B](lazy.Of(err))(refinement.GetOption),
	)

	return F.Pipe1(
		isAtoB,
		readereither.ChainLeft(reader.Of[error](isB)),
	)

}

// FromRefinement creates a Type codec from a Refinement (Prism).
//
// A Refinement[A, B] represents the concept that B is a specialized/refined version of A.
// For example, PositiveInt is a refinement of int, or NonEmptyString is a refinement of string.
// This function converts a Prism[A, B] into a Type[B, A, A] codec that can validate and transform
// between the base type A and the refined type B.
//
// Type Parameters:
//   - A: The base/broader type (e.g., int, string)
//   - B: The refined/specialized type (e.g., PositiveInt, NonEmptyString)
//
// Parameters:
//   - refinement: A Refinement[A, B] (which is a Prism[A, B]) that defines:
//   - GetOption: A → Option[B] - attempts to refine A to B (may fail if refinement conditions aren't met)
//   - ReverseGet: B → A - converts refined type back to base type (always succeeds)
//
// Returns:
//   - A Type[B, A, A] codec where:
//   - Decode: A → Validation[B] - validates that A satisfies refinement conditions and produces B
//   - Encode: B → A - converts refined type back to base type using ReverseGet
//   - Is: Checks if a value is of type B
//   - Name: Descriptive name including the refinement's string representation
//
// The codec:
//   - Uses the refinement's GetOption for validation during decoding
//   - Returns validation errors if the refinement conditions are not met
//   - Uses the refinement's ReverseGet for encoding (always succeeds)
//   - Provides context-aware error messages indicating why refinement failed
//
// Example:
//
//	// Define a refinement for positive integers
//	positiveIntPrism := prism.MakePrismWithName(
//	    func(n int) option.Option[int] {
//	        if n > 0 {
//	            return option.Some(n)
//	        }
//	        return option.None[int]()
//	    },
//	    func(n int) int { return n },
//	    "PositiveInt",
//	)
//
//	// Create a codec from the refinement
//	positiveIntCodec := codec.FromRefinement[int, int](positiveIntPrism)
//
//	// Decode: validates the refinement condition
//	result := positiveIntCodec.Decode(42)   // Success: Right(42)
//	result = positiveIntCodec.Decode(-5)    // Failure: validation error
//	result = positiveIntCodec.Decode(0)     // Failure: validation error
//
//	// Encode: converts back to base type
//	encoded := positiveIntCodec.Encode(42)  // Returns: 42
//
// Use cases:
//   - Creating codecs for refined types (positive numbers, non-empty strings, etc.)
//   - Validating that values meet specific constraints
//   - Building type-safe APIs with refined types
//   - Composing refinements with other codecs using Pipe
//
// Common refinement patterns:
//   - Numeric constraints: PositiveInt, NonNegativeFloat, BoundedInt
//   - String constraints: NonEmptyString, EmailAddress, URL
//   - Collection constraints: NonEmptyArray, UniqueElements
//   - Domain-specific constraints: ValidAge, ValidZipCode, ValidCreditCard
//
// Note: The refinement's GetOption returning None will result in a validation error
// with a message indicating the type cannot be refined. For more specific error messages,
// consider using MakeType directly with custom validation logic.
func FromRefinement[A, B any](refinement Refinement[A, B]) Type[B, A, A] {
	return MakeType(
		fmt.Sprintf("FromRefinement(%s)", refinement),
		isFromRefinement(refinement),
		validateFromRefinement(refinement),
		refinement.ReverseGet,
	)
}
