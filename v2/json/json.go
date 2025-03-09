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

package json

import (
	"encoding/json"

	E "github.com/IBM/fp-go/v2/either"
)

// Unmarshal parses a JSON data structure from bytes
func Unmarshal[A any](data []byte) Either[A] {
	var result A
	err := json.Unmarshal(data, &result)
	return E.TryCatchError(result, err)
}

// Marshal converts a data structure to json
func Marshal[A any](a A) Either[[]byte] {
	return E.TryCatchError(json.Marshal(a))
}
