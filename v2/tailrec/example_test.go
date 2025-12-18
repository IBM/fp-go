// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tailrec_test

import (
	"fmt"

	"github.com/IBM/fp-go/v2/tailrec"
)

// ExampleBounce demonstrates creating a trampoline that continues computation.
func ExampleBounce() {
	// Create a bounce trampoline with value 42
	tramp := tailrec.Bounce[string](42)

	// Access fields directly to inspect the state
	fmt.Printf("Is computation complete? %v\n", tramp.Landed)
	fmt.Printf("Next value to process: %d\n", tramp.Bounce)

	// Output:
	// Is computation complete? false
	// Next value to process: 42
}

// ExampleLand demonstrates creating a trampoline that completes computation.
func ExampleLand() {
	// Create a land trampoline with final result
	tramp := tailrec.Land[int]("done")

	// Access fields directly to inspect the state
	fmt.Printf("Is computation complete? %v\n", tramp.Landed)
	fmt.Printf("Final result: %s\n", tramp.Land)

	// Output:
	// Is computation complete? true
	// Final result: done
}

// Example_fieldAccess demonstrates accessing trampoline fields directly.
func Example_fieldAccess() {
	// Create a bounce trampoline
	bounceTramp := tailrec.Bounce[string](42)
	fmt.Printf("Bounce: value=%d, landed=%v\n", bounceTramp.Bounce, bounceTramp.Landed)

	// Create a land trampoline
	landTramp := tailrec.Land[int]("result")
	fmt.Printf("Land: value=%s, landed=%v\n", landTramp.Land, landTramp.Landed)

	// Output:
	// Bounce: value=42, landed=false
	// Land: value=result, landed=true
}

// Example_factorial demonstrates implementing factorial using trampolines.
func Example_factorial() {
	type State struct {
		n   int
		acc int
	}

	// Define the step function
	factorialStep := func(state State) tailrec.Trampoline[State, int] {
		if state.n <= 1 {
			return tailrec.Land[State](state.acc)
		}
		return tailrec.Bounce[int](State{state.n - 1, state.acc * state.n})
	}

	// Execute the trampoline
	factorial := func(n int) int {
		current := tailrec.Bounce[int](State{n, 1})
		for {
			if current.Landed {
				return current.Land
			}
			current = factorialStep(current.Bounce)
		}
	}

	// Calculate factorial of 5
	result := factorial(5)
	fmt.Printf("5! = %d\n", result)

	// Output:
	// 5! = 120
}

// Example_fibonacci demonstrates computing Fibonacci numbers using trampolines.
func Example_fibonacci() {
	type State struct {
		n    int
		curr int
		prev int
	}

	// Define the step function
	fibStep := func(state State) tailrec.Trampoline[State, int] {
		if state.n <= 0 {
			return tailrec.Land[State](state.curr)
		}
		return tailrec.Bounce[int](State{
			n:    state.n - 1,
			curr: state.prev + state.curr,
			prev: state.curr,
		})
	}

	// Execute the trampoline
	fibonacci := func(n int) int {
		current := tailrec.Bounce[int](State{n, 1, 0})
		for {
			if current.Landed {
				return current.Land
			}
			current = fibStep(current.Bounce)
		}
	}

	// Calculate 10th Fibonacci number
	result := fibonacci(10)
	fmt.Printf("fib(10) = %d\n", result)

	// Output:
	// fib(10) = 89
}

// Example_sumList demonstrates processing a list using trampolines.
func Example_sumList() {
	type State struct {
		list []int
		sum  int
	}

	// Define the step function
	sumStep := func(state State) tailrec.Trampoline[State, int] {
		if len(state.list) == 0 {
			return tailrec.Land[State](state.sum)
		}
		return tailrec.Bounce[int](State{
			list: state.list[1:],
			sum:  state.sum + state.list[0],
		})
	}

	// Execute the trampoline
	sumList := func(list []int) int {
		current := tailrec.Bounce[int](State{list, 0})
		for {
			if current.Landed {
				return current.Land
			}
			current = sumStep(current.Bounce)
		}
	}

	// Sum a list of numbers
	numbers := []int{1, 2, 3, 4, 5}
	result := sumList(numbers)
	fmt.Printf("sum([1,2,3,4,5]) = %d\n", result)

	// Output:
	// sum([1,2,3,4,5]) = 15
}

// Example_countdown demonstrates a simple countdown using trampolines.
func Example_countdown() {
	// Define the step function
	countdownStep := func(n int) tailrec.Trampoline[int, int] {
		if n <= 0 {
			return tailrec.Land[int](0)
		}
		return tailrec.Bounce[int](n - 1)
	}

	// Execute the trampoline
	countdown := func(n int) int {
		current := tailrec.Bounce[int](n)
		for {
			if current.Landed {
				return current.Land
			}
			current = countdownStep(current.Bounce)
		}
	}

	// Countdown from 5
	result := countdown(5)
	fmt.Printf("countdown(5) = %d\n", result)

	// Output:
	// countdown(5) = 0
}

// Example_gcd demonstrates computing greatest common divisor using trampolines.
func Example_gcd() {
	type State struct {
		a int
		b int
	}

	// Define the step function (Euclidean algorithm)
	gcdStep := func(state State) tailrec.Trampoline[State, int] {
		if state.b == 0 {
			return tailrec.Land[State](state.a)
		}
		return tailrec.Bounce[int](State{state.b, state.a % state.b})
	}

	// Execute the trampoline
	gcd := func(a, b int) int {
		current := tailrec.Bounce[int](State{a, b})
		for {
			if current.Landed {
				return current.Land
			}
			current = gcdStep(current.Bounce)
		}
	}

	// Calculate GCD of 48 and 18
	result := gcd(48, 18)
	fmt.Printf("gcd(48, 18) = %d\n", result)

	// Output:
	// gcd(48, 18) = 6
}

// Example_deepRecursion demonstrates handling deep recursion without stack overflow.
func Example_deepRecursion() {
	// Define the step function
	countdownStep := func(n int) tailrec.Trampoline[int, int] {
		if n <= 0 {
			return tailrec.Land[int](0)
		}
		return tailrec.Bounce[int](n - 1)
	}

	// Execute the trampoline
	countdown := func(n int) int {
		current := tailrec.Bounce[int](n)
		for {
			if current.Landed {
				return current.Land
			}
			current = countdownStep(current.Bounce)
		}
	}

	// This would cause stack overflow with regular recursion
	// but works fine with trampolines
	result := countdown(100000)
	fmt.Printf("countdown(100000) = %d (no stack overflow!)\n", result)

	// Output:
	// countdown(100000) = 0 (no stack overflow!)
}

// Example_collatz demonstrates the Collatz conjecture using trampolines.
func Example_collatz() {
	// Define the step function
	collatzStep := func(n int) tailrec.Trampoline[int, int] {
		if n <= 1 {
			return tailrec.Land[int](n)
		}
		if n%2 == 0 {
			return tailrec.Bounce[int](n / 2)
		}
		return tailrec.Bounce[int](3*n + 1)
	}

	// Execute the trampoline and count steps
	collatzSteps := func(n int) int {
		current := tailrec.Bounce[int](n)
		steps := 0
		for {
			if current.Landed {
				return steps
			}
			current = collatzStep(current.Bounce)
			steps++
		}
	}

	// Count steps for Collatz sequence starting at 27
	result := collatzSteps(27)
	fmt.Printf("Collatz(27) takes %d steps to reach 1\n", result)

	// Output:
	// Collatz(27) takes 112 steps to reach 1
}
