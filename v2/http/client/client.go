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

// Package client provides functional lenses for reading and updating fields of
// an *http.Client in an immutable, composable way.
//
// Each lens follows the standard Lens[S, A] contract:
//   - Get: extracts the focused field from the source
//   - Set: returns a copy of the source with the focused field replaced
//
// Both lenses in this package operate on *http.Client (a pointer type).
// The Set operation always copies the struct before mutating it, so the
// original pointer is never modified.
//
// Lenses:
//
//	TimeoutLens()   — focuses on the Timeout field (time.Duration)
//	TransportLens() — focuses on the Transport field (http.RoundTripper)
//
// Example — compose two lenses to build a configured client:
//
//	import (
//	    C  "github.com/IBM/fp-go/v2/http/client"
//	    EN "github.com/IBM/fp-go/v2/endomorphism"
//	)
//
//	configure := EN.Concat(
//	    C.TimeoutLens().Set(30 * time.Second),
//	    C.TransportLens().Set(myTransport),
//	)
//	client := configure(&http.Client{})
package client

import (
	"net/http"
	"time"

	"github.com/IBM/fp-go/v2/internal/common"
)

// TimeoutLens returns a Lens that focuses on the Timeout field of an *http.Client.
//
// The lens follows the standard immutable contract: Set copies the *http.Client
// before updating Timeout, so the original pointer is not modified.
//
// Returns:
//   - Lens[*http.Client, time.Duration]: a lens for the Timeout field
func TimeoutLens() Lens[*http.Client, time.Duration] {
	return common.MakeLensRefWithName(
		func(c *http.Client) time.Duration {
			return c.Timeout
		},
		func(c *http.Client, t time.Duration) *http.Client {
			c.Timeout = t
			return c
		},
		"HttpTimeout",
	)
}

// TransportLens returns a Lens that focuses on the Transport field of an *http.Client.
//
// The lens follows the standard immutable contract: Set copies the *http.Client
// before updating Transport, so the original pointer is not modified.
//
// Returns:
//   - Lens[*http.Client, http.RoundTripper]: a lens for the Transport field
func TransportLens() Lens[*http.Client, http.RoundTripper] {
	return common.MakeLensRefWithName(
		func(c *http.Client) http.RoundTripper {
			return c.Transport
		},
		func(c *http.Client, t http.RoundTripper) *http.Client {
			c.Transport = t
			return c
		},
		"HttpTransport",
	)
}
