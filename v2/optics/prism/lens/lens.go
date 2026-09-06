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

// Package lens provides utilities for composing prisms with lenses to create optionals.
//
// This package enables composition of a Prism (which focuses on a variant within a sum type)
// with a Lens (which focuses on a field within a product type) to create an Optional that
// combines both focusing operations.
package lens

import "github.com/IBM/fp-go/v2/internal/common"

// Compose composes a Prism with a Lens to create an Optional.
//
// Deprecated: Use github.com/IBM/fp-go/v2/optics/prism.ComposeLens instead.
func Compose[S, A, B any](l Lens[A, B]) func(Prism[S, A]) Optional[S, B] {
	return common.PrismComposeLens[S](l)
}
