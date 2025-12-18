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

package readereither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/tailrec"
)

//go:inline
func TailRec[R, E, A, B any](f Kleisli[R, E, A, tailrec.Trampoline[A, B]]) Kleisli[R, E, A, B] {
	return func(a A) ReaderEither[R, E, B] {
		initialReader := f(a)
		return func(r R) either.Either[E, B] {
			current := initialReader(r)
			for {
				rec, e := either.Unwrap(current)
				if either.IsLeft(current) {
					return either.Left[B](e)
				}
				if rec.Landed {
					return either.Right[E](rec.Land)
				}
				current = f(rec.Bounce)(r)
			}
		}
	}
}
