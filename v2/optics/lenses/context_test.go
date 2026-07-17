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

package lenses

import (
	"context"
	"fmt"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// contextKey is a private key type to avoid collisions in tests.
type contextKey string

// TestAtContext_Get_SomeWhenPresent verifies that Get returns Some when the key
// is present in the context and the stored value has the expected type.
func TestAtContext_Get_SomeWhenPresent(t *testing.T) {
	lens := AtContext[string, contextKey]("user")
	ctx := context.WithValue(context.Background(), contextKey("user"), "alice")

	result := lens.Get(ctx)

	assert.Equal(t, O.Some("alice"), result)
}

// TestAtContext_Get_NoneWhenAbsent verifies that Get returns None when the key
// has never been set in the context.
func TestAtContext_Get_NoneWhenAbsent(t *testing.T) {
	lens := AtContext[string, contextKey]("user")
	ctx := context.Background()

	result := lens.Get(ctx)

	assert.Equal(t, O.None[string](), result)
}

// TestAtContext_Get_NoneOnTypeMismatch verifies that Get returns None when a
// value is stored under the key but has an incompatible type.
func TestAtContext_Get_NoneOnTypeMismatch(t *testing.T) {
	lens := AtContext[int, contextKey]("value")
	// store a string, but the lens expects int
	ctx := context.WithValue(context.Background(), contextKey("value"), "not-an-int")

	result := lens.Get(ctx)

	assert.Equal(t, O.None[int](), result)
}

// TestAtContext_Set_SomeAddsValue verifies that Set(Some(v)) produces a new context
// that carries the value and leaves the original context unchanged.
func TestAtContext_Set_SomeAddsValue(t *testing.T) {
	lens := AtContext[string, contextKey]("user")
	original := context.Background()

	updated := lens.Set(O.Some("bob"))(original)

	// The new context has the value.
	assert.Equal(t, "bob", updated.Value(contextKey("user")))

	// The original context is unmodified.
	assert.Nil(t, original.Value(contextKey("user")))
}

// TestAtContext_Set_NonePreservesContext verifies that Set(None) returns the original
// context without any modification.
func TestAtContext_Set_NonePreservesContext(t *testing.T) {
	lens := AtContext[string, contextKey]("user")
	original := context.WithValue(context.Background(), contextKey("user"), "alice")

	result := lens.Set(O.None[string]())(original)

	// Value is still present, context is the same object.
	assert.Equal(t, "alice", result.Value(contextKey("user")))
}

// TestAtContext_SetGet_LensLaw verifies the SetGet law: Get(Set(Some(a))(s)) == Some(a).
func TestAtContext_SetGet_LensLaw(t *testing.T) {
	lens := AtContext[string, contextKey]("token")
	ctx := context.Background()

	withToken := lens.Set(O.Some("secret"))(ctx)
	retrieved := lens.Get(withToken)

	assert.Equal(t, O.Some("secret"), retrieved)
}

// TestAtContext_GetSet_LensLaw verifies the GetSet law: Set(Get(s))(s) produces a
// context that contains the same value for the lens key.
func TestAtContext_GetSet_LensLaw(t *testing.T) {
	lens := AtContext[string, contextKey]("role")
	ctx := context.WithValue(context.Background(), contextKey("role"), "admin")

	result := lens.Set(lens.Get(ctx))(ctx)

	assert.Equal(t, O.Some("admin"), lens.Get(result))
}

// TestAtContext_SetSet_LensLaw verifies the SetSet law: the second Set wins.
func TestAtContext_SetSet_LensLaw(t *testing.T) {
	lens := AtContext[string, contextKey]("role")
	ctx := context.Background()

	result := lens.Set(O.Some("manager"))(lens.Set(O.Some("admin"))(ctx))

	assert.Equal(t, O.Some("manager"), lens.Get(result))
}

// TestAtContext_IntValue verifies that the lens works with integer values.
func TestAtContext_IntValue(t *testing.T) {
	lens := AtContext[int, contextKey]("count")
	ctx := context.WithValue(context.Background(), contextKey("count"), 42)

	assert.Equal(t, O.Some(42), lens.Get(ctx))
}

// TestAtContext_StructValue verifies that the lens works with struct values.
func TestAtContext_StructValue(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	lens := AtContext[User, contextKey]("currentUser")
	user := User{ID: 1, Name: "alice"}
	ctx := context.WithValue(context.Background(), contextKey("currentUser"), user)

	assert.Equal(t, O.Some(user), lens.Get(ctx))
}

// TestAtContext_StringKey verifies that plain string keys work alongside custom
// key types.
func TestAtContext_StringKey(t *testing.T) {
	lens := AtContext[string]("name")
	ctx := context.WithValue(context.Background(), "name", "bob") //nolint:staticcheck

	assert.Equal(t, O.Some("bob"), lens.Get(ctx))
}

// TestAtContext_Name verifies that the lens carries the expected human-readable name.
func TestAtContext_Name(t *testing.T) {
	key := contextKey("session")
	lens := AtContext[string](key)

	assert.Equal(t, fmt.Sprintf("AtContext[%v]", key), fmt.Sprintf("%v", lens))
}

// TestAtContext_MultipleKeys verifies that two lenses with different keys are
// completely independent.
func TestAtContext_MultipleKeys(t *testing.T) {
	userLens := AtContext[string, contextKey]("user")
	tokenLens := AtContext[string, contextKey]("token")

	ctx := context.Background()
	ctx = userLens.Set(O.Some("alice"))(ctx)
	ctx = tokenLens.Set(O.Some("tok-123"))(ctx)

	assert.Equal(t, O.Some("alice"), userLens.Get(ctx))
	assert.Equal(t, O.Some("tok-123"), tokenLens.Get(ctx))
}
