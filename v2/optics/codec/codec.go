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
func (t *typeImpl[A, O, I]) Validate(i I) Reader[Context, Validation[A]] {
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

func validateFromIs[A any](
	is ReaderResult[any, A],
	msg string,
) Reader[any, Reader[Context, Validation[A]]] {
	return func(u any) Reader[Context, Validation[A]] {
		return F.Pipe2(
			u,
			is,
			result.Fold(
				validation.FailureWithError[A](u, msg),
				F.Flow2(
					validation.Success[A],
					reader.Of[Context],
				),
			),
		)
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
	var zero A
	name := fmt.Sprintf("%T", zero)
	is := Is[A]()

	return MakeType(
		name,
		is,
		validateFromIs(is, name),
		F.Identity[A],
	)
}

func String() Type[string, string, any] {
	return MakeSimpleType[string]()
}

func Int() Type[int, int, any] {
	return MakeSimpleType[int]()
}

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

func validateArray[T any](item Type[T, T, any]) func(u any) Reader[Context, Validation[[]T]] {

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

	return func(u any) Reader[Context, Validation[[]T]] {
		val := reflect.ValueOf(u)
		if !val.IsValid() {
			return validation.FailureWithMessage[[]T](val, "invalid value")
		}
		kind := val.Kind()

		switch kind {
		case reflect.Array, reflect.Slice, reflect.String:

			return func(c Context) Validation[[]T] {

				return F.Pipe1(
					R.MonadReduceWithIndex(val, func(i int, p validationPair[[]T], v reflect.Value) validationPair[[]T] {
						return either.MonadFold(
							item.Validate(v)(appendContext(strconv.Itoa(i), itemName, v)(c)),
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
