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
//go:build either_pointers

package either

import "fmt"

type Either[E, A any] struct {
	left  *E
	right *A
}

// String prints some debug info for the object
//
//go:noinline
func eitherString[E, A any](s *Either[E, A]) string {
	if s.right != nil {
		return fmt.Sprintf("Right[%T](%v)", *s.right, *s.right)
	}
	return fmt.Sprintf("Left[%T](%v)", *s.left, *s.left)
}

// Format prints some debug info for the object
//
//go:noinline
func eitherFormat[E, A any](e *Either[E, A], f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, eitherString(e))
	default:
		fmt.Fprint(f, eitherString(e))
	}
}

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	return eitherString(&s)
}

// Format prints some debug info for the object
func (s Either[E, A]) Format(f fmt.State, c rune) {
	eitherFormat(&s, f, c)
}

//go:inline
func Left[A, E any](value E) Either[E, A] {
	return Either[E, A]{left: &value}
}

//go:inline
func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{right: &value}
}

//go:inline
func IsLeft[E, A any](e Either[E, A]) bool {
	return e.left != nil
}

//go:inline
func IsRight[E, A any](e Either[E, A]) bool {
	return e.right != nil
}

//go:inline
func MonadFold[E, A, B any](ma Either[E, A], onLeft func(E) B, onRight func(A) B) B {
	if ma.right != nil {
		return onRight(*ma.right)
	}
	return onLeft(*ma.left)
}

//go:inline
func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	if ma.right != nil {
		var e E
		return *ma.right, e
	}
	var a A
	return a, *ma.left
}
