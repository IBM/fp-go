package codec

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/decoder"
	"github.com/IBM/fp-go/v2/optics/encoder"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

type (
	ReaderResult[R, A any] = readerresult.ReaderResult[R, A]

	// Lazy represents a lazily evaluated value.
	Lazy[A any] = lazy.Lazy[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Either represents a value that can be one of two types: Left (error) or Right (success).
	Either[E, A any] = either.Either[E, A]

	// Result represents a computation that may fail with an error.
	Result[A any] = result.Result[A]

	// Codec combines a Decoder and an Encoder for bidirectional transformations.
	// It can decode input I to type A and encode type A to output O.
	Codec[I, O, A any] struct {
		Decode decoder.Decoder[I, A]
		Encode encoder.Encoder[O, A]
	}

	Validation[A any] = validation.Validation[A]

	Context = validation.Context

	// Validate is a function that validates input I to produce type A.
	// It takes an input and returns a Reader that depends on the validation Context.
	Validate[I, A any] = Reader[I, Reader[Context, Validation[A]]]

	// Decode is a function that decodes input I to type A with validation.
	// It returns a Validation result directly.
	Decode[I, A any] = Reader[I, Validation[A]]

	// Encode is a function that encodes type A to output O.
	Encode[A, O any] = Reader[A, O]

	// Decoder is an interface for types that can decode and validate input.
	Decoder[I, A any] interface {
		Name() string
		Validate(I) Reader[Context, Validation[A]]
		Decode(I) Validation[A]
	}

	// Encoder is an interface for types that can encode values.
	Encoder[A, O any] interface {
		// Encode transforms a value of type A into output format O.
		Encode(A) O
	}
	// Type is a bidirectional codec that combines encoding, decoding, validation,
	// and type checking capabilities. It represents a complete specification of
	// how to work with a particular type.
	Type[A, O, I any] interface {
		Decoder[I, A]
		Encoder[A, O]
		AsDecoder() Decoder[I, A]
		AsEncoder() Encoder[A, O]
		Is(any) Result[A]
	}

	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Pair[L, R any] = pair.Pair[L, R]
)
