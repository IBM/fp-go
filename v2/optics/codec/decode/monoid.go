package decode

import "github.com/IBM/fp-go/v2/monoid"

// ApplicativeMonoid creates a Monoid instance for Decode[I, A] given a Monoid for A.
// This allows combining decoders where both the decoded values and validation errors
// are combined according to their respective monoid operations.
//
// The resulting monoid enables:
//   - Combining multiple decoders that produce monoidal values
//   - Accumulating validation errors when any decoder fails
//   - Building complex decoders from simpler ones through composition
//
// **Behavior**:
//   - Empty: Returns a decoder that always succeeds with the empty value from the inner monoid
//   - Concat: Combines two decoders:
//   - Both succeed: Combines decoded values using the inner monoid
//   - Any fails: Accumulates all validation errors using the Errors monoid
//
// This is particularly useful for:
//   - Aggregating results from multiple independent decoders
//   - Building decoders that combine partial results
//   - Validating and combining configuration from multiple sources
//   - Parallel validation with result accumulation
//
// Example - Combining string decoders:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	// Create a monoid for decoders that produce strings
//	m := ApplicativeMonoid[map[string]any](S.Monoid)
//
//	decoder1 := func(data map[string]any) Validation[string] {
//	    if name, ok := data["firstName"].(string); ok {
//	        return validation.Of(name)
//	    }
//	    return either.Left[string](validation.Errors{
//	        {Messsage: "missing firstName"},
//	    })
//	}
//
//	decoder2 := func(data map[string]any) Validation[string] {
//	    if name, ok := data["lastName"].(string); ok {
//	        return validation.Of(" " + name)
//	    }
//	    return either.Left[string](validation.Errors{
//	        {Messsage: "missing lastName"},
//	    })
//	}
//
//	// Combine decoders - will concatenate strings if both succeed
//	combined := m.Concat(decoder1, decoder2)
//	result := combined(map[string]any{
//	    "firstName": "John",
//	    "lastName":  "Doe",
//	}) // Success("John Doe")
//
// Example - Error accumulation:
//
//	// If any decoder fails, errors are accumulated
//	result := combined(map[string]any{}) // Failures with both error messages
//
// Example - Numeric aggregation:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intMonoid := monoid.MakeMonoid(N.Add[int], 0)
//	m := ApplicativeMonoid[string](intMonoid)
//
//	decoder1 := func(input string) Validation[int] {
//	    return validation.Of(10)
//	}
//	decoder2 := func(input string) Validation[int] {
//	    return validation.Of(32)
//	}
//
//	combined := m.Concat(decoder1, decoder2)
//	result := combined("input") // Success(42) - values are added
func ApplicativeMonoid[I, A any](m Monoid[A]) Monoid[Decode[I, A]] {
	return monoid.ApplicativeMonoid(
		Of[I, A],
		MonadMap[I, A, Endomorphism[A]],
		MonadAp[A, I, A],
		m,
	)
}

// AlternativeMonoid creates a Monoid instance for Decode[I, A] using the Alternative pattern.
// This combines applicative error-accumulation behavior with alternative fallback behavior,
// allowing you to both accumulate errors and provide fallback alternatives when combining decoders.
//
// The Alternative pattern provides two key operations:
//   - Applicative operations (Of, Map, Ap): accumulate errors when combining decoders
//   - Alternative operation (Alt): provide fallback when a decoder fails
//
// This monoid is particularly useful when you want to:
//   - Try multiple decoding strategies and fall back to alternatives
//   - Combine successful values using the provided monoid
//   - Accumulate all errors from failed attempts
//   - Build decoding pipelines with fallback logic
//
// **Behavior**:
//   - Empty: Returns a decoder that always succeeds with the empty value from the inner monoid
//   - Concat: Combines two decoders using both applicative and alternative semantics:
//   - If first succeeds and second succeeds: combines decoded values using inner monoid
//   - If first fails: tries second as fallback (alternative behavior)
//   - If both fail: **accumulates all errors from both decoders**
//
// **Error Aggregation**: When both decoders fail, all validation errors from both attempts
// are combined using the Errors monoid. This provides complete visibility into why all
// alternatives failed, which is essential for debugging and user feedback.
//
// Type Parameters:
//   - I: The input type being decoded
//   - A: The output type after successful decoding
//
// Parameters:
//   - m: The monoid for combining successful decoded values of type A
//
// Returns:
//
//	A Monoid[Decode[I, A]] that combines applicative and alternative behaviors
//
// Example - Combining successful decoders:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	m := AlternativeMonoid[string](S.Monoid)
//
//	decoder1 := func(input string) Validation[string] {
//	    return validation.Of("Hello")
//	}
//	decoder2 := func(input string) Validation[string] {
//	    return validation.Of(" World")
//	}
//
//	combined := m.Concat(decoder1, decoder2)
//	result := combined("input")
//	// Result: Success("Hello World") - values combined using string monoid
//
// Example - Fallback behavior:
//
//	m := AlternativeMonoid[string](S.Monoid)
//
//	failing := func(input string) Validation[string] {
//	    return either.Left[string](validation.Errors{
//	        {Value: input, Messsage: "primary failed"},
//	    })
//	}
//	fallback := func(input string) Validation[string] {
//	    return validation.Of("fallback value")
//	}
//
//	combined := m.Concat(failing, fallback)
//	result := combined("input")
//	// Result: Success("fallback value") - second decoder used as fallback
//
// Example - Error accumulation when both fail:
//
//	m := AlternativeMonoid[string](S.Monoid)
//
//	failing1 := func(input string) Validation[string] {
//	    return either.Left[string](validation.Errors{
//	        {Value: input, Messsage: "error 1"},
//	    })
//	}
//	failing2 := func(input string) Validation[string] {
//	    return either.Left[string](validation.Errors{
//	        {Value: input, Messsage: "error 2"},
//	    })
//	}
//
//	combined := m.Concat(failing1, failing2)
//	result := combined("input")
//	// Result: Failures with accumulated errors: ["error 1", "error 2"]
//
// Example - Building decoder with multiple fallbacks:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	m := AlternativeMonoid[string](N.MonoidSum[int]())
//
//	// Try to parse from different formats
//	parseJSON := func(input string) Validation[int] { /* ... */ }
//	parseYAML := func(input string) Validation[int] { /* ... */ }
//	parseINI := func(input string) Validation[int] { /* ... */ }
//
//	// Combine with fallback chain
//	decoder := m.Concat(m.Concat(parseJSON, parseYAML), parseINI)
//	// Uses first successful parser, or accumulates all errors if all fail
//
// Example - Combining multiple configuration sources:
//
//	type Config struct{ Port int }
//	configMonoid := monoid.MakeMonoid(
//	    func(a, b Config) Config {
//	        if b.Port != 0 { return b }
//	        return a
//	    },
//	    Config{Port: 0},
//	)
//
//	m := AlternativeMonoid[map[string]any](configMonoid)
//
//	fromEnv := func(data map[string]any) Validation[Config] { /* ... */ }
//	fromFile := func(data map[string]any) Validation[Config] { /* ... */ }
//	fromDefault := func(data map[string]any) Validation[Config] {
//	    return validation.Of(Config{Port: 8080})
//	}
//
//	// Try env, then file, then default
//	decoder := m.Concat(m.Concat(fromEnv, fromFile), fromDefault)
//	// Returns first successful config, or all errors if all fail
func AlternativeMonoid[I, A any](m Monoid[A]) Monoid[Decode[I, A]] {
	return monoid.AlternativeMonoid(
		Of[I, A],
		MonadMap[I, A, func(A) A],
		MonadAp[A, I, A],
		MonadAlt[I, A],
		m,
	)
}

// AltMonoid creates a Monoid instance for Decode[I, A] using the Alt (alternative) operation.
// This monoid provides a way to combine decoders with fallback behavior, where the second
// decoder is used as an alternative if the first one fails.
//
// The Alt operation implements the "try first, fallback to second" pattern, which is useful
// for decoding scenarios where you want to attempt multiple decoding strategies in sequence
// and use the first one that succeeds.
//
// **Behavior**:
//   - Empty: Returns the provided zero value (a lazy computation that produces a Decode[I, A])
//   - Concat: Combines two decoders using Alt semantics:
//   - If first succeeds: returns the first result (second is never evaluated)
//   - If first fails: tries the second decoder as fallback
//   - If both fail: **aggregates errors from both decoders**
//
// **Error Aggregation**: When both decoders fail, all validation errors from both attempts
// are combined using the Errors monoid. This ensures complete visibility into why all
// alternatives failed.
//
// This is different from [AlternativeMonoid] in that:
//   - AltMonoid uses a custom zero value (provided by the user)
//   - AlternativeMonoid derives the zero from an inner monoid
//   - AltMonoid is simpler and only provides fallback behavior
//   - AlternativeMonoid combines applicative and alternative behaviors
//
// Type Parameters:
//   - I: The input type being decoded
//   - A: The output type after successful decoding
//
// Parameters:
//   - zero: A lazy computation that produces the identity/empty Decode[I, A].
//     This is typically a decoder that always succeeds with a default value, or could be
//     a decoder that always fails representing "no decoding attempted"
//
// Returns:
//
//	A Monoid[Decode[I, A]] that combines decoders with fallback behavior
//
// Example - Using default value as zero:
//
//	m := AltMonoid(func() Decode[string, int] {
//	    return Of[string](0)
//	})
//
//	failing := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{
//	        {Value: input, Messsage: "failed"},
//	    })
//	}
//	succeeding := func(input string) Validation[int] {
//	    return validation.Of(42)
//	}
//
//	combined := m.Concat(failing, succeeding)
//	result := combined("input")
//	// Result: Success(42) - falls back to second decoder
//
//	empty := m.Empty()
//	result2 := empty("input")
//	// Result: Success(0) - the provided zero value
//
// Example - Chaining multiple fallbacks:
//
//	m := AltMonoid(func() Decode[string, Config] {
//	    return Of[string](defaultConfig)
//	})
//
//	primary := parseFromPrimarySource    // Fails
//	secondary := parseFromSecondarySource // Fails
//	tertiary := parseFromTertiarySource   // Succeeds
//
//	// Chain fallbacks
//	decoder := m.Concat(m.Concat(primary, secondary), tertiary)
//	result := decoder("input")
//	// Result: Success from tertiary - uses first successful decoder
//
// Example - Error aggregation when all fail:
//
//	m := AltMonoid(func() Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        return either.Left[int](validation.Errors{
//	            {Messsage: "no default available"},
//	        })
//	    }
//	})
//
//	failing1 := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{
//	        {Value: input, Messsage: "error 1"},
//	    })
//	}
//	failing2 := func(input string) Validation[int] {
//	    return either.Left[int](validation.Errors{
//	        {Value: input, Messsage: "error 2"},
//	    })
//	}
//
//	combined := m.Concat(failing1, failing2)
//	result := combined("input")
//	// Result: Failures with accumulated errors: ["error 1", "error 2"]
//
// Example - Building a decoder pipeline with fallbacks:
//
//	m := AltMonoid(func() Decode[string, Config] {
//	    return Of[string](defaultConfig)
//	})
//
//	// Try multiple decoding sources in order
//	decoders := []Decode[string, Config]{
//	    loadFromFile("config.json"),      // Try file first
//	    loadFromEnv,                       // Then environment
//	    loadFromRemote("api.example.com"), // Then remote API
//	}
//
//	// Fold using the monoid to get first successful config
//	result := array.MonoidFold(m)(decoders)
//	// Result: First successful config, or defaultConfig if all fail
//
// Example - Comparing with AlternativeMonoid:
//
//	// AltMonoid - simple fallback with custom zero
//	altM := AltMonoid(func() Decode[string, int] {
//	    return Of[string](0)
//	})
//
//	// AlternativeMonoid - combines values when both succeed
//	import N "github.com/IBM/fp-go/v2/number"
//	altMonoid := AlternativeMonoid[string](N.MonoidSum[int]())
//
//	decoder1 := Of[string](10)
//	decoder2 := Of[string](32)
//
//	// AltMonoid: returns first success (10)
//	result1 := altM.Concat(decoder1, decoder2)("input")
//	// Result: Success(10)
//
//	// AlternativeMonoid: combines both successes (10 + 32 = 42)
//	result2 := altMonoid.Concat(decoder1, decoder2)("input")
//	// Result: Success(42)
func AltMonoid[I, A any](zero Lazy[Decode[I, A]]) Monoid[Decode[I, A]] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[I, A],
	)
}
