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

package client_test

import (
	"fmt"
	"net/http"
	"time"

	C "github.com/IBM/fp-go/v2/http/client"
)

// ExampleTimeoutLens_get demonstrates reading the Timeout from an *http.Client.
func ExampleTimeoutLens_get() {
	lens := C.TimeoutLens()
	client := &http.Client{Timeout: 30 * time.Second}

	timeout := lens.Get(client)
	fmt.Println(timeout)
	// Output: 30s
}

// ExampleTimeoutLens_set demonstrates setting the Timeout on an *http.Client.
// Set always returns a new copy; the original is not modified.
func ExampleTimeoutLens_set() {
	lens := C.TimeoutLens()
	original := &http.Client{}

	updated := lens.Set(10 * time.Second)(original)

	fmt.Println(original.Timeout) // original unchanged
	fmt.Println(updated.Timeout)  // new copy has the updated value
	// Output:
	// 0s
	// 10s
}

// ExampleTransportLens_get demonstrates reading the Transport from an *http.Client.
func ExampleTransportLens_get() {
	lens := C.TransportLens()
	client := &http.Client{Transport: http.DefaultTransport}

	transport := lens.Get(client)
	fmt.Println(transport != nil)
	// Output: true
}

// ExampleTransportLens_set demonstrates setting the Transport on an *http.Client.
// Set always returns a new copy; the original is not modified.
func ExampleTransportLens_set() {
	lens := C.TransportLens()
	original := &http.Client{}

	updated := lens.Set(http.DefaultTransport)(original)

	fmt.Println(original.Transport == nil) // original unchanged
	fmt.Println(updated.Transport != nil)  // new copy has the transport set
	// Output:
	// true
	// true
}

// ExampleTimeoutLens_compose demonstrates composing both lenses to configure a
// client in a single pipeline.
func ExampleTimeoutLens_compose() {
	tl := C.TimeoutLens()
	trl := C.TransportLens()

	configure := func(c *http.Client) *http.Client {
		return trl.Set(http.DefaultTransport)(tl.Set(5 * time.Second)(c))
	}

	client := configure(&http.Client{})

	fmt.Println(client.Timeout)
	fmt.Println(client.Transport != nil)
	// Output:
	// 5s
	// true
}
