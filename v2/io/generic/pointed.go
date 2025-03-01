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

package generic

import (
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type ioPointed[A any, GA ~func() A] struct{}

func (o *ioPointed[A, GA]) Of(a A) GA {
	return Of[GA, A](a)
}

// Pointed implements the pointedic operations for [IO]
func Pointed[A any, GA ~func() A]() pointed.Pointed[A, GA] {
	return &ioPointed[A, GA]{}
}
