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

package identity

type (
	// Kleisli represents a Kleisli arrow for the Identity monad.
	// It's simply a function from A to B, as Identity has no computational context.
	Kleisli[A, B any] = func(A) B

	// Operator represents a function that transforms values.
	// In the Identity monad, it's equivalent to Kleisli since there's no wrapping context.
	Operator[A, B any] = Kleisli[A, B]
)
