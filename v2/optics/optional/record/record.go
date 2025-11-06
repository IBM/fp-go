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

package record

import (
	OP "github.com/IBM/fp-go/v2/optics/optional"
	G "github.com/IBM/fp-go/v2/optics/optional/record/generic"
)

// FromProperty returns a Optional that gets and sets properties of a map
func AtKey[K comparable, V any](key K) OP.Optional[map[K]V, V] {
	return G.AtKey[map[K]V](key)
}
