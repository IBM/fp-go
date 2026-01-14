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

package file

import "github.com/IBM/fp-go/v2/endomorphism"

type (
	// Endomorphism represents a function from a type to itself: A -> A.
	// This is a type alias for endomorphism.Endomorphism[A].
	//
	// In the context of the file package, this is used for functions that
	// transform strings (paths) into strings (paths), such as the Join function.
	//
	// An endomorphism has useful algebraic properties:
	//   - Identity: There exists an identity endomorphism (the identity function)
	//   - Composition: Endomorphisms can be composed to form new endomorphisms
	//   - Associativity: Composition is associative
	//
	// Example:
	//
	//	import F "github.com/IBM/fp-go/v2/function"
	//
	//	// Join returns an Endomorphism[string]
	//	addConfig := file.Join("config.json")  // Endomorphism[string]
	//	addLogs := file.Join("logs")           // Endomorphism[string]
	//
	//	// Compose endomorphisms
	//	addConfigLogs := F.Flow2(addLogs, addConfig)
	//	result := addConfigLogs("/var")
	//	// result is "/var/logs/config.json"
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
