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

// Package retry provides utilities for implementing retry logic with configurable policies.
// It includes type aliases and functions for handling retryable operations with exponential backoff,
// maximum retry attempts, and other retry strategies.
package retry

import (
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Option is a type alias for option.Option, representing an optional value that may or may not be present.
	// It is used throughout the retry package to represent optional configuration values and results.
	//
	// An Option[A] can be either:
	//   - Some(value): Contains a value of type A
	//   - None: Represents the absence of a value
	//
	// This type is commonly used for:
	//   - Optional retry delays
	//   - Optional maximum retry counts
	//   - Results that may or may not succeed after retries
	Option[A any] = option.Option[A]
)
