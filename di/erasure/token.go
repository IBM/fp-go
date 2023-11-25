// Copyright (c) 2023 IBM Corp.
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

package erasure

import (
	"fmt"

	O "github.com/IBM/fp-go/option"
)

const (
	BehaviourMask = 0x0f
	Identity      = 0 // required dependency
	Option        = 1 // optional dependency
	IOEither      = 2 // lazy and required
	IOOption      = 3 // lazy and optional

	TypeMask = 0xf0
	Multi    = 1 << 4 // array of implementations
	Item     = 2 << 4 // item of a multi token
)

// Dependency describes the relationship to a service
type Dependency interface {
	fmt.Stringer
	// Id returns a unique identifier for a token that can be used as a cache key
	Id() string
	// Flag returns a tag that identifies the behaviour of the dependency
	Flag() int
	// ProviderFactory optionally returns an attached [ProviderFactory] that represents the default for this dependency
	ProviderFactory() O.Option[ProviderFactory]
}
