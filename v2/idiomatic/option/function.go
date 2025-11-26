// Copyright (c) 2025 IBM Corp.
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

package option

// Pipe1 takes an initial value t0 and successively applies 1 function where the input of a function is the return value of the previous function.
// The final return value is the result of the last function application.
//
// Example:
//
//	result := Pipe1(42, func(x int) (int, bool) { return x * 2, true }) // (84, true)
//
//go:inline
func Pipe1[F1 ~func(T0) (T1, bool), T0, T1 any](t0 T0, f1 F1) (T1, bool) {
	return f1(t0)
}

// Flow1 creates a function that takes an initial value t0 and successively applies 1 function where the input of a function is the return value of the previous function.
// The final return value is the result of the last function application.
//
// Example:
//
//	double := Flow1(func(x int, ok bool) (int, bool) { return x * 2, ok })
//	result := double(42, true) // (84, true)
//
//go:inline
func Flow1[F1 ~func(T0, bool) (T1, bool), T0, T1 any](f1 F1) func(T0, bool) (T1, bool) {
	return f1
}

// Pipe2 takes an initial value t0 and successively applies 2 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Pipe2[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), T0, T1, T2 any](t0 T0, f1 F1, f2 F2) (T2, bool) {
	return f2(f1(t0))
}

// Flow2 creates a function that takes an initial value t0 and successively applies 2 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Flow2[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), T0, T1, T2 any](f1 F1, f2 F2) func(T0, bool) (T2, bool) {
	return func(t0 T0, t0ok bool) (T2, bool) {
		return f2(f1(t0, t0ok))
	}
}

// Pipe3 takes an initial value t0 and successively applies 3 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Pipe3[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), T0, T1, T2, T3 any](t0 T0, f1 F1, f2 F2, f3 F3) (T3, bool) {
	return f3(f2(f1(t0)))
}

// Flow3 creates a function that takes an initial value t0 and successively applies 3 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Flow3[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), T0, T1, T2, T3 any](f1 F1, f2 F2, f3 F3) func(T0, bool) (T3, bool) {
	return func(t0 T0, t0ok bool) (T3, bool) {
		return f3(f2(f1(t0, t0ok)))
	}
}

// Pipe4 takes an initial value t0 and successively applies 4 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Pipe4[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), T0, T1, T2, T3, T4 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4) (T4, bool) {
	return f4(f3(f2(f1(t0))))
}

// Flow4 creates a function that takes an initial value t0 and successively applies 4 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Flow4[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), T0, T1, T2, T3, T4 any](f1 F1, f2 F2, f3 F3, f4 F4) func(T0, bool) (T4, bool) {
	return func(t0 T0, t0ok bool) (T4, bool) {
		return f4(f3(f2(f1(t0, t0ok))))
	}
}

// Pipe5 takes an initial value t0 and successively applies 5 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Pipe5[F1 ~func(T0) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), F5 ~func(T4, bool) (T5, bool), T0, T1, T2, T3, T4, T5 any](t0 T0, f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) (T5, bool) {
	return f5(f4(f3(f2(f1(t0)))))
}

// Flow5 creates a function that takes an initial value t0 and successively applies 5 functions where the input of a function is the return value of the previous function
// The final return value is the result of the last function application
//
//go:inline
func Flow5[F1 ~func(T0, bool) (T1, bool), F2 ~func(T1, bool) (T2, bool), F3 ~func(T2, bool) (T3, bool), F4 ~func(T3, bool) (T4, bool), F5 ~func(T4, bool) (T5, bool), T0, T1, T2, T3, T4, T5 any](f1 F1, f2 F2, f3 F3, f4 F4, f5 F5) func(T0, bool) (T5, bool) {
	return func(t0 T0, t0ok bool) (T5, bool) {
		return f5(f4(f3(f2(f1(t0, t0ok)))))
	}
}
