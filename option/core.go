// Copyright (c) 2023 IBM Corp.
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

package option

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var (
	jsonNull = []byte("null")
)

// Option defines a data structure that logically holds a value or not
type Option[A any] struct {
	isSome bool
	some   A
}

// String prints some debug info for the object
func (s Option[A]) String() string {
	if s.isSome {
		return fmt.Sprintf("Some[%T](%v)", s.some, s.some)
	}
	return fmt.Sprintf("None[%T]", s.some)
}

// Format prints some debug info for the object
func (s Option[A]) Format(f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, s.String())
	default:
		fmt.Fprint(f, s.String())
	}
}

func (s Option[A]) MarshalJSON() ([]byte, error) {
	if IsSome(s) {
		return json.Marshal(s.some)
	}
	return jsonNull, nil
}

func (s *Option[A]) UnmarshalJSON(data []byte) error {
	// decode the value
	if bytes.Equal(data, jsonNull) {
		s.isSome = false
		s.some = *new(A)
		return nil
	}
	s.isSome = true
	return json.Unmarshal(data, &s.some)
}

func IsNone[T any](val Option[T]) bool {
	return !val.isSome
}

func Some[T any](value T) Option[T] {
	return Option[T]{isSome: true, some: value}
}

func Of[T any](value T) Option[T] {
	return Some(value)
}

func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

func IsSome[T any](val Option[T]) bool {
	return val.isSome
}

func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B {
	if IsSome(ma) {
		return onSome(ma.some)
	}
	return onNone()
}

func Unwrap[A any](ma Option[A]) (A, bool) {
	return ma.some, ma.isSome
}
