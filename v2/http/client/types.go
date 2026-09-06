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

package client

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/common"
)

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	//
	// Used as the return type of Lens.Set: applying a lens setter produces an
	// Endomorphism that can be composed with other transformations via
	// endomorphism.Concat.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lens is an optic that focuses on a single field A inside a larger
	// structure S. It provides two operations:
	//
	//   - Get(s S) A       — extract the focused field
	//   - Set(a A) func(S) S — return a copy of s with the field replaced by a
	//
	// Both lenses in this package use *http.Client as S, so Set always copies
	// the pointed-to struct before mutating it.
	Lens[S, A any] = common.Lens[S, A]
)
