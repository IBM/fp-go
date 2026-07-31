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

// Package lenses contains the pure utility functions used by the lens code
// generator. Everything here operates on plain data structures (no CLI flags,
// no filesystem side-effects) so that each function can be unit-tested in
// isolation.
package lenses

// StructInfo holds information about a struct that needs lens generation.
type StructInfo struct {
	Name           string
	QualifiedName  string // e.g., "openai.ChatCompletionUserMessageParam" when package differs
	TypeParams     string // e.g., "[T any]" or "[K comparable, V any]" - for type declarations
	TypeParamNames string // e.g., "[T]" or "[K, V]" - for type usage in function signatures
	Fields         []FieldInfo
	Imports        map[string]string // package path -> alias
}

// FieldInfo holds information about a struct field.
type FieldInfo struct {
	Name         string
	TypeName     string
	BaseType     string // TypeName without leading * for pointer types
	IsOptional   bool   // true if field is a pointer or has json omitempty tag
	IsComparable bool   // true if the type is comparable (can use ==)
	IsOption     bool   // true if the field type is already an option.Option — LensO must not be generated (would produce Option[Option[A]])
	IsEmbedded   bool   // true if this field comes from an embedded struct
	IsDeprecated bool   // true if the field is marked as deprecated
}
