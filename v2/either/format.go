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

package either

import (
	"fmt"
	"log/slog"
)

const (
	leftFmtTemplate  = "Left[%T](%v)"
	rightFmtTemplate = "Right[%T](%v)"
)

// String prints some debug info for the object
//
//go:noinline
func (s Either[E, A]) String() string {
	if !s.isLeft {
		return fmt.Sprintf(rightFmtTemplate, s.r, s.r)
	}
	return fmt.Sprintf(leftFmtTemplate, s.l, s.l)
}

// Format implements fmt.Formatter for Either.
// Supports all standard format verbs:
//   - %s, %v, %+v, %q, and all other verbs: uses String() representation
//
//go:noinline
func (s Either[E, A]) Format(f fmt.State, c rune) {
	switch c {
	case 'q':
		fmt.Fprintf(f, "%q", s.String())
	default:
		fmt.Fprint(f, s.String())
	}
}

// LogValue implements slog.LogValuer for Either.
// Returns a slog.Value that represents the Either for structured logging.
// Returns a group value with "right" key for Right values and "left" key for Left values.
//
// Example:
//
//	logger := slog.Default()
//	result := either.Right[error](42)
//	logger.Info("result", "value", result)
//	// Logs: {"msg":"result","value":{"right":42}}
//
//	err := either.Left[int](errors.New("failed"))
//	logger.Error("error", "value", err)
//	// Logs: {"msg":"error","value":{"left":"failed"}}
//
//go:noinline
func (s Either[E, A]) LogValue() slog.Value {
	if !s.isLeft {
		return slog.GroupValue(slog.Any("right", s.r))
	}
	return slog.GroupValue(slog.Any("left", s.l))
}
