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

// Package generic implements traversal endomorphisms for arbitrary effect types.
//
// A traversal endomorphism is a Traversal[S, A, HKTES, HKTA] whose result effect HKTES
// wraps an Endomorphism[S] instead of S. The package provides:
//
//   - Empty, Concat and MakeMonoid: the monoid that combines traversal endomorphisms
//     focusing on different parts of the same structure
//   - Compose: sequential composition for nested access
//   - ToTraversal: conversion into a regular traversal
//
// All functions take the type class operations (Of, Map, Ap) of the effect explicitly,
// so they work with any applicative such as identity, option, result or
// context/readerioresult.
package generic

//go:generate go run ../../../main.go lens --dir . --filename gen_lens.go --include-test-files
