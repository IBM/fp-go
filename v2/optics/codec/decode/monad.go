package decode

import (
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

// Of creates a Decode that always succeeds with the given value.
// This is the pointed functor operation that lifts a pure value into the Decode context.
//
// Example:
//
//	decoder := decode.Of[string](42)
//	result := decoder("any input") // Always returns validation.Success(42)
func Of[I, A any](a A) Decode[I, A] {
	return reader.Of[I](validation.Of(a))
}

// MonadChain sequences two decode operations, passing the result of the first to the second.
// This is the monadic bind operation that enables sequential composition of decoders.
//
// Example:
//
//	decoder1 := decode.Of[string](42)
//	decoder2 := decode.MonadChain(decoder1, func(n int) Decode[string, string] {
//	    return decode.Of[string](fmt.Sprintf("Number: %d", n))
//	})
func MonadChain[I, A, B any](fa Decode[I, A], f Kleisli[I, A, B]) Decode[I, B] {
	return readert.MonadChain(
		validation.MonadChain,
		fa,
		f,
	)
}

// Chain creates an operator that sequences decode operations.
// This is the curried version of MonadChain, useful for composition pipelines.
//
// Example:
//
//	chainOp := decode.Chain(func(n int) Decode[string, string] {
//	    return decode.Of[string](fmt.Sprintf("Number: %d", n))
//	})
//	decoder := chainOp(decode.Of[string](42))
func Chain[I, A, B any](f Kleisli[I, A, B]) Operator[I, A, B] {
	return readert.Chain[Decode[I, A]](
		validation.Chain,
		f,
	)
}

// ChainLeft transforms the error channel of a decoder, enabling error recovery and context addition.
// This is the left-biased monadic chain operation that operates on validation failures.
//
// **Key behaviors**:
//   - Success values pass through unchanged - the handler is never called
//   - On failure, the handler receives the errors and can recover or add context
//   - When the handler also fails, **both original and new errors are aggregated**
//   - The handler returns a Decode[I, A], giving it access to the original input
//
// **Error Aggregation**: Unlike standard Either operations, when the transformation function
// returns a failure, both the original errors AND the new errors are combined using the
// Errors monoid. This ensures no validation errors are lost.
//
// Use cases:
//   - Adding contextual information to validation errors
//   - Recovering from specific error conditions
//   - Transforming error messages while preserving original errors
//   - Implementing conditional recovery based on error types
//
// Example - Error recovery:
//
//	failingDecoder := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{
//	        {Value: input, Messsage: "not found"},
//	    })
//	}
//
//	recoverFromNotFound := ChainLeft(func(errs Errors) Decode[string, int] {
//	    for _, err := range errs {
//	        if err.Messsage == "not found" {
//	            return Of[string](0) // recover with default
//	        }
//	    }
//	    return func(input string) Validation[int] {
//	        return either.Left[int](errs)
//	    }
//	})
//
//	decoder := recoverFromNotFound(failingDecoder)
//	result := decoder("input") // Success(0) - recovered from failure
//
// Example - Adding context:
//
//	addContext := ChainLeft(func(errs Errors) Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{
//	            {
//	                Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
//	                Messsage: "failed to decode user age",
//	            },
//	        })
//	    }
//	})
//	// Result will contain BOTH original error and context error
func ChainLeft[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return readert.Chain[Decode[I, A]](
		validation.ChainLeft,
		f,
	)
}

// OrElse provides fallback decoding logic when the primary decoder fails.
// This is an alias for ChainLeft with a more semantic name for fallback scenarios.
//
// **OrElse is exactly the same as ChainLeft** - they are aliases with identical implementations
// and behavior. The choice between them is purely about code readability and semantic intent:
//   - Use **OrElse** when emphasizing fallback/alternative decoding logic
//   - Use **ChainLeft** when emphasizing technical error channel transformation
//
// **Key behaviors** (identical to ChainLeft):
//   - Success values pass through unchanged - the handler is never called
//   - On failure, the handler receives the errors and can provide an alternative
//   - When the handler also fails, **both original and new errors are aggregated**
//   - The handler returns a Decode[I, A], giving it access to the original input
//
// The name "OrElse" reads naturally in code: "try this decoder, or else try this alternative."
// This makes it ideal for expressing fallback logic and default values.
//
// Use cases:
//   - Providing default values when decoding fails
//   - Trying alternative decoding strategies
//   - Implementing fallback chains with multiple alternatives
//   - Input-dependent recovery (using access to original input)
//
// Example - Simple fallback:
//
//	primaryDecoder := func(input string) Validation[int] {
//	    n, err := strconv.Atoi(input)
//	    if err != nil {
//	        return either.Left[int](validation.Errors{
//	            {Value: input, Messsage: "not a valid integer"},
//	        })
//	    }
//	    return validation.Of(n)
//	}
//
//	withDefault := OrElse(func(errs Errors) Decode[string, int] {
//	    return Of[string](0) // default to 0 if decoding fails
//	})
//
//	decoder := withDefault(primaryDecoder)
//	result1 := decoder("42")  // Success(42)
//	result2 := decoder("abc") // Success(0) - fallback
//
// Example - Input-dependent fallback:
//
//	smartDefault := OrElse(func(errs Errors) Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        // Access original input to determine appropriate default
//	        if strings.Contains(input, "http") {
//	            return validation.Of(80)
//	        }
//	        if strings.Contains(input, "https") {
//	            return validation.Of(443)
//	        }
//	        return validation.Of(8080)
//	    }
//	})
//
//	decoder := smartDefault(decodePort)
//	result1 := decoder("http-server")  // Success(80)
//	result2 := decoder("https-server") // Success(443)
//	result3 := decoder("other")        // Success(8080)
func OrElse[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return ChainLeft(f)
}

// MonadMap transforms the decoded value using the provided function.
// This is the functor map operation that applies a transformation to successful decode results.
//
// Example:
//
//	decoder := decode.Of[string](42)
//	mapped := decode.MonadMap(decoder, func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
func MonadMap[I, A, B any](fa Decode[I, A], f func(A) B) Decode[I, B] {
	return readert.MonadMap[
		Decode[I, A],
		Decode[I, B]](
		validation.MonadMap,
		fa,
		f,
	)
}

// Map creates an operator that transforms decoded values.
// This is the curried version of MonadMap, useful for composition pipelines.
//
// Example:
//
//	mapOp := decode.Map(func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	decoder := mapOp(decode.Of[string](42))
func Map[I, A, B any](f func(A) B) Operator[I, A, B] {
	return readert.Map[
		Decode[I, A],
		Decode[I, B]](
		validation.Map,
		f,
	)
}

// MonadAp applies a decoder containing a function to a decoder containing a value.
// This is the applicative apply operation that enables parallel composition of decoders.
//
// Example:
//
//	decoderFn := decode.Of[string](func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	decoderVal := decode.Of[string](42)
//	result := decode.MonadAp(decoderFn, decoderVal)
func MonadAp[B, I, A any](fab Decode[I, func(A) B], fa Decode[I, A]) Decode[I, B] {
	return readert.MonadAp[
		Decode[I, A],
		Decode[I, B],
		Decode[I, func(A) B], I, A](
		validation.MonadAp[B, A],
		fab,
		fa,
	)
}

// Ap creates an operator that applies a function decoder to a value decoder.
// This is the curried version of MonadAp, useful for composition pipelines.
//
// Example:
//
//	apOp := decode.Ap[string](decode.Of[string](42))
//	decoderFn := decode.Of[string](func(n int) string {
//	    return fmt.Sprintf("Number: %d", n)
//	})
//	result := apOp(decoderFn)
func Ap[B, I, A any](fa Decode[I, A]) Operator[I, func(A) B, B] {
	return readert.Ap[
		Decode[I, A],
		Decode[I, B],
		Decode[I, func(A) B], I, A](
		validation.Ap[B, A],
		fa,
	)
}
