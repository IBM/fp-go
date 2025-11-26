// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"iter"

	"github.com/IBM/fp-go/v2/endomorphism"
)

type (
	// Seq is an iterator sequence type alias for working with Go 1.23+ iterators.
	Seq[T any] = iter.Seq[T]

	// Endomorphism represents a function from type T to type T.
	// It is commonly used for transformations that preserve the type.
	Endomorphism[T any] = endomorphism.Endomorphism[T]
)
