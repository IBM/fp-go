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

package readerreaderioresult

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// idExConfig is the outer environment used by the Identity examples.
type idExConfig struct {
	Host string
	Port int
}

// idExPort projects the field the examples care about.
func idExPort(c idExConfig) int { return c.Port }

// idExWide is a wider environment the narrow one is projected out of.
type idExWide struct {
	Cfg  idExConfig
	User string
}

// idExNarrow projects the wide environment onto the narrow one.
func idExNarrow(w idExWide) idExConfig { return w.Cfg }

// ExampleIdentity demonstrates the defining property of Identity: it succeeds
// with the outer environment unchanged, without touching the runtime context.
func ExampleIdentity() {
	id := Identity[idExConfig]()

	fmt.Println(id(idExConfig{Host: "localhost", Port: 8080})(context.Background())())
	// Output:
	// Right[readerreaderioresult.idExConfig]({localhost 8080})
}

// ExampleIdentity_projection demonstrates that Asks is exactly Identity
// followed by Map, so Identity is the point every projection is built from.
func ExampleIdentity_projection() {
	viaIdentity := F.Pipe1(Identity[idExConfig](), Map[idExConfig](idExPort))
	viaAsks := Asks(idExPort)

	cfg := idExConfig{Host: "localhost", Port: 8080}
	ctx := context.Background()

	fmt.Println(viaIdentity(cfg)(ctx)(), viaAsks(cfg)(ctx)())
	// Output:
	// Right[int](8080) Right[int](8080)
}

// ExampleIdentity_local demonstrates that rewiring the identity arrow turns it
// into the accessor for the projection it was rewired with.
func ExampleIdentity_local() {
	viaLocal := Local[idExConfig](idExNarrow)(Identity[idExConfig]())
	viaAsks := Asks(idExNarrow)

	wide := idExWide{Cfg: idExConfig{Host: "localhost", Port: 8080}, User: "ada"}
	ctx := context.Background()

	fmt.Println(viaLocal(wide)(ctx)(), viaAsks(wide)(ctx)())
	// Output:
	// Right[readerreaderioresult.idExConfig]({localhost 8080}) Right[readerreaderioresult.idExConfig]({localhost 8080})
}
