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

package readerioresult

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// idExKey is the key type the Identity examples store values under.
type idExKey string

const idExTenantKey idExKey = "tenant"

// idExTenant reads the tenant out of the context.
func idExTenant(ctx context.Context) string {
	tenant, _ := ctx.Value(idExTenantKey).(string)
	return tenant
}

// ExampleIdentity demonstrates the defining property of Identity: it succeeds
// with the [context.Context] it was given, unchanged.
func ExampleIdentity() {
	ctx := context.WithValue(context.Background(), idExTenantKey, "acme")

	// the context flows straight through, so reading the value back out of the
	// result yields what was put in
	same := F.Pipe1(Identity(), Map(idExTenant))

	fmt.Println(same(ctx)())
	// Output:
	// Right[string](acme)
}

// ExampleIdentity_projection demonstrates that Identity is the point every
// context projection is built from: Map over it reads a single value out.
func ExampleIdentity_projection() {
	ctx := context.WithValue(context.Background(), idExTenantKey, "acme")

	upper := F.Pipe1(Identity(), Map(F.Flow2(idExTenant, func(s string) string {
		return "tenant=" + s
	})))

	fmt.Println(upper(ctx)())
	// Output:
	// Right[string](tenant=acme)
}
