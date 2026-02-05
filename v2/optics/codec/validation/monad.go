package validation

import (
	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/applicative"
)

var errorsMonoid = ErrorsMonoid()

// Of creates a successful validation result containing the given value.
// This is the pure/return operation for the Validation monad.
//
// Example:
//
//	valid := Of(42) // Validation[int] containing 42
func Of[A any](a A) Validation[A] {
	return either.Of[Errors](a)
}

// Ap applies a validation containing a function to a validation containing a value.
// This is the applicative apply operation that accumulates errors from both validations.
// If either validation fails, all errors are collected. If both succeed, the function is applied.
//
// This enables combining multiple validations while collecting all errors:
//
// Example:
//
//	// Validate multiple fields and collect all errors
//	validateUser := Ap(Ap(Of(func(name string) func(age int) User {
//		return func(age int) User { return User{name, age} }
//	}))(validateName))(validateAge)
func Ap[B, A any](fa Validation[A]) Operator[func(A) B, B] {
	return either.ApV[B, A](errorsMonoid)(fa)
}

// MonadAp applies a validation containing a function to a validation containing a value.
// This is the applicative apply operation that **accumulates errors** from both validations.
//
// **Key behavior**: Unlike Either's MonadAp which fails fast (returns first error),
// this validation-specific implementation **accumulates all errors** using the Errors monoid.
// When both the function validation and value validation fail, all errors from both are combined.
//
// This error accumulation is the defining characteristic of the Validation applicative,
// making it ideal for scenarios where you want to collect all validation failures at once
// rather than stopping at the first error.
//
// Behavior:
//   - Both succeed: applies the function to the value → Success(result)
//   - Function fails, value succeeds: returns function's errors → Failure(func errors)
//   - Function succeeds, value fails: returns value's errors → Failure(value errors)
//   - Both fail: **combines all errors** → Failure(func errors + value errors)
//
// This is particularly useful for:
//   - Form validation: collect all field errors at once
//   - Configuration validation: report all invalid settings together
//   - Data validation: accumulate all constraint violations
//   - Multi-field validation: validate independent fields in parallel
//
// Example - Both succeed:
//
//	double := func(x int) int { return x * 2 }
//	result := MonadAp(Of(double), Of(21))
//	// Result: Success(42)
//
// Example - Error accumulation (key feature):
//
//	funcValidation := Failures[func(int) int](Errors{
//	    &ValidationError{Messsage: "function error"},
//	})
//	valueValidation := Failures[int](Errors{
//	    &ValidationError{Messsage: "value error"},
//	})
//	result := MonadAp(funcValidation, valueValidation)
//	// Result: Failure with BOTH errors: ["function error", "value error"]
//
// Example - Validating multiple fields:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	makeUser := func(name string) func(int) User {
//	    return func(age int) User { return User{name, age} }
//	}
//
//	nameValidation := validateName("ab")  // Fails: too short
//	ageValidation := validateAge(16)      // Fails: too young
//
//	// First apply name
//	step1 := MonadAp(Of(makeUser), nameValidation)
//	// Then apply age
//	result := MonadAp(step1, ageValidation)
//	// Result contains ALL validation errors from both fields
func MonadAp[B, A any](fab Validation[func(A) B], fa Validation[A]) Validation[B] {
	return either.MonadApV[B, A](errorsMonoid)(fab, fa)
}

// Map transforms the value inside a successful validation using the provided function.
// If the validation is a failure, the errors are preserved unchanged.
// This is the functor map operation for Validation.
//
// Map is used for transforming successful values without changing the validation context.
// It's the most basic operation for working with validated values and forms the foundation
// for more complex validation pipelines.
//
// Behavior:
//   - Success: applies function to value → Success(f(value))
//   - Failure: preserves errors unchanged → Failure(same errors)
//
// This is useful for:
//   - Type transformations: converting validated values to different types
//   - Value transformations: normalizing, formatting, or computing derived values
//   - Pipeline composition: chaining multiple transformations
//   - Preserving validation context: errors pass through unchanged
//
// Example - Transform successful value:
//
//	doubled := Map(func(x int) int { return x * 2 })(Of(21))
//	// Result: Success(42)
//
// Example - Failure preserved:
//
//	result := Map(func(x int) int { return x * 2 })(
//	    Failures[int](Errors{&ValidationError{Messsage: "invalid"}}),
//	)
//	// Result: Failure with same error: ["invalid"]
//
// Example - Type transformation:
//
//	toString := Map(func(x int) string { return fmt.Sprintf("%d", x) })
//	result := toString(Of(42))
//	// Result: Success("42")
//
// Example - Chaining transformations:
//
//	result := F.Pipe3(
//	    Of(5),
//	    Map(func(x int) int { return x + 10 }),  // 15
//	    Map(func(x int) int { return x * 2 }),   // 30
//	    Map(func(x int) string { return fmt.Sprintf("%d", x) }),  // "30"
//	)
//	// Result: Success("30")
func Map[A, B any](f func(A) B) Operator[A, B] {
	return either.Map[Errors](f)
}

// MonadMap transforms the value inside a successful validation using the provided function.
// If the validation is a failure, the errors are preserved unchanged.
// This is the non-curried version of [Map].
//
// MonadMap is useful when you have both the validation and the transformation function
// available at the same time, rather than needing to create a reusable operator.
//
// Behavior:
//   - Success: applies function to value → Success(f(value))
//   - Failure: preserves errors unchanged → Failure(same errors)
//
// Example - Transform successful value:
//
//	result := MonadMap(Of(21), func(x int) int { return x * 2 })
//	// Result: Success(42)
//
// Example - Failure preserved:
//
//	result := MonadMap(
//	    Failures[int](Errors{&ValidationError{Messsage: "invalid"}}),
//	    func(x int) int { return x * 2 },
//	)
//	// Result: Failure with same error: ["invalid"]
//
// Example - Type transformation:
//
//	result := MonadMap(Of(42), func(x int) string {
//	    return fmt.Sprintf("Value: %d", x)
//	})
//	// Result: Success("Value: 42")
//
// Example - Computing derived values:
//
//	type User struct { FirstName, LastName string }
//	result := MonadMap(
//	    Of(User{"John", "Doe"}),
//	    func(u User) string { return u.FirstName + " " + u.LastName },
//	)
//	// Result: Success("John Doe")
func MonadMap[A, B any](fa Validation[A], f func(A) B) Validation[B] {
	return either.MonadMap(fa, f)
}

// Chain is the curried version of [MonadChain].
// Sequences two validation computations where the second depends on the first.
//
// Example:
//
//	validatePositive := func(x int) Validation[int] {
//	    if x > 0 { return Success(x) }
//	    return Failure("must be positive")
//	}
//	result := Chain(validatePositive)(Success(42)) // Success(42)
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return either.Chain(f)
}

// MonadChain sequences two validation computations where the second depends on the first.
// If the first validation fails, returns the failure without executing the second.
// This is the monadic bind operation for Validation.
//
// Example:
//
//	result := MonadChain(
//	    Success(42),
//	    func(x int) Validation[string] {
//	        return Success(fmt.Sprintf("Value: %d", x))
//	    },
//	) // Success("Value: 42")
func MonadChain[A, B any](fa Validation[A], f Kleisli[A, B]) Validation[B] {
	return either.MonadChain(fa, f)
}

// chainErrors is an internal helper that chains error transformations while accumulating errors.
// When the transformation function f returns a failure, it concatenates the original errors (e1)
// with the new errors (e2) using the Errors monoid, ensuring all validation errors are preserved.
func chainErrors[A any](f Kleisli[Errors, A]) func(Errors) Validation[A] {
	return func(e1 Errors) Validation[A] {
		return either.MonadFold(
			f(e1),
			function.Flow2(array.Concat(e1), either.Left[A]),
			Of[A],
		)
	}
}

// ChainLeft is the curried version of [MonadChainLeft].
// Returns a function that transforms validation failures while preserving successes.
//
// Unlike the standard Either ChainLeft which replaces errors, this validation-specific
// implementation **aggregates errors** using the Errors monoid. When the transformation
// function returns a failure, both the original errors and the new errors are combined,
// ensuring no validation errors are lost.
//
// This is particularly useful for:
//   - Error recovery with fallback validation
//   - Adding contextual information to existing errors
//   - Transforming error types while preserving all error details
//   - Building error handling pipelines that accumulate failures
//
// Key behavior:
//   - Success values pass through unchanged
//   - When transforming failures, if the transformation also fails, **all errors are aggregated**
//   - If the transformation succeeds, it recovers from the original failure
//
// Example - Error recovery with aggregation:
//
//	recoverFromNotFound := ChainLeft(func(errs Errors) Validation[int] {
//	    // Check if this is a "not found" error
//	    for _, err := range errs {
//	        if err.Messsage == "not found" {
//	            return Success(0) // recover with default
//	        }
//	    }
//	    // Add context to existing errors
//	    return Failures[int](Errors{
//	        &ValidationError{Messsage: "recovery failed"},
//	    })
//	    // Result will contain BOTH original errors AND "recovery failed"
//	})
//
//	result := recoverFromNotFound(Failures[int](Errors{
//	    &ValidationError{Messsage: "database error"},
//	}))
//	// Result contains: ["database error", "recovery failed"]
//
// Example - Adding context to errors:
//
//	addContext := ChainLeft(func(errs Errors) Validation[string] {
//	    // Add contextual information
//	    return Failures[string](Errors{
//	        &ValidationError{
//	            Messsage: "validation failed in user.email field",
//	        },
//	    })
//	    // Original errors are preserved and new context is added
//	})
//
//	result := F.Pipe1(
//	    Failures[string](Errors{
//	        &ValidationError{Messsage: "invalid format"},
//	    }),
//	    addContext,
//	)
//	// Result contains: ["invalid format", "validation failed in user.email field"]
//
// Example - Success values pass through:
//
//	handler := ChainLeft(func(errs Errors) Validation[int] {
//	    return Failures[int](Errors{
//	        &ValidationError{Messsage: "never called"},
//	    })
//	})
//	result := handler(Success(42)) // Success(42) - unchanged
func ChainLeft[A any](f Kleisli[Errors, A]) Operator[A, A] {
	return either.Fold(
		chainErrors(f),
		Of[A],
	)
}

// MonadChainLeft sequences a computation on the failure (Left) channel of a Validation.
// If the Validation is a failure, applies the function to transform or recover from the errors.
// If the Validation is a success, returns the success value unchanged.
//
// **Critical difference from Either.MonadChainLeft**: This validation-specific implementation
// **aggregates errors** using the Errors monoid. When the transformation function returns a
// failure, both the original errors and the new errors are combined, ensuring comprehensive
// error reporting.
//
// This is the dual of [MonadChain] - while Chain operates on success values, ChainLeft
// operates on failure values. It's particularly useful for:
//   - Error recovery: converting specific errors into successful values
//   - Error enrichment: adding context or transforming error messages
//   - Fallback logic: providing alternative validations when the first fails
//   - Error aggregation: combining multiple validation failures
//
// The function parameter receives the collection of validation errors and must return
// a new Validation[A]. This allows you to:
//   - Recover by returning Success(value)
//   - Transform errors by returning Failures(newErrors) - **original errors are preserved**
//   - Implement conditional error handling based on error content
//
// Example - Error recovery:
//
//	result := MonadChainLeft(
//	    Failures[int](Errors{
//	        &ValidationError{Messsage: "not found"},
//	    }),
//	    func(errs Errors) Validation[int] {
//	        // Check if we can recover
//	        for _, err := range errs {
//	            if err.Messsage == "not found" {
//	                return Success(0) // recover with default value
//	            }
//	        }
//	        return Failures[int](errs) // propagate errors
//	    },
//	) // Success(0)
//
// Example - Error aggregation (key feature):
//
//	result := MonadChainLeft(
//	    Failures[string](Errors{
//	        &ValidationError{Messsage: "error 1"},
//	        &ValidationError{Messsage: "error 2"},
//	    }),
//	    func(errs Errors) Validation[string] {
//	        // Transformation also fails
//	        return Failures[string](Errors{
//	            &ValidationError{Messsage: "error 3"},
//	        })
//	    },
//	)
//	// Result contains ALL errors: ["error 1", "error 2", "error 3"]
//	// This is different from Either.MonadChainLeft which would only keep "error 3"
//
// Example - Adding context to errors:
//
//	result := MonadChainLeft(
//	    Failures[int](Errors{
//	        &ValidationError{Value: "abc", Messsage: "invalid number"},
//	    }),
//	    func(errs Errors) Validation[int] {
//	        // Add contextual information
//	        contextErrors := Errors{
//	            &ValidationError{
//	                Context:  []ContextEntry{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
//	                Messsage: "failed to parse user age",
//	            },
//	        }
//	        return Failures[int](contextErrors)
//	    },
//	)
//	// Result contains both original error and context:
//	// ["invalid number", "failed to parse user age"]
//
// Example - Success values pass through:
//
//	result := MonadChainLeft(
//	    Success(42),
//	    func(errs Errors) Validation[int] {
//	        return Failures[int](Errors{
//	            &ValidationError{Messsage: "never called"},
//	        })
//	    },
//	) // Success(42) - unchanged
func MonadChainLeft[A any](fa Validation[A], f Kleisli[Errors, A]) Validation[A] {
	return either.MonadFold(
		fa,
		chainErrors(f),
		Of[A],
	)
}

// Applicative creates an Applicative instance for Validation with error accumulation.
//
// This returns a lawful Applicative that accumulates validation errors using the Errors monoid.
// Unlike the standard Either applicative which fails fast, this validation applicative collects
// all errors when combining independent validations with Ap.
//
// The returned instance satisfies all applicative laws:
//   - Identity: Ap(Of(identity))(v) == v
//   - Homomorphism: Ap(Of(f))(Of(x)) == Of(f(x))
//   - Interchange: Ap(Of(f))(u) == Ap(Map(f => f(y))(u))(Of(y))
//   - Composition: Ap(Ap(Map(compose)(f))(g))(x) == Ap(f)(Ap(g)(x))
//
// Key behaviors:
//   - Of: lifts a value into a successful Validation (Right)
//   - Map: transforms successful values, preserves failures (standard functor)
//   - Ap: when both operands fail, combines all errors using the Errors monoid
//
// This is particularly useful for form validation, configuration validation, and any scenario
// where you want to collect all validation errors at once rather than stopping at the first failure.
//
// Example - Validating Multiple Fields:
//
//	app := Applicative[string, User]()
//
//	// Validate individual fields
//	validateName := func(name string) Validation[string] {
//		if len(name) < 3 {
//			return Failure("Name must be at least 3 characters")
//		}
//		return Success(name)
//	}
//
//	validateAge := func(age int) Validation[int] {
//		if age < 18 {
//			return Failure("Must be 18 or older")
//		}
//		return Success(age)
//	}
//
//	// Create a curried constructor
//	makeUser := func(name string) func(int) User {
//		return func(age int) User {
//			return User{Name: name, Age: age}
//		}
//	}
//
//	// Combine validations - all errors are collected
//	name := validateName("ab")  // Failure: name too short
//	age := validateAge(16)      // Failure: age too low
//
//	result := app.Ap(age)(app.Ap(name)(app.Of(makeUser)))
//	// result contains both validation errors:
//	// - "Name must be at least 3 characters"
//	// - "Must be 18 or older"
//
// Type Parameters:
//   - A: The input value type (Right value)
//   - B: The output value type after transformation
//
// Returns:
//
//	An Applicative instance with Of, Map, and Ap operations that accumulate errors
func Applicative[A, B any]() applicative.Applicative[A, B, Validation[A], Validation[B], Validation[func(A) B]] {
	return either.ApplicativeV[Errors, A, B](
		errorsMonoid,
	)
}

//go:inline
func OrElse[A any](f Kleisli[Errors, A]) Operator[A, A] {
	return ChainLeft(f)
}

// MonadAlt implements the Alternative operation for Validation, providing fallback behavior.
// If the first validation fails, it evaluates and returns the second validation as an alternative.
// If the first validation succeeds, it returns the first validation without evaluating the second.
//
// This is the fundamental operation for the Alt typeclass, enabling "try first, fallback to second"
// semantics. It's particularly useful for:
//   - Providing default values when validation fails
//   - Trying multiple validation strategies in sequence
//   - Building validation pipelines with fallback logic
//   - Implementing optional validation with defaults
//
// **Key behavior**: When both validations fail, MonadAlt DOES accumulate errors from both
// validations using the Errors monoid. This is different from standard Either Alt behavior.
// The error accumulation happens through the underlying ChainLeft/chainErrors mechanism.
//
// The second parameter is lazy (Lazy[Validation[A]]) to avoid unnecessary computation when
// the first validation succeeds. The second validation is only evaluated if needed.
//
// Behavior:
//   - First succeeds: returns first validation (second is not evaluated)
//   - First fails, second succeeds: returns second validation
//   - Both fail: aggregates errors from both validations
//
// This is useful for:
//   - Fallback values: provide defaults when primary validation fails
//   - Alternative strategies: try different validation approaches
//   - Optional validation: make validation optional with a default
//   - Chaining attempts: try multiple sources until one succeeds
//
// Type Parameters:
//   - A: The type of the successful value
//
// Parameters:
//   - first: The primary validation to try
//   - second: A lazy computation producing the fallback validation (only evaluated if first fails)
//
// Returns:
//
//	The first validation if it succeeds, otherwise the second validation
//
// Example - Fallback to default:
//
//	primary := parseConfig("config.json")  // Fails
//	fallback := func() Validation[Config] {
//	    return Success(defaultConfig)
//	}
//	result := MonadAlt(primary, fallback)
//	// Result: Success(defaultConfig)
//
// Example - First succeeds (second not evaluated):
//
//	primary := Success(42)
//	fallback := func() Validation[int] {
//	    panic("never called") // This won't execute
//	}
//	result := MonadAlt(primary, fallback)
//	// Result: Success(42)
//
// Example - Chaining multiple alternatives:
//
//	result := MonadAlt(
//	    parseFromEnv("API_KEY"),
//	    func() Validation[string] {
//	        return MonadAlt(
//	            parseFromFile(".env"),
//	            func() Validation[string] {
//	                return Success("default-key")
//	            },
//	        )
//	    },
//	)
//	// Tries: env var → file → default (uses first that succeeds)
//
// Example - Error accumulation when both fail:
//
//	v1 := Failures[int](Errors{
//	    &ValidationError{Messsage: "error 1"},
//	    &ValidationError{Messsage: "error 2"},
//	})
//	v2 := func() Validation[int] {
//	    return Failures[int](Errors{
//	        &ValidationError{Messsage: "error 3"},
//	    })
//	}
//	result := MonadAlt(v1, v2)
//	// Result: Failures with ALL errors ["error 1", "error 2", "error 3"]
//	// The errors from v1 are aggregated with errors from v2
func MonadAlt[A any](first Validation[A], second Lazy[Validation[A]]) Validation[A] {
	return MonadChainLeft(first, function.Ignore1of1[Errors](second))
}

// Alt is the curried version of [MonadAlt].
// Returns a function that provides fallback behavior for a Validation.
//
// This is useful for creating reusable fallback operators that can be applied
// to multiple validations, or for use in function composition pipelines.
//
// The returned function takes a validation and returns either that validation
// (if successful) or the provided alternative (if the validation fails).
//
// Type Parameters:
//   - A: The type of the successful value
//
// Parameters:
//   - second: A lazy computation producing the fallback validation
//
// Returns:
//
//	A function that takes a Validation[A] and returns a Validation[A] with fallback behavior
//
// Example - Creating a reusable fallback operator:
//
//	withDefault := Alt(func() Validation[int] {
//	    return Success(0)
//	})
//
//	result1 := withDefault(parseNumber("42"))    // Success(42)
//	result2 := withDefault(parseNumber("abc"))   // Success(0) - fallback
//	result3 := withDefault(parseNumber("123"))   // Success(123)
//
// Example - Using in a pipeline:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe2(
//	    parseFromEnv("CONFIG_PATH"),
//	    Alt(func() Validation[string] {
//	        return parseFromFile("config.json")
//	    }),
//	    Alt(func() Validation[string] {
//	        return Success("./default-config.json")
//	    }),
//	)
//	// Tries: env var → file → default path
//
// Example - Combining with Map:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe2(
//	    validatePositive(-5),  // Fails
//	    Alt(func() Validation[int] { return Success(1) }),
//	    Map(func(x int) int { return x * 2 }),
//	)
//	// Result: Success(2) - uses fallback value 1, then doubles it
//
// Example - Multiple fallback layers:
//
//	primaryFallback := Alt(func() Validation[Config] {
//	    return loadFromFile("backup.json")
//	})
//	secondaryFallback := Alt(func() Validation[Config] {
//	    return Success(defaultConfig)
//	})
//
//	result := F.Pipe2(
//	    loadFromFile("config.json"),
//	    primaryFallback,
//	    secondaryFallback,
//	)
//	// Tries: config.json → backup.json → default
func Alt[A any](second Lazy[Validation[A]]) Operator[A, A] {
	return ChainLeft(function.Ignore1of1[Errors](second))
}
