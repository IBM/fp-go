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

import "fmt"

type TokenType int

const (
	Identity TokenType = iota // required dependency
	Option                    // optional dependency
	IOEither                  // lazy and required
	IOOption                  // lazy and optional
	Multi                     // array of implementations
	Item                      // item of a multi token
	IOMulti                   // lazy and multi
	Unknown
)

// Dependency describes the relationship to a service
type Dependency interface {
	fmt.Stringer
	// Id returns a unique identifier for a token that can be used as a cache key
	Id() string
	// Type returns a tag that identifies the behaviour of the dependency
	Type() TokenType
}

func AsDependency[T Dependency](t T) Dependency {
	return t
}
