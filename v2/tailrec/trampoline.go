package tailrec

import "fmt"

// Bounce creates a Trampoline that indicates the computation should continue
// with a new intermediate state.
//
// This represents a recursive call in the original algorithm. The computation
// will continue by processing the provided state value in the next iteration.
//
// Type Parameters:
//   - L: The final result type (land type)
//   - B: The intermediate state type (bounce type)
//
// Parameters:
//   - b: The new intermediate state to process in the next step
//
// Returns:
//   - A Trampoline in the "bounce" state containing the intermediate value
//
// Example:
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
	return Trampoline[B, L]{Bounce: b, Landed: false}
}

// Land creates a Trampoline that indicates the computation is complete
// with a final result.
//
// This represents the base case in the original recursive algorithm. When
// a Land trampoline is encountered, the executor should stop iterating and
// return the final result.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Parameters:
//   - l: The final result value
//
// Returns:
//   - A Trampoline in the "land" state containing the final result
//
// Example:
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
// Returns a human-readable string representation of the trampoline state.
func (t Trampoline[B, L]) String() string {
	if t.Landed {
		return fmt.Sprintf("Land(%v)", t.Land)
	}
	return fmt.Sprintf("Bounce(%v)", t.Bounce)
}

// Format implements fmt.Formatter for Trampoline.
// Supports various formatting verbs for detailed output.
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
		} else if f.Flag('#') {
			// %#v: Go-syntax representation (delegates to GoString)
			fmt.Fprint(f, t.GoString())
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

// GoString implements fmt.GoStringer for Trampoline.
// Returns a Go-syntax representation that could be used to recreate the value.
func (t Trampoline[B, L]) GoString() string {
	if t.Landed {
		return fmt.Sprintf("tailrec.Land[%T](%#v)", t.Bounce, t.Land)
	}
	return fmt.Sprintf("tailrec.Bounce[%T](%#v)", t.Land, t.Bounce)
}
