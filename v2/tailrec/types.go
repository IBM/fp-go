package tailrec

type (
	// Trampoline represents a step in a tail-recursive computation.
	//
	// A Trampoline can be in one of two states:
	//   - Bounce: The computation should continue with a new intermediate state (type B)
	//   - Land: The computation is complete with a final result (type L)
	//
	// Type Parameters:
	//   - B: The "bounce" type - intermediate state passed between recursive steps
	//   - L: The "land" type - the final result type when computation completes
	//
	// The trampoline pattern allows converting recursive algorithms into iterative ones,
	// preventing stack overflow for deep recursion while maintaining code clarity.
	//
	// Example:
	//
	//	// Factorial using trampolines
	//	type State struct { n, acc int }
	//
	//	func factorialStep(state State) Trampoline[State, int] {
	//	    if state.n <= 1 {
	//	        return Land[State](state.acc)  // Base case
	//	    }
	//	    return Bounce[int](State{state.n - 1, state.acc * state.n})  // Recursive case
	//	}
	//
	// See package documentation for more examples and usage patterns.
	Trampoline[B, L any] struct {
		// Land holds the final result value when the computation has completed.
		// This field is only meaningful when Landed is true.
		Land L

		// Bounce holds the intermediate state for the next recursive step.
		// This field is only meaningful when Landed is false.
		Bounce B

		// Landed indicates whether the computation has completed.
		// When true, the Land field contains the final result.
		// When false, the Bounce field contains the state for the next iteration.
		Landed bool
	}
)
