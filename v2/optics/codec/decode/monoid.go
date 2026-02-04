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

// func AlternativeMonoid[I, A any](m Monoid[A]) Monoid[Decode[I, A]] {
// 	return monoid.AlternativeMonoid(
// 		Of[I, A],
// 		MonadMap[I, A, func(A) A],
// 		MonadAp[A, I, A],
// 		MonadAlt[I, A],
// 		m,
// 	)
// }
