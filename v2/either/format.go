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

func eitherString(tmp string, value any) string {
	return fmt.Sprintf(tmp, value, value)
}

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	if !s.l {
		return eitherString(rightFmtTemplate, s.a)
	}
	return eitherString(leftFmtTemplate, s.e)
}

func eitherFormat(f fmt.State, c rune, s string) {
	switch c {
	case 'q':
		fmt.Fprintf(f, "%q", s)
	default:
		fmt.Fprint(f, s)
	}
}

// Format implements fmt.Formatter for Either.
// Supports all standard format verbs:
//   - %s, %v, %+v, %q, and all other verbs: uses String() representation
//

func (s Either[E, A]) Format(f fmt.State, c rune) {
	eitherFormat(f, c, s.String())
}

func eitherLogValue(label string, value any) slog.Value {
	return slog.GroupValue(slog.Any(label, value))
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

func (s Either[E, A]) LogValue() slog.Value {
	if !s.l {
		return eitherLogValue("right", s.a)
	}
	return eitherLogValue("left", s.e)
}
