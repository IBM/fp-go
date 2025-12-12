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

package either

//go:inline
func TailRec[E, A, B any](f Kleisli[E, A, Either[A, B]]) Kleisli[E, A, B] {
	return func(a A) Either[E, B] {
		current := f(a)
		for {
			rec, e := Unwrap(current)
			if IsLeft(current) {
				return Left[B](e)
			}
			b, a := Unwrap(rec)
			if IsRight(rec) {
				return Right[E](b)
			}
			current = f(a)
		}
	}
}
