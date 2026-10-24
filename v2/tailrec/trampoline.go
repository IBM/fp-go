package tailrec

import "fmt"

// Bounce creates a Trampoline that indicates the computation should continue
// with a new intermediate state.
//
// This represents a recursive call in the original algorithm. The computation
// will continue by processing the provided state value in the next iteration.
//
// # Type Parameters
//
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// # Parameters
//
//   - b: The new intermediate state to process in the next step
//
// # Returns
//
//   - Trampoline[B, L]: A Trampoline in the "bounce" state containing the intermediate value
//
// # Example
//
//	// Countdown that bounces until reaching zero
//	func countdownStep(n int) Trampoline[int, int] {
//	    if n <= 0 {
//	        return Land[int](0)
//	    }
//	    return Bounce[int](n - 1)  // Continue with n-1
//	}
//
//go:inline
func Bounce[L, B any](b B) Trampoline[B, L] {
	return Trampoline[B, L]{Bounce: b}
}

// Land creates a Trampoline that indicates the computation is complete
// with a final result.
//
// This represents the base case in the original recursive algorithm. When
// a Land trampoline is encountered, the executor should stop iterating and
// return the final result.
//
// # Type Parameters
//
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// # Parameters
//
//   - l: The final result value
//
// # Returns
//
//   - Trampoline[B, L]: A Trampoline in the "land" state containing the final result
//
// # Example
//
//	// Factorial base case
//	func factorialStep(state State) Trampoline[State, int] {
//	    if state.n <= 1 {
//	        return Land[State](state.acc)  // Computation complete
//	    }
//	    return Bounce[int](State{state.n - 1, state.acc * state.n})
//	}
//
//go:inline
func Land[B, L any](l L) Trampoline[B, L] {
	return Trampoline[B, L]{Land: l, Landed: true}
}

// String implements fmt.Stringer for Trampoline.
//
// Returns a human-readable string representation of the trampoline state.
// For bounce states, returns "Bounce(value)". For land states, returns "Land(value)".
// The exact format is not a stable contract and may change across versions.
func (t Trampoline[B, L]) String() string {
	if t.Landed {
		return fmt.Sprintf("Land(%v)", t.Land)
	}
	return fmt.Sprintf("Bounce(%v)", t.Bounce)
}

// Format implements fmt.Formatter for Trampoline.
//
// Supports various formatting verbs for detailed output:
//   - %v: Default format (delegates to String)
//   - %+v: Detailed format with type information
//   - %s: String format
//   - %q: Quoted string format
//
// The exact output format is not a stable contract and may change across versions.
func (t Trampoline[B, L]) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			// %+v: detailed format with type information
			if t.Landed {
				fmt.Fprintf(f, "Trampoline[Land]{Land: %+v, Landed: true}", t.Land)
			} else {
				fmt.Fprintf(f, "Trampoline[Bounce]{Bounce: %+v, Landed: false}", t.Bounce)
			}
		} else {
			// %v: default format (delegates to String)
			fmt.Fprint(f, t.String())
		}
	case 's':
		// %s: string format
		fmt.Fprint(f, t.String())
	case 'q':
		// %q: quoted string format
		fmt.Fprintf(f, "%q", t.String())
	default:
		// Unknown verb: print with %!verb notation
		fmt.Fprintf(f, "%%!%c(Trampoline[B, L]=%s)", verb, t.String())
	}
}

// IsLand reports whether the Trampoline has landed, i.e. the computation is
// complete and carries a final result.
//
// Returns:
//   - bool: true if the Trampoline is in the "land" state
//
// See Also:
//   - IsBounce: the complementary check
//   - GetLand: extracts the final result
//
//go:inline
func (t Trampoline[B, L]) IsLand() bool {
	return t.Landed
}

// IsBounce reports whether the Trampoline bounces, i.e. the computation
// continues with an intermediate state.
//
// Returns:
//   - bool: true if the Trampoline is in the "bounce" state
//
// See Also:
//   - IsLand: the complementary check
//   - GetBounce: extracts the intermediate state
//
//go:inline
func (t Trampoline[B, L]) IsBounce() bool {
	return !t.Landed
}

// GetLand returns the final result of a landed Trampoline as an idiomatic
// comma-ok tuple.
//
// Returns:
//   - L: The final result, or the zero value of L if the Trampoline bounces
//   - bool: true if the Trampoline is in the "land" state
//
// See Also:
//   - GetBounce: extracts the intermediate state
//   - IsLand: checks the state without extracting the value
//
//go:inline
func (t Trampoline[B, L]) GetLand() (L, bool) {
	if t.Landed {
		return t.Land, true
	}
	var l L
	return l, false
}

// GetBounce returns the intermediate state of a bouncing Trampoline as an
// idiomatic comma-ok tuple.
//
// Returns:
//   - B: The intermediate state, or the zero value of B if the Trampoline has landed
//   - bool: true if the Trampoline is in the "bounce" state
//
// See Also:
//   - GetLand: extracts the final result
//   - IsBounce: checks the state without extracting the value
//
//go:inline
func (t Trampoline[B, L]) GetBounce() (B, bool) {
	if t.Landed {
		var b B
		return b, false
	}
	return t.Bounce, true
}

// IsLand reports whether a Trampoline has landed, i.e. the computation is
// complete and carries a final result.
//
// It is the standalone form of Trampoline.IsLand and can be used as a predicate.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Parameters:
//   - t: The Trampoline to inspect
//
// Returns:
//   - bool: true if t is in the "land" state
//
// See Also:
//   - IsBounce: the complementary check
//   - GetLand: extracts the final result
//
//go:inline
func IsLand[B, L any](t Trampoline[B, L]) bool {
	return t.IsLand()
}

// IsBounce reports whether a Trampoline bounces, i.e. the computation
// continues with an intermediate state.
//
// It is the standalone form of Trampoline.IsBounce and can be used as a predicate.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Parameters:
//   - t: The Trampoline to inspect
//
// Returns:
//   - bool: true if t is in the "bounce" state
//
// See Also:
//   - IsLand: the complementary check
//   - GetBounce: extracts the intermediate state
//
//go:inline
func IsBounce[B, L any](t Trampoline[B, L]) bool {
	return t.IsBounce()
}

// GetLand returns the final result of a landed Trampoline as an idiomatic
// comma-ok tuple.
//
// It is the standalone form of Trampoline.GetLand. Combine it with
// option.FromValidation to obtain an Option.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Parameters:
//   - t: The Trampoline to inspect
//
// Returns:
//   - L: The final result, or the zero value of L if t bounces
//   - bool: true if t is in the "land" state
//
// See Also:
//   - GetBounce: extracts the intermediate state
//   - Land: the inverse, constructs a landed Trampoline
//
//go:inline
func GetLand[B, L any](t Trampoline[B, L]) (L, bool) {
	return t.GetLand()
}

// GetBounce returns the intermediate state of a bouncing Trampoline as an
// idiomatic comma-ok tuple.
//
// It is the standalone form of Trampoline.GetBounce. Combine it with
// option.FromValidation to obtain an Option.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Parameters:
//   - t: The Trampoline to inspect
//
// Returns:
//   - B: The intermediate state, or the zero value of B if t has landed
//   - bool: true if t is in the "bounce" state
//
// See Also:
//   - GetLand: extracts the final result
//   - Bounce: the inverse, constructs a bouncing Trampoline
//
//go:inline
func GetBounce[B, L any](t Trampoline[B, L]) (B, bool) {
	return t.GetBounce()
}
