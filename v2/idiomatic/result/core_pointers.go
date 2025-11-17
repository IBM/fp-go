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

package result

import "fmt"

type Either[A any] struct {
	left  *E
	right *A
}

// String prints some debug info for the object
//
//go:noinline
func eitherString[A any](s *Either[A]) string {
	if s.right != nil {
		return fmt.Sprintf("Right[%T](%v)", *s.right, *s.right)
	}
	return fmt.Sprintf("Left[%T](%v)", *s.left, *s.left)
}

// Format prints some debug info for the object
//
//go:noinline
func eitherFormat[A any](e *Either[A], f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, eitherString(e))
	default:
		fmt.Fprint(f, eitherString(e))
	}
}

// String prints some debug info for the object
func (s Either[A]) String() string {
	return eitherString(&s)
}

// Format prints some debug info for the object
func (s Either[A]) Format(f fmt.State, c rune) {
	eitherFormat(&s, f, c)
}

//go:inline
func Left[A, E any](value E) Either[A] {
	return Either[A]{left: &value}
}

//go:inline
func Right[A any](value A) Either[A] {
	return Either[A]{right: &value}
}

//go:inline
func IsLeft[A any](e Either[A]) bool {
	return e.left != nil
}

//go:inline
func IsRight[A any](e Either[A]) bool {
	return e.right != nil
}

//go:inline
func MonadFold[A, B any](ma Either[A], onLeft func(E) B, onRight func(A) B) B {
	if ma.right != nil {
		return onRight(*ma.right)
	}
	return onLeft(*ma.left)
}

//go:inline
func Unwrap[A any](ma Either[A]) (A, E) {
	if ma.right != nil {
		var e E
		return *ma.right, e
	}
	var a A
	return a, *ma.left
}
