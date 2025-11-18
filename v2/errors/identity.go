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

// Package errors provides functional utilities for working with Go errors.
// It includes functions for error creation, transformation, and type conversion
// that integrate well with functional programming patterns.
package errors

import (
	F "github.com/IBM/fp-go/v2/function"
)

// Identity is the identity function specialized for error types.
// It returns the error unchanged, useful in functional composition where
// an error needs to be passed through without modification.
//
// Example:
//
//	err := errors.New("something went wrong")
//	same := Identity(err) // returns the same error
var Identity = F.Identity[error]
