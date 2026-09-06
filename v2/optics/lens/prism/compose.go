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

package prism

import (
	"github.com/IBM/fp-go/v2/internal/common"
)

// Compose composes a Lens with a Prism to create an Optional.
//
// Deprecated: Use github.com/IBM/fp-go/v2/optics/lens.ComposePrism instead.
func Compose[S, A, B any](p Prism[A, B]) func(Lens[S, A]) Optional[S, B] {
	return common.LensComposePrism[S](p)
}

// ComposeRef composes a Lens operating on pointer types with a Prism to create an Optional.
func ComposeRef[S, A, B any](p Prism[A, B]) func(Lens[*S, A]) Optional[*S, B] {
	return common.LensComposePrismRef[S](p)
}
