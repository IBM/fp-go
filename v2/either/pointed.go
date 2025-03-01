// Copyright (c) 2024 IBM Corp.
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
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type eitherPointed[E, A any] struct{}

func (o *eitherPointed[E, A]) Of(a A) Either[E, A] {
	return Of[E, A](a)
}

// Pointed implements the pointedic operations for [Either]
func Pointed[E, A any]() pointed.Pointed[A, Either[E, A]] {
	return &eitherPointed[E, A]{}
}
