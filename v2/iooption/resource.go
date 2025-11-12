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

package iooption

import "github.com/IBM/fp-go/v2/function"

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	R, A, ANY any](onCreate IOOption[R], onRelease Kleisli[R, ANY]) Kleisli[Kleisli[R, A], A] {
	// simply map to implementation of bracket
	return function.Bind13of3(Bracket[R, A, ANY])(onCreate, function.Ignore2of2[Option[A]](onRelease))
}
