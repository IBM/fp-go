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

package headers

import (
	"net/http"
	"net/textproto"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	LA "github.com/IBM/fp-go/v2/optics/lens/array"
	LRG "github.com/IBM/fp-go/v2/optics/lens/record/generic"
	RG "github.com/IBM/fp-go/v2/record/generic"
)

// HTTP headers
const (
	Accept        = "Accept"
	Authorization = "Authorization"
	ContentType   = "Content-Type"
	ContentLength = "Content-Length"
)

var (
	// Monoid is a [M.Monoid] to concatenate [http.Header] maps
	Monoid = RG.UnionMonoid[http.Header](A.Semigroup[string]())

	// AtValues is a [L.Lens] that focusses on the values of a header
	AtValues = F.Flow2(
		textproto.CanonicalMIMEHeaderKey,
		LRG.AtRecord[http.Header, []string],
	)

	composeHead = F.Pipe1(
		LA.AtHead[string](),
		L.ComposeOptions[http.Header, string](A.Empty[string]()),
	)

	// AtValue is a [L.Lens] that focusses on first value of a header
	AtValue = F.Flow2(
		AtValues,
		composeHead,
	)
)
