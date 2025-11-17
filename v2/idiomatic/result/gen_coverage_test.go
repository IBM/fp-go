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

package result

import (
	"errors"
	"testing"
)

// Test TraverseTuple1
func TestTraverseTuple1(t *testing.T) {
	f := func(x int) (string, error) {
		if x > 0 {
			return Right("positive")
		}
		return Left[string](errors.New("negative"))
	}

	result, err := TraverseTuple1[func(int) (string, error), error](f)(5)
	AssertEq(Right("positive"))(result, err)(t)

	result, err = TraverseTuple1[func(int) (string, error), error](f)(-1)
	AssertEq(Left[string](errors.New("negative")))(result, err)(t)
}

// Test TraverseTuple2
func TestTraverseTuple2(t *testing.T) {
	f1 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f2 := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, err := TraverseTuple2[func(int) (int, error), func(int) (int, error), error](f1, f2)(1, 2)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
}

// Test TraverseTuple3
func TestTraverseTuple3(t *testing.T) {
	f1 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f2 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f3 := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, err := TraverseTuple3[func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f1, f2, f3)(1, 2, 3)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
}

// Test TraverseTuple4
func TestTraverseTuple4(t *testing.T) {
	f1 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f2 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f3 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f4 := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, err := TraverseTuple4[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f1, f2, f3, f4)(1, 2, 3, 4)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
}

// Test TraverseTuple5
func TestTraverseTuple5(t *testing.T) {
	f1 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f2 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f3 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f4 := func(x int) (int, error) {
		return Right(x * 2)
	}
	f5 := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, err := TraverseTuple5[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f1, f2, f3, f4, f5)(1, 2, 3, 4, 5)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
}

// Test TraverseTuple6
func TestTraverseTuple6(t *testing.T) {
	f := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, r6, err := TraverseTuple6[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f, f, f, f, f, f)(1, 2, 3, 4, 5, 6)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
	AssertEq(Right(12))(r6, err)(t)
}

// Test TraverseTuple7
func TestTraverseTuple7(t *testing.T) {
	f := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, r6, r7, err := TraverseTuple7[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f, f, f, f, f, f, f)(1, 2, 3, 4, 5, 6, 7)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
	AssertEq(Right(12))(r6, err)(t)
	AssertEq(Right(14))(r7, err)(t)
}

// Test TraverseTuple8
func TestTraverseTuple8(t *testing.T) {
	f := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, r6, r7, r8, err := TraverseTuple8[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f, f, f, f, f, f, f, f)(1, 2, 3, 4, 5, 6, 7, 8)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
	AssertEq(Right(12))(r6, err)(t)
	AssertEq(Right(14))(r7, err)(t)
	AssertEq(Right(16))(r8, err)(t)
}

// Test TraverseTuple9
func TestTraverseTuple9(t *testing.T) {
	f := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, r6, r7, r8, r9, err := TraverseTuple9[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f, f, f, f, f, f, f, f, f)(1, 2, 3, 4, 5, 6, 7, 8, 9)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
	AssertEq(Right(12))(r6, err)(t)
	AssertEq(Right(14))(r7, err)(t)
	AssertEq(Right(16))(r8, err)(t)
	AssertEq(Right(18))(r9, err)(t)
}

// Test TraverseTuple10
func TestTraverseTuple10(t *testing.T) {
	f := func(x int) (int, error) {
		return Right(x * 2)
	}

	r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, err := TraverseTuple10[func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), func(int) (int, error), error](f, f, f, f, f, f, f, f, f, f)(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	AssertEq(Right(2))(r1, err)(t)
	AssertEq(Right(4))(r2, err)(t)
	AssertEq(Right(6))(r3, err)(t)
	AssertEq(Right(8))(r4, err)(t)
	AssertEq(Right(10))(r5, err)(t)
	AssertEq(Right(12))(r6, err)(t)
	AssertEq(Right(14))(r7, err)(t)
	AssertEq(Right(16))(r8, err)(t)
	AssertEq(Right(18))(r9, err)(t)
	AssertEq(Right(20))(r10, err)(t)
}
