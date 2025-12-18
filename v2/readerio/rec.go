package readerio

// TailRec implements stack-safe tail recursion for the ReaderIO monad.
//
// This function enables recursive computations that depend on an environment (Reader aspect)
// and perform side effects (IO aspect) without risking stack overflow. It uses an iterative
// loop to execute the recursion, making it safe for deep or unbounded recursion.
//
// # How It Works
//
// TailRec takes a Kleisli arrow that returns Trampoline[A, B]:
//   - Bounce(A): Continue recursion with the new state A
//   - Land(B): Terminate recursion and return the final result B
//
// The function iteratively applies the Kleisli arrow, passing the environment R to each
// iteration, until a Land(B) value is produced. This combines:
//   - Environment dependency (Reader monad): Access to configuration, context, or dependencies
//   - Side effects (IO monad): Logging, file I/O, network calls, etc.
//   - Stack safety: Iterative execution prevents stack overflow
//
// # Type Parameters
//
//   - R: The environment type (Reader context) - e.g., Config, Logger, Database connection
//   - A: The state type that changes during recursion
//   - B: The final result type when recursion terminates
//
// # Parameters
//
//   - f: A Kleisli arrow (A => ReaderIO[R, Either[A, B]]) that:
//   - Takes the current state A
//   - Returns a ReaderIO that depends on environment R
//   - Produces Either[A, B] to control recursion flow
//
// # Returns
//
// A Kleisli arrow (A => ReaderIO[R, B]) that:
//   - Takes an initial state A
//   - Returns a ReaderIO that requires environment R
//   - Produces the final result B after recursion completes
//
// # Comparison with Other Monads
//
// Unlike IOEither and IOOption tail recursion:
//   - No error channel (like IOEither's Left error case)
//   - No failure case (like IOOption's None case)
//   - Adds environment dependency that's available throughout recursion
//   - Environment R is passed to every recursive step
//
// # Use Cases
//
//  1. Environment-dependent recursive algorithms:
//     - Recursive computations that need configuration at each step
//     - Algorithms that log progress using an environment-provided logger
//     - Recursive operations that access shared resources from environment
//
//  2. Stateful computations with context:
//     - Tree traversals that need environment context
//     - Graph algorithms with configuration-dependent behavior
//     - Recursive parsers with environment-based rules
//
//  3. Recursive operations with side effects:
//     - File system traversals with logging
//     - Network operations with retry configuration
//     - Database operations with connection pooling
//
// # Example: Factorial with Logging
//
//	type Env struct {
//	    Logger func(string)
//	}
//
//	// Factorial that logs each step
//	factorialStep := func(state struct{ n, acc int }) readerio.ReaderIO[Env, either.Either[struct{ n, acc int }, int]] {
//	    return func(env Env) io.IO[either.Either[struct{ n, acc int }, int]] {
//	        return func() either.Either[struct{ n, acc int }, int] {
//	            if state.n <= 0 {
//	                env.Logger(fmt.Sprintf("Factorial complete: %d", state.acc))
//	                return either.Right[struct{ n, acc int }](state.acc)
//	            }
//	            env.Logger(fmt.Sprintf("Computing: %d * %d", state.n, state.acc))
//	            return either.Left[int](struct{ n, acc int }{state.n - 1, state.acc * state.n})
//	        }
//	    }
//	}
//
//	factorial := readerio.TailRec(factorialStep)
//	env := Env{Logger: func(msg string) { fmt.Println(msg) }}
//	result := factorial(struct{ n, acc int }{5, 1})(env)() // Returns 120, logs each step
//
// # Example: Countdown with Configuration
//
//	type Config struct {
//	    MinValue int
//	    Step     int
//	}
//
//	countdownStep := func(n int) readerio.ReaderIO[Config, either.Either[int, int]] {
//	    return func(cfg Config) io.IO[either.Either[int, int]] {
//	        return func() either.Either[int, int] {
//	            if n <= cfg.MinValue {
//	                return either.Right[int](n)
//	            }
//	            return either.Left[int](n - cfg.Step)
//	        }
//	    }
//	}
//
//	countdown := readerio.TailRec(countdownStep)
//	config := Config{MinValue: 0, Step: 2}
//	result := countdown(10)(config)() // Returns 0 (10 -> 8 -> 6 -> 4 -> 2 -> 0)
//
// # Stack Safety
//
// The iterative implementation ensures that even deeply recursive computations
// (thousands or millions of iterations) will not cause stack overflow:
//
//	// Safe for very large inputs
//	sumToZero := readerio.TailRec(func(n int) readerio.ReaderIO[Env, tailrec.Trampoline[int, int]] {
//	    return func(env Env) io.IO[tailrec.Trampoline[int, int]] {
//	        return func() tailrec.Trampoline[int, int] {
//	            if n <= 0 {
//	                return tailrec.Land[int](0)
//	            }
//	            return tailrec.Bounce[int](n - 1)
//	        }
//	    }
//	})
//	result := sumToZero(1000000)(env)() // Safe, no stack overflow
//
// # Performance Considerations
//
//   - Each iteration creates a new IO action by calling f(a)(r)()
//   - The environment R is passed to every iteration
//   - For performance-critical code, consider if the environment access is necessary
//   - Memoization of environment-derived values may improve performance
//
// # See Also
//
//   - [ioeither.TailRec]: Tail recursion with error handling
//   - [iooption.TailRec]: Tail recursion with optional results
//   - [Chain]: For sequencing ReaderIO computations
//   - [Ask]: For accessing the environment
//   - [Asks]: For extracting values from the environment
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B] {
	return func(a A) ReaderIO[R, B] {
		initialReader := f(a)
		return func(r R) IO[B] {
			initialB := initialReader(r)
			return func() B {
				current := initialB()
				for {
					if current.Landed {
						return current.Land
					}
					current = f(current.Bounce)(r)()
				}
			}
		}
	}
}
