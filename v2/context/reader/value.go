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
	IC "github.com/IBM/fp-go/v2/internal/context"
	"github.com/IBM/fp-go/v2/option"
)

// AskValue creates a Reader that looks up a typed value stored in the context under key.
//
// The value is retrieved with [context.Context.Value] and type-asserted to V.
// The result is Some(v) if the key is present and holds a V, and None if the key
// is absent or the stored value has a different type. This avoids the manual
// "v, ok := ctx.Value(key).(V)" boilerplate and never panics.
//
// AskValue is the read counterpart of [WithValue].
//
// Type Parameters:
//   - V: The expected type of the stored value (must be given explicitly)
//   - K: The type of the key (inferred from the argument)
//
// Parameters:
//   - key: The context key to look up
//
// Returns:
//   - A Reader[Option[V]] that reads the value from the context
//
// Example:
//
//	type contextKey string
//	const userKey contextKey = "user"
//
//	getUser := reader.AskValue[User](userKey)
//
//	ctx := reader.WithValue[User](userKey)(User{Name: "Alice"})(context.Background())
//	getUser(ctx)                  // Some(User{Name: "Alice"})
//	getUser(context.Background()) // None
//
// See Also:
//   - WithValue: Stores a value in the context
//   - option.InstanceOf: The type-safe assertion used internally
func AskValue[V, K any](key K) Reader[option.Option[V]] {
	return IC.AskValue[V](key)
}
