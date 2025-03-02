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

package io

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
)

type (
	ioApplicative[A, B any] struct{}
	IOApplicative[A, B any] = applicative.Applicative[A, B, IO[A], IO[B], IO[func(A) B]]
)

func (o *ioApplicative[A, B]) Of(a A) IO[A] {
	return Of(a)
}

func (o *ioApplicative[A, B]) Map(f func(A) B) func(IO[A]) IO[B] {
	return Map(f)
}

func (o *ioApplicative[A, B]) Ap(fa IO[A]) func(IO[func(A) B]) IO[B] {
	return Ap[B](fa)
}

// Applicative implements the applicativeic operations for [IO]
func Applicative[A, B any]() IOApplicative[A, B] {
	return &ioApplicative[A, B]{}
}
