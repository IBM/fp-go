// Copyright (c) 2023 IBM Corp.
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

package either

import (
	"fmt"
)

type (
	// Either defines a data structure that logically holds either an E or an A. The flag discriminates the cases
	Either[E, A any] struct {
		isLeft bool
		left   E
		right  A
	}
)

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	if s.isLeft {
		return fmt.Sprintf("Left[%T, %T](%v)", s.left, s.right, s.left)
	}
	return fmt.Sprintf("Right[%T, %T](%v)", s.left, s.right, s.right)
}

// Format prints some debug info for the object
func (s Either[E, A]) Format(f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, s.String())
	default:
		fmt.Fprint(f, s.String())
	}
}

// IsLeft tests if the [Either] is a left value. Rather use [Fold] if you need to access the values. Inverse is [IsRight].
func IsLeft[E, A any](val Either[E, A]) bool {
	return val.isLeft
}

// IsLeft tests if the [Either] is a right value. Rather use [Fold] if you need to access the values. Inverse is [IsLeft].
func IsRight[E, A any](val Either[E, A]) bool {
	return !val.isLeft
}

// Left creates a new instance of an [Either] representing the left value.
func Left[A, E any](value E) Either[E, A] {
	return Either[E, A]{isLeft: true, left: value}
}

// Right creates a new instance of an [Either] representing the right value.
func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{isLeft: false, right: value}
}

// MonadFold extracts the values from an [Either] by invoking the [onLeft] callback or the [onRight] callback depending on the case
func MonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if ma.isLeft {
		return onLeft(ma.left)
	}
	return onRight(ma.right)
}

// Unwrap converts an [Either] into the idiomatic tuple
func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	return ma.right, ma.left
}
