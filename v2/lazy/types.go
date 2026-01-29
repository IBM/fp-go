package lazy

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Lazy represents a synchronous computation without side effects.
	// It is a function that takes no arguments and returns a value of type A.
	//
	// Lazy computations are evaluated only when their result is needed (lazy evaluation).
	// This allows for:
	//   - Deferring expensive computations until they're actually required
	//   - Creating infinite data structures
	//   - Implementing memoization patterns
	//   - Composing pure computations in a functional style
	//
	// Example:
	//
	//	// Create a lazy computation
	//	computation := lazy.Of(42)
	//
	//	// Transform it (not evaluated yet)
	//	doubled := lazy.Map(N.Mul(2))(computation)
	//
	//	// Evaluate when needed
	//	result := doubled() // 84
	//
	// Note: Lazy is an alias for io.IO[A] but represents pure computations
	// without side effects, whereas IO represents computations that may have side effects.
	Lazy[A any] = func() A

	// Kleisli represents a function that takes a value of type A and returns
	// a lazy computation producing a value of type B.
	//
	// Kleisli arrows are used for composing monadic computations. They allow
	// you to chain operations where each step depends on the result of the previous step.
	//
	// Example:
	//
	//	// A Kleisli arrow that doubles a number lazily
	//	double := func(x int) lazy.Lazy[int] {
	//	    return lazy.Of(x * 2)
	//	}
	//
	//	// Chain it with another operation
	//	result := lazy.Chain(double)(lazy.Of(5))() // 10
	Kleisli[A, B any] = func(A) Lazy[B]

	// Operator represents a function that takes a lazy computation of type A
	// and returns a lazy computation of type B.
	//
	// Operators are used to transform lazy computations. They are essentially
	// Kleisli arrows where the input is already wrapped in a Lazy context.
	//
	// Example:
	//
	//	// An operator that doubles the value in a lazy computation
	//	doubleOp := lazy.Map(N.Mul(2))
	//
	//	// Apply it to a lazy computation
	//	result := doubleOp(lazy.Of(5))() // 10
	Operator[A, B any] = Kleisli[Lazy[A], B]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
