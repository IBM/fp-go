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

package http

import (
	H "net/http"

	T "github.com/IBM/fp-go/tuple"
)

type (
	// FullResponse represents a full http response, including headers and body
	FullResponse = T.Tuple2[*H.Response, []byte]
)

var (
	Response = T.First[*H.Response, []byte]
	Body     = T.Second[*H.Response, []byte]
)
