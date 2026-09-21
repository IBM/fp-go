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

package reader

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// idExConfig is the environment used by the Identity examples.
type idExConfig struct {
	Host string
	Port int
}

// idExPort projects the field the examples care about.
func idExPort(c idExConfig) int { return c.Port }

// ExampleIdentity demonstrates the defining property of Identity: it hands the
// environment back unchanged.
func ExampleIdentity() {
	id := Identity[idExConfig]()

	fmt.Println(id(idExConfig{Host: "localhost", Port: 8080}))
	// Output:
	// {localhost 8080}
}

// ExampleIdentity_projection demonstrates that Asks is exactly Identity
// followed by Map, so Identity is the point every projection is built from.
func ExampleIdentity_projection() {
	viaIdentity := F.Pipe1(Identity[idExConfig](), Map[idExConfig](idExPort))
	viaAsks := Asks(idExPort)

	cfg := idExConfig{Host: "localhost", Port: 8080}

	fmt.Println(viaIdentity(cfg), viaAsks(cfg), viaIdentity(cfg) == viaAsks(cfg))
	// Output:
	// 8080 8080 true
}

// ExampleIdentity_neutralElement demonstrates that Identity is the neutral
// element of Reader composition: composing with it on either side changes
// nothing.
func ExampleIdentity_neutralElement() {
	before := F.Flow2(Identity[idExConfig](), idExPort)
	after := F.Flow2(idExPort, Identity[int]())

	cfg := idExConfig{Host: "localhost", Port: 8080}

	fmt.Println(idExPort(cfg), before(cfg), after(cfg))
	// Output:
	// 8080 8080 8080
}
