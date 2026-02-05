package decode

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/readereither"
)

// Of creates a Decode that always succeeds with the given value.
// This is the pointed functor operation that lifts a pure value into the Decode context.
//
// Example:
//
//	decoder := decode.Of[string](42)
//	result := decoder("any input") // Always returns validation.Success(42)
func Of[I, A any](a A) Decode[I, A] {
	return readereither.Of[I, Errors](a)
}

// Left creates a Decode that always fails with the given validation errors.
// This is the dual of Of - while Of lifts a success value, Left lifts failure errors
// into the Decode context.
//
// Left is useful for:
//   - Creating decoders that represent known failure states
//   - Short-circuiting decode pipelines with specific errors
//   - Building custom validation error responses
//   - Testing error handling paths
//
// The returned decoder ignores its input and always returns a validation failure
// containing the provided errors. This makes it the identity element for the
// Alt/OrElse operations when used as a fallback.
//
// Type signature: func(Errors) Decode[I, A]
//   - Takes validation errors
//   - Returns a decoder that always fails with those errors
//   - The decoder ignores its input of type I
//   - The failure type A can be any type (phantom type)
//
// Example - Creating a failing decoder:
//
//	failDecoder := decode.Left[string, int](validation.Errors{
//	    &validation.ValidationError{
//	        Value:    nil,
//	        Messsage: "operation not supported",
//	    },
//	})
//	result := failDecoder("any input") // Always fails with the error
//
// Example - Short-circuiting with specific errors:
//
//	validateAge := func(age int) Decode[map[string]any, int] {
//	    if age < 0 {
//	        return decode.Left[map[string]any, int](validation.Errors{
//	            &validation.ValidationError{
//	                Value:    age,
//	                Context:  validation.Context{{Key: "age", Type: "int"}},
//	                Messsage: "age cannot be negative",
//	            },
//	        })
//	    }
//	    return decode.Of[map[string]any](age)
//	}
//
// Example - Building error responses:
//
//	notFoundError := decode.Left[string, User](validation.Errors{
//	    &validation.ValidationError{
//	        Messsage: "user not found",
//	    },
//	})
//
//	decoder := decode.MonadAlt(
//	    tryFindUser,
//	    func() Decode[string, User] { return notFoundError },
//	)
//
// Example - Testing error paths:
//
//	// Create a decoder that always fails for testing
//	alwaysFails := decode.Left[string, int](validation.Errors{
//	    &validation.ValidationError{Messsage: "test error"},
//	})
//
//	// Test error recovery logic
//	recovered := decode.OrElse(func(errs Errors) Decode[string, int] {
//	    return decode.Of[string](0) // recover with default
//	})(alwaysFails)
//
//	result := recovered("input") // Success(0)
func Left[I, A any](err Errors) Decode[I, A] {
	return readereither.Left[I, A](err)
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

// MonadChainLeft transforms the error channel of a decoder, enabling error recovery and context addition.
// This is the uncurried version of ChainLeft, taking both the decoder and the transformation function directly.
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
// This function is the direct, uncurried form of ChainLeft. Use ChainLeft when you need
// a curried operator for composition pipelines, and use MonadChainLeft when you have both
// the decoder and transformation function available at once.
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
//	recoverFromNotFound := func(errs Errors) Decode[string, int] {
//	    for _, err := range errs {
//	        if err.Messsage == "not found" {
//	            return Of[string](0) // recover with default
//	        }
//	    }
//	    return func(input string) Validation[int] {
//	        return either.Left[int](errs)
//	    }
//	}
//
//	decoder := MonadChainLeft(failingDecoder, recoverFromNotFound)
//	result := decoder("input") // Success(0) - recovered from failure
//
// Example - Adding context:
//
//	addContext := func(errs Errors) Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{
//	            {
//	                Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
//	                Messsage: "failed to decode user age",
//	            },
//	        })
//	    }
//	}
//
//	decoder := MonadChainLeft(failingDecoder, addContext)
//	result := decoder("abc")
//	// Result will contain BOTH original error and context error
//
// Example - Comparison with ChainLeft:
//
//	// MonadChainLeft - direct application
//	result1 := MonadChainLeft(decoder, handler)("input")
//
//	// ChainLeft - curried for pipelines
//	result2 := ChainLeft(handler)(decoder)("input")
//
//	// Both produce identical results
func MonadChainLeft[I, A any](fa Decode[I, A], f Kleisli[I, Errors, A]) Decode[I, A] {
	return readert.MonadChain(
		validation.MonadChainLeft,
		fa,
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

// MonadAlt provides alternative/fallback decoding with error aggregation.
// This is the Alternative pattern's core operation that tries the first decoder,
// and if it fails, tries the second decoder as a fallback.
//
// **Key behaviors**:
//   - If first succeeds: returns the first result (second is never evaluated)
//   - If first fails and second succeeds: returns the second result
//   - If both fail: **aggregates errors from both decoders**
//
// **Error Aggregation**: Unlike simple fallback patterns, when both decoders fail,
// MonadAlt combines ALL errors from both attempts using the Errors monoid. This ensures
// complete visibility into why all alternatives failed, which is crucial for debugging
// and providing comprehensive error messages to users.
//
// The name "Alt" comes from the Alternative type class in functional programming,
// which represents computations with a notion of choice and failure.
//
// Use cases:
//   - Trying multiple decoding strategies for the same input
//   - Providing fallback decoders when primary decoder fails
//   - Building validation pipelines with multiple alternatives
//   - Implementing "try this, or else try that" logic
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
//	fallbackDecoder := func() Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        // Try parsing as float and converting to int
//	        f, err := strconv.ParseFloat(input, 64)
//	        if err != nil {
//	            return either.Left[int](validation.Errors{
//	                {Value: input, Messsage: "not a valid number"},
//	            })
//	        }
//	        return validation.Of(int(f))
//	    }
//	}
//
//	decoder := MonadAlt(primaryDecoder, fallbackDecoder)
//	result1 := decoder("42")    // Success(42) - primary succeeds
//	result2 := decoder("42.5")  // Success(42) - fallback succeeds
//	result3 := decoder("abc")   // Failures with both errors aggregated
//
// Example - Multiple alternatives:
//
//	decoder1 := parseAsJSON
//	decoder2 := func() Decode[string, Config] { return parseAsYAML }
//	decoder3 := func() Decode[string, Config] { return parseAsINI }
//
//	// Try JSON, then YAML, then INI
//	decoder := MonadAlt(MonadAlt(decoder1, decoder2), decoder3)
//	// If all fail, errors from all three attempts are aggregated
//
// Example - Error aggregation:
//
//	failing1 := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{
//	        {Messsage: "primary decoder failed"},
//	    })
//	}
//	failing2 := func() Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{
//	            {Messsage: "fallback decoder failed"},
//	        })
//	    }
//	}
//
//	decoder := MonadAlt(failing1, failing2)
//	result := decoder("input")
//	// Result contains BOTH errors: ["primary decoder failed", "fallback decoder failed"]
func MonadAlt[I, A any](first Decode[I, A], second Lazy[Decode[I, A]]) Decode[I, A] {
	return MonadChainLeft(first, function.Ignore1of1[Errors](second))
}

// Alt creates an operator that provides alternative/fallback decoding with error aggregation.
// This is the curried version of MonadAlt, useful for composition pipelines.
//
// **Key behaviors** (identical to MonadAlt):
//   - If first succeeds: returns the first result (second is never evaluated)
//   - If first fails and second succeeds: returns the second result
//   - If both fail: **aggregates errors from both decoders**
//
// The Alt operator enables building reusable fallback chains that can be applied
// to different decoders. It reads naturally in pipelines: "apply this decoder,
// with this alternative if it fails."
//
// Use cases:
//   - Creating reusable fallback strategies
//   - Building decoder combinators with alternatives
//   - Composing multiple fallback layers
//   - Implementing retry logic with different strategies
//
// Example - Creating a reusable fallback:
//
//	// Create an operator that falls back to a default value
//	withDefault := Alt(func() Decode[string, int] {
//	    return Of[string](0)
//	})
//
//	// Apply to any decoder
//	decoder1 := withDefault(parseInteger)
//	decoder2 := withDefault(parseFromJSON)
//
//	result1 := decoder1("42")  // Success(42)
//	result2 := decoder1("abc") // Success(0) - fallback
//
// Example - Composing multiple alternatives:
//
//	tryYAML := Alt(func() Decode[string, Config] { return parseAsYAML })
//	tryINI := Alt(func() Decode[string, Config] { return parseAsINI })
//	useDefault := Alt(func() Decode[string, Config] {
//	    return Of[string](defaultConfig)
//	})
//
//	// Build a pipeline: try JSON, then YAML, then INI, then default
//	decoder := useDefault(tryINI(tryYAML(parseAsJSON)))
//
// Example - Error aggregation in pipeline:
//
//	failing1 := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{{Messsage: "error 1"}})
//	}
//	failing2 := func() Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{{Messsage: "error 2"}})
//	    }
//	}
//	failing3 := func() Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{{Messsage: "error 3"}})
//	    }
//	}
//
//	// Chain multiple alternatives
//	decoder := Alt(failing3)(Alt(failing2)(failing1))
//	result := decoder("input")
//	// Result contains ALL errors: ["error 1", "error 2", "error 3"]
func Alt[I, A any](second Lazy[Decode[I, A]]) Operator[I, A, A] {
	return ChainLeft(function.Ignore1of1[Errors](second))
}
